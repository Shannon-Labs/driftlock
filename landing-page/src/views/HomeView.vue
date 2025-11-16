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
                        Spin Up Driftlock in Minutes.
                    </h1>
                    <p class="mt-6 text-lg leading-8 font-sans text-gray-200 sm:text-xl">
                        Launch the HTTP API + Postgres demo locally, hit <code class="text-white">/v1/detect</code>, and show regulators persisted anomalies with compression math—not black-box ML. Deterministic evidence for EU DORA &amp; US audits.
                    </p>
                    <p class="mt-4 text-base leading-7 font-sans text-gray-200/90 max-w-3xl mx-auto">
                        Designed as a <strong class="font-semibold">developer-first anomaly API</strong>: a simple HTTPS endpoint with a generous free tier and usage-based pricing that targets roughly <span class="font-mono">$1 per million anomaly checks</span>, with volume discounts and enterprise plans in the low-thousands per month.
                    </p>
                    <div class="mt-10 flex items-center justify-center gap-x-6">
                        <a href="#api-demo" class="rounded-md cta-gradient px-6 py-3 text-base font-sans font-semibold text-white shadow-lg transition-opacity hover:opacity-90 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600">
                            Run the API Demo
                        </a>
                        <a href="#contact" class="group flex items-center gap-x-2 rounded-md border border-blue-200 bg-white/90 px-6 py-3 text-base font-sans font-semibold text-gray-900 shadow-sm hover:bg-gray-50">
                            Request Private Pilot
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-500 transition-transform group-hover:translate-x-1" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
                            </svg>
                        </a>
                    </div>
                </div>
            </div>
        </section>

        <!-- API Demo Section -->
        <section id="api-demo" class="bg-slate-950 text-white py-20">
            <div class="container mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
                <div class="grid gap-12 lg:grid-cols-2 lg:items-center">
                    <div>
                        <p class="text-sm font-semibold uppercase tracking-[0.2em] text-blue-400">New onboarding</p>
                        <h2 class="mt-3 text-4xl font-mono font-bold">HTTP API Demo with Docker + Postgres</h2>
                        <p class="mt-4 text-lg text-slate-200">
                            The README, docs, and YC narrative now lead with a real API experience. Prospects run one script, provision a tenant/key, hit <code>/v1/detect</code>, and inspect anomalies in Postgres—no black-box magic.
                        </p>
                        <ol class="mt-6 space-y-4 text-slate-100">
                            <li class="flex items-start gap-3">
                                <span class="mt-1 inline-flex h-8 w-8 items-center justify-center rounded-full bg-blue-500/30 text-blue-200 font-bold">1</span>
                                <div>
                                    <p class="font-semibold">Clone, build, run</p>
                                    <p class="text-sm text-slate-300"><code>git clone</code>, <code>cargo build --release</code>, then <code>DRIFTLOCK_DEV_MODE=true ./scripts/run-api-demo.sh</code>.</p>
                                </div>
                            </li>
                            <li class="flex items-start gap-3">
                                <span class="mt-1 inline-flex h-8 w-8 items-center justify-center rounded-full bg-blue-500/30 text-blue-200 font-bold">2</span>
                                <div>
                                    <p class="font-semibold">Script output = live instructions</p>
                                    <p class="text-sm text-slate-300">We print <code>curl</code> + <code>psql</code> snippets so visitors can replay <code>/v1/detect</code>, fetch <code>/v1/anomalies/{id}</code>, and verify <code>/healthz</code> showing license + queue + DB status.</p>
                                </div>
                            </li>
                            <li class="flex items-start gap-3">
                                <span class="mt-1 inline-flex h-8 w-8 items-center justify-center rounded-full bg-blue-500/30 text-blue-200 font-bold">3</span>
                                <div>
                                    <p class="font-semibold">Proof in Postgres</p>
                                    <p class="text-sm text-slate-300">`psql` output shows anomalies persisted with NCD, p-values, and explanations—ready for compliance teams.</p>
                                </div>
                            </li>
                        </ol>
                        <p class="mt-8 text-sm text-slate-400">Dev mode bypass (`DRIFTLOCK_DEV_MODE=true`) is local-only. Production needs real keys; `/healthz` reports violations.</p>
                    </div>
                    <div class="space-y-6">
                        <video src="/api-demo-demo.mp4" autoplay muted loop playsinline controls poster="/api-demo-terminal.png" class="w-full rounded-2xl shadow-2xl border border-white/10">
                            Sorry, your browser doesn't support embedded videos. Watch the terminal capture in the README.
                        </video>
                        <div class="grid gap-4 sm:grid-cols-2">
                            <img src="/api-demo-detect.png" alt="Sample /v1/detect response" class="rounded-2xl shadow-xl border border-white/5" />
                            <img src="/api-demo-psql.png" alt="psql anomalies view" class="rounded-2xl shadow-xl border border-white/5" />
                        </div>
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

        <!-- ROI Calculator Section -->
        <section id="roi" class="py-24 sm:py-32 bg-slate-950 text-white">
            <div class="container mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                <div class="mx-auto max-w-3xl text-center">
                    <h2 class="text-base font-sans font-semibold leading-7 text-blue-300">Financial Impact</h2>
                    <p class="mt-2 text-4xl font-mono font-bold tracking-tight text-white sm:text-5xl">Rough ROI, In Plain Numbers</p>
                    <p class="mt-4 text-base font-sans text-slate-200">
                        This is a directional model, not a quote. It helps you reason about how anomaly detection that actually explains itself compares to fraud losses, compliance fines, and legacy tooling.
                    </p>
                </div>
                <div class="mt-12 grid gap-10 lg:grid-cols-2 lg:items-start">
                    <!-- Inputs -->
                    <div class="rounded-3xl bg-slate-900/70 p-8 shadow-2xl border border-white/10">
                        <h3 class="text-xl font-mono font-semibold text-white">Your Assumptions</h3>
                        <p class="mt-2 text-sm font-sans text-slate-300">Adjust these to match your environment. We never send this data anywhere.</p>
                        <div class="mt-6 space-y-6">
                            <div>
                                <label for="roi-events" class="flex items-center justify-between text-sm font-sans font-semibold text-slate-200">
                                    <span>Events per month monitored</span>
                                    <span class="text-xs font-mono text-slate-400">{{ roiInputs.eventsPerMonth.toLocaleString() }}</span>
                                </label>
                                <input
                                    id="roi-events"
                                    v-model.number="roiInputs.eventsPerMonth"
                                    type="range"
                                    min="1000000"
                                    max="1000000000"
                                    step="1000000"
                                    class="mt-2 w-full accent-blue-500"
                                />
                            </div>
                            <div>
                                <label for="roi-fraud" class="flex items-center justify-between text-sm font-sans font-semibold text-slate-200">
                                    <span>Annual fraud / operational loss at risk (USD)</span>
                                    <span class="text-xs font-mono text-slate-400">${{ roiInputs.annualLoss.toLocaleString() }}</span>
                                </label>
                                <input
                                    id="roi-fraud"
                                    v-model.number="roiInputs.annualLoss"
                                    type="range"
                                    min="500000"
                                    max="50000000"
                                    step="500000"
                                    class="mt-2 w-full accent-blue-500"
                                />
                            </div>
                            <div>
                                <label for="roi-reduction" class="flex items-center justify-between text-sm font-sans font-semibold text-slate-200">
                                    <span>Expected reduction in losses with Driftlock</span>
                                    <span class="text-xs font-mono text-slate-400">{{ roiInputs.lossReductionPct }}%</span>
                                </label>
                                <input
                                    id="roi-reduction"
                                    v-model.number="roiInputs.lossReductionPct"
                                    type="range"
                                    min="1"
                                    max="30"
                                    step="1"
                                    class="mt-2 w-full accent-blue-500"
                                />
                            </div>
                            <div>
                                <label for="roi-fine" class="flex items-center justify-between text-sm font-sans font-semibold text-slate-200">
                                    <span>Probabilistic annual fine exposure (USD)</span>
                                    <span class="text-xs font-mono text-slate-400">${{ roiInputs.fineExposure.toLocaleString() }}</span>
                                </label>
                                <input
                                    id="roi-fine"
                                    v-model.number="roiInputs.fineExposure"
                                    type="range"
                                    min="0"
                                    max="100000000"
                                    step="1000000"
                                    class="mt-2 w-full accent-blue-500"
                                />
                            </div>
                            <div>
                                <label for="roi-price" class="flex items-center justify-between text-sm font-sans font-semibold text-slate-200">
                                    <span>Assumed Driftlock price per 1M anomaly checks</span>
                                    <span class="text-xs font-mono text-slate-400">${{ roiInputs.pricePerMillion.toFixed(2) }}</span>
                                </label>
                                <input
                                    id="roi-price"
                                    v-model.number="roiInputs.pricePerMillion"
                                    type="range"
                                    min="0.5"
                                    max="3"
                                    step="0.1"
                                    class="mt-2 w-full accent-blue-500"
                                />
                            </div>
                        </div>
                        <p class="mt-4 text-xs font-sans text-slate-500">All values are illustrative. Actual pricing and discounts are set by contract.</p>
                    </div>

                    <!-- Outputs -->
                    <div class="rounded-3xl bg-slate-900/70 p-8 shadow-2xl border border-blue-500/40">
                        <h3 class="text-xl font-mono font-semibold text-white">Modeled Impact (Annual)</h3>
                        <div class="mt-6 grid gap-6 sm:grid-cols-2">
                            <div class="rounded-2xl bg-slate-900 p-4 border border-white/10">
                                <p class="text-xs font-sans uppercase tracking-[0.18em] text-slate-400">Estimated Driftlock spend</p>
                                <p class="mt-2 text-3xl font-mono font-bold text-blue-300">${{ roiOutputs.annualDriftlockSpend.toLocaleString() }}</p>
                                <p class="mt-2 text-xs font-sans text-slate-400">Based on your events/month and assumed price per 1M anomaly checks.</p>
                            </div>
                            <div class="rounded-2xl bg-slate-900 p-4 border border-white/10">
                                <p class="text-xs font-sans uppercase tracking-[0.18em] text-slate-400">Losses avoided</p>
                                <p class="mt-2 text-3xl font-mono font-bold text-emerald-300">${{ roiOutputs.lossesAvoided.toLocaleString() }}</p>
                                <p class="mt-2 text-xs font-sans text-slate-400">Annual fraud/operational losses avoided if Driftlock reduces losses by your chosen percentage.</p>
                            </div>
                            <div class="rounded-2xl bg-slate-900 p-4 border border-white/10">
                                <p class="text-xs font-sans uppercase tracking-[0.18em] text-slate-400">Fine exposure considered</p>
                                <p class="mt-2 text-3xl font-mono font-bold text-amber-300">${{ roiOutputs.fineExposureConsidered.toLocaleString() }}</p>
                                <p class="mt-2 text-xs font-sans text-slate-400">You can treat this as expected regulatory downside in a given year.</p>
                            </div>
                            <div class="rounded-2xl bg-slate-900 p-4 border border-white/10">
                                <p class="text-xs font-sans uppercase tracking-[0.18em] text-slate-400">Modeled net benefit</p>
                                <p class="mt-2 text-3xl font-mono font-bold" :class="roiOutputs.netBenefit >= 0 ? 'text-emerald-300' : 'text-red-300'">
                                    ${{ roiOutputs.netBenefit.toLocaleString() }}
                                </p>
                                <p class="mt-2 text-xs font-sans text-slate-400">(Losses avoided + fine exposure) − Driftlock spend. Positive values suggest the math is in your favor.</p>
                            </div>
                        </div>
                        <div class="mt-8 rounded-2xl border border-white/10 bg-slate-900 p-4 text-xs font-sans text-slate-300">
                            <p><strong class="text-white">Disclaimer:</strong> This is a simplified model using your assumptions. It is not financial advice and not a binding price quote. For real pricing, talk to us about your data, risk profile, and regulatory requirements.</p>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <!-- CTA Section -->
        <section id="contact" class="bg-gray-900">
            <div class="container mx-auto max-w-7xl px-4 py-24 sm:px-6 sm:py-32 lg:px-8">
                <div class="grid gap-12 lg:grid-cols-2 lg:items-start">
                    <div>
                        <p class="text-sm font-semibold uppercase tracking-[0.2em] text-blue-400">Pilot Program</p>
                        <h2 class="mt-4 text-3xl font-mono font-bold tracking-tight text-white sm:text-4xl">
                            Don't Risk Your Next Audit.
                            <br>
                            Pilot Driftlock for Q1 2026.
                        </h2>
                        <p class="mt-6 text-lg leading-8 font-sans text-gray-300">
                            We are actively seeking <strong class="font-semibold text-white">3 pilot banks</strong> for Q1 2026 to validate real-world performance. Be one of the first to have a regulator-proof, explainable fraud detection system. Introductions to compliance and risk owners are especially helpful.
                        </p>
                        <div class="mt-8 rounded-2xl border border-white/10 bg-white/5 p-6 text-gray-200">
                            <h3 class="text-lg font-semibold text-white">What to expect</h3>
                            <ul class="mt-4 space-y-3 text-sm font-sans text-gray-300">
                                <li class="flex items-start gap-3">
                                    <span class="mt-0.5 inline-flex h-6 w-6 flex-shrink-0 items-center justify-center rounded-full bg-blue-500/30 text-xs font-bold text-blue-200">1</span>
                                    <span>Qualify your dataset + compliance needs. We sign mutual NDA if required.</span>
                                </li>
                                <li class="flex items-start gap-3">
                                    <span class="mt-0.5 inline-flex h-6 w-6 flex-shrink-0 items-center justify-center rounded-full bg-blue-500/30 text-xs font-bold text-blue-200">2</span>
                                    <span>Deploy HTTP API + Postgres demo in your sandbox (guided runbook included).</span>
                                </li>
                                <li class="flex items-start gap-3">
                                    <span class="mt-0.5 inline-flex h-6 w-6 flex-shrink-0 items-center justify-center rounded-full bg-blue-500/30 text-xs font-bold text-blue-200">3</span>
                                    <span>Review anomalies + math proof with compliance teams and plan production pilot.</span>
                                </li>
                            </ul>
                            <p class="mt-4 text-sm text-gray-400">Prefer email? <a :href="fallbackMailto" class="font-semibold text-white underline underline-offset-4">hunter@shannonlabs.dev</a></p>
                        </div>
                    </div>
                    <form @submit.prevent="handleContactSubmit" class="rounded-3xl bg-white p-8 shadow-2xl" novalidate>
                        <fieldset :disabled="isSubmitting" class="space-y-6">
                            <legend class="text-2xl font-mono font-bold text-gray-900">Tell us about your team</legend>
                            <div class="grid gap-6 md:grid-cols-2">
                                <div>
                                    <label for="contact-name" class="block text-sm font-sans font-semibold text-gray-700">Full Name *</label>
                                    <input
                                        id="contact-name"
                                        v-model="contactForm.name"
                                        type="text"
                                        required
                                        autocomplete="name"
                                        class="mt-2 w-full rounded-2xl border border-gray-200 bg-gray-50 px-4 py-3 text-base font-sans text-gray-900 shadow-sm focus:border-blue-500 focus:bg-white focus:outline-none focus:ring-2 focus:ring-blue-500/40"
                                        placeholder="Alex Compliance"
                                    />
                                </div>
                                <div>
                                    <label for="contact-email" class="block text-sm font-sans font-semibold text-gray-700">Work Email *</label>
                                    <input
                                        id="contact-email"
                                        v-model="contactForm.email"
                                        type="email"
                                        required
                                        autocomplete="email"
                                        class="mt-2 w-full rounded-2xl border border-gray-200 bg-gray-50 px-4 py-3 text-base font-sans text-gray-900 shadow-sm focus:border-blue-500 focus:bg-white focus:outline-none focus:ring-2 focus:ring-blue-500/40"
                                        placeholder="alex@bank.com"
                                    />
                                </div>
                            </div>
                            <div>
                                <label for="contact-company" class="block text-sm font-sans font-semibold text-gray-700">Company</label>
                                <input
                                    id="contact-company"
                                    v-model="contactForm.company"
                                    type="text"
                                    autocomplete="organization"
                                    class="mt-2 w-full rounded-2xl border border-gray-200 bg-gray-50 px-4 py-3 text-base font-sans text-gray-900 shadow-sm focus:border-blue-500 focus:bg-white focus:outline-none focus:ring-2 focus:ring-blue-500/40"
                                    placeholder="Global Bank Corp"
                                />
                            </div>
                            <div>
                                <label for="contact-message" class="block text-sm font-sans font-semibold text-gray-700">How can we help? *</label>
                                <textarea
                                    id="contact-message"
                                    v-model="contactForm.message"
                                    required
                                    rows="4"
                                    class="mt-2 w-full rounded-2xl border border-gray-200 bg-gray-50 px-4 py-3 text-base font-sans text-gray-900 shadow-sm focus:border-blue-500 focus:bg-white focus:outline-none focus:ring-2 focus:ring-blue-500/40"
                                    placeholder="We need help preparing for a DORA readiness audit..."
                                ></textarea>
                            </div>
                            <div class="flex items-center justify-between gap-3">
                                <button
                                    type="submit"
                                    :disabled="!isFormValid || isSubmitting"
                                    class="inline-flex w-full items-center justify-center rounded-2xl bg-gradient-to-r from-blue-600 to-indigo-600 px-6 py-3 text-base font-sans font-semibold text-white shadow-lg transition focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 disabled:cursor-not-allowed disabled:opacity-50"
                                >
                                    <svg v-if="isSubmitting" class="mr-3 h-5 w-5 animate-spin text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"></path>
                                    </svg>
                                    {{ isSubmitting ? 'Sending...' : 'Request a Pilot Call' }}
                                </button>
                            </div>
                            <p class="text-xs font-sans text-gray-500">By submitting, you agree to be contacted about Driftlock pilots. No spam.</p>
                        </fieldset>
                        <div class="mt-6 min-h-[48px]" aria-live="polite">
                            <div v-if="submissionState === 'success'" class="flex items-center rounded-2xl border border-green-200 bg-green-50 px-4 py-3 text-sm font-sans text-green-900">
                                <svg class="mr-2 h-5 w-5 text-green-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-7.778 7.778a1 1 0 01-1.414 0L3.293 9.293a1 1 0 011.414-1.414l3.102 3.102 7.071-7.071a1 1 0 011.414 0z" clip-rule="evenodd" />
                                </svg>
                                Thanks! We received your note and will respond within 1 business day.
                            </div>
                            <div v-else-if="submissionState === 'error'" class="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm font-sans text-red-900">
                                {{ submissionMessage }}
                                <span class="block text-xs text-red-700">Or email <a :href="fallbackMailto" class="underline font-semibold">hunter@shannonlabs.dev</a></span>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
        </section>
    </main>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import PlaygroundShell from '../components/playground/PlaygroundShell.vue'

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
    if (!isFormValid.value) {
        submissionState.value = 'error'
        submissionMessage.value = 'Please fill out the required fields before submitting.'
        return
    }

    isSubmitting.value = true
    submissionState.value = 'idle'
    submissionMessage.value = ''

    try {
        const response = await fetch(contactEndpoint, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                name: contactForm.name.trim(),
                email: contactForm.email.trim(),
                company: contactForm.company.trim(),
                message: contactForm.message.trim(),
            }),
        })

        if (!response.ok) {
            throw new Error(`Request failed with status ${response.status}`)
        }

        submissionState.value = 'success'
        resetForm()
    } catch (error) {
        console.error('Contact form submission failed', error)
        submissionState.value = 'error'
        submissionMessage.value = 'Something went wrong. Please try again or reach out via email.'
    } finally {
        isSubmitting.value = false
    }
}

// ROI calculator state (purely client-side, illustrative only)
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
