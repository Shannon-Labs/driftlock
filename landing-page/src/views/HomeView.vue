<template>
    <main class="bg-background text-foreground font-sans selection:bg-black selection:text-white">
        <!-- Hero Section -->
        <section class="relative w-full overflow-hidden border-b border-black bg-white text-black lg:h-[90vh] min-h-[800px] flex flex-col">
            <!-- Brutalist Grid Background -->
            <div class="absolute inset-0 z-0 opacity-[0.03]"
                 style="background-image: linear-gradient(#000 1px, transparent 1px), linear-gradient(90deg, #000 1px, transparent 1px); background-size: 40px 40px;">
            </div>
            
            <div class="container-padding relative z-10 grid lg:grid-cols-2 gap-8 h-full py-12 lg:py-0">
                <!-- Left: Content -->
                <div class="flex flex-col justify-center z-20">
                    <div class="mb-8">
                        <div class="inline-flex items-center gap-2 border border-black px-3 py-1 mb-6 bg-black text-white shadow-[4px_4px_0px_0px_rgba(0,0,0,0.2)]">
                            <span class="w-2 h-2 bg-green-400 animate-pulse rounded-full"></span>
                            <span class="text-xs font-bold font-mono tracking-widest uppercase">Early Access</span>
                        </div>

                        <h1 class="text-6xl sm:text-7xl lg:text-7xl xl:text-8xl font-sans font-bold tracking-tight leading-none mb-8 uppercase break-words w-full">
                            DORA-Compliant<br>Transaction Monitoring.
                        </h1>

                        <p class="text-xl font-serif leading-relaxed text-gray-800 border-l-4 border-black pl-6 max-w-xl mb-10">
                            Detect financial anomalies with audit-ready explainability.
                            <strong>EU DORA Article 10 compliance built-in. Every incident includes reproducible metrics regulators can verify.</strong>
                        </p>

                        <div class="flex flex-col sm:flex-row gap-4 items-start sm:items-center">
                            <a href="#signup" class="brutalist-button-primary text-lg px-8 py-4">
                                Join Waitlist
                            </a>
                            <a href="/docs" class="text-sm font-mono underline underline-offset-4 hover:no-underline">
                                Read the docs →
                            </a>
                        </div>
                    </div>
                </div>

                <!-- Right: Entropy Compression Visualization -->
                <div class="hidden lg:flex flex-col justify-center lg:col-span-1 h-full z-10">
                    <div class="relative w-full max-w-2xl scale-75 origin-top lg:origin-center">
                        <!-- Entropy Compression Visualization -->
                        <svg viewBox="0 0 400 400" class="w-full h-auto">
                            <!-- Background frame -->
                            <rect x="20" y="20" width="360" height="360" fill="none" stroke="black" stroke-width="4"/>

                            <!-- Raw data stream (high entropy) -->
                            <g id="raw-data">
                                <text x="40" y="50" font-family="monospace" font-size="12" font-weight="bold" fill="black">RAW STREAM</text>
                                <!-- Random-looking data points -->
                                <rect x="40" y="60" width="8" height="8" fill="black"/>
                                <rect x="52" y="65" width="8" height="8" fill="black"/>
                                <rect x="64" y="58" width="8" height="8" fill="black"/>
                                <rect x="76" y="70" width="8" height="8" fill="black"/>
                                <rect x="88" y="62" width="8" height="8" fill="black"/>
                                <rect x="100" y="68" width="8" height="8" fill="black"/>
                                <rect x="112" y="60" width="8" height="8" fill="black"/>
                                <rect x="124" y="72" width="8" height="8" fill="black"/>
                                <rect x="136" y="65" width="8" height="8" fill="black"/>
                                <rect x="148" y="58" width="8" height="8" fill="black"/>
                                <rect x="160" y="70" width="8" height="8" fill="black"/>
                                <rect x="172" y="63" width="8" height="8" fill="black"/>
                                <rect x="184" y="69" width="8" height="8" fill="black"/>
                                <rect x="196" y="61" width="8" height="8" fill="black"/>
                                <rect x="208" y="73" width="8" height="8" fill="black"/>
                                <rect x="220" y="66" width="8" height="8" fill="black"/>
                                <rect x="232" y="59" width="8" height="8" fill="black"/>
                                <rect x="244" y="71" width="8" height="8" fill="black"/>
                                <rect x="256" y="64" width="8" height="8" fill="black"/>
                                <rect x="268" y="70" width="8" height="8" fill="black"/>
                                <rect x="280" y="62" width="8" height="8" fill="black"/>
                                <rect x="292" y="68" width="8" height="8" fill="black"/>
                                <rect x="304" y="60" width="8" height="8" fill="black"/>
                                <rect x="316" y="72" width="8" height="8" fill="black"/>
                                <rect x="328" y="65" width="8" height="8" fill="black"/>
                                <rect x="340" y="58" width="8" height="8" fill="black"/>
                                <rect x="352" y="70" width="8" height="8" fill="black"/>
                            </g>

                            <!-- Compression arrow -->
                            <g id="compress-arrow">
                                <path d="M 200 100 L 200 140" stroke="black" stroke-width="4" fill="none"/>
                                <path d="M 192 132 L 200 140 L 208 132" stroke="black" stroke-width="4" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
                                <text x="210" y="125" font-family="sans-serif" font-size="14" font-weight="bold" fill="black">COMPRESS</text>
                            </g>

                            <!-- Entropy detection box -->
                            <g id="entropy-box">
                                <rect x="80" y="160" width="240" height="80" fill="black" stroke="black" stroke-width="4"/>
                                <text x="200" y="190" font-family="sans-serif" font-size="16" font-weight="bold" fill="white" text-anchor="middle">ENTROPY ANALYSIS</text>
                                <!-- Entropy formula visualization -->
                                <text x="200" y="215" font-family="monospace" font-size="14" fill="white" text-anchor="middle">H = -Σp(x)log₂p(x)</text>
                                <text x="200" y="230" font-family="monospace" font-size="10" fill="#999" text-anchor="middle">Anomaly = High Entropy Resistance</text>
                            </g>

                            <!-- Detection arrow -->
                            <g id="detect-arrow">
                                <path d="M 200 250 L 200 290" stroke="black" stroke-width="4" fill="none"/>
                                <path d="M 192 282 L 200 290 L 208 282" stroke="black" stroke-width="4" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
                                <text x="210" y="275" font-family="sans-serif" font-size="14" font-weight="bold" fill="black">DETECT</text>
                            </g>

                            <!-- Anomaly signal (clean pattern) -->
                            <g id="anomaly-signal">
                                <text x="40" y="320" font-family="monospace" font-size="12" font-weight="bold" fill="black">ANOMALY DETECTED</text>
                                <!-- Clean, structured pattern for anomaly -->
                                <rect x="40" y="330" width="320" height="40" fill="none" stroke="black" stroke-width="2"/>
                                <!-- Anomaly pattern - regular with disruption -->
                                <line x1="50" y1="350" x2="150" y2="350" stroke="black" stroke-width="3"/>
                                <!-- Disruption (anomaly) -->
                                <rect x="150" y="340" width="40" height="20" fill="black"/>
                                <line x1="190" y1="350" x2="290" y2="350" stroke="black" stroke-width="3"/>
                                <!-- Signal indicators -->
                                <circle cx="170" cy="350" r="8" fill="white" stroke="black" stroke-width="2"/>
                                <text x="165" y="355" font-family="monospace" font-size="10" font-weight="bold" fill="black">!</text>
                            </g>

                            <!-- Compression ratio indicator -->
                            <g id="ratio-indicator">
                                <rect x="300" y="40" width="80" height="60" fill="none" stroke="black" stroke-width="2" stroke-dasharray="5,5"/>
                                <text x="340" y="60" font-family="monospace" font-size="10" font-weight="bold" fill="black" text-anchor="middle">RATIO</text>
                                <text x="340" y="75" font-family="monospace" font-size="16" font-weight="bold" fill="black" text-anchor="middle">2.7x</text>
                                <text x="340" y="90" font-family="monospace" font-size="8" fill="#666" text-anchor="middle">Low = Normal</text>
                                <text x="340" y="100" font-family="monospace" font-size="8" fill="#666" text-anchor="middle">High = Anomaly</text>
                            </g>
                        </svg>
                    </div>
                </div>
            </div>
        </section>

        <!-- How It Works -->
        <section class="section-padding bg-gray-50">
            <div class="container-padding">
                <div class="max-w-4xl mx-auto text-center mb-16">
                    <span class="text-xs font-bold font-mono uppercase tracking-widest mb-4 block text-gray-500">How It Works</span>
                    <h2 class="text-4xl md:text-5xl font-sans font-bold tracking-tight uppercase mb-8">
                        Three Steps to<br>Anomaly Detection.
                    </h2>
                    <p class="font-serif text-xl text-gray-600 max-w-3xl mx-auto">
                        Our compression-based approach detects anomalies without training data or ML models.
                    </p>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-3 gap-0 border-t border-l border-black max-w-5xl mx-auto">
                    <!-- Step 1 -->
                    <div class="border-b border-r border-black p-8 hover:bg-white transition-colors group">
                        <div class="mb-6 group-hover:-translate-y-1 transition-transform">
                            <!-- Data flow icon -->
                            <svg class="w-12 h-12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                                <path d="M4 12h8m0 0l-3-3m3 3l-3 3m1 6h10a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v5" stroke-linecap="round" stroke-linejoin="round"/>
                            </svg>
                        </div>
                        <div class="text-4xl font-mono font-bold text-gray-300 mb-4">01</div>
                        <h3 class="text-xl font-bold font-sans uppercase mb-4">Collect</h3>
                        <p class="font-serif text-gray-600 leading-relaxed">
                            Ingest OpenTelemetry data from your existing infrastructure. No changes to your current stack needed.
                        </p>
                    </div>

                    <!-- Step 2 -->
                    <div class="border-b border-r border-black p-8 hover:bg-white transition-colors group">
                        <div class="mb-6 group-hover:-translate-y-1 transition-transform">
                            <!-- Compression icon -->
                            <svg class="w-12 h-12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                                <rect x="3" y="3" width="18" height="7" rx="1"/>
                                <rect x="3" y="14" width="18" height="7" rx="1"/>
                                <path d="M8 7v7m8-7v7" stroke-linecap="round"/>
                                <path d="M11 10h2" stroke-linecap="round"/>
                                <path d="M11 18h2" stroke-linecap="round"/>
                            </svg>
                        </div>
                        <div class="text-4xl font-mono font-bold text-gray-300 mb-4">02</div>
                        <h3 class="text-xl font-bold font-sans uppercase mb-4">Compress</h3>
                        <p class="font-serif text-gray-600 leading-relaxed">
                            Apply entropy-based compression algorithms to detect anomalies. Unusual patterns resist compression.
                        </p>
                    </div>

                    <!-- Step 3 -->
                    <div class="border-b border-r border-black p-8 hover:bg-white transition-colors group">
                        <div class="mb-6 group-hover:-translate-y-1 transition-transform">
                            <!-- Alert icon -->
                            <svg class="w-12 h-12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                                <path d="M14.857 17.082a23.848 23.848 0 005.454-1.31A8.967 8.967 0 0118 9.75v-.7V9A6 6 0 006 9v.75a8.967 8.967 0 01-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 01-5.714 0m5.714 0a3 3 0 11-5.714 0" stroke-linecap="round" stroke-linejoin="round"/>
                            </svg>
                        </div>
                        <div class="text-4xl font-mono font-bold text-gray-300 mb-4">03</div>
                        <h3 class="text-xl font-bold font-sans uppercase mb-4">Alert</h3>
                        <p class="font-serif text-gray-600 leading-relaxed">
                            Receive alerts with mathematical proofs explaining exactly why each anomaly was detected.
                        </p>
                    </div>
                </div>

                <!-- Key Differentiator -->
                <div class="mt-16 text-center">
                    <div class="inline-block border border-black bg-white px-8 py-6 shadow-[4px_4px_0px_0px_rgba(0,0,0,0.1)]">
                        <p class="font-serif text-lg">
                            <span class="font-bold font-sans uppercase">No training period.</span>
                            <span class="text-gray-600"> Start detecting anomalies immediately with minimal baseline data.</span>
                        </p>
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

        <!-- DORA Compliance Section -->
        <section class="py-20 bg-black text-white border-b border-black">
            <div class="container-padding">
                <div class="grid lg:grid-cols-2 gap-16 items-center">
                    <div>
                        <h2 class="text-4xl font-sans font-bold tracking-tight uppercase mb-6">
                            Built for DORA.<br>Ready for Auditors.
                        </h2>
                        <p class="font-serif text-xl text-gray-400 max-w-md mb-8">
                            EU DORA (Regulation 2022/2554) requires financial entities to detect, report, and document ICT-related incidents. Driftlock makes compliance automatic.
                        </p>
                        <div class="space-y-4">
                            <div class="flex items-center gap-3">
                                <span class="text-green-400 font-mono">[x]</span>
                                <span class="text-sm">Article 10: Incident classification & detection</span>
                            </div>
                            <div class="flex items-center gap-3">
                                <span class="text-green-400 font-mono">[x]</span>
                                <span class="text-sm">Article 17: ICT risk management</span>
                            </div>
                            <div class="flex items-center gap-3">
                                <span class="text-green-400 font-mono">[x]</span>
                                <span class="text-sm">Article 19: 24h notification deadlines</span>
                            </div>
                            <div class="flex items-center gap-3">
                                <span class="text-green-400 font-mono">[x]</span>
                                <span class="text-sm">Complete audit trail (Driftlog)</span>
                            </div>
                        </div>
                    </div>

                    <div class="grid sm:grid-cols-2 gap-4">
                        <div class="border border-white/20 p-6 hover:border-white transition-colors">
                            <div class="font-mono text-xs uppercase tracking-widest text-gray-500 mb-2">For Compliance Officers</div>
                            <h3 class="font-bold text-lg mb-2">Incident Register</h3>
                            <p class="text-sm text-gray-400">DORA-compliant incident management with automated severity classification and notification tracking.</p>
                        </div>
                        <div class="border border-white/20 p-6 hover:border-white transition-colors">
                            <div class="font-mono text-xs uppercase tracking-widest text-gray-500 mb-2">For Risk Teams</div>
                            <h3 class="font-bold text-lg mb-2">Transaction Monitoring</h3>
                            <p class="text-sm text-gray-400">Real-time fraud detection for SEPA, SWIFT, and card transactions with mathematical proofs.</p>
                        </div>
                        <div class="border border-white/20 p-6 hover:border-white transition-colors">
                            <div class="font-mono text-xs uppercase tracking-widest text-gray-500 mb-2">For Auditors</div>
                            <h3 class="font-bold text-lg mb-2">Complete Audit Trail</h3>
                            <p class="text-sm text-gray-400">The Driftlog records every detection decision with metrics, thresholds, and timestamps.</p>
                        </div>
                        <div class="border border-white/20 p-6 hover:border-white transition-colors">
                            <div class="font-mono text-xs uppercase tracking-widest text-gray-500 mb-2">For Engineering</div>
                            <h3 class="font-bold text-lg mb-2">API-First Integration</h3>
                            <p class="text-sm text-gray-400">REST API for transaction analysis. Drop-in replacement for existing fraud scoring systems.</p>
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
                            <p class="text-sm text-gray-500 mt-2 font-serif">For testing & evaluation.</p>
                        </div>
                        <ul class="space-y-3 text-sm font-mono border-t border-gray-100 pt-4 mb-8">
                            <li class="flex gap-2"><span>[x]</span> 10k events/mo</li>
                            <li class="flex gap-2"><span>[x]</span> 5 streams</li>
                            <li class="flex gap-2"><span>[x]</span> 14-day Driftlog</li>
                        </ul>
                        <a href="#signup" class="block w-full border border-black py-2 text-center text-sm font-bold uppercase hover:bg-black hover:text-white transition-colors">Get Started</a>
                    </div>

                    <!-- Pro (Highlighted) -->
                    <div class="bg-black text-white border border-black p-8 shadow-[12px_12px_0px_0px_rgba(0,0,0,0.2)] relative lg:-mt-4 lg:mb-4 transform scale-105 z-10">
                        <div class="absolute top-0 right-0 bg-white text-black text-[10px] font-bold uppercase px-2 py-1">Most Popular</div>
                        <div class="mb-4">
                            <h3 class="font-bold font-sans uppercase text-lg text-white">Pro</h3>
                            <div class="text-3xl font-bold font-mono mt-2">$99<span class="text-sm font-normal text-gray-400">/mo</span></div>
                            <p class="text-sm text-gray-400 mt-2 font-serif">For production compliance.</p>
                        </div>
                        <ul class="space-y-3 text-sm font-mono border-t border-gray-800 pt-4 mb-8">
                            <li class="flex gap-2"><span class="text-green-400">[x]</span> 500k events/mo</li>
                            <li class="flex gap-2"><span class="text-green-400">[x]</span> 20 streams</li>
                            <li class="flex gap-2"><span class="text-green-400">[x]</span> 90-day Driftlog</li>
                            <li class="flex gap-2"><span class="text-green-400">[x]</span> CSV/PDF exports</li>
                        </ul>
                        <a href="#signup" class="block w-full bg-white text-black py-3 text-center text-sm font-bold uppercase hover:bg-gray-200 transition-colors">
                            Join Waitlist
                        </a>
                    </div>

                    <!-- Team -->
                    <div class="bg-white border border-gray-200 p-8 hover:border-black transition-colors">
                        <div class="mb-4">
                            <h3 class="font-bold font-sans uppercase text-lg">Team</h3>
                            <div class="text-3xl font-bold font-mono mt-2">$199<span class="text-sm font-normal text-gray-500">/mo</span></div>
                            <p class="text-sm text-gray-500 mt-2 font-serif">For growing teams.</p>
                        </div>
                        <ul class="space-y-3 text-sm font-mono border-t border-gray-100 pt-4 mb-8">
                            <li class="flex gap-2"><span>[x]</span> 5M events/mo</li>
                            <li class="flex gap-2"><span>[x]</span> 100 streams</li>
                            <li class="flex gap-2"><span>[x]</span> 1-year Driftlog</li>
                            <li class="flex gap-2"><span>[x]</span> Priority support</li>
                        </ul>
                        <a href="#signup" class="block w-full border border-black py-2 text-center text-sm font-bold uppercase hover:bg-black hover:text-white transition-colors">
                            Join Waitlist
                        </a>
                    </div>

                    <!-- Enterprise -->
                    <div class="bg-white border border-gray-200 p-8 hover:border-black transition-colors">
                        <div class="mb-4">
                            <h3 class="font-bold font-sans uppercase text-lg">Enterprise</h3>
                            <div class="text-3xl font-bold font-mono mt-2">Custom</div>
                            <p class="text-sm text-gray-500 mt-2 font-serif">For banks & regulated entities.</p>
                        </div>
                        <ul class="space-y-3 text-sm font-mono border-t border-gray-100 pt-4 mb-8">
                            <li class="flex gap-2"><span>[x]</span> Unlimited events</li>
                            <li class="flex gap-2"><span>[x]</span> 500+ streams</li>
                            <li class="flex gap-2"><span>[x]</span> Self-hosting option</li>
                            <li class="flex gap-2"><span>[x]</span> EU data residency</li>
                        </ul>
                        <a href="mailto:sales@driftlock.net" class="block w-full border border-black py-2 text-center text-sm font-bold uppercase hover:bg-black hover:text-white transition-colors">
                            Contact Sales
                        </a>
                    </div>
                </div>
            </div>
        </section>

        <!-- Waitlist / Footer CTA -->
        <section id="signup" class="section-padding bg-white border-b-0">
            <div class="container-padding max-w-2xl mx-auto text-center">
                <h2 class="text-4xl font-sans font-bold tracking-tight uppercase mb-6">Get Early Access</h2>
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
