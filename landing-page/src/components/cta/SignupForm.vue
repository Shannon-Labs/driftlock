<template>
  <div class="mx-auto max-w-lg p-4 sm:p-8">
    <div class="bg-white rounded-xl shadow-xl border border-gray-100 p-8 sm:p-10">
      
      <!-- Authentication View (Login/Signup) -->
      <div v-if="authState === 'initial'" class="space-y-6">
        <div class="text-center">
          <h2 class="text-2xl font-bold text-gray-900">{{ isLoginMode ? 'Welcome back' : 'Initialize Radar' }}</h2>
          <p class="mt-2 text-sm text-gray-500">
            {{ isLoginMode ? 'Sign in to access your dashboard' : "Get instant access to Driftlock's anomaly detection API" }}
          </p>
        </div>

        <!-- Social Login Buttons -->
        <div class="grid grid-cols-2 gap-3">
          <button
            @click="handleSocialAuth('google')"
            type="button"
            class="flex w-full items-center justify-center gap-2 rounded-md border border-gray-300 bg-white px-3 py-2 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-1"
            :disabled="loading"
          >
            <svg class="h-5 w-5" viewBox="0 0 24 24">
              <path
                d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                fill="#4285F4"
              />
              <path
                d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                fill="#34A853"
              />
              <path
                d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.26-.19-.58z"
                fill="#FBBC05"
              />
              <path
                d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                fill="#EA4335"
              />
            </svg>
            <span v-if="!loading">Google</span>
             <span v-else class="opacity-0">Google</span>
          </button>

          <button
            @click="handleSocialAuth('github')"
            type="button"
            class="flex w-full items-center justify-center gap-2 rounded-md border border-gray-300 bg-white px-3 py-2 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-1"
            :disabled="loading"
          >
            <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
              <path
                fill-rule="evenodd"
                d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"
                clip-rule="evenodd"
              />
            </svg>
             <span v-if="!loading">GitHub</span>
             <span v-else class="opacity-0">GitHub</span>
          </button>
        </div>

        <div class="relative">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-gray-300"></div>
          </div>
          <div class="relative flex justify-center text-sm">
            <span class="bg-white px-2 text-gray-500">Or continue with email</span>
          </div>
        </div>

        <form @submit.prevent="handleEmailAuth" class="space-y-4">
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700">Work Email</label>
            <input
              id="email"
              v-model="email"
              type="email"
              placeholder="you@company.com"
              required
              class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-blue-500 sm:text-sm"
              :disabled="loading"
            />
          </div>

          <div v-if="!isLoginMode">
            <label for="company" class="block text-sm font-medium text-gray-700">Company Name</label>
            <input
              id="company"
              v-model="company"
              type="text"
              placeholder="Your Company"
              required
              minlength="2"
              maxlength="100"
              class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-blue-500 sm:text-sm"
              :disabled="loading"
            />
          </div>

          <!-- Password field needed for login and for real signup now -->
          <div v-if="isLoginMode || showPasswordField">
             <label for="password" class="block text-sm font-medium text-gray-700">Password</label>
             <input
              id="password"
              v-model="password"
              type="password"
              placeholder="••••••••"
              required
              minlength="6"
              class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-blue-500 sm:text-sm"
              :disabled="loading"
            />
            <div v-if="isLoginMode" class="text-right mt-1">
                <button type="button" @click="handleForgotPassword" class="text-xs text-blue-600 hover:text-blue-500">Forgot password?</button>
            </div>
          </div>

          <div v-if="error" class="rounded-md bg-red-50 p-4 text-sm text-red-700 flex items-center gap-2">
             <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            {{ error }}
          </div>
          
          <div v-if="successMessage" class="rounded-md bg-green-50 p-4 text-sm text-green-700 flex items-center gap-2">
             <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            {{ successMessage }}
          </div>

          <button
            type="submit"
            class="flex w-full justify-center rounded-md border border-transparent bg-blue-600 py-3 px-4 text-sm font-medium text-white shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-70 disabled:cursor-not-allowed transition-all"
            :disabled="loading"
          >
            <span v-if="!loading">{{ isLoginMode ? 'Sign In' : 'Start Pilot →' }}</span>
            <span v-else class="flex items-center gap-2">
              <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ isLoginMode ? 'Signing in...' : 'Creating account...' }}
            </span>
          </button>
        </form>

        <div class="text-center text-sm">
            <span class="text-gray-500">{{ isLoginMode ? "Don't have an account?" : "Already have an account?" }}</span>
            <button @click="toggleMode" class="ml-1 font-medium text-blue-600 hover:text-blue-500">
                {{ isLoginMode ? 'Sign up' : 'Sign in' }}
            </button>
        </div>
      </div>

      <!-- Pending Verification View -->
      <div v-else-if="authState === 'pending_verification'" class="text-center space-y-6">
        <div class="flex flex-col items-center">
          <div class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-blue-100">
            <svg class="h-8 w-8 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
          </div>
          <h3 class="mt-4 text-2xl font-bold text-gray-900">Check your email</h3>
          <p class="mt-2 text-sm text-gray-600 max-w-xs">
            We've sent a verification link to <strong class="text-gray-900">{{ pendingEmail }}</strong>
          </p>
        </div>

        <div class="rounded-lg bg-amber-50 border border-amber-200 p-4 text-left">
          <div class="flex">
            <svg class="h-5 w-5 text-amber-400 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <div class="ml-3">
              <h4 class="text-sm font-medium text-amber-800">Important</h4>
              <p class="mt-1 text-sm text-amber-700">
                Click the link in the email to verify your account and receive your API key. The link expires in 24 hours.
              </p>
            </div>
          </div>
        </div>

        <div class="text-sm text-gray-500">
          <p>Didn't receive the email?</p>
          <ul class="mt-2 space-y-1 text-left list-disc list-inside">
            <li>Check your spam folder</li>
            <li>Make sure you entered the correct email</li>
          </ul>
        </div>

        <button
          @click="authState = 'initial'"
          class="inline-flex items-center text-sm font-medium text-blue-600 hover:text-blue-500"
        >
          <svg class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          Back to sign up
        </button>
      </div>

      <!-- Success View (API Key Received) -->
      <div v-else-if="authState === 'submitted'" class="text-center space-y-8">
        <div class="flex flex-col items-center">
           <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-green-100">
              <svg class="h-6 w-6 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <h3 class="mt-3 text-2xl font-bold text-gray-900">Welcome back!</h3>
            <p class="mt-2 text-sm text-gray-500">Here is your API key:</p>
        </div>
        
        <div class="relative rounded-md bg-gray-50 p-4 border border-gray-200 flex items-center justify-between">
          <code class="font-mono text-sm text-blue-600 break-all">{{ apiKey }}</code>
          <button @click="copyToClipboard" class="ml-4 inline-flex items-center rounded-md border border-gray-300 bg-white px-3 py-2 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2" type="button">
            <svg v-if="!copied" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            {{ copied ? 'Copied!' : 'Copy' }}
          </button>
        </div>

        <div class="rounded-lg bg-blue-50 p-4 text-left">
          <h4 class="text-sm font-semibold text-blue-900">Next Steps:</h4>
          <ol class="mt-2 list-decimal list-inside text-sm text-blue-800 space-y-1">
            <li>Save your API key securely</li>
            <li>Read the <a href="/docs" target="_blank" class="underline font-medium">API documentation</a></li>
            <li>Make your first API call to <code>/v1/detect</code></li>
          </ol>
        </div>

        <div class="text-left">
          <h4 class="text-xs font-semibold uppercase tracking-wider text-gray-500 mb-2">Quick Start:</h4>
          <div class="bg-gray-900 rounded-lg p-4 overflow-x-auto">
              <pre class="text-xs text-gray-300 font-mono"><code>curl -X POST https://driftlock.net/api/v1/detect \
  -H "X-Api-Key: {{ apiKey }}" \
  -H "Content-Type: application/json" \
  -d '{"events": [...], "window_size": 50}'</code></pre>
          </div>
        </div>

        <div class="bg-gradient-to-r from-indigo-50 to-blue-50 p-4 rounded-lg border border-blue-100">
          <div class="flex flex-col space-y-3">
             <div>
                 <h4 class="text-sm font-bold text-blue-900">Upgrade Plan</h4>
                 <p class="text-xs text-blue-700 mt-1">Choose a plan to remove limits.</p>
             </div>
             <div class="flex gap-2">
                 <button @click="handleUpgrade('radar')" class="flex-1 inline-flex justify-center items-center rounded-md border border-blue-200 bg-white px-3 py-2 text-xs font-medium text-blue-700 shadow-sm hover:bg-blue-50 focus:outline-none" :disabled="upgrading">
                    <span v-if="!upgrading">Radar ($15)</span>
                    <span v-else>...</span>
                 </button>
                 <button @click="handleUpgrade('tensor')" class="flex-1 inline-flex justify-center items-center rounded-md border border-transparent bg-blue-600 px-3 py-2 text-xs font-medium text-white shadow-sm hover:bg-blue-700 focus:outline-none" :disabled="upgrading">
                    <span v-if="!upgrading">Pro ($100)</span>
                    <span v-else>...</span>
                 </button>
             </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { getFirebaseAuth } from '@/firebase'
