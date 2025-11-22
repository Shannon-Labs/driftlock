<template>
  <div class="bg-white border-2 border-black p-6 shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
    <h3 class="text-xl font-sans font-black uppercase tracking-tighter text-black mb-4 border-b-2 border-black pb-4">Parameters</h3>
    <div class="space-y-4">
      <div>
        <label class="block text-xs font-bold uppercase tracking-widest text-gray-500 mb-2">Baseline</label>
        <input 
          type="number" 
          v-model.number="local.baseline" 
          class="w-full border-2 border-black px-3 py-2 bg-white text-black font-mono focus:outline-none focus:ring-4 focus:ring-black/10" 
          min="1" 
        />
      </div>
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="block text-xs font-bold uppercase tracking-widest text-gray-500 mb-2">Window</label>
          <input 
            type="number" 
            v-model.number="local.window" 
            class="w-full border-2 border-black px-3 py-2 bg-white text-black font-mono focus:outline-none focus:ring-4 focus:ring-black/10" 
            min="1" 
          />
        </div>
        <div>
          <label class="block text-xs font-bold uppercase tracking-widest text-gray-500 mb-2">Hop</label>
          <input 
            type="number" 
            v-model.number="local.hop" 
            class="w-full border-2 border-black px-3 py-2 bg-white text-black font-mono focus:outline-none focus:ring-4 focus:ring-black/10" 
            min="1" 
          />
        </div>
      </div>
      <div>
        <label class="block text-xs font-bold uppercase tracking-widest text-gray-500 mb-2">Algorithm</label>
        <select 
          v-model="local.algo" 
          class="w-full border-2 border-black px-3 py-2 bg-white text-black font-mono focus:outline-none focus:ring-4 focus:ring-black/10"
        >
          <option value="zstd">zstd</option>
          <option value="lz4">lz4</option>
          <option value="gzip">gzip</option>
          <option value="openzl">openzl</option>
        </select>
      </div>
      <div class="pt-2 flex gap-3">
        <button 
          class="flex-1 px-4 py-2 border-2 border-black bg-white text-black font-bold uppercase tracking-widest text-xs hover:bg-black hover:text-white transition-colors" 
          @click="reset"
        >
          Reset
        </button>
        <button 
          class="flex-1 px-4 py-2 border-2 border-black bg-black text-white font-bold uppercase tracking-widest text-xs hover:bg-white hover:text-black transition-colors shadow-[4px_4px_0px_0px_rgba(0,0,0,0)] hover:shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]" 
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

