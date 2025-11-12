<template>
  <div class="overflow-x-auto">
    <table class="min-w-full text-sm">
      <thead>
        <tr class="text-left border-b-2 border-gray-200 dark:border-gray-700">
          <th class="py-3 pr-4 font-semibold text-gray-900 dark:text-white">#</th>
          <th class="py-3 pr-4 font-semibold text-gray-900 dark:text-white">NCD</th>
          <th class="py-3 pr-4 font-semibold text-gray-900 dark:text-white">p-value</th>
          <th class="py-3 pr-4 font-semibold text-gray-900 dark:text-white">Confidence</th>
          <th class="py-3 pr-4 font-semibold text-gray-900 dark:text-white">Why</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="row in items" :key="row.index" class="border-b border-gray-100 dark:border-gray-800 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors">
          <td class="py-3 pr-4 font-mono text-gray-900 dark:text-gray-100">{{ row.index }}</td>
          <td class="py-3 pr-4 font-mono text-gray-700 dark:text-gray-300">{{ fmt(row.metrics.ncd) }}</td>
          <td class="py-3 pr-4 font-mono text-gray-700 dark:text-gray-300">{{ fmt(row.metrics.p_value) }}</td>
          <td class="py-3 pr-4 font-mono text-blue-600 dark:text-blue-400 font-semibold">{{ pct(row.metrics.confidence_level) }}</td>
          <td class="py-3 pr-4 text-gray-600 dark:text-gray-400">{{ row.why }}</td>
        </tr>
      </tbody>
    </table>
    <div v-if="!items?.length" class="text-center text-gray-500 dark:text-gray-400 text-sm py-8">
      <div class="text-4xl mb-2">✓</div>
      <div>No anomalies detected</div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ items: any[] }>()
function fmt(n: number) { return typeof n === 'number' ? n.toFixed(3) : n }
function pct(n: number) { return typeof n === 'number' ? (n * 100).toFixed(1) + '%' : n }
</script>

