import { initializeApp, type FirebaseApp } from 'firebase/app'
import { getAuth, type Auth } from 'firebase/auth'

let app: FirebaseApp | null = null;
let auth: Auth | null = null;
let firebasePromise: Promise<void> | null = null;

const initializeFirebase = async () => {
  if (app) return;

  const buildConfigFromEnv = () => {
    const envConfig = {
      apiKey: import.meta.env.VITE_FIREBASE_API_KEY,
      authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
      projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID,
      storageBucket: import.meta.env.VITE_FIREBASE_STORAGE_BUCKET,
      messagingSenderId: import.meta.env.VITE_FIREBASE_MESSAGING_SENDER_ID,
      appId: import.meta.env.VITE_FIREBASE_APP_ID,
      measurementId: import.meta.env.VITE_FIREBASE_MEASUREMENT_ID,
    }

    const hasRequired = envConfig.apiKey && envConfig.authDomain && envConfig.projectId && envConfig.appId
    return hasRequired ? envConfig : null
  }

  let firebaseConfig: any = null

  try {
    // Fetch Firebase config from Firebase Function
    // In prod, this is served via Firebase Hosting rewrite → getFirebaseConfig function
    // In dev, you can run Firebase emulators or rely on env fallback
    const functionUrl = '/getFirebaseConfig'

    const response = await fetch(functionUrl)
    if (!response.ok) {
      throw new Error(`Failed to fetch Firebase config: ${response.statusText}`)
    }
    firebaseConfig = await response.json()

    const isValid = firebaseConfig?.apiKey && firebaseConfig?.authDomain && firebaseConfig?.projectId && firebaseConfig?.appId
    if (!isValid) {
      console.warn('Fetched Firebase configuration is incomplete; falling back to env config.')
      firebaseConfig = buildConfigFromEnv()
    }
  } catch (error) {
    console.error('Error initializing Firebase via function:', error)
    firebaseConfig = buildConfigFromEnv()
  }

  if (!firebaseConfig) {
    throw new Error('Firebase configuration unavailable (function + env fallback failed).')
  }

  app = initializeApp(firebaseConfig)
  auth = getAuth(app)
};

const getFirebaseAuth = async () => {
  if (!firebasePromise) {
    firebasePromise = initializeFirebase();
  }
  await firebasePromise;
  return auth;
};

export { getFirebaseAuth };