import { 
  createUserWithEmailAndPassword, 
  GoogleAuthProvider, 
  GithubAuthProvider, 
  signInWithPopup, 
  signInWithEmailAndPassword,
  sendPasswordResetEmail 
} from 'firebase/auth'

const email = ref('')
const company = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')
const successMessage = ref('')
const authState = ref<'initial' | 'submitted' | 'pending_verification'>('initial')
const apiKey = ref('')
const copied = ref(false)
const upgrading = ref(false)
const isLoginMode = ref(false)
const showPasswordField = ref(true)
const pendingEmail = ref('') // Email to show in verification pending message
let currentUser: any = null

const toggleMode = () => {
    isLoginMode.value = !isLoginMode.value
    error.value = ''
    successMessage.value = ''
    password.value = ''
}

const handleForgotPassword = async () => {
    if (!email.value) {
        error.value = 'Please enter your email address first.'
        return
    }
    
    loading.value = true
    error.value = ''
    successMessage.value = ''
    
    try {
        const auth = await getFirebaseAuth()
        if (!auth) throw new Error('Auth not ready')
        
        await sendPasswordResetEmail(auth, email.value)
        successMessage.value = 'Password reset email sent! Check your inbox.'
    } catch (err: any) {
        error.value = err.message || 'Failed to send reset email.'
    } finally {
        loading.value = false
    }
}

