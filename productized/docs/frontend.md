# DriftLock Frontend Documentation

This document provides information about the DriftLock React frontend application.

## Overview

The DriftLock frontend is built with React and Material-UI, providing a modern dashboard for monitoring anomalies, managing subscriptions, and configuring settings.

## Architecture

```
Public Directory (Static Assets)
    |
React App (src/)
    |
├── Components (Reusable UI components)
├── Pages (Route-specific views)
├── API (API calls and integrations)
├── Contexts (Global state management)
├── Utils (Helper functions)
└── Styles (CSS and theme files)
```

## Technologies

- React 18
- React Router v6
- Material-UI (MUI) v5
- Axios for API requests
- Styled Components
- Stripe.js for payment integration
- Recharts for data visualization

## Setup

### Prerequisites
- Node.js 18+
- npm or yarn

### Installation
```bash
cd frontend
npm install
```

### Environment Variables

Create a `.env` file in the frontend directory:

```env
REACT_APP_API_URL=http://localhost:8080/api/v1
REACT_APP_STRIPE_PUBLISHABLE_KEY=pk_test_...
REACT_APP_CLOUDFLARE_SITE_KEY=your_cloudflare_site_key
```

### Running Development Server
```bash
npm start
```

## Project Structure

```
frontend/
├── public/             # Static assets
├── src/
│   ├── api/            # API service files
│   ├── components/     # Reusable UI components
│   ├── contexts/       # React context providers
│   ├── pages/          # Route-specific components
│   ├── styles/         # CSS and theme files
│   ├── utils/          # Helper functions
│   ├── App.js          # Main application component
│   └── index.js        # Application entry point
├── package.json        # Dependencies and scripts
└── .env                # Environment variables
```

## Components

### Layout Components
- `Layout.js` - Main layout with navigation
- `PrivateRoute.js` - Route protection wrapper

### Contexts
- `AuthContext.js` - User authentication state
- `AnomalyContext.js` - Anomaly data management

### Pages
- `Dashboard.js` - Main dashboard view
- `Anomalies.js` - Anomaly list and details
- `Billing.js` - Subscription management
- `Login.js` - Authentication login
- `Register.js` - User registration
- `Profile.js` - User profile management

## API Integration

### Base API Service
Located in `src/api/index.js`, handles:
- Base URL configuration
- Authentication headers
- Error handling
- Token refresh

### API Hooks
All API calls use the service files grouped by feature:
- `authAPI` - Authentication endpoints
- `anomalyAPI` - Anomaly detection endpoints
- `billingAPI` - Subscription and payment endpoints
- `emailAPI` - Email service endpoints
- `onboardingAPI` - Onboarding flow endpoints

## Styling

The application uses Material-UI for consistent styling with a custom theme defined in `App.js`. Component-specific styles are handled through MUI's `sx` prop and system.

## Routing

The application uses React Router with:
- Public routes (login, register)
- Protected routes (dashboard, billing, etc.)
- Nested routing within sections

## State Management

State is managed using:
- React Context for global state (auth, anomalies)
- Local component state for component-specific data
- React Router for navigation state

## Deployment

The frontend is deployed to Cloudflare Pages with:
- Environment-specific configuration
- Automatic builds from GitHub
- Custom domain and SSL
- CDN for asset delivery