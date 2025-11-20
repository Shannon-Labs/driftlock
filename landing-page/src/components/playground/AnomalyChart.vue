<template>
  <div class="relative overflow-hidden rounded-3xl border border-slate-800 bg-gradient-to-b from-slate-950 via-slate-950/95 to-slate-900 px-6 py-6 shadow-[0_40px_120px_rgba(15,23,42,0.65)]">
    <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <p class="text-xs font-semibold uppercase tracking-[0.3em] text-cyan-300/80">
          cbad core · live signal
        </p>
        <h3 class="text-2xl font-semibold text-white">
          Compression Distance Monitor
        </h3>
        <p class="text-sm text-slate-400">
          Streaming NCD scores with anomaly spikes highlighted in real time.
        </p>
      </div>
      <div class="inline-flex items-center gap-3 rounded-2xl border border-white/10 bg-white/5 px-4 py-2 text-xs font-semibold uppercase tracking-[0.3em] text-white">
        <span class="h-2 w-2 animate-pulse rounded-full bg-emerald-400"></span>
        Threshold {{ (threshold ?? 0.5).toFixed(2) }}
      </div>
    </div>

    <div class="relative mt-6 h-64 w-full sm:h-72 lg:h-80">
      <Line
        v-if="displayedSeries.length"
        :data="chartData"
        :options="chartOptions"
      />
      <div
        v-else
        class="absolute inset-0 flex flex-col items-center justify-center text-center text-slate-400"
      >
        <span class="text-3xl">⌁</span>
        <p class="mt-3 text-sm uppercase tracking-[0.3em]">Awaiting signal</p>
        <p class="text-xs text-slate-500">Load a sample or run the API to stream data</p>
      </div>
    </div>

    <div class="mt-5 flex flex-wrap gap-6 text-xs font-medium uppercase tracking-[0.3em] text-slate-500">
      <div class="flex items-center gap-2">
        <span class="h-1.5 w-1.5 rounded-full bg-cyan-300 animate-pulse"></span>
        Baseline tuned
      </div>
      <div class="flex items-center gap-2">
        <span class="h-1.5 w-1.5 rounded-full bg-fuchsia-300 animate-pulse"></span>
        Gemini insights optional
      </div>
      <div class="flex items-center gap-2">
        <span class="h-1.5 w-1.5 rounded-full bg-rose-400 animate-pulse"></span>
        Glass-box anomalies
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  LineElement,
  PointElement,
  LinearScale,
  CategoryScale,
  Tooltip,
  Filler,
  Legend,
} from 'chart.js'
import type { ChartOptions } from 'chart.js'
import { computed, onBeforeUnmount, ref, watch } from 'vue'

ChartJS.register(LineElement, PointElement, LinearScale, CategoryScale, Tooltip, Filler, Legend)

interface SeriesPoint {
  index: number
  score: number
  anomaly?: boolean
  record?: any
}

const props = withDefaults(defineProps<{
  series: SeriesPoint[]
  threshold?: number
}>(), {
  series: () => [],
  threshold: 0.5,
})

const emit = defineEmits<{ (e: 'select', payload: { index: number, record?: any }): void }>()

const displayedSeries = ref<SeriesPoint[]>([])
const streamTimer = ref<number | null>(null)

function hydrateSeries(points: SeriesPoint[]) {
  if (streamTimer.value) {
    clearInterval(streamTimer.value)
    streamTimer.value = null
  }
  displayedSeries.value = []
  if (!points?.length) return
  let cursor = 0
  streamTimer.value = window.setInterval(() => {
    displayedSeries.value = [...displayedSeries.value, points[cursor]]
    cursor += 1
    if (cursor >= points.length) {
      if (streamTimer.value) {
        clearInterval(streamTimer.value)
        streamTimer.value = null
      }
    }
  }, 80)
}

watch(() => props.series, (next) => hydrateSeries(next ?? []), { immediate: true, deep: true })

onBeforeUnmount(() => {
  if (streamTimer.value) {
    clearInterval(streamTimer.value)
  }
})

const chartData = computed(() => ({
  labels: displayedSeries.value.map((point) => `#${point.index}`),
  datasets: [
    {
      label: 'NCD Score',
      data: displayedSeries.value.map((point) => point.score),
      fill: true,
      tension: 0.35,
      pointRadius: displayedSeries.value.map((point) => point.anomaly ? 6 : 3),
      pointHoverRadius: displayedSeries.value.map((point) => point.anomaly ? 7 : 4),
      pointBorderWidth: displayedSeries.value.map((point) => point.anomaly ? 2 : 0),
      pointBackgroundColor: displayedSeries.value.map((point) =>
        point.anomaly ? '#f87171' : '#38bdf8'
      ),
      pointBorderColor: displayedSeries.value.map((point) =>
        point.anomaly ? '#fecaca' : '#bae6fd'
      ),
      borderColor: '#38bdf8',
      segment: {
        borderColor: (context: any) => {
          const current = context?.p1?.raw
          return current?.anomaly ? '#f87171' : '#38bdf8'
        },
      },
      backgroundColor: (context: any) => {
        const chart = context.chart
        const { ctx, chartArea } = chart
        if (!chartArea) return 'rgba(14,165,233,0.15)'
        const gradient = ctx.createLinearGradient(0, chartArea.bottom, 0, chartArea.top)
        gradient.addColorStop(0, 'rgba(15,23,42,0)')
        gradient.addColorStop(0.9, 'rgba(14,165,233,0.5)')
        return gradient
      },
    },
    {
      label: 'Threshold',
      data: displayedSeries.value.map(() => props.threshold ?? 0.5),
      borderDash: [6, 6],
      borderWidth: 1.2,
      borderColor: 'rgba(226,232,240,0.6)',
      pointRadius: 0,
      fill: false,
    },
  ],
}))

const chartOptions = computed<ChartOptions<'line'>>(() => ({
  responsive: true,
  maintainAspectRatio: false,
  animation: {
    duration: 450,
    easing: 'easeOutQuart',
  },
  plugins: {
    legend: {
      display: false,
    },
    tooltip: {
      backgroundColor: 'rgba(15,23,42,0.9)',
      borderColor: 'rgba(148,163,184,0.4)',
      borderWidth: 1,
      callbacks: {
        label: (context: any) => {
          const point = displayedSeries.value[context.dataIndex]
          const base = `NCD: ${Number(context.parsed.y).toFixed(3)}`
          return point?.anomaly ? `${base} · anomaly spike` : base
        },
      },
    },
  },
  scales: {
    x: {
      ticks: {
        color: 'rgba(191,219,254,0.7)',
      },
      grid: {
        color: 'rgba(30,41,59,0.6)',
      },
    },
    y: {
      min: 0,
      max: 1,
      ticks: {
        stepSize: 0.1,
        color: 'rgba(191,219,254,0.7)',
      },
      grid: {
        color: 'rgba(30,41,59,0.6)',
      },
    },
  },
  onClick: (_event, elements: any[]) => {
    if (!elements?.length) return
    const index = elements[0].index
    const target = displayedSeries.value[index]
    if (target) {
      emit('select', { index: target.index, record: target.record })
    }
  },
}) as ChartOptions<'line'>)
</script>


