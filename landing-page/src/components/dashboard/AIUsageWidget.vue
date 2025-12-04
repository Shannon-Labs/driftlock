<template>
  <div class="bg-white border-2 border-black shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
    <div class="px-6 py-6">
      <!-- Header -->
      <div class="flex justify-between items-center border-b-2 border-black pb-4 mb-4">
        <h3 class="text-xl font-bold uppercase tracking-wide text-black">AI Usage & Costs</h3>
        <div
          class="px-3 py-1 text-xs font-bold uppercase tracking-widest border-2 border-black"
          :class="modelBadgeClass"
        >
          {{ modelDisplayName }}
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="flex items-center justify-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-2 border-black border-t-transparent"></div>
      </div>

      <!-- Content -->
      <div v-else>
        <!-- Current Period -->
        <div class="mb-6">
          <div class="flex justify-between items-center mb-3">
            <span class="text-xs font-bold uppercase tracking-widest text-gray-500">Current Period</span>
            <span class="text-xs font-mono text-gray-500">{{ periodDates }}</span>
          </div>

          <!-- AI Calls Usage -->
          <div class="mb-4">
            <div class="flex justify-between mb-1">
              <span class="text-sm font-bold">AI Calls</span>
              <span class="text-sm font-mono">{{ usage.callsUsed.toLocaleString() }} / {{ usage.callsLimit.toLocaleString() }}</span>
            </div>
            <div class="h-3 bg-gray-200 border-2 border-black">
              <div
                class="h-full transition-all duration-300"
                :class="getUsageBarClass(usage.callsPercent)"
                :style="{ width: `${Math.min(usage.callsPercent, 100)}%` }"
              ></div>
            </div>
          </div>

          <!-- Cost Usage -->
          <div class="mb-4">
            <div class="flex justify-between mb-1">
              <span class="text-sm font-bold">AI Costs</span>
              <span class="text-sm font-mono">${{ usage.costUsed.toFixed(2) }} / ${{ usage.costLimit.toFixed(2) }}</span>
            </div>
            <div class="h-3 bg-gray-200 border-2 border-black">
              <div
                class="h-full transition-all duration-300"
                :class="getUsageBarClass(usage.costPercent)"
                :style="{ width: `${Math.min(usage.costPercent, 100)}%` }"
              ></div>
            </div>
          </div>

          <!-- Token Breakdown -->
          <div class="flex justify-between bg-gray-50 border-2 border-black p-3">
            <div class="text-center">
              <div class="text-xs font-bold uppercase tracking-widest text-gray-500">Input</div>
              <div class="text-sm font-mono font-bold">{{ formatNumber(usage.inputTokens) }}</div>
            </div>
            <div class="text-center">
              <div class="text-xs font-bold uppercase tracking-widest text-gray-500">Output</div>
              <div class="text-sm font-mono font-bold">{{ formatNumber(usage.outputTokens) }}</div>
            </div>
            <div class="text-center">
              <div class="text-xs font-bold uppercase tracking-widest text-gray-500">Total</div>
              <div class="text-sm font-mono font-bold">{{ formatNumber(usage.inputTokens + usage.outputTokens) }}</div>
            </div>
          </div>
        </div>

        <!-- Forecast -->
        <div class="border-t-2 border-black pt-4 mb-4">
          <div class="flex justify-between items-center mb-2">
            <span class="text-sm font-bold">Monthly Forecast</span>
            <button
              @click="showBreakdown = !showBreakdown"
              class="text-xs font-bold uppercase tracking-widest text-gray-500 hover:text-black transition-colors"
            >
              {{ showBreakdown ? 'Hide' : 'Show' }} Details
            </button>
          </div>
          <div class="flex justify-between items-baseline">
            <span class="text-xs text-gray-500">Estimated Total:</span>
            <span class="text-xl font-bold font-mono">${{ forecastTotal.toFixed(2) }}</span>
          </div>

          <!-- Forecast Warning -->
          <div
            v-if="forecastWarning"
            class="mt-2 flex items-center gap-2 bg-yellow-100 border-2 border-yellow-600 p-2 text-xs font-bold text-yellow-800"
          >
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            {{ forecastWarning }}
          </div>
        </div>

        <!-- Cost Breakdown (collapsible) -->
        <div v-if="showBreakdown" class="bg-gray-50 border-2 border-black p-4 mb-4">
          <div class="text-xs font-bold uppercase tracking-widest text-gray-500 mb-3">Cost Breakdown</div>
          <div class="space-y-2 text-sm">
            <div class="flex justify-between">
              <span>AI Processing</span>
              <span class="font-mono">${{ usage.costUsed.toFixed(2) }}</span>
            </div>
            <div class="flex justify-between">
              <span>Margin (15%)</span>
              <span class="font-mono">${{ marginCost.toFixed(2) }}</span>
            </div>
            <div class="flex justify-between border-t border-gray-300 pt-2 font-bold">
              <span>Total Charges</span>
              <span class="font-mono">${{ totalCost.toFixed(2) }}</span>
            </div>
          </div>
        </div>

        <!-- Alerts -->
        <div v-if="alerts.length > 0" class="space-y-2 mb-4">
          <div
            v-for="alert in alerts"
            :key="alert.id"
            class="flex items-center gap-2 p-3 text-sm font-bold border-2"
            :class="alertClass(alert.severity)"
          >
            <svg v-if="alert.severity === 'error'" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <svg v-else class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            {{ alert.message }}
          </div>
        </div>

        <!-- Upgrade Prompt -->
        <div
          v-if="showUpgradePrompt"
          class="bg-black text-white border-2 border-black p-4 text-center"
        >
          <p class="text-sm mb-3">You're approaching your AI usage limits. Upgrade for more AI-powered insights.</p>
          <button
            @click="$emit('upgrade')"
            class="bg-white text-black border-2 border-white px-4 py-2 text-xs font-bold uppercase tracking-widest hover:bg-gray-100 transition-colors"
          >
            Upgrade Plan
          </button>
        </div>

        <!-- Configure AI Button -->
        <div class="border-t-2 border-black pt-4">
          <button
            @click="showConfig = !showConfig"
            class="flex items-center gap-2 text-xs font-bold uppercase tracking-widest text-gray-500 hover:text-black transition-colors"
          >
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            Configure AI
          </button>

          <!-- Config Panel -->
          <div v-if="showConfig" class="mt-4 bg-gray-50 border-2 border-black p-4 space-y-4">
            <div>
              <label class="block text-xs font-bold uppercase tracking-widest text-gray-500 mb-2">
                Analysis Threshold
              </label>
              <div class="flex items-center gap-3">
                <input
                  type="range"
                  min="0.1"
                  max="1.0"
                  step="0.1"
                  v-model.number="localConfig.threshold"
                  @change="updateConfig"
                  class="flex-1"
                />
                <span class="text-sm font-mono w-8">{{ localConfig.threshold }}</span>
              </div>
            </div>
            <div>
              <label class="block text-xs font-bold uppercase tracking-widest text-gray-500 mb-2">
                Optimize For
              </label>
              <select
                v-model="localConfig.optimizeFor"
                @change="updateConfig"
                class="w-full border-2 border-black px-3 py-2 text-sm font-mono focus:outline-none"
              >
                <option value="cost">Cost</option>
                <option value="speed">Speed</option>
                <option value="accuracy">Accuracy</option>
              </select>
            </div>
            <div v-if="canEditMaxCost">
              <label class="block text-xs font-bold uppercase tracking-widest text-gray-500 mb-2">
                Max Cost/Month ($)
              </label>
              <input
                type="number"
                v-model.number="localConfig.maxCost"
                @change="updateConfig"
                class="w-full border-2 border-black px-3 py-2 text-sm font-mono focus:outline-none"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

