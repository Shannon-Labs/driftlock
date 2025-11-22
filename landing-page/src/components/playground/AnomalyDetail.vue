<template>
  <div class="border-2 border-black bg-white p-6 shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
    <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between border-b-2 border-black pb-4 mb-6">
      <div>
        <p class="text-xs font-bold uppercase tracking-widest text-gray-500 mb-1">
          Glass-box inspector
        </p>
        <h3 class="text-xl font-sans font-black uppercase tracking-tighter text-black">Anomaly Evidence</h3>
        <p class="text-sm font-serif text-gray-800">
          Raw payload + deterministic math for auditors and LLM ops.
        </p>
      </div>
      <div v-if="anomaly" class="inline-flex items-center gap-2 border-2 border-black bg-black px-4 py-2 text-xs font-bold uppercase tracking-widest text-white">
        <span class="h-2 w-2 border border-white bg-red-500"></span>
        {{ displayIndex }}
      </div>
      <div v-else-if="loading" class="inline-flex items-center gap-2 border-2 border-black bg-white px-4 py-2 text-xs font-bold uppercase tracking-widest text-black">
        <span class="h-2 w-2 animate-spin border-2 border-black border-t-transparent"></span>
        Listening
      </div>
    </div>

    <div v-if="anomaly" class="space-y-6">
      <div class="grid gap-4">
        <div class="border-2 border-black bg-gray-50 p-4">
          <div class="flex items-center justify-between mb-3">
            <h4 class="text-sm font-bold uppercase tracking-widest text-gray-500">Raw Event</h4>
            <span class="border border-black px-3 py-1 text-xs font-bold uppercase tracking-wider bg-white">
              immutable
            </span>
          </div>
          <pre class="overflow-x-auto bg-black p-4 text-xs leading-relaxed text-green-400 font-mono border border-black">
{{ prettyEvent }}
          </pre>
        </div>

        <div class="space-y-4">
          <div class="border-2 border-black bg-white p-4">
            <div class="flex items-center justify-between mb-3">
              <h4 class="text-sm font-bold uppercase tracking-widest text-gray-500">Math Breakdown</h4>
              <span class="border border-black px-3 py-1 text-xs font-bold uppercase tracking-wider bg-green-100">
                deterministic
              </span>
            </div>
            <dl class="grid grid-cols-2 gap-4 text-sm">
              <div>
                <dt class="text-xs font-bold uppercase tracking-widest text-gray-500 mb-1">NCD</dt>
                <dd class="text-lg font-mono font-bold text-black">{{ fmt(metrics.ncd) }}</dd>
              </div>
              <div>
                <dt class="text-xs font-bold uppercase tracking-widest text-gray-500 mb-1">p-value</dt>
                <dd class="text-lg font-mono font-bold text-black">{{ fmt(metrics.p_value) }}</dd>
              </div>
              <div>
                <dt class="text-xs font-bold uppercase tracking-widest text-gray-500 mb-1">Confidence</dt>
                <dd class="text-lg font-mono font-bold text-black">{{ pct(metrics.confidence_level) }}</dd>
              </div>
              <div>
                <dt class="text-xs font-bold uppercase tracking-widest text-gray-500 mb-1">Compression Δ</dt>
                <dd class="text-lg font-mono font-bold text-red-600">
                  {{ compressionDelta }}
                </dd>
              </div>
            </dl>
          </div>

          <div class="border-2 border-black bg-gray-50 p-4">
            <div class="flex items-center justify-between mb-3">
              <h4 class="text-sm font-bold uppercase tracking-widest text-gray-500">Narrative</h4>
              <span class="border border-black px-3 py-1 text-xs font-bold uppercase tracking-wider bg-white">
                AI optional
              </span>
            </div>
            <p class="text-sm leading-relaxed font-serif text-black">
              {{ explanation }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="py-12 flex flex-col items-center justify-center text-center text-gray-500">
      <span class="text-4xl">⌁</span>
      <p class="mt-3 text-sm font-bold uppercase tracking-widest">Run analyzer to inspect anomalies</p>
      <p class="text-xs font-mono mt-2">Mathematical context appears automatically when Driftlock flags an event.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  anomaly: any | null
  loading?: boolean
}>()

const metrics = computed(() => props.anomaly?.metrics ?? {})

const prettyEvent = computed(() => {
  const payload = props.anomaly?.event ?? props.anomaly?.raw_event ?? props.anomaly?.record
  if (!payload) return '{\n  // awaiting anomaly payload\n}'
  try {
    return JSON.stringify(payload, null, 2)
  } catch {
    return String(payload)
  }
})

const compressionDelta = computed(() => {
  const baseline = metrics.value?.baseline_bytes ?? metrics.value?.baseline ?? null
  const windowBytes = metrics.value?.window_bytes ?? metrics.value?.compressed ?? null
  if (baseline == null || windowBytes == null) return '—'
  const delta = (windowBytes - baseline)
  return `${delta > 0 ? '+' : ''}${delta.toFixed(0)} bytes`
})

const explanation = computed(() => {
  return props.anomaly?.ai_explanation
    || props.anomaly?.why
    || 'Upgrade to Pro for Gemini Flash narratives, or rely on the deterministic math above.'
})

const displayIndex = computed(() => {
  const id = props.anomaly?.event_id ?? props.anomaly?.index ?? props.anomaly?.sequence ?? '—'
  return `event ${id}`
})

function fmt(input: number | undefined) {
  if (typeof input !== 'number') return '—'
  return input.toFixed(3)
}

function pct(input: number | undefined) {
  if (typeof input !== 'number') return '—'
  return `${(input * 100).toFixed(1)}%`
}
</script>


