<template>
  <div class="flex flex-col">
    <div class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
      <div class="inline-block min-w-full py-2 align-middle md:px-6 lg:px-8">
        <div class="overflow-hidden shadow ring-1 ring-black ring-opacity-5 md:rounded-lg">
          <table class="min-w-full divide-y divide-gray-300">
            <thead class="bg-gray-50">
              <tr>
                <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-6">Name</th>
                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Email</th>
                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Plan</th>
                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Status</th>
                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Created</th>
                <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-6">
                  <span class="sr-only">Actions</span>
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 bg-white">
              <tr v-for="tenant in tenants" :key="tenant.id">
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6">{{ tenant.name }}</td>
                <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{{ tenant.email }}</td>
                <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500 uppercase">{{ tenant.plan }}</td>
                <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                  <span :class="statusClass(tenant.status)" class="inline-flex rounded-full px-2 text-xs font-semibold leading-5">
                    {{ tenant.status }}
                  </span>
                </td>
                <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{{ new Date(tenant.created_at).toLocaleDateString() }}</td>
                <td class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-6">
                  <button @click="$emit('view', tenant)" class="text-indigo-600 hover:text-indigo-900">View</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  tenants: any[]
}>()

defineEmits(['view'])

const statusClass = (status: string) => {
  switch (status) {
    case 'active':
      return 'bg-green-100 text-green-800'
    case 'suspended':
      return 'bg-red-100 text-red-800'
    case 'pending_verification':
      return 'bg-yellow-100 text-yellow-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}
</script>
