<template>
  <div id="app" :class="{ 'dark': isDarkMode }">
    <!-- Navigation -->
    <nav class="w-full border-b border-black bg-background z-50 relative">
      <div class="container-padding mx-auto">
        <div class="flex items-center justify-between h-20">
          <!-- Logo -->
          <router-link to="/" class="flex items-center space-x-3 group">
            <div class="border border-black p-1 bg-white">
               <img src="/logo.svg" alt="Driftlock Logo" class="h-6 w-6 grayscale" style="object-fit: contain;">
            </div>
            <span class="text-xl font-bold font-sans tracking-widest uppercase">Driftlock</span>
          </router-link>

          <!-- Desktop Menu -->
          <div class="hidden md:flex items-center space-x-8">
            <router-link to="/docs" class="text-sm font-bold font-sans uppercase tracking-widest hover:underline underline-offset-4 decoration-1">Docs</router-link>
            <a href="#problem" class="text-sm font-bold font-sans uppercase tracking-widest hover:underline underline-offset-4 decoration-1">Problem</a>
            <a href="#solution" class="text-sm font-bold font-sans uppercase tracking-widest hover:underline underline-offset-4 decoration-1">Solution</a>
            <a href="#how" class="text-sm font-bold font-sans uppercase tracking-widest hover:underline underline-offset-4 decoration-1">How It Works</a>
            <router-link :to="{ path: '/', hash: '#playground' }" class="text-sm font-bold font-sans uppercase tracking-widest hover:underline underline-offset-4 decoration-1">Playground</router-link>
            
            <a href="#contact" class="brutalist-button">
              Partner
            </a>
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
          <a href="#problem" class="text-lg font-bold font-sans uppercase border-b border-gray-200 py-2" @click="isMobileMenuOpen = false">Problem</a>
          <a href="#solution" class="text-lg font-bold font-sans uppercase border-b border-gray-200 py-2" @click="isMobileMenuOpen = false">Solution</a>
          <a href="#how" class="text-lg font-bold font-sans uppercase border-b border-gray-200 py-2" @click="isMobileMenuOpen = false">How It Works</a>
          <router-link :to="{ path: '/', hash: '#playground' }" class="text-lg font-bold font-sans uppercase border-b border-gray-200 py-2" @click="isMobileMenuOpen = false">Playground</router-link>
          <a href="#contact" class="brutalist-button text-center w-full" @click="isMobileMenuOpen = false">Become a Partner</a>
        </div>
      </div>
    </nav>

    <!-- Main Content -->
    <main>
      <router-view />
    </main>

    <!-- Footer -->
    <footer class="border-t border-black bg-background py-12">
        <div class="container-padding mx-auto">
            <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-8">
                <div class="flex flex-col space-y-4">
                    <div class="flex items-center space-x-3">
                        <div class="border border-black p-1 bg-white">
                            <img src="/logo.svg" alt="Driftlock Logo" class="h-5 w-5 grayscale" style="object-fit: contain;">
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
                        <a href="#features" class="text-sm font-serif hover:underline underline-offset-2">Features</a>
                        <a href="#pricing" class="text-sm font-serif hover:underline underline-offset-2">Pricing</a>
                        <a href="/docs" class="text-sm font-serif hover:underline underline-offset-2">Documentation</a>
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
import { ref, onMounted, onUnmounted } from 'vue'
import { Menu } from 'lucide-vue-next'
import { useRoute } from 'vue-router'

const route = useRoute()
const isDarkMode = ref(false) // Force light mode for brutalist style initially
const isMobileMenuOpen = ref(false)

// Removed scroll handler as sticky nav is removed
</script>

<style scoped>
/* No scoped styles needed, using utility classes */
</style>