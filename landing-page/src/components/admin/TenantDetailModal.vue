<template>
  <div class="relative z-10" aria-labelledby="modal-title" role="dialog" aria-modal="true">
    <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"></div>

    <div class="fixed inset-0 z-10 overflow-y-auto">
      <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
        <div class="relative transform overflow-hidden rounded-lg bg-white px-4 pt-5 pb-4 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-4xl sm:p-6">
          <div class="absolute top-0 right-0 hidden pt-4 pr-4 sm:block">
            <button @click="$emit('close')" type="button" class="rounded-md bg-white text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
              <span class="sr-only">Close</span>
              <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <div v-if="loading" class="text-center py-12">
            Loading...
          </div>
          <div v-else-if="tenant">
            <div class="sm:flex sm:items-start">
              <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left w-full">
                <h3 class="text-lg font-medium leading-6 text-gray-900" id="modal-title">{{ tenant.name }}</h3>
                <div class="mt-2 text-sm text-gray-500">
                  <p>ID: {{ tenant.id }}</p>
                  <p>Email: {{ tenant.email }}</p>
                  <p>Plan: {{ tenant.plan }}</p>
                  <p>Created: {{ new Date(tenant.created_at).toLocaleString() }}</p>
                </div>

                <div class="mt-6 border-t border-gray-200 pt-4">
                  <h4 class="text-md font-medium text-gray-900">Status</h4>
                  <div class="mt-2 flex items-center space-x-4">
                    <span class="inline-flex rounded-full px-2 text-xs font-semibold leading-5" :class="statusClass(tenant.status)">
                      {{ tenant.status }}
                    </span>
                    <button v-if="tenant.status === 'active'" @click="updateStatus('suspended')" class="text-red-600 hover:text-red-800 text-sm">Suspend</button>
                    <button v-if="tenant.status === 'suspended'" @click="updateStatus('active')" class="text-green-600 hover:text-green-800 text-sm">Activate</button>
                  </div>
                </div>

                <div class="mt-6 border-t border-gray-200 pt-4">
                  <h4 class="text-md font-medium text-gray-900">Usage (Last 30 Days)</h4>
                  <div class="mt-4 h-48 flex items-end space-x-1">
                    <div v-for="day in usage" :key="day.date" class="flex-1 flex flex-col items-center group relative">
                      <div
                        class="w-full bg-indigo-500 hover:bg-indigo-600 rounded-t"
                        :style="{ height: barHeight(day.event_count) + '%' }"
                      ></div>
                      <span class="absolute bottom-full mb-1 hidden group-hover:block bg-black text-white text-xs rounded p-1 z-10 whitespace-nowrap">
                        {{ day.date }}: {{ day.event_count }} events
                      </span>
                    </div>
                  </div>
                  <div class="mt-2 text-center text-xs text-gray-500">
                    Events per day
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'

const props = defineProps<{
  tenantId: string
}>()

const emit = defineEmits(['close'])

const tenant = ref<any>(null)
const loading = ref(true)

const usage = computed(() => tenant.value?.usage_metrics || [])

const barHeight = (count: number) => {
  if (usage.value.length === 0) return 0
  const max = Math.max(...usage.value.map((u: any) => u.event_count))
  if (max === 0) return 0
  return (count / max) * 100
}

const statusClass = (status: string) => {
  switch (status) {
    case 'active':
      return 'bg-green-100 text-green-800'
    case 'suspended':
      return 'bg-red-100 text-red-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

const fetchTenant = async () => {
  loading.value = true
  try {
    const key = localStorage.getItem('driftlock_admin_key')
    const res = await fetch(`/v1/admin/tenants/${props.tenantId}`, {
      headers: { 'X-Admin-Key': key || '' }
    })
    if (res.ok) {
      tenant.value = await res.json()
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const updateStatus = async (newStatus: string) => {
  try {
    const key = localStorage.getItem('driftlock_admin_key')
    const res = await fetch(`/v1/admin/tenants/${props.tenantId}/status`, {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
        'X-Admin-Key': key || ''
      },
      body: JSON.stringify({ status: newStatus })
    })
    if (res.ok) {
      tenant.value.status = newStatus
    }
  } catch (e) {
    console.error(e)
  }
}

onMounted(fetchTenant)
</script>
