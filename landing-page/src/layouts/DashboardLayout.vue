<template>
  <div class="min-h-screen bg-gray-50 flex">
    <!-- Sidebar -->
    <div class="hidden md:flex md:w-64 md:flex-col md:fixed md:inset-y-0 bg-gray-900 text-white">
      <div class="flex items-center h-16 flex-shrink-0 px-4 bg-gray-900 border-b border-gray-800">
        <img class="h-8 w-auto" src="/logo.svg" alt="Driftlock" />
        <span class="ml-2 font-mono font-bold text-lg tracking-tight text-white">Driftlock</span>
      </div>
      <div class="flex-1 flex flex-col overflow-y-auto">
        <nav class="flex-1 px-2 py-4 space-y-1">
          <router-link to="/dashboard" class="group flex items-center px-2 py-2 text-sm font-medium rounded-md bg-gray-800 text-white">
            <svg class="mr-3 h-6 w-6 text-gray-300" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
            </svg>
            Dashboard
          </router-link>
          <!-- Add more links here later: Anomalies, Settings, etc. -->
        </nav>
      </div>
      <div class="flex-shrink-0 flex bg-gray-800 p-4">
        <div class="flex items-center w-full">
          <div class="ml-3 w-full">
            <p class="text-sm font-medium text-white truncate">{{ authStore.user?.email }}</p>
            <button @click="logout" class="text-xs font-medium text-gray-400 hover:text-white mt-1">Sign out</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="flex-1 flex flex-col md:pl-64">
      <main class="flex-1">
        <slot></slot>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const router = useRouter()

const logout = async () => {
  await authStore.logout()
  router.push('/login')
}
</script>


