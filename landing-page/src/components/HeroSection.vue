<template>
  <section class="relative min-h-screen flex items-center justify-center overflow-hidden bg-gradient-to-br from-blue-50 via-white to-cyan-50 dark:from-gray-900 dark:via-gray-900 dark:to-gray-800">
    <!-- Background Pattern -->
    <div class="absolute inset-0 opacity-5 dark:opacity-10">
      <div class="absolute inset-0" :style="backgroundPattern"></div>
    </div>

    <!-- Animated Elements -->
    <div class="absolute inset-0 overflow-hidden">
      <div class="absolute top-20 left-10 w-20 h-20 bg-primary-200 dark:bg-primary-800 rounded-full animate-float opacity-20"></div>
      <div class="absolute top-40 right-20 w-16 h-16 bg-cyan-200 dark:bg-cyan-800 rounded-full animate-float opacity-20" style="animation-delay: 1s;"></div>
      <div class="absolute bottom-20 left-20 w-24 h-24 bg-blue-200 dark:bg-blue-800 rounded-full animate-float opacity-20" style="animation-delay: 2s;"></div>
      <div class="absolute top-60 left-1/3 w-12 h-12 bg-primary-300 dark:bg-primary-700 rounded-full animate-float opacity-15" style="animation-delay: 0.5s;"></div>
    </div>

    <div class="relative z-10 container-padding mx-auto text-center">
      <!-- Main Headline -->
      <h1 class="text-5xl md:text-6xl lg:text-7xl font-bold mb-8 animate-fade-in">
        <span class="gradient-text">Stop €50M</span>
        <br />
        <span class="text-gray-900 dark:text-white">DORA Fines</span>
      </h1>

      <!-- Subheadline -->
      <p class="text-xl md:text-2xl lg:text-3xl text-gray-600 dark:text-gray-300 mb-12 max-w-4xl mx-auto animate-slide-up leading-relaxed">
        Explainable fraud detection that <strong>auditors love</strong>.
        Glass-box compliance with mathematical proof for every flagged transaction.
      </p>

      <!-- Key Benefits -->
      <div class="grid md:grid-cols-3 gap-8 mb-16 max-w-4xl mx-auto">
        <div class="bg-white/80 dark:bg-gray-800/80 backdrop-blur-sm rounded-xl p-6 border border-gray-200 dark:border-gray-700 animate-slide-up" style="animation-delay: 0.2s;">
          <Shield class="w-12 h-12 text-primary-600 mx-auto mb-4" />
          <h3 class="font-bold text-lg mb-2">Glass-Box AI</h3>
          <p class="text-gray-600 dark:text-gray-400 text-sm">Mathematical explanations for every decision, no black boxes</p>
        </div>

        <div class="bg-white/80 dark:bg-gray-800/80 backdrop-blur-sm rounded-xl p-6 border border-gray-200 dark:border-gray-700 animate-slide-up" style="animation-delay: 0.4s;">
          <Zap class="w-12 h-12 text-primary-600 mx-auto mb-4" />
          <h3 class="font-bold text-lg mb-2">Real-Time Detection</h3>
          <p class="text-gray-600 dark:text-gray-400 text-sm">Process 1000+ events/second with sub-second latency</p>
        </div>

        <div class="bg-white/80 dark:bg-gray-800/80 backdrop-blur-sm rounded-xl p-6 border border-gray-200 dark:border-gray-700 animate-slide-up" style="animation-delay: 0.6s;">
          <FileText class="w-12 h-12 text-primary-600 mx-auto mb-4" />
          <h3 class="font-bold text-lg mb-2">Audit-Ready Reports</h3>
          <p class="text-gray-600 dark:text-gray-400 text-sm">One-click compliance reports for regulator submissions</p>
        </div>
      </div>

      <!-- CTA Buttons -->
      <div class="flex flex-col sm:flex-row gap-6 justify-center items-center animate-slide-up" style="animation-delay: 0.8s;">
        <a
          href="#cta"
          class="group relative px-8 py-4 bg-primary-600 hover:bg-primary-700 text-white font-bold rounded-xl transition-all duration-300 transform hover:scale-105 shadow-lg hover:shadow-xl"
        >
          <span class="relative z-10 flex items-center gap-2">
            Request Compliance Audit
            <ArrowRight class="w-5 h-5 group-hover:translate-x-1 transition-transform" />
          </span>
        </a>

        <a
          href="#proof"
          class="group px-8 py-4 bg-white dark:bg-gray-800 text-primary-600 dark:text-primary-400 font-bold rounded-xl border-2 border-primary-600 dark:border-primary-400 hover:bg-primary-50 dark:hover:bg-primary-900/20 transition-all duration-300 flex items-center gap-2"
        >
          <Play class="w-5 h-5" />
          See 30-Second Demo
        </a>
      </div>

      <!-- Trust Indicators -->
      <div class="mt-20 pt-12 border-t border-gray-200 dark:border-gray-700">
        <div class="grid grid-cols-2 md:grid-cols-4 gap-8">
          <div class="text-center">
            <div class="text-3xl font-bold text-primary-600 mb-1">€50M</div>
            <div class="text-sm text-gray-600 dark:text-gray-400">Potential fines per bank</div>
          </div>
          <div class="text-center">
            <div class="text-3xl font-bold text-primary-600 mb-1">< 100ms</div>
            <div class="text-sm text-gray-600 dark:text-gray-400">Detection latency</div>
          </div>
          <div class="text-center">
            <div class="text-3xl font-bold text-primary-600 mb-1">95%+</div>
            <div class="text-sm text-gray-600 dark:text-gray-400">Accuracy rate</div>
          </div>
          <div class="text-center">
            <div class="text-3xl font-bold text-primary-600 mb-1">10x</div>
            <div class="text-sm text-gray-600 dark:text-gray-400">Less storage than ML</div>
          </div>
        </div>
      </div>

      <!-- Scroll Indicator -->
      <div class="absolute bottom-8 left-1/2 transform -translate-x-1/2 animate-bounce">
        <ChevronDown class="w-6 h-6 text-gray-400" />
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Shield, Zap, FileText, ArrowRight, Play, ChevronDown } from 'lucide-vue-next'

const backgroundPattern = computed(() => ({
  backgroundImage: `url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%231a73e8' fill-opacity='0.05'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")`,
}))
</script>

<style scoped>
@keyframes float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-20px); }
}
</style>