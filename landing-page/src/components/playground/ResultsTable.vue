<template>
  <div class="overflow-x-auto">
    <table class="min-w-full text-sm">
      <thead>
        <tr class="border-b-2 border-gray-200 text-left dark:border-gray-700">
          <th class="py-3 pr-4 font-semibold text-gray-900 dark:text-white">#</th>
          <th class="py-3 pr-4 font-semibold text-gray-900 dark:text-white">NCD</th>
          <th class="py-3 pr-4 font-semibold text-gray-900 dark:text-white">p-value</th>
          <th class="py-3 pr-4 font-semibold text-gray-900 dark:text-white">Confidence</th>
          <th class="py-3 pr-4 font-semibold text-gray-900 dark:text-white">Why</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="row in items"
          :key="row.index"
          :class="rowClasses(row)"
          @click="emit('select', row)"
        >
          <td class="py-3 pr-4 font-mono text-gray-900 dark:text-gray-100">{{ row.index }}</td>
          <td class="py-3 pr-4 font-mono text-gray-700 dark:text-gray-300">{{ fmt(row.metrics?.ncd) }}</td>
          <td class="py-3 pr-4 font-mono text-gray-700 dark:text-gray-300">{{ fmt(row.metrics?.p_value) }}</td>
          <td class="py-3 pr-4 font-mono font-semibold text-blue-600 dark:text-blue-400">{{ pct(row.metrics?.confidence_level) }}</td>
          <td class="py-3 pr-4 text-gray-600 dark:text-gray-400">{{ row.why || '—' }}</td>
        </tr>
      </tbody>
    </table>
    <div v-if="!items?.length" class="py-8 text-center text-sm text-gray-500 dark:text-gray-400">
      <div class="mb-2 text-4xl">✓</div>
      <div>No anomalies detected</div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  items: any[]
  selectedIndex?: number | null
}>(), {
  items: () => [],
  selectedIndex: null,
})

const emit = defineEmits<{ (e: 'select', row: any): void }>()

function fmt(n: number | undefined) { return typeof n === 'number' ? n.toFixed(3) : '—' }
function pct(n: number | undefined) { return typeof n === 'number' ? (n * 100).toFixed(1) + '%' : '—' }

function rowClasses(row: any) {
  const isActive = props.selectedIndex != null && (row.index === props.selectedIndex)
  const base = 'border-b border-gray-100 dark:border-gray-800 transition-colors cursor-pointer'
  return isActive
    ? `${base} bg-blue-50/80 dark:bg-blue-900/30`
    : `${base} hover:bg-gray-50 dark:hover:bg-gray-800/50`
}
</script>


