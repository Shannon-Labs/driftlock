<template>
  <div class="signup-form-container">
    <form @submit.prevent="handleSignup" class="signup-form">
      <div v-if="!submitted" class="form-fields">
        <h2 class="form-title">Start Your Free Trial</h2>
        <p class="form-subtitle">Get instant access to Driftlock's anomaly detection API</p>
        
        <div class="form-group">
          <label for="email" class="form-label">Work Email</label>
          <input
            id="email"
            v-model="email"
            type="email"
            placeholder="you@company.com"
            required
            class="form-input"
            :disabled="loading"
          />
        </div>

        <div class="form-group">
          <label for="company" class="form-label">Company Name</label>
          <input
            id="company"
            v-model="company"
            type="text"
            placeholder="Your Company"
            required
            minlength="2"
            maxlength="100"
            class="form-input"
            :disabled="loading"
          />
        </div>

        <div v-if="error" class="error-message">
          <svg xmlns="http://www.w3.org/2000/svg" class="error-icon" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          {{ error }}
        </div>

        <button
          type="submit"
          class="submit-button"
          :disabled="loading"
        >
          <span v-if="!loading">Start Free Trial →</span>
          <span v-else class="loading-spinner">
            <svg class="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Creating account...
          </span>
        </button>

        <p class="form-footer">
          Free trial includes 10,000 events. No credit card required.
        </p>
      </div>

      <div v-else class="success-message">
        <svg xmlns="http://www.w3.org/2000/svg" class="success-icon" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <h3 class="success-title">Welcome to Driftlock!</h3>
        <p class="success-subtitle">Your account has been created. Here's your API key:</p>
        
        <div class="api-key-container">
          <code class="api-key">{{ apiKey }}</code>
          <button @click="copyToClipboard" class="copy-button" type="button">
            <svg v-if="!copied" xmlns="http://www.w3.org/2000/svg" class="copy-icon" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" class="copy-icon" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            {{ copied ? 'Copied!' : 'Copy' }}
          </button>
        </div>

        <div class="next-steps">
          <h4 class="next-steps-title">Next Steps:</h4>
          <ol class="next-steps-list">
            <li>Save your API key securely (you won't see it again)</li>
            <li>Read the <a href="/docs" target="_blank" class="docs-link">API documentation</a></li>
            <li>Make your first API call to <code>/v1/detect</code></li>
          </ol>
        </div>

        <div class="quick-start">
          <h4 class="quick-start-title">Quick Start:</h4>
          <pre class="code-block"><code>curl -X POST https://driftlock.net/api/v1/detect \
  -H "X-Api-Key: {{ apiKey }}" \
  -H "Content-Type: application/json" \
  -d '{"events": [...], "window_size": 50}'</code></pre>
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

const handleSignup = async () => {
  error.value = ''
  loading.value = true

  try {
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
</script>

<style scoped>
.signup-form-container {
  max-width: 500px;
  margin: 0 auto;
  padding: 2rem;
}

.signup-form {
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1);
  padding: 2.5rem;
}

.form-title {
  font-size: 1.875rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  color: #111827;
}

.form-subtitle {
  color: #6b7280;
  margin-bottom: 2rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  margin-bottom: 0.5rem;
  color: #374151;
}

.form-input {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 1px solid #d1d5db;
  border-radius: 0.5rem;
  font-size: 1rem;
  transition: all 0.2s;
}

.form-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-input:disabled {
  background-color: #f9fafb;
  cursor: not-allowed;
}

.error-message {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  background-color: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 0.5rem;
  color: #dc2626;
  font-size: 0.875rem;
  margin-bottom: 1rem;
}

.error-icon {
  width: 1.25rem;
  height: 1.25rem;
  flex-shrink: 0;
}

.submit-button {
  width: 100%;
  padding: 0.875rem 1.5rem;
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  color: white;
  border: none;
  border-radius: 0.5rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.submit-button:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1);
}

.submit-button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.loading-spinner {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.form-footer {
  margin-top: 1rem;
  text-align: center;
  font-size: 0.875rem;
  color: #6b7280;
}

/* Success state */
.success-message {
  text-align: center;
}

.success-icon {
  width: 4rem;
  height: 4rem;
  margin: 0 auto 1rem;
  color: #10b981;
}

.success-title {
  font-size: 1.875rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  color: #111827;
}

.success-subtitle {
  color: #6b7280;
  margin-bottom: 1.5rem;
}

.api-key-container {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 2rem;
  padding: 1rem;
  background-color: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 0.5rem;
}

.api-key {
  flex: 1;
  font-family: 'Courier New', monospace;
  font-size: 0.875rem;
  word-break: break-all;
  color: #111827;
}

.copy-button {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.5rem 1rem;
  background-color: white;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: #374151;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.copy-button:hover {
  background-color: #f9fafb;
  border-color: #9ca3af;
}

.copy-icon {
  width: 1rem;
  height: 1rem;
}

.next-steps, .quick-start {
  text-align: left;
  margin-top: 2rem;
  padding: 1.5rem;
  background-color: #f9fafb;
  border-radius: 0.5rem;
}

.next-steps-title, .quick-start-title {
  font-size: 1rem;
  font-weight: 600;
  margin-bottom: 1rem;
  color: #111827;
}

.next-steps-list {
  margin: 0;
  padding-left: 1.5rem;
  color: #374151;
}

.next-steps-list li {
  margin-bottom: 0.5rem;
}

.docs-link {
  color: #3b82f6;
  text-decoration: underline;
}

.code-block {
  background-color: #1f2937;
  color: #f3f4f6;
  padding: 1rem;
  border-radius: 0.375rem;
  overflow-x: auto;
  font-size: 0.75rem;
  line-height: 1.5;
}

.code-block code {
  font-family: 'Courier New', monospace;
}
</style>

