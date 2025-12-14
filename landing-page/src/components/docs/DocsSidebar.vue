<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRoute } from 'vue-router';
import { ChevronDownIcon, ChevronRightIcon } from '@heroicons/vue/24/outline';

const route = useRoute();

const navigation = [
  {
    title: 'Getting Started',
    items: [
      { title: 'Quickstart', href: '/docs/user-guide/getting-started/quickstart' },
      { title: 'Core Concepts', href: '/docs/user-guide/getting-started/concepts' },
      { title: 'Authentication', href: '/docs/user-guide/getting-started/authentication' },
    ]
  },
  {
    title: 'API Reference',
    items: [
      { title: 'POST /v1/detect', href: '/docs/user-guide/api/endpoints/detect' },
      { title: 'POST /v1/demo/detect', href: '/docs/user-guide/api/endpoints/demo' },
      { title: 'GET /v1/anomalies', href: '/docs/user-guide/api/endpoints/anomalies' },
      { title: 'Compression Algorithms', href: '/docs/architecture/ALGORITHMS' },
      { title: 'Error Codes', href: '/docs/user-guide/api/errors' },
    ]
  },
  {
    title: 'Examples',
    items: [
      { title: 'cURL', href: '/docs/user-guide/api/examples/curl-examples' },
      { title: 'Python', href: '/docs/user-guide/api/examples/python-examples' },
    ]
  },
  {
    title: 'SDKs',
    items: [
      { title: 'Overview', href: '/docs/sdk/README' },
      { title: 'Node.js', href: '/docs/sdk/nodejs' },
      { title: 'Python', href: '/docs/sdk/python' },
    ]
  },
  {
    title: 'Compliance',
    items: [
      { title: 'DORA', href: '/docs/compliance/COMPLIANCE_DORA' },
      { title: 'NIS2', href: '/docs/compliance/COMPLIANCE_NIS2' },
      { title: 'AI Act', href: '/docs/compliance/COMPLIANCE_RUNTIME_AI' },
    ]
  },
];

const openSections = ref<Set<string>>(new Set(navigation.map(n => n.title)));

const toggleSection = (title: string) => {
  if (openSections.value.has(title)) {
    openSections.value.delete(title);
  } else {
    openSections.value.add(title);
  }
};

const isActive = (href: string) => {
    // Handle potential trailing slashes or .md extensions if needed
    return route.path === href || route.path === href.replace('.md', '');
};
</script>

<template>
  <nav class="w-64 flex-shrink-0 border-r border-black h-[calc(100vh-5rem)] overflow-y-auto sticky top-20 bg-white py-6 px-4 hidden lg:block">
    <div class="space-y-8">
      <div v-for="section in navigation" :key="section.title">
        <button 
          @click="toggleSection(section.title)"
          class="flex items-center justify-between w-full text-sm font-bold font-sans uppercase tracking-widest text-black mb-4 hover:underline decoration-1 underline-offset-2"
        >
          {{ section.title }}
          <component 
            :is="openSections.has(section.title) ? ChevronDownIcon : ChevronRightIcon" 
            class="h-4 w-4 text-black border border-black p-0.5"
          />
        </button>
        
        <ul v-show="openSections.has(section.title)" class="space-y-1 pl-2 border-l border-black ml-1">
          <li v-for="item in section.items" :key="item.href">
            <router-link 
              :to="item.href"
              class="block px-3 py-1 text-sm transition-all duration-200 font-mono border-l-2 border-transparent"
              :class="[
                isActive(item.href) 
                  ? 'border-black font-bold pl-4 bg-gray-100' 
                  : 'text-gray-600 hover:text-black hover:pl-4 hover:bg-gray-50'
              ]"
            >
              {{ item.title }}
            </router-link>
          </li>
        </ul>
      </div>
    </div>
  </nav>
</template>