import { expect, test } from '@playwright/test'

const HORIZONS = [
  { button: /Financial Fraud/i, samplePath: '/samples/fraud.json' },
  { button: /Market Crash/i, samplePath: '/samples/terra.json' },
  { button: /Aviation Ops/i, samplePath: '/samples/airline.json' },
  { button: /AI Safety/i, samplePath: '/samples/safety.json' },
]

test.describe('Horizon Showcase datasets', () => {
  for (const horizon of HORIZONS) {
    test(`loads dataset for ${horizon.samplePath}`, async ({ page }) => {
      await page.goto('/', { waitUntil: 'networkidle' })

      const playgroundTextarea = page.locator('#playground textarea').first()
      await expect(playgroundTextarea).toBeVisible()

      const showcase = page.locator('#showcase')
      await showcase.scrollIntoViewIfNeeded()
      await expect(showcase).toBeVisible()

      await showcase.getByRole('button', { name: horizon.button }).click()
      await page.locator('#showcase').getByRole('button', { name: 'Load Data →' }).click()

      const sampleResponse = await page.request.get(horizon.samplePath)
      expect(sampleResponse.ok()).toBeTruthy()
      const sample = await sampleResponse.text()
      expect(sample).not.toBe('')

      const snippet = escapeForRegex(sample.slice(0, 200))
      await expect(playgroundTextarea).toHaveValue(new RegExp(snippet, 's'), { timeout: 15000 })
    })
  }
})

function escapeForRegex(value: string) {
  return value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}
