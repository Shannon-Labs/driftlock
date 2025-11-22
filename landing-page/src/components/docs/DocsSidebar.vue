<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRoute } from 'vue-router';
import { ChevronDownIcon, ChevronRightIcon } from '@heroicons/vue/24/outline';

const route = useRoute();

const navigation = [
  {
    title: 'Getting Started',
    items: [
      { title: 'Quickstart', href: '/docs/deployment/DEPLOYMENT_QUICKSTART' },
      { title: 'Core Concepts', href: '/docs/architecture/ALGORITHMS' },
      { title: 'Authentication', href: '/docs/architecture/API' },
    ]
  },
  {
    title: 'Compliance',
    items: [
      { title: 'DORA', href: '/docs/compliance/COMPLIANCE_DORA' },
      { title: 'Runtime AI', href: '/docs/compliance/COMPLIANCE_RUNTIME_AI' },
      { title: 'US Regulations', href: '/docs/compliance/COMPLIANCE_US' },
    ]
  },
  {
    title: 'Architecture',
    items: [
      { title: 'Overview', href: '/docs/architecture/ARCHITECTURE' },
      { title: 'Deployment', href: '/docs/deployment/DEPLOYMENT' },
      { title: 'Kafka Setup', href: '/docs/deployment/KAFKA_SETUP' },
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