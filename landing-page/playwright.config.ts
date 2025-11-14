import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './tests',
  timeout: 60 * 1000,
  expect: {
    timeout: 20 * 1000,
  },
  use: {
    baseURL: process.env.DRIFTLOCK_URL ?? 'https://driftlock.net/',
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'chromium',
      use: {
        ...devices['Desktop Chrome'],
      },
    },
  ],
})

