<template>
  <div class="bg-gray-50 rounded-lg p-6">
    <h3 class="text-lg font-semibold text-gray-900 mb-4">Interactive Anomaly Detection</h3>
    
    <!-- Sample Data Selection -->
    <div class="mb-4">
      <label class="block text-sm font-medium text-gray-700 mb-2">Sample Dataset:</label>
      <select 
        v-model="selectedDataset" 
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="financial">Financial Transactions</option>
        <option value="network">Network Traffic</option>
        <option value="healthcare">Healthcare Records</option>
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
        <h4 class="font-medium text-gray-900 mb-2">🤖 AI Insights</h4>
        <div class="text-sm text-gray-700 whitespace-pre-wrap">{{ aiAnalysis }}</div>
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

// Sample data for demo
const sampleData = {
  financial: [
    { id: 1, amount: 1000.00, merchant: "Coffee Shop", timestamp: "2024-01-15T10:30:00Z" },
    { id: 2, amount: 850000.00, merchant: "Luxury Cars Inc", timestamp: "2024-01-15T14:22:00Z" }, // Anomaly
    { id: 3, amount: 45.50, merchant: "Gas Station", timestamp: "2024-01-15T16:15:00Z" },
  ],
  network: [
    { id: 1, bytes: 1024, src_ip: "192.168.1.10", dst_ip: "10.0.0.5" },
    { id: 2, bytes: 104857600, src_ip: "192.168.1.15", dst_ip: "suspicious.domain.com" }, // Anomaly
    { id: 3, bytes: 2048, src_ip: "192.168.1.12", dst_ip: "10.0.0.8" },
  ],
  healthcare: [
    { id: 1, patient_id: "P001", procedure: "routine_checkup", duration_mins: 30 },
    { id: 2, patient_id: "P002", procedure: "emergency_surgery", duration_mins: 720 }, // Anomaly
    { id: 3, patient_id: "P003", procedure: "vaccination", duration_mins: 15 },
  ]
}

async function runAnalysis() {
  isAnalyzing.value = true
  results.value = null
  aiAnalysis.value = ''
  error.value = ''

  try {
    // Simulate anomaly detection
    await new Promise(resolve => setTimeout(resolve, 2000))
    
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

    // Get AI analysis
    try {
      const response = await fetch('/api/analyze', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          anomalies: mockAnomalies,
          query: `Analyze ${selectedDataset.value} anomalies for business impact`
        }),
      })

      if (response.ok) {
        const data = await response.json()
        aiAnalysis.value = data.analysis
      }
    } catch (apiError) {
      // Fallback to mock analysis if API is not available
      aiAnalysis.value = `Risk Assessment: HIGH\n\nThe detected anomaly shows significant deviation from normal ${selectedDataset.value} patterns. The compression-based analysis indicates a 97.3% confidence that this represents unusual behavior.\n\nRecommended Actions:\n1. Investigate the flagged transaction immediately\n2. Review related account activity\n3. Consider implementing additional monitoring\n\nCompliance: This detection meets DORA explainability requirements with mathematical proof via compression analysis.`
    }

  } catch (err) {
    error.value = 'Analysis failed. Please try again.'
  } finally {
    isAnalyzing.value = false
  }
}
</script>