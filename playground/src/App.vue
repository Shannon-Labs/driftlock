<template>
  <div class="max-w-6xl mx-auto p-6">
    <header class="mb-6">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Driftlock Playground</h1>
          <p class="text-gray-600 dark:text-gray-300">Upload or paste JSON/NDJSON, tune parameters, and run detection.</p>
        </div>
        <div class="flex items-center gap-2">
          <div class="flex items-center gap-2 px-3 py-1.5 rounded-lg border" :class="apiStatus === 'connected' ? 'bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-800' : apiStatus === 'checking' ? 'bg-yellow-50 dark:bg-yellow-900/20 border-yellow-200 dark:border-yellow-800' : 'bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800'">
            <div class="w-2 h-2 rounded-full" :class="apiStatus === 'connected' ? 'bg-green-500' : apiStatus === 'checking' ? 'bg-yellow-500 animate-pulse' : 'bg-red-500'"></div>
            <span class="text-sm font-medium" :class="apiStatus === 'connected' ? 'text-green-700 dark:text-green-400' : apiStatus === 'checking' ? 'text-yellow-700 dark:text-yellow-400' : 'text-red-700 dark:text-red-400'">
              {{ apiStatus === 'connected' ? 'API Connected' : apiStatus === 'checking' ? 'Checking...' : 'API Unavailable' }}
            </span>
          </div>
        </div>
      </div>
    </header>

    <section class="grid md:grid-cols-3 gap-6 mb-6">
      <div class="md:col-span-2 space-y-4">
        <UploadPanel @data="onData" />
        <SamplePicker @load="onSample" />
      </div>
      <div>
        <ParamsForm :params="params" @update="onParams" @run="runDetect" />
      </div>
    </section>

    <section v-if="loading" class="mt-4 p-6 rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900">
      <div class="flex items-center justify-center gap-3">
        <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
        <span class="text-gray-600 dark:text-gray-300">Processing detection...</span>
      </div>
    </section>

    <section v-if="response" class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-700 p-4">
      <div class="flex items-center justify-between mb-3">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white">Results</h2>
        <CurlSnippet :curl="curlCmd" />
      </div>
      <div class="text-sm text-gray-600 dark:text-gray-400 mb-4">
        <span class="mr-4">Total: {{ response.total_events }}</span>
        <span class="mr-4">Anomalies: {{ response.anomaly_count }}</span>
        <span class="mr-4">Algo: {{ response.compression_algo }}</span>
        <span>Time: {{ response.processing_time }}</span>
      </div>
      <ResultsTable :items="response.anomalies" />
      <div class="mt-4">
        <button class="px-3 py-2 text-sm bg-gray-100 dark:bg-gray-800 rounded border" @click="downloadJSON">Download JSON</button>
      </div>
    </section>

    <section v-if="error" class="mt-4 p-3 rounded border border-red-300 bg-red-50 text-red-700 dark:bg-red-900/20 dark:text-red-300">
      {{ error }}
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import UploadPanel from './components/UploadPanel.vue'
import ParamsForm from './components/ParamsForm.vue'
import ResultsTable from './components/ResultsTable.vue'
import SamplePicker from './components/SamplePicker.vue'
import CurlSnippet from './components/CurlSnippet.vue'

const apiBase = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

const raw = ref<string>('')          // user data
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
    const timeoutId = setTimeout(() => controller.abort(), 5000) // 5 second timeout
    
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
  // Recheck every 30 seconds
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
  
  // Check API health before making request
  if (apiStatus.value !== 'connected') {
    await checkApiHealth()
    if (apiStatus.value !== 'connected') {
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


