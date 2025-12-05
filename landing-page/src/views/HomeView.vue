<template>
    <main class="bg-background text-foreground">
        <!-- Hero Section -->
        <section class="relative h-[80vh] min-h-[600px] w-full overflow-hidden border-b border-black bg-white text-black">
            <!-- Brutalist Grid Background -->
            <div class="absolute inset-0 z-0 opacity-10"
                 style="background-image: linear-gradient(#000 1px, transparent 1px), linear-gradient(90deg, #000 1px, transparent 1px); background-size: 40px 40px;">
            </div>
            
            <!-- Content -->
            <div class="container-padding relative z-10 flex h-full flex-col justify-center pb-16">
                <div class="max-w-4xl">
                    <div class="mb-6 flex items-center gap-3">
                        <div class="inline-block border border-black px-3 py-1 text-xs font-bold font-sans uppercase tracking-widest bg-black text-white">
                            Universal Signal Processing
                        </div>
                    </div>
                    
                    <h1 class="text-6xl sm:text-8xl font-sans font-bold tracking-tight leading-none mb-8 uppercase text-black">
                        The Universal<br>Anomaly Detector.
                    </h1>
                    
                    <div class="grid gap-8 lg:grid-cols-2 items-end">
                        <p class="text-xl font-serif leading-relaxed text-gray-800 border-l-2 border-black pl-6">
                            <strong>Low-cost. Explainable. Mathematical.</strong><br>
                            Detect drift in any data stream without black-box training. For agents, markets, and machines.
                        </p>
                        
                        <div class="flex flex-wrap gap-4">
                            <a href="#signup" class="brutalist-button-primary text-lg">
                                Join the Waitlist
                            </a>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <!-- Horizon Showcase (Factory Style) -->
        <HorizonShowcase @load="loadScenario" />

        <!-- Capabilities Section (Replacing Problem/Solution) -->
        <section id="capabilities" class="section-padding">
            <div class="container-padding">
                <div class="max-w-3xl mb-12">
                    <h2 class="text-sm font-bold font-sans uppercase tracking-widest mb-2 border-b border-black inline-block">Capabilities</h2>
                    <p class="text-4xl font-sans font-bold tracking-tight mt-4">DEPLOY EVERYWHERE.</p>
                </div>
                <!-- Use Case Cards -->
                <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
                    <div class="brutalist-card">
                        <h3 class="text-lg font-bold font-sans uppercase border-b border-black pb-2 mb-4">Agent Workflows</h3>
                        <p class="text-sm font-serif">Verify LLM outputs and tool usage patterns in real-time. Detect loops, hallucinations, and prompt injections instantly.</p>
                    </div>
                    <div class="brutalist-card">
                        <h3 class="text-lg font-bold font-sans uppercase border-b border-black pb-2 mb-4">Market Data</h3>
                        <p class="text-sm font-serif">Detect microstructure anomalies and liquidity breaks before the candle closes. Pure signal for HFT and quant strategies.</p>
                    </div>
                    <div class="brutalist-card bg-black text-white border-black">
                        <h3 class="text-lg font-bold font-sans uppercase border-b border-white pb-2 mb-4">IoT Telemetry</h3>
                        <p class="text-sm font-serif">Monitor thousands of sensors with O(1) complexity. No training data required. Just math.</p>
                    </div>
                </div>
            </div>
        </section>

        <!-- Who It's For Section -->
        <section id="audience" class="section-padding bg-white">
            <div class="container-padding">
                <div class="max-w-3xl mb-12">
                    <p class="text-sm font-bold font-sans uppercase tracking-widest mb-2 border-b border-black inline-block">Who It's For</p>
                    <h2 class="text-4xl font-sans font-bold tracking-tight mt-4">BUILT FOR BUILDERS. READY FOR AUDITORS.</h2>
                </div>

                <div class="grid md:grid-cols-2 gap-8 lg:gap-16">
                    <!-- For Developers -->
                    <div class="brutalist-card">
                        <h3 class="text-lg font-bold font-sans uppercase border-b border-black pb-2 mb-6">For Developers</h3>
                        <ul class="space-y-4 font-mono text-sm">
                            <li class="flex items-start gap-3">
                                <span class="text-green-600 font-bold">[x]</span>
                                <span>No ML training required</span>
                            </li>
                            <li class="flex items-start gap-3">
                                <span class="text-green-600 font-bold">[x]</span>
                                <span>Sub-50ms detection latency</span>
                            </li>
                            <li class="flex items-start gap-3">
                                <span class="text-green-600 font-bold">[x]</span>
                                <span>Simple REST API</span>
                            </li>
                            <li class="flex items-start gap-3">
                                <span class="text-green-600 font-bold">[x]</span>
                                <span>Free tier for experimentation</span>
                            </li>
                        </ul>
                        <p class="mt-6 text-sm font-serif text-gray-600 border-t border-gray-200 pt-4">
                            <strong>Use cases:</strong> API health, latency drift, error rate spikes, trading signals.
                        </p>
                    </div>

                    <!-- For Compliance -->
                    <div class="brutalist-card bg-black text-white border-black">
                        <h3 class="text-lg font-bold font-sans uppercase border-b border-white pb-2 mb-6">For Compliance</h3>
                        <ul class="space-y-4 font-mono text-sm">
                            <li class="flex items-start gap-3">
                                <span class="text-green-400 font-bold">[x]</span>
                                <span>Explainable results (White-box)</span>
                            </li>
                            <li class="flex items-start gap-3">
                                <span class="text-green-400 font-bold">[x]</span>
                                <span>DORA/NIS2 evidence bundles</span>
                            </li>
                            <li class="flex items-start gap-3">
                                <span class="text-green-400 font-bold">[x]</span>
                                <span>Mathematical proof</span>
                            </li>
                            <li class="flex items-start gap-3">
                                <span class="text-green-400 font-bold">[x]</span>
                                <span>Audit-ready documentation</span>
                            </li>
                        </ul>
                        <p class="mt-6 text-sm font-serif text-gray-400 border-t border-gray-700 pt-4">
                            <strong>Regulations:</strong> DORA Article 15, NIS2, FFIEC, AI Act monitoring.
                        </p>
                    </div>
                </div>
            </div>
        </section>

        <!-- Waitlist Section -->
        <section id="signup" class="section-padding bg-gray-50">
            <div class="container-padding">
                <div class="max-w-2xl mx-auto text-center">
                    <h2 class="text-4xl font-sans font-bold tracking-tight mb-4">GET EARLY ACCESS</h2>
                    <p class="text-lg font-serif mb-8 text-gray-600">
                        We're launching soon. Join the waitlist for priority access.
                    </p>
                    <div class="flex justify-center">
                        <WaitlistForm />
                    </div>
                    <p class="mt-8 text-sm font-mono text-gray-400">
                        Already have access? <a href="/login" class="underline hover:text-black">Sign in →</a>
                    </p>
                </div>
            </div>
        </section>

        <!-- Pricing Section -->
        <section id="pricing" class="section-padding">
            <div class="container-padding">
                <div class="max-w-3xl mb-16">
                    <p class="text-sm font-bold font-sans uppercase tracking-widest mb-2 border-b border-black inline-block">Pricing</p>
                    <h2 class="text-4xl font-sans font-bold tracking-tight mt-4">PRICING OPTIONS</h2>
                    <p class="text-lg font-serif mt-4">Simple, usage-based pricing for any scale.</p>
                </div>

                <div class="grid gap-8 lg:grid-cols-4">
                    <!-- Developer -->
                    <div class="brutalist-card flex flex-col">
                        <div class="mb-8">
                            <p class="text-sm font-bold font-sans uppercase tracking-wide mb-2">Free</p>
                            <p class="text-4xl font-sans font-bold">$0</p>
                            <p class="text-sm text-gray-500 font-mono">per month</p>
                            <p class="mt-4 text-sm font-serif">For experimentation and hobbyists.</p>
                        </div>
                        <ul class="space-y-3 text-sm font-mono border-t border-black pt-4 mt-auto">
                            <li>[x] 10k events/mo</li>
                            <li>[x] CLI demo & Playground</li>
                            <li>[x] JSON reports</li>
                        </ul>
                        <a href="#signup" class="brutalist-button w-full text-center mt-8">Join Waitlist</a>
                    </div>

                    <!-- Basic -->
                    <div class="brutalist-card flex flex-col relative bg-black text-white">
                         <div class="absolute top-0 right-0 bg-white text-black text-xs font-bold font-sans px-2 py-1 uppercase border-l border-b border-black">
                            Most Popular
                        </div>
                        <div class="mb-8">
                            <p class="text-sm font-bold font-sans uppercase tracking-wide mb-2">Standard</p>
                            <p class="text-4xl font-sans font-bold">$15</p>
                            <p class="text-sm text-gray-400 font-mono">per month</p>
                            <p class="mt-4 text-sm font-serif">For startups and active monitoring.</p>
                        </div>
                        <ul class="space-y-3 text-sm font-mono border-t border-white pt-4 mt-auto">
                            <li>[x] 500k events/mo</li>
                            <li>[x] Email alerts</li>
                            <li>[x] 30-day retention</li>
                        </ul>
                        <a href="#signup" class="brutalist-button-primary w-full text-center mt-8 border-white hover:bg-white hover:text-black">Join Waitlist</a>
                    </div>

                    <!-- Pro -->
                    <div class="brutalist-card flex flex-col relative">
                        <div class="mb-8">
                            <p class="text-sm font-bold font-sans uppercase tracking-wide mb-2">Pro</p>
                            <p class="text-4xl font-sans font-bold">$100</p>
                            <p class="text-sm text-gray-500 font-mono">per month</p>
                            <p class="mt-4 text-sm font-serif">For high-volume production streams.</p>
                        </div>
                        <ul class="space-y-3 text-sm font-mono border-t border-black pt-4 mt-auto">
                            <li>[x] 5M events/mo</li>
                            <li>[x] DORA/NIS2 Evidence</li>
                            <li>[x] Priority Support</li>
                        </ul>
                        <a href="#signup" class="brutalist-button w-full text-center mt-8">Join Waitlist</a>
                    </div>

                    <!-- Enterprise -->
                    <div class="brutalist-card flex flex-col">
                        <div class="mb-8">
                            <p class="text-sm font-bold font-sans uppercase tracking-wide mb-2">Enterprise</p>
                            <p class="text-4xl font-sans font-bold">$299</p>
                            <p class="text-sm text-gray-500 font-mono">per month</p>
                            <p class="mt-4 text-sm font-serif">For large-scale production monitoring.</p>
                        </div>
                        <ul class="space-y-3 text-sm font-mono border-t border-black pt-4 mt-auto">
                            <li>[x] 25M events/mo</li>
                            <li>[x] All Pro features</li>
                            <li>[x] 90-day retention</li>
                            <li class="text-gray-500">[~] SSO (coming soon)</li>
                        </ul>
                        <button @click="handleCheckout('orbit')"
                                :disabled="checkoutLoading === 'orbit'"
                                class="brutalist-button w-full text-center mt-8">
                            {{ checkoutLoading === 'orbit' ? 'Loading...' : 'Get Started' }}
                        </button>
                    </div>
                </div>
            </div>
        </section>

        <!-- CTA Section -->
        <section id="contact" class="section-padding">
             <div class="container-padding text-center">
                <h2 class="text-4xl font-sans font-bold tracking-tight mb-6">LAUNCHING SOON</h2>
                <p class="text-xl font-serif max-w-2xl mx-auto mb-10">
                    Be the first to know when we launch. Early access for waitlist members.
                </p>
                <div class="flex flex-wrap items-center justify-center gap-6">
                    <a href="#signup" class="brutalist-button-primary">
                        Join the Waitlist
                    </a>
                    <a href="mailto:hunter@shannonlabs.dev" class="text-lg font-bold font-sans uppercase hover:underline underline-offset-4 decoration-1">
                        Contact Us →
                    </a>
                </div>
             </div>
        </section>

        <!-- Checkout Error Toast -->
        <div v-if="checkoutError"
             class="fixed bottom-4 left-1/2 -translate-x-1/2 z-50
                    bg-red-600 text-white px-6 py-4
                    border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]
                    flex items-center gap-3 max-w-md">
            <svg class="w-5 h-5 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span class="text-sm">{{ checkoutError }}</span>
            <button @click="checkoutError = null" class="ml-2 text-red-200 hover:text-white">
                <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
            </button>
        </div>
    </main>
