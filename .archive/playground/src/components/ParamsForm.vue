<template>
  <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-700 p-4">
    <div class="flex items-start justify-between gap-3 mb-4">
      <div>
        <h3 class="font-semibold text-gray-900 dark:text-white">Detection settings</h3>
        <p class="text-sm text-gray-500 dark:text-gray-400">Derived automatically from the data you provide.</p>
      </div>
      <button
        class="text-sm px-3 py-1.5 border rounded-lg bg-gray-50 dark:bg-gray-800"
        @click="open = !open"
      >
        {{ open ? 'Hide' : 'Show' }} advanced
      </button>
    </div>

    <dl class="grid grid-cols-2 gap-3 text-sm mb-4">
      <div class="p-3 rounded-lg bg-gray-50 dark:bg-gray-800">
        <dt class="text-gray-500 dark:text-gray-400">Baseline</dt>
        <dd class="text-lg font-semibold text-gray-900 dark:text-white">{{ props.derived.baseline }}</dd>
      </div>
      <div class="p-3 rounded-lg bg-gray-50 dark:bg-gray-800">
        <dt class="text-gray-500 dark:text-gray-400">Window / Hop</dt>
        <dd class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ props.derived.window }} / {{ props.derived.hop }}
        </dd>
      </div>
    </dl>

    <transition name="fade">
      <div v-if="open" class="space-y-3 pt-2 border-t border-gray-200 dark:border-gray-700">
        <p class="text-xs text-gray-500 dark:text-gray-400">
          Override the auto values if you know you need something specific. Leave blank to keep auto tuning.
        </p>
        <div>
          <label class="text-sm text-gray-600 dark:text-gray-300">Baseline override</label>
          <input
            type="number"
            :value="stringValue(local.baseline)"
            @input="updateNumber('baseline', $event)"
            class="w-full border rounded px-2 py-1 bg-white dark:bg-gray-800"
            min="1"
            :placeholder="`Auto (${props.derived.baseline})`"
          />
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="text-sm text-gray-600 dark:text-gray-300">Window override</label>
            <input
              type="number"
              :value="stringValue(local.window)"
              @input="updateNumber('window', $event)"
              class="w-full border rounded px-2 py-1 bg-white dark:bg-gray-800"
              min="1"
              :placeholder="`Auto (${props.derived.window})`"
            />
          </div>
          <div>
            <label class="text-sm text-gray-600 dark:text-gray-300">Hop override</label>
            <input
              type="number"
              :value="stringValue(local.hop)"
              @input="updateNumber('hop', $event)"
              class="w-full border rounded px-2 py-1 bg-white dark:bg-gray-800"
              min="1"
              :placeholder="`Auto (${props.derived.hop})`"
            />
          </div>
        </div>
        <div>
          <label class="text-sm text-gray-600 dark:text-gray-300">Algorithm override</label>
          <select
            :value="local.algo ?? ''"
            @change="updateAlgo($event)"
            class="w-full border rounded px-2 py-1 bg-white dark:bg-gray-800"
          >
            <option value="">Auto ({{ props.derived.algo }})</option>
            <option value="zstd">zstd</option>
            <option value="lz4">lz4</option>
            <option value="gzip">gzip</option>
            <option value="openzl">openzl</option>
          </select>
        </div>
        <div class="pt-1">
          <button class="text-xs text-primary-600" @click="clearOverrides">Clear overrides</button>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue'

type DerivedParams = {
  baseline: number
  window: number
  hop: number
  algo: string
}

type Overrides = Partial<{
  baseline: number
  window: number
  hop: number
  algo: string
}>

const props = defineProps<{
  derived: DerivedParams
  overrides: Overrides
}>()
const emit = defineEmits<{ (e: 'change', payload: Overrides): void }>()

const open = ref(false)
const local = reactive<{
  baseline: number | null
  window: number | null
  hop: number | null
  algo: string | null
}>({
  baseline: props.overrides.baseline ?? null,
  window: props.overrides.window ?? null,
  hop: props.overrides.hop ?? null,
  algo: props.overrides.algo ?? null,
})

watch(
  () => props.overrides,
  (ov) => {
    local.baseline = ov.baseline ?? null
    local.window = ov.window ?? null
    local.hop = ov.hop ?? null
    local.algo = ov.algo ?? null
  },
  { deep: true }
)

watch(
  local,
  () => {
    emit('change', cleanedOverrides())
  },
  { deep: true }
)

function cleanedOverrides(): Overrides {
  const next: Overrides = {}
  if (typeof local.baseline === 'number' && !Number.isNaN(local.baseline)) next.baseline = local.baseline
  if (typeof local.window === 'number' && !Number.isNaN(local.window)) next.window = local.window
  if (typeof local.hop === 'number' && !Number.isNaN(local.hop)) next.hop = local.hop
  if (typeof local.algo === 'string' && local.algo.length) next.algo = local.algo
  return next
}

function updateNumber(field: 'baseline' | 'window' | 'hop', event: Event) {
  const input = event.target as HTMLInputElement
  const value = input.value.trim()
  const parsed = value === '' ? null : Number(value)
  ;(local as any)[field] = Number.isFinite(parsed as number) ? (parsed as number) : null
}

function updateAlgo(event: Event) {
  const select = event.target as HTMLSelectElement
  local.algo = select.value === '' ? null : select.value
}

function clearOverrides() {
  local.baseline = null
  local.window = null
  local.hop = null
  local.algo = null
}

function stringValue(value: number | null) {
  return value === null ? '' : String(value)
}
</script>

