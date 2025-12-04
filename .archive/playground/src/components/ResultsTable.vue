<template>
  <div class="overflow-x-auto">
    <table class="min-w-full text-sm">
      <thead>
        <tr class="text-left border-b border-gray-200 dark:border-gray-700">
          <th class="py-2 pr-4">#</th>
          <th class="py-2 pr-4">NCD</th>
          <th class="py-2 pr-4">p-value</th>
          <th class="py-2 pr-4">Confidence</th>
          <th class="py-2 pr-4">Why</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="row in items" :key="row.index" class="border-b border-gray-100 dark:border-gray-800">
          <td class="py-2 pr-4">{{ row.index }}</td>
          <td class="py-2 pr-4">{{ fmt(row.metrics.ncd) }}</td>
          <td class="py-2 pr-4">{{ fmt(row.metrics.p_value) }}</td>
          <td class="py-2 pr-4">{{ pct(row.metrics.confidence_level) }}</td>
          <td class="py-2 pr-4">{{ row.why }}</td>
        </tr>
      </tbody>
    </table>
    <div v-if="!items?.length" class="text-gray-500 text-sm py-3">No anomalies.</div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ items: any[] }>()
function fmt(n: number) { return typeof n === 'number' ? n.toFixed(3) : n }
function pct(n: number) { return typeof n === 'number' ? (n * 100).toFixed(1) + '%' : n }
</script>


