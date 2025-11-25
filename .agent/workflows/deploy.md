---
description: Deploy the application
---

# Deploy Workflow

1.  **Pre-Deployment Checks**
    -   Ensure all tests pass: `just test`
    -   Ensure build succeeds: `just build`

2.  **Deploy Landing Page & Functions**
    -   Command: `just deploy`
    -   This runs `firebase deploy`.
    -   Requires `firebase-tools` login or `FIREBASE_TOKEN`.

3.  **Deploy Landing Page Only**
    -   Command: `just deploy-landing`
    -   Runs `scripts/deploy-landing.sh`.

4.  **Verification**
    -   Visit the deployed URL and verify functionality.
