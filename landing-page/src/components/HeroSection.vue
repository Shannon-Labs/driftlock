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

    <div class="relative z-10 container-padding mx-auto">
      <div class="grid gap-12 items-start lg:grid-cols-2">
        <div class="text-center lg:text-left max-w-3xl mx-auto lg:mx-0">
          <!-- Main Headline -->
          <h1 class="text-5xl md:text-6xl lg:text-7xl xl:text-8xl font-bold mb-8 animate-fade-in leading-tight tracking-tight">
            <span class="gradient-text">Stop Multi-Million</span>
            <br />
            <span class="text-gray-900 dark:text-white">Dollar Fines</span>
          </h1>

          <!-- Alpha Badge -->
          <div class="mb-10 animate-fade-in" style="animation-delay: 0.2s;">
            <span class="inline-flex items-center px-5 py-2.5 bg-gradient-to-r from-orange-500 to-red-500 text-white text-sm font-semibold rounded-full shadow-lg hover:shadow-xl transition-all duration-300 hover:scale-105">
              <Sparkles class="w-4 h-4 mr-2" />
              Alpha Phase - Seeking Design Partners
            </span>
          </div>

          <!-- Subheadline -->
          <p class="text-xl md:text-2xl lg:text-3xl text-gray-600 dark:text-gray-300 mb-12 animate-slide-up leading-relaxed font-light" style="animation-delay: 0.3s;">
            Mathematically explainable fraud detection that <strong class="font-semibold text-gray-900 dark:text-white">auditors love</strong>.
            <br class="hidden sm:block" />
            Glass-box compliance for EU DORA & US regulations with mathematical proof for every flagged transaction.
          </p>

          <!-- Key Benefits -->
          <div class="grid md:grid-cols-2 gap-6 mb-12">
            <div class="bg-white/90 dark:bg-gray-800/90 backdrop-blur-md rounded-2xl p-6 border border-gray-200 dark:border-gray-700 animate-slide-up card-hover premium-shadow" style="animation-delay: 0.2s;">
              <div class="bg-gradient-to-br from-blue-500 to-indigo-600 w-14 h-14 rounded-xl flex items-center justify-center mb-4">
                <Shield class="w-7 h-7 text-white" />
              </div>
              <h3 class="font-bold text-xl mb-2 text-gray-900 dark:text-white">Glass-Box Algorithm</h3>
              <p class="text-gray-600 dark:text-gray-400 leading-relaxed">Mathematical explanations for every decision, no black boxes</p>
            </div>

            <div class="bg-white/90 dark:bg-gray-800/90 backdrop-blur-md rounded-2xl p-6 border border-gray-200 dark:border-gray-700 animate-slide-up card-hover premium-shadow" style="animation-delay: 0.4s;">
              <div class="bg-gradient-to-br from-cyan-500 to-blue-600 w-14 h-14 rounded-xl flex items-center justify-center mb-4">
                <Zap class="w-7 h-7 text-white" />
              </div>
              <h3 class="font-bold text-xl mb-2 text-gray-900 dark:text-white">Real-Time Detection</h3>
              <p class="text-gray-600 dark:text-gray-400 leading-relaxed">Process 1000+ events/second with sub-second latency</p>
            </div>

            <div class="bg-white/90 dark:bg-gray-800/90 backdrop-blur-md rounded-2xl p-6 border border-gray-200 dark:border-gray-700 animate-slide-up card-hover premium-shadow" style="animation-delay: 0.6s.">
              <div class="bg-gradient-to-br from-indigo-500 to-purple-600 w-14 h-14 rounded-xl flex items-center justify-center mb-4">
                <FileText class="w-7 h-7 text-white" />
              </div>
              <h3 class="font-bold text-xl mb-2 text-gray-900 dark:text-white">Audit-Ready Reports</h3>
              <p class="text-gray-600 dark:text-gray-400 leading-relaxed">One-click compliance reports for regulator submissions</p>
            </div>

            <div class="bg-white/90 dark:bg-gray-800/90 backdrop-blur-md rounded-2xl p-6 border border-gray-200 dark:border-gray-700 animate-slide-up card-hover premium-shadow" style="animation-delay: 0.8s;">
              <div class="bg-gradient-to-br from-emerald-500 to-green-600 w-14 h-14 rounded-xl flex items-center justify-center mb-4">
                <Clock class="w-7 h-7 text-white" />
              </div>
              <h3 class="font-bold text-xl mb-2 text-gray-900 dark:text-white">Compliance Proven</h3>
              <p class="text-gray-600 dark:text-gray-400 leading-relaxed">Ready for EU DORA and US regulators with traceable math</p>
            </div>
          </div>

          <!-- CTA Buttons -->
          <div class="flex flex-col sm:flex-row gap-4 sm:justify-start items-stretch sm:items-center animate-slide-up" style="animation-delay: 0.8s;">
            <a
              href="#cta"
              class="group relative px-8 py-4 bg-gradient-to-r from-red-600 to-orange-600 hover:from-red-700 hover:to-orange-700 text-white font-bold rounded-2xl transition-all duration-300 transform hover:scale-105 premium-shadow-lg text-lg text-center"
            >
              <span class="relative z-10 flex items-center justify-center gap-3">
                Get Compliant Now
                <ArrowRight class="w-6 h-6 group-hover:translate-x-1 transition-transform" />
              </span>
            </a>
            
            <a
              href="#playground"
              class="group px-8 py-4 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white font-semibold rounded-2xl transition-all duration-300 transform hover:scale-105 flex items-center justify-center gap-3 shadow-lg"
            >
              <Play class="w-5 h-5" />
              Try Playground
            </a>
            
            <a
              href="#proof"
              class="group px-8 py-4 bg-white/80 dark:bg-gray-800/80 text-primary-600 dark:text-primary-400 font-semibold rounded-2xl border-2 border-gray-300 dark:border-gray-600 hover:border-primary-600 dark:hover:border-primary-400 transition-all duration-300 flex items-center justify-center gap-3"
            >
              <FileText class="w-5 h-5" />
              Learn More
            </a>
          </div>

          <!-- Urgency Banner -->
          <div class="mt-6 inline-flex items-center gap-2 px-4 py-2 bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400 rounded-full">
            <Clock class="w-4 h-4" />
            <span class="text-sm font-semibold">DORA Enforcement Active Since January 2025 | US Regulations Ongoing</span>
          </div>

          <!-- Trust Indicators -->
          <div class="mt-16 pt-12 border-t border-gray-200 dark:border-gray-700">
            <div class="grid grid-cols-2 gap-8">
              <div class="text-left group">
                <div class="text-4xl md:text-5xl font-bold bg-gradient-to-r from-red-600 to-red-500 bg-clip-text text-transparent mb-2 transition-transform duration-300 group-hover:scale-105">Up to 2%</div>
                <div class="text-sm md:text-base text-gray-600 dark:text-gray-400 font-medium">of annual turnover (DORA)</div>
              </div>
              <div class="text-left group">
                <div class="text-4xl md:text-5xl font-bold bg-gradient-to-r from-blue-600 to-cyan-600 bg-clip-text text-transparent mb-2 transition-transform duration-300 group-hover:scale-105">&lt; 100ms</div>
                <div class="text-sm md:text-base text-gray-600 dark:text-gray-400 font-medium">Detection latency</div>
              </div>
              <div class="text-left group">
                <div class="text-4xl md:text-5xl font-bold bg-gradient-to-r from-green-600 to-emerald-600 bg-clip-text text-transparent mb-2 transition-transform duration-300 group-hover:scale-105">95%+</div>
                <div class="text-sm md:text-base text-gray-600 dark:text-gray-400 font-medium">Accuracy rate</div>
              </div>
              <div class="text-left group">
                <div class="text-4xl md:text-5xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent mb-2 transition-transform duration-300 group-hover:scale-105">10x</div>
                <div class="text-sm md:text-base text-gray-600 dark:text-gray-400 font-medium">Less storage than model-driven systems</div>
              </div>
            </div>
          </div>
        </div>

        <div id="playground" class="w-full max-w-3xl mx-auto lg:mx-0 lg:sticky lg:top-24">
          <PlaygroundShell variant="embedded" />
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
import { Shield, Zap, FileText, ArrowRight, Play, ChevronDown, Sparkles, Clock } from 'lucide-vue-next'
import PlaygroundShell from './playground/PlaygroundShell.vue'

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