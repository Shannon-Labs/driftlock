<template>
  <div :class="containerClasses">
    <div v-if="variant === 'panel'" class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Sample Data</h3>
      <span class="text-xs font-semibold uppercase tracking-[0.3em] text-gray-400 dark:text-gray-500">
        safe fixtures
      </span>
    </div>
    <div :class="chipWrapperClasses">
      <button
        v-for="sample in samples"
        :key="sample.name"
        :class="chipClasses(sample.url)"
        @click="selectSample(sample)"
      >
        <div class="flex items-center gap-2">
          <span
            v-if="sample.badge"
            class="rounded-full border px-2 py-0.5 text-[10px] font-semibold uppercase tracking-[0.3em]"
            :class="activeSample === sample.url
              ? 'border-cyan-300/60 text-cyan-200'
              : 'border-gray-400/60 text-gray-500 dark:text-gray-400'"
          >
            {{ sample.badge }}
          </span>
          <span class="text-sm font-semibold" :class="activeSample === sample.url ? activeTextClasses : baseTextClasses">
            {{ sample.name }}
          </span>
        </div>
        <p v-if="sample.description && variant === 'panel'" class="text-xs text-gray-500 dark:text-gray-400">
          {{ sample.description }}
        </p>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

const props = withDefaults(defineProps<{
  variant?: 'panel' | 'inline'
}>(), {
  variant: 'panel',
})

const emit = defineEmits<{ (e: 'load', payload: string): void }>()

const samples = [
  {
    name: 'Financial Live Demo',
    url: '/samples/demo-financial.ndjson',
    badge: 'NEW',
    description: '60 fintech events with 3 anomaly spikes',
  },
  {
    name: 'NDJSON (small)',
    url: '/samples/small.ndjson',
    badge: 'core math',
    description: 'Minimal dataset for quick testing',
  },
  {
    name: 'JSON Array (small)',
    url: '/samples/small.json',
    badge: 'array',
    description: 'JSON array for browser uploads',
  },
]

const activeSample = ref<string | null>(null)

const containerClasses = computed(() => props.variant === 'panel'
  ? 'rounded-2xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 p-6 shadow-sm'
  : 'flex flex-col gap-2'
)

const chipWrapperClasses = computed(() => props.variant === 'panel'
  ? 'flex flex-col gap-3'
  : 'flex flex-wrap gap-3'
)

const baseTextClasses = computed(() => props.variant === 'panel'
  ? 'text-gray-700 dark:text-gray-200'
  : 'text-gray-200'
)

const activeTextClasses = computed(() => props.variant === 'panel'
  ? 'text-blue-600'
  : 'text-cyan-300'
)

function chipClasses(url: string) {
  const isActive = activeSample.value === url
  const base = props.variant === 'panel'
    ? 'flex flex-col gap-1 rounded-xl border px-4 py-3 text-left transition-all duration-200'
    : 'inline-flex items-center gap-2 rounded-full border px-4 py-2 text-sm font-medium transition-all duration-200'

  const palette = props.variant === 'panel'
    ? isActive
      ? 'border-blue-400 bg-blue-50 dark:bg-blue-950/40 shadow-inner'
      : 'border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/40 hover:border-blue-400'
    : isActive
      ? 'border-cyan-300/60 bg-cyan-500/10 text-cyan-200 shadow-inner'
      : 'border-white/10 bg-white/5 text-white/80 hover:border-cyan-400/50'

  return `${base} ${palette}`
}

function selectSample(sample: { url: string }) {
  activeSample.value = sample.url
  emit('load', sample.url)
}
</script>

