import { initializeApp, type FirebaseApp } from 'firebase/app'
import { getAuth, type Auth } from 'firebase/auth'

let app: FirebaseApp | null = null;
let auth: Auth | null = null;
let firebasePromise: Promise<void> | null = null;

const initializeFirebase = async () => {
  if (app) {
    return;
  }

  try {
    // In production, the function is at the root. In dev, it's proxied by vite.
    const functionUrl = import.meta.env.PROD
      ? '/getFirebaseConfig' // Relative path for prod
      : 'http://127.0.0.1:5001/driftlock/us-central1/getFirebaseConfig'; // Local emulator

    const response = await fetch(functionUrl);
    if (!response.ok) {
      throw new Error(`Failed to fetch Firebase config: ${response.statusText}`);
    }
    const firebaseConfig = await response.json();

    // Validate fetched Firebase configuration
    const isFirebaseConfigValid = firebaseConfig.apiKey &&
                                  firebaseConfig.authDomain &&
                                  firebaseConfig.projectId &&
                                  firebaseConfig.appId;

    if (!isFirebaseConfigValid) {
      console.warn('Fetched Firebase configuration is incomplete.');
      return;
    }

    app = initializeApp(firebaseConfig);
    auth = getAuth(app);
  } catch (error) {
    console.error("Error initializing Firebase:", error);
  }
};

const getFirebaseAuth = async () => {
  if (!firebasePromise) {
    firebasePromise = initializeFirebase();
  }
  await firebasePromise;
  return auth;
};

export { getFirebaseAuth };


