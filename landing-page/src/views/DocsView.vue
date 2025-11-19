<template>
  <div class="min-h-screen bg-white">
    <!-- Navbar -->
    <header class="sticky top-0 z-50 w-full border-b border-gray-200 bg-white/95 backdrop-blur transition-all">
      <div class="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
        <div class="flex items-center gap-8">
          <router-link to="/" class="flex items-center gap-2">
            <img class="h-8 w-auto" src="/logo-icon.svg" alt="Driftlock" />
            <span class="font-mono font-bold text-gray-900">Driftlock Docs</span>
          </router-link>
        </div>
        <div class="flex items-center gap-4">
          <router-link to="/dashboard" class="text-sm font-medium text-gray-700 hover:text-blue-600">Dashboard</router-link>
          <a href="https://github.com/Shannon-Labs/driftlock" target="_blank" class="text-gray-500 hover:text-gray-900">
            <span class="sr-only">GitHub</span>
            <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
              <path fill-rule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" clip-rule="evenodd" />
            </svg>
          </a>
        </div>
      </div>
    </header>

    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="flex lg:gap-10">
        <!-- Sidebar Navigation -->
        <aside class="fixed inset-y-0 left-0 z-40 hidden w-64 -translate-x-full flex-col overflow-y-auto border-r border-gray-200 bg-white px-6 pb-10 pt-24 lg:static lg:flex lg:w-72 lg:translate-x-0 lg:pt-10">
          <nav class="space-y-8">
            <div v-for="(section, title) in navigation" :key="title">
              <h3 class="font-mono text-sm font-semibold text-gray-900">{{ title }}</h3>
              <ul class="mt-2 space-y-2">
                <li v-for="link in section" :key="link.href">
                  <router-link 
                    :to="link.href" 
                    class="block text-sm text-gray-600 hover:text-blue-600"
                    active-class="font-medium text-blue-600"
                  >
                    {{ link.label }}
                  </router-link>
                </li>
              </ul>
            </div>
          </nav>
        </aside>

        <!-- Main Content -->
        <main class="flex-1 py-10 lg:pl-8">
          <div v-if="loading" class="flex items-center justify-center py-20">
            <div class="h-8 w-8 animate-spin rounded-full border-2 border-blue-500 border-t-transparent"></div>
          </div>
          <div v-else-if="error" class="rounded-lg bg-red-50 p-4 text-red-700">
            {{ error }}
          </div>
          <article v-else class="prose prose-blue max-w-none prose-headings:font-mono prose-headings:font-bold prose-a:text-blue-600" v-html="content"></article>
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'

const route = useRoute()
const content = ref('')
const loading = ref(true)
const error = ref<string | null>(null)

// Configure marked with highlight.js
marked.setOptions({
  highlight: function(code, lang) {
    const language = hljs.getLanguage(lang) ? lang : 'plaintext';
    return hljs.highlight(code, { language }).value;
  },
  langPrefix: 'hljs language-',
})

const navigation = {
  'Getting Started': [
    { label: 'Introduction', href: '/docs' },
    { label: 'Quickstart', href: '/docs/quickstart' },
    { label: 'Architecture', href: '/docs/architecture' },
  ],
  'API Reference': [
    { label: 'HTTP API', href: '/docs/api' },
    { label: 'Authentication', href: '/docs/auth' },
  ],
  'Concepts': [
    { label: 'Algorithms (CBAD)', href: '/docs/algorithms' },
    { label: 'Compliance', href: '/docs/compliance' },
    { label: 'Use Cases', href: '/docs/use-cases' },
  ],
  'Deployment': [
    { label: 'Docker', href: '/docs/docker' },
    { label: 'Production Guide', href: '/docs/production' },
  ]
}

// Map routes to repo markdown files (raw content)
// In a real app, we'd fetch these from a CDN or bundle them.
// For this demo, we'll mock fetching from GitHub raw or local static files if moved.
// Since we can't dynamically import md files easily without a plugin in existing setup,
// we will fetch from GitHub raw for now, or use a few hardcoded placeholders.

const docMap: Record<string, string> = {
  'undefined': 'README.md', // /docs -> README
  'quickstart': 'docs/DEPLOYMENT_QUICKSTART.md',
  'architecture': 'docs/ARCHITECTURE.md',
  'api': 'docs/API.md',
  'auth': 'docs/api/API-AUTH.md', // Assuming this exists or we fallback
  'algorithms': 'docs/ALGORITHMS.md',
  'compliance': 'docs/COMPLIANCE_DORA.md',
  'docker': 'docs/ai-agents/DOCKER-BUILD-STATUS.md', // Closest docker doc
  'production': 'docs/COMPLETE_DEPLOYMENT_PLAN.md',
  'use-cases': 'docs/USE_CASES.md',
}

const fetchDoc = async (slug: string | undefined) => {
  loading.value = true
  error.value = null
  content.value = ''

  try {
    const path = docMap[slug || 'undefined'] || 'README.md'
    // Use raw.githubusercontent.com to fetch docs from the main branch for now
    // This ensures we always show latest docs without rebuilding frontend
    const repoBase = 'https://raw.githubusercontent.com/Shannon-Labs/driftlock/main/'
    
    const res = await fetch(`${repoBase}${path}`)
    if (!res.ok) throw new Error(`Failed to load document: ${path}`)
    
    const text = await res.text()
    
    // Fix relative links in markdown
    // Replace (docs/...) with (/docs/...) basically
    // This is a naive replacement for demo purposes
    const fixedText = text.replace(/\]\(docs\/(.*?)\.md\)/g, '](/docs/$1)')
                          .replace(/\]\(\.\.\/docs\/(.*?)\.md\)/g, '](/docs/$1)')
                          
    content.value = marked.parse(fixedText)
  } catch (e: any) {
    error.value = `Could not load documentation: ${e.message}`
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchDoc(route.params.slug as string)
})

watch(() => route.params.slug, (newSlug) => {
  fetchDoc(newSlug as string)
})
</script>

<style>
/* Add github-markdown-css styles if imported globally or scoped here */
/* We rely on @tailwindcss/typography (prose) which is installed in tailwind.config usually */
/* If not, we might need to add it or manually style. Assuming typography plugin is present or minimal styling. */
</style>


