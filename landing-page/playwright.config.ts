import { defineConfig, devices } from '@playwright/test'

const localBaseURL = 'http://127.0.0.1:4173'
const targetBaseURL = process.env.DRIFTLOCK_URL ?? localBaseURL
const shouldStartPreview = !process.env.DRIFTLOCK_URL

export default defineConfig({
  testDir: './tests',
  timeout: 60 * 1000,
  expect: {
    timeout: 20 * 1000,
  },
  outputDir: 'tests/artifacts',
  use: {
    baseURL: targetBaseURL,
    trace: 'on-first-retry',
  },
  webServer: shouldStartPreview
    ? {
        command: 'npm run preview -- --host 127.0.0.1 --port 4173',
        port: 4173,
        reuseExistingServer: !process.env.CI,
        timeout: 120 * 1000,
      }
    : undefined,
  projects: [
    {
      name: 'chromium',
      use: {
        ...devices['Desktop Chrome'],
      },
    },
  ],
})

