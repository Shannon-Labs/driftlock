<template>
  <div id="app" :class="{ 'dark': isDarkMode }">
    <!-- Navigation -->
    <nav class="fixed top-0 w-full bg-white/80 dark:bg-gray-900/80 backdrop-blur-md border-b border-gray-200 dark:border-gray-700 z-50">
      <div class="container-padding mx-auto">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center space-x-2">
            <Shield class="w-8 h-8 text-primary-600" />
            <span class="font-bold text-xl">DriftLock</span>
          </div>

          <div class="hidden md:flex items-center space-x-8">
            <a href="#problem" class="text-gray-600 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400 transition-colors">
              Problem
            </a>
            <a href="#solution" class="text-gray-600 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400 transition-colors">
              Solution
            </a>
            <a href="#proof" class="text-gray-600 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400 transition-colors">
              Proof
            </a>
            <a href="#cta" class="btn-primary">
              Request Demo
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
            Request Demo
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
    <footer class="bg-gray-50 dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700">
      <div class="container-padding mx-auto py-12">
        <div class="grid md:grid-cols-4 gap-8">
          <div>
            <div class="flex items-center space-x-2 mb-4">
              <Shield class="w-6 h-6 text-primary-600" />
              <span class="font-bold text-lg">DriftLock</span>
            </div>
            <p class="text-sm text-gray-600 dark:text-gray-400">
              Explainable anomaly detection for regulated industries.
            </p>
          </div>

          <div>
            <h3 class="font-semibold mb-4">Product</h3>
            <ul class="space-y-2 text-sm text-gray-600 dark:text-gray-400">
              <li><a href="#solution" class="hover:text-primary-600 dark:hover:text-primary-400">Technology</a></li>
              <li><a href="#proof" class="hover:text-primary-600 dark:hover:text-primary-400">Demo</a></li>
              <li><a href="#comparison" class="hover:text-primary-600 dark:hover:text-primary-400">Comparison</a></li>
            </ul>
          </div>

          <div>
            <h3 class="font-semibold mb-4">Company</h3>
            <ul class="space-y-2 text-sm text-gray-600 dark:text-gray-400">
              <li><a href="#" class="hover:text-primary-600 dark:hover:text-primary-400">About</a></li>
              <li><a href="#" class="hover:text-primary-600 dark:hover:text-primary-400">Blog</a></li>
              <li><a href="#" class="hover:text-primary-600 dark:hover:text-primary-400">Careers</a></li>
            </ul>
          </div>

          <div>
            <h3 class="font-semibold mb-4">Legal</h3>
            <ul class="space-y-2 text-sm text-gray-600 dark:text-gray-400">
              <li><a href="#" class="hover:text-primary-600 dark:hover:text-primary-400">Privacy</a></li>
              <li><a href="#" class="hover:text-primary-600 dark:hover:text-primary-400">Terms</a></li>
              <li><a href="#" class="hover:text-primary-600 dark:hover:text-primary-400">Compliance</a></li>
            </ul>
          </div>
        </div>

        <div class="border-t border-gray-200 dark:border-gray-600 mt-8 pt-8 text-center text-sm text-gray-600 dark:text-gray-400">
          <p>&copy; 2024 Shannon Labs. Licensed under Apache 2.0.</p>
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
  @apply bg-primary-600 hover:bg-primary-700 text-white px-6 py-2 rounded-lg font-medium transition-colors;
}
</style>