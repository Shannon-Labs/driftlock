<template>
  <div class="min-h-screen bg-white">
    <DashboardLayout>
      <div class="px-4 sm:px-6 lg:px-8 py-8 space-y-8">
        <div class="flex flex-col gap-3 border-b-4 border-black pb-6">
          <div class="flex items-center gap-3">
            <div class="h-10 w-10 flex items-center justify-center border-2 border-black bg-yellow-200 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]">
              <svg class="h-6 w-6 text-black" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7h16M4 12h16m-7 5h7" />
              </svg>
            </div>
            <div>
              <h1 class="text-3xl font-black uppercase tracking-tight text-black">Upload & Analyze</h1>
              <p class="text-sm font-mono text-gray-600">Upload JSON/NDJSON or paste events, then run /v1/detect (max 10 MB / 10k events).</p>
            </div>
          </div>
          <div class="flex flex-wrap gap-2 text-xs font-mono text-gray-600">
            <span class="px-2 py-1 border border-black bg-gray-100 uppercase">Auth required</span>
            <span class="px-2 py-1 border border-black bg-gray-100 uppercase">10 MB / 10k events</span>
            <span class="px-2 py-1 border border-black bg-gray-100 uppercase">Adds idempotency keys</span>
          </div>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <!-- Form -->
          <div class="lg:col-span-2 space-y-6">
            <div class="border-2 border-black p-6 shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] bg-white">
              <h2 class="text-xl font-bold uppercase tracking-wide text-black mb-4">Configuration</h2>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label class="block text-xs font-bold uppercase tracking-widest text-gray-500 mb-2">API Key (required)</label>
                  <input
                    v-model="apiKey"
                    type="password"
                    placeholder="dlk_live_..."
                    class="w-full border-2 border-black px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-black focus:ring-offset-2"
                  />
                  <p class="mt-1 text-[11px] text-gray-500 font-mono">Not stored on server. Kept locally for this session.</p>
                </div>
                <div>
                  <label class="block text-xs font-bold uppercase tracking-widest text-gray-500 mb-2">Stream ID</label>
                  <input
                    v-model="streamId"
                    type="text"
                    placeholder="default"
                    class="w-full border-2 border-black px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-black focus:ring-offset-2"
                  />
                  <p class="mt-1 text-[11px] text-gray-500 font-mono">Leave blank for default stream.</p>
                </div>
                <div>
                  <label class="block text-xs font-bold uppercase tracking-widest text-gray-500 mb-2">Sensitivity</label>
                  <select
                    v-model="sensitivity"
                    class="w-full border-2 border-black px-3 py-2 text-sm font-bold uppercase tracking-wide focus:outline-none focus:ring-2 focus:ring-black focus:ring-offset-2"
                  >
                    <option value="balanced">Balanced (default)</option>
                    <option value="sensitive">High (more anomalies)</option>
                    <option value="strict">Low (fewer false positives)</option>
                  </select>
                  <p class="mt-1 text-[11px] text-gray-500 font-mono">
                    Sends config_override when not balanced.
                  </p>
                </div>
                <div class="flex items-center gap-3 mt-6">
                  <input id="idempotency" type="checkbox" v-model="addIdempotency" class="h-4 w-4 border-2 border-black text-black" />
                  <label for="idempotency" class="text-xs font-bold uppercase tracking-widest text-gray-700">Add idempotency keys to events</label>
                </div>
              </div>
            </div>

            <div class="border-2 border-black p-6 shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] bg-white space-y-4">
              <div class="flex items-center justify-between">
                <div>
                  <h2 class="text-xl font-bold uppercase tracking-wide text-black">Upload or Paste Events</h2>
                  <p class="text-sm font-mono text-gray-600">JSON array, { events: [...] }, or NDJSON. Max 10 MB / 10k events.</p>
                </div>
                <div class="text-xs font-mono text-gray-600">
                  <p>Selected: {{ fileName || 'None' }}</p>
                  <p v-if="fileSize">Size: {{ formatBytes(fileSize) }}</p>
                </div>
              </div>

              <div class="flex flex-wrap gap-3">
                <label class="inline-flex items-center px-4 py-2 border-2 border-black bg-white text-sm font-bold uppercase tracking-widest cursor-pointer hover:bg-black hover:text-white transition-colors">
                  <input type="file" class="hidden" accept=".json,.ndjson,application/json" @change="onFileChange" />
                  Upload File
                </label>
                <button
                  type="button"
                  @click="clearFile"
                  class="inline-flex items-center px-4 py-2 border-2 border-black bg-gray-100 text-sm font-bold uppercase tracking-widest hover:bg-white transition-colors"
                  :disabled="!fileContent && !fileName"
                >
                  Clear File
                </button>
              </div>

              <textarea
                v-model="eventsText"
                placeholder='[{"timestamp":"2025-01-01T10:00:00Z","type":"metric","body":{"latency_ms":120}}]'
                class="w-full border-2 border-black px-3 py-2 text-sm font-mono min-h-[180px] focus:outline-none focus:ring-2 focus:ring-black focus:ring-offset-2"
              ></textarea>

              <div class="flex flex-wrap gap-3">
                <button
                  @click="runDetect"
                  :disabled="loading"
                  class="inline-flex items-center justify-center px-5 py-3 border-2 border-black bg-black text-white text-sm font-bold uppercase tracking-widest hover:bg-white hover:text-black transition-colors disabled:opacity-50"
                >
                  <svg v-if="loading" class="h-4 w-4 mr-2 animate-spin" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke-width="4"></circle>
                    <path class="opacity-75" stroke-linecap="round" stroke-linejoin="round" stroke-width="4" d="M4 12a8 8 0 018-8"></path>
                  </svg>
                  {{ loading ? 'Analyzing...' : 'Run Detection' }}
                </button>
                <button
                  type="button"
                  @click="prefillSample"
                  class="inline-flex items-center justify-center px-4 py-3 border-2 border-black bg-white text-sm font-bold uppercase tracking-widest hover:bg-black hover:text-white transition-colors"
                >
                  Load Sample Events
                </button>
              </div>

              <div v-if="error" class="border-2 border-red-600 bg-red-50 text-red-800 px-4 py-3 text-sm font-mono">
                {{ error }}
              </div>
              <div v-if="info" class="border-2 border-yellow-500 bg-yellow-50 text-yellow-800 px-4 py-3 text-sm font-mono">
                {{ info }}
              </div>
            </div>
          </div>

          <!-- Side panel -->
          <div class="space-y-4">
            <div class="border-2 border-black p-5 bg-gray-50 shadow-[6px_6px_0px_0px_rgba(0,0,0,1)]">
              <h3 class="text-lg font-bold uppercase tracking-wide text-black mb-2">Limits</h3>
              <ul class="text-sm font-mono text-gray-700 space-y-1">
                <li>• Max size: 10 MB</li>
                <li>• Max events: 10,000</li>
                <li>• Accepted: JSON array, {"events":[]}, NDJSON</li>
                <li>• Auth: X-Api-Key required</li>
              </ul>
            </div>

            <div class="border-2 border-black p-5 bg-white shadow-[6px_6px_0px_0px_rgba(0,0,0,1)]">
              <h3 class="text-lg font-bold uppercase tracking-wide text-black mb-2">cURL Example</h3>
              <pre class="bg-black text-green-400 text-xs font-mono p-3 overflow-x-auto">
