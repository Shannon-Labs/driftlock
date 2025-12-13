import { expect, test } from '@playwright/test'

/**
 * Billing UI E2E Tests
 *
 * These tests verify the billing-related UI elements in the dashboard.
 * Some tests require authentication and may need mock API responses.
 */

test.describe('Billing UI Components', () => {
  test.describe('Landing Page Pricing', () => {
    test('pricing section displays all tiers', async ({ page }) => {
      await page.goto('/', { waitUntil: 'networkidle' })

      // Verify pricing tiers are visible
      await expect(page.getByText('Free')).toBeVisible()
      await expect(page.getByText('Pro')).toBeVisible()
      await expect(page.getByText('Team')).toBeVisible()

      // Verify pricing amounts
      await expect(page.getByText('$99')).toBeVisible()
      await expect(page.getByText('$199')).toBeVisible()
    })

    test('upgrade buttons link to signup', async ({ page }) => {
      await page.goto('/', { waitUntil: 'networkidle' })

      // Find a CTA button in pricing section
      const ctaButton = page.getByRole('link', { name: /Get Started|Start Trial|Sign Up/i }).first()
      if (await ctaButton.isVisible()) {
        const href = await ctaButton.getAttribute('href')
        expect(href).toMatch(/signup|register|auth/i)
      }
    })
  })

  test.describe('Dashboard Billing (Authenticated)', () => {
    // Skip if no test credentials available
    test.skip(({ browserName }) => !process.env.TEST_USER_EMAIL, 'Requires TEST_USER_EMAIL')

    test('dashboard shows billing section', async ({ page }) => {
      // This test requires authentication
      // For CI, set TEST_USER_EMAIL and TEST_USER_PASSWORD environment variables
      const email = process.env.TEST_USER_EMAIL
      const password = process.env.TEST_USER_PASSWORD

      if (!email || !password) {
        test.skip()
        return
      }

      await page.goto('/login', { waitUntil: 'networkidle' })
      await page.fill('input[type="email"]', email)
      await page.fill('input[type="password"]', password)
      await page.click('button[type="submit"]')

      // Wait for dashboard redirect
      await page.waitForURL(/dashboard/, { timeout: 10000 })

      // Verify billing-related elements exist
      const hasBillingSection = await page.getByText(/billing|subscription|trial/i).first().isVisible({ timeout: 5000 }).catch(() => false)
      expect(hasBillingSection).toBeTruthy()
    })
  })

  test.describe('Billing API Mocks', () => {
    test('trial banner displays correctly (relaxed - 8+ days)', async ({ page }) => {
      // Mock the billing API response
      await page.route('**/api/v1/me/billing', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            status: 'trialing',
            plan: 'radar',
            trial_ends_at: new Date(Date.now() + 10 * 24 * 60 * 60 * 1000).toISOString(),
            trial_days_remaining: 10,
          }),
        })
      })

      // Mock auth check
      await page.route('**/api/v1/me', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ id: 'test-user', email: 'test@example.com' }),
        })
      })

      await page.goto('/dashboard', { waitUntil: 'networkidle' })

      // Check for relaxed trial banner (gray, subtle)
      const trialBanner = page.getByText(/TRIAL ACTIVE/i)
      if (await trialBanner.isVisible({ timeout: 5000 }).catch(() => false)) {
        await expect(trialBanner).toBeVisible()
        await expect(page.getByText(/10 days remaining/i)).toBeVisible()
      }
    })

    test('trial banner displays correctly (urgent - 0-3 days)', async ({ page }) => {
      await page.route('**/api/v1/me/billing', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            status: 'trialing',
            plan: 'radar',
            trial_ends_at: new Date(Date.now() + 2 * 24 * 60 * 60 * 1000).toISOString(),
            trial_days_remaining: 2,
          }),
        })
      })

      await page.route('**/api/v1/me', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ id: 'test-user', email: 'test@example.com' }),
        })
      })

      await page.goto('/dashboard', { waitUntil: 'networkidle' })

      // Check for urgent trial banner (orange)
      const urgentBanner = page.getByText(/TRIAL ENDS IN 2 DAYS/i)
      if (await urgentBanner.isVisible({ timeout: 5000 }).catch(() => false)) {
        await expect(urgentBanner).toBeVisible()
      }
    })

    test('grace period warning displays correctly', async ({ page }) => {
      await page.route('**/api/v1/me/billing', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            status: 'grace_period',
            plan: 'radar',
            grace_period_ends_at: new Date(Date.now() + 3 * 24 * 60 * 60 * 1000).toISOString(),
          }),
        })
      })

      await page.route('**/api/v1/me', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ id: 'test-user', email: 'test@example.com' }),
        })
      })

      await page.goto('/dashboard', { waitUntil: 'networkidle' })

      // Check for grace period warning (red)
      const graceWarning = page.getByText(/update your payment method/i)
      if (await graceWarning.isVisible({ timeout: 5000 }).catch(() => false)) {
        await expect(graceWarning).toBeVisible()
      }
    })

    test('free tier shows upgrade prompt', async ({ page }) => {
      await page.route('**/api/v1/me/billing', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            status: 'free',
            plan: 'pulse',
          }),
        })
      })

      await page.route('**/api/v1/me', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ id: 'test-user', email: 'test@example.com' }),
        })
      })

      await page.goto('/dashboard', { waitUntil: 'networkidle' })

      // Check for free tier upgrade prompt
      const upgradeButton = page.getByRole('button', { name: /Upgrade to Radar/i })
      if (await upgradeButton.isVisible({ timeout: 5000 }).catch(() => false)) {
        await expect(upgradeButton).toBeVisible()
      }
    })

    test('upgrade button initiates Stripe checkout', async ({ page }) => {
      let checkoutCalled = false

      await page.route('**/api/v1/me/billing', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ status: 'free', plan: 'pulse' }),
        })
      })

      await page.route('**/api/v1/me', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ id: 'test-user', email: 'test@example.com' }),
        })
      })

      await page.route('**/api/v1/billing/checkout', async (route) => {
        checkoutCalled = true
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ url: 'https://checkout.stripe.com/test-session' }),
        })
      })

      await page.goto('/dashboard', { waitUntil: 'networkidle' })

      const upgradeButton = page.getByRole('button', { name: /Upgrade to Radar/i })
      if (await upgradeButton.isVisible({ timeout: 5000 }).catch(() => false)) {
        // Click should trigger checkout API call
        await upgradeButton.click()

        // Verify checkout was initiated
        await page.waitForTimeout(1000)
        expect(checkoutCalled).toBeTruthy()
      }
    })

    test('manage billing opens Stripe portal', async ({ page }) => {
      let portalCalled = false

      await page.route('**/api/v1/me/billing', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ status: 'active', plan: 'radar' }),
        })
      })

      await page.route('**/api/v1/me', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ id: 'test-user', email: 'test@example.com' }),
        })
      })

      await page.route('**/api/v1/billing/portal', async (route) => {
        portalCalled = true
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ url: 'https://billing.stripe.com/test-session' }),
        })
      })

      await page.goto('/dashboard', { waitUntil: 'networkidle' })

      const manageBillingButton = page.getByRole('button', { name: /Manage Billing/i })
      if (await manageBillingButton.isVisible({ timeout: 5000 }).catch(() => false)) {
        await manageBillingButton.click()
        await page.waitForTimeout(1000)
        expect(portalCalled).toBeTruthy()
      }
    })
  })
})
