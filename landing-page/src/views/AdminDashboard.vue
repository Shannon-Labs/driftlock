<template>
    <main class="min-h-screen bg-white">
        <!-- Header -->
        <header class="bg-white border-b-4 border-black">
            <div class="container mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                <div class="flex h-20 items-center justify-between">
                    <h1 class="text-2xl font-sans font-black uppercase tracking-tighter text-black">Driftlock Admin</h1>
                    <div class="flex items-center gap-4">
                        <span class="text-sm font-mono font-bold text-gray-600 uppercase tracking-wider">{{ tenants.length }} tenants</span>
                        <button
                            @click="refreshData"
                            :disabled="loading"
                            class="border-2 border-black bg-black px-4 py-2 text-sm font-bold uppercase tracking-widest text-white hover:bg-white hover:text-black transition-colors shadow-[4px_4px_0px_0px_rgba(0,0,0,0)] hover:shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                            {{ loading ? 'Loading...' : 'Refresh' }}
                        </button>
                    </div>
                </div>
            </div>
        </header>

        <!-- Auth Modal -->
        <div v-if="!authenticated" class="fixed inset-0 bg-black/80 flex items-center justify-center z-50 backdrop-blur-sm">
            <div class="bg-white border-4 border-black p-8 shadow-[8px_8px_0px_0px_rgba(255,255,255,1)] w-full max-w-md">
                <h2 class="text-2xl font-sans font-black uppercase tracking-tighter text-black mb-2">Admin Authentication</h2>
                <p class="text-sm font-serif text-gray-800 mb-6">Enter your admin key to access the dashboard.</p>
                <form @submit.prevent="authenticate" class="space-y-4">
                    <div>
                        <label for="admin-key" class="block text-xs font-bold uppercase tracking-widest text-black mb-1">Admin Key</label>
                        <input
                            id="admin-key"
                            v-model="adminKey"
                            type="password"
                            required
                            class="w-full border-2 border-black px-4 py-3 focus:outline-none focus:ring-4 focus:ring-black/20 font-mono"
                            placeholder="Enter admin key"
                        />
                    </div>
                    <button
                        type="submit"
                        class="w-full border-2 border-black bg-black px-4 py-3 font-bold uppercase tracking-widest text-white hover:bg-white hover:text-black transition-colors shadow-[4px_4px_0px_0px_rgba(0,0,0,0)] hover:shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]"
                    >
                        Sign In
                    </button>
                    <p v-if="authError" class="text-sm font-bold text-red-600 border border-red-600 bg-red-50 p-2">{{ authError }}</p>
                </form>
            </div>
        </div>

        <!-- Main Content -->
        <div v-if="authenticated" class="container mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
            <!-- Stats Cards -->
            <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-4 mb-8">
                <div class="bg-white border-2 border-black p-6 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]">
                    <p class="text-xs font-bold uppercase tracking-widest text-gray-500">Total Tenants</p>
                    <p class="mt-2 text-4xl font-mono font-bold text-black">{{ tenants.length }}</p>
                </div>
                <div class="bg-white border-2 border-black p-6 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]">
                    <p class="text-xs font-bold uppercase tracking-widest text-gray-500">Verified</p>
                    <p class="mt-2 text-4xl font-mono font-bold text-black">{{ verifiedCount }}</p>
                </div>
                <div class="bg-white border-2 border-black p-6 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]">
                    <p class="text-xs font-bold uppercase tracking-widest text-gray-500">Trial</p>
                    <p class="mt-2 text-4xl font-mono font-bold text-black">{{ trialCount }}</p>
                </div>
                <div class="bg-white border-2 border-black p-6 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]">
                    <p class="text-xs font-bold uppercase tracking-widest text-gray-500">Paid</p>
                    <p class="mt-2 text-4xl font-mono font-bold text-black">{{ paidCount }}</p>
                </div>
            </div>

            <!-- Search and Filter -->
            <div class="mb-8 flex flex-col sm:flex-row gap-4">
                <div class="flex-1">
                    <input
                        v-model="searchQuery"
                        type="text"
                        placeholder="SEARCH BY NAME, EMAIL, OR SLUG..."
                        class="w-full border-2 border-black px-4 py-3 focus:outline-none focus:ring-4 focus:ring-black/20 font-mono text-sm placeholder-gray-500 uppercase"
                    />
                </div>
                <select
                    v-model="filterPlan"
                    class="border-2 border-black px-4 py-3 focus:outline-none focus:ring-4 focus:ring-black/20 font-bold uppercase text-sm bg-white"
                >
                    <option value="">All Plans</option>
                    <option value="trial">Trial</option>
                    <option value="starter">Starter</option>
                    <option value="growth">Growth</option>
                    <option value="enterprise">Enterprise</option>
                </select>
            </div>

            <!-- Tenants Table -->
            <div class="bg-white border-2 border-black shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] overflow-hidden">
                <div class="overflow-x-auto">
                    <table class="w-full border-collapse">
                        <thead class="bg-black text-white">
                            <tr>
                                <th class="px-6 py-4 text-left text-xs font-bold uppercase tracking-widest border-b-2 border-black">
                                    Tenant
                                </th>
                                <th class="px-6 py-4 text-left text-xs font-bold uppercase tracking-widest border-b-2 border-black">
                                    Email
                                </th>
                                <th class="px-6 py-4 text-left text-xs font-bold uppercase tracking-widest border-b-2 border-black">
                                    Plan
                                </th>
                                <th class="px-6 py-4 text-left text-xs font-bold uppercase tracking-widest border-b-2 border-black">
                                    Status
                                </th>
                                <th class="px-6 py-4 text-left text-xs font-bold uppercase tracking-widest border-b-2 border-black">
                                    Created
                                </th>
                                <th class="px-6 py-4 text-left text-xs font-bold uppercase tracking-widest border-b-2 border-black">
                                    Actions
                                </th>
                            </tr>
                        </thead>
                        <tbody class="divide-y-2 divide-black">
                            <tr v-for="tenant in filteredTenants" :key="tenant.id" class="hover:bg-gray-50 transition-colors">
                                <td class="px-6 py-4 whitespace-nowrap border-r border-gray-200">
                                    <div>
                                        <p class="font-bold text-black uppercase">{{ tenant.name }}</p>
                                        <p class="text-xs text-gray-500 font-mono">{{ tenant.slug }}</p>
                                    </div>
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-800 border-r border-gray-200">
                                    {{ tenant.email || '-' }}
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap border-r border-gray-200">
                                    <span :class="planBadgeClass(tenant.plan)" class="px-2 py-1 text-xs font-bold uppercase tracking-wide border border-black">
                                        {{ tenant.plan }}
                                    </span>
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap border-r border-gray-200">
                                    <span :class="statusBadgeClass(tenant.status)" class="px-2 py-1 text-xs font-bold uppercase tracking-wide border border-black">
                                        {{ tenant.status }}
                                    </span>
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-600 border-r border-gray-200">
                                    {{ formatDate(tenant.created_at) }}
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap">
                                    <button
                                        @click="viewTenantUsage(tenant.id)"
                                        class="text-black hover:text-white hover:bg-black px-2 py-1 text-xs font-bold uppercase tracking-wider border border-black transition-colors"
                                    >
                                        View Usage
                                    </button>
                                </td>
                            </tr>
                            <tr v-if="filteredTenants.length === 0">
                                <td colspan="6" class="px-6 py-12 text-center text-gray-500 font-mono uppercase tracking-widest">
                                    No tenants found
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>

            <!-- Usage Modal -->
            <div v-if="selectedTenant" class="fixed inset-0 bg-black/80 flex items-center justify-center z-50 backdrop-blur-sm">
                <div class="bg-white border-4 border-black p-8 shadow-[8px_8px_0px_0px_rgba(255,255,255,1)] w-full max-w-lg">
                    <div class="flex items-center justify-between mb-6 border-b-2 border-black pb-4">
                        <h2 class="text-xl font-sans font-black uppercase tracking-tighter text-black">Usage Details</h2>
                        <button @click="selectedTenant = null" class="text-black hover:bg-black hover:text-white border-2 border-black p-1 transition-colors">
                            <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="square" stroke-linejoin="miter" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button>
                    </div>
                    <div v-if="usageLoading" class="text-center py-12">
                        <div class="animate-spin h-8 w-8 border-4 border-black border-t-transparent rounded-full mx-auto mb-4"></div>
                        <p class="text-gray-500 font-mono uppercase text-sm">Loading usage data...</p>
                    </div>
                    <div v-else-if="usageData" class="space-y-4">
                        <div class="border-2 border-black p-4 bg-gray-50">
                            <p class="text-xs font-bold uppercase tracking-widest text-gray-500">Event Count (30 days)</p>
                            <p class="mt-1 text-3xl font-mono font-bold text-black">{{ usageData.event_count.toLocaleString() }}</p>
                        </div>
                        <div class="border-2 border-black p-4 bg-gray-50">
                            <p class="text-xs font-bold uppercase tracking-widest text-gray-500">Anomalies Detected</p>
                            <p class="mt-1 text-3xl font-mono font-bold text-black">{{ usageData.anomaly_count.toLocaleString() }}</p>
                        </div>
                        <div class="border-2 border-black p-4 bg-gray-50">
                            <p class="text-xs font-bold uppercase tracking-widest text-gray-500">API Requests</p>
                            <p class="mt-1 text-3xl font-mono font-bold text-black">{{ usageData.api_requests.toLocaleString() }}</p>
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
        case 'trial': return 'bg-white text-black'
        case 'starter': return 'bg-gray-200 text-black'
        case 'growth': return 'bg-gray-800 text-white'
        case 'enterprise': return 'bg-black text-white'
        default: return 'bg-white text-black'
    }
}

const statusBadgeClass = (status: string) => {
    switch (status) {
        case 'verified': return 'bg-black text-white'
        case 'pending': return 'bg-white text-black border-dashed'
        default: return 'bg-white text-black'
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
