<template>
  <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-700 p-4">
    <div class="flex items-center justify-between mb-3">
      <h3 class="font-semibold text-gray-900 dark:text-white">Samples</h3>
      <span class="text-xs text-gray-500 dark:text-gray-400">Auto-runs on load</span>
    </div>
    <div class="flex flex-wrap gap-2">
      <button
        v-for="s in samples"
        :key="s.name"
        class="px-3 py-2 text-sm border rounded bg-gray-50 dark:bg-gray-800 hover:bg-gray-100 dark:hover:bg-gray-700 transition"
        @click="$emit('load', s.url)"
      >
        <div class="font-medium">{{ s.name }}</div>
        <div class="text-[11px] text-gray-500 dark:text-gray-400">{{ s.desc }}</div>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
const samples = [
  { name: 'Payments (1K events)', url: '/samples/payments.ndjson', desc: 'Financial demo stream' },
  { name: 'NDJSON (small)', url: '/samples/small.ndjson', desc: '10 events' },
  { name: 'JSON Array (small)', url: '/samples/small.json', desc: '10 events' },
]
</script>

