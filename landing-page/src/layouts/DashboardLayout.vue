<template>
  <div class="min-h-screen bg-white flex">
    <!-- Sidebar -->
    <div class="hidden md:flex md:w-64 md:flex-col md:fixed md:inset-y-0 bg-white border-r-2 border-black">
      <div class="flex items-center h-20 flex-shrink-0 px-4 bg-white border-b-2 border-black">
        <div class="h-8 w-8 bg-black flex items-center justify-center border border-black mr-3">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M2 12H22M12 2V22M4.5 4.5L19.5 19.5M4.5 19.5L19.5 4.5" stroke="white" stroke-width="3"/>
          </svg>
        </div>
        <span class="font-sans font-black text-xl tracking-tighter uppercase text-black">Driftlock</span>
      </div>
      <div class="flex-1 flex flex-col overflow-y-auto">
        <nav class="flex-1 px-4 py-6 space-y-2">
          <router-link to="/dashboard" class="group flex items-center px-3 py-3 text-sm font-bold uppercase tracking-widest border-2 border-black bg-black text-white hover:bg-white hover:text-black transition-colors shadow-[4px_4px_0px_0px_rgba(0,0,0,0)] hover:shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]">
            <svg class="mr-3 h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="square" stroke-linejoin="miter" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
            </svg>
            Dashboard
          </router-link>
          <!-- Add more links here later: Anomalies, Settings, etc. -->
        </nav>
      </div>
      <div class="flex-shrink-0 flex bg-white border-t-2 border-black p-4">
        <div class="flex items-center w-full">
          <div class="w-full">
            <p class="text-xs font-mono font-bold text-black truncate uppercase mb-2">Logged in as:</p>
            <p class="text-sm font-bold text-black truncate border-b-2 border-black pb-1 mb-2">{{ authStore.user?.email }}</p>
            <button @click="logout" class="w-full text-xs font-bold uppercase tracking-widest text-black border-2 border-black py-2 hover:bg-black hover:text-white transition-colors">Sign out</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="flex-1 flex flex-col md:pl-64">
      <main class="flex-1 bg-white">
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


