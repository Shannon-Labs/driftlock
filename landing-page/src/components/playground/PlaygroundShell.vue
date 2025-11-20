<template>
  <div :class="rootClasses">
    <header class="mb-10">
      <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <p class="text-xs font-semibold uppercase tracking-[0.4em] text-cyan-200/80">
            Live Playground
          </p>
          <h2 class="text-3xl lg:text-4xl font-semibold text-white">
            CBAD Mission Control
          </h2>
          <p class="text-sm text-slate-300">
            Paste NDJSON, stream a synthetic burst, and watch Driftlock surface anomalies with audit-ready math.
          </p>
        </div>
        <div class="flex items-center gap-3">
          <div class="flex items-center gap-2 rounded-2xl border border-white/10 bg-white/5 px-4 py-2 text-sm font-semibold text-white">
            <span
              class="h-2.5 w-2.5 rounded-full"
              :class="{
                'bg-emerald-400 animate-pulse': apiStatus === 'connected',
                'bg-amber-300 animate-pulse': apiStatus === 'checking',
                'bg-rose-400': apiStatus === 'disconnected'
              }"
            ></span>
            {{ apiStatusLabel }}
          </div>
        </div>
      </div>
    </header>

    <section class="rounded-3xl border border-white/10 bg-white/5 p-6 text-white/80 shadow-[0_35px_120px_rgba(15,23,42,0.55)]">
      <div class="flex flex-col gap-6 lg:flex-row lg:items-center lg:justify-between">
        <div class="flex flex-1 flex-col gap-3">
          <p class="text-xs font-semibold uppercase tracking-[0.4em] text-slate-400">Datasets</p>
          <SamplePicker variant="inline" @load="onSample" />
        </div>
        <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
          <button
            class="inline-flex items-center justify-center rounded-full bg-gradient-to-r from-cyan-500 to-blue-600 px-6 py-3 text-sm font-semibold text-white shadow-lg shadow-cyan-500/30 transition hover:scale-105 disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="loading"
            @click="runDetect"
          >
            {{ loading ? 'Scanning…' : 'Run Analyzer' }}
          </button>
          <div class="rounded-2xl border border-white/10 px-4 py-2 text-xs uppercase tracking-[0.3em] text-slate-300">
            Status · {{ loaderLabel }}
          </div>
        </div>
      </div>

      <div class="mt-6 grid gap-6 lg:grid-cols-3">
        <div class="lg:col-span-2 space-y-6">
          <UploadPanel @data="onData" />
        </div>
        <div class="space-y-6">
          <ParamsForm :params="params" @update="onParams" @run="runDetect" />
          <div class="rounded-2xl border border-white/10 bg-slate-950/60 p-4">
            <p class="text-xs font-semibold uppercase tracking-[0.4em] text-slate-400">API Health</p>
            <p class="mt-2 text-lg font-semibold text-white">{{ apiStatusLabel }}</p>
            <p class="text-sm text-slate-400">Last check: {{ lastHealthCheckCopy }}</p>
          </div>
        </div>
      </div>
    </section>

    <section v-if="loading" class="mt-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <div
        v-for="(step, idx) in loaderStates"
        :key="idx"
        :class="[
          'rounded-2xl border px-4 py-4 text-sm uppercase tracking-[0.4em]',
          step.state === 'done' ? 'border-emerald-500/40 bg-emerald-500/10 text-emerald-200' :
          step.state === 'active' ? 'border-cyan-400/50 bg-cyan-400/10 text-cyan-200 animate-pulse' :
          'border-white/10 bg-white/5 text-slate-400'
        ]"
      >
        {{ step.label }}
      </div>
    </section>

    <section class="mt-10">
      <AnomalyChart :series="chartSeries" :threshold="0.5" @select="handleChartSelect" />
    </section>

    <section
      v-if="response"
      class="mt-10 rounded-3xl border border-white/10 bg-white/5 p-6 text-white shadow-[0_35px_120px_rgba(15,23,42,0.45)]"
    >
      <div class="mb-6 flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <h3 class="text-2xl font-semibold">Detection Results</h3>
          <p class="text-sm text-slate-300">Deterministic compression metrics ready for auditors.</p>
        </div>
        <CurlSnippet :curl="curlCmd" />
      </div>

      <div class="mb-6 flex flex-wrap gap-6 text-sm text-slate-200">
        <div class="rounded-2xl border border-white/10 px-4 py-3">
          <p class="text-xs uppercase tracking-[0.3em] text-slate-400">Events</p>
          <p class="text-2xl font-mono">{{ response.total_events }}</p>
        </div>
        <div class="rounded-2xl border border-white/10 px-4 py-3">
          <p class="text-xs uppercase tracking-[0.3em] text-slate-400">Anomalies</p>
          <p :class="['text-2xl font-mono', (response.anomaly_count || 0) > 0 ? 'text-rose-300' : 'text-emerald-300']">
            {{ response.anomaly_count || 0 }}
          </p>
        </div>
        <div class="rounded-2xl border border-white/10 px-4 py-3">
          <p class="text-xs uppercase tracking-[0.3em] text-slate-400">Algorithm</p>
          <p class="text-2xl font-mono">{{ response.compression_algo }}</p>
        </div>
        <div class="rounded-2xl border border-white/10 px-4 py-3">
          <p class="text-xs uppercase tracking-[0.3em] text-slate-400">Latency</p>
          <p class="text-2xl font-mono">{{ response.processing_time }}</p>
        </div>
      </div>

      <div class="grid gap-6 lg:grid-cols-3">
        <div class="lg:col-span-2 rounded-2xl border border-white/10 bg-slate-950/60 p-4">
          <ResultsTable
            :items="response.anomalies || []"
            :selected-index="selectedIndex"
            @select="selectFromTable"
          />
          <div class="mt-4 flex justify-end">
            <button
              class="inline-flex items-center gap-2 rounded-full border border-white/20 px-4 py-2 text-sm font-semibold text-white/90 transition hover:border-cyan-400 hover:text-white"
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
      class="mt-8 rounded-2xl border border-rose-500/40 bg-rose-500/10 p-4 text-sm text-rose-100"
    >
      <div class="flex items-center gap-2 font-semibold">
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
  const base = 'rounded-[32px] border border-white/10 bg-gradient-to-b from-slate-950 via-slate-950/90 to-slate-900 p-6 lg:p-10 text-white'
  return props.variant === 'embedded' ? `${base} backdrop-blur-xl shadow-[0_60px_140px_rgba(2,6,23,0.75)]` : base
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

