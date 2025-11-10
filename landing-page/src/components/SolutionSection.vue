<template>
  <section id="solution" class="section-padding bg-white dark:bg-gray-900">
    <div class="container-padding mx-auto">
      <!-- Section Header -->
      <div class="text-center mb-16">
        <h2 class="text-4xl md:text-5xl font-bold mb-6">
          <span class="gradient-text">Glass-Box Solution</span>
          <span class="text-gray-900 dark:text-white block">That Auditors Love</span>
        </h2>
        <p class="text-xl text-gray-600 dark:text-gray-300 max-w-3xl mx-auto">
          Driftlock uses compression-based anomaly detection to provide mathematical proof for every flagged transaction.
        </p>
      </div>

      <!-- How It Works -->
      <div class="mb-20">
        <h3 class="text-3xl font-bold text-center text-gray-900 dark:text-white mb-12">How It Works</h3>

        <div class="grid md:grid-cols-4 gap-8">
          <div class="text-center group">
            <div class="relative mb-6">
              <div class="w-20 h-20 bg-primary-100 dark:bg-primary-900/30 rounded-full flex items-center justify-center mx-auto group-hover:scale-110 transition-transform">
                <Database class="w-10 h-10 text-primary-600" />
              </div>
              <div class="absolute -bottom-2 left-1/2 transform -translate-x-1/2 w-0.5 h-8 bg-primary-200 dark:bg-primary-800 md:hidden"></div>
            </div>
            <h4 class="text-lg font-bold text-gray-900 dark:text-white mb-3">1. Build Baseline</h4>
            <p class="text-gray-600 dark:text-gray-400">Analyzes normal transaction patterns (first ~400 events) to establish compression baseline</p>
          </div>

          <div class="text-center group">
            <div class="relative mb-6">
              <div class="w-20 h-20 bg-primary-100 dark:bg-primary-900/30 rounded-full flex items-center justify-center mx-auto group-hover:scale-110 transition-transform">
                <Search class="w-10 h-10 text-primary-600" />
              </div>
              <div class="absolute -bottom-2 left-1/2 transform -translate-x-1/2 w-0.5 h-8 bg-primary-200 dark:bg-primary-800 md:hidden"></div>
            </div>
            <h4 class="text-lg font-bold text-gray-900 dark:text-white mb-3">2. Detect Anomalies</h4>
            <p class="text-gray-600 dark:text-gray-400">Compares new transactions vs baseline using Normalized Compression Distance (NCD)</p>
          </div>

          <div class="text-center group">
            <div class="relative mb-6">
              <div class="w-20 h-20 bg-primary-100 dark:bg-primary-900/30 rounded-full flex items-center justify-center mx-auto group-hover:scale-110 transition-transform">
                <FileText class="w-10 h-10 text-primary-600" />
              </div>
              <div class="absolute -bottom-2 left-1/2 transform -translate-x-1/2 w-0.5 h-8 bg-primary-200 dark:bg-primary-800 md:hidden"></div>
            </div>
            <h4 class="text-lg font-bold text-gray-900 dark:text-white mb-3">3. Generate Proof</h4>
            <p class="text-gray-600 dark:text-gray-400">Creates mathematical explanations with NCD scores, p-values, and statistical significance</p>
          </div>

          <div class="text-center group">
            <div class="relative mb-6">
              <div class="w-20 h-20 bg-primary-100 dark:bg-primary-900/30 rounded-full flex items-center justify-center mx-auto group-hover:scale-110 transition-transform">
                <CheckCircle class="w-10 h-10 text-primary-600" />
              </div>
            </div>
            <h4 class="text-lg font-bold text-gray-900 dark:text-white mb-3">4. Audit Ready</h4>
            <p class="text-gray-600 dark:text-gray-400">One-click export of regulator-friendly reports with complete audit trails</p>
          </div>
        </div>

        <!-- Desktop Arrow Flow -->
        <div class="hidden md:flex items-center justify-center mt-8 mb-12">
          <ArrowRight class="w-6 h-6 text-primary-400" />
          <ArrowRight class="w-6 h-6 text-primary-400 mx-2" />
          <ArrowRight class="w-6 h-6 text-primary-400" />
        </div>
      </div>

      <!-- Technology Deep Dive -->
      <div class="grid lg:grid-cols-2 gap-12 mb-20">
        <!-- CBAD Algorithm -->
        <div class="bg-gradient-to-br from-blue-50 to-indigo-50 dark:from-blue-900/20 dark:to-indigo-900/20 rounded-2xl p-8">
          <div class="flex items-center gap-3 mb-6">
            <div class="w-12 h-12 bg-primary-600 rounded-lg flex items-center justify-center">
              <Cpu class="w-6 h-6 text-white" />
            </div>
            <h3 class="text-2xl font-bold text-gray-900 dark:text-white">CBAD Algorithm</h3>
          </div>

          <div class="space-y-6">
            <div>
              <h4 class="font-bold text-lg text-gray-900 dark:text-white mb-2">Normalized Compression Distance</h4>
              <p class="text-gray-600 dark:text-gray-400 mb-4">
                Uses compression mathematics to measure similarity between data patterns:
              </p>
              <div class="bg-white dark:bg-gray-800 rounded-lg p-4 font-mono text-sm">
                NCD(x,y) = C(xy) - min(C(x), C(y)) / max(C(x), C(y))
              </div>
            </div>

            <div>
              <h4 class="font-bold text-lg text-gray-900 dark:text-white mb-2">Why It Works Better</h4>
              <ul class="space-y-2 text-gray-600 dark:text-gray-400">
                <li class="flex items-start gap-2">
                  <Check class="w-5 h-5 text-green-500 mt-0.5 flex-shrink-0" />
                  <span>No training data required - works with any data format</span>
                </li>
                <li class="flex items-start gap-2">
                  <Check class="w-5 h-5 text-green-500 mt-0.5 flex-shrink-0" />
                  <span>Explainable results with mathematical proof</span>
                </li>
                <li class="flex items-start gap-2">
                  <Check class="w-5 h-5 text-green-500 mt-0.5 flex-shrink-0" />
                  <span>Computationally efficient and scalable</span>
                </li>
                <li class="flex items-start gap-2">
                  <Check class="w-5 h-5 text-green-500 mt-0.5 flex-shrink-0" />
                  <span>Industry standard method with proven track record</span>
                </li>
              </ul>
            </div>
          </div>
        </div>

        <!-- Performance Metrics -->
        <div class="bg-gradient-to-br from-green-50 to-emerald-50 dark:from-green-900/20 dark:to-emerald-900/20 rounded-2xl p-8">
          <div class="flex items-center gap-3 mb-6">
            <div class="w-12 h-12 bg-green-600 rounded-lg flex items-center justify-center">
              <Zap class="w-6 h-6 text-white" />
            </div>
            <h3 class="text-2xl font-bold text-gray-900 dark:text-white">Performance Metrics</h3>
          </div>

          <div class="space-y-6">
            <!-- Metric Cards -->
            <div class="grid grid-cols-2 gap-4">
              <div class="bg-white dark:bg-gray-800 rounded-lg p-4 text-center">
                <div class="text-2xl font-bold text-green-600 mb-1">1000+</div>
                <div class="text-sm text-gray-600 dark:text-gray-400">Events/Second</div>
              </div>
              <div class="bg-white dark:bg-gray-800 rounded-lg p-4 text-center">
                <div class="text-2xl font-bold text-green-600 mb-1">&lt;100ms</div>
                <div class="text-sm text-gray-600 dark:text-gray-400">Latency (P99)</div>
              </div>
              <div class="bg-white dark:bg-gray-800 rounded-lg p-4 text-center">
                <div class="text-2xl font-bold text-green-600 mb-1">95%+</div>
                <div class="text-sm text-gray-600 dark:text-gray-400">Accuracy Rate</div>
              </div>
              <div class="bg-white dark:bg-gray-800 rounded-lg p-4 text-center">
                <div class="text-2xl font-bold text-green-600 mb-1">10x</div>
                <div class="text-sm text-gray-600 dark:text-gray-400">Less Storage</div>
              </div>
            </div>

            <div>
              <h4 class="font-bold text-lg text-gray-900 dark:text-white mb-2">Integration Capabilities</h4>
              <div class="flex flex-wrap gap-2">
                <span class="px-3 py-1 bg-white dark:bg-gray-800 rounded-full text-sm text-gray-700 dark:text-gray-300">OpenTelemetry</span>
                <span class="px-3 py-1 bg-white dark:bg-gray-800 rounded-full text-sm text-gray-700 dark:text-gray-300">Kafka</span>
                <span class="px-3 py-1 bg-white dark:bg-gray-800 rounded-full text-sm text-gray-700 dark:text-gray-300">REST API</span>
                <span class="px-3 py-1 bg-white dark:bg-gray-800 rounded-full text-sm text-gray-700 dark:text-gray-300">Prometheus</span>
                <span class="px-3 py-1 bg-white dark:bg-gray-800 rounded-full text-sm text-gray-700 dark:text-gray-300">Grafana</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- What Banks Get -->
      <div class="bg-gradient-to-r from-primary-50 to-cyan-50 dark:from-primary-900/20 dark:to-cyan-900/20 rounded-2xl p-8 md:p-12">
        <h3 class="text-2xl font-bold text-center text-gray-900 dark:text-white mb-8">What Banks Get</h3>

        <div class="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
          <div class="text-center">
            <div class="w-16 h-16 bg-primary-100 dark:bg-primary-900/30 rounded-full flex items-center justify-center mx-auto mb-4">
              <RefreshCw class="w-8 h-8 text-primary-600" />
            </div>
            <h4 class="font-bold text-gray-900 dark:text-white mb-2">Drop-in Replacement</h4>
            <p class="text-sm text-gray-600 dark:text-gray-400">Replaces black-box fraud detection without infrastructure changes</p>
          </div>

          <div class="text-center">
            <div class="w-16 h-16 bg-primary-100 dark:bg-primary-900/30 rounded-full flex items-center justify-center mx-auto mb-4">
              <Calculator class="w-8 h-8 text-primary-600" />
            </div>
            <h4 class="font-bold text-gray-900 dark:text-white mb-2">Mathematical Proof</h4>
            <p class="text-sm text-gray-600 dark:text-gray-400">Statistical evidence for every flagged transaction</p>
          </div>

          <div class="text-center">
            <div class="w-16 h-16 bg-primary-100 dark:bg-primary-900/30 rounded-full flex items-center justify-center mx-auto mb-4">
              <FileCheck class="w-8 h-8 text-primary-600" />
            </div>
            <h4 class="font-bold text-gray-900 dark:text-white mb-2">Audit-Ready Reports</h4>
            <p class="text-sm text-gray-600 dark:text-gray-400">One-click generation of compliance documentation</p>
          </div>

          <div class="text-center">
            <div class="w-16 h-16 bg-primary-100 dark:bg-primary-900/30 rounded-full flex items-center justify-center mx-auto mb-4">
              <Database class="w-8 h-8 text-primary-600" />
            </div>
            <h4 class="font-bold text-gray-900 dark:text-white mb-2">Works with Existing Data</h4>
            <p class="text-sm text-gray-600 dark:text-gray-400">No data migration or system overhaul required</p>
          </div>
        </div>
      </div>

      <!-- Bottom CTA -->
      <div class="text-center mt-16">
        <a href="#proof" class="btn-primary">
          See It In Action
          <ArrowRight class="w-5 h-5 ml-2" />
        </a>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { Database, Search, FileText, CheckCircle, ArrowRight, Cpu, Check, Zap, RefreshCw, Calculator, FileCheck } from 'lucide-vue-next'
</script>

<style scoped>
.btn-primary {
  @apply inline-flex items-center bg-primary-600 hover:bg-primary-700 text-white px-8 py-3 rounded-lg font-semibold transition-colors;
}
</style>