<template>
  <div class="min-h-screen bg-white">
    <!-- Checkout Success Toast -->
    <div v-if="showCheckoutSuccess"
         class="fixed top-4 right-4 z-50 bg-green-600 text-white px-6 py-4
                border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]
                flex items-center gap-3 max-w-md">
      <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
      </svg>
      <div>
        <p class="font-bold uppercase tracking-wider">Subscription Active!</p>
        <p class="text-sm">Your 14-day trial has begun.</p>
      </div>
      <button @click="showCheckoutSuccess = false" class="ml-2 text-green-200 hover:text-white">
        <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <!-- Toast Notifications -->
    <div v-if="toastMessage"
         :class="[
           'fixed top-4 right-4 z-50 px-6 py-4 border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] flex items-center gap-3 max-w-md',
           toastType === 'success' ? 'bg-green-600 text-white' : 'bg-red-600 text-white'
         ]">
      <svg v-if="toastType === 'success'" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
      </svg>
      <svg v-else class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <span class="font-bold uppercase tracking-wider">{{ toastMessage }}</span>
      <button @click="toastMessage = ''" class="ml-2 opacity-70 hover:opacity-100">
        <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <!-- Create Key Modal -->
    <div v-if="showCreateKeyModal" class="fixed inset-0 z-50 overflow-y-auto">
      <div class="flex min-h-full items-center justify-center p-4">
        <div class="fixed inset-0 bg-black/50" @click="showCreateKeyModal = false"></div>
        <div class="relative bg-white border-2 border-black shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] w-full max-w-md p-6">
          <h3 class="text-xl font-bold uppercase tracking-wide text-black border-b-2 border-black pb-2 mb-4">
            Create API Key
          </h3>
          <div class="space-y-4">
            <div>
              <label class="block text-xs font-bold uppercase tracking-widest text-gray-500 mb-2">
                Key Name
              </label>
              <input
                v-model="newKeyName"
                type="text"
                placeholder="e.g., Production API"
                class="w-full border-2 border-black px-4 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-black focus:ring-offset-2"
              />
            </div>
            <div>
              <label class="block text-xs font-bold uppercase tracking-widest text-gray-500 mb-2">
                Role
              </label>
              <select
                v-model="newKeyRole"
                class="w-full border-2 border-black px-4 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-black focus:ring-offset-2"
              >
                <option value="admin">Admin - Full access</option>
                <option value="stream">Stream - Detection only</option>
              </select>
            </div>
          </div>
          <div class="mt-6 flex gap-4">
            <button
              @click="showCreateKeyModal = false"
              class="flex-1 border-2 border-black bg-white px-4 py-2 text-sm font-bold uppercase tracking-widest text-black hover:bg-gray-100 transition-colors"
            >
              Cancel
            </button>
            <button
              @click="createKey"
              :disabled="createKeyLoading"
              class="flex-1 border-2 border-black bg-black px-4 py-2 text-sm font-bold uppercase tracking-widest text-white hover:bg-gray-800 transition-colors disabled:opacity-50"
            >
              {{ createKeyLoading ? 'Creating...' : 'Create Key' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- New Key Success Modal -->
    <div v-if="showNewKeyModal" class="fixed inset-0 z-50 overflow-y-auto">
      <div class="flex min-h-full items-center justify-center p-4">
        <div class="fixed inset-0 bg-black/50"></div>
        <div class="relative bg-white border-2 border-black shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] w-full max-w-lg p-6">
          <div class="flex items-center gap-2 border-b-2 border-black pb-2 mb-4">
            <svg class="w-6 h-6 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            <h3 class="text-xl font-bold uppercase tracking-wide text-black">
              Key Created!
            </h3>
          </div>
          <div class="bg-yellow-50 border-2 border-yellow-400 p-4 mb-4">
            <p class="text-sm font-bold text-yellow-800 uppercase tracking-wider mb-1">
              Copy your API key now
            </p>
            <p class="text-xs text-yellow-700">
              This is the only time you'll see the full key. Store it securely.
            </p>
          </div>
          <div class="bg-gray-900 p-4 font-mono text-sm text-green-400 break-all border-2 border-black">
            {{ newApiKey }}
          </div>
          <button
            @click="copyNewKey"
            class="mt-4 w-full border-2 border-black bg-black px-4 py-3 text-sm font-bold uppercase tracking-widest text-white hover:bg-gray-800 transition-colors flex items-center justify-center gap-2"
          >
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
            {{ keyCopied ? 'Copied!' : 'Copy to Clipboard' }}
          </button>
          <button
            @click="closeNewKeyModal"
            class="mt-2 w-full border-2 border-black bg-white px-4 py-2 text-sm font-bold uppercase tracking-widest text-black hover:bg-gray-100 transition-colors"
          >
            I've Saved My Key
          </button>
        </div>
      </div>
    </div>

    <!-- Revoke Key Confirmation Modal -->
    <div v-if="showRevokeModal" class="fixed inset-0 z-50 overflow-y-auto">
      <div class="flex min-h-full items-center justify-center p-4">
        <div class="fixed inset-0 bg-black/50" @click="showRevokeModal = false"></div>
        <div class="relative bg-white border-2 border-black shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] w-full max-w-md p-6">
          <div class="flex items-center gap-2 border-b-2 border-black pb-2 mb-4">
            <svg class="w-6 h-6 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <h3 class="text-xl font-bold uppercase tracking-wide text-black">
              Revoke Key?
            </h3>
          </div>
          <p class="text-sm text-gray-600 mb-2">
            Are you sure you want to revoke <span class="font-bold text-black">{{ keyToRevoke?.name || 'this key' }}</span>?
          </p>
          <p class="text-sm text-red-600 font-bold">
            This action cannot be undone. Any applications using this key will stop working immediately.
          </p>
          <div class="mt-6 flex gap-4">
            <button
              @click="showRevokeModal = false"
              class="flex-1 border-2 border-black bg-white px-4 py-2 text-sm font-bold uppercase tracking-widest text-black hover:bg-gray-100 transition-colors"
            >
              Cancel
            </button>
            <button
              @click="revokeKey"
              :disabled="revokeKeyLoading"
              class="flex-1 border-2 border-red-600 bg-red-600 px-4 py-2 text-sm font-bold uppercase tracking-widest text-white hover:bg-red-700 transition-colors disabled:opacity-50"
            >
              {{ revokeKeyLoading ? 'Revoking...' : 'Revoke Key' }}
            </button>
          </div>
        </div>
      </div>
    </div>

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

        <!-- Onboarding Wizard for new users -->
        <OnboardingWizard 
            v-if="!keys.length && !billingLoading" 
            :api-key="newApiKey || 'Generate a key below'"
            :api-url="apiUrl"
            @complete="showCreateKeyModal = false"
        />

        <!-- Billing Error -->
        <div v-if="billingError" class="mb-8 border-2 border-gray-300 bg-gray-50 px-6 py-4 flex items-center justify-between">
          <span class="text-sm text-gray-600">{{ billingError }}</span>
          <button @click="fetchBillingStatus" class="text-sm font-bold uppercase tracking-wider hover:underline">
            Retry
          </button>
        </div>

        <!-- Trial: Relaxed (8+ days) - Gray/subtle -->
        <div v-if="billing?.status === 'trialing' && trialUrgency === 'relaxed'"
             class="mb-8 border-2 border-gray-300 bg-gray-50 p-4">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-sm font-bold uppercase tracking-widest text-gray-500">TRIAL ACTIVE</span>
              <p class="text-sm font-mono text-gray-600 mt-1">
                {{ billing.trial_days_remaining }} days remaining to explore Driftlock.
              </p>
            </div>
            <button @click="manageBilling"
                    class="text-sm font-bold uppercase tracking-wider text-gray-600 hover:underline">
              View Plans
            </button>
          </div>
        </div>

        <!-- Trial: Reminder (4-7 days) - Yellow/warning -->
        <div v-else-if="billing?.status === 'trialing' && trialUrgency === 'reminder'"
             class="mb-8 border-2 border-black bg-yellow-100 p-4 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-sm font-bold uppercase tracking-widest text-yellow-800">{{ billing.trial_days_remaining }} DAYS LEFT</span>
              <p class="text-sm font-mono text-gray-800 mt-1">
                Add a payment method to continue using Driftlock after your trial.
              </p>
            </div>
            <button @click="manageBilling"
                    class="border-2 border-black bg-black px-4 py-2 text-sm font-bold uppercase text-white hover:bg-white hover:text-black transition-colors">
              Add Payment
            </button>
          </div>
        </div>

        <!-- Trial: Urgent (0-3 days) - Orange -->
        <div v-else-if="billing?.status === 'trialing' && trialUrgency === 'urgent'"
             class="mb-8 border-2 border-orange-500 bg-orange-100 p-4 shadow-[4px_4px_0px_0px_rgba(249,115,22,0.5)]">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-sm font-bold uppercase tracking-widest text-orange-800">
                {{ billing.trial_days_remaining === 0 ? 'TRIAL ENDS TODAY' : `TRIAL ENDS IN ${billing.trial_days_remaining} DAYS` }}
              </span>
              <p class="text-sm font-mono text-orange-900 mt-1">
                Subscribe now to avoid service interruption.
              </p>
            </div>
            <button @click="manageBilling"
                    class="border-2 border-orange-600 bg-orange-600 px-4 py-2 text-sm font-bold uppercase text-white hover:bg-white hover:text-orange-600 transition-colors">
              Subscribe Now
            </button>
          </div>
        </div>

        <!-- Grace Period Warning (Red) -->
        <div v-else-if="billing?.status === 'grace_period'"
             class="mb-8 border-2 border-red-600 bg-red-100 p-4 shadow-[4px_4px_0px_0px_rgba(220,38,38,0.5)]">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-sm font-bold uppercase tracking-widest text-red-800">PAYMENT FAILED</span>
              <p class="text-sm font-mono text-red-700 mt-1">
                Please update your payment method. Access will be restricted after grace period ends.
              </p>
            </div>
            <button @click="manageBilling"
                    class="border-2 border-red-800 bg-red-800 px-4 py-2 text-sm font-bold uppercase text-white hover:bg-white hover:text-red-800 transition-colors">
              Update Payment
            </button>
          </div>
        </div>

        <!-- Free Tier Upgrade Prompt (Gray) -->
	        <div v-else-if="billing?.status === 'free'"
	             class="mb-8 border-2 border-black bg-gray-100 p-4 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]">
	          <div class="flex items-center justify-between">
	            <div>
	              <span class="text-sm font-bold uppercase tracking-widest text-gray-600">FREE</span>
	              <p class="text-sm font-mono text-gray-600 mt-1">
	                Upgrade to unlock higher limits and longer retention.
	              </p>
	            </div>
	            <button @click="handleUpgrade('starter')"
	                    class="border-2 border-black bg-black px-4 py-2 text-sm font-bold uppercase text-white hover:bg-white hover:text-black transition-colors">
	              Upgrade to Starter
	            </button>
	          </div>
	        </div>

        <!-- Bento Grid Layout -->
        <div class="grid grid-cols-1 gap-8 lg:grid-cols-3">
          
          <!-- Usage Card with Chart -->
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

              <!-- Daily Usage Chart -->
              <div class="mt-8 border-t-2 border-black pt-4">
                <h4 class="text-sm font-bold uppercase tracking-widest text-gray-500 mb-4">Daily Activity (Last 30 Days)</h4>
                <div v-if="usageDetailsLoading" class="h-64 flex items-center justify-center">
                  <div class="text-sm font-mono text-gray-500">Loading chart...</div>
                </div>
                <UsageChart
                  v-else-if="usageDetails?.daily_usage?.length"
                  :data="usageDetails.daily_usage"
                />
                <div v-else class="h-64 flex items-center justify-center border-2 border-dashed border-gray-300">
                  <div class="text-center">
                    <p class="text-sm font-mono text-gray-500">No usage data yet</p>
                    <p class="text-xs text-gray-400 mt-1">Send your first events to see analytics</p>
                  </div>
                </div>
              </div>

              <div class="mt-8 border-t-2 border-black pt-4">
                 <dl class="grid grid-cols-1 gap-x-4 gap-y-4 sm:grid-cols-2">
                    <div class="sm:col-span-1">
                      <dt class="text-xs font-bold uppercase tracking-widest text-gray-500">Current Plan</dt>
                      <dd class="mt-1 text-2xl font-sans font-black text-black uppercase">{{ usage.plan || 'Free' }}</dd>
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
              <div class="mt-4 relative">
                <div class="bg-gray-900 p-4 font-mono text-xs text-green-400 overflow-x-auto border border-gray-700">
                  curl -X POST {{ apiUrl }}/v1/detect \<br>
                  &nbsp;&nbsp;-H "Authorization: Bearer {{ firstKey }}" \<br>
                  &nbsp;&nbsp;-d @events.json
                </div>
                <button
                  @click="copyCurlExample"
                  class="absolute top-2 right-2 p-1.5 bg-gray-800 hover:bg-gray-700 text-gray-400 hover:text-white rounded transition-colors"
                  title="Copy curl command"
                >
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                  </svg>
                </button>
              </div>
              <div class="mt-6">
                <router-link to="/docs/quickstart" class="text-sm font-bold uppercase tracking-widest text-white hover:underline decoration-2 underline-offset-4">
                  View full guide &rarr;
                </router-link>
              </div>
            </div>
          </div>

	          <!-- AI Usage Widget - show for paid plans -->
	          <AIUsageWidget
	            v-if="billing?.plan && !['free','pulse','trial','starter','basic'].includes(billing.plan)"
	            :plan="billing?.plan"
	            class="lg:col-span-1"
	            @upgrade="handleUpgrade"
	            @config-changed="onAIConfigChanged"
	          />

          <!-- API Keys Card -->
          <div class="bg-white border-2 border-black lg:col-span-3 shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
            <div class="px-6 py-6">
              <div class="sm:flex sm:items-center sm:justify-between border-b-2 border-black pb-4 mb-4">
                <h3 class="text-xl font-bold uppercase tracking-wide text-black">API Keys</h3>
                <div class="mt-4 sm:mt-0">
                  <button
                    type="button"
                    @click="showCreateKeyModal = true"
                    class="inline-flex items-center justify-center border-2 border-black bg-white px-4 py-2 text-sm font-bold uppercase tracking-widest text-black hover:bg-black hover:text-white transition-colors shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:shadow-none hover:translate-x-[2px] hover:translate-y-[2px]"
                  >
                    + Create Key
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
                              <div class="flex items-center gap-2">
                                <span class="bg-gray-100 px-2 py-1 border border-gray-300 rounded-none">{{ key.prefix || 'dlk_...' }}••••••••••••</span>
                                <span class="text-[10px] uppercase text-gray-400 font-bold tracking-wider">(Hidden)</span>
                              </div>
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
                              <button
                                @click="confirmRevokeKey(key)"
                                class="text-red-600 font-bold uppercase tracking-wider hover:text-red-900 hover:underline decoration-2 underline-offset-4"
                              >
                                Revoke
                              </button>
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

          <!-- Stream Breakdown -->
          <div v-if="usageDetails?.stream_breakdown?.length" class="bg-white border-2 border-black lg:col-span-3 shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
            <div class="px-6 py-6">
              <div class="border-b-2 border-black pb-4 mb-4">
                <h3 class="text-xl font-bold uppercase tracking-wide text-black">Usage by Stream</h3>
                <p class="text-sm font-serif text-gray-600 mt-1">Breakdown of events processed per stream this month</p>
              </div>
              <div class="overflow-x-auto">
                <table class="min-w-full divide-y-2 divide-black border-2 border-black">
                  <thead class="bg-black">
                    <tr>
                      <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-bold uppercase tracking-widest text-white sm:pl-6">Stream</th>
                      <th scope="col" class="px-3 py-3.5 text-right text-sm font-bold uppercase tracking-widest text-white">Events</th>
                      <th scope="col" class="px-3 py-3.5 text-right text-sm font-bold uppercase tracking-widest text-white">Requests</th>
                      <th scope="col" class="px-3 py-3.5 text-right text-sm font-bold uppercase tracking-widest text-white">Anomalies</th>
                      <th scope="col" class="px-3 py-3.5 text-right text-sm font-bold uppercase tracking-widest text-white sm:pr-6">Rate</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y-2 divide-black bg-white">
                    <tr v-for="stream in usageDetails.stream_breakdown" :key="stream.stream_id" class="hover:bg-gray-50">
                      <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-bold text-black sm:pl-6">
                        {{ stream.stream_name }}
                      </td>
                      <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-600 font-mono text-right">
                        {{ stream.event_count.toLocaleString() }}
                      </td>
                      <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-600 font-mono text-right">
                        {{ stream.request_count.toLocaleString() }}
                      </td>
                      <td class="whitespace-nowrap px-3 py-4 text-sm font-mono text-right">
                        <span :class="stream.anomaly_count > 0 ? 'text-red-600 font-bold' : 'text-gray-600'">
                          {{ stream.anomaly_count.toLocaleString() }}
                        </span>
                      </td>
                      <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-600 font-mono text-right sm:pr-6">
                        {{ stream.event_count > 0 ? ((stream.anomaly_count / stream.event_count) * 100).toFixed(2) : '0.00' }}%
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>

          <!-- Detection Settings -->
          <div class="bg-white border-2 border-black lg:col-span-3 shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
            <div class="px-6 py-5">
              <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                <!-- Learning Status Indicator -->
                <div class="flex items-center gap-3">
                  <div class="w-2 h-2 rounded-full bg-green-500 animate-pulse"></div>
                  <div>
                    <span class="text-sm font-bold uppercase tracking-wide text-black">Auto-Tuning Active</span>
                    <span v-if="feedbackCount > 0" class="ml-2 text-xs font-mono text-gray-600">
                      ({{ feedbackCount }} feedback{{ feedbackCount === 1 ? '' : 's' }} this session)
                    </span>
                  </div>
                </div>

                <!-- Profile Picker -->
                <div class="flex items-center gap-3">
                  <span class="text-xs font-bold uppercase tracking-widest text-gray-500">Sensitivity:</span>
                  <div class="flex border-2 border-black">
                    <button
                      @click="updateDetectionProfile('strict')"
                      :disabled="profileUpdateLoading"
                      :class="[
                        'px-3 py-1.5 text-xs font-bold uppercase tracking-wide transition-colors',
                        detectionProfile === 'strict'
                          ? 'bg-black text-white'
                          : 'bg-white text-black hover:bg-gray-100'
                      ]"
                    >
                      Low
                    </button>
                    <button
                      @click="updateDetectionProfile('balanced')"
                      :disabled="profileUpdateLoading"
                      :class="[
                        'px-3 py-1.5 text-xs font-bold uppercase tracking-wide border-l-2 border-r-2 border-black transition-colors',
                        detectionProfile === 'balanced'
                          ? 'bg-black text-white'
                          : 'bg-white text-black hover:bg-gray-100'
                      ]"
                    >
                      Med
                    </button>
                    <button
                      @click="updateDetectionProfile('sensitive')"
                      :disabled="profileUpdateLoading"
                      :class="[
                        'px-3 py-1.5 text-xs font-bold uppercase tracking-wide transition-colors',
                        detectionProfile === 'sensitive'
                          ? 'bg-black text-white'
                          : 'bg-white text-black hover:bg-gray-100'
                      ]"
                    >
                      High
                    </button>
                  </div>
                  <div v-if="profileUpdateLoading" class="w-4 h-4 border-2 border-black border-t-transparent rounded-full animate-spin"></div>
                </div>
              </div>

              <!-- Explanation text -->
              <p class="mt-3 text-xs text-gray-500 font-mono">
                {{ detectionProfile === 'sensitive' ? 'High: Catches more anomalies, may have more false positives' :
                   detectionProfile === 'strict' ? 'Low: Only flags clear anomalies, fewer false positives' :
                   'Medium: Balanced detection for most use cases' }}
              </p>
            </div>
          </div>

          <!-- Recent Anomalies -->
          <div class="bg-white border-2 border-black lg:col-span-3 shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
             <div class="px-6 py-5 border-b-2 border-black flex items-center justify-between">
               <h3 class="text-xl font-bold uppercase tracking-wide text-black">Recent Anomalies</h3>
               <router-link
                 v-if="recentAnomalies.length > 0"
                 to="/anomalies"
                 class="text-sm font-bold uppercase tracking-wider text-black hover:underline decoration-2 underline-offset-4"
               >
                 View All &rarr;
               </router-link>
             </div>

             <!-- Loading state -->
             <div v-if="anomaliesLoading" class="px-6 py-12 text-center">
               <div class="text-sm font-mono text-gray-500">Loading anomalies...</div>
             </div>

             <!-- Empty state -->
             <!-- Empty state with Action -->
             <div v-else-if="recentAnomalies.length === 0" class="px-6 py-12 flex flex-col items-center text-center">
                <div class="mb-4 text-4xl">📡</div>
                <h4 class="text-lg font-bold uppercase tracking-wide text-black mb-2">No Signal Detected</h4>
                <p class="text-sm font-serif text-gray-600 max-w-md mb-6">
                  Driftlock is listening, but hasn't heard anything yet. Send a test event to verify your connection.
                </p>
                <div class="w-full max-w-xl bg-black text-green-400 p-4 font-mono text-xs text-left relative group">
                  <div class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
                    <button @click="copyCurlExample" class="bg-white text-black px-2 py-1 text-[10px] font-bold uppercase">Copy</button>
                  </div>
                  curl -X POST {{ apiUrl }}/v1/detect \<br>
                  &nbsp;&nbsp;-H "Authorization: Bearer {{ firstKey }}" \<br>
                  &nbsp;&nbsp;-d @events.json
                </div>
             </div>

             <!-- Anomalies table -->
             <div v-else class="overflow-x-auto">
               <table class="min-w-full divide-y-2 divide-black">
                 <thead class="bg-gray-50">
                   <tr>
                     <th scope="col" class="py-3.5 pl-6 pr-3 text-left text-xs font-bold uppercase tracking-widest text-gray-500">Time</th>
                     <th scope="col" class="px-3 py-3.5 text-left text-xs font-bold uppercase tracking-widest text-gray-500">Stream</th>
                     <th scope="col" class="px-3 py-3.5 text-left text-xs font-bold uppercase tracking-widest text-gray-500">NCD Score</th>
                     <th scope="col" class="px-3 py-3.5 text-left text-xs font-bold uppercase tracking-widest text-gray-500">Confidence</th>
                     <th scope="col" class="px-3 py-3.5 text-left text-xs font-bold uppercase tracking-widest text-gray-500">Explanation</th>
                     <th scope="col" class="px-3 py-3.5 text-center text-xs font-bold uppercase tracking-widest text-gray-500">Feedback</th>
                   </tr>
                 </thead>
                 <tbody class="divide-y divide-gray-200 bg-white">
                   <tr v-for="anomaly in recentAnomalies" :key="anomaly.id" class="hover:bg-gray-50">
                     <td class="whitespace-nowrap py-4 pl-6 pr-3 text-sm font-mono text-gray-600">
                       {{ formatRelativeTime(anomaly.detected_at) }}
                     </td>
                     <td class="whitespace-nowrap px-3 py-4 text-sm font-bold text-black">
                       {{ anomaly.stream_name || 'default' }}
                     </td>
                     <td class="whitespace-nowrap px-3 py-4 text-sm">
                       <span :class="getNcdBadgeClass(anomaly.ncd)">
                         {{ anomaly.ncd?.toFixed(2) || 'N/A' }}
                       </span>
                     </td>
                     <td class="whitespace-nowrap px-3 py-4 text-sm font-mono text-gray-600">
                       {{ anomaly.confidence ? (anomaly.confidence * 100).toFixed(1) + '%' : 'N/A' }}
                     </td>
                     <td class="px-3 py-4 text-sm text-gray-600 max-w-xs truncate">
                       {{ anomaly.explanation || 'Anomaly detected' }}
                     </td>
                     <td class="whitespace-nowrap px-3 py-4 text-center">
                       <!-- Show feedback status if already given -->
                       <span v-if="anomaly.feedback_given === 'confirmed'" class="inline-flex items-center text-xs font-bold uppercase tracking-wide text-green-700">
                         <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/></svg>
                         Confirmed
                       </span>
                       <span v-else-if="anomaly.feedback_given === 'false_positive'" class="inline-flex items-center text-xs font-bold uppercase tracking-wide text-gray-500">
                         <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"/></svg>
                         False Pos
                       </span>
                       <!-- Show buttons if no feedback yet -->
                       <div v-else class="flex items-center justify-center gap-2">
                         <button
                           @click="submitFeedback(anomaly.id, 'confirmed')"
                           :disabled="feedbackLoading.has(anomaly.id)"
                           class="p-1.5 border-2 border-black bg-white hover:bg-green-100 disabled:opacity-50"
                           title="Confirm anomaly"
                         >
                           <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20"><path d="M2 10.5a1.5 1.5 0 113 0v6a1.5 1.5 0 01-3 0v-6zM6 10.333v5.43a2 2 0 001.106 1.79l.05.025A4 4 0 008.943 18h5.416a2 2 0 001.962-1.608l1.2-6A2 2 0 0015.56 8H12V4a2 2 0 00-2-2 1 1 0 00-1 1v.667a4 4 0 01-.8 2.4L6.8 7.933a4 4 0 00-.8 2.4z"/></svg>
                         </button>
                         <button
                           @click="submitFeedback(anomaly.id, 'false_positive')"
                           :disabled="feedbackLoading.has(anomaly.id)"
                           class="p-1.5 border-2 border-black bg-white hover:bg-red-100 disabled:opacity-50"
                           title="Mark as false positive"
                         >
                           <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20"><path d="M18 9.5a1.5 1.5 0 11-3 0v-6a1.5 1.5 0 013 0v6zM14 9.667v-5.43a2 2 0 00-1.106-1.79l-.05-.025A4 4 0 0011.057 2H5.64a2 2 0 00-1.962 1.608l-1.2 6A2 2 0 004.44 12H8v4a2 2 0 002 2 1 1 0 001-1v-.667a4 4 0 01.8-2.4l1.4-1.866a4 4 0 00.8-2.4z"/></svg>
                         </button>
                       </div>
                     </td>
                   </tr>
                 </tbody>
               </table>
             </div>
          </div>

        </div>
      </div>
    </DashboardLayout>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import DashboardLayout from '../layouts/DashboardLayout.vue'
