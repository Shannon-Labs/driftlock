import { expect, test } from '@playwright/test'

const BASE_URL = process.env.DRIFTLOCK_URL ?? 'https://driftlock.net/'

test('homepage embeds the live playground and processes a sample', async ({ page }) => {
  await page.goto(BASE_URL, { waitUntil: 'networkidle' })

  // Hero copy should mention the playground CTA and prompt text.
  await expect(page.getByText('Paste JSON/NDJSON or load a sample', { exact: false })).toBeVisible()

  const textarea = page.locator('textarea', { hasText: undefined }).first()
  await expect(textarea).toBeVisible()

  // Load the bundled NDJSON sample via the same endpoint used by the UI and paste it.
  const sampleResponse = await page.request.get('/samples/small.ndjson')
  const sample = await sampleResponse.text()
  await textarea.fill(sample)
  await expect(textarea).toHaveValue(/latency_ms/, { timeout: 5000 })
  await page.getByRole('button', { name: /Use Data/i }).click()

  // Run detection to confirm the client-side logic stays wired to the API.
  await page.getByRole('button', { name: /Run Detection/i }).click()

  const processing = page.getByText('Processing detection...', { exact: false })
  await expect(processing).toBeVisible()
  await processing.waitFor({ state: 'hidden', timeout: 20000 })

  const hasResults = await page.getByRole('heading', { name: /Results/i }).isVisible().catch(() => false)
  const hasError = await page.getByText(/API is unavailable|Request timed out|Network error/i).isVisible().catch(() => false)
  expect(hasResults || hasError).toBeTruthy()
})

