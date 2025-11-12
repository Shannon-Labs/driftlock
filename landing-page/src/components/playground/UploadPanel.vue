<template>
  <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-700 p-6 shadow-sm">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Data</h3>
      <div class="text-sm text-gray-500 dark:text-gray-400">Format:
        <select v-model="format" class="ml-2 border rounded-lg px-3 py-1.5 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
          <option value="ndjson">NDJSON</option>
          <option value="json">JSON Array</option>
        </select>
      </div>
    </div>
    <div class="mb-4">
      <textarea 
        v-model="text" 
        rows="10" 
        class="w-full border rounded-lg p-3 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 font-mono text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500" 
        placeholder="Paste NDJSON or JSON array..."
      ></textarea>
    </div>
    <div class="flex items-center gap-3">
      <label class="px-4 py-2 bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 rounded-lg cursor-pointer hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors text-sm font-medium">
        <input type="file" @change="onFile" class="hidden" accept=".json,.jsonl,.ndjson" />
        Choose File
      </label>
      <button 
        class="px-4 py-2 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white rounded-lg font-medium transition-all duration-200 transform hover:scale-105 shadow-sm" 
        @click="emitData"
      >
        Use Data
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const emit = defineEmits<{ (e: 'data', p: { text: string, format: 'ndjson' | 'json' }): void }>()

const text = ref('')
const format = ref<'ndjson' | 'json'>('ndjson')

function emitData() {
  emit('data', { text: text.value, format: format.value })
}

async function onFile(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  const content = await file.text()
  text.value = content
  if (file.name.endsWith('.json')) format.value = 'json'
  if (file.name.endsWith('.jsonl') || file.name.endsWith('.ndjson')) format.value = 'ndjson'
  emitData()
}

watch(text, () => {
  // lightweight auto-detect on paste
  const t = text.value.trim()
  if (t.startsWith('[')) format.value = 'json'
  else format.value = 'ndjson'
})
</script>

