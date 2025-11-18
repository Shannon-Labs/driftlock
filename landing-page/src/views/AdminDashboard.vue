<template>
    <main class="min-h-screen bg-gray-100">
        <!-- Header -->
        <header class="bg-white shadow">
            <div class="container mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                <div class="flex h-16 items-center justify-between">
                    <h1 class="text-xl font-mono font-bold text-gray-900">Driftlock Admin</h1>
                    <div class="flex items-center gap-4">
                        <span class="text-sm text-gray-500">{{ tenants.length }} tenants</span>
                        <button
                            @click="refreshData"
                            :disabled="loading"
                            class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white hover:bg-blue-700 disabled:opacity-50"
                        >
                            {{ loading ? 'Loading...' : 'Refresh' }}
                        </button>
                    </div>
                </div>
            </div>
        </header>

        <!-- Auth Modal -->
        <div v-if="!authenticated" class="fixed inset-0 bg-gray-900/50 flex items-center justify-center z-50">
            <div class="bg-white rounded-2xl p-8 shadow-2xl w-full max-w-md">
                <h2 class="text-xl font-mono font-bold text-gray-900">Admin Authentication</h2>
                <p class="mt-2 text-sm text-gray-600">Enter your admin key to access the dashboard.</p>
                <form @submit.prevent="authenticate" class="mt-6 space-y-4">
                    <div>
                        <label for="admin-key" class="block text-sm font-semibold text-gray-700">Admin Key</label>
                        <input
                            id="admin-key"
                            v-model="adminKey"
                            type="password"
                            required
                            class="mt-2 w-full rounded-xl border border-gray-200 px-4 py-3 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/40"
                            placeholder="Enter admin key"
                        />
                    </div>
                    <button
                        type="submit"
                        class="w-full rounded-xl bg-blue-600 px-4 py-3 font-semibold text-white hover:bg-blue-700"
                    >
                        Sign In
                    </button>
                    <p v-if="authError" class="text-sm text-red-600">{{ authError }}</p>
                </form>
            </div>
        </div>

        <!-- Main Content -->
        <div v-if="authenticated" class="container mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
            <!-- Stats Cards -->
            <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-4 mb-8">
                <div class="rounded-xl bg-white p-6 shadow">
                    <p class="text-sm font-semibold text-gray-500">Total Tenants</p>
                    <p class="mt-2 text-3xl font-mono font-bold text-gray-900">{{ tenants.length }}</p>
                </div>
                <div class="rounded-xl bg-white p-6 shadow">
                    <p class="text-sm font-semibold text-gray-500">Verified</p>
                    <p class="mt-2 text-3xl font-mono font-bold text-green-600">{{ verifiedCount }}</p>
                </div>
                <div class="rounded-xl bg-white p-6 shadow">
                    <p class="text-sm font-semibold text-gray-500">Trial</p>
                    <p class="mt-2 text-3xl font-mono font-bold text-blue-600">{{ trialCount }}</p>
                </div>
                <div class="rounded-xl bg-white p-6 shadow">
                    <p class="text-sm font-semibold text-gray-500">Paid</p>
                    <p class="mt-2 text-3xl font-mono font-bold text-purple-600">{{ paidCount }}</p>
                </div>
            </div>

            <!-- Search and Filter -->
            <div class="mb-6 flex flex-col sm:flex-row gap-4">
                <div class="flex-1">
                    <input
                        v-model="searchQuery"
                        type="text"
                        placeholder="Search by name, email, or slug..."
                        class="w-full rounded-xl border border-gray-200 px-4 py-3 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/40"
                    />
                </div>
                <select
                    v-model="filterPlan"
                    class="rounded-xl border border-gray-200 px-4 py-3 focus:border-blue-500 focus:outline-none"
                >
                    <option value="">All Plans</option>
                    <option value="trial">Trial</option>
                    <option value="starter">Starter</option>
                    <option value="growth">Growth</option>
                    <option value="enterprise">Enterprise</option>
                </select>
            </div>

            <!-- Tenants Table -->
            <div class="rounded-xl bg-white shadow overflow-hidden">
                <div class="overflow-x-auto">
                    <table class="w-full">
                        <thead class="bg-gray-50 border-b border-gray-200">
                            <tr>
                                <th class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
                                    Tenant
                                </th>
                                <th class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
                                    Email
                                </th>
                                <th class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
                                    Plan
                                </th>
                                <th class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
                                    Status
                                </th>
                                <th class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
                                    Created
                                </th>
                                <th class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
                                    Actions
                                </th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-gray-200">
                            <tr v-for="tenant in filteredTenants" :key="tenant.id" class="hover:bg-gray-50">
                                <td class="px-6 py-4 whitespace-nowrap">
                                    <div>
                                        <p class="font-semibold text-gray-900">{{ tenant.name }}</p>
                                        <p class="text-xs text-gray-500 font-mono">{{ tenant.slug }}</p>
                                    </div>
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                                    {{ tenant.email || '-' }}
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap">
                                    <span :class="planBadgeClass(tenant.plan)" class="px-2 py-1 text-xs font-semibold rounded-full">
                                        {{ tenant.plan }}
                                    </span>
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap">
                                    <span :class="statusBadgeClass(tenant.status)" class="px-2 py-1 text-xs font-semibold rounded-full">
                                        {{ tenant.status }}
                                    </span>
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                                    {{ formatDate(tenant.created_at) }}
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap">
                                    <button
                                        @click="viewTenantUsage(tenant.id)"
                                        class="text-blue-600 hover:text-blue-800 text-sm font-semibold"
                                    >
                                        View Usage
                                    </button>
                                </td>
                            </tr>
                            <tr v-if="filteredTenants.length === 0">
                                <td colspan="6" class="px-6 py-8 text-center text-gray-500">
                                    No tenants found
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>

            <!-- Usage Modal -->
            <div v-if="selectedTenant" class="fixed inset-0 bg-gray-900/50 flex items-center justify-center z-50">
                <div class="bg-white rounded-2xl p-8 shadow-2xl w-full max-w-lg">
                    <div class="flex items-center justify-between mb-6">
                        <h2 class="text-xl font-mono font-bold text-gray-900">Usage Details</h2>
                        <button @click="selectedTenant = null" class="text-gray-400 hover:text-gray-600">
                            <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button>
                    </div>
                    <div v-if="usageLoading" class="text-center py-8">
                        <p class="text-gray-500">Loading usage data...</p>
                    </div>
                    <div v-else-if="usageData" class="space-y-4">
                        <div class="rounded-xl bg-gray-50 p-4">
                            <p class="text-sm font-semibold text-gray-500">Event Count (30 days)</p>
                            <p class="mt-1 text-2xl font-mono font-bold text-gray-900">{{ usageData.event_count.toLocaleString() }}</p>
                        </div>
                        <div class="rounded-xl bg-gray-50 p-4">
                            <p class="text-sm font-semibold text-gray-500">Anomalies Detected</p>
                            <p class="mt-1 text-2xl font-mono font-bold text-gray-900">{{ usageData.anomaly_count.toLocaleString() }}</p>
                        </div>
                        <div class="rounded-xl bg-gray-50 p-4">
                            <p class="text-sm font-semibold text-gray-500">API Requests</p>
                            <p class="mt-1 text-2xl font-mono font-bold text-gray-900">{{ usageData.api_requests.toLocaleString() }}</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </main>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

