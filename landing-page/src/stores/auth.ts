import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  signOut as firebaseSignOut,
  onAuthStateChanged,
  type User,
  createUserWithEmailAndPassword,
  signInWithEmailAndPassword,
  signInWithPopup,
  GoogleAuthProvider,
  sendPasswordResetEmail,
} from 'firebase/auth'
import { getFirebaseAuth } from '../firebase'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const loading = ref(true)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => !!user.value)

  // Initialize auth listener
  const init = async () => {
    try {
      const auth = await getFirebaseAuth();
      return new Promise<void>((resolve) => {
        if (!auth) {
          console.warn('Firebase auth not initialized. Please check your Firebase configuration.')
          loading.value = false
          resolve()
          return
        }

        onAuthStateChanged(auth, (u) => {
          user.value = u
          loading.value = false
          resolve()
        })
      })
    } catch (e) {
      console.error('Failed to initialize auth:', e)
      loading.value = false
      return Promise.resolve()
    }
  }

  // Email/Password Sign Up
  const signUpWithEmail = async (email: string, password: string) => {
    const auth = await getFirebaseAuth();
    if (!auth) {
      error.value = 'Firebase auth not initialized. Please check your configuration.'
      throw new Error(error.value)
    }

    loading.value = true
    error.value = null
    try {
      const result = await createUserWithEmailAndPassword(auth, email, password)
      user.value = result.user
      return result.user
    } catch (e: any) {
      error.value = getErrorMessage(e.code)
      throw e
    } finally {
      loading.value = false
    }
  }

  // Email/Password Sign In
  const signInWithEmail = async (email: string, password: string) => {
    const auth = await getFirebaseAuth();
    if (!auth) {
      error.value = 'Firebase auth not initialized. Please check your configuration.'
      throw new Error(error.value)
    }

    loading.value = true
    error.value = null
    try {
      const result = await signInWithEmailAndPassword(auth, email, password)
      user.value = result.user
      return result.user
    } catch (e: any) {
      error.value = getErrorMessage(e.code)
      throw e
    } finally {
      loading.value = false
    }
  }

  // Google Sign In
  const signInWithGoogle = async () => {
    const auth = await getFirebaseAuth();
    if (!auth) {
      error.value = 'Firebase auth not initialized. Please check your configuration.'
      throw new Error(error.value)
    }

    loading.value = true
    error.value = null
    try {
      const provider = new GoogleAuthProvider()
      const result = await signInWithPopup(auth, provider)
      user.value = result.user
      return result.user
    } catch (e: any) {
      // Don't show error if user closed the popup
      if (e.code !== 'auth/popup-closed-by-user' && e.code !== 'auth/cancelled-popup-request') {
        error.value = getErrorMessage(e.code)
      }
      throw e
    } finally {
      loading.value = false
    }
  }

  // Password Reset
  const resetPassword = async (email: string) => {
    const auth = await getFirebaseAuth();
    if (!auth) {
      error.value = 'Firebase auth not initialized. Please check your configuration.'
      throw new Error(error.value)
    }

    loading.value = true
    error.value = null
    try {
      await sendPasswordResetEmail(auth, email)
    } catch (e: any) {
      error.value = getErrorMessage(e.code)
      throw e
    } finally {
      loading.value = false
    }
  }

  const logout = async () => {
    const auth = await getFirebaseAuth();
    if (!auth) {
      user.value = null
      return
    }

    await firebaseSignOut(auth)
    user.value = null
  }

  const getToken = async () => {
    if (!user.value) return null
    return await user.value.getIdToken()
  }

  const clearError = () => {
    error.value = null
  }

  // Convert Firebase error codes to user-friendly messages
  const getErrorMessage = (code: string): string => {
    switch (code) {
      case 'auth/email-already-in-use':
        return 'This email is already registered. Try signing in instead.'
      case 'auth/invalid-email':
        return 'Please enter a valid email address.'
      case 'auth/operation-not-allowed':
        return 'This sign-in method is not enabled.'
      case 'auth/weak-password':
        return 'Password should be at least 6 characters.'
      case 'auth/user-disabled':
        return 'This account has been disabled.'
      case 'auth/user-not-found':
      case 'auth/wrong-password':
      case 'auth/invalid-credential':
        return 'Invalid email or password.'
      case 'auth/too-many-requests':
        return 'Too many attempts. Please try again later.'
      case 'auth/network-request-failed':
        return 'Network error. Please check your connection.'
      default:
        return 'An error occurred. Please try again.'
    }
  }

  return {
    user,
    loading,
    error,
    isAuthenticated,
    init,
    signUpWithEmail,
    signInWithEmail,
    signInWithGoogle,
    resetPassword,
    logout,
    getToken,
    clearError,
  }
})
