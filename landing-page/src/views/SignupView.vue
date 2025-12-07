<template>
  <div class="flex min-h-screen flex-col justify-center py-12 sm:px-6 lg:px-8 bg-white">
    <div class="sm:mx-auto sm:w-full sm:max-w-md">
      <router-link to="/" class="flex justify-center">
        <div class="h-16 w-16 hover:scale-105 transition-transform duration-300">
          <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M6 11V4H18V11" stroke="black" stroke-width="4"/>
            <rect x="2" y="11" width="20" height="11" fill="black"/>
            <path d="M4 16.5H7L9.5 12.5L12.5 19.5L15 15.5L17 16.5H20" stroke="white" stroke-width="2" stroke-linejoin="bevel"/>
          </svg>
        </div>
      </router-link>
      <h2 class="mt-6 text-center text-3xl font-sans font-black uppercase tracking-tighter text-black">
        Create your account
      </h2>
      <p class="mt-2 text-center text-sm text-gray-600 font-mono">
        Already have an account?
        <router-link to="/login" class="font-bold text-black underline decoration-2 underline-offset-4 hover:bg-black hover:text-white transition-colors px-1">Sign in</router-link>
      </p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
      <div class="bg-white py-8 px-4 border-2 border-black sm:px-10 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]">
        <form class="space-y-6" @submit.prevent="handleSignup">
          <div>
            <label for="email" class="block text-sm font-bold uppercase tracking-wide text-black">Email address</label>
            <div class="mt-1">
              <input
                id="email"
                v-model="email"
                name="email"
                type="email"
                autocomplete="email"
                required
                class="block w-full appearance-none border-2 border-black px-3 py-3 placeholder-gray-500 shadow-none focus:outline-none focus:ring-2 focus:ring-black focus:ring-offset-2 sm:text-sm font-mono"
                placeholder="you@example.com"
              />
            </div>
          </div>

          <div>
            <label for="password" class="block text-sm font-bold uppercase tracking-wide text-black">Password</label>
            <div class="mt-1">
              <input
                id="password"
                v-model="password"
                name="password"
                type="password"
                autocomplete="new-password"
                required
                minlength="6"
                class="block w-full appearance-none border-2 border-black px-3 py-3 placeholder-gray-500 shadow-none focus:outline-none focus:ring-2 focus:ring-black focus:ring-offset-2 sm:text-sm font-mono"
                placeholder="••••••••"
              />
            </div>
            <p class="mt-1 text-xs text-gray-500 font-mono">At least 6 characters</p>
          </div>

          <div>
            <label for="confirmPassword" class="block text-sm font-bold uppercase tracking-wide text-black">Confirm password</label>
            <div class="mt-1">
              <input
                id="confirmPassword"
                v-model="confirmPassword"
                name="confirmPassword"
                type="password"
                autocomplete="new-password"
                required
                class="block w-full appearance-none border-2 border-black px-3 py-3 placeholder-gray-500 shadow-none focus:outline-none focus:ring-2 focus:ring-black focus:ring-offset-2 sm:text-sm font-mono"
                placeholder="••••••••"
              />
            </div>
          </div>

          <div v-if="localError" class="text-red-600 text-sm font-bold border border-red-600 p-2 bg-red-50">
            {{ localError }}
          </div>

          <div v-if="authStore.error" class="text-red-600 text-sm font-bold border border-red-600 p-2 bg-red-50">
            {{ authStore.error }}
          </div>

          <div>
            <button
              type="submit"
              :disabled="authStore.loading"
              class="flex w-full justify-center border-2 border-black bg-black py-3 px-4 text-sm font-bold uppercase tracking-widest text-white hover:bg-white hover:text-black transition-colors focus:outline-none focus:ring-2 focus:ring-black focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed shadow-[4px_4px_0px_0px_rgba(0,0,0,0)] hover:shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none"
            >
              <span v-if="authStore.loading" class="flex items-center gap-2">
                <svg class="animate-spin h-4 w-4 text-current" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Creating account...
              </span>
              <span v-else>Create Account</span>
            </button>
          </div>

          <!-- Divider -->
          <div class="relative">
            <div class="absolute inset-0 flex items-center">
              <div class="w-full border-t-2 border-black"></div>
            </div>
            <div class="relative flex justify-center text-sm">
              <span class="bg-white px-4 font-bold uppercase tracking-wide text-black">Or</span>
            </div>
          </div>

          <!-- Google Sign Up -->
          <div>
            <button
              type="button"
              @click="handleGoogleSignIn"
              :disabled="authStore.loading"
              class="flex w-full items-center justify-center gap-3 border-2 border-black bg-white py-3 px-4 text-sm font-bold uppercase tracking-widest text-black hover:bg-black hover:text-white transition-colors focus:outline-none focus:ring-2 focus:ring-black focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:shadow-[4px_4px_0px_0px_rgba(0,0,0,0)] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none"
            >
              <svg class="h-5 w-5" viewBox="0 0 24 24">
                <path fill="currentColor" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                <path fill="currentColor" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                <path fill="currentColor" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                <path fill="currentColor" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
              </svg>
              Continue with Google
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const localError = ref<string | null>(null)

const handleSignup = async () => {
  localError.value = null
  authStore.clearError()

  if (!email.value || !password.value || !confirmPassword.value) return

  if (password.value !== confirmPassword.value) {
    localError.value = 'Passwords do not match.'
    return
  }

  if (password.value.length < 6) {
    localError.value = 'Password must be at least 6 characters.'
    return
  }

  try {
    await authStore.signUpWithEmail(email.value, password.value)
    router.push('/dashboard')
  } catch (e) {
    // Error handled in store/UI
  }
}

const handleGoogleSignIn = async () => {
  try {
    await authStore.signInWithGoogle()
    router.push('/dashboard')
  } catch (e) {
    // Error handled in store/UI (except popup closed)
  }
}
</script>
