<template>
  <div class="w-full max-w-md">
    <!-- Success State -->
    <div v-if="submitted" class="text-center py-8">
      <div class="w-16 h-16 mx-auto border-2 border-current flex items-center justify-center mb-4">
        <svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
      </div>
      <p class="text-xl font-bold font-sans uppercase mb-2">You're on the list</p>
      <p class="text-sm font-serif text-gray-600">We'll notify you when we launch.</p>
    </div>

    <!-- Form -->
    <form v-else @submit.prevent="handleSubmit" class="space-y-4">
      <div>
        <input
          v-model="email"
          type="email"
          placeholder="you@company.com"
          required
          class="w-full px-4 py-3 border-2 border-black bg-white text-black
                 placeholder-gray-500 font-mono text-sm
                 focus:outline-none focus:ring-0 focus:border-black
                 disabled:opacity-50 disabled:cursor-not-allowed"
          :disabled="loading"
        />
      </div>

      <div v-if="error" class="border-2 border-red-600 bg-red-50 p-3 text-sm text-red-700 font-mono">
        {{ error }}
      </div>

      <button
        type="submit"
        class="brutalist-button-primary w-full py-3 text-sm uppercase tracking-wider
               disabled:opacity-50 disabled:cursor-not-allowed"
        :disabled="loading || !email"
      >
        <span v-if="!loading">Notify Me</span>
        <span v-else class="flex items-center justify-center gap-2">
          <svg class="animate-spin h-4 w-4" viewBox="0 0 24 24" fill="none">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
          Joining...
        </span>
      </button>

      <p class="text-xs text-gray-500 font-mono text-center">
        No spam. Unsubscribe anytime.
      </p>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const email = ref('')
const loading = ref(false)
const error = ref('')
const submitted = ref(false)

const handleSubmit = async () => {
  if (!email.value || loading.value) return

  error.value = ''
  loading.value = true

  try {
    const res = await fetch('/api/v1/waitlist', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        email: email.value,
        source: 'website'
      })
    })

    const data = await res.json()

    if (res.ok && data.success) {
      submitted.value = true
      // Track conversion
      if (typeof window !== 'undefined' && (window as any).gtag) {
        (window as any).gtag('event', 'waitlist_signup', { method: 'email' })
      }
    } else {
      error.value = data.error?.message || data.message || 'Something went wrong. Please try again.'
    }
  } catch (e) {
    error.value = 'Network error. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>