import UsageChart from '../components/dashboard/UsageChart.vue'
import AIUsageWidget from '../components/dashboard/AIUsageWidget.vue'
import OnboardingWizard from '../components/dashboard/OnboardingWizard.vue'
import { useAuthStore } from '../stores/auth'
import type { APIKey } from '@/types/api'

const authStore = useAuthStore()
const route = useRoute()

const keys = ref<APIKey[]>([])
const usage = ref({
  current_period_usage: 0,
  plan_limit: 50000,
  plan: 'free'
})

// Detailed usage data
interface DailyUsage {
  date: string
  event_count: number
  request_count: number
  anomaly_count: number
}

interface StreamUsage {
  stream_id: string
  stream_name: string
  event_count: number
  request_count: number
  anomaly_count: number
}

interface Anomaly {
  id: string
  stream_id: string
  stream_name?: string
  ncd: number
  confidence: number
  explanation?: string
  detected_at: string
  feedback_given?: 'confirmed' | 'false_positive' | null // Track user feedback
}

const usageDetails = ref<{
  daily_usage: DailyUsage[]
  stream_breakdown: StreamUsage[]
  usage_percent: number
  period_start: string
  period_end: string
} | null>(null)
const usageDetailsLoading = ref(true)

// Recent anomalies state
const recentAnomalies = ref<Anomaly[]>([])
const anomaliesLoading = ref(true)

