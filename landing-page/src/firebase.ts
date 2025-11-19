import { initializeApp } from 'firebase/app'
import { getAuth, type Auth } from 'firebase/auth'

// Your web app's Firebase configuration
// For Firebase JS SDK v7.20.0 and later, measurementId is optional
// NOTE: These should be environment variables in production
const firebaseConfig = {
  apiKey: import.meta.env.VITE_FIREBASE_API_KEY || '',
  authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN || '',
  projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID || '',
  storageBucket: import.meta.env.VITE_FIREBASE_STORAGE_BUCKET || '',
  messagingSenderId: import.meta.env.VITE_FIREBASE_MESSAGING_SENDER_ID || '',
  appId: import.meta.env.VITE_FIREBASE_APP_ID || '',
  measurementId: import.meta.env.VITE_FIREBASE_MEASUREMENT_ID || ''
}

// Validate Firebase configuration
const isFirebaseConfigValid = firebaseConfig.apiKey && 
                              firebaseConfig.authDomain && 
                              firebaseConfig.projectId && 
                              firebaseConfig.appId

if (!isFirebaseConfigValid) {
  console.warn('Firebase configuration is incomplete. Please set all VITE_FIREBASE_* environment variables.')
}

// Initialize Firebase only if config is valid
let app
let auth: Auth | null = null

if (isFirebaseConfigValid) {
  app = initializeApp(firebaseConfig)
  auth = getAuth(app)
}

export { auth }


