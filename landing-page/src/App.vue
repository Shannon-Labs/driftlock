<template>
  <div id="app" :class="{ 'dark': isDarkMode }">
    <!-- Navigation - Hide when in dashboard -->
    <nav v-if="!isInDashboard" class="w-full border-b border-black bg-background z-50 relative">
      <div class="container-padding mx-auto">
        <div class="flex items-center justify-between h-20">
          <!-- Logo -->
          <router-link to="/" class="flex items-center space-x-3 group">
            <!-- Replaced Logo with Brutalist Driftlock Visual -->
             <div class="w-10 h-10 group-hover:scale-110 transition-transform duration-300">
                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M6 11V4H18V11" stroke="black" stroke-width="4"/>
                  <rect x="2" y="11" width="20" height="11" fill="black"/>
                  <path d="M4 16.5H7L9.5 12.5L12.5 19.5L15 15.5L17 16.5H20" stroke="white" stroke-width="2" stroke-linejoin="bevel"/>
                </svg>
              </div>
            <span class="text-xl font-bold font-sans tracking-widest uppercase">Driftlock</span>
          </router-link>

          <!-- Desktop Menu -->
          <div class="hidden md:flex items-center space-x-8">
            <router-link to="/docs" class="text-sm font-bold font-sans uppercase tracking-widest hover:underline underline-offset-4 decoration-1">Docs</router-link>
            <a href="/#features" class="text-sm font-bold font-sans uppercase tracking-widest hover:underline underline-offset-4 decoration-1">Features</a>
            <a href="/#pricing" class="text-sm font-bold font-sans uppercase tracking-widest hover:underline underline-offset-4 decoration-1">Pricing</a>
            <router-link to="/playground" class="text-sm font-bold font-sans uppercase tracking-widest hover:underline underline-offset-4 decoration-1">Playground</router-link>

            <router-link v-if="!authStore.isAuthenticated" to="/login" class="brutalist-button">
              Sign In
            </router-link>
            <router-link v-else to="/dashboard" class="brutalist-button">
              Dashboard
            </router-link>
          </div>

          <!-- Mobile Menu Button -->
          <div class="flex items-center space-x-4 md:hidden">
            <button 
              class="p-2 border border-black hover:bg-black hover:text-white transition-colors" 
              @click="isMobileMenuOpen = !isMobileMenuOpen"
              aria-label="Toggle menu"
            >
              <Menu class="w-6 h-6" />
            </button>
          </div>
        </div>
      </div>

      <!-- Mobile Menu -->
      <div v-if="isMobileMenuOpen" class="md:hidden border-t border-black bg-background">
        <div class="container-padding mx-auto py-4 flex flex-col space-y-4">
          <router-link to="/docs" class="text-lg font-bold font-sans uppercase border-b border-gray-200 py-2" @click="isMobileMenuOpen = false">Docs</router-link>
          <a href="/#features" class="text-lg font-bold font-sans uppercase border-b border-gray-200 py-2" @click="isMobileMenuOpen = false">Features</a>
          <a href="/#pricing" class="text-lg font-bold font-sans uppercase border-b border-gray-200 py-2" @click="isMobileMenuOpen = false">Pricing</a>
          <router-link to="/playground" class="text-lg font-bold font-sans uppercase border-b border-gray-200 py-2" @click="isMobileMenuOpen = false">Playground</router-link>
          <router-link v-if="!authStore.isAuthenticated" to="/login" class="brutalist-button text-center w-full" @click="isMobileMenuOpen = false">Sign In</router-link>
          <router-link v-else to="/dashboard" class="brutalist-button text-center w-full" @click="isMobileMenuOpen = false">Dashboard</router-link>
        </div>
      </div>
    </nav>

    <!-- Main Content -->
    <main>
      <router-view />
    </main>

    <!-- Footer - Hide when in dashboard -->
    <footer v-if="!isInDashboard" class="border-t border-black bg-background py-12">
        <div class="container-padding mx-auto">
            <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-8">
                <div class="flex flex-col space-y-4">
                    <div class="flex items-center space-x-3">
                        <!-- Replaced Logo with Brutalist Driftlock Visual -->
                         <div class="w-8 h-8">
                            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                              <path d="M6 11V4H18V11" stroke="black" stroke-width="4"/>
                              <rect x="2" y="11" width="20" height="11" fill="black"/>
                              <path d="M4 16.5H7L9.5 12.5L12.5 19.5L15 15.5L17 16.5H20" stroke="white" stroke-width="2" stroke-linejoin="bevel"/>
                            </svg>
                          </div>
                        <span class="text-lg font-bold font-sans uppercase tracking-widest">Driftlock</span>
                    </div>
                    <p class="text-sm font-serif text-foreground max-w-xs leading-relaxed">
                        Anomaly detection for high-velocity data streams.
                        <br>
                        © 2025 Shannon Labs.
                    </p>
                </div>
                
                <div class="flex flex-col md:flex-row gap-8 md:gap-16">
                    <div class="flex flex-col space-y-2">
                        <h4 class="font-bold font-sans uppercase tracking-widest text-sm border-b border-black pb-1 mb-2">Product</h4>
                        <a href="#showcase" class="text-sm font-serif hover:underline underline-offset-2">Radar</a>
                        <a href="#pricing" class="text-sm font-serif hover:underline underline-offset-2">Pricing</a>
                        <router-link to="/docs" class="text-sm font-serif hover:underline underline-offset-2">Documentation</router-link>
                    </div>
                    <div class="flex flex-col space-y-2">
                        <h4 class="font-bold font-sans uppercase tracking-widest text-sm border-b border-black pb-1 mb-2">Company</h4>
                        <a href="https://shannonlabs.dev" class="text-sm font-serif hover:underline underline-offset-2">Shannon Labs</a>
                        <a href="https://github.com/Shannon-Labs/driftlock" class="text-sm font-serif hover:underline underline-offset-2">GitHub</a>
                        <a href="mailto:contact@shannonlabs.dev" class="text-sm font-serif hover:underline underline-offset-2">Contact</a>
                    </div>
                </div>
            </div>
        </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Menu } from 'lucide-vue-next'
import { useRoute } from 'vue-router'
import { useAuthStore } from './stores/auth'

const authStore = useAuthStore()
const route = useRoute()
const isDarkMode = ref(false) // Force light mode for brutalist style initially
const isMobileMenuOpen = ref(false)

// Check if current route is in dashboard
const isInDashboard = computed(() => route.path.startsWith('/dashboard'))
</script>

<style scoped>
/* No scoped styles needed, using utility classes */
</style>