// Billing state
const billing = ref<{
  status: string
  plan: string
  trial_ends_at?: string
  trial_days_remaining?: number
  grace_period_ends_at?: string
  payment_failure_count: number
} | null>(null)
const billingLoading = ref(true)
const billingError = ref<string | null>(null)
const showCheckoutSuccess = ref(false)
const upgradeError = ref<string | null>(null)

// Detection settings state
const detectionProfile = ref<'sensitive' | 'balanced' | 'strict'>('balanced')
const profileUpdateLoading = ref(false)
const feedbackCount = ref(0) // Tracks total feedback given this session

// Toast notifications
const toastMessage = ref('')
const toastType = ref<'success' | 'error'>('success')

const showToast = (message: string, type: 'success' | 'error' = 'success') => {
  toastMessage.value = message
  toastType.value = type
  setTimeout(() => { toastMessage.value = '' }, 4000)
}

// Create Key Modal state
const showCreateKeyModal = ref(false)
const newKeyName = ref('')
const newKeyRole = ref('admin')
const createKeyLoading = ref(false)

// New Key Success Modal state
const showNewKeyModal = ref(false)
const newApiKey = ref('')
const keyCopied = ref(false)

// Revoke Key Modal state
const showRevokeModal = ref(false)
const keyToRevoke = ref<APIKey | null>(null)
const revokeKeyLoading = ref(false)

