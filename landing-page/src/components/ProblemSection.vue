<template>
  <section id="problem" class="section-padding bg-gray-50 dark:bg-gray-800">
    <div class="container-padding mx-auto">
      <!-- Section Header -->
      <div class="text-center mb-16">
        <h2 class="text-4xl md:text-5xl font-bold mb-6">
          <span class="text-gray-900 dark:text-white">The €50M Problem</span>
          <span class="gradient-text block">Starting January 2025</span>
        </h2>
        <p class="text-xl text-gray-600 dark:text-gray-300 max-w-3xl mx-auto">
          EU regulators are auditing bank mathematical systems for DORA compliance. Black-box fraud detection means automatic failure.
        </p>
      </div>

      <!-- Main Problem Cards -->
      <div class="grid md:grid-cols-2 gap-12 mb-16">
        <!-- DORA Compliance Card -->
        <div class="bg-white dark:bg-gray-900 rounded-2xl p-8 shadow-xl border border-red-100 dark:border-red-900/30 relative overflow-hidden">
          <div class="absolute top-0 right-0 w-32 h-32 bg-red-500/10 rounded-full -mr-16 -mt-16"></div>

          <div class="relative z-10">
            <div class="flex items-center gap-3 mb-6">
              <AlertTriangle class="w-8 h-8 text-red-500" />
              <h3 class="text-2xl font-bold text-gray-900 dark:text-white">DORA Compliance</h3>
            </div>

            <ul class="space-y-4 text-gray-700 dark:text-gray-300">
              <li class="flex items-start gap-3">
                <X class="w-5 h-5 text-red-500 mt-0.5 flex-shrink-0" />
                <span><strong class="text-red-600 dark:text-red-400">€50M-€200M fines</strong> per violation for non-compliance</span>
              </li>
              <li class="flex items-start gap-3">
                <X class="w-5 h-5 text-red-500 mt-0.5 flex-shrink-0" />
                <span><strong class="text-red-600 dark:text-red-400">Explainable algorithms required</strong> - must show reasoning for every decision</span>
              </li>
              <li class="flex items-start gap-3">
                <X class="w-5 h-5 text-red-500 mt-0.5 flex-shrink-0" />
                <span><strong class="text-red-600 dark:text-red-400">Audit trails mandatory</strong> for all automated fraud detection</span>
              </li>
              <li class="flex items-start gap-3">
                <X class="w-5 h-5 text-red-500 mt-0.5 flex-shrink-0" />
                <span><strong class="text-red-600 dark:text-red-400">Enforcement active since January 2025</strong> - fines being issued now</span>
              </li>
            </ul>
          </div>
        </div>

        <!-- Black Box Problem Card -->
        <div class="bg-white dark:bg-gray-900 rounded-2xl p-8 shadow-xl border border-gray-200 dark:border-gray-700 relative overflow-hidden">
          <div class="absolute top-0 right-0 w-32 h-32 bg-gray-500/10 rounded-full -mr-16 -mt-16"></div>

          <div class="relative z-10">
            <div class="flex items-center gap-3 mb-6">
              <EyeOff class="w-8 h-8 text-gray-500" />
              <h3 class="text-2xl font-bold text-gray-900 dark:text-white">Black-Box Algorithm Problems</h3>
            </div>

            <ul class="space-y-4 text-gray-700 dark:text-gray-300">
              <li class="flex items-start gap-3">
                <X class="w-5 h-5 text-red-500 mt-0.5 flex-shrink-0" />
                <span><strong class="text-red-600 dark:text-red-400">No explanations</strong> - "Why was this flagged?" has no answer</span>
              </li>
              <li class="flex items-start gap-3">
                <X class="w-5 h-5 text-red-500 mt-0.5 flex-shrink-0" />
                <span><strong class="text-red-600 dark:text-red-400">Regulatory rejection</strong> - auditors automatically fail black boxes</span>
              </li>
              <li class="flex items-start gap-3">
                <X class="w-5 h-5 text-red-500 mt-0.5 flex-shrink-0" />
                <span><strong class="text-red-600 dark:text-red-400">High false positives</strong> - can't explain why errors occur</span>
              </li>
              <li class="flex items-start gap-3">
                <X class="w-5 h-5 text-red-500 mt-0.5 flex-shrink-0" />
                <span><strong class="text-red-600 dark:text-red-400">Manual investigation</strong> - teams spend weeks on each case</span>
              </li>
            </ul>
          </div>
        </div>
      </div>

      <!-- Quote Section -->
      <div class="bg-gradient-to-r from-red-50 to-orange-50 dark:from-red-900/20 dark:to-orange-900/20 rounded-2xl p-8 md:p-12 text-center border border-red-100 dark:border-red-800/30">
        <Quote class="w-12 h-12 text-red-500 mx-auto mb-6" />

        <blockquote class="text-2xl md:text-3xl font-bold text-gray-900 dark:text-white mb-6 max-w-4xl mx-auto">
          "When your ML system flags a transaction as suspicious, regulators demand: <strong class="text-red-600 dark:text-red-400">'Show your work.'</strong> Black-box models can't answer that."
        </blockquote>

        <div class="text-lg text-gray-600 dark:text-gray-300">
          <strong>Article 22, GDPR</strong> + <strong>Basel III Requirements</strong> + <strong>DORA Mandate</strong>
        </div>
      </div>

      <!-- Timeline -->
      <div class="mt-16">
        <h3 class="text-2xl font-bold text-center text-gray-900 dark:text-white mb-12">Compliance Timeline</h3>

        <div class="relative">
          <!-- Timeline Line -->
          <div class="absolute left-1/2 transform -translate-x-1/2 w-0.5 h-full bg-gray-300 dark:bg-gray-600"></div>

          <div class="space-y-12">
            <div class="flex items-center justify-between">
              <div class="w-5/12 text-right pr-8">
                <div class="text-lg font-bold text-gray-900 dark:text-white">DORA Approved</div>
                <div class="text-gray-600 dark:text-gray-400">EU passes Digital Operational Resilience Act</div>
              </div>
              <div class="w-2/12 flex justify-center">
                <div class="w-4 h-4 bg-primary-600 rounded-full border-4 border-white dark:border-gray-900"></div>
              </div>
              <div class="w-5/12 pl-8 text-gray-500">Jan 2023</div>
            </div>

            <div class="flex items-center justify-between">
              <div class="w-5/12 text-right pr-8 text-gray-500">Throughout 2024</div>
              <div class="w-2/12 flex justify-center">
                <div class="w-4 h-4 bg-orange-500 rounded-full border-4 border-white dark:border-gray-900"></div>
              </div>
              <div class="w-5/12 pl-8">
                <div class="text-lg font-bold text-gray-900 dark:text-white">Preparation Period</div>
                <div class="text-gray-600 dark:text-gray-400">Banks must implement compliant systems</div>
              </div>
            </div>

            <div class="flex items-center justify-between">
              <div class="w-5/12 text-right pr-8">
                <div class="text-lg font-bold text-red-600 dark:text-red-400">Enforcement Begins</div>
                <div class="text-gray-600 dark:text-gray-400">€50M-€200M fines per violation</div>
              </div>
              <div class="w-2/12 flex justify-center">
                <div class="w-4 h-4 bg-red-500 rounded-full border-4 border-white dark:border-gray-900 animate-pulse"></div>
              </div>
              <div class="w-5/12 pl-8 text-gray-500 text-red-600 dark:text-red-400 font-bold">ENFORCEMENT ACTIVE</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Call to Action -->
      <div class="mt-16 text-center">
        <div class="inline-flex items-center gap-2 px-6 py-3 bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400 rounded-full mb-6">
          <Clock class="w-5 h-5" />
          <span class="font-semibold">DORA Enforcement Active Since January 2025</span>
        </div>
        <p class="text-xl text-gray-600 dark:text-gray-300 mb-8">
          Banks are being audited now. Every day without explainable fraud detection increases your compliance risk.
        </p>
        <a href="#solution" class="btn-primary">
          See How We Solve This
          <ArrowRight class="w-5 h-5 ml-2" />
        </a>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { AlertTriangle, X, EyeOff, Quote, Clock, ArrowRight } from 'lucide-vue-next'
</script>

<style scoped>
.btn-primary {
  @apply inline-flex items-center bg-primary-600 hover:bg-primary-700 text-white px-8 py-3 rounded-lg font-semibold transition-colors;
}
</style>