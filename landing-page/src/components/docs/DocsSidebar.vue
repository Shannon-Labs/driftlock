<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRoute } from 'vue-router';
import { ChevronDownIcon, ChevronRightIcon } from '@heroicons/vue/24/outline';

const route = useRoute();

const navigation = [
  {
    title: 'Getting Started',
    items: [
      { title: 'Quickstart', href: '/docs/getting-started/quickstart' },
      { title: 'Core Concepts', href: '/docs/getting-started/concepts' },
      { title: 'Authentication', href: '/docs/getting-started/authentication' },
    ]
  },
  {
    title: 'REST API',
    items: [
      { title: 'Overview', href: '/docs/api/rest-api' },
      { title: 'POST /detect', href: '/docs/api/endpoints/detect' },
      { title: 'GET /anomalies', href: '/docs/api/endpoints/anomalies' },
      { title: 'GET /anomalies/{id}', href: '/docs/api/endpoints/anomaly-detail' },
    ]
  },
  {
    title: 'Code Examples',
    items: [
      { title: 'cURL', href: '/docs/api/examples/curl-examples' },
      { title: 'Python', href: '/docs/api/examples/python-examples' },
    ]
  },
  {
    title: 'SDKs',
    items: [
      { title: 'Data Connect', href: '/docs/sdks/data-connect' },
    ]
  }
];

const openSections = ref<Set<string>>(new Set(navigation.map(n => n.title)));

const toggleSection = (title: string) => {
  if (openSections.value.has(title)) {
    openSections.value.delete(title);
  } else {
    openSections.value.add(title);
  }
};

const isActive = (href: string) => route.path === href;
</script>

<template>
  <nav class="w-64 flex-shrink-0 border-r border-gray-200 dark:border-gray-800 h-[calc(100vh-4rem)] overflow-y-auto sticky top-16 bg-white dark:bg-gray-900 py-6 px-4 hidden lg:block">
    <div class="space-y-6">
      <div v-for="section in navigation" :key="section.title">
        <button 
          @click="toggleSection(section.title)"
          class="flex items-center justify-between w-full text-sm font-semibold text-gray-900 dark:text-white mb-2 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
        >
          {{ section.title }}
          <component 
            :is="openSections.has(section.title) ? ChevronDownIcon : ChevronRightIcon" 
            class="h-4 w-4 text-gray-500"
          />
        </button>
        
        <ul v-show="openSections.has(section.title)" class="space-y-1">
          <li v-for="item in section.items" :key="item.href">
            <router-link 
              :to="item.href"
              class="block px-2 py-1.5 text-sm rounded-md transition-colors duration-200"
              :class="[
                isActive(item.href) 
                  ? 'bg-blue-50 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300 font-medium' 
                  : 'text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-800 hover:text-gray-900 dark:hover:text-white'
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
