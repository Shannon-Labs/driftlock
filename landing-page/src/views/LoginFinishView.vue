<template>
  <div class="min-h-screen flex items-center justify-center bg-white">
    <div class="text-center p-8 border-2 border-black shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
      <div v-if="authStore.loading" class="mx-auto">
        <div class="animate-spin h-12 w-12 border-4 border-black border-t-transparent rounded-full mx-auto"></div>
      </div>
      <div v-else>
        <h2 class="mt-6 text-3xl font-sans font-black uppercase tracking-tighter text-black">Verifying login...</h2>
        <p class="mt-2 text-sm font-mono text-gray-600">Please wait while we securely log you in.</p>
      </div>
      <div v-if="authStore.error" class="mt-6 text-red-600 font-bold border border-red-600 p-2 bg-red-50">
        {{ authStore.error }}
        <p class="mt-4">
          <router-link to="/login" class="text-black hover:text-white hover:bg-black px-4 py-2 border-2 border-black font-bold uppercase tracking-widest transition-colors inline-block">Return to login</router-link>
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





