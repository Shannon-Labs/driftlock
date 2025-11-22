<template>
  <div :class="rootClasses">
    <header class="mb-10 border-b-4 border-black pb-6">
      <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <p class="text-xs font-bold uppercase tracking-widest text-gray-500 mb-2">
            Live Playground
          </p>
          <h2 class="text-3xl lg:text-4xl font-sans font-black uppercase tracking-tighter text-black mb-2">
            CBAD Mission Control
          </h2>
          <p class="text-sm font-serif text-gray-800 max-w-2xl">
            Paste NDJSON, stream a synthetic burst, and watch Driftlock surface anomalies with audit-ready math.
          </p>
        </div>
        <div class="flex items-center gap-3">
          <div class="flex items-center gap-2 border-2 border-black bg-white px-4 py-2 text-sm font-bold uppercase tracking-wider text-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]">
            <span
              class="h-3 w-3 border border-black"
              :class="{
                'bg-green-500': apiStatus === 'connected',
                'bg-yellow-400 animate-pulse': apiStatus === 'checking',
                'bg-red-500': apiStatus === 'disconnected'
              }"
            ></span>
            {{ apiStatusLabel }}
          </div>
        </div>
      </div>
    </header>

    <section class="border-2 border-black bg-white p-6 shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] mb-8">
      <div class="flex flex-col gap-6 lg:flex-row lg:items-center lg:justify-between mb-6 border-b-2 border-black pb-6">
        <div class="flex flex-1 flex-col gap-3">
          <p class="text-xs font-bold uppercase tracking-widest text-gray-500">Datasets</p>
          <SamplePicker variant="inline" @load="onSample" />
        </div>
        <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
          <button
            class="inline-flex items-center justify-center border-2 border-black bg-black px-6 py-3 text-sm font-bold uppercase tracking-widest text-white transition-all hover:bg-white hover:text-black shadow-[4px_4px_0px_0px_rgba(0,0,0,0)] hover:shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:shadow-none disabled:hover:bg-black disabled:hover:text-white"
            :disabled="loading"
            @click="runDetect"
          >
            {{ loading ? 'Scanning…' : 'Run Analyzer' }}
          </button>
          <div class="border-2 border-black px-4 py-2 text-xs font-bold uppercase tracking-widest text-black bg-gray-100">
            Status · {{ loaderLabel }}
          </div>
        </div>
      </div>

      <div class="grid gap-6 lg:grid-cols-3">
        <div class="lg:col-span-2 space-y-6">
          <UploadPanel @data="onData" />
        </div>
        <div class="space-y-6">
          <ParamsForm :params="params" @update="onParams" @run="runDetect" />
          <div class="border-2 border-black bg-gray-50 p-4">
            <p class="text-xs font-bold uppercase tracking-widest text-gray-500 mb-2">API Health</p>
            <p class="text-lg font-mono font-bold text-black">{{ apiStatusLabel }}</p>
            <p class="text-xs font-mono text-gray-500 mt-1">Last check: {{ lastHealthCheckCopy }}</p>
          </div>
        </div>
      </div>
    </section>

    <section v-if="loading" class="mt-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-4 mb-8">
      <div
        v-for="(step, idx) in loaderStates"
        :key="idx"
        :class="[
          'border-2 px-4 py-4 text-sm font-bold uppercase tracking-widest transition-colors',
          step.state === 'done' ? 'border-black bg-black text-white' :
          step.state === 'active' ? 'border-black bg-white text-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]' :
          'border-gray-300 bg-white text-gray-400'
        ]"
      >
        {{ step.label }}
      </div>
    </section>

    <section class="mt-10 mb-10 border-2 border-black p-4 bg-white shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
      <AnomalyChart :series="chartSeries" :threshold="0.5" @select="handleChartSelect" />
    </section>

    <section
      v-if="response"
      class="mt-10 border-2 border-black bg-white p-6 shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]"
    >
      <div class="mb-6 flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between border-b-2 border-black pb-6">
        <div>
          <h3 class="text-2xl font-sans font-black uppercase tracking-tighter text-black">Detection Results</h3>
          <p class="text-sm font-serif text-gray-800">Deterministic compression metrics ready for auditors.</p>
        </div>
        <CurlSnippet :curl="curlCmd" />
      </div>

      <div class="mb-6 flex flex-wrap gap-6">
        <div class="border-2 border-black px-4 py-3 bg-gray-50 min-w-[140px]">
          <p class="text-xs font-bold uppercase tracking-widest text-gray-500 mb-1">Events</p>
          <p class="text-2xl font-mono font-bold text-black">{{ response.total_events }}</p>
        </div>
        <div class="border-2 border-black px-4 py-3 bg-gray-50 min-w-[140px]">
          <p class="text-xs font-bold uppercase tracking-widest text-gray-500 mb-1">Anomalies</p>
          <p :class="['text-2xl font-mono font-bold', (response.anomaly_count || 0) > 0 ? 'text-red-600' : 'text-green-600']">
            {{ response.anomaly_count || 0 }}
          </p>
        </div>
        <div class="border-2 border-black px-4 py-3 bg-gray-50 min-w-[140px]">
          <p class="text-xs font-bold uppercase tracking-widest text-gray-500 mb-1">Algorithm</p>
          <p class="text-2xl font-mono font-bold text-black">{{ response.compression_algo }}</p>
        </div>
        <div class="border-2 border-black px-4 py-3 bg-gray-50 min-w-[140px]">
          <p class="text-xs font-bold uppercase tracking-widest text-gray-500 mb-1">Latency</p>
          <p class="text-2xl font-mono font-bold text-black">{{ response.processing_time }}</p>
        </div>
      </div>

      <div class="grid gap-6 lg:grid-cols-3">
        <div class="lg:col-span-2 border-2 border-black bg-white p-4">
          <ResultsTable
            :items="response.anomalies || []"
            :selected-index="selectedIndex"
            @select="selectFromTable"
          />
          <div class="mt-4 flex justify-end">
            <button
              class="inline-flex items-center gap-2 border-2 border-black px-4 py-2 text-sm font-bold uppercase tracking-widest text-black hover:bg-black hover:text-white transition-colors"
              :disabled="!response"
              @click="downloadJSON"
            >
              Download JSON
            </button>
          </div>
        </div>
        <AnomalyDetail :anomaly="selectedAnomaly" :loading="loading" />
      </div>
    </section>

    <section
      v-if="error"
      class="mt-8 border-2 border-red-600 bg-red-50 p-4 text-sm font-bold text-red-600 shadow-[4px_4px_0px_0px_rgba(220,38,38,1)]"
    >
      <div class="flex items-center gap-2">
        <span>⚠</span>
        <span>{{ error }}</span>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import UploadPanel from './UploadPanel.vue'