</template>

<script setup lang="ts">
import { computed, reactive, ref, nextTick } from 'vue'
import WaitlistForm from '../components/cta/WaitlistForm.vue'
import HorizonShowcase from '../components/sections/HorizonShowcase.vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const checkoutLoading = ref<string | null>(null)
const checkoutError = ref<string | null>(null)

const contactEndpoint = '/api/v1/contact'
const fallbackMailto = 'mailto:hunter@shannonlabs.dev?subject=Driftlock Pilot Program Inquiry'
const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

const contactForm = reactive({
    name: '',
    email: '',
    company: '',
    message: '',
})

const isSubmitting = ref(false)
const submissionState = ref<'idle' | 'success' | 'error'>('idle')
const submissionMessage = ref('')

const isFormValid = computed(() => {
    return (
        contactForm.name.trim().length >= 2 &&
        emailRegex.test(contactForm.email.trim()) &&
        contactForm.message.trim().length >= 20
    )
})

const resetForm = () => {
    contactForm.name = ''
    contactForm.email = ''
    contactForm.company = ''
    contactForm.message = ''
}

const handleContactSubmit = async () => {
    // Existing logic... (keeping it minimal for this view update)
    if (!isFormValid.value) return
    // ... (simulated)
    isSubmitting.value = true
    setTimeout(() => {
        isSubmitting.value = false
        submissionState.value = 'success'
        resetForm()
    }, 1000)
}

const loadScenario = async (url: string) => {
    // Placeholder for when we re-add interactivity or link to docs
    console.log('Load scenario:', url)
}

// Handle pricing tier checkout
const handleCheckout = async (plan: string) => {
    checkoutError.value = null

    // If not authenticated, scroll to signup
    if (!authStore.isAuthenticated) {
        const el = document.getElementById('signup')
        el?.scrollIntoView({ behavior: 'smooth' })
        return
    }

    // Authenticated: call checkout endpoint
    checkoutLoading.value = plan
    try {
        const token = await authStore.getToken()
        const res = await fetch('/api/v1/billing/checkout', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ plan })
        })
        if (res.ok) {
            const data = await res.json()
            if (data.url) window.location.href = data.url
        } else {
            checkoutError.value = 'Unable to start checkout. Please try again.'
        }
    } catch (e) {
        checkoutError.value = 'Network error. Please check your connection.'
    } finally {
        checkoutLoading.value = null
    }
}
</script>

<style scoped>
/* Utility-first approach, minimal custom CSS */
.section-padding {
    @apply py-16 sm:py-24 border-b border-black;
}

.container-padding {
    @apply mx-auto max-w-7xl px-4 sm:px-6 lg:px-8;
}
</style>