interface AIUsage {
  callsUsed: number
  callsLimit: number
  callsPercent: number
  costUsed: number
  costLimit: number
  costPercent: number
  inputTokens: number
  outputTokens: number
  modelType: string
}

interface Alert {
  id: string
  severity: 'error' | 'warning' | 'info'
  message: string
}

interface Props {
  plan?: string
}

const props = withDefaults(defineProps<Props>(), {
  plan: 'radar'
})

const emit = defineEmits<{
  (e: 'upgrade'): void
  (e: 'config-changed', config: { threshold: number; optimizeFor: string; maxCost: number }): void
}>()

// State
const loading = ref(true)
const showBreakdown = ref(false)
const showConfig = ref(false)
const usage = ref<AIUsage>({
  callsUsed: 0,
  callsLimit: 0,
  callsPercent: 0,
  costUsed: 0,
  costLimit: 0,
  costPercent: 0,
  inputTokens: 0,
  outputTokens: 0,
  modelType: 'none'
})

const localConfig = ref({
  threshold: 0.7,
  optimizeFor: 'cost',
  maxCost: 100
})

const alerts = ref<Alert[]>([])
let refreshInterval: ReturnType<typeof setInterval> | null = null

// Computed
const periodDates = computed(() => {
  const now = new Date()
  const start = new Date(now.getFullYear(), now.getMonth(), 1)
  const end = new Date(now.getFullYear(), now.getMonth() + 1, 0)
  return `${start.toLocaleDateString()} - ${end.toLocaleDateString()}`
})

