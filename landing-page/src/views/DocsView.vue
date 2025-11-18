<template>
  <div class="min-h-screen bg-white">
    <!-- Navigation -->
    <nav class="border-b border-gray-200 bg-white sticky top-0 z-50">
      <div class="container mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex h-16 justify-between">
          <div class="flex">
            <div class="flex flex-shrink-0 items-center">
              <a href="/" class="text-xl font-mono font-bold text-gray-900">Driftlock Docs</a>
            </div>
            <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
              <a href="#quickstart" class="inline-flex items-center border-b-2 border-blue-500 px-1 pt-1 text-sm font-medium text-gray-900">Quick Start</a>
              <a href="#api" class="inline-flex items-center border-b-2 border-transparent px-1 pt-1 text-sm font-medium text-gray-500 hover:border-gray-300 hover:text-gray-700">API Reference</a>
              <a href="#agents" class="inline-flex items-center border-b-2 border-transparent px-1 pt-1 text-sm font-medium text-gray-500 hover:border-gray-300 hover:text-gray-700">AI Agents</a>
            </div>
          </div>
          <div class="flex items-center">
             <a href="/dashboard" class="text-sm font-semibold leading-6 text-gray-900">Dashboard <span aria-hidden="true">&rarr;</span></a>
          </div>
        </div>
      </div>
    </nav>

    <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <div class="mx-auto max-w-4xl">
        
        <section id="quickstart" class="mb-16">
          <h1 class="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl mb-6">Quick Start</h1>
          <p class="text-lg leading-8 text-gray-600 mb-8">
            Get up and running with Driftlock in less than 5 minutes.
          </p>

          <div class="bg-slate-900 rounded-xl p-6 mb-8 text-slate-200 font-mono text-sm">
            <div class="flex items-center justify-between mb-4 border-b border-slate-700 pb-2">
              <span>bash</span>
            </div>
            <p class="whitespace-pre-wrap"># 1. Install the CLI (Hypothetical)
npm install -g driftlock-cli

# 2. Login and get your key
driftlock login

# 3. Start monitoring
driftlock watch --url http://localhost:8000</p>
          </div>

          <p class="text-gray-600 mb-4">
            Alternatively, use the HTTP API directly:
          </p>
        </section>

        <section id="api" class="mb-16">
          <h2 class="text-2xl font-bold tracking-tight text-gray-900 mb-6">API Reference</h2>
          
          <h3 class="text-lg font-semibold text-gray-900 mb-4">POST /v1/detect</h3>
          <p class="text-gray-600 mb-4">Send a transaction or event log to be analyzed.</p>
          
          <div class="bg-slate-900 rounded-xl p-6 mb-8 text-slate-200 font-mono text-sm">
             <p>curl -X POST https://api.driftlock.net/v1/detect \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "event_type": "agent_tool_call",
    "payload": {
      "tool": "stripe_charge",
      "amount": 50000,
      "currency": "usd"
    }
  }'</p>
          </div>
        </section>

        <section id="agents" class="mb-16">
           <h2 class="text-2xl font-bold tracking-tight text-gray-900 mb-6">Integration with AI Agents</h2>
           <p class="text-gray-600 mb-4">
             Driftlock is designed to be the safety layer for autonomous agents built with LangChain, AutoGPT, or custom loops.
           </p>
           <p class="text-gray-600 mb-4">
             <strong>Why monitor agents?</strong> Agents can "drift" into repetitive loops, hallucinate arguments, or be hijacked. Driftlock catches this by analyzing the compression distance of the agent's output stream against its "normal" baseline.
           </p>
        </section>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// No script needed for static docs view
</script>