// Handle Email/Password Auth (Sign Up or Sign In)
const handleEmailAuth = async () => {
  error.value = ''
  successMessage.value = ''
  loading.value = true

  try {
    const auth = await getFirebaseAuth()
    if (!auth) throw new Error('Auth not ready')

    if (isLoginMode.value) {
        // SIGN IN FLOW
        const userCredential = await signInWithEmailAndPassword(auth, email.value, password.value)
        currentUser = userCredential.user
        await fetchApiKeys(currentUser)
    } else {
        // SIGN UP FLOW
        const userCredential = await createUserWithEmailAndPassword(auth, email.value, password.value)
        const idToken = await userCredential.user.getIdToken()
        await submitOnboarding(idToken, email.value, company.value, 'email_form', userCredential.user)
    }

  } catch (err: any) {
    handleAuthError(err)
  } finally {
    loading.value = false
  }
}

// Handle Social Signup/Login
const handleSocialAuth = async (providerName: 'google' | 'github') => {
  error.value = ''
  successMessage.value = ''
  loading.value = true

  try {
    const auth = await getFirebaseAuth()
    if (!auth) throw new Error('Auth not ready')

    let provider
    if (providerName === 'google') {
      provider = new GoogleAuthProvider()
    } else {
      provider = new GithubAuthProvider()
    }

    const result = await signInWithPopup(auth, provider)
    currentUser = result.user
    
    // Attempt to fetch keys first (Login scenario)
    try {
        await fetchApiKeys(currentUser, true) 
    } catch (e) {
        // If no keys found, likely new user. Auto-register with derived name.
        // Derive company name from user profile or email
        let derivedCompany = currentUser.displayName 
          ? `${currentUser.displayName}'s Org` 
          : (currentUser.email ? `${currentUser.email.split('@')[0]}'s Org` : 'My Organization')
        
        // Remove special chars that might look weird
        derivedCompany = derivedCompany.replace(/[^\w\s'-]/g, '').trim()
        if (derivedCompany.length < 2) derivedCompany = "My Organization"

        const idToken = await currentUser.getIdToken()
        await submitOnboarding(idToken, currentUser.email || '', derivedCompany, 'social_auth', currentUser)
    }
    
  } catch (err: any) {
    handleAuthError(err)
    loading.value = false
  }
}

// Fetch API Keys for existing user
const fetchApiKeys = async (user: any, suppressError = false) => {
    try {
        const idToken = await user.getIdToken()
        const response = await fetch('/api/v1/me/keys', {
            headers: {
                'Authorization': `Bearer ${idToken}`
            }
        })
        
        if (response.ok) {
            const data = await response.json()
            if (data.keys && data.keys.length > 0) {
                apiKey.value = data.keys[0].key || data.keys[0].token || 'Error retrieving key' 
                authState.value = 'submitted'
                loading.value = false
                return
            }
        }
        
        if (!suppressError) {
             throw new Error('No account found or setup incomplete. Please sign up.')
        } else {
            throw new Error('No keys') 
        }
    } catch (err) {
        if (!suppressError) throw err
        throw err
    }
}

// Shared logic to call backend for signup
const submitOnboarding = async (idToken: string, userEmail: string, companyName: string, source: string, userObj: any) => {
    const response = await fetch('/api/v1/onboard/signup', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${idToken}`
      },
      body: JSON.stringify({
        email: userEmail,
        company_name: companyName,
        plan: 'trial',
        source: source
      }),
    })

    // Check content type to avoid parsing HTML as JSON
    const contentType = response.headers.get("content-type");
    if (!contentType || !contentType.includes("application/json")) {
        const text = await response.text().catch(() => 'Unknown error');
        console.error('Non-JSON response:', text);
        throw new Error(`Server returned a non-JSON response (${response.status}). Please try again later.`);
    }

    const data = await response.json()

    if (!response.ok) {
      let errorMessage = 'Failed to create account. Please try again.';
      if (data.error) {
          if (typeof data.error === 'string') {
              errorMessage = data.error;
          } else if (typeof data.error === 'object' && data.error.message) {
              errorMessage = data.error.message;
          } else {
              errorMessage = JSON.stringify(data.error);
          }
      }
      throw new Error(errorMessage)
    }

    if (data.success) {
      if (data.pending_verification) {
        // New flow: email verification required
        pendingEmail.value = userEmail
        authState.value = 'pending_verification'
      } else if (data.api_key) {
        // Legacy flow: immediate API key (shouldn't happen anymore)
        apiKey.value = data.api_key
        authState.value = 'submitted'
      }

      if (typeof window !== 'undefined' && (window as any).gtag) {
        (window as any).gtag('event', 'signup', { method: source })
      }
    } else {
      throw new Error('Invalid response from server')
    }
}

const handleAuthError = (err: any) => {
  console.error('Auth error details:', err);
  
  if (err.code === 'auth/email-already-in-use') {
      error.value = 'This email is already registered. Please sign in.'
    } else if (err.code === 'auth/invalid-email') {
      error.value = 'Please enter a valid email address.'
    } else if (err.code === 'auth/user-not-found' || err.code === 'auth/wrong-password' || err.code === 'auth/invalid-credential') {
      error.value = 'Invalid email or password.'
    } else if (err.code === 'auth/popup-closed-by-user') {
      error.value = 'Sign in cancelled.'
    } else {
      // Handle cases where err.message might be an object or missing
      let msg = err.message;
      if (typeof msg === 'object') {
          msg = JSON.stringify(msg);
      }
      if (!msg || msg === '[object Object]') {
          msg = 'Something went wrong. Please try again.';
      }
      error.value = msg;
    }
}

const copyToClipboard = async () => {
  try {
    await navigator.clipboard.writeText(apiKey.value)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}

const handleUpgrade = async (plan: string) => {
  upgrading.value = true
  try {
    const response = await fetch('/api/v1/billing/checkout', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Api-Key': apiKey.value
      },
      body: JSON.stringify({ plan })
    })
    
    if (!response.ok) throw new Error('Failed to start checkout')
    
    const data = await response.json()
    if (data.url) window.location.href = data.url
  } catch (err) {
    console.error('Upgrade error:', err)
    alert('Failed to start upgrade process.')
  } finally {
    upgrading.value = false
  }
}
</script>