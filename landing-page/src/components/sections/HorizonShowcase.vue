<template>
  <section id="showcase" class="section-padding bg-gray-50 border-b border-black">
    <div class="container-padding">
      <div class="mb-12 max-w-3xl">
        <h2 class="text-sm font-bold font-sans uppercase tracking-widest mb-2 border-b border-black inline-block">
          The Universal Radar
        </h2>
        <p class="text-4xl font-sans font-bold tracking-tight mt-4 uppercase">
          Tested on Reality.
        </p>
        <p class="text-xl font-serif mt-4 text-gray-600">
          Driftlock detects anomalies across every horizon. Select a scenario to see the volatility signature.
        </p>
      </div>

      <!-- The Bento Grid / Tabs Layout -->
      <div class="flex flex-col lg:flex-row gap-8">
        <!-- Sidebar (Tabs) -->
        <div class="lg:w-1/3 flex flex-col gap-2">
          <button 
            v-for="h in horizons" 
            :key="h.id"
            @click="activeId = h.id"
            class="text-left p-6 border-2 transition-all duration-200 group relative overflow-hidden"
            :class="activeId === h.id 
              ? 'border-black bg-black text-white shadow-[8px_8px_0px_0px_rgba(0,0,0,0.2)] translate-x-2' 
              : 'border-gray-200 bg-white text-gray-500 hover:border-black hover:text-black'"
          >
            <span class="text-xs font-bold uppercase tracking-widest mb-1 block opacity-60">
              Horizon 0{{ h.index }}
            </span>
            <span class="text-2xl font-sans font-bold uppercase block">
              {{ h.title }}
            </span>
             <span class="text-sm font-serif mt-2 block opacity-80" v-if="activeId === h.id">
              {{ h.subtitle }}
            </span>
          </button>
        </div>

        <!-- Main Display (The Card) -->
        <div class="lg:w-2/3 relative min-h-[500px]">
          <transition name="fade" mode="out-in">
            <div :key="activeHorizon.id" class="border-2 border-black bg-white p-8 shadow-[12px_12px_0px_0px_rgba(0,0,0,1)] h-full flex flex-col">
              
              <!-- Header -->
              <div class="flex justify-between items-start mb-8 border-b border-black pb-4">
                <div>
                  <div class="inline-block bg-black text-white px-2 py-1 text-xs font-bold uppercase tracking-widest mb-2">
                    {{ activeHorizon.tag }}
                  </div>
                  <h3 class="text-3xl font-sans font-black uppercase leading-none">
                    {{ activeHorizon.headline }}
                  </h3>
                </div>
                <div class="text-right hidden sm:block">
                  <div class="text-xs font-bold uppercase tracking-widest text-gray-400">Detection Latency</div>
                  <div class="text-xl font-mono font-bold">{{ activeHorizon.latency }}</div>
                </div>
              </div>

              <!-- Visual/Graph Area (Simulated) -->
              <div class="bg-white border border-black p-6 mb-8 flex-grow relative overflow-hidden group min-h-[200px]">
                 <!-- Background Grid -->
                 <div class="absolute inset-0 opacity-10 pointer-events-none" 
                      style="background-image: radial-gradient(#000 1px, transparent 1px); background-size: 20px 20px;">
                 </div>
                 
                 <!-- Simulated Graph Bars -->
                 <div class="flex items-end justify-between h-full gap-1 mb-4 pt-10">
                    <div v-for="n in 40" :key="n" 
                         class="bg-gray-300 w-full transition-all duration-500 origin-bottom"
                         :style="{ 
                           height: getBarHeight(n, activeHorizon.pattern) + '%',
                           opacity: n > 30 ? '1' : '0.5',
                           backgroundColor: n > 35 && activeHorizon.hasAnomaly ? '#000' : undefined
                         }"
                    ></div>
                 </div>

                 <!-- Overlay Verdict -->
                 <div class="absolute top-4 right-4 bg-white border border-black px-4 py-3 shadow-sm max-w-xs z-10">
                    <div class="flex items-center gap-2 mb-1">
                        <div class="w-2 h-2 rounded-full" :class="activeHorizon.hasAnomaly ? 'bg-black animate-pulse' : 'border-2 border-black bg-white'"></div>
                        <span class="text-xs font-bold uppercase tracking-widest">Forensic Verdict</span>
                    </div>
                    <p class="font-mono text-xs leading-relaxed mb-2">
                        {{ activeHorizon.verdict }}
                    </p>
                    <div class="border-t border-black pt-2 mt-2">
                         <div class="flex items-center gap-1 mb-1">
                            <span class="text-[10px] font-bold uppercase tracking-widest text-black">Gemini Insight</span>
                        </div>
                        <p class="font-serif text-xs italic text-gray-600">
                            "{{ activeHorizon.geminiInsight }}"
                        </p>
                    </div>
                 </div>
              </div>

              <!-- Footer/Action -->
              <div class="mt-auto flex flex-col sm:flex-row items-center justify-between gap-6">
                <div class="max-w-md">
                    <p class="font-serif text-sm text-gray-600 mb-1">
                        {{ activeHorizon.description }}
                    </p>
                    <p class="text-[10px] font-mono text-gray-400 uppercase tracking-wide">
                        Verified on Kaggle: {{ activeHorizon.kaggleRef }}
                    </p>
                </div>
                <button 
                  @click="$emit('load', activeHorizon.sampleUrl)"
                  class="w-full sm:w-auto whitespace-nowrap bg-black text-white px-6 py-3 text-sm font-bold uppercase tracking-widest hover:bg-gray-800 transition-colors border border-transparent hover:border-black hover:bg-white hover:text-black"
                >
                  Load Data →
                </button>
              </div>

            </div>
          </transition>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

