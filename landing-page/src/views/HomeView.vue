<template>
    <main>
        <!-- Hero Section -->
        <section class="relative overflow-hidden py-24 sm:py-32" style="background-image: url('/hero-driftlock.jpg'); background-size: cover; background-position: center;">
            <!-- Overlay -->
            <div class="absolute inset-0 bg-gray-900/70" aria-hidden="true"></div>
            
            <div class="relative z-10 container mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                <div class="mx-auto max-w-3xl text-center">
                    <span class="rounded-full bg-blue-500 px-4 py-1 text-sm font-sans font-medium text-white">
                        Alpha Phase - Seeking Design Partners
                    </span>
                    <h1 class="mt-6 text-5xl font-mono font-extrabold tracking-tight text-white sm:text-7xl">
                        Stop Multi-Million Dollar Fines.
                    </h1>
                    <p class="mt-6 text-lg leading-8 font-sans text-gray-200 sm:text-xl">
                        Driftlock is regulator-proof, explainable anomaly detection for global banks. Pass EU DORA & US regulatory audits with <strong class="font-semibold text-white">math-based proof</strong>, not black-box ML.
                    </p>
                    <div class="mt-10 flex items-center justify-center gap-x-6">
                        <a href="#contact" class="rounded-md cta-gradient px-6 py-3 text-base font-sans font-semibold text-white shadow-lg transition-opacity hover:opacity-90 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600">
                            Request Private Demo
                        </a>
                        <a href="#playground" class="group flex items-center gap-x-2 rounded-md border border-gray-300 bg-white px-6 py-3 text-base font-sans font-semibold text-gray-900 shadow-sm hover:bg-gray-50">
                            See Live Playground
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-500 transition-transform group-hover:translate-x-1" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
                            </svg>
                        </a>
                    </div>
                </div>
            </div>
        </section>

        <!-- Problem Section (The $ Risk) -->
        <section id="problem" class="py-24 sm:py-32 bg-gray-50">
            <div class="container mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                <div class="mx-auto max-w-3xl text-center">
                    <h2 class="text-base font-sans font-semibold leading-7 text-blue-600">The Problem</h2>
                    <p class="mt-2 text-4xl font-mono font-bold tracking-tight text-gray-900 sm:text-5xl">Your Black Box is an Audit Failure.</p>
                    <p class="mt-6 text-lg leading-8 font-sans text-gray-700">
                        Regulators worldwide are auditing bank systems for explainability. When your system flags a transaction, they demand: <strong class="font-semibold">"Show your work."</strong> Black-box ML models are an automatic failure, leading to massive fines.
                    </p>
                </div>
                <!-- Regulation Cards -->
                <div class="mt-16 grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-3">
                    <div class="rounded-2xl border border-gray-200 bg-white p-8 shadow-sm">
                        <h3 class="text-xl font-mono font-semibold leading-7 text-gray-900">EU DORA</h3>
                        <p class="mt-4 text-3xl font-mono font-bold tracking-tight text-red-600">Up to 2%</p>
                        <p class="mt-1 text-base font-sans text-gray-600">of annual global turnover in fines for non-compliance. Article 15 requires full explainability for all automated fraud detection.</p>
                    </div>
                    <div class="rounded-2xl border border-gray-200 bg-white p-8 shadow-sm">
                        <h3 class="text-xl font-mono font-semibold leading-7 text-gray-900">US Regulations</h3>
                        <p class="mt-4 text-3xl font-mono font-bold tracking-tight text-red-600">CFPB, NYDFS, FFIEC</p>
                        <p class="mt-1 text-base font-sans text-gray-600">Fair Lending, Part 500, and FFIEC guidelines all mandate model explainability, auditability, and fairness. Black boxes are a direct violation.</p>
                    </div>
                    <div class="rounded-2xl border border-gray-200 bg-white p-8 shadow-sm md:col-span-2 lg:col-span-1">
                        <h3 class="text-xl font-mono font-semibold leading-7 text-gray-900">The Compliance Gap</h3>
                        <p class="mt-4 text-3xl font-mono font-bold tracking-tight text-red-600">Automatic Failure</p>
                        <p class="mt-1 text-base font-sans text-gray-600">If you can't provide a mathematical proof for why a transaction was flagged, you have already failed your audit. Fines start January 2025.</p>
                    </div>
                </div>
            </div>
        </section>

        <!-- Solution Section -->
        <section id="solution" class="py-24 sm:py-32">
            <div class="container mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                <div class="lg:grid lg:grid-cols-12 lg:gap-16">
                    <!-- Left Column (Content) -->
                    <div class="lg:col-span-7">
                        <h2 class="text-base font-sans font-semibold leading-7 text-blue-600">The Solution</h2>
                        <p class="mt-2 text-4xl font-mono font-bold tracking-tight text-gray-900 sm:text-5xl">Math, Not Magic.</p>
                        <p class="mt-6 text-xl leading-8 font-sans text-gray-700">
                            Driftlock replaces opaque, black-box models with <strong class="font-semibold text-gray-900">compression-based anomaly detection.</strong>
                        </p>
                        <p class="mt-4 text-lg font-sans text-gray-600">
                            Instead of a "computer says no" answer, Driftlock provides a complete, human-readable mathematical proof for every flagged transaction. It's the "show your work" that auditors love, because it's verifiable math (NCD), not a probabilistic guess (ML).
                        </p>
                        <!-- Key Features (from screenshot) -->
                        <div class="mt-12 grid grid-cols-1 gap-x-8 gap-y-10 sm:grid-cols-2">
                            <div class="flex gap-x-4">
                                <div class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-lg bg-blue-100 text-blue-600">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                                    </svg>
                                </div>
                                <div>
                                    <h3 class="text-base font-mono font-semibold text-gray-900">Glass-Box Algorithm</h3>
                                    <p class="mt-1 text-sm font-sans text-gray-600">Mathematical explanations for every decision, no black boxes.</p>
                                </div>
                            </div>
                            <div class="flex gap-x-4">
                                <div class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-lg bg-blue-100 text-blue-600">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
                                    </svg>
                                </div>
                                <div>
                                    <h3 class="text-base font-mono font-semibold text-gray-900">Real-Time Detection</h3>
                                    <p class="mt-1 text-sm font-sans text-gray-600">Process 1000+ events/second with sub-second latency.</p>
                                </div>
                            </div>
                            <div class="flex gap-x-4">
                                <div class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-lg bg-blue-100 text-blue-600">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                                    </svg>
                                </div>
                                <div>
                                    <h3 class="text-base font-mono font-semibold text-gray-900">Audit-Ready Reports</h3>
                                    <p class="mt-1 text-sm font-sans text-gray-600">One-click compliance reports for regulator submissions.</p>
                                </div>
                            </div>
                            <div class="flex gap-x-4">
                                <div class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-lg bg-blue-100 text-blue-600">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3-.895 3-2-1.343-2-3-2z" />
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M12 14c-1.657 0-3 .895-3 2s1.343 2 3 2 3-.895 3-2-1.343-2-3-2z" />
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M12 2a10 10 0 100 20 10 10 0 000-20z" />
                                    </svg>
                                </div>
                                <div>
                                    <h3 class="text-base font-mono font-semibold text-gray-900">Compliance Proven</h3>
                                    <p class="mt-1 text-sm font-sans text-gray-600">Ready for EU DORA and US regulators with traceable math.</p>
                                </div>
                            </div>
                        </div>
                    </div>
                    
                    <!-- Right Column (Visual) -->
                    <div class="mt-16 lg:col-span-5 lg:mt-0">
                        <div class="relative h-full rounded-2xl border border-gray-200 bg-white p-2 shadow-xl">
                            <div class="absolute -top-4 -left-4 z-10 rounded-full bg-yellow-300 px-3 py-1 text-sm font-sans font-semibold text-gray-900 shadow-lg">
                                Audit-Ready Report
                            </div>
                            <!-- Placeholder for the `demo-output.html` screenshot -->
                            <div class="flex h-full min-h-[400px] flex-col items-center justify-center rounded-lg bg-gray-50 text-center">
                                <img src="/demo-anomaly-card.png" 
                                     alt="Screenshot of an audit-ready anomaly report card" 
                                     class="w-full h-auto rounded-lg shadow-md border border-gray-300">
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <!-- How It Works Section -->
        <section id="how" class="py-24 sm:py-32 bg-gray-900 text-white">
            <div class="container mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                <div class="mx-auto max-w-3xl text-center">
                    <h2 class="text-base font-sans font-semibold leading-7 text-blue-300">How It Works</h2>
                    <p class="mt-2 text-4xl font-mono font-bold tracking-tight text-white sm:text-5xl">From Data to Decision in Milliseconds</p>
                </div>
                <!-- 3-Step Process -->
                <div class="mt-16 grid grid-cols-1 gap-12 md:grid-cols-3">
                    <!-- Step 1 -->
                    <div class="flex flex-col items-center text-center">
                        <div class="flex h-16 w-16 items-center justify-center rounded-full border-2 border-blue-400 bg-blue-900 text-2xl font-mono font-bold text-white">1</div>
                        <h3 class="mt-6 text-xl font-mono font-semibold text-white">Build Baseline</h3>
                        <p class="mt-2 text-base font-sans text-gray-300">Driftlock ingests your normal transaction telemetry (e.g., first 400 events) to build a mathematical baseline of "normal" behavior.</p>
                    </div>
                    <!-- Step 2 -->
                    <div class="flex flex-col items-center text-center">
                        <div class="flex h-16 w-16 items-center justify-center rounded-full border-2 border-blue-400 bg-blue-900 text-2xl font-mono font-bold text-white">2</div>
                        <h3 class="mt-6 text-xl font-mono font-semibold text-white">Detect Anomalies</h3>
                        <p class="mt-2 text-base font-sans text-gray-300">New events are compared to the baseline using Normalized Compression Distance (NCD). Significant deviations are flagged as anomalies.</p>
                    </div>
                    <!-- Step 3 -->
                    <div class="flex flex-col items-center text-center">
                        <div class="flex h-16 w-16 items-center justify-center rounded-full border-2 border-blue-400 bg-blue-900 text-2xl font-mono font-to-bold text-white">3</div>
                        <h3 class="mt-6 text-xl font-mono font-semibold text-white">Generate Explanation</h3>
                        <p class="mt-2 text-base font-sans text-gray-300">An HTML report is generated with NCD scores, p-values, and statistical significance—a verifiable, audit-friendly mathematical proof.</p>
                    </div>
                </div>
            </div>
        </section>

        <!-- Live Playground Section -->
        <section id="playground" class="py-24 sm:py-32">
            <div class="container mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                <PlaygroundShell variant="embedded" />
            </div>
        </section>

        <!-- CTA Section -->
        <section id="contact" class="bg-gray-900">
            <div class="container mx-auto max-w-7xl px-4 py-24 sm:px-6 sm:py-32 lg:px-8">
                <div class="mx-auto max-w-3xl text-center">
                    <h2 class="text-3xl font-mono font-bold tracking-tight text-white sm:text-4xl">
                        Don't Risk Your Next Audit.
                        <br>
                        Pilot Driftlock for Q1 2026.
                    </h2>
                    <p class="mt-6 text-lg leading-8 font-sans text-gray-300">
                        We are actively seeking <strong class="font-semibold text-white">3 pilot banks</strong> for Q1 2026 to validate real-world performance. Be one of the first to have a regulator-proof, explainable fraud detection system. We need introductions to bank compliance teams facing regulatory audits.
                    </p>
                    <div class="mt-10">
                        <a href="mailto:hunter@shannonlabs.dev?subject=Driftlock Pilot Program Inquiry" class="rounded-md cta-gradient px-8 py-3 text-lg font-sans font-semibold text-white shadow-lg transition-opacity hover:opacity-90 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600">
                            Contact Us for a Pilot
                        </a>
                        <p class="mt-4 text-sm font-sans text-gray-400">Or reach out at hunter@shannonlabs.dev</p>
                    </div>
                </div>
            </div>
        </section>
    </main>
</template>

<script setup lang="ts">
import PlaygroundShell from '../components/playground/PlaygroundShell.vue'
</script>