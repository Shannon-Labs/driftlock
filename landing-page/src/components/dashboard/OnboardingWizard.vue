<template>
  <div class="bg-white border-2 border-black p-8 shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] mb-12">
    <div class="mb-8">
      <h2 class="text-2xl font-bold uppercase tracking-wide mb-2">Welcome to Driftlock</h2>
      <p class="font-serif text-gray-600">Let's get you set up in minutes.</p>
    </div>

    <!-- Steps Indicator -->
    <div class="flex items-center justify-between mb-8 relative">
        <div class="absolute left-0 top-1/2 -translate-y-1/2 w-full h-0.5 bg-gray-200 -z-10"></div>
        <div class="absolute left-0 top-1/2 -translate-y-1/2 h-0.5 bg-black transition-all duration-300 -z-10" :style="{ width: progressWidth }"></div>

        <div v-for="s in [1, 2, 3]" :key="s" 
             class="w-8 h-8 flex items-center justify-center font-bold border-2 transition-colors duration-300 bg-white"
             :class="[
               step >= s ? 'border-black text-black' : 'border-gray-300 text-gray-300',
               step > s ? 'bg-black text-white' : ''
             ]">
            <span v-if="step > s">✓</span>
            <span v-else>{{ s }}</span>
        </div>
    </div>

    <!-- Step 1: Get API Key -->
    <div v-if="step === 1" class="animate-in slide-in-from-right-4 fade-in duration-300">
      <h3 class="text-lg font-bold uppercase mb-4">1. Get Your API Key</h3>
      <p class="text-sm font-serif mb-4">You'll need this key to send data to Driftlock.</p>
      
      <div class="bg-black text-green-400 p-4 font-mono text-sm mb-6 flex items-center justify-between border border-gray-800 relative group">
        <span class="break-all">{{ visibleKey ? apiKey : maskedKey }}</span>
        <div class="flex gap-2">
             <button @click="visibleKey = !visibleKey" class="text-xs uppercase font-bold text-gray-500 hover:text-white">
                {{ visibleKey ? 'Hide' : 'Reveal' }}
            </button>
            <button @click="copyKey" class="text-xs uppercase font-bold text-gray-500 hover:text-white">
                {{ copied ? 'Copied' : 'Copy' }}
            </button>
        </div>
      </div>

      <button @click="step = 2" class="brutalist-button-primary w-full sm:w-auto">
        I've copied it &rarr;
      </button>
    </div>

    <!-- Step 2: Send Request -->
    <div v-if="step === 2" class="animate-in slide-in-from-right-4 fade-in duration-300">
      <h3 class="text-lg font-bold uppercase mb-4">2. Send Your First Event</h3>
      <p class="text-sm font-serif mb-4">Copy and run this command in your terminal to send a test anomaly.</p>

      <div class="bg-gray-100 border border-black p-4 font-mono text-xs overflow-x-auto mb-6 relative">
<pre class="whitespace-pre-wrap break-all">curl -X POST {{ apiUrl }}/v1/detect \
  -H "Authorization: Bearer {{ apiKey }}" \
  -d '{
    "stream": "test-stream",
    "events": [
      { "msg": "Normal behavior", "value": 10 },
      { "msg": "Normal behavior", "value": 12 },
      { "msg": "CRITICAL FAILURE", "value": 9999 }
    ]
  }'</pre>
        <button @click="copyCurl" class="absolute top-2 right-2 text-xs font-bold uppercase bg-white border border-black px-2 py-1 hover:bg-black hover:text-white transition-colors">
            {{ curlCopied ? 'Copied!' : 'Copy Command' }}
        </button>
      </div>

      <div class="flex gap-4">
        <button @click="step = 3" class="brutalist-button-primary w-full sm:w-auto">
          I've sent it &rarr;
        </button>
        <button @click="step = 1" class="text-sm font-bold uppercase hover:underline p-3">
          Back
        </button>
      </div>
    </div>

    <!-- Step 3: View Results -->
    <div v-if="step === 3" class="animate-in slide-in-from-right-4 fade-in duration-300 text-center">
      <div class="text-6xl mb-4">🎉</div>
      <h3 class="text-xl font-bold uppercase mb-2">You're Ready!</h3>
      <p class="text-sm font-serif mb-6 text-gray-600">You've successfully set up Driftlock. Your dashboard is ready.</p>
      
      <button @click="$emit('complete')" class="brutalist-button-primary w-full sm:w-auto">
        Go to Dashboard
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const props = defineProps<{
  apiKey: string
  apiUrl?: string
}>()

const emit = defineEmits(['complete'])

const step = ref(1)
const visibleKey = ref(false)
const copied = ref(false)
const curlCopied = ref(false)

const maskedKey = computed(() => {
    if (!props.apiKey) return 'dlk_................'
    return props.apiKey.substring(0, 8) + '................' + props.apiKey.substring(props.apiKey.length - 4)
})

const progressWidth = computed(() => {
    return ((step.value - 1) / 2) * 100 + '%'
})

const copyKey = () => {
    navigator.clipboard.writeText(props.apiKey)
    copied.value = true
    setTimeout(() => copied.value = false, 2000)
}

const copyCurl = () => {
    const cmd = `curl -X POST ${props.apiUrl || window.location.origin}/v1/detect \\
  -H "Authorization: Bearer ${props.apiKey}" \\
  -d '{
    "stream": "test-stream",
    "events": [
      { "msg": "Normal behavior", "value": 10 },
      { "msg": "Normal behavior", "value": 12 },
      { "msg": "CRITICAL FAILURE", "value": 9999 }
    ]
  }'`
    navigator.clipboard.writeText(cmd)
    curlCopied.value = true
    setTimeout(() => curlCopied.value = false, 2000)
}
</script>

<style scoped>
.brutalist-button-primary {
    @apply inline-flex items-center justify-center px-6 py-3 border-2 border-black bg-black text-white font-sans uppercase tracking-wider font-bold text-sm transition-all hover:bg-white hover:text-black;
}
</style>
