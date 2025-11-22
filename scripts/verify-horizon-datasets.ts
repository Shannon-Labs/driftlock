#!/usr/bin/env tsx

interface HorizonSample {
  id: string
  path: string
  description: string
  expectJson: boolean
}

const SAMPLES: HorizonSample[] = [
  { id: 'fintech', path: '/samples/fraud.json', description: 'Financial Fraud dataset', expectJson: true },
  { id: 'defi', path: '/samples/terra.json', description: 'Terra/Luna collapse dataset', expectJson: true },
  { id: 'aviation', path: '/samples/airline.json', description: 'Airline telemetry dataset', expectJson: true },
  { id: 'ai-safety', path: '/samples/safety.json', description: 'Prompt injection dataset', expectJson: true },
  { id: 'network', path: '/samples/network.json', description: 'Network intrusion dataset', expectJson: true },
  { id: 'supply-chain', path: '/samples/supply.json', description: 'Supply chain dataset', expectJson: true },
  { id: 'cloud', path: '/samples/cloud.json', description: 'Cloud outage dataset', expectJson: true },
  { id: 'ndjson-demo', path: '/samples/demo-financial.ndjson', description: 'NDJSON financial sample', expectJson: false }
]

type Result = { sample: HorizonSample; ok: boolean; message: string }

async function main() {
  const base = normalizeBase(process.env.FIREBASE_HOSTING_URL || process.argv[2] || 'https://driftlock.net')
  const results: Result[] = []

  for (const sample of SAMPLES) {
    const url = `${base}${sample.path}`
    try {
      const res = await fetch(url, { headers: { 'Accept': 'application/json,text/plain' } })
      if (!res.ok) {
        results.push({ sample, ok: false, message: `HTTP ${res.status} ${res.statusText}` })
        continue
      }
      const body = await res.text()
      if (!body.trim()) {
        results.push({ sample, ok: false, message: 'Empty response body' })
        continue
      }

      if (sample.expectJson) {
        JSON.parse(body)
      } else {
        validateNdjson(body)
      }

      results.push({ sample, ok: true, message: `Fetched ${formatBytes(body.length)} from ${url}` })
    } catch (err) {
      results.push({ sample, ok: false, message: (err as Error).message })
    }
  }

  const failed = results.filter(r => !r.ok)
  console.log(`\nDriftlock Horizon dataset verification against ${base}`)
  for (const { sample, ok, message } of results) {
    const status = ok ? '✅' : '❌'
    console.log(`${status} [${sample.id}] ${sample.description}: ${message}`)
  }

  if (failed.length) {
    console.error(`\n${failed.length} dataset(s) failed validation.`)
    process.exit(1)
  }

  console.log('\nAll datasets validated successfully.')
}

function normalizeBase(value: string): string {
  if (!value) return 'https://driftlock.net'
  return value.replace(/\/$/, '')
}

function validateNdjson(payload: string) {
  const lines = payload.split(/\r?\n/).filter(Boolean)
  if (!lines.length) {
    throw new Error('NDJSON payload contained no lines')
  }
  for (const [index, line] of lines.entries()) {
    try {
      JSON.parse(line)
      return
    } catch (err) {
      if (index === lines.length - 1) {
        throw new Error(`NDJSON failed to parse: ${(err as Error).message}`)
      }
    }
  }
}

function formatBytes(size: number) {
  if (size < 1024) return `${size} B`
  const kb = size / 1024
  if (kb < 1024) return `${kb.toFixed(1)} KB`
  const mb = kb / 1024
  return `${mb.toFixed(1)} MB`
}

main().catch(err => {
  console.error(err)
  process.exit(1)
})
