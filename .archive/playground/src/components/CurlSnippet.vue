<template>
  <div class="text-xs bg-gray-100 dark:bg-gray-800 rounded p-2 overflow-x-auto">
    <code>{{ curl }}</code>
    <button class="ml-2 px-2 py-0.5 border rounded" @click="copy">Copy</button>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ curl: string }>()
async function copy() {
  try { await navigator.clipboard.writeText(props.curl) } catch {}
}
</script>


