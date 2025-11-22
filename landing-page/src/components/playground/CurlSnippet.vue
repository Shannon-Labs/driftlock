<template>
  <details class="text-xs">
    <summary class="cursor-pointer text-black hover:underline font-bold uppercase tracking-wider">
      Show cURL
    </summary>
    <div class="mt-2 bg-gray-50 border-2 border-black p-3">
      <div class="flex items-start justify-between gap-2">
        <code class="flex-1 text-black font-mono break-all text-xs">{{ curl }}</code>
        <button 
          class="px-3 py-1 border-2 border-black bg-white text-black hover:bg-black hover:text-white transition-colors text-xs font-bold uppercase tracking-wider whitespace-nowrap" 
          @click="copy"
        >
          Copy
        </button>
      </div>
    </div>
  </details>
</template>

<script setup lang="ts">
const props = defineProps<{ curl: string }>()
async function copy() {
  try { 
    await navigator.clipboard.writeText(props.curl)
    // Could add a toast notification here
  } catch {}
}
</script>