const emit = defineEmits(['load'])

const activeId = ref('finance')

const horizons = [
  {
    id: 'finance',
    index: 1,
    title: 'Financial Fraud',
    subtitle: 'Credit Card Transactions',
    tag: 'FinTech',
    headline: 'The Hidden Pattern',
    latency: '12ms',
    hasAnomaly: true,
    pattern: 'spike',
    verdict: 'ANOMALY DETECTED: Compression ratio spike (2.4x). Entropy variance exceeds 3σ baseline.',
    description: 'Detects fraud patterns that rules engines miss by analyzing the entropy of transaction metadata.',
    kaggleRef: 'Credit Card Fraud Detection (ULB)',
    geminiInsight: 'Gemini 3 Pro analysis detects coordinated high-frequency micro-transactions consistent with skimming. Geographical entropy mismatch confirmed.',
    sampleUrl: '/samples/fraud.json'
  },
  {
    id: 'crypto',
    index: 2,
    title: 'Market Crash',
    subtitle: 'Terra/Luna Collapse',
    tag: 'DeFi',
    headline: 'The Death Spiral',
    latency: '45ms',
    hasAnomaly: true,
    pattern: 'volatility',
    verdict: 'CRITICAL DRIFT: Market structure collapse. NCD 0.92 > 0.5 threshold. Algorithmic de-pegging detected.',
    description: 'Identifies structural breaks in market data before price aggregators update.',
    kaggleRef: 'Terra Luna Crash Data (2022)',
    geminiInsight: 'Liquidity pool entropy variance > 3σ. Algorithmic de-peg imminent. Oracle latency exploit detected 45s pre-impact.',
    sampleUrl: '/samples/terra.json'
  },
  {
    id: 'aviation',
    index: 3,
    title: 'Aviation Ops',
    subtitle: 'Turbofan Degradation',
    tag: 'Critical Infra',
    headline: 'Operational Drift',
    latency: '120ms',
    hasAnomaly: false,
    pattern: 'wave',
    verdict: 'NOMINAL: Fleet telemetry within normal compression bounds. No mechanical divergence.',
    description: 'Monitors thousands of sensors. If a turbine vibrates differently, the compression ratio changes.',
    kaggleRef: 'NASA Turbofan Degradation',
    geminiInsight: 'Turbofan vibration harmonics nominal. Sensor fusion entropy within 0.5% of fleet baseline. No predictive maintenance required.',
    sampleUrl: '/samples/airline.json'
  },
  {
    id: 'safety',
    index: 4,
    title: 'AI Safety',
    subtitle: 'Prompt Injection',
    tag: 'GenAI',
    headline: 'Jailbreak Attempt',
    latency: '8ms',
    hasAnomaly: true,
    pattern: 'chaos',
    verdict: 'SECURITY EVENT: Input entropy matches known adversarial patterns. Prompt injection blocked.',
    description: 'Protects LLMs by detecting the statistical signature of jailbreak attempts.',
    kaggleRef: 'Jailbreak Prompt Dataset (HuggingFace)',
    geminiInsight: 'Adversarial prompt detected. Statistical signature matches "DAN 12.0" jailbreak pattern. Request rejected by entropy filter.',
    sampleUrl: '/samples/safety.json'
  }
]

const activeHorizon = computed(() => horizons.find(h => h.id === activeId.value) || horizons[0])

function getBarHeight(n: number, pattern: string) {
    // Simple simulation for the visual
    const seed = n * 0.5
    if (pattern === 'spike') {
        return n > 35 ? 90 : 20 + Math.sin(seed) * 10
    }
    if (pattern === 'volatility') {
        return n > 30 ? 40 + Math.random() * 50 : 30 + Math.sin(seed) * 5
    }
    if (pattern === 'wave') {
        return 40 + Math.sin(n * 0.2) * 20
    }
    // chaos
    return Math.random() * 80 + 10
}
</script>

<style scoped>
.section-padding {
    @apply py-16 sm:py-24;
}
.container-padding {
    @apply mx-auto max-w-7xl px-4 sm:px-6 lg:px-8;
}
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