import ParamsForm from './ParamsForm.vue'
import ResultsTable from './ResultsTable.vue'
import SamplePicker from './SamplePicker.vue'
import CurlSnippet from './CurlSnippet.vue'
import AnomalyChart from './AnomalyChart.vue'
import AnomalyDetail from './AnomalyDetail.vue'

const props = withDefaults(defineProps<{
  variant?: 'full' | 'embedded'
}>(), {
  variant: 'full',
})

const apiBase = import.meta.env.VITE_API_BASE_URL || 'https://driftlock.net/api/v1'

const raw = ref<string>('')
const mime = ref<'ndjson' | 'json'>('ndjson')
const params = ref({ baseline: 400, window: 1, hop: 1, algo: 'zstd' })
const response = ref<any | null>(null)
const error = ref<string | null>(null)
const loading = ref<boolean>(false)
const apiStatus = ref<'checking' | 'connected' | 'disconnected'>('checking')
const lastHealthCheck = ref<Date | null>(null)
const chartSeries = ref<SeriesPoint[]>([])
const selectedAnomaly = ref<any | null>(null)
const loaderIndex = ref(0)
let loaderTimer: number | null = null
let healthInterval: number | null = null

const rootClasses = computed(() => {
  // Brutalist container style
  return 'bg-white text-black'
})

const loaderSteps = [
  'init core',
  'calibrate baseline',
  'compress + compare',
  'explain anomaly',
]

const loaderDurations = [400, 800, 900, 600]

const loaderStates = computed(() => loaderSteps.map((label, idx) => ({
  label,
  state: idx < loaderIndex.value ? 'done' : idx === loaderIndex.value ? 'active' : 'pending',
})))

const loaderLabel = computed(() => loaderSteps[Math.min(loaderIndex.value, loaderSteps.length - 1)])

const lastHealthCheckCopy = computed(() => {
  if (!lastHealthCheck.value) return 'pending…'
  return lastHealthCheck.value.toLocaleTimeString()
})

const apiStatusLabel = computed(() => {
  if (apiStatus.value === 'connected') return 'API Connected'
  if (apiStatus.value === 'checking') return 'Checking…'
  return 'API Unavailable'
})

const selectedIndex = computed(() => {
  if (!selectedAnomaly.value) return null
  return selectedAnomaly.value.index
    ?? selectedAnomaly.value.sequence
    ?? selectedAnomaly.value.window_index
    ?? null
})

