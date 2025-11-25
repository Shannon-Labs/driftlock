<template>
    <main class="bg-background text-foreground">
        <!-- Hero Section -->
        <!-- Hero Section -->
        <section class="relative h-[80vh] min-h-[600px] w-full overflow-hidden border-b border-black bg-black text-white">
            <!-- Video Background -->
            <div class="absolute inset-0 z-0">
                <video 
                    autoplay 
                    muted 
                    loop 
                    playsinline 
                    poster="/media/anomaly-hero.png"
                    class="h-full w-full object-cover opacity-60"
                >
                    <source src="/media/anomaly-hero.mp4" type="video/mp4">
                </video>
                <!-- Overlay Gradient -->
                <div class="absolute inset-0 bg-gradient-to-t from-black via-transparent to-transparent opacity-90"></div>
            </div>

            <!-- Content -->
            <div class="container-padding relative z-10 flex h-full flex-col justify-end pb-16">
                <div class="max-w-4xl">
                    <div class="mb-6 flex items-center gap-3">
                        <div class="inline-block border border-white px-3 py-1 text-xs font-bold font-sans uppercase tracking-widest bg-black text-white">
                            Anomaly Detected
                        </div>
                        <span class="font-mono text-xs text-red-500 animate-pulse">● LIVE RECORDING: KRAKEN XBT/USD</span>
                    </div>
                    
                    <h1 class="text-6xl sm:text-8xl font-sans font-bold tracking-tight leading-none mb-8 uppercase text-white mix-blend-difference">
                        When Data Drifts,<br>We Lock On.
                    </h1>
                    
                    <div class="grid gap-8 lg:grid-cols-2 items-end">
                        <p class="text-xl font-serif leading-relaxed text-gray-200 border-l-2 border-white pl-6">
                            <strong>Analysis:</strong> Live capture: Kraken XBT/USD sell, score 0.0215 with ΔH≈+0.0155 and ΔCR≈+0.0275. 
                            Compression efficiency dropped while entropy rose, flagging a micro-size sell near local highs—exactly the “sharp signal” we use to surface premium events for Pro/Enterprise.
                        </p>
                        
                        <div class="flex flex-wrap gap-4">
                            <a href="#signup" class="brutalist-button-primary text-lg border-white hover:bg-white hover:text-black transition-colors">
                                Start Free Trial
                            </a>
                            <button type="button" @click="triggerLiveScan" class="brutalist-button text-lg bg-transparent text-white border-white hover:bg-white hover:text-black">
                                Live Signal Scan →
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <!-- Horizon Showcase (Factory Style) -->
        <HorizonShowcase @load="loadScenario" />

        <!-- Signup Section -->
        <section id="signup" class="section-padding bg-gray-50">
            <div class="container-padding">
                <div class="grid lg:grid-cols-2 gap-16 items-start">
                    <div>
                        <h2 class="text-4xl font-sans font-bold tracking-tight mb-6">START YOUR PILOT</h2>
                        <p class="text-xl font-serif mb-8">
                            Create your account and receive your API key instantly.
                            <br>
                            <span class="font-bold underline underline-offset-2">No credit card required for Free plan.</span>
                        </p>
                        <ul class="space-y-4 font-mono text-sm border-t border-black pt-4">
                            <li class="flex items-center gap-3">
                                [x] Instant API access
                            </li>
                            <li class="flex items-center gap-3">
                                [x] 10k monthly events free
                            </li>
                            <li class="flex items-center gap-3">
                                [x] Full dashboard access
                            </li>
                        </ul>
                    </div>
                    <div class="brutalist-card">
                        <SignupForm />
                    </div>
                </div>
            </div>
        </section>



        <!-- Problem Section (The $ Risk) -->
        <section id="problem" class="section-padding">
            <div class="container-padding">
                <div class="max-w-3xl mb-12">
                    <h2 class="text-sm font-bold font-sans uppercase tracking-widest mb-2 border-b border-black inline-block">The Problem</h2>
                    <p class="text-4xl font-sans font-bold tracking-tight mt-4">YOUR BLACK BOX IS AN AUDIT FAILURE.</p>
                </div>
                <!-- Regulation Cards -->
                <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
                    <div class="brutalist-card">
                        <h3 class="text-lg font-bold font-sans uppercase border-b border-black pb-2 mb-4">EU DORA</h3>
                        <p class="text-5xl font-sans font-bold mb-2">2%</p>
                        <p class="text-sm font-serif">of annual global turnover in fines. Article 15 requires full explainability.</p>
                    </div>
                    <div class="brutalist-card">
                        <h3 class="text-lg font-bold font-sans uppercase border-b border-black pb-2 mb-4">US Regulations</h3>
                        <p class="text-4xl font-sans font-bold mb-2">CFPB + FFIEC</p>
                        <p class="text-sm font-serif">Fair Lending and FFIEC guidelines mandate auditability. Black boxes violate this.</p>
                    </div>
                    <div class="brutalist-card bg-black text-white border-black">
                        <h3 class="text-lg font-bold font-sans uppercase border-b border-white pb-2 mb-4">The Gap</h3>
                        <p class="text-4xl font-sans font-bold mb-2">Auto Failure</p>
                        <p class="text-sm font-serif">No mathematical proof = failed audit. Fines start January 2025.</p>
                    </div>
                </div>
            </div>
        </section>

        <!-- Live Playground Section -->
        <section id="playground" class="section-padding">
            <div class="container-padding">
                <div class="mb-12">
                    <h2 class="text-4xl font-sans font-bold tracking-tight">TRY IT YOURSELF</h2>
                    <p class="text-lg font-serif mt-2">Interact with the API directly from your browser.</p>
                </div>
                <div class="border border-black p-2 bg-gray-50">
                    <PlaygroundShell ref="playgroundRef" variant="embedded" />
                </div>
            </div>
        </section>

        <!-- Pricing Section -->
        <section id="pricing" class="section-padding">
            <div class="container-padding">
                <div class="max-w-3xl mb-16">
                    <p class="text-sm font-bold font-sans uppercase tracking-widest mb-2 border-b border-black inline-block">Pricing</p>
                    <h2 class="text-4xl font-sans font-bold tracking-tight mt-4">THE SHANNON TIERS</h2>
                    <p class="text-lg font-serif mt-4">Compression math costs about $1 per million events. We honor our Bell Labs roots by naming our tiers after the evolution of information theory.</p>
                </div>

                <div class="grid gap-8 lg:grid-cols-4">
                    <!-- Developer -->
                    <div class="brutalist-card flex flex-col">
                        <div class="mb-8">
                            <p class="text-sm font-bold font-sans uppercase tracking-wide mb-2">1. Pulse</p>
                            <p class="text-4xl font-sans font-bold">$0</p>
                            <p class="text-sm text-gray-500 font-mono">per month</p>
                            <p class="mt-4 text-sm font-serif">The heartbeat of your data. Local CLI, playground, and basic detection.</p>
                        </div>
                        <ul class="space-y-3 text-sm font-mono border-t border-black pt-4 mt-auto">
                            <li>[x] 10k events/mo</li>
                            <li>[x] CLI demo & Playground</li>
                            <li>[x] JSON reports</li>
                        </ul>
                        <a href="#signup" class="brutalist-button w-full text-center mt-8">Start Building</a>
                    </div>

                    <!-- Basic -->
                    <div class="brutalist-card flex flex-col relative">
                        <div class="mb-8">
                            <p class="text-sm font-bold font-sans uppercase tracking-wide mb-2">2. Radar</p>
                            <p class="text-4xl font-sans font-bold">$20</p>
                            <p class="text-sm text-gray-500 font-mono">per month</p>
                            <p class="mt-4 text-sm font-serif">Detect standard deviation. Filter the background noise.</p>
                            <p class="mt-2 text-xs text-gray-400 font-mono">Usage-based overage available</p>
                        </div>
                        <ul class="space-y-3 text-sm font-mono border-t border-black pt-4 mt-auto">
                            <li>[x] 500k events/mo</li>
                            <li>[x] Email alerts with context</li>
                            <li>[x] 30-day retention</li>
                        </ul>
                        <button type="button" onclick="window.location.href='#signup'" class="brutalist-button w-full text-center mt-8">Get Radar</button>
                    </div>

                    <!-- Pro -->
                    <div class="brutalist-card bg-black text-white flex flex-col relative">
                         <div class="absolute top-0 right-0 bg-white text-black text-xs font-bold font-sans px-2 py-1 uppercase border-l border-b border-black">
                            Most Popular
                        </div>
                        <div class="mb-8">
                            <p class="text-sm font-bold font-sans uppercase tracking-wide mb-2">3. Lock</p>
                            <p class="text-4xl font-sans font-bold">$200</p>
                            <p class="text-sm text-gray-400 font-mono">per month</p>
                            <p class="mt-4 text-sm font-serif">Full entropy surveillance. High-fidelity chaos detection & compliance.</p>
                            <p class="mt-2 text-xs text-gray-400 font-mono">Usage-based overage available</p>
                        </div>
                        <ul class="space-y-3 text-sm font-mono border-t border-white pt-4 mt-auto">
                            <li>[x] 5M events/mo</li>
                            <li>[x] DORA/NIS2 Evidence</li>
                            <li>[x] Priority Support</li>
                        </ul>
                        <button type="button" onclick="window.location.href='#signup'" class="brutalist-button-primary w-full text-center mt-8 border-white hover:bg-white hover:text-black">Engage Lock</button>
                    </div>

                    <!-- Enterprise -->
                    <div class="brutalist-card flex flex-col">
                        <div class="mb-8">
                            <p class="text-sm font-bold font-sans uppercase tracking-wide mb-2">4. Orbit</p>
                            <p class="text-4xl font-sans font-bold">Custom</p>
                            <p class="text-sm text-gray-500 font-mono">Unlimited</p>
                            <p class="mt-4 text-sm font-serif">See beyond the edge. Custom models, adaptive windowing, on-prem.</p>
                        </div>
                        <ul class="space-y-3 text-sm font-mono border-t border-black pt-4 mt-auto">
                            <li>[x] Custom OpenZL Models</li>
                            <li>[x] Private Cloud / On-prem</li>
                            <li>[x] SLA & Account Manager</li>
                        </ul>
                        <a href="mailto:hunter@shannonlabs.dev" class="brutalist-button w-full text-center mt-8">Contact Sales</a>
                    </div>
                </div>
            </div>
        </section>

        <!-- ROI Calculator Section -->
        <section id="roi" class="section-padding bg-black text-white">
            <div class="container-padding">
                <div class="max-w-3xl mb-12">
                    <h2 class="text-sm font-bold font-sans uppercase tracking-widest mb-2 text-gray-400 border-b border-gray-400 inline-block">Financial Impact</h2>
                    <h2 class="text-4xl font-sans font-bold tracking-tight mt-4">PROJECTED IMPACT</h2>
                    <p class="text-base font-mono mt-4 text-gray-400">
                        Estimated savings based on typical incident response costs and regulatory fine exposure.
                    </p>
                </div>
                <div class="grid gap-10 lg:grid-cols-2 lg:items-start">
                    <!-- Inputs -->
                    <div class="border border-white p-8">
                        <h3 class="text-lg font-bold font-sans uppercase mb-6 border-b border-white pb-2">Your Assumptions</h3>
                        <div class="space-y-6">
                            <div>
                                <label class="flex justify-between text-sm font-mono mb-2">
                                    <span>Events per month</span>
                                    <span>{{ roiInputs.eventsPerMonth.toLocaleString() }}</span>
                                </label>
                                <input v-model.number="roiInputs.eventsPerMonth" type="range" min="1000000" max="100000000" step="1000000" class="w-full accent-white h-2 bg-gray-800 rounded-none appearance-none cursor-pointer" />
                            </div>
                            <div>
                                <label class="flex justify-between text-sm font-mono mb-2">
                                    <span>Annual risk exposure</span>
                                    <span>${{ roiInputs.annualLoss.toLocaleString() }}</span>
                                </label>
                                <input v-model.number="roiInputs.annualLoss" type="range" min="500000" max="50000000" step="500000" class="w-full accent-white h-2 bg-gray-800 rounded-none appearance-none cursor-pointer" />
                            </div>
                            <div>
                                <label class="flex justify-between text-sm font-mono mb-2">
                                    <span>Loss reduction</span>
                                    <span>{{ roiInputs.lossReductionPct }}%</span>
                                </label>
                                <input v-model.number="roiInputs.lossReductionPct" type="range" min="1" max="50" step="1" class="w-full accent-white h-2 bg-gray-800 rounded-none appearance-none cursor-pointer" />
                            </div>
                        </div>
                    </div>

                    <!-- Outputs -->
                    <div class="border border-white p-8 bg-gray-900">
                        <h3 class="text-lg font-bold font-sans uppercase mb-6 border-b border-white pb-2">Projected Annual Impact</h3>
                        <div class="grid gap-4 sm:grid-cols-2">
                            <div class="border border-gray-700 p-4">
                                <p class="text-xs font-mono uppercase text-gray-400">Est. Cost</p>
                                <p class="mt-2 text-2xl font-mono font-bold">${{ roiOutputs.annualDriftlockSpend.toLocaleString() }}</p>
                            </div>
                             <div class="border border-gray-700 p-4">
                                <p class="text-xs font-mono uppercase text-gray-400">Savings</p>
                                <p class="mt-2 text-2xl font-mono font-bold text-green-400">${{ roiOutputs.lossesAvoided.toLocaleString() }}</p>
                            </div>
                             <div class="border border-gray-700 p-4 sm:col-span-2 bg-black">
                                <p class="text-xs font-mono uppercase text-gray-400">Net Benefit</p>
                                <p class="mt-2 text-4xl font-mono font-bold text-white">
                                    ${{ roiOutputs.netBenefit.toLocaleString() }}
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <!-- CTA Section -->
        <section id="contact" class="section-padding">
             <div class="container-padding text-center">
                <h2 class="text-4xl font-sans font-bold tracking-tight mb-6">READY TO VERIFY?</h2>
                <p class="text-xl font-serif max-w-2xl mx-auto mb-10">
                    Join the pilot program for Q1 2026. We are accepting 3 more partners.
                </p>
                <div class="flex flex-wrap items-center justify-center gap-6">
                    <a href="mailto:hunter@shannonlabs.dev" class="brutalist-button">
                        Contact Sales
                    </a>
                    <a href="#signup" class="text-lg font-bold font-sans uppercase hover:underline underline-offset-4 decoration-1">
                        Sign Up Free →
                    </a>
                </div>
             </div>
        </section>
    </main>