const apiUrl = window.location.origin // Assuming API is on same domain
const firstKey = computed(() => keys.value.length > 0 ? (keys.value[0].prefix || 'dlk_...') : 'YOUR_API_KEY')

const usagePercentage = computed(() => {
  if (usage.value.plan_limit === 0) return 0
  return Math.min(100, (usage.value.current_period_usage / usage.value.plan_limit) * 100)
})

// Trial urgency: escalates as trial ends
const trialUrgency = computed(() => {
  if (!billing.value || billing.value.status !== 'trialing') return null
  const days = billing.value.trial_days_remaining ?? 0
  if (days >= 8) return 'relaxed'      // Gray - subtle
  if (days >= 4) return 'reminder'     // Yellow - noticeable
  return 'urgent'                       // Orange - action needed
})

// Fetch billing status with error handling
const fetchBillingStatus = async () => {
  billingLoading.value = true
  billingError.value = null
  try {
    const token = await authStore.getToken()
    if (!token) return

    const res = await fetch('/api/v1/me/billing', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (res.ok) {
      billing.value = await res.json()
    } else {
      billingError.value = 'Unable to load billing status'
    }
  } catch (e) {
    billingError.value = 'Network error'
  } finally {
    billingLoading.value = false
  }
}

// Upgrade handler for free tier users
const handleUpgrade = async (plan?: string) => {
  // Default to next tier if no plan specified
  const targetPlan = plan || getNextTier(billing.value?.plan || 'free')
  upgradeError.value = null
  try {
    const token = await authStore.getToken()
    const res = await fetch('/api/v1/billing/checkout', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ plan: targetPlan })
    })
    if (res.ok) {
      const data = await res.json()
      if (data.url) window.location.href = data.url
    } else {
      upgradeError.value = 'Unable to start checkout'
    }
  } catch (e) {
    upgradeError.value = 'Network error. Please try again.'
  }
}

