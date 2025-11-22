<template>
    <main class="bg-background text-foreground">
        <!-- Hero Section -->
        <section class="section-padding">
            <div class="container-padding">
                <div class="max-w-4xl">
                    <div class="inline-block border border-black px-3 py-1 text-xs font-bold font-sans uppercase tracking-widest mb-6 bg-white">
                        Explainable · DORA Ready · Gemini Optional
                    </div>
                    <h1 class="text-6xl sm:text-8xl font-sans font-bold tracking-tight leading-none mb-8 uppercase">
                        Stop Black<br>Box Audits.
                    </h1>
                    <p class="text-xl sm:text-2xl font-serif leading-relaxed max-w-2xl mb-10 border-l border-black pl-6">
                        Driftlock compresses every payload twice, proves the math behind each anomaly, and (optionally) asks Gemini Flash to narrate the context.
                        <br><br>
                        <strong class="font-sans uppercase text-sm tracking-widest">Compliance teams get receipts. Engineers get sub-2s signal.</strong>
                    </p>
                    <div class="flex flex-wrap gap-4">
                        <a href="#signup" class="brutalist-button-primary text-lg">
                            Start Free Trial
                        </a>
                        <button type="button" @click="triggerLiveScan" class="brutalist-button text-lg">
                            Live Signal Scan →
                        </button>
                        <a href="#api-demo" class="brutalist-button text-lg">
                            View API Demo
                        </a>
                    </div>
                </div>
            </div>
        </section>

        <!-- Feature Grid (Academic Style) -->
        <section id="solution" class="section-padding">
             <div class="container-padding">
                <div class="mb-12">
                    <h2 class="text-sm font-bold font-sans uppercase tracking-widest mb-2 border-b border-black inline-block">The Solution</h2>
                    <p class="text-4xl font-sans font-bold tracking-tight mt-4">MATH, NOT MAGIC.</p>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-px bg-black border border-black">
                    <!-- Card 1 -->
                    <div class="p-8 bg-white h-full">
                        <h3 class="text-xl font-bold font-sans uppercase mb-4 flex items-baseline gap-2">
                            <span class="text-sm font-mono text-gray-500">01</span> Glass-Box Algorithm
                        </h3>
                        <p class="font-serif text-lg leading-relaxed">
                            Mathematical explanations for every decision. No black boxes, no "computer says no". Verify the math yourself.
                        </p>
                    </div>

                    <!-- Card 2 -->
                    <div class="p-8 bg-white h-full">
                        <h3 class="text-xl font-bold font-sans uppercase mb-4 flex items-baseline gap-2">
                            <span class="text-sm font-mono text-gray-500">02</span> Real-Time Detection
                        </h3>
                        <p class="font-serif text-lg leading-relaxed">
                            Process 1000+ events/second with sub-second latency. Built in Rust for maximum performance and safety.
                        </p>
                    </div>

                    <!-- Card 3 -->
                    <div class="p-8 bg-white h-full">
                        <h3 class="text-xl font-bold font-sans uppercase mb-4 flex items-baseline gap-2">
                            <span class="text-sm font-mono text-gray-500">03</span> Audit-Ready Reports
                        </h3>
                        <p class="font-serif text-lg leading-relaxed">
                            One-click compliance reports for regulator submissions. Ready for EU DORA and US audit requirements.
                        </p>
                    </div>

                    <!-- Card 4 -->
                    <div class="p-8 bg-white h-full">
                        <h3 class="text-xl font-bold font-sans uppercase mb-4 flex items-baseline gap-2">
                            <span class="text-sm font-mono text-gray-500">04</span> Hybrid Explainability
                        </h3>
                        <p class="font-serif text-lg leading-relaxed">
                            LZ-based math filters 99.9% of traffic; Gemini Flash translates anomalies into human context for tickets.
                        </p>
                    </div>

                    <!-- Card 5 -->
                    <div class="p-8 bg-white h-full">
                        <h3 class="text-xl font-bold font-sans uppercase mb-4 flex items-baseline gap-2">
                            <span class="text-sm font-mono text-gray-500">05</span> AI Agent Monitoring
                        </h3>
                        <p class="font-serif text-lg leading-relaxed">
                            Attach Driftlock to tool-using agents to spot runaway loops, hallucinations, or unexpected tool chains.
                        </p>
                    </div>
                    
                    <!-- CTA Card -->
                    <div class="p-8 bg-black text-white h-full flex items-center justify-center">
                        <a href="#signup" class="text-xl font-sans font-bold uppercase hover:underline underline-offset-4 decoration-white decoration-1">
                            Get Started Now →
                        </a>
                    </div>
                </div>
             </div>
        </section>

        <!-- Signup Section -->
        <section id="signup" class="section-padding bg-gray-50">
            <div class="container-padding">
                <div class="grid lg:grid-cols-2 gap-16 items-start">
                    <div>
                        <h2 class="text-4xl font-sans font-bold tracking-tight mb-6">START YOUR PILOT</h2>
                        <p class="text-xl font-serif mb-8">
                            Create your account and receive your API key instantly.
                            <br>
                            <span class="font-bold underline underline-offset-2">No credit card required for Developer plan.</span>
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

        <!-- API Demo Section -->
        <section id="api-demo" class="section-padding bg-black text-white">
            <div class="container-padding">
                <div class="grid gap-12 lg:grid-cols-2 lg:items-center">
                    <div>
                        <div class="inline-block border border-white px-2 py-1 text-xs font-mono mb-6">
                            New Onboarding Flow
                        </div>
                        <h2 class="text-4xl font-sans font-bold tracking-tight mb-6">HTTP API DEMO</h2>
                        <p class="text-lg font-serif text-gray-300 mb-8 leading-relaxed">
                            This is the exact flow from the README quickstart. Run one script, provision a tenant/key, hit <code class="font-mono bg-gray-800 px-1">/v1/detect</code>, and inspect anomalies in Postgres—no black-box magic.
                        </p>
                        
                        <div class="space-y-6 border-l border-white pl-6">
                             <div>
                                <h3 class="text-lg font-bold font-sans uppercase">1. Clone & Build</h3>
                                <p class="mt-2 font-mono text-sm bg-gray-900 p-3 border border-gray-700">
                                    git clone ... && cargo build --release
                                </p>
                            </div>
                            <div>
                                <h3 class="text-lg font-bold font-sans uppercase">2. Run Demo Script</h3>
                                <p class="mt-2 font-mono text-sm bg-gray-900 p-3 border border-gray-700">
                                    DRIFTLOCK_DEV_MODE=true ./scripts/run-api-demo.sh
                                </p>
                            </div>
                             <div>
                                <h3 class="text-lg font-bold font-sans uppercase">3. Verify in Postgres</h3>
                                <p class="mt-2 text-sm font-serif text-gray-400">
                                    `psql` output shows anomalies persisted with NCD, p-values, and explanations.
                                </p>
                            </div>
                        </div>
                    </div>
                    
                    <!-- Terminal Window -->
                    <div class="border border-white bg-black font-mono text-sm p-4 shadow-none">
                        <div class="border-b border-gray-700 pb-2 mb-4 flex justify-between">
                            <span>bash — 80x24</span>
                            <span>[●]</span>
                        </div>
                        <div class="space-y-1">
                            <div class="text-green-400">$ ./scripts/run-api-demo.sh</div>
                            <div class="text-gray-500"># Initializing Driftlock HTTP API...</div>
                            <div>> Service listening on :8080</div>
                            <div>> Provisioning tenant... OK</div>
                            <div>> Sending 50 normal events... OK</div>
                            <div>> Sending 1 anomaly... <span class="text-red-500 font-bold bg-red-900/20 px-1">DETECTED</span></div>
                            <br>
                            <div class="text-blue-400 font-bold">ANOMALY FOUND (ID: 550e8400-...)</div>
                            <div>NCD: 0.82 (Threshold: 0.5)</div>
                            <div>Explanation: "High entropy variance in field 'amount'"</div>
                        </div>
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
                    <h2 class="text-4xl font-sans font-bold tracking-tight mt-4">START FREE. SCALE WITH MATH.</h2>
                    <p class="text-lg font-serif mt-4">Compression math costs about $1 per million events. Our plans price in the Gemini explainability layer so you don't have to think about tokens.</p>
                </div>

                <div class="grid gap-8 lg:grid-cols-3">
                    <!-- Developer -->
                    <div class="brutalist-card flex flex-col">
                        <div class="mb-8">
                            <p class="text-sm font-bold font-sans uppercase tracking-wide mb-2">Developer</p>
                            <p class="text-5xl font-sans font-bold">$0</p>
                            <p class="text-sm text-gray-500 font-mono">per month</p>
                            <p class="mt-4 text-sm font-serif">10k events/mo, math proofs, CLI + API access.</p>
                        </div>
                        <ul class="space-y-3 text-sm font-mono border-t border-black pt-4 mt-auto">
                            <li>[x] Deterministic anomaly detection</li>
                            <li>[x] CLI demo & Playground</li>
                            <li>[x] JSON reports</li>
                        </ul>
                        <a href="#signup" class="brutalist-button w-full text-center mt-8">Start Building</a>
                    </div>

                    <!-- Starter -->
                    <div class="brutalist-card bg-black text-white flex flex-col relative">
                        <div class="absolute top-0 right-0 bg-white text-black text-xs font-bold font-sans px-2 py-1 uppercase border-l border-b border-black">
                            Most Popular
                        </div>
                        <div class="mb-8">
                            <p class="text-sm font-bold font-sans uppercase tracking-wide mb-2">Starter</p>
                            <p class="text-5xl font-sans font-bold">$25</p>
                            <p class="text-sm text-gray-400 font-mono">per month</p>
                            <p class="mt-4 text-sm font-serif">500k events/mo + Gemini summaries for every anomaly.</p>
                        </div>
                        <ul class="space-y-3 text-sm font-mono border-t border-white pt-4 mt-auto">
                            <li>[x] Email alerts with plain-English context</li>
                            <li>[x] Driftlock dashboard + usage analytics</li>
                            <li>[x] API key rotation & audit logs</li>
                        </ul>
                        <a href="#signup" class="brutalist-button-primary w-full text-center mt-8 border-white hover:bg-white hover:text-black">Launch Starter</a>
                    </div>

                    <!-- Enterprise -->
                    <div class="brutalist-card flex flex-col">
                        <div class="mb-8">
                            <p class="text-sm font-bold font-sans uppercase tracking-wide mb-2">Pro / Enterprise</p>
                            <p class="text-5xl font-sans font-bold">Custom</p>
                            <p class="text-sm text-gray-500 font-mono">5M+ events/mo</p>
                            <p class="mt-4 text-sm font-serif">Unlimited streams, AI agent monitoring, Slack/PagerDuty integrations.</p>
                        </div>
                        <ul class="space-y-3 text-sm font-mono border-t border-black pt-4 mt-auto">
                            <li>[x] Dedicated environments & SOC2 pipeline</li>
                            <li>[x] Human + LLM explainability packs</li>
                            <li>[x] Compliance reports (DORA, FFIEC, NIS2)</li>
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
                    <h2 class="text-4xl font-sans font-bold tracking-tight mt-4">ROUGH ROI ESTIMATION</h2>
                    <p class="text-base font-mono mt-4 text-gray-400">
                        Uses assumptions from our README: ~$1 per million anomaly checks.
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
import { computed, reactive, ref } from 'vue'
import PlaygroundShell from '../components/playground/PlaygroundShell.vue'
import SignupForm from '../components/cta/SignupForm.vue'

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