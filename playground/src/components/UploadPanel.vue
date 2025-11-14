<template>
  <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-700 p-4">
    <div class="flex items-center justify-between mb-3">
      <div>
        <h3 class="font-semibold text-gray-900 dark:text-white">Data</h3>
        <p class="text-sm text-gray-500 dark:text-gray-400">Paste data or drop a file - Driftlock will auto-tune and run.</p>
      </div>
      <div class="text-sm text-gray-500">Format:
        <select v-model="format" class="ml-2 border rounded px-2 py-1 bg-white dark:bg-gray-800">
          <option value="ndjson">NDJSON</option>
          <option value="json">JSON Array</option>
        </select>
      </div>
    </div>
    <div class="mb-3">
      <textarea
        v-model="text"
        rows="12"
        class="w-full border rounded p-2 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
        placeholder="Paste NDJSON or JSON array..."
      ></textarea>
    </div>
    <div class="flex items-center gap-3 text-sm text-gray-500 dark:text-gray-400">
      <input type="file" @change="onFile" class="text-sm" />
      <span>Changes trigger detection automatically.</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const emit = defineEmits<{ (e: 'data', p: { text: string, format: 'ndjson' | 'json' }): void }>()

const text = ref('')
const format = ref<'ndjson' | 'json'>('ndjson')
let debounceHandle: number | null = null

function scheduleEmit(immediate = false) {
  if (debounceHandle) {
    clearTimeout(debounceHandle)
    debounceHandle = null
  }
  const trigger = () => emit('data', { text: text.value, format: format.value })
  if (immediate) {
    trigger()
    return
  }
  debounceHandle = window.setTimeout(trigger, 500)
}

async function onFile(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  const content = await file.text()
  text.value = content
  if (file.name.endsWith('.json')) format.value = 'json'
  if (file.name.endsWith('.jsonl') || file.name.endsWith('.ndjson')) format.value = 'ndjson'
  scheduleEmit(true)
}

watch(text, () => {
  const t = text.value.trim()
  if (t.startsWith('[')) format.value = 'json'
  else if (t.length) format.value = 'ndjson'
  scheduleEmit()
})

watch(format, () => scheduleEmit())
</script>

