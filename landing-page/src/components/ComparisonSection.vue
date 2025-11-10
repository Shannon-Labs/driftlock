<template>
  <section id="comparison" class="section-padding bg-white dark:bg-gray-900">
    <div class="container-padding mx-auto">
      <!-- Section Header -->
      <div class="text-center mb-16">
        <h2 class="text-4xl md:text-5xl font-bold mb-6">
          <span class="gradient-text">Glass-Box vs Black-Box</span>
          <span class="text-gray-900 dark:text-white block">The Clear Winner</span>
        </h2>
        <p class="text-xl text-gray-600 dark:text-gray-300 max-w-3xl mx-auto">
          Compare Driftlock's explainable approach against traditional black-box algorithmic solutions.
        </p>
      </div>

      <!-- Comparison Table -->
      <div class="mb-16 overflow-x-auto">
        <div class="min-w-[800px] bg-gray-50 dark:bg-gray-800 rounded-2xl p-6 md:p-8">
          <table class="w-full">
            <thead>
              <tr class="border-b border-gray-200 dark:border-gray-700">
                <th class="text-left py-4 px-4 text-gray-900 dark:text-white font-bold">Feature</th>
                <th class="text-center py-4 px-4">
                  <div class="flex flex-col items-center">
                    <div class="w-12 h-12 bg-primary-100 dark:bg-primary-900/30 rounded-lg flex items-center justify-center mb-2">
                      <Shield class="w-6 h-6 text-primary-600" />
                    </div>
                    <span class="font-bold text-primary-600">Driftlock</span>
                  </div>
                </th>
                <th class="text-center py-4 px-4">
                  <div class="flex flex-col items-center">
                    <div class="w-12 h-12 bg-gray-200 dark:bg-gray-700 rounded-lg flex items-center justify-center mb-2">
                      <EyeOff class="w-6 h-6 text-gray-500" />
                    </div>
                    <span class="font-bold text-gray-600 dark:text-gray-400">Black-Box Algorithms</span>
                  </div>
                </th>
              </tr>
            </thead>
            <tbody>
              <tr class="border-b border-gray-200 dark:border-gray-700">
                <td class="py-6 px-4 font-medium text-gray-900 dark:text-white">Explainability</td>
                <td class="py-6 px-4 text-center">
                  <div class="flex items-center justify-center gap-2 text-green-600 dark:text-green-400">
                    <CheckCircle class="w-5 h-5" />
                    <span>Glass-box</span>
                  </div>
                </td>
                <td class="py-6 px-4 text-center">
                  <div class="flex items-center justify-center gap-2 text-red-600 dark:text-red-400">
                    <X class="w-5 h-5" />
                    <span>Black-box</span>
                  </div>
                </td>
              </tr>

              <tr class="border-b border-gray-200 dark:border-gray-700">
                <td class="py-6 px-4 font-medium text-gray-900 dark:text-white">Algorithm</td>
                <td class="py-6 px-4 text-center">
                  <span class="px-3 py-1 bg-primary-100 dark:bg-primary-900/30 text-primary-700 dark:text-primary-400 rounded-full text-sm">CBAD</span>
                </td>
                <td class="py-6 px-4 text-center">
                  <span class="px-3 py-1 bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-400 rounded-full text-sm">Neural Network</span>
                </td>
              </tr>

              <tr class="border-b border-gray-200 dark:border-gray-700">
                <td class="py-6 px-4 font-medium text-gray-900 dark:text-white">Performance</td>
                <td class="py-6 px-4 text-center">
                  <div class="flex items-center justify-center gap-2 text-green-600 dark:text-green-400">
                    <Zap class="w-5 h-5" />
                    <span>Sub-second</span>
                  </div>
                </td>
                <td class="py-6 px-4 text-center">
                  <span class="text-gray-600 dark:text-gray-400">Minutes</span>
                </td>
              </tr>

              <tr class="border-b border-gray-200 dark:border-gray-700">
                <td class="py-6 px-4 font-medium text-gray-900 dark:text-white">Storage</td>
                <td class="py-6 px-4 text-center">
                  <div class="flex items-center justify-center gap-2 text-green-600 dark:text-green-400">
                    <TrendingDown class="w-5 h-5" />
                    <span>10x efficient</span>
                  </div>
                </td>
                <td class="py-6 px-4 text-center">
                  <span class="text-gray-600 dark:text-gray-400">Standard</span>
                </td>
              </tr>

              <tr class="border-b border-gray-200 dark:border-gray-700">
                <td class="py-6 px-4 font-medium text-gray-900 dark:text-white">Integration</td>
                <td class="py-6 px-4 text-center">
                  <div class="flex items-center justify-center gap-2 text-green-600 dark:text-green-400">
                    <Plug class="w-5 h-5" />
                    <span>Native OpenTelemetry</span>
                  </div>
                </td>
                <td class="py-6 px-4 text-center">
                  <span class="text-gray-600 dark:text-gray-400">Custom</span>
                </td>
              </tr>

              <tr class="border-b border-gray-200 dark:border-gray-700">
                <td class="py-6 px-4 font-medium text-gray-900 dark:text-white">Compliance</td>
                <td class="py-6 px-4 text-center">
                  <div class="flex items-center justify-center gap-2 text-green-600 dark:text-green-400">
                    <CheckCircle class="w-5 h-5" />
                    <span>Built-in</span>
                  </div>
                </td>
                <td class="py-6 px-4 text-center">
                  <div class="flex items-center justify-center gap-2 text-orange-600 dark:text-orange-400">
                    <AlertTriangle class="w-5 h-5" />
                    <span>Add-on</span>
                  </div>
                </td>
              </tr>

              <tr>
                <td class="py-6 px-4 font-medium text-gray-900 dark:text-white">DORA Ready</td>
                <td class="py-6 px-4 text-center">
                  <div class="flex items-center justify-center gap-2 text-green-600 dark:text-green-400">
                    <CheckCircle class="w-5 h-5" />
                    <span>Yes</span>
                  </div>
                </td>
                <td class="py-6 px-4 text-center">
                  <div class="flex items-center justify-center gap-2 text-red-600 dark:text-red-400">
                    <X class="w-5 h-5" />
                    <span>No</span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Regulatory Compliance Comparison -->
      <div class="grid lg:grid-cols-2 gap-12 mb-16">
        <!-- DORA Compliance Matrix -->
        <div class="bg-gradient-to-br from-green-50 to-emerald-50 dark:from-green-900/20 dark:to-emerald-900/20 rounded-2xl p-8">
          <div class="flex items-center gap-3 mb-6">
            <CheckCircle class="w-8 h-8 text-green-600" />
            <h3 class="text-2xl font-bold text-gray-900 dark:text-white">DORA Compliance Matrix</h3>
          </div>

          <div class="space-y-4">
            <div class="bg-white dark:bg-gray-800 rounded-lg p-4">
              <div class="flex items-center justify-between mb-2">
                <span class="font-medium text-gray-900 dark:text-white">Explainable Algorithms</span>
                <span class="px-2 py-1 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 rounded text-sm">100% Compliant</span>
              </div>
              <p class="text-sm text-gray-600 dark:text-gray-400">Mathematical explanations for every decision</p>
            </div>

            <div class="bg-white dark:bg-gray-800 rounded-lg p-4">
              <div class="flex items-center justify-between mb-2">
                <span class="font-medium text-gray-900 dark:text-white">Audit Trails</span>
                <span class="px-2 py-1 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 rounded text-sm">Automatic</span>
              </div>
              <p class="text-sm text-gray-600 dark:text-gray-400">Immutable logs with tamper-proof timestamps</p>
            </div>

            <div class="bg-white dark:bg-gray-800 rounded-lg p-4">
              <div class="flex items-center justify-between mb-2">
                <span class="font-medium text-gray-900 dark:text-white">Documentation</span>
                <span class="px-2 py-1 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 rounded text-sm">One-click</span>
              </div>
              <p class="text-sm text-gray-600 dark:text-gray-400">Regulator-ready reports in seconds</p>
            </div>

            <div class="bg-white dark:bg-gray-800 rounded-lg p-4">
              <div class="flex items-center justify-between mb-2">
                <span class="font-medium text-gray-900 dark:text-white">Testing & Validation</span>
                <span class="px-2 py-1 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 rounded text-sm">Built-in</span>
              </div>
              <p class="text-sm text-gray-600 dark:text-gray-400">Statistical significance testing included</p>
            </div>
          </div>
        </div>

        <!-- Risk Assessment -->
        <div class="bg-gradient-to-br from-red-50 to-orange-50 dark:from-red-900/20 dark:to-orange-900/20 rounded-2xl p-8">
          <div class="flex items-center gap-3 mb-6">
            <AlertTriangle class="w-8 h-8 text-red-600" />
            <h3 class="text-2xl font-bold text-gray-900 dark:text-white">Black-Box Risk Assessment</h3>
          </div>

          <div class="space-y-4">
            <div class="bg-white dark:bg-gray-800 rounded-lg p-4 border-l-4 border-red-500">
              <div class="flex items-center justify-between mb-2">
                <span class="font-medium text-gray-900 dark:text-white">Regulatory Rejection</span>
                <span class="px-2 py-1 bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400 rounded text-sm">High Risk</span>
              </div>
              <p class="text-sm text-gray-600 dark:text-gray-400">Auditors automatically fail non-explainable systems</p>
            </div>

            <div class="bg-white dark:bg-gray-800 rounded-lg p-4 border-l-4 border-red-500">
              <div class="flex items-center justify-between mb-2">
                <span class="font-medium text-gray-900 dark:text-white">Financial Penalties</span>
                <span class="px-2 py-1 bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400 rounded text-sm">€50M-€200M</span>
              </div>
              <p class="text-sm text-gray-600 dark:text-gray-400">Fines per violation starting January 2025</p>
            </div>

            <div class="bg-white dark:bg-gray-800 rounded-lg p-4 border-l-4 border-orange-500">
              <div class="flex items-center justify-between mb-2">
                <span class="font-medium text-gray-900 dark:text-white">Operational Risk</span>
                <span class="px-2 py-1 bg-orange-100 dark:bg-orange-900/30 text-orange-700 dark:text-orange-400 rounded text-sm">Medium</span>
              </div>
              <p class="text-sm text-gray-600 dark:text-gray-400">Cannot debug or improve model performance</p>
            </div>

            <div class="bg-white dark:bg-gray-800 rounded-lg p-4 border-l-4 border-orange-500">
              <div class="flex items-center justify-between mb-2">
                <span class="font-medium text-gray-900 dark:text-white">Reputation Risk</span>
                <span class="px-2 py-1 bg-orange-100 dark:bg-orange-900/30 text-orange-700 dark:text-orange-400 rounded text-sm">Medium</span>
              </div>
              <p class="text-sm text-gray-600 dark:text-gray-400">Customer trust issues with opaque decisions</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Key Differentiators -->
      <div class="bg-gradient-to-r from-primary-50 to-cyan-50 dark:from-primary-900/20 dark:to-cyan-900/20 rounded-2xl p-8 md:p-12">
        <h3 class="text-2xl font-bold text-center text-gray-900 dark:text-white mb-12">Why Driftlock Wins Every Time</h3>

        <div class="grid md:grid-cols-3 gap-8">
          <div class="text-center">
            <div class="w-20 h-20 bg-white dark:bg-gray-800 rounded-2xl flex items-center justify-center mx-auto mb-6 shadow-lg">
              <Eye class="w-10 h-10 text-primary-600" />
            </div>
            <h4 class="text-xl font-bold text-gray-900 dark:text-white mb-4">100% Transparency</h4>
            <p class="text-gray-600 dark:text-gray-400">
              Every decision can be explained with mathematical proof. No "trust me, it works" - just show the math.
            </p>
          </div>

          <div class="text-center">
            <div class="w-20 h-20 bg-white dark:bg-gray-800 rounded-2xl flex items-center justify-center mx-auto mb-6 shadow-lg">
              <Scale class="w-10 h-10 text-primary-600" />
            </div>
            <h4 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Regulatory Proof</h4>
            <p class="text-gray-600 dark:text-gray-400">
              Built specifically for DORA compliance. Audit trails, statistical significance, and regulator-ready reports.
            </p>
          </div>

          <div class="text-center">
            <div class="w-20 h-20 bg-white dark:bg-gray-800 rounded-2xl flex items-center justify-center mx-auto mb-6 shadow-lg">
              <TrendingUp class="w-10 h-10 text-primary-600" />
            </div>
            <h4 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Better Performance</h4>
            <p class="text-gray-600 dark:text-gray-400">
              Faster processing, less storage, higher accuracy, and no training data required. Just plug and play.
            </p>
          </div>
        </div>
      </div>

      <!-- Bottom CTA -->
      <div class="text-center mt-16">
        <div class="inline-flex items-center gap-2 px-6 py-3 bg-primary-100 dark:bg-primary-900/30 text-primary-700 dark:text-primary-400 rounded-full mb-6">
          <Shield class="w-5 h-5" />
          <span class="font-semibold">The Choice is Clear</span>
        </div>
        <p class="text-xl text-gray-600 dark:text-gray-300 mb-8">
          Don't risk €50M fines. Choose explainable algorithms that regulators will actually approve.
        </p>
        <a href="#cta" class="btn-primary">
          Get Started Today
          <ArrowRight class="w-5 h-5 ml-2" />
        </a>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { Shield, EyeOff, CheckCircle, X, Zap, TrendingDown, Plug, AlertTriangle, Eye, Scale, TrendingUp, ArrowRight } from 'lucide-vue-next'
</script>

<style scoped>
.btn-primary {
  @apply inline-flex items-center bg-primary-600 hover:bg-primary-700 text-white px-8 py-3 rounded-lg font-semibold transition-colors;
}
</style>