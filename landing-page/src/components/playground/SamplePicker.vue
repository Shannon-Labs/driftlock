<template>
  <div :class="containerClasses">
    <div v-if="variant === 'panel'" class="flex items-center justify-between mb-4 border-b-2 border-black pb-4">
      <h3 class="text-xl font-sans font-black uppercase tracking-tighter text-black">Sample Data</h3>
      <span class="text-xs font-bold uppercase tracking-widest text-gray-500">
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
            class="border border-black px-2 py-0.5 text-[10px] font-bold uppercase tracking-widest"
            :class="activeSample === sample.url
              ? 'bg-black text-white'
              : 'bg-white text-gray-600'"
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
  ? 'border-2 border-black bg-white p-6 shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]'
  : 'flex flex-col gap-2'
)

const chipWrapperClasses = computed(() => props.variant === 'panel'
  ? 'flex flex-col gap-3'
  : 'flex flex-wrap gap-3'
)

const baseTextClasses = computed(() => props.variant === 'panel'
  ? 'text-black'
  : 'text-black'
)

const activeTextClasses = computed(() => props.variant === 'panel'
  ? 'text-white'
  : 'text-white'
)

function chipClasses(url: string) {
  const isActive = activeSample.value === url
  const base = props.variant === 'panel'
    ? 'flex flex-col gap-1 border-2 px-4 py-3 text-left transition-colors cursor-pointer'
    : 'inline-flex items-center gap-2 border-2 px-4 py-2 text-sm font-bold uppercase tracking-wider transition-colors cursor-pointer'

  const palette = props.variant === 'panel'
    ? isActive
      ? 'border-black bg-black text-white'
      : 'border-black bg-white hover:bg-gray-50'
    : isActive
      ? 'border-black bg-black text-white'
      : 'border-black bg-white text-black hover:bg-gray-50'

  return `${base} ${palette}`
}

function selectSample(sample: { url: string }) {
  activeSample.value = sample.url
  emit('load', sample.url)
}
</script>

