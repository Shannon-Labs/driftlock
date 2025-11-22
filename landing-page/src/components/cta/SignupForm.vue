<template>
  <div class="mx-auto max-w-lg p-4 sm:p-8">
    <form @submit.prevent="handleSignup" class="bg-white rounded-xl shadow-xl border border-gray-100 p-8 sm:p-10">
      <div v-if="!submitted" class="space-y-6">
        <div class="text-center">
          <h2 class="text-2xl font-bold text-gray-900">Start Your Free Trial</h2>
          <p class="mt-2 text-sm text-gray-500">Get instant access to Driftlock's anomaly detection API</p>
        </div>
        
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

        <div>
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

        <div v-if="error" class="rounded-md bg-red-50 p-4 text-sm text-red-700 flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          {{ error }}
        </div>

        <button
          type="submit"
          class="flex w-full justify-center rounded-md border border-transparent bg-blue-600 py-3 px-4 text-sm font-medium text-white shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-70 disabled:cursor-not-allowed transition-all"
          :disabled="loading"
        >
          <span v-if="!loading">Start Free Trial →</span>
          <span v-else class="flex items-center gap-2">
            <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Creating account...
          </span>
        </button>

        <p class="text-center text-xs text-gray-500">
          Free trial includes 10,000 events. No credit card required.
        </p>
      </div>

      <div v-else class="text-center space-y-8">
        <div class="flex flex-col items-center">
           <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-green-100">
              <svg class="h-6 w-6 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <h3 class="mt-3 text-2xl font-bold text-gray-900">Welcome to Driftlock!</h3>
            <p class="mt-2 text-sm text-gray-500">Your account has been created. Here's your API key:</p>
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
            <li>Save your API key securely (you won't see it again)</li>
            <li>Read the <a href="/docs" target="_blank" class="underline font-medium">API documentation</a></li>
            <li>Make your first API call to <code>/v1/detect</code></li>
          </ol>
        </div>

        <div class="text-left">
          <h4 class="text-xs font-semibold uppercase tracking-wider text-gray-500 mb-2">Quick Start:</h4>
          <div class="bg-gray-900 rounded-lg p-4 overflow-x-auto">
              <pre class="text-xs text-gray-300 font-mono"><code>curl -X POST https://driftlock.web.app/api/v1/detect \
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
                 <button @click="handleUpgrade('basic')" class="flex-1 inline-flex justify-center items-center rounded-md border border-blue-200 bg-white px-3 py-2 text-xs font-medium text-blue-700 shadow-sm hover:bg-blue-50 focus:outline-none" :disabled="upgrading">
                    <span v-if="!upgrading">Basic ($20)</span>
                    <span v-else>...</span>
                 </button>
                 <button @click="handleUpgrade('pro')" class="flex-1 inline-flex justify-center items-center rounded-md border border-transparent bg-blue-600 px-3 py-2 text-xs font-medium text-white shadow-sm hover:bg-blue-700 focus:outline-none" :disabled="upgrading">
                    <span v-if="!upgrading">Pro ($200)</span>
                    <span v-else>...</span>
                 </button>
             </div>
          </div>
        </div>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const email = ref('')
const company = ref('')
const loading = ref(false)
const error = ref('')
const submitted = ref(false)
const apiKey = ref('')
const copied = ref(false)
const upgrading = ref(false)

const handleSignup = async () => {
  error.value = ''
  loading.value = true

  try {
    // Use relative path to leverage Firebase Hosting rewrites
    const response = await fetch('/api/v1/onboard/signup', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: email.value,
        company_name: company.value,
        plan: 'trial',
        source: 'landing_page'
      }),
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.error || 'Failed to create account. Please try again.')
    }

    if (data.success && data.api_key) {
      apiKey.value = data.api_key
      submitted.value = true
      
      // Track signup event (if analytics is set up)
      if (typeof window !== 'undefined' && (window as any).gtag) {
        (window as any).gtag('event', 'signup', {
          method: 'landing_page'
        })
      }
    } else {
      throw new Error('Invalid response from server')
    }
  } catch (err: any) {
    error.value = err.message || 'Something went wrong. Please try again.'
    console.error('Signup error:', err)
  } finally {
    loading.value = false
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
    // Use proxy to route through Firebase Functions to Cloud Run
    const response = await fetch('/api/proxy/v1/billing/checkout', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Api-Key': apiKey.value
      },
      body: JSON.stringify({ plan })
    })
    
    if (!response.ok) {
      throw new Error('Failed to start checkout')
    }
    
    const data = await response.json()
    if (data.url) {
      window.location.href = data.url
    }
  } catch (err) {
    console.error('Upgrade error:', err)
    alert('Failed to start upgrade process. Please try again.')
  } finally {
    upgrading.value = false
  }
}
</script>