interface Tenant {
    id: string
    name: string
    slug: string
    email: string
    plan: string
    status: string
    created_at: string
}

interface UsageData {
    tenant_id: string
    event_count: number
    anomaly_count: number
    api_requests: number
    period: string
}

const adminKey = ref('')
const authenticated = ref(false)
const authError = ref('')
const loading = ref(false)
const tenants = ref<Tenant[]>([])
const searchQuery = ref('')
const filterPlan = ref('')
const selectedTenant = ref<string | null>(null)
const usageData = ref<UsageData | null>(null)
const usageLoading = ref(false)

const verifiedCount = computed(() => tenants.value.filter(t => t.status === 'verified').length)
const trialCount = computed(() => tenants.value.filter(t => t.plan === 'trial').length)
const paidCount = computed(() => tenants.value.filter(t => ['starter', 'growth', 'enterprise'].includes(t.plan)).length)

const filteredTenants = computed(() => {
    return tenants.value.filter(t => {
        const matchesSearch = !searchQuery.value ||
            t.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
            t.email?.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
            t.slug.toLowerCase().includes(searchQuery.value.toLowerCase())
        const matchesPlan = !filterPlan.value || t.plan === filterPlan.value
        return matchesSearch && matchesPlan
    })
})

const authenticate = async () => {
    authError.value = ''
    try {
        const response = await fetch('/api/v1/admin/tenants', {
            headers: { 'X-Admin-Key': adminKey.value }
        })
        if (!response.ok) {
            throw new Error('Invalid admin key')
        }
        authenticated.value = true
        const data = await response.json()
        tenants.value = data.tenants || []
    } catch (error) {
        authError.value = error instanceof Error ? error.message : 'Authentication failed'
    }
}

const refreshData = async () => {
    if (!authenticated.value) return
    loading.value = true
    try {
        const response = await fetch('/api/v1/admin/tenants', {
            headers: { 'X-Admin-Key': adminKey.value }
        })
        if (response.ok) {
            const data = await response.json()
            tenants.value = data.tenants || []
        }
    } catch (error) {
        console.error('Failed to refresh data:', error)
    } finally {
        loading.value = false
    }
}

const viewTenantUsage = async (tenantId: string) => {
    selectedTenant.value = tenantId
    usageLoading.value = true
    usageData.value = null
    try {
        const response = await fetch(`/api/v1/admin/tenants/${tenantId}/usage`, {
            headers: { 'X-Admin-Key': adminKey.value }
        })
        if (response.ok) {
            usageData.value = await response.json()
        }
    } catch (error) {
        console.error('Failed to load usage:', error)
    } finally {
        usageLoading.value = false
    }
}

const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    })
}

const planBadgeClass = (plan: string) => {
    switch (plan) {
        case 'trial': return 'bg-gray-100 text-gray-800'
        case 'starter': return 'bg-blue-100 text-blue-800'
        case 'growth': return 'bg-purple-100 text-purple-800'
        case 'enterprise': return 'bg-amber-100 text-amber-800'
        default: return 'bg-gray-100 text-gray-800'
    }
}

const statusBadgeClass = (status: string) => {
    switch (status) {
        case 'verified': return 'bg-green-100 text-green-800'
        case 'pending': return 'bg-yellow-100 text-yellow-800'
        default: return 'bg-gray-100 text-gray-800'
    }
}

onMounted(() => {
    // Check for stored admin key
    const storedKey = localStorage.getItem('driftlock_admin_key')
    if (storedKey) {
        adminKey.value = storedKey
        authenticate()
    }
})
</script>
