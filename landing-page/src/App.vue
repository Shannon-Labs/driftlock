<template>
  <div id="app" :class="{ 'dark': isDarkMode }">
    <!-- Navigation -->
    <nav class="fixed top-0 w-full bg-white/90 dark:bg-gray-900/90 backdrop-blur-xl border-b border-gray-200/50 dark:border-gray-700/50 z-50 premium-shadow">
      <div class="container-padding mx-auto">
        <div class="flex items-center justify-between h-20">
          <div class="flex items-center space-x-3">
            <div class="bg-gradient-to-br from-blue-600 to-indigo-600 p-2 rounded-xl">
              <Shield class="w-6 h-6 text-white" />
            </div>
            <span class="font-bold text-2xl bg-gradient-to-r from-gray-900 to-gray-700 dark:from-white dark:to-gray-300 bg-clip-text text-transparent">Driftlock</span>
          </div>

          <div class="hidden md:flex items-center space-x-8">
            <a href="#problem" class="text-gray-700 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400 transition-colors font-medium">
              Problem
            </a>
            <a href="#solution" class="text-gray-700 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400 transition-colors font-medium">
              Solution
            </a>
            <a href="#proof" class="text-gray-700 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400 transition-colors font-medium">
              Proof
            </a>
            <a href="#cta" class="btn-primary px-6 py-3 rounded-xl">
              Become a Partner
            </a>
          </div>

          <div class="flex items-center space-x-4">
            <button
              @click="toggleDarkMode"
              class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
              aria-label="Toggle dark mode"
            >
              <Sun v-if="isDarkMode" class="w-5 h-5" />
              <Moon v-else class="w-5 h-5" />
            </button>

            <button class="md:hidden p-2" @click="isMobileMenuOpen = !isMobileMenuOpen">
              <Menu class="w-6 h-6" />
            </button>
          </div>
        </div>
      </div>

      <!-- Mobile Menu -->
      <div v-if="isMobileMenuOpen" class="md:hidden bg-white dark:bg-gray-900 border-t border-gray-200 dark:border-gray-700">
        <div class="container-padding mx-auto py-4 space-y-3">
          <a href="#problem" class="block text-gray-600 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400">
            Problem
          </a>
          <a href="#solution" class="block text-gray-600 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400">
            Solution
          </a>
          <a href="#proof" class="block text-gray-600 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400">
            Proof
          </a>
          <a href="#cta" class="btn-primary w-full text-center">
            Become a Partner
          </a>
        </div>
      </div>
    </nav>

    <!-- Main Content -->
    <main class="pt-16">
      <HeroSection />
      <ProblemSection />
      <SolutionSection />
      <ProofSection />
      <ComparisonSection />
      <CTASection />
    </main>

    <!-- Footer -->
    <footer class="bg-gradient-to-b from-gray-50 to-white dark:from-gray-800 dark:to-gray-900 border-t border-gray-200 dark:border-gray-700">
      <div class="container-padding mx-auto py-16">
        <div class="grid md:grid-cols-4 gap-12 mb-12">
          <div>
            <div class="flex items-center space-x-3 mb-6">
              <div class="bg-gradient-to-br from-blue-600 to-indigo-600 p-2 rounded-xl">
                <Shield class="w-5 h-5 text-white" />
              </div>
              <span class="font-bold text-xl bg-gradient-to-r from-gray-900 to-gray-700 dark:from-white dark:to-gray-300 bg-clip-text text-transparent">Driftlock</span>
            </div>
            <p class="text-sm text-gray-600 dark:text-gray-400 leading-relaxed mb-4">
              Explainable anomaly detection for regulated industries. Stop €50M DORA fines with glass-box algorithm.
            </p>
            <p class="text-xs text-gray-500 dark:text-gray-500">
              Built by Shannon Labs
            </p>
          </div>

          <div>
            <h3 class="font-bold text-gray-900 dark:text-white mb-4 text-sm uppercase tracking-wider">Product</h3>
            <ul class="space-y-3 text-sm text-gray-600 dark:text-gray-400">
              <li><a href="#solution" class="hover:text-primary-600 dark:hover:text-primary-400 transition-colors">Technology</a></li>
              <li><a href="#proof" class="hover:text-primary-600 dark:hover:text-primary-400 transition-colors">Demo</a></li>
              <li><a href="#comparison" class="hover:text-primary-600 dark:hover:text-primary-400 transition-colors">Comparison</a></li>
            </ul>
          </div>

          <div>
            <h3 class="font-bold text-gray-900 dark:text-white mb-4 text-sm uppercase tracking-wider">Company</h3>
            <ul class="space-y-3 text-sm text-gray-600 dark:text-gray-400">
              <li><a href="#" class="hover:text-primary-600 dark:hover:text-primary-400 transition-colors">About</a></li>
              <li><a href="#" class="hover:text-primary-600 dark:hover:text-primary-400 transition-colors">Blog</a></li>
              <li><a href="#" class="hover:text-primary-600 dark:hover:text-primary-400 transition-colors">Careers</a></li>
            </ul>
          </div>

          <div>
            <h3 class="font-bold text-gray-900 dark:text-white mb-4 text-sm uppercase tracking-wider">Legal</h3>
            <ul class="space-y-3 text-sm text-gray-600 dark:text-gray-400">
              <li><a href="#" class="hover:text-primary-600 dark:hover:text-primary-400 transition-colors">Privacy</a></li>
              <li><a href="#" class="hover:text-primary-600 dark:hover:text-primary-400 transition-colors">Terms</a></li>
              <li><a href="#" class="hover:text-primary-600 dark:hover:text-primary-400 transition-colors">Compliance</a></li>
            </ul>
          </div>
        </div>

        <div class="border-t border-gray-200 dark:border-gray-700 pt-8 flex flex-col md:flex-row justify-between items-center gap-4">
          <p class="text-sm text-gray-600 dark:text-gray-400">
            &copy; 2025 Shannon Labs. Licensed under Apache 2.0.
          </p>
          <div class="flex items-center gap-6 text-sm text-gray-600 dark:text-gray-400">
            <span class="flex items-center gap-2">
              <span class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></span>
              All systems operational
            </span>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Shield, Sun, Moon, Menu } from 'lucide-vue-next'
import HeroSection from './components/HeroSection.vue'
import ProblemSection from './components/ProblemSection.vue'
import SolutionSection from './components/SolutionSection.vue'
import ProofSection from './components/ProofSection.vue'
import ComparisonSection from './components/ComparisonSection.vue'
import CTASection from './components/CTASection.vue'

const isDarkMode = ref(false)
const isMobileMenuOpen = ref(false)

const toggleDarkMode = () => {
  isDarkMode.value = !isDarkMode.value
  if (isDarkMode.value) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}

onMounted(() => {
  // Check for saved theme preference or default to light mode
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme === 'dark') {
    isDarkMode.value = true
    document.documentElement.classList.add('dark')
  }
})
</script>

<style scoped>
.btn-primary {
  @apply bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white font-semibold transition-all duration-300 transform hover:scale-105;
  box-shadow:
    0 10px 15px -3px rgba(0, 0, 0, 0.1),
    0 4px 6px -2px rgba(0, 0, 0, 0.05),
    0 0 0 1px rgba(0, 0, 0, 0.05);
}
</style>