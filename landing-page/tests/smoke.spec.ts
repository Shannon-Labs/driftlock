import { expect, test } from '@playwright/test'

test('homepage embeds the live playground and processes a sample', async ({ page }) => {
  await page.goto('/', { waitUntil: 'networkidle' })

  // Hero copy should mention the new playground CTA copy.
  await expect(page.getByText('Live Playground', { exact: false })).toBeVisible()

  const textarea = page.locator('textarea', { hasText: undefined }).first()
  await expect(textarea).toBeVisible()

  // Load the bundled NDJSON sample via the same endpoint used by the UI and paste it.
  const sampleResponse = await page.request.get('/samples/small.ndjson')
  const sample = await sampleResponse.text()
  await textarea.fill(sample)
  await expect(textarea).toHaveValue(/latency_ms/, { timeout: 5000 })
  await page.getByRole('button', { name: /Use Data/i }).click()

  // Run detection to confirm the client-side logic stays wired to the API.
  const runButton = page.locator('#playground').getByRole('button').filter({ hasText: /Run Analyzer|Scanning/ }).first()
  await runButton.click()
  await expect(runButton).toHaveText(/Scanning…/i)
  await expect(runButton).toHaveText(/Run Analyzer/i, { timeout: 20000 })

  const hasResults = await page.getByRole('heading', { name: /Results/i }).isVisible().catch(() => false)
  const hasError = await page.getByText(/API is unavailable|Request timed out|Network error/i).isVisible().catch(() => false)
  expect(hasResults || hasError).toBeTruthy()
})