curl -X POST https://api.driftlock.net/v1/detect \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: {{ apiKey || 'YOUR_API_KEY' }}" \
  -d '{ "stream_id": "{{ streamId || 'default' }}", "events": [...] }'
              </pre>
            </div>

            <div v-if="recentRuns.length" class="border-2 border-black p-5 bg-white shadow-[6px_6px_0px_0px_rgba(0,0,0,1)]">
              <h3 class="text-lg font-bold uppercase tracking-wide text-black mb-3">Recent Runs</h3>
              <ul class="space-y-2 text-sm font-mono text-gray-700">
                <li v-for="run in recentRuns" :key="run.id" class="flex items-center justify-between">
                  <div>
                    <p class="font-bold text-black">{{ run.stream || 'default' }} • {{ run.anomaly_count }} anomalies</p>
                    <p class="text-[11px] text-gray-500">{{ run.events }} events • {{ new Date(run.at).toLocaleString() }}</p>
                  </div>
                  <span class="text-[11px] uppercase px-2 py-1 border border-black bg-gray-100">{{ run.sensitivity }}</span>
                </li>
              </ul>
            </div>
          </div>
        </div>

        <!-- Results -->
        <div v-if="response || loading" class="border-2 border-black bg-white shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] p-6">
          <div class="flex items-center justify-between gap-3 flex-wrap border-b-2 border-black pb-4 mb-4">
            <div>
              <h2 class="text-xl font-bold uppercase tracking-wide text-black">Results</h2>
              <p class="text-sm font-mono text-gray-600">
                {{ response ? `${response.anomaly_count || 0} anomalies • ${response.total_events || 0} events` : 'Awaiting results' }}
              </p>
            </div>
            <div class="flex gap-2">
              <button
                :disabled="!response"
                @click="downloadJson"
                class="px-3 py-2 border-2 border-black bg-white text-xs font-bold uppercase tracking-widest hover:bg-black hover:text-white transition-colors disabled:opacity-50"
              >
                Download JSON
              </button>
              <button
                :disabled="!response"
                @click="copyResponseCurl"
                class="px-3 py-2 border-2 border-black bg-white text-xs font-bold uppercase tracking-widest hover:bg-black hover:text-white transition-colors disabled:opacity-50"
              >
                Copy cURL
              </button>
            </div>
          </div>

          <div v-if="loading" class="flex items-center gap-3 text-sm font-mono text-gray-600">
            <div class="h-4 w-4 border-2 border-black border-t-transparent rounded-full animate-spin"></div>
            Running detection...
          </div>

          <div v-else-if="response">
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
              <div class="border-2 border-black bg-gray-50 p-3">
                <p class="text-[11px] font-bold uppercase tracking-widest text-gray-500">Anomalies</p>
                <p class="text-2xl font-black text-black">{{ response.anomaly_count ?? 0 }}</p>
              </div>
              <div class="border-2 border-black bg-gray-50 p-3">
                <p class="text-[11px] font-bold uppercase tracking-widest text-gray-500">Total Events</p>
                <p class="text-2xl font-black text-black">{{ response.total_events ?? 0 }}</p>
              </div>
              <div class="border-2 border-black bg-gray-50 p-3">
                <p class="text-[11px] font-bold uppercase tracking-widest text-gray-500">Processing Time</p>
                <p class="text-2xl font-black text-black">{{ response.processing_time || '—' }}</p>
              </div>
              <div class="border-2 border-black bg-gray-50 p-3">
                <p class="text-[11px] font-bold uppercase tracking-widest text-gray-500">Compression</p>
                <p class="text-xl font-black text-black">{{ response.compression_algo || '—' }}</p>
                <p class="text-[11px] text-gray-600" v-if="response.fallback_from_algo">Fallback from {{ response.fallback_from_algo }}</p>
              </div>
            </div>

            <div v-if="(response.anomalies || []).length === 0" class="border-2 border-black bg-green-50 text-green-800 px-4 py-3 font-mono">
              No anomalies detected.
            </div>

            <div v-else class="overflow-x-auto">
              <table class="min-w-full divide-y-2 divide-black">
                <thead class="bg-gray-50">
                  <tr>
                    <th class="px-3 py-2 text-left text-xs font-bold uppercase tracking-widest text-gray-500">Index</th>
                    <th class="px-3 py-2 text-left text-xs font-bold uppercase tracking-widest text-gray-500">NCD</th>
                    <th class="px-3 py-2 text-left text-xs font-bold uppercase tracking-widest text-gray-500">p-value</th>
                    <th class="px-3 py-2 text-left text-xs font-bold uppercase tracking-widest text-gray-500">Confidence</th>
                    <th class="px-3 py-2 text-left text-xs font-bold uppercase tracking-widest text-gray-500">Why</th>
                    <th class="px-3 py-2 text-left text-xs font-bold uppercase tracking-widest text-gray-500">Event</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-200 bg-white">
                  <tr v-for="anom in response.anomalies" :key="anom.id || anom.index" class="hover:bg-gray-50">
                    <td class="px-3 py-2 text-sm font-mono text-gray-700">{{ anom.index ?? '—' }}</td>
                    <td class="px-3 py-2 text-sm font-mono text-gray-700">{{ formatNumber(anom.metrics?.ncd) }}</td>
                    <td class="px-3 py-2 text-sm font-mono text-gray-700">{{ formatNumber(anom.metrics?.p_value) }}</td>
                    <td class="px-3 py-2 text-sm font-mono text-gray-700">{{ formatPercent(anom.metrics?.confidence) }}</td>
                    <td class="px-3 py-2 text-sm text-gray-700 max-w-xs truncate">{{ anom.why || '—' }}</td>
                    <td class="px-3 py-2 text-sm font-mono text-gray-700 max-w-xs truncate">
                      {{ previewEvent(anom.event || anom.body) }}
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </DashboardLayout>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DashboardLayout from '../layouts/DashboardLayout.vue'
import { buildDetectPayload, parseEventsInput, PayloadError } from '../utils/payload'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()

