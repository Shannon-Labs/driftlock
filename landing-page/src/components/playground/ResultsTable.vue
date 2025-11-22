<template>
  <div class="overflow-x-auto">
    <table class="min-w-full text-sm border-collapse">
      <thead class="bg-black text-white">
        <tr class="border-b-2 border-black text-left">
          <th class="py-3 pr-4 font-bold uppercase tracking-wider">#</th>
          <th class="py-3 pr-4 font-bold uppercase tracking-wider">NCD</th>
          <th class="py-3 pr-4 font-bold uppercase tracking-wider">p-value</th>
          <th class="py-3 pr-4 font-bold uppercase tracking-wider">Confidence</th>
          <th class="py-3 pr-4 font-bold uppercase tracking-wider">Why</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="row in items"
          :key="row.index"
          :class="rowClasses(row)"
          @click="emit('select', row)"
        >
          <td class="py-3 pr-4 font-mono font-bold text-black">{{ row.index }}</td>
          <td class="py-3 pr-4 font-mono text-gray-800">{{ fmt(row.metrics?.ncd) }}</td>
          <td class="py-3 pr-4 font-mono text-gray-800">{{ fmt(row.metrics?.p_value) }}</td>
          <td class="py-3 pr-4 font-mono font-bold text-black">{{ pct(row.metrics?.confidence_level) }}</td>
          <td class="py-3 pr-4 text-gray-600 font-serif">{{ row.why || '—' }}</td>
        </tr>
      </tbody>
    </table>
    <div v-if="!items?.length" class="py-12 text-center text-sm text-gray-500">
      <div class="mb-2 text-4xl">✓</div>
      <div class="font-mono uppercase tracking-widest">No anomalies detected</div>
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
  const base = 'border-b-2 border-black transition-colors cursor-pointer'
  return isActive
    ? `${base} bg-black text-white`
    : `${base} hover:bg-gray-50`
}
</script>


