<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50">
    <div class="text-center">
      <svg v-if="authStore.loading" class="mx-auto h-12 w-12 animate-spin text-blue-600" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"></path>
      </svg>
      <div v-else>
        <h2 class="mt-6 text-3xl font-bold tracking-tight text-gray-900">Verifying login...</h2>
        <p class="mt-2 text-sm text-gray-600">Please wait while we securely log you in.</p>
      </div>
      <div v-if="authStore.error" class="mt-4 text-red-600">
        {{ authStore.error }}
        <p class="mt-2">
          <router-link to="/login" class="text-blue-600 hover:underline">Return to login</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const router = useRouter()

onMounted(async () => {
  try {
    const success = await authStore.completeMagicLinkLogin()
    if (success) {
      router.push('/dashboard')
    } else {
      router.push('/login')
    }
  } catch (e) {
    // Error displayed in template
  }
})
</script>





