<template>
  <div class="h-64">
    <Line :data="chartData" :options="chartOptions" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

interface DailyUsage {
  date: string
  event_count: number
  request_count: number
  anomaly_count: number
}

const props = defineProps<{
  data: DailyUsage[]
}>()

const chartData = computed(() => {
  const labels = props.data.map(d => {
    const date = new Date(d.date)
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
  })

  return {
    labels,
    datasets: [
      {
        label: 'Events',
        data: props.data.map(d => d.event_count),
        borderColor: '#000000',
        backgroundColor: 'rgba(0, 0, 0, 0.1)',
        fill: true,
        tension: 0.3,
        pointRadius: 3,
        pointHoverRadius: 5,
        pointBackgroundColor: '#000000',
        borderWidth: 2,
      },
      {
        label: 'Anomalies',
        data: props.data.map(d => d.anomaly_count),
        borderColor: '#DC2626',
        backgroundColor: 'rgba(220, 38, 38, 0.1)',
        fill: false,
        tension: 0.3,
        pointRadius: 3,
        pointHoverRadius: 5,
        pointBackgroundColor: '#DC2626',
        borderWidth: 2,
      },
    ],
  }
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  interaction: {
    intersect: false,
    mode: 'index' as const,
  },
  plugins: {
    legend: {
      position: 'top' as const,
      align: 'end' as const,
      labels: {
        font: {
          family: 'monospace',
          size: 11,
          weight: 'bold' as const,
        },
        usePointStyle: true,
        padding: 16,
      },
    },
    tooltip: {
      backgroundColor: '#000000',
      titleFont: {
        family: 'monospace',
        size: 12,
        weight: 'bold' as const,
      },
      bodyFont: {
        family: 'monospace',
        size: 11,
      },
      padding: 12,
      cornerRadius: 0,
      displayColors: true,
    },
  },
  scales: {
    x: {
      grid: {
        display: false,
      },
      ticks: {
        font: {
          family: 'monospace',
          size: 10,
        },
        maxRotation: 45,
        minRotation: 45,
      },
      border: {
        color: '#000000',
        width: 2,
      },
    },
    y: {
      beginAtZero: true,
      grid: {
        color: '#E5E7EB',
      },
      ticks: {
        font: {
          family: 'monospace',
          size: 10,
        },
        precision: 0,
      },
      border: {
        color: '#000000',
        width: 2,
      },
    },
  },
}
</script>