async function checkApiHealth() {
  apiStatus.value = 'checking'
  try {
    const controller = new AbortController()
    const timeoutId = window.setTimeout(() => controller.abort(), 5000)
    const res = await fetch(`${apiBase}/healthz`, {
      method: 'GET',
      signal: controller.signal,
    })
    window.clearTimeout(timeoutId)
    lastHealthCheck.value = new Date()
    if (res.ok) {
      const data = await res.json()
      if (data.success === true || data.ok === true) {
        apiStatus.value = 'connected'
        return
      }
    }
    apiStatus.value = 'disconnected'
  } catch {
    apiStatus.value = 'disconnected'
  }
}

onMounted(() => {
  checkApiHealth()
  healthInterval = window.setInterval(checkApiHealth, 30000)
})

onUnmounted(() => {
  if (healthInterval) {
    clearInterval(healthInterval)
  }
  stopLoaderSequence()
})

watch(response, (payload) => {
  chartSeries.value = buildSeries(payload)
  selectedAnomaly.value = payload?.anomalies?.[0] ?? null
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

function startLoaderSequence() {
  stopLoaderSequence()
  loaderIndex.value = 0
  const advance = (step: number) => {
    loaderTimer = window.setTimeout(() => {
      if (loaderIndex.value < loaderSteps.length - 1) {
        loaderIndex.value += 1
        advance(step + 1)
      }
    }, loaderDurations[Math.min(step, loaderDurations.length - 1)])
  }
  advance(0)
}

function stopLoaderSequence() {
  if (loaderTimer) {
    clearTimeout(loaderTimer)
    loaderTimer = null
  }
  loaderIndex.value = loaderSteps.length - 1
}

async function runDetect() {
  if (!raw.value.trim()) {
    error.value = 'Please load or paste data before running detection.'
    return
  }
  error.value = null
  response.value = null
  loading.value = true
  startLoaderSequence()

  let status = apiStatus.value
  if (status !== 'connected') {
    await checkApiHealth()
    status = apiStatus.value
    if (status !== 'connected') {
      error.value = 'API is unavailable. Please ensure the Driftlock API server is reachable.'
      loading.value = false
      stopLoaderSequence()
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
    if (e?.name === 'AbortError' || e?.name === 'TimeoutError') {
      error.value = 'Request timed out. The API may be slow or unavailable.'
    } else if (e?.message?.includes('Failed to fetch') || e?.message?.includes('NetworkError')) {
      error.value = `Network error. Please ensure the API server is running at ${apiBase}`
      apiStatus.value = 'disconnected'
    } else {
      error.value = e?.message || String(e)
    }
  } finally {
    loading.value = false
    stopLoaderSequence()
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

function buildSeries(payload: any | null): SeriesPoint[] {
  const total = Number(payload?.total_events) || (payload?.anomalies?.length ?? 0)
  if (!total) return []
  const baseSeries: SeriesPoint[] = Array.from({ length: total }).map((_, idx) => {
    const signal = 0.22 + Math.sin((idx + 1) / 3.5) * 0.08 + ((idx + 1) % 7 === 0 ? 0.05 : 0)
    const clamped = Math.min(Math.max(signal, 0.05), 0.7)
    return { index: idx + 1, score: Number(clamped.toFixed(3)), anomaly: false }
  })

  const anomalies = payload?.anomalies ?? []
  anomalies.forEach((entry: any) => {
    const rawIndex = Number(entry.index ?? entry.sequence ?? entry.window_index ?? entry.position ?? 1)
    const slot = Math.min(Math.max(rawIndex || 1, 1), baseSeries.length)
    baseSeries[slot - 1] = {
      index: slot,
      score: Number((entry.metrics?.ncd ?? 0.82).toFixed(3)),
      anomaly: true,
      record: entry,
    }
  })

  return baseSeries
}

function handleChartSelect(payload: { index: number, record?: any }) {
  selectedAnomaly.value = payload.record ?? findAnomalyByIndex(payload.index)
}

function selectFromTable(row: any) {
  selectedAnomaly.value = row
}

function findAnomalyByIndex(index: number) {
  return response.value?.anomalies?.find((item: any) => {
    const idx = item.index ?? item.sequence ?? item.window_index ?? item.position
    return idx === index
  }) ?? null
}

async function runFinancialDemo() {
  await onSample('/samples/demo-financial.ndjson')
  await nextTick()
  runDetect()
}

defineExpose({ runFinancialDemo })

interface SeriesPoint {
  index: number
  score: number
  anomaly?: boolean
  record?: any
}
</script>