const apiKey = ref(localStorage.getItem('analyze_api_key') || '')
const streamId = ref(localStorage.getItem('analyze_stream_id') || 'default')
const sensitivity = ref<'sensitive' | 'balanced' | 'strict'>(
  (localStorage.getItem('analyze_sensitivity') as 'sensitive' | 'balanced' | 'strict') || 'balanced'
)
const addIdempotency = ref(true)

const eventsText = ref('')
const fileContent = ref('')
const fileName = ref('')
const fileSize = ref<number | null>(null)

const loading = ref(false)
const error = ref('')
const info = ref('')
const response = ref<any | null>(null)

interface RunSummary {
  id: string
  at: string
  stream: string
  anomaly_count: number
  events: number
  sensitivity: string
}
const recentRuns = ref<RunSummary[]>([])

const MAX_SIZE_BYTES = 10 * 1024 * 1024
const MAX_EVENTS = 10_000

const formatBytes = (bytes: number) => {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${units[i]}`
}

const formatNumber = (val?: number) => (typeof val === 'number' ? val.toFixed(3) : '—')
const formatPercent = (val?: number) => (typeof val === 'number' ? `${(val * 100).toFixed(1)}%` : '—')

const previewEvent = (evt: unknown) => {
  try {
    return JSON.stringify(evt)?.slice(0, 120) || ''
  } catch {
    return ''
  }
}

const persistPrefs = () => {
  localStorage.setItem('analyze_api_key', apiKey.value)
  localStorage.setItem('analyze_stream_id', streamId.value)
  localStorage.setItem('analyze_sensitivity', sensitivity.value)
}

const loadHistory = () => {
  try {
    const stored = localStorage.getItem('analyze_history')
    if (stored) {
      recentRuns.value = JSON.parse(stored)
    }
  } catch {
    recentRuns.value = []
  }
}

const saveHistory = (run: RunSummary) => {
  const next = [run, ...recentRuns.value].slice(0, 5)
  recentRuns.value = next
  localStorage.setItem('analyze_history', JSON.stringify(next))
}

const validateSize = (payload: string) => {
  const bytes = new TextEncoder().encode(payload).length
  if (bytes > MAX_SIZE_BYTES) {
    throw new PayloadError(`Payload too large (${formatBytes(bytes)}). Limit is 10 MB.`)
  }
}

const onFileChange = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  error.value = ''
  info.value = ''
  fileContent.value = ''

  if (!file) return
  if (file.size > MAX_SIZE_BYTES) {
    error.value = `File too large (${formatBytes(file.size)}). Max 10 MB.`
    target.value = ''
    return
  }

  fileName.value = file.name
  fileSize.value = file.size
  try {
    const text = await file.text()
    fileContent.value = text
    info.value = `Loaded ${file.name}`
  } catch {
    error.value = 'Failed to read file.'
  }
}

const clearFile = () => {
  fileName.value = ''
  fileSize.value = null
  fileContent.value = ''
}

const pickInputText = () => {
  if (fileContent.value) return fileContent.value
  return eventsText.value
}

const runDetect = async () => {
  error.value = ''
  info.value = ''
  response.value = null

  if (!apiKey.value.trim()) {
    error.value = 'API key is required.'
    return
  }

  const payloadText = pickInputText()
  if (!payloadText.trim()) {
    error.value = 'Provide events via upload or paste.'
    return
  }

  try {
    validateSize(payloadText)
    const { events } = parseEventsInput(payloadText, MAX_EVENTS)
    const payload = buildDetectPayload(events, streamId.value || 'default', sensitivity.value, addIdempotency.value)

    loading.value = true
    const res = await fetch('/api/v1/detect', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Api-Key': apiKey.value.trim(),
      },
      body: JSON.stringify(payload),
    })

    const data = await res.json().catch(() => null)
    if (!res.ok) {
      const msg =
        data?.error?.message ||
        (res.status === 429 ? 'Rate limit exceeded. Please retry shortly.' : 'Detection failed. Check your payload.')
      throw new Error(msg)
    }

    response.value = data
    persistPrefs()
    saveHistory({
      id: data.batch_id || crypto.randomUUID(),
      at: new Date().toISOString(),
      stream: data.stream_id || streamId.value || 'default',
      anomaly_count: data.anomaly_count || 0,
      events: data.total_events || events.length,
      sensitivity: sensitivity.value,
    })
  } catch (e: any) {
    error.value = e instanceof PayloadError ? e.message : e?.message || 'Something went wrong.'
  } finally {
    loading.value = false
  }
}

const downloadJson = () => {
  if (!response.value) return
  const blob = new Blob([JSON.stringify(response.value, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'detect-response.json'
  a.click()
  URL.revokeObjectURL(url)
}

const copyResponseCurl = async () => {
  if (!response.value) return
  const curl = [
    'curl -X POST https://api.driftlock.net/v1/detect',
    `  -H "Content-Type: application/json"`,
    `  -H "X-Api-Key: ${apiKey.value || 'YOUR_API_KEY'}"`,
    `  -d '{ "stream_id": "${streamId.value || 'default'}", "events": [...] }'`,
  ].join(' \\\n')
  await navigator.clipboard.writeText(curl)
  info.value = 'Copied cURL to clipboard'
}

const prefillSample = () => {
  const sample = JSON.stringify(
    [
      { timestamp: '2025-01-01T10:00:00Z', type: 'metric', body: { latency_ms: 120, cpu: 45 } },
      { timestamp: '2025-01-01T10:01:00Z', type: 'metric', body: { latency_ms: 118, cpu: 44 } },
      { timestamp: '2025-01-01T10:02:00Z', type: 'metric', body: { latency_ms: 950, cpu: 92 } },
    ],
    null,
    2
  )
  eventsText.value = sample
  clearFile()
}

onMounted(() => {
  if (authStore.loading) {
    authStore.init()
  }
  loadHistory()
})
</script>
