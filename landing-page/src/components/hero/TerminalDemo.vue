<template>
  <div class="font-mono text-xs sm:text-sm bg-black text-gray-300 p-4 rounded-none border-2 border-gray-800 shadow-[8px_8px_0px_0px_rgba(0,0,0,0.2)] h-full min-h-[300px] flex flex-col relative overflow-hidden">
    <!-- Terminal Header -->
    <div class="flex items-center justify-between border-b border-gray-800 pb-2 mb-2">
      <div class="flex gap-2">
        <div class="w-3 h-3 rounded-full bg-red-500 opacity-50"></div>
        <div class="w-3 h-3 rounded-full bg-yellow-500 opacity-50"></div>
        <div class="w-3 h-3 rounded-full bg-green-500 opacity-50"></div>
      </div>
      <div class="text-gray-500 text-[10px] uppercase tracking-widest">driftlock-cli — v1.0.4</div>
    </div>

    <!-- Terminal Body -->
    <div class="flex-1 overflow-hidden relative font-mono">
      <div v-for="(log, index) in visibleLogs" :key="index" class="mb-1 leading-relaxed break-all">
        <span class="text-gray-500 mr-2">[{{ log.time }}]</span>
        <span :class="getLevelColor(log.level)" class="mr-2">[{{ log.level }}]</span>
        <span :class="{'text-white font-bold': log.highlight}">{{ log.message }}</span>
      </div>
      <!-- Cursor -->
      <div class="inline-block w-2 h-4 bg-gray-500 animate-pulse align-middle"></div>
    </div>
    
    <!-- Scan Overlay (Only appears on anomaly) -->
    <div v-if="scanning" class="absolute inset-0 bg-red-900/10 pointer-events-none flex items-center justify-center backdrop-blur-[1px]">
        <div class="border border-red-500 bg-black text-red-500 px-4 py-2 font-bold uppercase tracking-widest animate-pulse border-2 shadow-[4px_4px_0px_0px_#ef4444]">
            Anomaly Detected
        </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const scanning = ref(false)

interface Log {
  time: string
  level: string
  message: string
  highlight?: boolean
}

const allLogs = [
  { level: 'INFO', message: 'Connecting to stream wss://api.prod.svc/v1/metrics...' },
  { level: 'INFO', message: 'Connection established. Latency: 4ms' },
  { level: 'INFO', message: 'Loaded 4096-bit compression model' },
  { level: 'DEBUG', message: 'Sampler rate set to 100Hz' },
  { level: 'INFO', message: 'Ingesting frame #104922 (Size: 2KB)' },
  { level: 'INFO', message: 'Ingesting frame #104923 (Size: 2KB)' },
  { level: 'INFO', message: 'Ingesting frame #104924 (Size: 2KB)' },
  { level: 'INFO', message: 'Ingesting frame #104925 (Size: 2KB)' },
  { level: 'WARN', message: 'Entropy spike detected in frame #104926', highlight: true },
  { level: 'CRIT', message: 'ANOMALY CONFIRMED: Compression ratio drop (0.8 -> 0.4)', highlight: true },
  { level: 'INFO', message: 'Generating evidence bundle drift-88a2...' },
  { level: 'INFO', message: 'Alert sent to PagerDuty' },
  { level: 'INFO', message: 'Resuming monitoring...' },
  { level: 'INFO', message: 'Ingesting frame #104927 (Size: 2KB)' },
]

const visibleLogs = ref<Log[]>([])
let intervalId: any

const getLevelColor = (level: string) => {
  switch (level) {
    case 'INFO': return 'text-blue-400'
    case 'DEBUG': return 'text-gray-500'
    case 'WARN': return 'text-yellow-400'
    case 'CRIT': return 'text-red-500'
    default: return 'text-gray-300'
  }
}

const startSimulation = () => {
    let i = 0
    visibleLogs.value = []
    
    intervalId = setInterval(() => {
        if (i >= allLogs.length) {
            // Reset
            visibleLogs.value = []
            i = 0
            scanning.value = false
            return
        }

        const now = new Date()
        const timeStr = `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}:${now.getSeconds().toString().padStart(2, '0')}`
        
        visibleLogs.value.push({
            ...allLogs[i],
            time: timeStr
        })

        // Scroll to bottom
        // (In a real app we'd use a ref to the container, but Vue reactivity handles the render fast enough for this simple demo)
        
        // Trigger visual alarm
        if (allLogs[i].level === 'CRIT') {
            scanning.value = true
            setTimeout(() => { scanning.value = false }, 2000)
        }

        i++
    }, 800) // Speed of logs
}

onMounted(() => {
    startSimulation()
})

onUnmounted(() => {
    clearInterval(intervalId)
})
</script>
