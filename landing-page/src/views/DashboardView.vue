<template>
  <div class="min-h-screen bg-gray-50">
    <DashboardLayout>
      <div class="px-4 sm:px-6 lg:px-8 py-8">
        <!-- Header -->
        <div class="sm:flex sm:items-center sm:justify-between mb-8">
          <div>
            <h1 class="text-3xl font-mono font-bold text-gray-900">Dashboard</h1>
            <p class="mt-2 text-sm text-gray-600">
              Manage your API keys and monitor usage for <span class="font-semibold text-gray-900">{{ authStore.user?.email }}</span>.
            </p>
          </div>
          <div class="mt-4 sm:mt-0 flex space-x-3">
            <a
              href="https://docs.driftlock.net"
              target="_blank"
              class="inline-flex items-center justify-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
            >
              Documentation
            </a>
            <button
              @click="manageBilling"
              class="inline-flex items-center justify-center rounded-md border border-transparent bg-blue-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
            >
              Manage Billing
            </button>
          </div>
        </div>

        <!-- Bento Grid Layout -->
        <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
          
          <!-- Usage Card -->
          <div class="bg-white overflow-hidden shadow rounded-xl border border-gray-100 lg:col-span-2">
            <div class="p-6">
              <h3 class="text-base font-semibold leading-6 text-gray-900">Monthly Usage</h3>
              <div class="mt-2 max-w-xl text-sm text-gray-500">
                <p>Events processed this billing period.</p>
              </div>
              <div class="mt-6">
                <!-- Simple Progress Bar if chart fails or for simplicity -->
                <div class="flex items-center justify-between text-sm font-medium text-gray-900 mb-1">
                  <span>{{ usage.current_period_usage.toLocaleString() }} events</span>
                  <span>{{ usage.plan_limit.toLocaleString() }} limit</span>
                </div>
                <div class="w-full bg-gray-200 rounded-full h-4 overflow-hidden">
                  <div 
                    class="bg-blue-600 h-4 rounded-full transition-all duration-500 ease-out" 
                    :style="{ width: usagePercentage + '%' }"
                  ></div>
                </div>
                <p class="mt-2 text-xs text-gray-500 text-right">
                  {{ usagePercentage.toFixed(1) }}% used
                </p>
                
                <!-- Chart (Placeholder for now, can add vue-chartjs if needed for historical) -->
                <!-- We will stick to the progress bar for "Current Month" as it's cleaner for MVP -->
              </div>
              <div class="mt-6 border-t border-gray-100 pt-4">
                 <dl class="grid grid-cols-1 gap-x-4 gap-y-4 sm:grid-cols-2">
                    <div class="sm:col-span-1">
                      <dt class="text-sm font-medium text-gray-500">Current Plan</dt>
                      <dd class="mt-1 text-lg font-mono font-semibold text-gray-900 uppercase">{{ usage.plan || 'Developer' }}</dd>
                    </div>
                    <div class="sm:col-span-1">
                      <dt class="text-sm font-medium text-gray-500">Next Reset</dt>
                      <dd class="mt-1 text-sm text-gray-900">First of next month</dd>
                    </div>
                 </dl>
              </div>
            </div>
          </div>

          <!-- Integration Card -->
          <div class="bg-slate-900 overflow-hidden shadow rounded-xl border border-slate-800 text-white">
            <div class="p-6">
              <h3 class="text-base font-semibold leading-6 text-white">Quick Integration</h3>
              <p class="mt-2 text-sm text-slate-400">Send your first event in seconds.</p>
              <div class="mt-4 rounded-md bg-black/50 p-3 font-mono text-xs text-blue-300 overflow-x-auto border border-slate-700">
                curl -X POST {{ apiUrl }}/v1/detect \<br>
                &nbsp;&nbsp;-H "Authorization: Bearer {{ firstKey }}" \<br>
                &nbsp;&nbsp;-d @events.json
              </div>
              <div class="mt-4">
                <router-link to="/docs/quickstart" class="text-sm font-medium text-blue-400 hover:text-blue-300">
                  View full guide &rarr;
                </router-link>
              </div>
            </div>
          </div>

          <!-- API Keys Card -->
          <div class="bg-white overflow-hidden shadow rounded-xl border border-gray-100 lg:col-span-3">
            <div class="px-4 py-5 sm:p-6">
              <div class="sm:flex sm:items-center sm:justify-between">
                <h3 class="text-lg font-medium leading-6 text-gray-900">API Keys</h3>
                <div class="mt-4 sm:mt-0">
                  <button
                    type="button"
                    disabled
                    class="inline-flex items-center justify-center rounded-md border border-transparent bg-blue-100 px-4 py-2 text-sm font-medium text-blue-700 hover:bg-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Create New Key (Coming Soon)
                  </button>
                </div>
              </div>
              
              <div class="mt-4 flex flex-col">
                <div class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
                  <div class="inline-block min-w-full py-2 align-middle md:px-6 lg:px-8">
                    <div class="overflow-hidden shadow ring-1 ring-black ring-opacity-5 md:rounded-lg">
                      <table class="min-w-full divide-y divide-gray-300">
                        <thead class="bg-gray-50">
                          <tr>
                            <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-6">Name</th>
                            <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Key Prefix</th>
                            <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Created</th>
                            <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Status</th>
                            <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-6">
                              <span class="sr-only">Actions</span>
                            </th>
                          </tr>
                        </thead>
                        <tbody class="divide-y divide-gray-200 bg-white">
                          <tr v-if="keys.length === 0">
                            <td colspan="5" class="py-8 text-center text-sm text-gray-500">
                              No API keys found. This is unusual.
                            </td>
                          </tr>
                          <tr v-for="key in keys" :key="key.id">
                            <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6">
                              {{ key.name || 'Default Key' }}
                            </td>
                            <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500 font-mono">
                              {{ key.prefix || 'dlk_...' }}
                            </td>
                            <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                              {{ new Date(key.created_at).toLocaleDateString() }}
                            </td>
                            <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                              <span class="inline-flex items-center rounded-full bg-green-100 px-2.5 py-0.5 text-xs font-medium text-green-800">
                                {{ key.status }}
                              </span>
                            </td>
                            <td class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-6">
                              <!-- Revoke button placeholder -->
                              <button class="text-red-600 hover:text-red-900">Revoke</button>
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
          <div class="bg-white overflow-hidden shadow rounded-xl border border-gray-100 lg:col-span-3">
             <div class="px-4 py-5 sm:px-6 border-b border-gray-100">
               <h3 class="text-lg font-medium leading-6 text-gray-900">Recent Anomalies</h3>
             </div>
             <div class="px-4 py-5 sm:p-6 text-center text-gray-500 text-sm">
                No anomalies detected in the last 24 hours.
                <br>
                <span class="text-xs text-gray-400">Real-time feed integration coming soon.</span>
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
