<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 via-white to-blue-50/30 dark:from-gray-900 dark:via-gray-900 dark:to-gray-800">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <header class="mb-8">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-4xl font-bold bg-gradient-to-r from-blue-600 via-indigo-600 to-cyan-600 bg-clip-text text-transparent mb-2">
              Driftlock Playground
            </h1>
            <p class="text-gray-600 dark:text-gray-300">Upload or paste JSON/NDJSON, tune parameters, and run anomaly detection.</p>
          </div>
          <div class="flex items-center gap-3">
            <div 
              class="flex items-center gap-2 px-4 py-2 rounded-lg border shadow-sm transition-all duration-200" 
              :class="apiStatus === 'connected' ? 'bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-800' : apiStatus === 'checking' ? 'bg-yellow-50 dark:bg-yellow-900/20 border-yellow-200 dark:border-yellow-800' : 'bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800'"
            >
              <div 
                class="w-2.5 h-2.5 rounded-full" 
                :class="apiStatus === 'connected' ? 'bg-green-500 animate-pulse' : apiStatus === 'checking' ? 'bg-yellow-500 animate-pulse' : 'bg-red-500'"
              ></div>
              <span 
                class="text-sm font-medium" 
                :class="apiStatus === 'connected' ? 'text-green-700 dark:text-green-400' : apiStatus === 'checking' ? 'text-yellow-700 dark:text-yellow-400' : 'text-red-700 dark:text-red-400'"
              >
                {{ apiStatus === 'connected' ? 'API Connected' : apiStatus === 'checking' ? 'Checking...' : 'API Unavailable' }}
              </span>
            </div>
          </div>
        </div>
      </header>

      <section class="grid lg:grid-cols-3 gap-6 mb-6">
        <div class="lg:col-span-2 space-y-6">
          <UploadPanel @data="onData" />
          <SamplePicker @load="onSample" />
        </div>
        <div>
          <ParamsForm :params="params" @update="onParams" @run="runDetect" />
        </div>
      </section>

      <section v-if="loading" class="mt-6 p-8 rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 shadow-sm">
        <div class="flex items-center justify-center gap-4">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          <span class="text-gray-600 dark:text-gray-300 text-lg font-medium">Processing detection...</span>
        </div>
      </section>

      <section v-if="response" class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-700 p-6 shadow-sm">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-2xl font-bold text-gray-900 dark:text-white">Results</h2>
          <CurlSnippet :curl="curlCmd" />
        </div>
        <div class="flex flex-wrap gap-6 mb-6 p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
          <div class="flex items-center gap-2">
            <span class="text-sm text-gray-600 dark:text-gray-400">Total Events:</span>
            <span class="text-lg font-semibold text-gray-900 dark:text-white">{{ response.total_events }}</span>
          </div>
          <div class="flex items-center gap-2">
            <span class="text-sm text-gray-600 dark:text-gray-400">Anomalies:</span>
            <span class="text-lg font-semibold" :class="response.anomaly_count > 0 ? 'text-red-600 dark:text-red-400' : 'text-green-600 dark:text-green-400'">
              {{ response.anomaly_count || 0 }}
            </span>
          </div>
          <div class="flex items-center gap-2">
            <span class="text-sm text-gray-600 dark:text-gray-400">Algorithm:</span>
            <span class="text-lg font-semibold text-gray-900 dark:text-white">{{ response.compression_algo }}</span>
          </div>
          <div class="flex items-center gap-2">
            <span class="text-sm text-gray-600 dark:text-gray-400">Processing Time:</span>
            <span class="text-lg font-semibold text-gray-900 dark:text-white">{{ response.processing_time }}</span>
          </div>
        </div>
        <ResultsTable :items="response.anomalies || []" />
        <div class="mt-6">
          <button 
            class="px-4 py-2 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white rounded-lg font-medium transition-all duration-200 transform hover:scale-105 shadow-sm" 
            @click="downloadJSON"
          >
            Download JSON Results
          </button>
        </div>
      </section>

      <section v-if="error" class="mt-6 p-4 rounded-xl border border-red-300 dark:border-red-800 bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-300 shadow-sm">
        <div class="flex items-center gap-2">
          <span class="text-xl">⚠️</span>
          <span class="font-medium">{{ error }}</span>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import UploadPanel from '../components/playground/UploadPanel.vue'
