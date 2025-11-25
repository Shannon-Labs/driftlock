<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50 flex items-center justify-center p-4">
    <div class="mx-auto max-w-lg w-full">
      <div class="bg-white rounded-xl shadow-xl border border-gray-100 p-8 sm:p-10">

        <!-- Loading State -->
        <div v-if="loading" class="text-center space-y-6">
          <div class="flex justify-center">
            <svg class="animate-spin h-12 w-12 text-blue-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </div>
          <p class="text-gray-600">Verifying your email...</p>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="text-center space-y-6">
          <div class="flex flex-col items-center">
            <div class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-red-100">
              <svg class="h-8 w-8 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </div>
            <h3 class="mt-4 text-2xl font-bold text-gray-900">Verification Failed</h3>
            <p class="mt-2 text-sm text-gray-600">{{ error }}</p>
          </div>

          <div class="space-y-3">
            <router-link
              to="/login"
              class="block w-full text-center rounded-md border border-transparent bg-blue-600 py-3 px-4 text-sm font-medium text-white shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
            >
              Go to Sign Up
            </router-link>
            <p class="text-xs text-gray-500">
              Your verification link may have expired. Please sign up again to receive a new link.
            </p>
          </div>
        </div>

        <!-- Success State -->
        <div v-else class="text-center space-y-8">
          <div class="flex flex-col items-center">
            <div class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-green-100">
              <svg class="h-8 w-8 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <h3 class="mt-4 text-2xl font-bold text-gray-900">Email Verified!</h3>
            <p class="mt-2 text-sm text-gray-600">Your account is now active. Here's your API key:</p>
          </div>

          <div class="rounded-lg bg-amber-50 border border-amber-200 p-4">
            <div class="flex">
              <svg class="h-5 w-5 text-amber-400 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              <p class="ml-3 text-sm text-amber-800">
                <strong>Save this key now!</strong> It will not be shown again.
              </p>
            </div>
          </div>

          <div class="relative rounded-md bg-gray-50 p-4 border border-gray-200">
            <code class="font-mono text-sm text-blue-600 break-all block pr-20">{{ apiKey }}</code>
            <button
              @click="copyToClipboard"
              class="absolute right-2 top-1/2 -translate-y-1/2 inline-flex items-center rounded-md border border-gray-300 bg-white px-3 py-2 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
              type="button"
            >
              <svg v-if="!copied" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
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
              <li>Make your first API call to <code class="bg-blue-100 px-1 rounded">/v1/detect</code></li>
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

          <router-link
            to="/dashboard"
            class="block w-full text-center rounded-md border border-transparent bg-blue-600 py-3 px-4 text-sm font-medium text-white shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
          >
            Go to Dashboard
          </router-link>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const loading = ref(true)
const error = ref('')
const apiKey = ref('')
const copied = ref(false)

onMounted(async () => {
  const token = route.query.token as string

  if (!token) {
    error.value = 'No verification token provided.'
    loading.value = false
    return
  }

  try {
    const response = await fetch(`/api/v1/onboard/verify?token=${encodeURIComponent(token)}`)

    // Check content type to avoid parsing HTML as JSON
    const contentType = response.headers.get("content-type")
    if (!contentType || !contentType.includes("application/json")) {
      throw new Error('Server returned an unexpected response. Please try again.')
    }

    const data = await response.json()

    if (!response.ok) {
      let errorMessage = 'Verification failed. Please try again.'
      if (data.error) {
        if (typeof data.error === 'string') {
          errorMessage = data.error
        } else if (typeof data.error === 'object' && data.error.message) {
          errorMessage = data.error.message
        }
      }
      throw new Error(errorMessage)
    }

    if (data.success && data.api_key) {
      apiKey.value = data.api_key

      // Track conversion
      if (typeof window !== 'undefined' && (window as any).gtag) {
        (window as any).gtag('event', 'email_verified')
      }
    } else {
      throw new Error('Invalid response from server')
    }
  } catch (err: any) {
    error.value = err.message || 'Verification failed. Please try again.'
  } finally {
    loading.value = false
  }
})

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
</script>
