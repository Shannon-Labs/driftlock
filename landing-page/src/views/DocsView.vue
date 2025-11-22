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
  <div class="min-h-screen bg-background pt-0">
    <div class="flex max-w-8xl mx-auto">
      <!-- Sidebar -->
      <DocsSidebar />
      
      <!-- Main Content -->
      <main class="flex-1 min-w-0 px-4 sm:px-6 lg:px-8 py-10">
        <div v-if="loading" class="flex justify-center py-20">
          <div class="animate-spin h-8 w-8 border-2 border-black border-t-transparent"></div>
        </div>
        
        <div v-else-if="error" class="text-center py-20">
          <h2 class="text-2xl font-bold font-sans uppercase text-foreground mb-4">404 - Not Found</h2>
          <p class="text-gray-600 font-mono mb-6">{{ error }}</p>
          <router-link to="/docs" class="brutalist-button">
            Return to Documentation Home
          </router-link>
        </div>
        
        <div v-else class="prose max-w-4xl mx-auto">
          <div v-html="content"></div>
        </div>
      </main>
    </div>
  </div>
</template>

<style>
/* Custom styles for markdown content - Brutalist Academic Edition */
.prose {
  color: var(--color-foreground);
  font-family: "EB Garamond", serif;
  line-height: 1.6;
}

.prose pre {
  background-color: #1a1a1a;
  border: 1px solid #000;
  border-radius: 0;
  padding: 1rem;
  overflow-x: auto;
  margin: 1.5rem 0;
}

.prose code {
  color: #d63384;
  background-color: #f3f4f6;
  padding: 0.2em 0.4em;
  border: 1px solid #e5e7eb;
  font-family: "JetBrains Mono", monospace;
  font-size: 0.875em;
}

.prose pre code {
  color: inherit;
  background-color: transparent;
  padding: 0;
  border: none;
  font-size: 0.875em;
}

.prose h1 {
  font-family: "Inter", sans-serif;
  font-size: 3rem;
  font-weight: 800;
  letter-spacing: -0.02em;
  margin-bottom: 2rem;
  text-transform: uppercase;
  border-bottom: 1px solid #000;
  padding-bottom: 1rem;
}

.prose h2 {
  font-family: "Inter", sans-serif;
  font-size: 1.8rem;
  font-weight: 700;
  margin-top: 3rem;
  margin-bottom: 1.5rem;
  border-bottom: 1px solid #000;
  padding-bottom: 0.5rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.prose h3 {
  font-family: "Inter", sans-serif;
  font-size: 1.4rem;
  font-weight: 600;
  margin-top: 2rem;
  margin-bottom: 1rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.prose p {
  margin-bottom: 1.25rem;
  font-size: 1.125rem;
}

.prose a {
  color: #000;
  text-decoration: underline;
  text-decoration-thickness: 1px;
  text-underline-offset: 4px;
  font-weight: 600;
  transition: all 0.2s;
}

.prose a:hover {
  background-color: #000;
  color: #fff;
  text-decoration: none;
}

.prose ul, .prose ol {
  margin: 1.5rem 0;
  padding-left: 1.5rem;
}

.prose li {
  margin-bottom: 0.5rem;
}

.prose table {
  width: 100%;
  border-collapse: collapse;
  margin: 2rem 0;
  border: 1px solid #000;
  font-family: "JetBrains Mono", monospace;
  font-size: 0.85rem;
}

.prose th, .prose td {
  padding: 0.75rem;
  border: 1px solid #000;
  text-align: left;
}

.prose th {
  background-color: #000;
  color: #fff;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.prose tr:nth-child(even) {
  background-color: #f9fafb;
}

.prose blockquote {
  border-left: 2px solid #000;
  padding-left: 1.5rem;
  margin: 1.5rem 0;
  font-style: italic;
  background-color: #f3f4f6;
  padding: 1rem;
}
</style>
