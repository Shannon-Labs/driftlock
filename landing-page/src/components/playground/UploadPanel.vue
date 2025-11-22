<template>
  <div class="bg-white border-2 border-black p-6 shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
    <div class="flex items-center justify-between mb-4 border-b-2 border-black pb-4">
      <h3 class="text-xl font-sans font-black uppercase tracking-tighter text-black">Data Input</h3>
      <div class="text-xs font-bold uppercase tracking-widest text-gray-500 flex items-center">Format:
        <select v-model="format" class="ml-2 border-2 border-black px-2 py-1 bg-white text-black font-mono text-xs focus:outline-none focus:ring-2 focus:ring-black">
          <option value="ndjson">NDJSON</option>
          <option value="json">JSON Array</option>
        </select>
      </div>
    </div>
    <div class="mb-4">
      <textarea 
        v-model="text" 
        rows="10" 
        class="w-full border-2 border-black p-3 bg-gray-50 text-black font-mono text-sm focus:outline-none focus:ring-4 focus:ring-black/10 placeholder-gray-500" 
        placeholder="PASTE NDJSON OR JSON ARRAY..."
      ></textarea>
    </div>
    <div class="flex items-center gap-3">
      <label class="px-4 py-2 border-2 border-black bg-white text-black font-bold uppercase tracking-widest text-xs cursor-pointer hover:bg-black hover:text-white transition-colors shadow-[4px_4px_0px_0px_rgba(0,0,0,0)] hover:shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]">
        <input type="file" @change="onFile" class="hidden" accept=".json,.jsonl,.ndjson" />
        Choose File
      </label>
      <button 
        class="px-4 py-2 border-2 border-black bg-black text-white font-bold uppercase tracking-widest text-xs hover:bg-white hover:text-black transition-colors shadow-[4px_4px_0px_0px_rgba(0,0,0,0)] hover:shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]" 
        @click="emitData"
      >
        Use Data
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  text: string
  format: 'ndjson' | 'json'
}>()

const emit = defineEmits<{
  (e: 'data', p: { text: string, format: 'ndjson' | 'json' }): void
  (e: 'update:text', value: string): void
  (e: 'update:format', value: 'ndjson' | 'json'): void
}>()

const text = ref(props.text)
const format = ref<'ndjson' | 'json'>(props.format)

watch(() => props.text, value => {
  if (value !== text.value) {
    text.value = value
  }
})

watch(() => props.format, value => {
  if (value !== format.value) {
    format.value = value
  }
})

watch(text, value => {
  emit('update:text', value)
  const trimmed = value.trim()
  const inferred: 'ndjson' | 'json' = trimmed.startsWith('[') ? 'json' : 'ndjson'
  if (inferred !== format.value) {
    format.value = inferred
  }
})

watch(format, value => {
  emit('update:format', value)
})

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
</script>

