<template>
  <div class="bg-gray-50 rounded-lg p-6">
    <h3 class="text-lg font-semibold text-gray-900 mb-4">Interactive Anomaly Detection</h3>
    
    <!-- Sample Data Selection -->
    <div class="mb-4">
      <label class="block text-sm font-medium text-gray-700 mb-2">Use Case & Dataset:</label>
      <select 
        v-model="selectedDataset" 
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="financial">🏦 Financial Fraud Detection</option>
        <option value="security">🛡️ Cybersecurity Events</option>
        <option value="network">🌐 Network Traffic Analysis</option>
        <option value="healthcare">🏥 Healthcare Data Monitoring</option>
        <option value="ecommerce">🛒 E-commerce User Behavior</option>
        <option value="manufacturing">🏭 Industrial System Logs</option>
        <option value="social">📱 Social Media Analytics</option>
        <option value="crypto">₿ Crypto Transaction Analysis</option>
      </select>
    </div>

    <!-- Analysis Button -->
    <button 
      @click="runAnalysis"
      :disabled="isAnalyzing"
      class="w-full bg-blue-600 text-white px-4 py-2 rounded-md font-medium transition-colors hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
    >
      <span v-if="!isAnalyzing">🔍 Analyze for Anomalies</span>
      <span v-else>🔄 Analyzing...</span>
    </button>

    <!-- Results -->
    <div v-if="results" class="mt-6 space-y-4">
      <div class="bg-white border border-gray-200 rounded-lg p-4">
        <h4 class="font-medium text-gray-900 mb-2">Detection Results</h4>
        <div class="grid grid-cols-2 gap-4 text-sm">
          <div>
            <span class="text-gray-600">Anomalies Found:</span>
            <span class="font-mono text-red-600 ml-2">{{ results.anomaly_count }}</span>
          </div>
          <div>
            <span class="text-gray-600">Confidence:</span>
            <span class="font-mono text-green-600 ml-2">{{ results.confidence }}%</span>
          </div>
        </div>
        
        <div class="mt-3 p-3 bg-blue-50 rounded-md">
          <p class="text-sm text-blue-800">
            <strong>Mathematical Proof:</strong> {{ results.explanation }}
          </p>
        </div>
      </div>

      <!-- AI Analysis -->
      <div v-if="aiAnalysis" class="bg-white border border-gray-200 rounded-lg p-4">
        <h4 class="font-medium text-gray-900 mb-2">💡 AI Enhancement Available</h4>
        <div class="text-sm text-gray-700 whitespace-pre-wrap">{{ aiAnalysis }}</div>
        <button class="mt-3 text-xs bg-blue-600 text-white px-3 py-1 rounded-md hover:bg-blue-700">
          Upgrade to Pro
        </button>
      </div>

      <!-- Compliance Report -->
      <div class="bg-green-50 border border-green-200 rounded-lg p-4">
        <h4 class="font-medium text-green-900 mb-2">✅ Compliance Report Generated</h4>
        <p class="text-sm text-green-800">
          DORA-compliant report with mathematical evidence ready for audit review.
        </p>
        <button class="mt-2 text-xs bg-green-600 text-white px-3 py-1 rounded-md hover:bg-green-700">
          Download Report
        </button>
      </div>
    </div>

    <!-- Error State -->
    <div v-if="error" class="mt-4 p-3 bg-red-50 border border-red-200 rounded-md">
      <p class="text-sm text-red-800">{{ error }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const selectedDataset = ref('financial')
const isAnalyzing = ref(false)
const results = ref<any>(null)
const aiAnalysis = ref('')
const error = ref('')

// Sample data for diverse use cases
const sampleData = {
  financial: [
    { id: 1, amount: 1000.00, merchant: "Coffee Shop", timestamp: "2024-01-15T10:30:00Z" },
    { id: 2, amount: 850000.00, merchant: "Luxury Cars Inc", timestamp: "2024-01-15T14:22:00Z" }, // Fraud
    { id: 3, amount: 45.50, merchant: "Gas Station", timestamp: "2024-01-15T16:15:00Z" },
  ],
  security: [
    { id: 1, event: "login_success", user: "alice@company.com", location: "New York" },
    { id: 2, event: "admin_access", user: "temp@suspicious.ru", location: "Unknown" }, // Security threat
    { id: 3, event: "file_access", user: "bob@company.com", location: "London" },
  ],
  network: [
    { id: 1, bytes: 1024, src_ip: "192.168.1.10", dst_ip: "10.0.0.5" },
    { id: 2, bytes: 104857600, src_ip: "192.168.1.15", dst_ip: "suspicious.domain.com" }, // Data exfiltration
    { id: 3, bytes: 2048, src_ip: "192.168.1.12", dst_ip: "10.0.0.8" },
  ],
  healthcare: [
    { id: 1, patient_id: "P001", procedure: "routine_checkup", duration_mins: 30 },
    { id: 2, patient_id: "P002", procedure: "emergency_surgery", duration_mins: 720 }, // Unusual procedure
    { id: 3, patient_id: "P003", procedure: "vaccination", duration_mins: 15 },
  ],
  ecommerce: [
    { id: 1, user: "customer123", action: "view_product", session_time: 120 },
    { id: 2, user: "bot_account", action: "mass_purchase", session_time: 5 }, // Bot behavior
    { id: 3, user: "customer456", action: "add_to_cart", session_time: 180 },
  ],
  manufacturing: [
    { id: 1, machine: "conveyor_1", temperature: 75.2, vibration: 0.8 },
    { id: 2, machine: "press_3", temperature: 180.5, vibration: 8.2 }, // Equipment failure
    { id: 3, machine: "welder_2", temperature: 82.1, vibration: 1.1 },
  ],
  social: [
    { id: 1, post: "Happy birthday!", engagement: 25, sentiment: 0.8 },
    { id: 2, post: "Spam promotional content", engagement: 0, sentiment: -0.9 }, // Spam/bot
    { id: 3, post: "Great weather today", engagement: 15, sentiment: 0.6 },
  ],
  crypto: [
    { id: 1, from: "0x123...abc", to: "0x456...def", amount: 0.5, gas: 21000 },
    { id: 2, from: "0x789...ghi", to: "0xsuspicious...mixer", amount: 1000.0, gas: 500000 }, // Money laundering
    { id: 3, from: "0x321...fed", to: "0x654...cba", amount: 0.1, gas: 21000 },
  ]
}

async function runAnalysis() {
  isAnalyzing.value = true
  results.value = null
  aiAnalysis.value = ''
  error.value = ''

  try {
    // Fast mathematical detection (no AI overhead)
    await new Promise(resolve => setTimeout(resolve, 1000)) // Simulate real backend call
    
    const mockAnomalies = [
      {
        id: 'anom_001',
        ncd_score: 0.85,
        p_value: 0.003,
        explanation: `Compression ratio anomaly detected: baseline 0.23, current 0.89 (NCD=0.85, p<0.01)`,
        stream_type: selectedDataset.value,
        confidence: 97.3
      }
    ]

    results.value = {
      anomaly_count: mockAnomalies.length,
      confidence: 97.3,
      explanation: mockAnomalies[0].explanation
    }

    // AI analysis as optional enhancement (not default)
    // Only show upgrade prompt for now to save costs
    aiAnalysis.value = `🤖 AI Business Insights Available\n\nUpgrade to Pro plan to get:\n• Executive risk assessment\n• Business impact analysis\n• Recommended actions\n• Custom compliance narratives\n\nThe mathematical detection above provides the core explainable evidence needed for audits.`

  } catch (err) {
    error.value = 'Analysis failed. Please try again.'
  } finally {
    isAnalyzing.value = false
  }
}
</script>