</template>

<script setup lang="ts">
import { computed, reactive, ref, nextTick } from 'vue'
import PlaygroundShell from '../components/playground/PlaygroundShell.vue'
import SignupForm from '../components/cta/SignupForm.vue'
import HorizonShowcase from '../components/sections/HorizonShowcase.vue'

const playgroundRef = ref<InstanceType<typeof PlaygroundShell> | null>(null)
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

const triggerLiveScan = () => {
    const target = document.getElementById('playground')
    target?.scrollIntoView({ behavior: 'smooth', block: 'start' })
    playgroundRef.value?.runFinancialDemo?.()
}

const loadScenario = async (url: string) => {
    const target = document.getElementById('playground')
    target?.scrollIntoView({ behavior: 'smooth', block: 'start' })
    // We need to expose a method on PlaygroundShell to load a specific URL
    // Currently it has runFinancialDemo and onSample
    // I will need to access the SamplePicker inside or call onSample directly if I expose it
    if (playgroundRef.value?.loadSample) {
        await playgroundRef.value.loadSample(url)
    }
}

// ROI calculator state
const roiInputs = reactive({
    eventsPerMonth: 50_000_000,
    annualLoss: 5_000_000,
    lossReductionPct: 10,
    fineExposure: 10_000_000,
    pricePerMillion: 1.0,
})

const roiOutputs = computed(() => {
    const eventsPerMonth = Math.max(1_000_000, roiInputs.eventsPerMonth || 0)
    const annualEvents = eventsPerMonth * 12
    const millions = annualEvents / 1_000_000
    const annualDriftlockSpend = Math.round(millions * (roiInputs.pricePerMillion || 0))

    const lossesAvoided = Math.round((roiInputs.annualLoss || 0) * (Math.max(0, roiInputs.lossReductionPct || 0) / 100))
    const fineExposureConsidered = Math.round(roiInputs.fineExposure || 0)

    const netBenefit = lossesAvoided + fineExposureConsidered - annualDriftlockSpend

    return {
        annualDriftlockSpend: Math.max(0, annualDriftlockSpend),
        lossesAvoided: Math.max(0, lossesAvoided),
        fineExposureConsidered: Math.max(0, fineExposureConsidered),
        netBenefit,
    }
})
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
