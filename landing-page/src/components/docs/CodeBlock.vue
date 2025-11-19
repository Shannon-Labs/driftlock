<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import hljs from 'highlight.js';
import 'highlight.js/styles/atom-one-dark.css';
import { ClipboardDocumentIcon, CheckIcon } from '@heroicons/vue/24/outline';

const props = defineProps<{
  code: string;
  language?: string;
}>();

const copied = ref(false);
const codeElement = ref<HTMLElement | null>(null);

const highlightCode = () => {
  if (codeElement.value && props.code) {
    // If language is provided, use it, otherwise auto-detect
    if (props.language && hljs.getLanguage(props.language)) {
      codeElement.value.innerHTML = hljs.highlight(props.code, { language: props.language }).value;
    } else {
      codeElement.value.innerHTML = hljs.highlightAuto(props.code).value;
    }
  }
};

onMounted(highlightCode);
watch(() => props.code, highlightCode);

const copyToClipboard = async () => {
  try {
    await navigator.clipboard.writeText(props.code);
    copied.value = true;
    setTimeout(() => {
      copied.value = false;
    }, 2000);
  } catch (err) {
    console.error('Failed to copy code', err);
  }
};
</script>

<template>
  <div class="relative group my-6 rounded-lg overflow-hidden bg-[#282c34] shadow-lg border border-gray-700">
    <div class="flex items-center justify-between px-4 py-2 bg-[#21252b] border-b border-gray-700">
      <span class="text-xs font-mono text-gray-400 uppercase">{{ language || 'text' }}</span>
      <button 
        @click="copyToClipboard"
        class="text-gray-400 hover:text-white transition-colors p-1 rounded hover:bg-gray-700"
        :title="copied ? 'Copied!' : 'Copy code'"
      >
        <component :is="copied ? CheckIcon : ClipboardDocumentIcon" class="h-4 w-4" />
      </button>
    </div>
    <div class="overflow-x-auto p-4">
      <pre><code ref="codeElement" class="font-mono text-sm leading-relaxed text-gray-300 bg-transparent p-0"></code></pre>
    </div>
  </div>
</template>

<style>
/* Override highlight.js background to match our container */
.hljs {
  background: transparent !important;
  padding: 0 !important;
}
</style>
