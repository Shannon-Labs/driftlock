<template>
  <div class="rounded-3xl border border-slate-800 bg-slate-950/95 p-6 text-white shadow-[0_40px_120px_rgba(15,23,42,0.65)]">
    <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <p class="text-xs font-semibold uppercase tracking-[0.3em] text-cyan-300/80">
          Glass-box inspector
        </p>
        <h3 class="text-xl font-semibold">Anomaly Evidence</h3>
        <p class="text-sm text-slate-400">
          Raw payload + deterministic math for auditors and LLM ops.
        </p>
      </div>
      <div v-if="anomaly" class="inline-flex items-center gap-2 rounded-2xl border border-white/10 bg-white/5 px-4 py-2 text-xs font-semibold uppercase tracking-[0.3em]">
        <span class="h-1.5 w-1.5 animate-pulse rounded-full bg-rose-400"></span>
        {{ displayIndex }}
      </div>
      <div v-else-if="loading" class="inline-flex items-center gap-2 rounded-2xl border border-white/10 bg-white/5 px-4 py-2 text-xs font-semibold uppercase tracking-[0.3em] text-slate-300">
        <span class="h-1.5 w-1.5 animate-spin rounded-full border border-cyan-200 border-t-transparent"></span>
        Listening
      </div>
    </div>

    <div v-if="anomaly" class="mt-6 space-y-6">
      <div class="grid gap-4">
        <div class="rounded-2xl border border-slate-800 bg-slate-900/60 p-4">
          <div class="flex items-center justify-between">
            <h4 class="text-sm font-semibold uppercase tracking-[0.4em] text-slate-400">Raw Event</h4>
            <span class="rounded-full border border-cyan-400/40 px-3 py-1 text-xs font-semibold text-cyan-200">
              immutable
            </span>
          </div>
          <pre class="mt-3 overflow-x-auto rounded-xl bg-black/40 p-4 text-xs leading-relaxed text-cyan-100">
{{ prettyEvent }}
          </pre>
        </div>

        <div class="space-y-4">
          <div class="rounded-2xl border border-slate-800 bg-slate-900/60 p-4">
            <div class="flex items-center justify-between">
              <h4 class="text-sm font-semibold uppercase tracking-[0.4em] text-slate-400">Math Breakdown</h4>
              <span class="rounded-full bg-emerald-400/10 px-3 py-1 text-xs font-semibold text-emerald-300">
                deterministic
              </span>
            </div>
            <dl class="mt-4 grid grid-cols-2 gap-4 text-sm">
              <div>
                <dt class="text-slate-400">NCD</dt>
                <dd class="text-lg font-mono text-white">{{ fmt(metrics.ncd) }}</dd>
              </div>
              <div>
                <dt class="text-slate-400">p-value</dt>
                <dd class="text-lg font-mono text-white">{{ fmt(metrics.p_value) }}</dd>
              </div>
              <div>
                <dt class="text-slate-400">Confidence</dt>
                <dd class="text-lg font-mono text-cyan-300">{{ pct(metrics.confidence_level) }}</dd>
              </div>
              <div>
                <dt class="text-slate-400">Compression Δ</dt>
                <dd class="text-lg font-mono text-rose-300">
                  {{ compressionDelta }}
                </dd>
              </div>
            </dl>
          </div>

          <div class="rounded-2xl border border-slate-800 bg-gradient-to-br from-slate-900/80 via-slate-900/60 to-slate-950/60 p-4">
            <div class="flex items-center justify-between">
              <h4 class="text-sm font-semibold uppercase tracking-[0.4em] text-slate-400">Narrative</h4>
              <span class="rounded-full border border-fuchsia-400/40 px-3 py-1 text-xs font-semibold text-fuchsia-200">
                AI optional
              </span>
            </div>
            <p class="mt-3 text-sm leading-relaxed text-slate-200">
              {{ explanation }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="mt-10 flex flex-col items-center justify-center text-center text-slate-400">
      <span class="text-4xl">⌁</span>
      <p class="mt-3 text-sm uppercase tracking-[0.3em]">Run analyzer to inspect anomalies</p>
      <p class="text-xs text-slate-500">Mathematical context appears automatically when Driftlock flags an event.</p>
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