// Get next tier for upgrade suggestions
const getNextTier = (currentPlan: string): string => {
  const normalized = currentPlan.toLowerCase()
  const canonical =
    ['pulse', 'trial', 'pilot'].includes(normalized) ? 'free'
    : ['basic'].includes(normalized) ? 'starter'
    : ['radar', 'signal'].includes(normalized) ? 'pro'
    : ['tensor', 'growth', 'lock'].includes(normalized) ? 'team'
    : ['orbit', 'horizon'].includes(normalized) ? 'enterprise'
    : normalized

  const tiers = ['free', 'starter', 'pro', 'team', 'scale']
  const currentIndex = tiers.indexOf(canonical)
  if (currentIndex < 0) return 'starter'
  if (currentIndex >= tiers.length - 1) return 'scale'
  return tiers[currentIndex + 1]
}

// Handle AI config changes
const onAIConfigChanged = (config: { threshold: number; optimizeFor: string; maxCost: number }) => {
  showToast('AI configuration updated', 'success')
}

// Fetch keys from API
const fetchKeys = async () => {
  try {
    const token = await authStore.getToken()
    if (!token) return
    const res = await fetch('/api/v1/api-keys', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (res.ok) {
      const data = await res.json()
      keys.value = Array.isArray(data) ? data : (data.keys || [])
    }
  } catch (e) {
    if (import.meta.env.DEV) {
      console.error('Failed to fetch keys', e)
    }
  }
}

// Create new API key
const createKey = async () => {
  createKeyLoading.value = true
  try {
    const token = await authStore.getToken()
    const res = await fetch('/api/v1/api-keys', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        name: newKeyName.value || 'API Key',
        role: newKeyRole.value
      })
    })
    if (res.ok) {
      const data = await res.json()
      // Close create modal, show success modal with the key
      showCreateKeyModal.value = false
      newApiKey.value = data.key || data.api_key
      showNewKeyModal.value = true
      keyCopied.value = false
      // Reset form
      newKeyName.value = ''
      newKeyRole.value = 'admin'
      // Refresh keys list
      await fetchKeys()
    } else {
      const err = await res.json().catch(() => ({}))
      showToast(err.message || 'Failed to create key', 'error')
    }
  } catch (e) {
    showToast('Network error. Please try again.', 'error')
  } finally {
    createKeyLoading.value = false
  }
}

