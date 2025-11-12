<template>
  <details class="text-xs">
    <summary class="cursor-pointer text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100 font-medium">
      Show cURL
    </summary>
    <div class="mt-2 bg-gray-50 dark:bg-gray-800 rounded-lg p-3 border border-gray-200 dark:border-gray-700">
      <div class="flex items-start justify-between gap-2">
        <code class="flex-1 text-gray-800 dark:text-gray-200 break-all">{{ curl }}</code>
        <button 
          class="px-3 py-1 border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-600 transition-colors text-xs font-medium whitespace-nowrap" 
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

