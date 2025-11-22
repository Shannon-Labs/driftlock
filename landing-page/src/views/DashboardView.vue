<template>
  <div class="min-h-screen bg-white">
    <DashboardLayout>
      <div class="px-4 sm:px-6 lg:px-8 py-8">
        <!-- Header -->
        <div class="sm:flex sm:items-center sm:justify-between mb-12 border-b-4 border-black pb-6">
          <div>
            <h1 class="text-4xl font-sans font-black uppercase tracking-tighter text-black">Dashboard</h1>
            <p class="mt-2 text-sm font-mono text-gray-600">
              Manage your API keys and monitor usage for <span class="font-bold text-black border-b-2 border-black">{{ authStore.user?.email }}</span>.
            </p>
          </div>
          <div class="mt-4 sm:mt-0 flex space-x-4">
            <a
              href="https://docs.driftlock.net"
              target="_blank"
              class="inline-flex items-center justify-center border-2 border-black bg-white px-6 py-3 text-sm font-bold uppercase tracking-widest text-black hover:bg-black hover:text-white transition-colors shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:shadow-none hover:translate-x-[2px] hover:translate-y-[2px]"
            >
              Documentation
            </a>
            <button
              @click="manageBilling"
              class="inline-flex items-center justify-center border-2 border-black bg-black px-6 py-3 text-sm font-bold uppercase tracking-widest text-white hover:bg-white hover:text-black transition-colors shadow-[4px_4px_0px_0px_rgba(0,0,0,0)] hover:shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]"
            >
              Manage Billing
            </button>
          </div>
        </div>

        <!-- Bento Grid Layout -->
        <div class="grid grid-cols-1 gap-8 lg:grid-cols-3">
          
          <!-- Usage Card -->
          <div class="bg-white border-2 border-black p-6 lg:col-span-2 shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
            <div>
              <h3 class="text-xl font-bold uppercase tracking-wide text-black border-b-2 border-black pb-2 mb-4">Monthly Usage</h3>
              <div class="mt-2 max-w-xl text-sm font-serif text-gray-800">
                <p>Events processed this billing period.</p>
              </div>
              <div class="mt-8">
                <!-- Simple Progress Bar -->
                <div class="flex items-center justify-between text-sm font-bold font-mono text-black mb-2">
                  <span>{{ usage.current_period_usage.toLocaleString() }} EVENTS</span>
                  <span>{{ usage.plan_limit.toLocaleString() }} LIMIT</span>
                </div>
                <div class="w-full bg-gray-200 border-2 border-black h-6 relative">
                  <div 
                    class="bg-black h-full absolute top-0 left-0" 
                    :style="{ width: usagePercentage + '%' }"
                  ></div>
                </div>
                <p class="mt-2 text-xs font-mono font-bold text-black text-right">
                  {{ usagePercentage.toFixed(1) }}% USED
                </p>
              </div>
              <div class="mt-8 border-t-2 border-black pt-4">
                 <dl class="grid grid-cols-1 gap-x-4 gap-y-4 sm:grid-cols-2">
                    <div class="sm:col-span-1">
                      <dt class="text-xs font-bold uppercase tracking-widest text-gray-500">Current Plan</dt>
                      <dd class="mt-1 text-2xl font-sans font-black text-black uppercase">{{ usage.plan || 'Developer' }}</dd>
                    </div>
                    <div class="sm:col-span-1">
                      <dt class="text-xs font-bold uppercase tracking-widest text-gray-500">Next Reset</dt>
                      <dd class="mt-1 text-sm font-mono font-bold text-black">First of next month</dd>
                    </div>
                 </dl>
              </div>
            </div>
          </div>

          <!-- Integration Card -->
          <div class="bg-black text-white border-2 border-black p-6 shadow-[8px_8px_0px_0px_rgba(128,128,128,1)]">
            <div>
              <h3 class="text-xl font-bold uppercase tracking-wide text-white border-b-2 border-white pb-2 mb-4">Quick Integration</h3>
              <p class="mt-2 text-sm font-serif text-gray-300">Send your first event in seconds.</p>
              <div class="mt-4 bg-gray-900 p-4 font-mono text-xs text-green-400 overflow-x-auto border border-gray-700">
                curl -X POST {{ apiUrl }}/v1/detect \<br>
                &nbsp;&nbsp;-H "Authorization: Bearer {{ firstKey }}" \<br>
                &nbsp;&nbsp;-d @events.json
              </div>
              <div class="mt-6">
                <router-link to="/docs/quickstart" class="text-sm font-bold uppercase tracking-widest text-white hover:underline decoration-2 underline-offset-4">
                  View full guide &rarr;
                </router-link>
              </div>
            </div>
          </div>

          <!-- API Keys Card -->
          <div class="bg-white border-2 border-black lg:col-span-3 shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
            <div class="px-6 py-6">
              <div class="sm:flex sm:items-center sm:justify-between border-b-2 border-black pb-4 mb-4">
                <h3 class="text-xl font-bold uppercase tracking-wide text-black">API Keys</h3>
                <div class="mt-4 sm:mt-0">
                  <button
                    type="button"
                    disabled
                    class="inline-flex items-center justify-center border-2 border-black bg-gray-100 px-4 py-2 text-sm font-bold uppercase tracking-widest text-gray-500 cursor-not-allowed"
                  >
                    Create New Key (Coming Soon)
                  </button>
                </div>
              </div>
              
              <div class="mt-4 flex flex-col">
                <div class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
                  <div class="inline-block min-w-full py-2 align-middle md:px-6 lg:px-8">
                    <div class="overflow-hidden">
                      <table class="min-w-full divide-y-2 divide-black border-2 border-black">
                        <thead class="bg-black">
                          <tr>
                            <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-bold uppercase tracking-widest text-white sm:pl-6">Name</th>
                            <th scope="col" class="px-3 py-3.5 text-left text-sm font-bold uppercase tracking-widest text-white">Key Prefix</th>
                            <th scope="col" class="px-3 py-3.5 text-left text-sm font-bold uppercase tracking-widest text-white">Created</th>
                            <th scope="col" class="px-3 py-3.5 text-left text-sm font-bold uppercase tracking-widest text-white">Status</th>
                            <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-6">
                              <span class="sr-only">Actions</span>
                            </th>
                          </tr>
                        </thead>
                        <tbody class="divide-y-2 divide-black bg-white">
                          <tr v-if="keys.length === 0">
                            <td colspan="5" class="py-8 text-center text-sm font-mono text-gray-500">
                              No API keys found. This is unusual.
                            </td>
                          </tr>
                          <tr v-for="key in keys" :key="key.id" class="hover:bg-gray-50">
                            <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-bold text-black sm:pl-6">
                              {{ key.name || 'Default Key' }}
                            </td>
                            <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-600 font-mono">
                              {{ key.prefix || 'dlk_...' }}
                            </td>
                            <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-600 font-mono">
                              {{ new Date(key.created_at).toLocaleDateString() }}
                            </td>
                            <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                              <span class="inline-flex items-center border border-black bg-green-100 px-2 py-0.5 text-xs font-bold uppercase tracking-wide text-green-800">
                                {{ key.status }}
                              </span>
                            </td>
                            <td class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-6">
                              <!-- Revoke button placeholder -->
                              <button class="text-red-600 font-bold uppercase tracking-wider hover:text-red-900 hover:underline decoration-2 underline-offset-4">Revoke</button>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Recent Anomalies -->
          <div class="bg-white border-2 border-black lg:col-span-3 shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
             <div class="px-6 py-5 border-b-2 border-black">
               <h3 class="text-xl font-bold uppercase tracking-wide text-black">Recent Anomalies</h3>
             </div>
             <div class="px-6 py-12 text-center text-gray-500 text-sm font-mono border-b border-transparent">
                No anomalies detected in the last 24 hours.
                <br>
                <span class="text-xs text-gray-400 uppercase tracking-widest mt-2 block">Real-time feed integration coming soon</span>
             </div>
          </div>

        </div>
      </div>
    </DashboardLayout>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import DashboardLayout from '../layouts/DashboardLayout.vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const keys = ref<any[]>([])
const usage = ref({
  current_period_usage: 0,
  plan_limit: 10000,
  plan: 'developer'
})

const apiUrl = window.location.origin // Assuming API is on same domain
const firstKey = computed(() => keys.value.length > 0 ? (keys.value[0].prefix || 'dlk_...') : 'YOUR_API_KEY')

const usagePercentage = computed(() => {
  if (usage.value.plan_limit === 0) return 0
  return Math.min(100, (usage.value.current_period_usage / usage.value.plan_limit) * 100)
})

onMounted(async () => {
  try {
    const token = await authStore.getToken()
    if (!token) return

    // Fetch Keys
    const resKeys = await fetch('/api/v1/me/keys', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (resKeys.ok) {
      const data = await resKeys.json()
      keys.value = data.keys || []
    }

    // Fetch Usage
    const resUsage = await fetch('/api/v1/me/usage', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (resUsage.ok) {
      const data = await resUsage.json()
      usage.value = data
    }

  } catch (e) {
    console.error('Failed to fetch dashboard data', e)
  }
})

const manageBilling = () => {
  window.location.href = '/api/v1/billing/portal'
}
</script>
