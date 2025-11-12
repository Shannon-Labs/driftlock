<template>
  <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-700 p-6 shadow-sm">
    <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Sample Data</h3>
    <div class="flex flex-wrap gap-3">
      <button 
        v-for="s in samples" 
        :key="s.name" 
        class="px-4 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-gray-50 dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 hover:border-blue-500 dark:hover:border-blue-400 transition-all duration-200 font-medium" 
        @click="$emit('load', s.url)"
      >
        {{ s.name }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
const samples = [
  { name: 'NDJSON (small)', url: '/samples/small.ndjson' },
  { name: 'JSON Array (small)', url: '/samples/small.json' },
]
</script>

