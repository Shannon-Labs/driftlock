
import { initializeApp } from 'firebase/app';
import { getAuth, connectAuthEmulator, createUserWithEmailAndPassword } from 'firebase/auth';
import fetch from 'node-fetch';

const firebaseConfig = {
  apiKey: "fake-api-key",
  authDomain: "localhost",
  projectId: "driftlock",
  storageBucket: "driftlock.appspot.com",
  messagingSenderId: "123456789",
  appId: "1:123456789:web:a1b2c3d4e5f6g7h8"
};

const app = initializeApp(firebaseConfig);
const auth = getAuth(app);
connectAuthEmulator(auth, "http://localhost:9099");

const email = `test-${Date.now()}@example.com`;
const password = "password123";

async function run() {
  console.log(`Creating user: ${email}`);
  try {
    const userCredential = await createUserWithEmailAndPassword(auth, email, password);
    const user = userCredential.user;
    const token = await user.getIdToken();
    console.log("Got ID token");

    const response = await fetch('http://localhost:8080/v1/auth/signup', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({
        email: email,
        company_name: "Headless Test Co",
        plan: "trial",
        source: "script"
      })
    });

    const data = await response.json();
    console.log("Signup Response:", JSON.stringify(data, null, 2));

    if (data.success && data.pending_verification) {
        console.log("Signup successful! Verification required.");
        // We could look up the token in Postgres here if we wanted to be 100% sure
    } else {
        console.log("Signup unexpected response");
    }

  } catch (error) {
    console.error("Error:", error);
    process.exit(1);
  }
}

run();
