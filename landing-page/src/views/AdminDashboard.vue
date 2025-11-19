<template>
  <div class="min-h-screen bg-gray-100">
    <nav class="bg-white shadow-sm">
      <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div class="flex h-16 justify-between">
          <div class="flex">
            <div class="flex flex-shrink-0 items-center">
              <span class="font-bold text-xl">Driftlock Admin</span>
            </div>
          </div>
          <div class="flex items-center">
            <button @click="logout" class="text-gray-500 hover:text-gray-700">Logout</button>
          </div>
        </div>
      </div>
    </nav>

    <div class="py-10">
      <header>
        <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <h1 class="text-3xl font-bold leading-tight tracking-tight text-gray-900">Dashboard</h1>
        </div>
      </header>
      <main>
        <div class="mx-auto max-w-7xl sm:px-6 lg:px-8">
          <!-- Stats -->
          <div class="mt-8 grid grid-cols-1 gap-5 sm:grid-cols-3">
            <div class="overflow-hidden rounded-lg bg-white shadow">
              <div class="px-4 py-5 sm:p-6">
                <dt class="truncate text-sm font-medium text-gray-500">Total Tenants</dt>
                <dd class="mt-1 text-3xl font-semibold tracking-tight text-gray-900">{{ stats.total_tenants }}</dd>
              </div>
            </div>
            <div class="overflow-hidden rounded-lg bg-white shadow">
              <div class="px-4 py-5 sm:p-6">
                <dt class="truncate text-sm font-medium text-gray-500">Active Subscriptions</dt>
                <dd class="mt-1 text-3xl font-semibold tracking-tight text-gray-900">{{ stats.active_subscriptions }}</dd>
              </div>
            </div>
            <div class="overflow-hidden rounded-lg bg-white shadow">
              <div class="px-4 py-5 sm:p-6">
                <dt class="truncate text-sm font-medium text-gray-500">Anomalies (24h)</dt>
                <dd class="mt-1 text-3xl font-semibold tracking-tight text-gray-900">{{ stats.total_anomalies_24h }}</dd>
              </div>
            </div>
          </div>

          <!-- Tenant List -->
          <div class="mt-8">
            <TenantTable :tenants="tenants" @view="openTenantModal" />
          </div>
        </div>
      </main>
    </div>

    <TenantDetailModal v-if="selectedTenantId" :tenant-id="selectedTenantId" @close="selectedTenantId = null" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import TenantTable from '../components/admin/TenantTable.vue'
import TenantDetailModal from '../components/admin/TenantDetailModal.vue'

const router = useRouter()
const stats = ref({
  total_tenants: 0,
  active_subscriptions: 0,
  total_anomalies_24h: 0
})
const tenants = ref([])
const selectedTenantId = ref<string | null>(null)

const fetchStats = async () => {
  try {
    const key = localStorage.getItem('driftlock_admin_key')
    const res = await fetch('/v1/admin/stats', {
      headers: { 'X-Admin-Key': key || '' }
    })
    if (res.ok) {
      stats.value = await res.json()
    }
  } catch (e) {
    console.error(e)
  }
}

const fetchTenants = async () => {
  try {
    const key = localStorage.getItem('driftlock_admin_key')
    const res = await fetch('/v1/admin/tenants', {
      headers: { 'X-Admin-Key': key || '' }
    })
    if (res.ok) {
      const data = await res.json()
      tenants.value = data.tenants
    }
  } catch (e) {
    console.error(e)
  }
}

const logout = () => {
  localStorage.removeItem('driftlock_admin_key')
  router.push('/admin/login')
}

const openTenantModal = (tenant: any) => {
  selectedTenantId.value = tenant.id
}

onMounted(() => {
  const key = localStorage.getItem('driftlock_admin_key')
  if (!key) {
    router.push('/admin/login')
    return
  }
  fetchStats()
  fetchTenants()
})
</script>
