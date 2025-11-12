<template>
  <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-700 p-6 shadow-sm">
    <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Parameters</h3>
    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Baseline</label>
        <input 
          type="number" 
          v-model.number="local.baseline" 
          class="w-full border rounded-lg px-3 py-2 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-blue-500" 
          min="1" 
        />
      </div>
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Window</label>
          <input 
            type="number" 
            v-model.number="local.window" 
            class="w-full border rounded-lg px-3 py-2 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-blue-500" 
            min="1" 
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Hop</label>
          <input 
            type="number" 
            v-model.number="local.hop" 
            class="w-full border rounded-lg px-3 py-2 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-blue-500" 
            min="1" 
          />
        </div>
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Algorithm</label>
        <select 
          v-model="local.algo" 
          class="w-full border rounded-lg px-3 py-2 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
        >
          <option value="zstd">zstd</option>
          <option value="lz4">lz4</option>
          <option value="gzip">gzip</option>
          <option value="openzl">openzl</option>
        </select>
      </div>
      <div class="pt-2 flex gap-3">
        <button 
          class="flex-1 px-4 py-2 bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 rounded-lg border border-gray-300 dark:border-gray-700 hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors font-medium" 
          @click="reset"
        >
          Reset
        </button>
        <button 
          class="flex-1 px-4 py-2 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white rounded-lg font-medium transition-all duration-200 transform hover:scale-105 shadow-sm" 
          @click="apply"
        >
          Run Detection
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, watch } from 'vue'

const props = defineProps<{ params: { baseline: number, window: number, hop: number, algo: string } }>()
const emit = defineEmits<{ (e: 'update', p: any): void, (e: 'run'): void }>()

const local = reactive({ ...props.params })

watch(() => props.params, (p) => Object.assign(local, p))

function reset() {
  Object.assign(local, { baseline: 400, window: 1, hop: 1, algo: 'zstd' })
  emit('update', { ...local })
}

function apply() {
  emit('update', { ...local })
  emit('run')
}
</script>

