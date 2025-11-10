<template>
  <section id="proof" class="section-padding bg-gray-50 dark:bg-gray-800">
    <div class="container-padding mx-auto">
      <!-- Section Header -->
      <div class="text-center mb-16">
        <h2 class="text-4xl md:text-5xl font-bold mb-6">
          <span class="text-gray-900 dark:text-white">See It Work</span>
          <span class="gradient-text block">In 30 Seconds</span>
        </h2>
        <p class="text-xl text-gray-600 dark:text-gray-300 max-w-3xl mx-auto">
          Our working prototype demonstrates explainable fraud detection with real-time processing and audit-ready output.
        </p>
      </div>

      <!-- Demo Integration -->
      <div class="mb-16">
        <div class="bg-white dark:bg-gray-900 rounded-2xl shadow-xl overflow-hidden">
          <!-- Demo Header -->
          <div class="bg-gradient-to-r from-primary-600 to-cyan-600 p-6 text-white">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-3">
                <Terminal class="w-6 h-6" />
                <span class="font-bold text-lg">DriftLock Demo</span>
              </div>
              <div class="flex items-center gap-4">
                <div class="flex items-center gap-2">
                  <div class="w-3 h-3 bg-green-400 rounded-full animate-pulse"></div>
                  <span class="text-sm">Ready</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Demo Content -->
          <div class="p-6">
            <!-- Quick Demo Stats -->
            <div class="grid md:grid-cols-4 gap-6 mb-8">
              <div class="text-center">
                <div class="text-2xl font-bold text-primary-600 mb-1">5,000</div>
                <div class="text-sm text-gray-600 dark:text-gray-400">Test Transactions</div>
              </div>
              <div class="text-center">
                <div class="text-2xl font-bold text-primary-600 mb-1">2,000</div>
                <div class="text-sm text-gray-600 dark:text-gray-400">Processed Events</div>
              </div>
              <div class="text-center">
                <div class="text-2xl font-bold text-orange-500 mb-1">~30</div>
                <div class="text-sm text-gray-600 dark:text-gray-400">Anomalies Detected</div>
              </div>
              <div class="text-center">
                <div class="text-2xl font-bold text-primary-600 mb-1">4-6s</div>
                <div class="text-sm text-gray-600 dark:text-gray-400">Processing Time</div>
              </div>
            </div>

            <!-- Interactive Demo Preview -->
            <div class="bg-gray-50 dark:bg-gray-800 rounded-xl p-6 mb-8">
              <h4 class="font-bold text-lg text-gray-900 dark:text-white mb-4">What Compliance Officers See</h4>

              <!-- Simulated Anomaly Card -->
              <div class="bg-white dark:bg-gray-900 rounded-lg border border-red-200 dark:border-red-800 p-6">
                <div class="flex items-start justify-between mb-4">
                  <div class="flex items-center gap-3">
                    <div class="w-3 h-3 bg-red-500 rounded-full"></div>
                    <span class="font-bold text-red-600 dark:text-red-400">High Risk Anomaly Detected</span>
                  </div>
                  <span class="text-sm text-gray-500">ID: TXN-2024-1127-001</span>
                </div>

                <div class="space-y-4">
                  <div class="grid md:grid-cols-2 gap-4">
                    <div>
                      <h5 class="font-semibold text-gray-900 dark:text-white mb-2">Regulatory Explanation</h5>
                      <p class="text-sm text-gray-600 dark:text-gray-400">
                        This payment is <strong>47x larger</strong> than customer's typical transactions and shows
                        <strong>8000ms latency</strong> (vs 50-100ms baseline). Pattern deviation indicates potential fraud.
                      </p>
                    </div>
                    <div>
                      <h5 class="font-semibold text-gray-900 dark:text-white mb-2">Mathematical Proof</h5>
                      <div class="bg-gray-50 dark:bg-gray-800 rounded p-3 font-mono text-xs">
                        <div>NCD Score: 0.847</div>
                        <div>P-value: 0.002</div>
                        <div>Statistical Significance: 99.8%</div>
                        <div>Compression Ratio: 1:47</div>
                      </div>
                    </div>
                  </div>

                  <div>
                    <h5 class="font-semibold text-gray-900 dark:text-white mb-2">Comparison to Normal Transactions</h5>
                    <div class="grid grid-cols-3 gap-2">
                      <div class="bg-green-50 dark:bg-green-900/20 rounded p-2 text-center text-xs">
                        <div class="text-green-600 dark:text-green-400">Normal</div>
                        <div>$127.50</div>
                        <div>67ms</div>
                      </div>
                      <div class="bg-green-50 dark:bg-green-900/20 rounded p-2 text-center text-xs">
                        <div class="text-green-600 dark:text-green-400">Normal</div>
                        <div>$89.99</div>
                        <div>52ms</div>
                      </div>
                      <div class="bg-red-50 dark:bg-red-900/20 rounded p-2 text-center text-xs border-2 border-red-200">
                        <div class="text-red-600 dark:text-red-400">This Tx</div>
                        <div>$5,992.50</div>
                        <div>8,047ms</div>
                      </div>
                    </div>
                  </div>

                  <div class="flex justify-end gap-3">
                    <button class="px-4 py-2 bg-primary-600 text-white rounded-lg text-sm hover:bg-primary-700 transition-colors">
                      Export Audit Trail
                    </button>
                    <button class="px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg text-sm hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors">
                      View Similar Cases
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <!-- Demo Commands -->
            <div class="bg-gray-900 rounded-lg p-6 mb-6">
              <h4 class="font-bold text-white mb-4">Try It Yourself</h4>
              <div class="space-y-3">
                <div class="bg-black rounded p-4 font-mono text-sm">
                  <div class="flex items-center gap-2 mb-2">
                    <span class="text-gray-500">$</span>
                    <span class="text-green-400">git clone https://github.com/Shannon-Labs/driftlock.git</span>
                  </div>
                  <div class="flex items-center gap-2 mb-2">
                    <span class="text-gray-500">$</span>
                    <span class="text-green-400">cd driftlock && make demo</span>
                  </div>
                  <div class="flex items-center gap-2 mb-2">
                    <span class="text-gray-500">$</span>
                    <span class="text-green-400">./driftlock-demo test-data/financial-demo.json</span>
                  </div>
                  <div class="flex items-center gap-2">
                    <span class="text-gray-500">$</span>
                    <span class="text-green-400">open demo-output.html</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- What This Demo Proves -->
      <div class="mb-16">
        <h3 class="text-2xl font-bold text-center text-gray-900 dark:text-white mb-12">What This Demo Proves</h3>

        <div class="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
          <div class="bg-white dark:bg-gray-900 rounded-xl p-6 text-center">
            <div class="w-12 h-12 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center mx-auto mb-4">
              <CheckCircle class="w-6 h-6 text-green-600" />
            </div>
            <h4 class="font-bold text-gray-900 dark:text-white mb-2">Core Algorithm Works</h4>
            <p class="text-sm text-gray-600 dark:text-gray-400">Compression distance + statistical testing detects real anomalies</p>
          </div>

          <div class="bg-white dark:bg-gray-900 rounded-xl p-6 text-center">
            <div class="w-12 h-12 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center mx-auto mb-4">
              <CheckCircle class="w-6 h-6 text-green-600" />
            </div>
            <h4 class="font-bold text-gray-900 dark:text-white mb-2">Explainable Output</h4>
            <p class="text-sm text-gray-600 dark:text-gray-400">Audit-friendly explanations with mathematical proof</p>
          </div>

          <div class="bg-white dark:bg-gray-900 rounded-xl p-6 text-center">
            <div class="w-12 h-12 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center mx-auto mb-4">
              <CheckCircle class="w-6 h-6 text-green-600" />
            </div>
            <h4 class="font-bold text-gray-900 dark:text-white mb-2">Fast Processing</h4>
            <p class="text-sm text-gray-600 dark:text-gray-400">2,000 transactions processed in 4-6 seconds</p>
          </div>

          <div class="bg-white dark:bg-gray-900 rounded-xl p-6 text-center">
            <div class="w-12 h-12 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center mx-auto mb-4">
              <CheckCircle class="w-6 h-6 text-green-600" />
            </div>
            <h4 class="font-bold text-gray-900 dark:text-white mb-2">Glass Box</h4>
            <p class="text-sm text-gray-600 dark:text-gray-400">Compression ratios, p-values, and statistical significance</p>
          </div>
        </div>
      </div>

      <!-- Demo Data Information -->
      <div class="bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-blue-900/20 dark:to-indigo-900/20 rounded-2xl p-8">
        <div class="grid md:grid-cols-2 gap-8 items-center">
          <div>
            <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Demo Data Specifications</h3>
            <div class="space-y-3 text-gray-600 dark:text-gray-400">
              <div class="flex items-start gap-3">
                <Database class="w-5 h-5 text-primary-600 mt-0.5 flex-shrink-0" />
                <div>
                  <strong class="text-gray-900 dark:text-white">5,000 synthetic payment transactions</strong>
                  <p class="text-sm">Realistic payment gateway telemetry data</p>
                </div>
              </div>
              <div class="flex items-start gap-3">
                <Zap class="w-5 h-5 text-primary-600 mt-0.5 flex-shrink-0" />
                <div>
                  <strong class="text-gray-900 dark:text-white">Normal patterns</strong>
                  <p class="text-sm">50-100ms processing, US/UK origins, /v1/charges endpoint</p>
                </div>
              </div>
              <div class="flex items-start gap-3">
                <AlertTriangle class="w-5 h-5 text-primary-600 mt-0.5 flex-shrink-0" />
                <div>
                  <strong class="text-gray-900 dark:text-white">Anomalies</strong>
                  <p class="text-sm">Latency spikes up to 8000ms and malformed endpoints</p>
                </div>
              </div>
              <div class="flex items-start gap-3">
                <Target class="w-5 h-5 text-primary-600 mt-0.5 flex-shrink-0" />
                <div>
                  <strong class="text-gray-900 dark:text-white">Detection rate</strong>
                  <p class="text-sm">Tuned to flag ~30 anomalies from 2,000 events (~1.5% detection rate)</p>
                </div>
              </div>
            </div>
          </div>

          <div class="bg-white dark:bg-gray-800 rounded-xl p-6">
            <h4 class="font-bold text-gray-900 dark:text-white mb-4">Ready for Production</h4>
            <p class="text-gray-600 dark:text-gray-400 mb-6">
              This demo proves our algorithm works. In production, DriftLock:
            </p>
            <ul class="space-y-2 text-gray-600 dark:text-gray-400">
              <li class="flex items-start gap-2">
                <ArrowRight class="w-4 h-4 text-primary-600 mt-1 flex-shrink-0" />
                <span>Integrates without code changes</span>
              </li>
              <li class="flex items-start gap-2">
                <ArrowRight class="w-4 h-4 text-primary-600 mt-1 flex-shrink-0" />
                <span>Detects fraud in real-time</span>
              </li>
              <li class="flex items-start gap-2">
                <ArrowRight class="w-4 h-4 text-primary-600 mt-1 flex-shrink-0" />
                <span>Alerts your team immediately</span>
              </li>
              <li class="flex items-start gap-2">
                <ArrowRight class="w-4 h-4 text-primary-600 mt-1 flex-shrink-0" />
                <span>Stores immutable audit trails</span>
              </li>
              <li class="flex items-start gap-2">
                <ArrowRight class="w-4 h-4 text-primary-600 mt-1 flex-shrink-0" />
                <span>Generates compliance reports</span>
              </li>
            </ul>
          </div>
        </div>
      </div>

      <!-- Bottom CTA -->
      <div class="text-center mt-16">
        <div class="inline-flex items-center gap-2 px-6 py-3 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 rounded-full mb-6">
          <CheckCircle class="w-5 h-5" />
          <span class="font-semibold">Production Ready</span>
        </div>
        <p class="text-xl text-gray-600 dark:text-gray-300 mb-8">
          Ready to pilot on your real payment data?
        </p>
        <a href="#cta" class="btn-primary">
          Schedule Pilot
          <ArrowRight class="w-5 h-5 ml-2" />
        </a>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { Terminal, CheckCircle, Database, Zap, AlertTriangle, Target, ArrowRight } from 'lucide-vue-next'
</script>

<style scoped>
.btn-primary {
  @apply inline-flex items-center bg-primary-600 hover:bg-primary-700 text-white px-8 py-3 rounded-lg font-semibold transition-colors;
}
</style>