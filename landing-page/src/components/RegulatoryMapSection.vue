<template>
  <section id="regulations" class="section-padding bg-white dark:bg-gray-900">
    <div class="container-padding mx-auto">
      <!-- Section Header -->
      <div class="text-center mb-16">
        <h2 class="text-4xl md:text-5xl font-bold mb-6">
          <span class="text-gray-900 dark:text-white">Global Regulatory</span>
          <span class="gradient-text block">Landscape</span>
        </h2>
        <p class="text-xl text-gray-600 dark:text-gray-300 max-w-3xl mx-auto">
          Explainability requirements are converging worldwide. Driftlock addresses compliance across EU, US, and other jurisdictions.
        </p>
      </div>

      <!-- Regulatory Comparison Table -->
      <div class="mb-16 overflow-x-auto">
        <div class="bg-gray-50 dark:bg-gray-800 rounded-2xl p-8">
          <div class="grid md:grid-cols-4 gap-6">
            <!-- Header Row -->
            <div class="font-bold text-gray-900 dark:text-white text-sm uppercase tracking-wider">Region</div>
            <div class="font-bold text-gray-900 dark:text-white text-sm uppercase tracking-wider">Regulation</div>
            <div class="font-bold text-gray-900 dark:text-white text-sm uppercase tracking-wider">Explainability Requirement</div>
            <div class="font-bold text-gray-900 dark:text-white text-sm uppercase tracking-wider">Penalties</div>

            <!-- EU Row -->
            <div class="text-gray-700 dark:text-gray-300 font-semibold">EU</div>
            <div class="text-gray-700 dark:text-gray-300">DORA Article 17</div>
            <div class="text-gray-600 dark:text-gray-400 text-sm">Mandatory explainability for all automated fraud detection systems</div>
            <div class="text-red-600 dark:text-red-400 font-semibold">Up to 2% annual turnover</div>

            <!-- US Federal Row -->
            <div class="text-gray-700 dark:text-gray-300 font-semibold">US Federal</div>
            <div class="text-gray-700 dark:text-gray-300">CFPB Fair Lending</div>
            <div class="text-gray-600 dark:text-gray-400 text-sm">Explainability required for AI/ML models in credit decisions</div>
            <div class="text-blue-600 dark:text-blue-400 font-semibold">Significant fines</div>

            <!-- NY Row -->
            <div class="text-gray-700 dark:text-gray-300 font-semibold">New York</div>
            <div class="text-gray-700 dark:text-gray-300">NYDFS Part 500</div>
            <div class="text-gray-600 dark:text-gray-400 text-sm">Cybersecurity requirements with audit trail mandates</div>
            <div class="text-blue-600 dark:text-blue-400 font-semibold">Enforcement actions</div>

            <!-- California Row -->
            <div class="text-gray-700 dark:text-gray-300 font-semibold">California</div>
            <div class="text-gray-700 dark:text-gray-300">CCPA + State Regs</div>
            <div class="text-gray-600 dark:text-gray-400 text-sm">Data privacy and algorithmic transparency requirements</div>
            <div class="text-blue-600 dark:text-blue-400 font-semibold">Per-violation fines</div>
          </div>
        </div>
      </div>

      <!-- How Driftlock Addresses Each -->
      <div class="bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-blue-900/20 dark:to-indigo-900/20 rounded-2xl p-8 md:p-12">
        <h3 class="text-2xl font-bold text-center text-gray-900 dark:text-white mb-8">How Driftlock Addresses Global Requirements</h3>
        
        <div class="grid md:grid-cols-2 gap-8">
          <div class="bg-white dark:bg-gray-800 rounded-xl p-6">
            <h4 class="font-bold text-lg text-gray-900 dark:text-white mb-4 flex items-center gap-2">
              <Shield class="w-5 h-5 text-primary-600" />
              EU DORA Compliance
            </h4>
            <ul class="space-y-2 text-gray-600 dark:text-gray-400 text-sm">
              <li class="flex items-start gap-2">
                <CheckCircle class="w-4 h-4 text-green-600 mt-0.5 flex-shrink-0" />
                <span>Mathematical explanations for every anomaly (Article 17)</span>
              </li>
              <li class="flex items-start gap-2">
                <CheckCircle class="w-4 h-4 text-green-600 mt-0.5 flex-shrink-0" />
                <span>Immutable audit trails with NCD scores and p-values</span>
              </li>
              <li class="flex items-start gap-2">
                <CheckCircle class="w-4 h-4 text-green-600 mt-0.5 flex-shrink-0" />
                <span>Regulator-ready HTML reports for compliance submissions</span>
              </li>
            </ul>
          </div>

          <div class="bg-white dark:bg-gray-800 rounded-xl p-6">
            <h4 class="font-bold text-lg text-gray-900 dark:text-white mb-4 flex items-center gap-2">
              <Shield class="w-5 h-5 text-primary-600" />
              US Regulatory Compliance
            </h4>
            <ul class="space-y-2 text-gray-600 dark:text-gray-400 text-sm">
              <li class="flex items-start gap-2">
                <CheckCircle class="w-4 h-4 text-green-600 mt-0.5 flex-shrink-0" />
                <span>CFPB-compliant explainability for fair lending requirements</span>
              </li>
              <li class="flex items-start gap-2">
                <CheckCircle class="w-4 h-4 text-green-600 mt-0.5 flex-shrink-0" />
                <span>NYDFS Part 500 audit trail compliance</span>
              </li>
              <li class="flex items-start gap-2">
                <CheckCircle class="w-4 h-4 text-green-600 mt-0.5 flex-shrink-0" />
                <span>GLBA and FFIEC operational resilience alignment</span>
              </li>
            </ul>
          </div>
        </div>
      </div>

      <!-- Key Insight -->
      <div class="mt-16 text-center">
        <div class="inline-flex items-center gap-2 px-6 py-3 bg-primary-100 dark:bg-primary-900/30 text-primary-700 dark:text-primary-400 rounded-full mb-6">
          <Globe class="w-5 h-5" />
          <span class="font-semibold">Universal Solution</span>
        </div>
        <p class="text-xl text-gray-600 dark:text-gray-300 max-w-3xl mx-auto">
          Driftlock's compression-based algorithm provides mathematical proof that satisfies regulators across jurisdictions. 
          One system, global compliance.
        </p>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { Shield, CheckCircle, Globe } from 'lucide-vue-next'
</script>

<style scoped>
</style>

