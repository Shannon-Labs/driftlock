<script setup lang="ts">
import { ref, watch, onMounted, nextTick } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { marked } from 'marked';
import hljs from 'highlight.js';
import 'highlight.js/styles/atom-one-dark.css';
import DocsSidebar from '../components/docs/DocsSidebar.vue';

const route = useRoute();
const router = useRouter();
const content = ref('');
const loading = ref(true);
const error = ref<string | null>(null);

// Configure marked
marked.setOptions({
  gfm: true,
  breaks: true
});

const fetchDoc = async (path: string) => {
  loading.value = true;
  error.value = null;
  
  try {
    // Remove '/docs' prefix and ensure .md extension
    let docPath = path.replace(/^\/docs/, '');
    if (docPath === '' || docPath === '/') docPath = '/README';
    if (!docPath.endsWith('.md')) docPath += '.md';
    
    const response = await fetch(`/docs${docPath}`);
    
    if (!response.ok) {
      throw new Error(`Documentation not found: ${docPath}`);
    }
    
    const text = await response.text();
    content.value = await marked(text);
    
    // Scroll to top
    window.scrollTo(0, 0);
    
    // Highlight code blocks after render
    nextTick(() => {
      document.querySelectorAll('pre code').forEach((block) => {
        hljs.highlightElement(block as HTMLElement);
      });
    });
    
  } catch (err) {
    console.error(err);
    error.value = 'Failed to load documentation. Please try again later.';
  } finally {
    loading.value = false;
  }
};

watch(() => route.path, (newPath) => {
  if (newPath.startsWith('/docs')) {
    fetchDoc(newPath);
  }
}, { immediate: true });

</script>

<template>
  <div class="min-h-screen bg-white dark:bg-gray-900 pt-16">
    <div class="flex max-w-8xl mx-auto">
      <!-- Sidebar -->
      <DocsSidebar />
      
      <!-- Main Content -->
      <main class="flex-1 min-w-0 px-4 sm:px-6 lg:px-8 py-10">
        <div v-if="loading" class="flex justify-center py-20">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
        
        <div v-else-if="error" class="text-center py-20">
          <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-4">404 - Not Found</h2>
          <p class="text-gray-600 dark:text-gray-400">{{ error }}</p>
          <router-link to="/docs" class="text-blue-600 hover:underline mt-4 inline-block">
            Return to Documentation Home
          </router-link>
        </div>
        
        <div v-else class="prose prose-blue dark:prose-invert max-w-4xl mx-auto">
          <div v-html="content"></div>
        </div>
      </main>
    </div>
  </div>
</template>

<style>
/* Custom styles for markdown content */
.prose pre {
  background-color: #282c34;
  border-radius: 0.5rem;
  padding: 1rem;
  overflow-x: auto;
}

.prose code {
  color: #e06c75;
  background-color: rgba(40, 44, 52, 0.1);
  padding: 0.2em 0.4em;
  border-radius: 0.25rem;
  font-size: 0.875em;
}

.dark .prose code {
  background-color: rgba(255, 255, 255, 0.1);
}

.prose pre code {
  color: inherit;
  background-color: transparent;
  padding: 0;
  font-size: 0.875em;
}

.prose h1 {
  font-size: 2.25rem;
  font-weight: 800;
  margin-bottom: 2rem;
  background: linear-gradient(90deg, #0284c7 0%, #3b82f6 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.prose h2 {
  font-size: 1.5rem;
  font-weight: 700;
  margin-top: 2.5rem;
  margin-bottom: 1rem;
  border-bottom: 1px solid #e5e7eb;
  padding-bottom: 0.5rem;
}

.dark .prose h2 {
  border-color: #374151;
}

.prose a {
  color: #2563eb;
  text-decoration: none;
}

.prose a:hover {
  text-decoration: underline;
}

.prose table {
  width: 100%;
  border-collapse: collapse;
  margin: 1.5rem 0;
}

.prose th, .prose td {
  padding: 0.75rem;
  border: 1px solid #e5e7eb;
  text-align: left;
}

.dark .prose th, .dark .prose td {
  border-color: #374151;
}

.prose th {
  background-color: #f9fafb;
  font-weight: 600;
}

.dark .prose th {
  background-color: #1f2937;
}
</style>