// Copy new key to clipboard
const copyNewKey = async () => {
  try {
    await navigator.clipboard.writeText(newApiKey.value)
    keyCopied.value = true
    setTimeout(() => { keyCopied.value = false }, 2000)
  } catch (e) {
    showToast('Failed to copy to clipboard', 'error')
  }
}

// Close new key modal
const closeNewKeyModal = () => {
  showNewKeyModal.value = false
  newApiKey.value = ''
}

// Generic copy to clipboard
const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    showToast('Copied to clipboard', 'success')
  } catch (e) {
    showToast('Failed to copy', 'error')
  }
}

// Copy curl example
const copyCurlExample = async () => {
  const curlCommand = `curl -X POST ${apiUrl}/v1/detect \\
  -H "Authorization: Bearer ${firstKey.value}" \\
  -d @events.json`
  await copyToClipboard(curlCommand)
}

// Open revoke confirmation
const confirmRevokeKey = (key: APIKey) => {
  keyToRevoke.value = key
  showRevokeModal.value = true
}

// Revoke API key
const revokeKey = async () => {
  if (!keyToRevoke.value) return
  revokeKeyLoading.value = true
  try {
    const token = await authStore.getToken()
    const res = await fetch(`/api/v1/api-keys/${keyToRevoke.value.id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    if (res.ok) {
      showRevokeModal.value = false
      keyToRevoke.value = null
      showToast('API key revoked', 'success')
      // Refresh keys list
      await fetchKeys()
    } else {
      const err = await res.json().catch(() => ({}))
      showToast(err.message || 'Failed to revoke key', 'error')
    }
  } catch (e) {
    showToast('Network error. Please try again.', 'error')
  } finally {
    revokeKeyLoading.value = false
  }
}

// Fetch recent anomalies
const fetchRecentAnomalies = async () => {
  anomaliesLoading.value = true
  try {
    // Use the first API key to fetch anomalies
    if (keys.value.length === 0) {
      anomaliesLoading.value = false
      return
    }
    const token = await authStore.getToken()
    const res = await fetch('/api/v1/anomalies?limit=5', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (res.ok) {
      const data = await res.json()
      recentAnomalies.value = data.anomalies || []
    }
  } catch (e) {
    if (import.meta.env.DEV) {
      console.error('Failed to fetch anomalies', e)
    }
  } finally {
    anomaliesLoading.value = false
  }
}

// Feedback loading state per anomaly
const feedbackLoading = ref<Set<string>>(new Set())

// Submit feedback for an anomaly (thumbs up = confirmed, thumbs down = false_positive)
const submitFeedback = async (anomalyId: string, feedbackType: 'confirmed' | 'false_positive') => {
  if (feedbackLoading.value.has(anomalyId)) return

  feedbackLoading.value.add(anomalyId)
  try {
    const token = await authStore.getToken()
    const res = await fetch(`/api/v1/anomalies/${anomalyId}/feedback`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ feedback_type: feedbackType })
    })
    if (res.ok) {
      // Update the local anomaly to show feedback was given
      const anomaly = recentAnomalies.value.find(a => a.id === anomalyId)
      if (anomaly) {
        anomaly.feedback_given = feedbackType
      }
      // Increment session feedback count for learning indicator
      feedbackCount.value++
      showToast(feedbackType === 'confirmed' ? 'Anomaly confirmed' : 'Marked as false positive', 'success')
    } else {
      const err = await res.json().catch(() => ({}))
      showToast(err.message || 'Failed to submit feedback', 'error')
    }
  } catch (e) {
    showToast('Network error. Please try again.', 'error')
  } finally {
    feedbackLoading.value.delete(anomalyId)
  }
}

// Update detection profile for all user's streams
const updateDetectionProfile = async (profile: 'sensitive' | 'balanced' | 'strict') => {
  if (profileUpdateLoading.value) return

  profileUpdateLoading.value = true
  try {
    const token = await authStore.getToken()

    // Get all streams from usageDetails
    const streams = usageDetails.value?.stream_breakdown || []
    if (streams.length === 0) {
      showToast('No streams found', 'error')
      return
    }

    // Update each stream's profile
    const updatePromises = streams.map(stream =>
      fetch(`/api/v1/streams/${stream.stream_id}/profile`, {
        method: 'PATCH',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ profile })
      })
    )

    const results = await Promise.all(updatePromises)
    const allSucceeded = results.every(res => res.ok)

    if (allSucceeded) {
      detectionProfile.value = profile
      const profileLabel = profile === 'sensitive' ? 'High' : profile === 'strict' ? 'Low' : 'Medium'
      showToast(`Sensitivity set to ${profileLabel}`, 'success')
    } else {
      showToast('Failed to update some streams', 'error')
    }
  } catch (e) {
    showToast('Network error. Please try again.', 'error')
  } finally {
    profileUpdateLoading.value = false
  }
}

// Format relative time (e.g., "2 hours ago")
const formatRelativeTime = (dateStr: string) => {
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)

  if (diffMins < 1) return 'just now'
  if (diffMins < 60) return `${diffMins}m ago`
  if (diffHours < 24) return `${diffHours}h ago`
  if (diffDays < 7) return `${diffDays}d ago`
  return date.toLocaleDateString()
}

// Get NCD badge class based on severity
const getNcdBadgeClass = (ncd: number) => {
  if (ncd >= 0.7) {
    return 'inline-flex items-center border border-red-600 bg-red-100 px-2 py-0.5 text-xs font-bold uppercase tracking-wide text-red-800'
  } else if (ncd >= 0.5) {
    return 'inline-flex items-center border border-orange-500 bg-orange-100 px-2 py-0.5 text-xs font-bold uppercase tracking-wide text-orange-800'
  } else if (ncd >= 0.3) {
    return 'inline-flex items-center border border-yellow-500 bg-yellow-100 px-2 py-0.5 text-xs font-bold uppercase tracking-wide text-yellow-800'
  }
  return 'inline-flex items-center border border-gray-300 bg-gray-100 px-2 py-0.5 text-xs font-bold uppercase tracking-wide text-gray-600'
}

onMounted(async () => {
  // Handle Stripe redirect success
  if (route.query.success === 'true') {
    showCheckoutSuccess.value = true
    window.history.replaceState({}, '', '/dashboard')
    setTimeout(() => { showCheckoutSuccess.value = false }, 5000)
  }

  // Clear canceled param
  if (route.query.canceled === 'true') {
    window.history.replaceState({}, '', '/dashboard')
  }

  try {
    const token = await authStore.getToken()
    if (!token) return

    // Fetch Billing Status
    await fetchBillingStatus()

    // Fetch Keys
    const resKeys = await fetch('/api/v1/api-keys', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (resKeys.ok) {
      const data = await resKeys.json()
      keys.value = Array.isArray(data) ? data : (data.keys || [])
    }

    // Fetch Usage
    const resUsage = await fetch('/api/v1/account/usage', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (resUsage.ok) {
      const data = await resUsage.json()
      usage.value = data
    }

    // Fetch Usage Details (chart data)
    usageDetailsLoading.value = true
    try {
      const resDetails = await fetch('/api/v1/me/usage/details', {
        headers: { 'Authorization': `Bearer ${token}` }
      })
      if (resDetails.ok) {
        usageDetails.value = await resDetails.json()
      }
    } finally {
      usageDetailsLoading.value = false
    }

    // Fetch Recent Anomalies
    await fetchRecentAnomalies()

  } catch (e) {
    if (import.meta.env.DEV) {
      console.error('Failed to fetch dashboard data', e)
    }
  }
})

const manageBilling = async () => {
  try {
    const token = await authStore.getToken()
    const res = await fetch('/api/v1/billing/portal', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (res.ok) {
      const data = await res.json()
      if (data.url) window.location.href = data.url
    } else {
      showToast('Unable to open billing portal', 'error')
    }
  } catch (e) {
    if (import.meta.env.DEV) console.error('Failed to open billing portal', e)
    showToast('Unable to open billing portal', 'error')
  }
}
</script>
