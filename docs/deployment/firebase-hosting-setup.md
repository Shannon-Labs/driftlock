# Firebase Hosting Setup

This guide covers deploying the Driftlock frontend (Landing Page + Dashboard) to Firebase Hosting.

## Prerequisites

1.  **Firebase CLI**: `npm install -g firebase-tools`
2.  **Firebase Login**: `firebase login`
3.  **Firebase Project**: Created in Firebase Console.

## Initialization

If setting up from scratch (already done in this repo):

```bash
firebase init hosting
```

Select the project and configure:
- **Public directory**: `landing-page/dist`
- **Configure as single-page app**: Yes
- **Set up automatic builds and deploys with GitHub**: Optional

## Deployment

The `deploy.sh` script handles building the Vue app and deploying it.

```bash
./deploy.sh
```

## Connecting to Backend

To connect the frontend to the Cloud Run backend, you can configure rewrites in `firebase.json` or use the backend URL directly.

### Option 1: Rewrites (Recommended)

Add a rewrite rule in `firebase.json` to proxy `/api` requests to your Cloud Run service.

```json
"rewrites": [
  {
    "source": "/api/**",
    "run": {
      "serviceId": "driftlock-api",
      "region": "us-central1"
    }
  },
  {
    "source": "**",
    "destination": "/index.html"
  }
]
```

*Note: This requires the Firebase project to be on the Blaze (pay-as-you-go) plan.*

### Option 2: Direct URL

If not using rewrites, ensure `VITE_API_URL` in the frontend is set to the Cloud Run URL.

