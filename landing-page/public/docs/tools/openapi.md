# OpenAPI Specification

Driftlock provides a complete OpenAPI 3.0 (formerly Swagger) specification for our REST API.

## Specification File

You can download the latest specification file here:

- **[driftlock-openapi.yaml](/downloads/driftlock-openapi.yaml)** (YAML format)
- **[driftlock-openapi.json](/downloads/driftlock-openapi.json)** (JSON format)

## Interactive Documentation

We host an interactive Swagger UI where you can explore endpoints and test requests directly in your browser.

[Launch Swagger UI](https://driftlock-api-o6kjgrsowq-uc.a.run.app/docs)

## Code Generation

You can use the OpenAPI specification to automatically generate client libraries for your preferred language using tools like [OpenAPI Generator](https://openapi-generator.tech/).

### Example: Generating a Ruby Client

```bash
# Install OpenAPI Generator
npm install @openapitools/openapi-generator-cli -g

# Generate Client
openapi-generator-cli generate \
  -i https://driftlock.web.app/downloads/driftlock-openapi.yaml \
  -g ruby \
  -o ./driftlock-ruby-client
```

### Supported Languages

OpenAPI Generator supports over 50 languages, including:
- Rust
- PHP
- C# / .NET
- Swift
- Kotlin
- Dart

## Schema Validation

You can also use the specification to validate requests and responses in your own testing pipelines.

```javascript
const OpenApiValidator = require('express-openapi-validator');

app.use(
  OpenApiValidator.middleware({
    apiSpec: './driftlock-openapi.yaml',
    validateRequests: true,
    validateResponses: true,
  }),
);
```
