<template>
  <div id="app" :class="{ 'dark': isDarkMode }">
    <!-- Skip to content link for accessibility -->
    <a 
      href="#main-content" 
      class="sr-only focus:not-sr-only focus:absolute focus:top-4 focus:left-4 focus:z-[100] focus:px-4 focus:py-2 focus:bg-primary-600 focus:text-white focus:rounded-lg focus:shadow-lg"
    >
      Skip to main content
    </a>

    <!-- Navigation -->
    <nav 
      class="fixed top-0 left-0 right-0 w-full z-50 transition-all duration-300"
      :class="{
        'bg-white/95 dark:bg-gray-900/95 backdrop-blur-xl shadow-lg': true,
        'bg-white/80 dark:bg-gray-900/80': isScrolled
      }"
      role="navigation"
      aria-label="Main navigation"
    >
      <div class="absolute inset-0 bg-gradient-to-r from-blue-50/50 via-transparent to-cyan-50/50 dark:from-blue-950/20 dark:via-transparent dark:to-cyan-950/20 pointer-events-none"></div>
      <div class="relative container-padding mx-auto">
        <div class="flex items-center justify-between h-20">
          <a href="#app" class="flex items-center space-x-3 group" aria-label="Driftlock home">
            <div class="bg-gradient-to-br from-blue-600 to-indigo-600 p-2.5 rounded-xl shadow-lg group-hover:shadow-xl group-hover:scale-105 transition-all duration-300">
              <Shield class="w-6 h-6 text-white" aria-hidden="true" />
            </div>
            <span class="font-bold text-2xl bg-gradient-to-r from-gray-900 via-blue-700 to-indigo-700 dark:from-white dark:via-blue-300 dark:to-indigo-300 bg-clip-text text-transparent group-hover:from-blue-600 group-hover:to-indigo-600 transition-all duration-300">Driftlock</span>
          </a>

          <div class="hidden md:flex items-center space-x-1" role="menubar">
            <a 
              href="#problem" 
              class="px-4 py-2 text-gray-700 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400 transition-all duration-200 font-medium rounded-lg hover:bg-gray-100/50 dark:hover:bg-gray-800/50 relative group focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
              role="menuitem"
            >
              Problem
              <span class="absolute bottom-0 left-0 right-0 h-0.5 bg-gradient-to-r from-blue-600 to-indigo-600 scale-x-0 group-hover:scale-x-100 transition-transform duration-200" aria-hidden="true"></span>
            </a>
            <a 
              href="#solution" 
              class="px-4 py-2 text-gray-700 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400 transition-all duration-200 font-medium rounded-lg hover:bg-gray-100/50 dark:hover:bg-gray-800/50 relative group focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
              role="menuitem"
            >
              Solution
              <span class="absolute bottom-0 left-0 right-0 h-0.5 bg-gradient-to-r from-blue-600 to-indigo-600 scale-x-0 group-hover:scale-x-100 transition-transform duration-200" aria-hidden="true"></span>
            </a>
            <a 
              href="#proof" 
              class="px-4 py-2 text-gray-700 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400 transition-all duration-200 font-medium rounded-lg hover:bg-gray-100/50 dark:hover:bg-gray-800/50 relative group focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
              role="menuitem"
            >
              Proof
              <span class="absolute bottom-0 left-0 right-0 h-0.5 bg-gradient-to-r from-blue-600 to-indigo-600 scale-x-0 group-hover:scale-x-100 transition-transform duration-200" aria-hidden="true"></span>
            </a>
            <a 
              href="#cta" 
              class="btn-primary px-6 py-2.5 rounded-xl ml-2 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
              role="menuitem"
            >
              Become a Partner
            </a>
          </div>

          <div class="flex items-center space-x-3">
            <button
              @click="toggleDarkMode"
              class="p-2.5 rounded-xl hover:bg-gray-100 dark:hover:bg-gray-800 transition-all duration-200 hover:scale-110 group focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
              aria-label="Toggle dark mode"
              :aria-pressed="isDarkMode"
            >
              <Sun v-if="isDarkMode" class="w-5 h-5 text-gray-700 dark:text-gray-300 group-hover:text-primary-600 dark:group-hover:text-primary-400 transition-colors" aria-hidden="true" />
              <Moon v-else class="w-5 h-5 text-gray-700 dark:text-gray-300 group-hover:text-primary-600 dark:group-hover:text-primary-400 transition-colors" aria-hidden="true" />
            </button>

            <button 
              class="md:hidden p-2.5 rounded-xl hover:bg-gray-100 dark:hover:bg-gray-800 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2" 
              @click="isMobileMenuOpen = !isMobileMenuOpen"
              :aria-label="isMobileMenuOpen ? 'Close menu' : 'Open menu'"
              :aria-expanded="isMobileMenuOpen"
            >
              <Menu class="w-6 h-6 text-gray-700 dark:text-gray-300" aria-hidden="true" />
            </button>
          </div>
        </div>
      </div>

      <!-- Mobile Menu -->
      <transition
        enter-active-class="transition-all duration-300 ease-out"
        enter-from-class="opacity-0 -translate-y-4"
        enter-to-class="opacity-100 translate-y-0"
        leave-active-class="transition-all duration-200 ease-in"
        leave-from-class="opacity-100 translate-y-0"
        leave-to-class="opacity-0 -translate-y-4"
      >
        <div 
          v-if="isMobileMenuOpen" 
          class="md:hidden bg-white/95 dark:bg-gray-900/95 backdrop-blur-xl border-t border-gray-200/50 dark:border-gray-700/50 shadow-lg"
          role="menu"
        >
          <div class="container-padding mx-auto py-4 space-y-2">
            <a 
              href="#problem" 
              class="block px-4 py-3 text-gray-700 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400 rounded-lg hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200 font-medium focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
              @click="isMobileMenuOpen = false"
              role="menuitem"
            >
              Problem
            </a>
            <a 
              href="#solution" 
              class="block px-4 py-3 text-gray-700 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400 rounded-lg hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200 font-medium focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
              @click="isMobileMenuOpen = false"
              role="menuitem"
            >
              Solution
            </a>
            <a 
              href="#proof" 
              class="block px-4 py-3 text-gray-700 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400 rounded-lg hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200 font-medium focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
              @click="isMobileMenuOpen = false"
              role="menuitem"
            >
              Proof
            </a>
            <a 
              href="#cta" 
              class="btn-primary w-full text-center mt-4 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
              @click="isMobileMenuOpen = false"
              role="menuitem"
            >
              Become a Partner
            </a>
          </div>
        </div>
      </transition>
    </nav>

    <!-- Main Content -->
    <main id="main-content" class="pt-20" role="main">
      <HeroSection />
      <ProblemSection />
      <RegulatoryMapSection />
      <SolutionSection />
      <ProofSection />
      <ComparisonSection />
      <CTASection />
    </main>

    <!-- Footer -->
    <footer class="bg-gradient-to-b from-gray-50 to-white dark:from-gray-800 dark:to-gray-900 border-t border-gray-200 dark:border-gray-700" role="contentinfo">
      <div class="container-padding mx-auto py-16">
        <div class="grid md:grid-cols-3 gap-12 mb-12">
          <div>
            <div class="flex items-center space-x-3 mb-6">
              <div class="bg-gradient-to-br from-blue-600 to-indigo-600 p-2 rounded-xl">
                <Shield class="w-5 h-5 text-white" aria-hidden="true" />
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
              <li><a href="#solution" class="hover:text-primary-600 dark:hover:text-primary-400 transition-colors focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 rounded">Technology</a></li>
              <li><a href="#proof" class="hover:text-primary-600 dark:hover:text-primary-400 transition-colors focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 rounded">Demo</a></li>
              <li><a href="#comparison" class="hover:text-primary-600 dark:hover:text-primary-400 transition-colors focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 rounded">Comparison</a></li>
            </ul>
          </div>

          <div>
            <h3 class="font-bold text-gray-900 dark:text-white mb-4 text-sm uppercase tracking-wider">Contact</h3>
            <ul class="space-y-3 text-sm text-gray-600 dark:text-gray-400">
              <li><a href="#cta" class="hover:text-primary-600 dark:hover:text-primary-400 transition-colors focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 rounded">Become a Partner</a></li>
              <li><a href="mailto:hunter@shannonlabs.dev" class="hover:text-primary-600 dark:hover:text-primary-400 transition-colors focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 rounded">Email Us</a></li>
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
import { ref, onMounted, onUnmounted } from 'vue'
import { Shield, Sun, Moon, Menu } from 'lucide-vue-next'
import HeroSection from './components/HeroSection.vue'
import ProblemSection from './components/ProblemSection.vue'
import RegulatoryMapSection from './components/RegulatoryMapSection.vue'
import SolutionSection from './components/SolutionSection.vue'
import ProofSection from './components/ProofSection.vue'
import ComparisonSection from './components/ComparisonSection.vue'
import CTASection from './components/CTASection.vue'

const isDarkMode = ref(false)
const isMobileMenuOpen = ref(false)
const isScrolled = ref(false)

const toggleDarkMode = () => {
  isDarkMode.value = !isDarkMode.value
  if (isDarkMode.value) {
    document.documentElement.classList.add('dark')
    localStorage.setItem('theme', 'dark')
  } else {
    document.documentElement.classList.remove('dark')
    localStorage.setItem('theme', 'light')
  }
}

const handleScroll = () => {
  isScrolled.value = window.scrollY > 20
}

onMounted(() => {
  // Check for saved theme preference or default to light mode
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme === 'dark') {
    isDarkMode.value = true
    document.documentElement.classList.add('dark')
  }
  
  // Handle scroll for navigation styling
  window.addEventListener('scroll', handleScroll)
  handleScroll()
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
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

.btn-primary:focus-visible {
  @apply ring-2 ring-primary-500 ring-offset-2;
}
</style>