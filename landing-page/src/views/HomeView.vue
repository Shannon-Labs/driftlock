<template>
    <main>
        <!-- Hero Section -->
        <section class="relative isolate overflow-hidden py-24 sm:py-32">
            <div class="absolute inset-0 -z-10 bg-slate-950">
                <div class="absolute inset-0 opacity-40 bg-[radial-gradient(circle_at_top,rgba(14,165,233,0.35),transparent_55%),linear-gradient(120deg,rgba(15,23,42,0.95),rgba(8,47,73,0.92))]"></div>
                <div class="absolute inset-0 bg-[linear-gradient(to_right,rgba(255,255,255,0.05)_1px,transparent_1px),linear-gradient(to_bottom,rgba(255,255,255,0.05)_1px,transparent_1px)] bg-[size:32px_32px] opacity-20 [mask-image:radial-gradient(ellipse_70%_60%_at_50%_0%,#fff_60%,transparent)]"></div>
                <div class="absolute inset-x-0 top-0 flex justify-center">
                    <div class="h-[120px] w-[70%] bg-gradient-to-r from-cyan-400/40 via-blue-500/30 to-transparent blur-3xl opacity-60"></div>
                </div>
                <div class="absolute inset-x-0 top-4 flex justify-around text-[9px] font-mono uppercase tracking-[0.8em] text-cyan-200/40">
                    <span class="animate-pulse">cbad core live</span>
                    <span class="animate-pulse delay-300">ncd signal stable</span>
                    <span class="animate-pulse delay-500">audit log signed</span>
                </div>
                <div class="pointer-events-none absolute inset-0">
                    <div class="absolute left-1/4 top-0 h-full w-px bg-gradient-to-b from-transparent via-cyan-400/50 to-transparent animate-[pulse_6s_ease-in-out_infinite]"></div>
                    <div class="absolute right-1/3 top-0 h-full w-px bg-gradient-to-b from-transparent via-indigo-400/40 to-transparent animate-[pulse_8s_ease-in-out_infinite]"></div>
                </div>
            </div>

            <div class="container relative z-10 mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                <div class="mx-auto max-w-3xl text-center">
                    <div class="mb-6 inline-flex items-center gap-3 rounded-full border border-white/10 bg-white/5 px-4 py-2 text-xs font-semibold uppercase tracking-[0.4em] text-cyan-200">
                        Explainable · DORA Ready · Gemini Optional
                    </div>
                    <h1 class="text-5xl font-mono font-extrabold tracking-tight text-white sm:text-7xl">
                        Stop Black Box Audits.
                    </h1>
                    <p class="mt-6 text-lg leading-8 text-slate-200 sm:text-xl">
                        Driftlock compresses every payload twice, proves the math behind each anomaly, and (optionally) asks Gemini Flash to narrate the context.
                        Compliance teams get receipts. Engineers get sub-2s signal.
                    </p>
                    <div class="mt-10 flex flex-wrap items-center justify-center gap-4">
                        <a href="#signup" class="rounded-full bg-blue-600 px-8 py-3.5 text-sm font-semibold text-white shadow-lg shadow-blue-500/30 transition hover:-translate-y-0.5 hover:bg-blue-500">
                            Start Free Trial
                        </a>
                        <button type="button" @click="triggerLiveScan" class="rounded-full border border-white/30 px-6 py-3.5 text-sm font-semibold text-white transition hover:border-cyan-300 hover:text-cyan-200">
                            Live Signal Scan →
                        </button>
                        <a href="#api-demo" class="text-sm font-semibold leading-6 text-cyan-200 hover:text-white transition">
                            View API Demo
                        </a>
                    </div>
                </div>
            </div>
        </section>

        <!-- Feature Grid (Bento Style) -->
        <section id="solution" class="py-24 bg-white">
             <div class="container mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                <div class="mx-auto max-w-2xl text-center">
                    <h2 class="text-base font-semibold leading-7 text-blue-600">The Solution</h2>
                    <p class="mt-2 text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">Math, Not Magic.</p>
                    <p class="mt-6 text-lg leading-8 text-gray-600">
                        Driftlock replaces opaque, black-box models with compression-based anomaly detection. It's the "show your work" that auditors love.
                    </p>
                </div>

                <div class="mx-auto mt-16 max-w-2xl sm:mt-20 lg:mt-24 lg:max-w-none">
                    <dl class="grid max-w-xl grid-cols-1 gap-x-8 gap-y-16 lg:max-w-none lg:grid-cols-5">
                        <!-- Card 1 -->
                        <div class="flex flex-col rounded-2xl border border-gray-200 bg-white p-8 shadow-sm transition-shadow hover:shadow-lg">
                            <dt class="flex items-center gap-x-3 text-base font-semibold leading-7 text-gray-900">
                                <div class="h-10 w-10 flex items-center justify-center rounded-lg bg-blue-600">
                                    <svg class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                                    </svg>
                                </div>
                                Glass-Box Algorithm
                            </dt>
                            <dd class="mt-4 flex flex-auto flex-col text-base leading-7 text-gray-600">
                                <p class="flex-auto">Mathematical explanations for every decision. No black boxes, no "computer says no". Verify the math yourself.</p>
                            </dd>
                        </div>

                         <!-- Card 2 -->
                         <div class="flex flex-col rounded-2xl border border-gray-200 bg-white p-8 shadow-sm transition-shadow hover:shadow-lg">
                            <dt class="flex items-center gap-x-3 text-base font-semibold leading-7 text-gray-900">
                                <div class="h-10 w-10 flex items-center justify-center rounded-lg bg-blue-600">
                                    <svg class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z" />
                                    </svg>
                                </div>
                                Real-Time Detection
                            </dt>
                            <dd class="mt-4 flex flex-auto flex-col text-base leading-7 text-gray-600">
                                <p class="flex-auto">Process 1000+ events/second with sub-second latency. Built in Rust for maximum performance and safety.</p>
                            </dd>
                        </div>

                        <!-- Card 3 -->
                        <div class="flex flex-col rounded-2xl border border-gray-200 bg-white p-8 shadow-sm transition-shadow hover:shadow-lg">
                            <dt class="flex items-center gap-x-3 text-base font-semibold leading-7 text-gray-900">
                                <div class="h-10 w-10 flex items-center justify-center rounded-lg bg-blue-600">
                                    <svg class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
                                    </svg>
                                </div>
                                Audit-Ready Reports
                            </dt>
                            <dd class="mt-4 flex flex-auto flex-col text-base leading-7 text-gray-600">
                                <p class="flex-auto">One-click compliance reports for regulator submissions. Ready for EU DORA and US audit requirements.</p>
                            </dd>
                        </div>

                        <!-- Card 4 -->
                        <div class="flex flex-col rounded-2xl border border-gray-200 bg-white p-8 shadow-sm transition-shadow hover:shadow-lg">
                            <dt class="flex items-center gap-x-3 text-base font-semibold leading-7 text-gray-900">
                                <div class="h-10 w-10 flex items-center justify-center rounded-lg bg-blue-600">
                                    <svg class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 12l8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v9.75A1.5 1.5 0 006 21h12a1.5 1.5 0 001.5-1.5V9.75M8.25 21v-6.375a.375.375 0 01.375-.375h6.75a.375.375 0 01.375.375V21" />
                                    </svg>
                                </div>
                                Hybrid Explainability
                            </dt>
                            <dd class="mt-4 flex flex-auto flex-col text-base leading-7 text-gray-600">
                                <p class="flex-auto">LZ-based math filters 99.9% of traffic; Gemini Flash translates anomalies into human context for tickets, alerts, or AI ops reviews.</p>
                            </dd>
                        </div>

                        <!-- Card 5 -->
                        <div class="flex flex-col rounded-2xl border border-gray-200 bg-white p-8 shadow-sm transition-shadow hover:shadow-lg">
                            <dt class="flex items-center gap-x-3 text-base font-semibold leading-7 text-gray-900">
                                <div class="h-10 w-10 flex items-center justify-center rounded-lg bg-blue-600">
                                    <svg class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M11.42 15.88l-2.3 1.79a1 1 0 01-1.62-.78v-3.68L5 11.31a1 1 0 01.58-1.81l3.11-.25 1.2-3.42a1 1 0 011.88 0l1.2 3.42 3.11.25a1 1 0 01.58 1.81l-2.5 1.9v3.68a1 1 0 01-1.62.78l-2.12-1.65z" />
                                    </svg>
                                </div>
                                AI Agent Monitoring
                            </dt>
                            <dd class="mt-4 flex flex-auto flex-col text-base leading-7 text-gray-600">
                                <p class="flex-auto">Attach Driftlock to tool-using agents to spot runaway loops, hallucinations, or unexpected tool chains. Auto-trigger circuit breakers before spend explodes.</p>
                            </dd>
                        </div>
                    </dl>
                </div>
             </div>
        </section>

        <!-- Signup Section -->
        <section id="signup" class="py-24 bg-gray-50">
            <div class="container mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                <div class="grid lg:grid-cols-2 gap-16 items-center">
                    <div>
                        <h2 class="text-4xl font-mono font-bold text-gray-900 sm:text-5xl">Start Your Pilot</h2>
                        <p class="mt-6 text-lg text-gray-600">
                            Create your account and receive your API key instantly.
                            <br>
                            <span class="text-blue-600 font-medium">No credit card required for Developer plan.</span>
                        </p>
                        <ul class="mt-8 space-y-4 text-gray-600">
                            <li class="flex items-center gap-3">
                                <svg class="h-5 w-5 text-blue-500" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/></svg>
                                Instant API access
                            </li>
                            <li class="flex items-center gap-3">
                                <svg class="h-5 w-5 text-blue-500" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/></svg>
                                10k monthly events free
                            </li>
                            <li class="flex items-center gap-3">
                                <svg class="h-5 w-5 text-blue-500" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/></svg>
                                Full dashboard access
                            </li>
                        </ul>
                    </div>
                    <div class="bg-white rounded-2xl shadow-xl p-8 border border-gray-100">
                        <SignupForm />
                    </div>
                </div>
            </div>
        </section>

        <!-- API Demo Section -->
        <section id="api-demo" class="bg-slate-950 text-white py-24">
            <div class="container mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
                <div class="grid gap-12 lg:grid-cols-2 lg:items-center">
                    <div>
                        <div class="inline-flex items-center rounded-full bg-blue-400/10 px-3 py-1 text-sm font-medium text-blue-400 mb-6">
                            New Onboarding Flow
                        </div>
                        <h2 class="text-3xl font-mono font-bold sm:text-4xl">HTTP API Demo with Docker + Postgres</h2>
                        <p class="mt-4 text-lg text-slate-300">
                            This is the exact flow from the README quickstart. Run one script, provision a tenant/key, hit <code>/v1/detect</code>, and inspect anomalies in Postgres—no black-box magic.
                        </p>
                        
                        <div class="mt-8 space-y-6 border-l-2 border-slate-800 pl-6">
                             <div class="relative">
                                <h3 class="text-lg font-semibold text-white">1. Clone & Build</h3>
                                <p class="mt-2 text-slate-400 font-mono text-sm bg-slate-900 p-3 rounded-lg border border-slate-800">
                                    git clone ... && cargo build --release
                                </p>
                            </div>
                            <div class="relative">
                                <h3 class="text-lg font-semibold text-white">2. Run Demo Script</h3>
                                <p class="mt-2 text-slate-400 font-mono text-sm bg-slate-900 p-3 rounded-lg border border-slate-800">
                                    DRIFTLOCK_DEV_MODE=true ./scripts/run-api-demo.sh
                                </p>
                            </div>
                             <div class="relative">
                                <h3 class="text-lg font-semibold text-white">3. Verify in Postgres</h3>
                                <p class="mt-2 text-slate-400 text-sm">
                                    `psql` output shows anomalies persisted with NCD, p-values, and explanations—ready for compliance teams.
                                </p>
                            </div>
                        </div>
                    </div>
                    <div class="relative">
                        <div class="absolute -inset-1 bg-gradient-to-r from-blue-600 to-cyan-600 rounded-2xl blur opacity-30"></div>
                        <div class="relative bg-slate-900 rounded-2xl border border-slate-800 shadow-2xl overflow-hidden">
                            <div class="flex items-center border-b border-slate-800 bg-slate-950/50 px-4 py-3">
                                <div class="flex space-x-2">
                                    <div class="h-3 w-3 rounded-full bg-red-500"></div>
                                    <div class="h-3 w-3 rounded-full bg-yellow-500"></div>
                                    <div class="h-3 w-3 rounded-full bg-green-500"></div>
                                </div>
                                <div class="ml-4 text-xs text-slate-500 font-mono">bash — 80x24</div>
                            </div>
                            <div class="p-6 font-mono text-sm text-slate-300">
                                <div class="text-green-400">$ ./scripts/run-api-demo.sh</div>
                                <div class="text-slate-500"># Initializing Driftlock HTTP API...</div>
                                <div>> Service listening on :8080</div>
                                <div>> Provisioning tenant... OK</div>
                                <div>> Sending 50 normal events... OK</div>
                                <div>> Sending 1 anomaly... <span class="text-red-400">DETECTED</span></div>
                                <br>
                                <div class="text-blue-300">ANOMALY FOUND (ID: 550e8400-...)</div>
                                <div>NCD: 0.82 (Threshold: 0.5)</div>
                                <div>Explanation: "High entropy variance in field 'amount'"</div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <!-- Problem Section (The $ Risk) -->
        <section id="problem" class="py-24 bg-gray-50">
            <div class="container mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                <div class="mx-auto max-w-3xl text-center">
                    <h2 class="text-base font-semibold leading-7 text-blue-600">The Problem</h2>
                    <p class="mt-2 text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">Your Black Box is an Audit Failure.</p>
                </div>
                <!-- Regulation Cards -->
                <div class="mt-16 grid grid-cols-1 gap-8 md:grid-cols-3">
                    <div class="rounded-2xl bg-white p-8 shadow-sm border border-gray-100">
                        <h3 class="text-lg font-semibold text-gray-900">EU DORA</h3>
                        <p class="mt-4 text-3xl font-bold text-red-600">Up to 2%</p>
                        <p class="mt-2 text-sm text-gray-600">of annual global turnover in fines. Article 15 requires full explainability.</p>
                    </div>
                    <div class="rounded-2xl bg-white p-8 shadow-sm border border-gray-100">
                        <h3 class="text-lg font-semibold text-gray-900">US Regulations</h3>
                        <p class="mt-4 text-3xl font-bold text-red-600">CFPB + FFIEC</p>
                        <p class="mt-2 text-sm text-gray-600">Fair Lending and FFIEC guidelines mandate auditability. Black boxes violate this.</p>
                    </div>
                    <div class="rounded-2xl bg-white p-8 shadow-sm border border-gray-100">
                        <h3 class="text-lg font-semibold text-gray-900">The Gap</h3>
                        <p class="mt-4 text-3xl font-bold text-red-600">Auto Failure</p>
                        <p class="mt-2 text-sm text-gray-600">No mathematical proof = failed audit. Fines start January 2025.</p>
                    </div>
                </div>
            </div>
        </section>

        <!-- Live Playground Section -->
        <section id="playground" class="py-24">
            <div class="container mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                <div class="text-center mb-12">
                    <h2 class="text-3xl font-bold tracking-tight text-gray-900">Try It Yourself</h2>
                    <p class="mt-4 text-lg text-gray-600">Interact with the API directly from your browser.</p>
                </div>
                <PlaygroundShell ref="playgroundRef" variant="embedded" />
            </div>
        </section>

        <!-- Pricing Section -->
        <section id="pricing" class="py-24 bg-white border-t border-gray-200">
            <div class="container mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                <div class="text-center max-w-2xl mx-auto">
                    <p class="text-sm font-semibold uppercase tracking-[0.2em] text-blue-600">Pricing</p>
                    <h2 class="mt-4 text-4xl font-mono font-bold text-gray-900">Start free. Scale with AI explainability.</h2>
                    <p class="mt-4 text-base text-gray-600">Compression math costs about $1 per million events. Our plans price in the Gemini explainability layer so you don't have to think about tokens.</p>
                </div>

                <div class="mt-16 grid gap-8 lg:grid-cols-3">
                    <div class="rounded-3xl border border-gray-200 bg-white p-8 shadow-sm flex flex-col">
                        <div>
                            <p class="text-sm font-semibold uppercase tracking-wide text-blue-600">Developer</p>
                            <p class="mt-4 text-4xl font-bold text-gray-900">$0</p>
                            <p class="text-sm text-gray-500">per month</p>
                            <p class="mt-4 text-sm text-gray-600">10k events/mo, math proofs, CLI + API access.</p>
                        </div>
                        <ul class="mt-8 space-y-3 text-sm text-gray-600">
                            <li>✔ Deterministic anomaly detection</li>
                            <li>✔ CLI demo & Playground</li>
                            <li>✔ JSON reports</li>
                        </ul>
                        <a href="#signup" class="mt-10 inline-flex justify-center rounded-full border border-blue-200 bg-blue-50 px-4 py-2 text-sm font-semibold text-blue-600 hover:bg-blue-100">Start Building</a>
                    </div>

                    <div class="rounded-3xl border-2 border-blue-600 bg-slate-950 text-white p-8 shadow-xl flex flex-col relative overflow-hidden">
                        <span class="absolute right-6 top-6 rounded-full bg-blue-500 px-3 py-1 text-xs font-semibold uppercase tracking-wide">Most Popular</span>
                        <div>
                            <p class="text-sm font-semibold uppercase tracking-wide text-blue-300">Starter</p>
                            <p class="mt-4 text-4xl font-bold">$25</p>
                            <p class="text-sm text-slate-300">per month</p>
                            <p class="mt-4 text-sm text-slate-300">500k events/mo + Gemini summaries for every anomaly.</p>
                        </div>
                        <ul class="mt-8 space-y-3 text-sm text-slate-200">
                            <li>✔ Email alerts with plain-English context</li>
                            <li>✔ Driftlock dashboard + usage analytics</li>
                            <li>✔ API key rotation & audit logs</li>
                        </ul>
                        <a href="#signup" class="mt-10 inline-flex justify-center rounded-full bg-white px-4 py-2 text-sm font-semibold text-slate-900 hover:bg-slate-100">Launch Starter</a>
                    </div>

                    <div class="rounded-3xl border border-gray-200 bg-white p-8 shadow-sm flex flex-col">
                        <div>
                            <p class="text-sm font-semibold uppercase tracking-wide text-blue-600">Pro / Enterprise</p>
                            <p class="mt-4 text-4xl font-bold text-gray-900">Custom</p>
                            <p class="text-sm text-gray-500">5M+ events/mo</p>
                            <p class="mt-4 text-sm text-gray-600">Unlimited streams, AI agent monitoring, Slack/PagerDuty integrations.</p>
                        </div>
                        <ul class="mt-8 space-y-3 text-sm text-gray-600">
                            <li>✔ Dedicated environments & SOC2 pipeline</li>
                            <li>✔ Human + LLM explainability packs</li>
                            <li>✔ Compliance reports (DORA, FFIEC, NIS2)</li>
                        </ul>
                        <a href="mailto:hunter@shannonlabs.dev" class="mt-10 inline-flex justify-center rounded-full border border-gray-300 px-4 py-2 text-sm font-semibold text-gray-900 hover:bg-gray-50">Contact Sales</a>
                    </div>
                </div>
            </div>
        </section>

        <!-- ROI Calculator Section -->
        <section id="roi" class="py-24 bg-slate-950 text-white border-t border-slate-900">
            <div class="container mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                <div class="mx-auto max-w-3xl text-center">
                    <h2 class="text-base font-semibold leading-7 text-blue-400">Financial Impact</h2>
                    <p class="mt-2 text-3xl font-bold tracking-tight text-white sm:text-4xl">Rough ROI Estimation</p>
                    <p class="mt-4 text-base text-slate-400">
                        Uses assumptions from our README: ~$1 per million anomaly checks.
                    </p>
                </div>
                <div class="mt-12 grid gap-10 lg:grid-cols-2 lg:items-start">
                    <!-- Inputs -->
                    <div class="rounded-3xl bg-slate-900 p-8 shadow-2xl border border-slate-800">
                        <h3 class="text-lg font-semibold text-white mb-6">Your Assumptions</h3>
                        <div class="space-y-6">
                            <div>
                                <label class="flex justify-between text-sm font-medium text-slate-300">
                                    <span>Events per month</span>
                                    <span class="font-mono text-blue-300">{{ roiInputs.eventsPerMonth.toLocaleString() }}</span>
                                </label>
                                <input v-model.number="roiInputs.eventsPerMonth" type="range" min="1000000" max="100000000" step="1000000" class="mt-2 w-full accent-blue-500 h-2 bg-slate-700 rounded-lg appearance-none cursor-pointer" />
                            </div>
                            <div>
                                <label class="flex justify-between text-sm font-medium text-slate-300">
                                    <span>Annual risk exposure</span>
                                    <span class="font-mono text-blue-300">${{ roiInputs.annualLoss.toLocaleString() }}</span>
                                </label>
                                <input v-model.number="roiInputs.annualLoss" type="range" min="500000" max="50000000" step="500000" class="mt-2 w-full accent-blue-500 h-2 bg-slate-700 rounded-lg appearance-none cursor-pointer" />
                            </div>
                            <div>
                                <label class="flex justify-between text-sm font-medium text-slate-300">
                                    <span>Loss reduction</span>
                                    <span class="font-mono text-blue-300">{{ roiInputs.lossReductionPct }}%</span>
                                </label>
                                <input v-model.number="roiInputs.lossReductionPct" type="range" min="1" max="50" step="1" class="mt-2 w-full accent-blue-500 h-2 bg-slate-700 rounded-lg appearance-none cursor-pointer" />
                            </div>
                        </div>
                    </div>

                    <!-- Outputs -->
                    <div class="rounded-3xl bg-slate-900 p-8 shadow-2xl border border-blue-900/30">
                        <h3 class="text-lg font-semibold text-white mb-6">Projected Annual Impact</h3>
                        <div class="grid gap-4 sm:grid-cols-2">
                            <div class="rounded-2xl bg-slate-950 p-5 border border-slate-800">
                                <p class="text-xs font-medium uppercase text-slate-500">Est. Cost</p>
                                <p class="mt-2 text-2xl font-mono font-bold text-white">${{ roiOutputs.annualDriftlockSpend.toLocaleString() }}</p>
                            </div>
                             <div class="rounded-2xl bg-slate-950 p-5 border border-slate-800">
                                <p class="text-xs font-medium uppercase text-slate-500">Savings</p>
                                <p class="mt-2 text-2xl font-mono font-bold text-emerald-400">${{ roiOutputs.lossesAvoided.toLocaleString() }}</p>
                            </div>
                             <div class="rounded-2xl bg-slate-950 p-5 border border-slate-800 sm:col-span-2">
                                <p class="text-xs font-medium uppercase text-slate-500">Net Benefit</p>
                                <p class="mt-2 text-3xl font-mono font-bold text-blue-400">
                                    ${{ roiOutputs.netBenefit.toLocaleString() }}
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <!-- CTA Section -->
        <section id="contact" class="bg-white py-24 border-t border-gray-200">
             <div class="container mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 text-center">
                <h2 class="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">Ready to Verify?</h2>
                <p class="mt-6 text-lg leading-8 text-gray-600 max-w-2xl mx-auto">
                    Join the pilot program for Q1 2026. We are accepting 3 more partners.
                </p>
                <div class="mt-10 flex items-center justify-center gap-x-6">
                    <a href="mailto:hunter@shannonlabs.dev" class="rounded-md bg-blue-600 px-8 py-3.5 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600">
                        Contact Sales
                    </a>
                    <a href="#signup" class="text-sm font-semibold leading-6 text-gray-900">
                        Sign Up Free <span aria-hidden="true">→</span>
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
