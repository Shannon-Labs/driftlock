<template>
  <div class="flex min-h-screen flex-col justify-center py-12 sm:px-6 lg:px-8 bg-white">
    <div class="sm:mx-auto sm:w-full sm:max-w-md">
      <router-link to="/" class="flex justify-center">
        <div class="h-16 w-16 bg-black flex items-center justify-center border-2 border-black hover:rotate-12 transition-transform duration-300">
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M2 12H22M12 2V22M4.5 4.5L19.5 19.5M4.5 19.5L19.5 4.5" stroke="white" stroke-width="3"/>
          </svg>
        </div>
      </router-link>
      <h2 class="mt-6 text-center text-3xl font-sans font-black uppercase tracking-tighter text-black">
        Sign in to dashboard
      </h2>
      <p class="mt-2 text-center text-sm text-gray-600 font-mono">
        Or
        <router-link to="/#signup" class="font-bold text-black underline decoration-2 underline-offset-4 hover:bg-black hover:text-white transition-colors px-1">initialize your pilot</router-link>
      </p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
      <div class="bg-white py-8 px-4 border-2 border-black sm:px-10 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]">
        <div v-if="emailSent" class="border-2 border-black bg-gray-100 p-4 mb-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-black" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-bold uppercase tracking-wide text-black">Check your email</h3>
              <div class="mt-2 text-sm text-gray-800 font-serif">
                <p>We sent a magic link to <strong>{{ email }}</strong>. Click it to log in.</p>
              </div>
              <div class="mt-4">
                <button 
                  @click="emailSent = false" 
                  class="text-sm font-bold text-black underline decoration-2 underline-offset-4 hover:bg-black hover:text-white transition-colors px-1 focus:outline-none"
                >
                  Didn't receive it? Try again.
                </button>
              </div>
            </div>
          </div>
        </div>

        <form v-else class="space-y-6" @submit.prevent="handleLogin">
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
                Sending...
              </span>
              <span v-else>Send Magic Link</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const email = ref('')
const emailSent = ref(false)

const handleLogin = async () => {
  if (!email.value) return
  
  try {
    await authStore.sendMagicLink(email.value)
    emailSent.value = true
  } catch (e) {
    // Error handled in store/UI
  }
}
</script>
