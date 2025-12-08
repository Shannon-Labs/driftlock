<template>
    <main class="bg-background text-foreground font-sans selection:bg-black selection:text-white">
        <!-- Hero Section -->
        <section class="relative w-full overflow-hidden border-b border-black bg-white text-black lg:h-[85vh] min-h-[700px] flex flex-col">
            <!-- Brutalist Grid Background -->
            <div class="absolute inset-0 z-0 opacity-[0.03]"
                 style="background-image: linear-gradient(#000 1px, transparent 1px), linear-gradient(90deg, #000 1px, transparent 1px); background-size: 40px 40px;">
            </div>
            
            <div class="container-padding relative z-10 grid lg:grid-cols-2 gap-12 h-full py-12 lg:py-0">
                <!-- Left: Content -->
                <div class="flex flex-col justify-center">
                    <div class="mb-8">
                        <div class="inline-flex items-center gap-2 border border-black px-3 py-1 mb-6 bg-black text-white shadow-[4px_4px_0px_0px_rgba(0,0,0,0.2)]">
                            <span class="w-2 h-2 bg-green-400 animate-pulse rounded-full"></span>
                            <span class="text-xs font-bold font-mono tracking-widest uppercase">System Online</span>
                        </div>
                        
                        <h1 class="text-6xl sm:text-7xl lg:text-8xl font-sans font-bold tracking-tight leading-[0.9] mb-8 uppercase">
                            Mathematical<br>Observability.
                        </h1>

                        <p class="text-xl font-serif leading-relaxed text-gray-800 border-l-4 border-black pl-6 max-w-xl mb-10">
                            Detect drift, anomalies, and outliers in any stream without training data. 
                            <strong>O(1) complexity signal processing for modern engineering teams.</strong>
                        </p>

                        <div class="flex flex-col sm:flex-row gap-4 items-start sm:items-center">
                            <a href="#signup" class="brutalist-button-primary text-lg px-8 py-4">
                                Join Waitlist
                            </a>
                        </div>
                    </div>
                </div>

                <!-- Right: Clean geometric visual -->
                <div class="hidden lg:flex flex-col justify-center lg:pl-12 h-full">
                    <div class="relative w-full max-w-md aspect-square">
                        <!-- Minimalist geometric pattern -->
                        <div class="absolute inset-0 border-2 border-black"></div>
                        <div class="absolute inset-4 border border-black"></div>
                        <div class="absolute inset-8 border border-black/50"></div>
                        <div class="absolute inset-12 border border-black/30"></div>
                        <div class="absolute inset-0 flex items-center justify-center">
                            <div class="text-center">
                                <div class="text-6xl font-mono font-bold tracking-tighter">O(1)</div>
                                <div class="text-xs font-mono uppercase tracking-widest mt-2 text-gray-500">Complexity</div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <!-- Use Cases / Capabilities (Features) -->
        <section id="features" class="section-padding bg-white">
            <div class="container-padding">
                <div class="flex flex-col md:flex-row md:items-end justify-between mb-16 border-b border-black pb-8">
                    <div class="max-w-2xl">
                        <span class="text-xs font-bold font-mono uppercase tracking-widest mb-2 block text-gray-500">Capabilities</span>
                        <h2 class="text-4xl md:text-5xl font-sans font-bold tracking-tight uppercase">
                            Pure Signal.<br>Zero Noise.
                        </h2>
                    </div>
                    <div class="mt-6 md:mt-0 font-serif text-lg max-w-xs text-right hidden md:block">
                        Replace "black box" ML with transparent signal processing.
                    </div>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-3 gap-0 border-t border-l border-black">
                    <!-- Feature 1 -->
                    <div class="border-b border-r border-black p-8 hover:bg-gray-50 transition-colors group">
                        <div class="mb-6 group-hover:-translate-y-1 transition-transform">
                            <svg class="w-10 h-10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                                <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z" stroke-linecap="round" stroke-linejoin="round"/>
                            </svg>
                        </div>
                        <h3 class="text-xl font-bold font-sans uppercase mb-4">0ms Training Time</h3>
                        <p class="font-serif text-gray-600 leading-relaxed">
                            No learning period. No labeled datasets. Driftlock uses entropy-based compression to detect anomalies from the very first event.
                        </p>
                    </div>

                    <!-- Feature 2 -->
                    <div class="border-b border-r border-black p-8 hover:bg-gray-50 transition-colors group">
                        <div class="mb-6 group-hover:-translate-y-1 transition-transform">
                            <svg class="w-10 h-10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                                <circle cx="11" cy="11" r="8"/>
                                <path d="M21 21l-4.35-4.35" stroke-linecap="round"/>
                            </svg>
                        </div>
                        <h3 class="text-xl font-bold font-sans uppercase mb-4">White-Box Proofs</h3>
                        <p class="font-serif text-gray-600 leading-relaxed">
                            Every alert comes with a mathematical proof of *why* it is anomalous. Perfect for compliance workflows (DORA, NIS2, AI Act).
                        </p>
                    </div>

                    <!-- Feature 3 -->
                    <div class="border-b border-r border-black p-8 hover:bg-gray-50 transition-colors group">
                        <div class="mb-6 group-hover:-translate-y-1 transition-transform">
                            <svg class="w-10 h-10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                                <circle cx="12" cy="12" r="10"/>
                                <path d="M2 12h20M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
                            </svg>
                        </div>
                        <h3 class="text-xl font-bold font-sans uppercase mb-4">OpenTelemetry Native</h3>
                        <p class="font-serif text-gray-600 leading-relaxed">
                            Drop-in replacement for your current collector. Compatible with Prometheus, Datadog, and Honeycomb formats.
                        </p>
                    </div>
                </div>

                <!-- Adaptive Detection Highlight -->
                <div class="mt-8 flex flex-wrap items-center gap-4 text-sm font-mono text-gray-500">
                    <span class="font-bold uppercase">Also included:</span>
                    <span class="border border-gray-300 px-3 py-1">Adaptive Profiles</span>
                    <span class="border border-gray-300 px-3 py-1">Auto-Tuning</span>
                    <span class="border border-gray-300 px-3 py-1">Feedback Learning</span>
                </div>
            </div>
        </section>

        <!-- Audience (No fake social proof, just honest targeting) -->
        <section class="py-20 bg-black text-white border-b border-black">
            <div class="container-padding">
                <div class="grid lg:grid-cols-2 gap-16 items-center">
                    <div>
                        <h2 class="text-4xl font-sans font-bold tracking-tight uppercase mb-6">
                            Built for Builders.<br>Ready for Auditors.
                        </h2>
                        <p class="font-serif text-xl text-gray-400 max-w-md">
                            Driftlock bridges the gap between fast-moving engineering teams and strict compliance requirements.
                        </p>
                    </div>
                    
                    <div class="grid sm:grid-cols-2 gap-4">
                        <div class="border border-white/20 p-6 hover:border-white transition-colors">
                            <div class="font-mono text-xs uppercase tracking-widest text-gray-500 mb-2">Use Case 01</div>
                            <h3 class="font-bold text-lg mb-2">Agent loops</h3>
                            <p class="text-sm text-gray-400">Detect infinite loops and hallucinated tool calls in AI agents.</p>
                        </div>
                        <div class="border border-white/20 p-6 hover:border-white transition-colors">
                            <div class="font-mono text-xs uppercase tracking-widest text-gray-500 mb-2">Use Case 02</div>
                            <h3 class="font-bold text-lg mb-2">Market Data</h3>
                            <p class="text-sm text-gray-400">Validate real-time feeds for microstructure anomalies.</p>
                        </div>
                        <div class="border border-white/20 p-6 hover:border-white transition-colors">
                            <div class="font-mono text-xs uppercase tracking-widest text-gray-500 mb-2">Use Case 03</div>
                            <h3 class="font-bold text-lg mb-2">IoT Telemetry</h3>
                            <p class="text-sm text-gray-400">Monitor thousands of sensors with O(1) resource usage.</p>
                        </div>
                        <div class="border border-white/20 p-6 hover:border-white transition-colors">
                            <div class="font-mono text-xs uppercase tracking-widest text-gray-500 mb-2">Use Case 04</div>
                            <h3 class="font-bold text-lg mb-2">Security</h3>
                            <p class="text-sm text-gray-400">Spot data exfiltration patterns cloaked as normal traffic.</p>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <!-- Pricing Section -->
        <section id="pricing" class="section-padding bg-gray-50">
            <div class="container-padding">
                <div class="max-w-3xl mb-16 text-center mx-auto">
                    <span class="text-xs font-bold font-mono uppercase tracking-widest mb-2 block text-gray-500">Pricing</span>
                    <h2 class="text-4xl font-sans font-bold tracking-tight uppercase">
                        Fair, Predictable Pricing
                    </h2>
                </div>

                <div class="grid gap-6 lg:grid-cols-4 items-start">
                    <!-- Free -->
                    <div class="bg-white border border-gray-200 p-8 hover:border-black transition-colors">
                        <div class="mb-4">
                            <h3 class="font-bold font-sans uppercase text-lg">Free</h3>
                            <div class="text-3xl font-bold font-mono mt-2">$0<span class="text-sm font-normal text-gray-500">/mo</span></div>
                            <p class="text-sm text-gray-500 mt-2 font-serif">For experimentation.</p>
                        </div>
                        <ul class="space-y-3 text-sm font-mono border-t border-gray-100 pt-4 mb-8">
                            <li class="flex gap-2"><span>[x]</span> 10k events/mo</li>
                            <li class="flex gap-2"><span>[x]</span> CLI & Playground</li>
                            <li class="flex gap-2"><span>[x]</span> JSON reports</li>
                        </ul>
                        <a href="#signup" class="block w-full border border-black py-2 text-center text-sm font-bold uppercase hover:bg-black hover:text-white transition-colors">Waitlist</a>
                    </div>

                    <!-- Standard (Highlighted) -->
                    <div class="bg-black text-white border border-black p-8 shadow-[12px_12px_0px_0px_rgba(0,0,0,0.2)] relative lg:-mt-4 lg:mb-4 transform scale-105 z-10">
                        <div class="absolute top-0 right-0 bg-white text-black text-[10px] font-bold uppercase px-2 py-1">Most Popular</div>
                        <div class="mb-4">
                            <h3 class="font-bold font-sans uppercase text-lg text-white">Standard</h3>
                            <div class="text-3xl font-bold font-mono mt-2">$15<span class="text-sm font-normal text-gray-400">/mo</span></div>
                            <p class="text-sm text-gray-400 mt-2 font-serif">For startups & active monitoring.</p>
                        </div>
                        <ul class="space-y-3 text-sm font-mono border-t border-gray-800 pt-4 mb-8">
                            <li class="flex gap-2"><span class="text-green-400">[x]</span> 500k events/mo</li>
                            <li class="flex gap-2"><span class="text-green-400">[x]</span> Email alerts</li>
                            <li class="flex gap-2"><span class="text-green-400">[x]</span> 30-day retention</li>
                        </ul>
                        <button @click="handleCheckout('standard')" class="block w-full bg-white text-black py-3 text-center text-sm font-bold uppercase hover:bg-gray-200 transition-colors">
                            Get Standard
                        </button>
                    </div>

                    <!-- Pro -->
                    <div class="bg-white border border-gray-200 p-8 hover:border-black transition-colors">
                        <div class="mb-4">
                            <h3 class="font-bold font-sans uppercase text-lg">Pro</h3>
                            <div class="text-3xl font-bold font-mono mt-2">$100<span class="text-sm font-normal text-gray-500">/mo</span></div>
                            <p class="text-sm text-gray-500 mt-2 font-serif">For high-volume production.</p>
                        </div>
                        <ul class="space-y-3 text-sm font-mono border-t border-gray-100 pt-4 mb-8">
                            <li class="flex gap-2"><span>[x]</span> 5M events/mo</li>
                            <li class="flex gap-2"><span>[x]</span> DORA/NIS2 Evidence</li>
                            <li class="flex gap-2"><span>[x]</span> Priority Support</li>
                        </ul>
                        <button @click="handleCheckout('pro')" class="block w-full border border-black py-2 text-center text-sm font-bold uppercase hover:bg-black hover:text-white transition-colors">
                            Get Pro
                        </button>
                    </div>

                    <!-- Enterprise -->
                    <div class="bg-white border border-gray-200 p-8 hover:border-black transition-colors">
                        <div class="mb-4">
                            <h3 class="font-bold font-sans uppercase text-lg">Enterprise</h3>
                            <div class="text-3xl font-bold font-mono mt-2">$299<span class="text-sm font-normal text-gray-500">/mo</span></div>
                            <p class="text-sm text-gray-500 mt-2 font-serif">Large-scale compliance.</p>
                        </div>
                        <ul class="space-y-3 text-sm font-mono border-t border-gray-100 pt-4 mb-8">
                            <li class="flex gap-2"><span>[x]</span> 25M events/mo</li>
                            <li class="flex gap-2"><span>[x]</span> 90-day retention</li>
                            <li class="flex gap-2"><span>[x]</span> SSO (Coming soon)</li>
                        </ul>
                        <button @click="handleCheckout('enterprise')" class="block w-full border border-black py-2 text-center text-sm font-bold uppercase hover:bg-black hover:text-white transition-colors">
                            Contact Us
                        </button>
                    </div>
                </div>
            </div>
        </section>

        <!-- Waitlist / Footer CTA -->
        <section id="signup" class="section-padding bg-white border-b-0">
            <div class="container-padding max-w-2xl mx-auto text-center">
                <h2 class="text-4xl font-sans font-bold tracking-tight uppercase mb-6">Start Monitoring Today</h2>
                <div class="bg-gray-50 border border-black p-8 shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
                    <WaitlistForm />
                </div>
                <p class="mt-8 text-sm font-mono text-gray-400">
                    Already have access? <a href="/login" class="underline hover:text-black">Sign in →</a>
                </p>
            </div>
        </section>

        <!-- Checkout Error Toast -->
        <Transition
            enter-active-class="transition duration-300 ease-out"
            enter-from-class="transform translate-y-5 opacity-0"
            enter-to-class="transform translate-y-0 opacity-100"
            leave-active-class="transition duration-200 ease-in"
            leave-from-class="transform translate-y-0 opacity-100"
            leave-to-class="transform translate-y-5 opacity-0"
        >
            <div v-if="checkoutError"
                 class="fixed bottom-4 left-1/2 -translate-x-1/2 z-50
                        bg-red-600 text-white px-6 py-4
                        border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]
                        flex items-center gap-3 max-w-md">
                <span class="font-bold">ERROR:</span>
                <span class="text-sm font-mono">{{ checkoutError }}</span>
                <button @click="checkoutError = null" class="ml-2 hover:opacity-75">✕</button>
            </div>
        </Transition>
    </main>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import WaitlistForm from '../components/cta/WaitlistForm.vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const checkoutLoading = ref<string | null>(null)
const checkoutError = ref<string | null>(null)

// Handle pricing tier checkout
const handleCheckout = async (plan: string) => {
    checkoutError.value = null

    // If not authenticated, scroll to signup
    if (!authStore.isAuthenticated) {
        // For demo purposes, we just scroll to waitlist
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
            checkoutError.value = 'INIT_FAIL: Payment gateway unreachable'
        }
    } catch (e) {
        checkoutError.value = 'NET_ERR: Connection timed out'
    } finally {
        checkoutLoading.value = null
    }
}
</script>

<style scoped>
.section-padding {
    @apply py-16 sm:py-24 border-b border-black;
}

.container-padding {
    @apply mx-auto max-w-7xl px-4 sm:px-6 lg:px-8;
}

/* Custom selection color for that 'hacker' feel */
::selection {
  background-color: #000;
  color: #fff;
}
</style>