import ParamsForm from '../components/playground/ParamsForm.vue'
import ResultsTable from '../components/playground/ResultsTable.vue'
import SamplePicker from '../components/playground/SamplePicker.vue'
import CurlSnippet from '../components/playground/CurlSnippet.vue'

const apiBase = import.meta.env.VITE_API_BASE_URL || 'https://driftlock.net/api/v1'

const raw = ref<string>('')
const mime = ref<'ndjson' | 'json'>('ndjson')
const params = ref({ baseline: 400, window: 1, hop: 1, algo: 'zstd' })
const response = ref<any | null>(null)
const error = ref<string | null>(null)
const loading = ref<boolean>(false)
const apiStatus = ref<'checking' | 'connected' | 'disconnected'>('checking')

async function checkApiHealth() {
  apiStatus.value = 'checking'
  try {
    const controller = new AbortController()
    const timeoutId = setTimeout(() => controller.abort(), 5000)
    
    const res = await fetch(`${apiBase}/healthz`, {
      method: 'GET',
      signal: controller.signal,
    })
    clearTimeout(timeoutId)
    
    if (res.ok) {
      const data = await res.json()
      if (data.success === true || data.ok === true) {
        apiStatus.value = 'connected'
        return
      }
    }
    apiStatus.value = 'disconnected'
  } catch (e) {
    apiStatus.value = 'disconnected'
  }
}

onMounted(() => {
  checkApiHealth()
  setInterval(checkApiHealth, 30000)
})

function onData(payload: { text: string, format: 'ndjson' | 'json' }) {
  raw.value = payload.text
  mime.value = payload.format
}

async function onSample(url: string) {
  const res = await fetch(url)
  const text = await res.text()
  raw.value = text
  mime.value = url.endsWith('.json') ? 'json' : 'ndjson'
}

function onParams(next: any) {
  params.value = next
}

const curlCmd = computed(() => {
  const q = new URLSearchParams({
    format: mime.value,
    baseline: String(params.value.baseline),
    window: String(params.value.window),
    hop: String(params.value.hop),
    algo: params.value.algo,
  })
  return `curl -s -X POST "${apiBase}/v1/detect?${q.toString()}" -H "Content-Type: application/json" --data-binary @your-file.${mime.value === 'ndjson' ? 'jsonl' : 'json'}`
})

async function runDetect() {
  error.value = null
  response.value = null
  loading.value = true
  
  const currentStatus = apiStatus.value
  if (currentStatus !== 'connected') {
    await checkApiHealth()
    const newStatus = apiStatus.value
    if (newStatus !== 'connected') {
      error.value = 'API is unavailable. Please check your connection and ensure the API server is running.'
      loading.value = false
      return
    }
  }
  
  try {
    const q = new URLSearchParams({
      format: mime.value,
      baseline: String(params.value.baseline),
      window: String(params.value.window),
      hop: String(params.value.hop),
      algo: params.value.algo,
    })
    const res = await fetch(`${apiBase}/v1/detect?${q.toString()}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: raw.value,
    })
    const json = await res.json()
    if (!res.ok) {
      error.value = json?.error || json?.message || `Request failed with status ${res.status}`
      return
    }
    response.value = json
  } catch (e: any) {
    if (e.name === 'AbortError' || e.name === 'TimeoutError') {
      error.value = 'Request timed out. The API may be slow or unavailable.'
    } else if (e.message?.includes('Failed to fetch') || e.message?.includes('NetworkError')) {
      error.value = 'Network error. Please check your connection and ensure the API server is running at ' + apiBase
      apiStatus.value = 'disconnected'
    } else {
      error.value = e?.message || String(e)
    }
  } finally {
    loading.value = false
  }
}

function downloadJSON() {
  if (!response.value) return
  const blob = new Blob([JSON.stringify(response.value, null, 2)], { type: 'application/json' })
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = 'driftlock-results.json'
  a.click()
}
</script>