const modelDisplayName = computed(() => {
  const names: Record<string, string> = {
    'haiku-4.5': 'Haiku 4.5',
    'sonnet-4.5': 'Sonnet 4.5',
    'opus-4.5': 'Opus 4.5',
    'none': 'AI Disabled'
  }
  return names[usage.value.modelType] || 'Unknown'
})

const modelBadgeClass = computed(() => {
  const classes: Record<string, string> = {
    'haiku-4.5': 'bg-green-100 text-green-800 border-green-600',
    'sonnet-4.5': 'bg-blue-100 text-blue-800 border-blue-600',
    'opus-4.5': 'bg-purple-100 text-purple-800 border-purple-600',
    'none': 'bg-gray-100 text-gray-600 border-gray-400'
  }
  return classes[usage.value.modelType] || 'bg-gray-100 text-gray-600 border-gray-400'
})

const showUpgradePrompt = computed(() => {
  return usage.value.callsPercent > 80 || usage.value.costPercent > 80
})

const canEditMaxCost = computed(() => {
  return props.plan === 'orbit'
})

const marginCost = computed(() => usage.value.costUsed * 0.15)
const totalCost = computed(() => usage.value.costUsed + marginCost.value)

const forecastTotal = computed(() => {
  const daysPassed = new Date().getDate()
  if (daysPassed === 0) return 0
  const daysInMonth = new Date(new Date().getFullYear(), new Date().getMonth() + 1, 0).getDate()
  const dailyRate = totalCost.value / daysPassed
  return dailyRate * daysInMonth
})

const forecastWarning = computed(() => {
  if (forecastTotal.value > usage.value.costLimit * 0.9) {
    return 'You may exceed your monthly budget at this rate'
  }
  return null
})

// Methods
const fetchUsage = async () => {
  try {
    const response = await fetch('/api/v1/me/usage/ai', {
      credentials: 'include'
    })

    if (!response.ok) {
      throw new Error('Failed to fetch AI usage')
    }

    const data = await response.json()

    usage.value = {
      callsUsed: data.calls_used || 0,
      callsLimit: data.calls_limit || 0,
      callsPercent: data.calls_limit > 0 ? (data.calls_used / data.calls_limit) * 100 : 0,
      costUsed: data.cost_used || 0,
      costLimit: data.cost_limit || 0,
      costPercent: data.cost_limit > 0 ? (data.cost_used / data.cost_limit) * 100 : 0,
      inputTokens: data.input_tokens || 0,
      outputTokens: data.output_tokens || 0,
      modelType: data.model_type || 'none'
    }

    updateAlerts()
  } catch (error) {
    console.error('Failed to fetch AI usage:', error)
  } finally {
    loading.value = false
  }
}

const updateAlerts = () => {
  alerts.value = []

  if (usage.value.callsPercent >= 100) {
    alerts.value.push({
      id: 'calls-exceeded',
      severity: 'error',
      message: 'You have reached your AI call limit for this month'
    })
  } else if (usage.value.callsPercent >= 90) {
    alerts.value.push({
      id: 'calls-warning',
      severity: 'warning',
      message: 'You have used 90% of your AI calls for this month'
    })
  }

  if (usage.value.costPercent >= 100) {
    alerts.value.push({
      id: 'cost-exceeded',
      severity: 'error',
      message: 'You have reached your AI budget for this month'
    })
  } else if (usage.value.costPercent >= 80) {
    alerts.value.push({
      id: 'cost-warning',
      severity: 'warning',
      message: 'You have used 80% of your AI budget for this month'
    })
  }
}

const getUsageBarClass = (percent: number) => {
  if (percent >= 100) return 'bg-red-600'
  if (percent >= 80) return 'bg-yellow-500'
  return 'bg-black'
}

const alertClass = (severity: string) => {
  const classes: Record<string, string> = {
    'error': 'bg-red-100 text-red-800 border-red-600',
    'warning': 'bg-yellow-100 text-yellow-800 border-yellow-600',
    'info': 'bg-blue-100 text-blue-800 border-blue-600'
  }
  return classes[severity] || classes.info
}

const formatNumber = (num: number) => {
  return new Intl.NumberFormat().format(num)
}

const updateConfig = async () => {
  try {
    await fetch('/api/v1/me/ai/config', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        analysis_threshold: localConfig.value.threshold,
        optimize_for: localConfig.value.optimizeFor,
        max_cost: localConfig.value.maxCost
      })
    })
    emit('config-changed', localConfig.value)
  } catch (error) {
    console.error('Failed to update AI config:', error)
  }
}

// Lifecycle
onMounted(() => {
  fetchUsage()
  // Refresh every 5 minutes
  refreshInterval = setInterval(fetchUsage, 5 * 60 * 1000)
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>
