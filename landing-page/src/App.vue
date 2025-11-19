<template>
  <div id="app" :class="{ 'dark': isDarkMode }">
    <!-- Navigation -->
    <nav 
      class="fixed top-0 left-0 right-0 w-full z-50 transition-all duration-300"
      :class="{
        'bg-white/95 dark:bg-gray-900/95 backdrop-blur-xl shadow-lg': true,
        'bg-white/80 dark:bg-gray-900/80': isScrolled
      }"
    >
      <div class="absolute inset-0 bg-gradient-to-r from-blue-50/50 via-transparent to-cyan-50/50 dark:from-blue-950/20 dark:via-transparent dark:to-cyan-950/20 pointer-events-none"></div>
      <div class="relative container-padding mx-auto">
        <div class="flex items-center justify-between h-20">
          <router-link to="/" class="flex items-center space-x-2 group hover:opacity-90 transition-opacity">
            <!-- Logo (recreated from screenshot) -->
            <div class="flex-shrink-0 flex items-center justify-center">
                <img src="/logo.svg" alt="Driftlock Logo" class="h-7 w-7 transition-transform duration-200 group-hover:scale-105" style="object-fit: contain;">
            </div>
            <span class="text-xl font-extrabold font-mono text-gray-900 dark:text-white tracking-tight leading-none">Driftlock</span>
          </router-link>

          <div class="hidden md:flex items-center space-x-2">
            <router-link 
              to="/docs" 
              class="px-4 py-2 font-sans text-gray-700 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400 transition-all duration-200 font-medium rounded-lg hover:bg-gray-100/50 dark:hover:bg-gray-800/50 relative group"
            >
              Docs
              <span class="absolute bottom-0 left-0 right-0 h-0.5 bg-gradient-to-r from-blue-600 to-indigo-600 scale-x-0 group-hover:scale-x-100 transition-transform duration-200"></span>
            </router-link>
            <a 
              href="#problem" 
              class="px-4 py-2 font-sans text-gray-700 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400 transition-all duration-200 font-medium rounded-lg hover:bg-gray-100/50 dark:hover:bg-gray-800/50 relative group"
            >
              The Problem
              <span class="absolute bottom-0 left-0 right-0 h-0.5 bg-gradient-to-r from-blue-600 to-indigo-600 scale-x-0 group-hover:scale-x-100 transition-transform duration-200"></span>
            </a>
            <a 
              href="#solution" 
              class="px-4 py-2 font-sans text-gray-700 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400 transition-all duration-200 font-medium rounded-lg hover:bg-gray-100/50 dark:hover:bg-gray-800/50 relative group"
            >
              Solution
              <span class="absolute bottom-0 left-0 right-0 h-0.5 bg-gradient-to-r from-blue-600 to-indigo-600 scale-x-0 group-hover:scale-x-100 transition-transform duration-200"></span>
            </a>
            <a 
              href="#how" 
              class="px-4 py-2 font-sans text-gray-700 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400 transition-all duration-200 font-medium rounded-lg hover:bg-gray-100/50 dark:hover:bg-gray-800/50 relative group"
            >
              How It Works
              <span class="absolute bottom-0 left-0 right-0 h-0.5 bg-gradient-to-r from-blue-600 to-indigo-600 scale-x-0 group-hover:scale-x-100 transition-transform duration-200"></span>
            </a>
            <router-link 
              :to="{ path: '/', hash: '#playground' }"
              class="px-4 py-2 font-sans text-gray-700 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400 transition-all duration-200 font-medium rounded-lg hover:bg-gray-100/50 dark:hover:bg-gray-800/50 relative group"
            >
              Playground
              <span class="absolute bottom-0 left-0 right-0 h-0.5 bg-gradient-to-r from-blue-600 to-indigo-600 scale-x-0 group-hover:scale-x-100 transition-transform duration-200"></span>
            </router-link>
            <a 
              href="#contact" 
              class="btn-primary px-6 py-2.5 font-sans rounded-xl ml-2"
            >
              Become a Partner
            </a>
          </div>

          <div class="flex items-center space-x-3">
            <button
              @click="toggleDarkMode"
              class="p-2.5 rounded-xl hover:bg-gray-100 dark:hover:bg-gray-800 transition-all duration-200 hover:scale-110 group"
              aria-label="Toggle dark mode"
            >
              <Sun v-if="isDarkMode" class="w-5 h-5 text-gray-700 dark:text-gray-300 group-hover:text-primary-600 dark:group-hover:text-primary-400 transition-colors" />
              <Moon v-else class="w-5 h-5 text-gray-700 dark:text-gray-300 group-hover:text-primary-600 dark:group-hover:text-primary-400 transition-colors" />
            </button>

            <button 
              class="md:hidden p-2.5 rounded-xl hover:bg-gray-100 dark:hover:bg-gray-800 transition-all duration-200" 
              @click="isMobileMenuOpen = !isMobileMenuOpen"
              aria-label="Toggle menu"
            >
              <Menu class="w-6 h-6 text-gray-700 dark:text-gray-300" />
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
        >
          <div class="container-padding mx-auto py-4 space-y-2">
            <router-link 
              to="/docs" 
              class="block px-4 py-3 font-sans text-gray-700 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400 rounded-lg hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200 font-medium"
              @click="isMobileMenuOpen = false"
            >
              Docs
            </router-link>
            <a 
              href="#problem" 
              class="block px-4 py-3 font-sans text-gray-700 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400 rounded-lg hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200 font-medium"
              @click="isMobileMenuOpen = false"
            >
              The Problem
            </a>
            <a 
              href="#solution" 
              class="block px-4 py-3 font-sans text-gray-700 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400 rounded-lg hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200 font-medium"
              @click="isMobileMenuOpen = false"
            >
              Solution
            </a>
            <a 
              href="#how" 
              class="block px-4 py-3 font-sans text-gray-700 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400 rounded-lg hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200 font-medium"
              @click="isMobileMenuOpen = false"
            >
              How It Works
            </a>
          <a 
            href="#contact" 
            class="btn-primary w-full font-sans text-center mt-4"
            @click="isMobileMenuOpen = false"
          >
            Become a Partner
          </a>
          <router-link 
            :to="{ path: '/', hash: '#playground' }" 
            class="block px-4 py-3 font-sans text-gray-700 hover:text-primary-600 dark:text-gray-300 dark:hover:text-primary-400 rounded-lg hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200 font-medium"
            @click="isMobileMenuOpen = false"
          >
            Playground
          </router-link>
          </div>
        </div>
      </transition>
    </nav>

    <!-- Main Content -->
    <main class="pt-20">
      <router-view />
    </main>

    <!-- Footer -->
    <footer class="bg-white border-t border-gray-200 dark:bg-gray-900 dark:border-gray-700">
        <div class="container-padding mx-auto py-12">
            <div class="flex flex-col items-center justify-between gap-4 md:flex-row md:items-center">
                <div class="flex items-center space-x-2">
                    <div class="flex-shrink-0 flex items-center justify-center">
                        <img src="/logo.svg" alt="Driftlock Logo" class="h-7 w-7" style="object-fit: contain;">
                    </div>
                    <span class="text-lg font-bold text-gray-900 dark:text-white">Driftlock</span>
                </div>
                <p class="text-sm text-gray-600 dark:text-gray-400 text-center md:text-left">
                    © 2025 Shannon Labs. Licensed under Apache 2.0.
                </p>
                <div class="flex items-center space-x-6">
                    <a href="https://github.com/Shannon-Labs/driftlock" class="text-gray-500 hover:text-gray-700 dark:hover:text-gray-300 transition-colors">
                        <span class="sr-only">GitHub</span>
                        <svg class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                            <path fill-rule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.009-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.026 2.747-1.026.546 1.379.202 2.398.1 2.65.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" clip-rule="evenodd" />
                        </svg>
                    </a>
                </div>
            </div>
        </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { Sun, Moon, Menu } from 'lucide-vue-next'
import { useRoute } from 'vue-router'

const route = useRoute()

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
</style>