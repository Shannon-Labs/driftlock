<template>
  <div class="min-h-screen bg-gray-50">
    <div class="bg-white border-b border-gray-200">
      <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-4 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <a href="/" class="text-xl font-mono font-bold text-gray-900">Driftlock</a>
          <span class="rounded-full bg-blue-100 px-2.5 py-0.5 text-xs font-medium text-blue-800">Dashboard</span>
        </div>
        <div class="flex items-center gap-4">
           <!-- Placeholder for Auth User -->
           <div class="h-8 w-8 rounded-full bg-gray-200 flex items-center justify-center text-xs font-bold text-gray-500">
             JD
           </div>
        </div>
      </div>
    </div>

    <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <header class="mb-8">
        <h1 class="text-3xl font-mono font-bold text-gray-900">Your API Keys</h1>
        <p class="mt-2 text-sm text-gray-600">Manage access to the Driftlock detection engine.</p>
      </header>

      <div class="bg-white shadow-sm ring-1 ring-gray-900/5 sm:rounded-xl md:col-span-2 mb-8">
        <div class="px-4 py-6 sm:p-8">
          <div class="flex items-center justify-between mb-6">
             <h2 class="text-base font-semibold leading-7 text-gray-900">Active Keys</h2>
             <button @click="generateKey" class="rounded-md bg-blue-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600">
               Generate New Key
             </button>
          </div>
          
          <div v-if="keys.length > 0" class="flow-root">
            <ul role="list" class="-my-5 divide-y divide-gray-100">
              <li v-for="key in keys" :key="key.id" class="py-5">
                <div class="flex items-center justify-between gap-x-6 py-5">
                  <div class="min-w-0">
                    <div class="flex items-start gap-x-3">
                      <p class="text-sm font-mono font-semibold leading-6 text-gray-900">{{ key.prefix }}...</p>
                      <p :class="[key.status === 'Active' ? 'text-green-700 bg-green-50 ring-green-600/20' : 'text-gray-600 bg-gray-50 ring-gray-500/10', 'rounded-md whitespace-nowrap mt-0.5 px-1.5 py-0.5 text-xs font-medium ring-1 ring-inset']">{{ key.status }}</p>
                    </div>
                    <div class="mt-1 flex items-center gap-x-2 text-xs leading-5 text-gray-500">
                      <p class="whitespace-nowrap">Created on <time :datetime="key.created">{{ key.createdDate }}</time></p>
                    </div>
                  </div>
                  <div class="flex flex-none items-center gap-x-4">
                    <button @click="copyKey(key.fullKey)" class="hidden rounded-md bg-white px-2.5 py-1.5 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:block">Copy</button>
                  </div>
                </div>
              </li>
            </ul>
          </div>
          <div v-else class="text-center py-12">
             <p class="text-sm text-gray-500">No API keys found. Generate one to get started.</p>
          </div>

        </div>
      </div>

      <header class="mb-8">
        <div class="flex items-center justify-between">
            <div>
                <h1 class="text-3xl font-mono font-bold text-gray-900">Live Stream</h1>
                <p class="mt-2 text-sm text-gray-600">Real-time anomalies from your agents.</p>
            </div>
            <div class="flex items-center gap-2">
                <span class="relative flex h-3 w-3">
                  <span v-if="isConnected" class="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                  <span :class="isConnected ? 'bg-green-500' : 'bg-red-500'" class="relative inline-flex rounded-full h-3 w-3"></span>
                </span>
                <span class="text-sm font-medium text-gray-700">{{ isConnected ? 'Connected' : 'Disconnected' }}</span>
            </div>
        </div>
      </header>

      <div class="bg-slate-950 shadow-sm ring-1 ring-white/10 sm:rounded-xl md:col-span-2 h-[400px] overflow-y-auto font-mono text-sm p-4 text-slate-300">
        <div v-if="liveAnomalies.length === 0" class="h-full flex flex-col items-center justify-center text-slate-500">
            <p>Waiting for anomalies...</p>
            <p class="text-xs mt-2">Send events to /v1/detect to see them here.</p>
        </div>
        <div v-else class="space-y-2">
            <div v-for="anomaly in liveAnomalies" :key="anomaly.id" class="p-3 rounded border border-red-500/20 bg-red-500/5 animate-fade-in-down">
                <div class="flex justify-between text-red-300 mb-1">
                    <span class="font-bold">⚠ ANOMALY DETECTED</span>
                    <span>{{ new Date(anomaly.detected_at).toLocaleTimeString() }}</span>
                </div>
                <div class="text-slate-300">{{ anomaly.explanation }}</div>
                <div class="mt-2 grid grid-cols-2 gap-2 text-xs text-slate-500">
                    <div>NCD: <span class="text-slate-200">{{ anomaly.ncd.toFixed(4) }}</span></div>
                    <div>Confidence: <span class="text-slate-200">{{ (anomaly.confidence * 100).toFixed(1) }}%</span></div>
                </div>
            </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useAnomalyStore, useAuthStore } from '../stores'

const anomalyStore = useAnomalyStore()
const authStore = useAuthStore()

const { anomalies: liveAnomalies, isConnected } = storeToRefs(anomalyStore)
const { keys } = storeToRefs(authStore)

onMounted(() => {
    // Connect using the first available key (demo mode)
    if (keys.value.length > 0) {
        anomalyStore.connect(keys.value[0].fullKey)
    }
})

onUnmounted(() => {
    anomalyStore.disconnect()
})

const generateKey = authStore.generateKey

const copyKey = (key: string) => {
  navigator.clipboard.writeText(key)
  alert('Copied to clipboard!')
}
</script>

<style scoped>
.animate-fade-in-down {
  animation: fadeInDown 0.5s ease-out;
}

@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>

