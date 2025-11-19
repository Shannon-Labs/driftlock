<template>
    <div class="rounded-3xl bg-white p-8 shadow-2xl">
        <form @submit.prevent="handleSignup" novalidate>
            <fieldset :disabled="isSubmitting" class="space-y-6">
                <legend class="text-2xl font-mono font-bold text-gray-900">Get Started with Driftlock</legend>
                <p class="text-sm text-gray-600">Create your free trial account and get an API key instantly.</p>

                <div>
                    <label for="signup-email" class="block text-sm font-sans font-semibold text-gray-700">Work Email *</label>
                    <input
                        id="signup-email"
                        v-model="form.email"
                        type="email"
                        required
                        autocomplete="email"
                        class="mt-2 w-full rounded-2xl border border-gray-200 bg-gray-50 px-4 py-3 text-base font-sans text-gray-900 shadow-sm focus:border-blue-500 focus:bg-white focus:outline-none focus:ring-2 focus:ring-blue-500/40"
                        placeholder="you@company.com"
                    />
                </div>

                <div>
                    <label for="signup-company" class="block text-sm font-sans font-semibold text-gray-700">Company Name *</label>
                    <input
                        id="signup-company"
                        v-model="form.companyName"
                        type="text"
                        required
                        autocomplete="organization"
                        class="mt-2 w-full rounded-2xl border border-gray-200 bg-gray-50 px-4 py-3 text-base font-sans text-gray-900 shadow-sm focus:border-blue-500 focus:bg-white focus:outline-none focus:ring-2 focus:ring-blue-500/40"
                        placeholder="Your Company"
                    />
                </div>

                <button
                    type="submit"
                    :disabled="!isFormValid || isSubmitting"
                    class="inline-flex w-full items-center justify-center rounded-2xl bg-gradient-to-r from-blue-600 to-indigo-600 px-6 py-3 text-base font-sans font-semibold text-white shadow-lg transition focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 disabled:cursor-not-allowed disabled:opacity-50"
                >
                    <svg v-if="isSubmitting" class="mr-3 h-5 w-5 animate-spin text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"></path>
                    </svg>
                    {{ isSubmitting ? 'Creating Account...' : 'Create Free Account' }}
                </button>

                <p class="text-xs font-sans text-gray-500 text-center">
                    Free trial includes 10,000 events/month for 14 days.
                </p>
            </fieldset>
        </form>

        <!-- Result display -->
        <div class="mt-6 min-h-[48px]" aria-live="polite">
            <div v-if="state === 'success'" class="rounded-2xl border border-green-200 bg-green-50 p-4">
                <div class="flex items-start">
                    <svg class="mr-3 h-5 w-5 text-green-500 mt-0.5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-7.778 7.778a1 1 0 01-1.414 0L3.293 9.293a1 1 0 011.414-1.414l3.102 3.102 7.071-7.071a1 1 0 011.414 0z" clip-rule="evenodd" />
                    </svg>
                    <div class="flex-1">
                        <p class="text-sm font-semibold text-green-900">Account Created Successfully!</p>
                        <p class="mt-1 text-xs text-green-700">Save your API key below - it won't be shown again.</p>
                    </div>
                </div>
                <div class="mt-4 rounded-xl bg-gray-900 p-4 font-mono text-sm text-green-400 break-all">
                    {{ apiKey }}
                </div>
                <button
                    @click="copyApiKey"
                    class="mt-3 w-full rounded-xl border border-green-300 bg-white px-4 py-2 text-sm font-semibold text-green-700 hover:bg-green-50 transition"
                >
                    {{ copied ? 'Copied!' : 'Copy API Key' }}
                </button>
                <div class="mt-4 p-3 rounded-xl bg-blue-50 border border-blue-200">
                    <p class="text-xs font-semibold text-blue-900">Next Steps:</p>
                    <ol class="mt-2 text-xs text-blue-700 space-y-1 list-decimal list-inside">
                        <li>Save your API key securely</li>
                        <li>Check out the <a href="#api-demo" class="underline font-semibold">API Demo</a></li>
                        <li>Make your first API call to <code class="bg-blue-100 px-1 rounded">/v1/detect</code></li>
                    </ol>
                </div>
            </div>

            <div v-else-if="state === 'error'" class="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm font-sans text-red-900">
                {{ errorMessage }}
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'

const signupEndpoint = '/api/v1/onboard/signup'
const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

const form = reactive({
    email: '',
    companyName: '',
})

const isSubmitting = ref(false)
const state = ref<'idle' | 'success' | 'error'>('idle')
const errorMessage = ref('')
const apiKey = ref('')
const copied = ref(false)

const isFormValid = computed(() => {
    return (
        emailRegex.test(form.email.trim()) &&
        form.companyName.trim().length >= 2
    )
})

const resetForm = () => {
    form.email = ''
    form.companyName = ''
}

const copyApiKey = async () => {
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

const handleSignup = async () => {
    if (!isFormValid.value) {
        state.value = 'error'
        errorMessage.value = 'Please fill out all required fields correctly.'
        return
    }

    isSubmitting.value = true
    state.value = 'idle'
    errorMessage.value = ''

    try {
        const response = await fetch(signupEndpoint, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                email: form.email.trim(),
                company_name: form.companyName.trim(),
                plan: 'trial',
            }),
        })

        const data = await response.json()

        if (!response.ok) {
            throw new Error(data.error?.message || `Request failed with status ${response.status}`)
        }

        apiKey.value = data.api_key
        state.value = 'success'
        resetForm()
    } catch (error) {
        console.error('Signup failed', error)
        state.value = 'error'
        errorMessage.value = error instanceof Error ? error.message : 'Something went wrong. Please try again.'
    } finally {
        isSubmitting.value = false
    }
}
</script>
