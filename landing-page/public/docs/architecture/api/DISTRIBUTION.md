# API Distribution & Client SDKs

Driftlock provides a REST API and real-time streaming endpoints. While the core product is the hosted API, we provide client SDKs to make integration easier for developers.

## Why No "PyPI Equivalent"?

Unlike a library that you `import` and run locally (like `pandas` or `requests`), Driftlock is a **SaaS Platform**. The heavy lifting (CBAD algorithm, state management) happens on our servers.

However, developers expect idiomatic clients. We distribute these via standard package managers:

*   **Python**: `pip install driftlock` (PyPI)
*   **Node.js**: `npm install driftlock` (npm)
*   **Go**: `go get github.com/driftlock/driftlock-go`

## Generating Clients

We use the **OpenAPI Specification** (`docs/api/openapi.yaml`) as the single source of truth. You can generate clients for any language using `openapi-generator`.

### Example: Generating a Python Client

```bash
# Install generator
npm install @openapitools/openapi-generator-cli -g

# Generate Python client
openapi-generator-cli generate \
  -i docs/api/openapi.yaml \
  -g python \
  -o clients/python \
  --packageName driftlock
```

## Interactive Documentation

For "easy to understand" usage, we recommend deploying a developer portal.

*   **Swagger UI**: Included in many API gateways.
*   **Scalar / Redoc**: Modern, beautiful API references.

### Recommended: Embed Scalar in Docs

You can embed the API reference directly in your Vue app using `@scalar/api-reference`.

```html
<script setup>
import { ApiReference } from '@scalar/api-reference'
</script>

<template>
  <ApiReference :spec="{ url: '/openapi.yaml' }" />
</template>
```

This gives you a "Stripe-like" documentation experience immediately.

