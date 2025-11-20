import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { 
  signOut as firebaseSignOut, 
  onAuthStateChanged, 
  type User,
  sendSignInLinkToEmail,
  isSignInWithEmailLink,
  signInWithEmailLink,
} from 'firebase/auth'
import { getFirebaseAuth } from '../firebase'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const loading = ref(true)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => !!user.value)

  // Initialize auth listener
  const init = async () => {
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
  }

  // Magic Link Login (Passwordless)
  const sendMagicLink = async (email: string) => {
    const auth = await getFirebaseAuth();
    if (!auth) {
      error.value = 'Firebase auth not initialized. Please check your configuration.'
      throw new Error(error.value)
    }

    loading.value = true
    error.value = null
    try {
      const actionCodeSettings = {
        // URL you want to redirect back to. The domain (www.example.com) for this
        // URL must be in the authorized domains list in the Firebase Console.
        url: window.location.origin + '/login/finish',
        handleCodeInApp: true,
      }
      await sendSignInLinkToEmail(auth, email, actionCodeSettings)
      window.localStorage.setItem('emailForSignIn', email)
    } catch (e: any) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  const completeMagicLinkLogin = async () => {
    const auth = await getFirebaseAuth();
    if (!auth) {
      error.value = 'Firebase auth not initialized. Please check your configuration.'
      return false
    }

    if (!isSignInWithEmailLink(auth, window.location.href)) {
      return false
    }

    let email = window.localStorage.getItem('emailForSignIn')
    if (!email) {
      email = window.prompt('Please provide your email for confirmation')
    }

    if (!email) return false

    loading.value = true
    try {
      const result = await signInWithEmailLink(auth, email, window.location.href)
      window.localStorage.removeItem('emailForSignIn')
      user.value = result.user
      return true
    } catch (e: any) {
      error.value = e.message
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

  return {
    user,
    loading,
    error,
    isAuthenticated,
    init,
    sendMagicLink,
    completeMagicLinkLogin,
    logout,
    getToken
  }
})





