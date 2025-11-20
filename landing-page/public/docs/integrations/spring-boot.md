# Spring Boot Integration

Integrate Driftlock anomaly detection into your Java Spring Boot application.

## Installation

Add the Driftlock dependency to your `pom.xml`:

```xml
<dependency>
    <groupId>io.driftlock</groupId>
    <artifactId>driftlock-spring-boot-starter</artifactId>
    <version>1.0.0</version>
</dependency>
```

Or `build.gradle`:

```gradle
implementation 'io.driftlock:driftlock-spring-boot-starter:1.0.0'
```

## Configuration

Configure Driftlock in your `application.properties` or `application.yml`.

```yaml
driftlock:
  api-key: ${DRIFTLOCK_API_KEY}
  stream-id: spring-app
  enabled: true
  async: true
```

## Automatic Request Monitoring

The starter automatically registers a `Filter` that intercepts all HTTP requests and sends metrics to Driftlock.

You can customize which paths to exclude:

```yaml
driftlock:
  exclude-paths:
    - /actuator/**
    - /health
    - /swagger-ui/**
```

## Manual Detection

Inject the `DriftlockClient` bean to perform manual detection.

```java
import io.driftlock.client.DriftlockClient;
import io.driftlock.model.DetectionRequest;
import io.driftlock.model.Event;
import org.springframework.stereotype.Service;

@Service
public class PaymentService {

    private final DriftlockClient driftlockClient;

    public PaymentService(DriftlockClient driftlockClient) {
        this.driftlockClient = driftlockClient;
    }

    public void processPayment(Payment payment) {
        // Create event
        Event event = Event.builder()
                .type("payment")
                .body(payment)
                .build();

        // Detect anomalies
        var result = driftlockClient.detect("payments", List.of(event));

        if (!result.getAnomalies().isEmpty()) {
            log.warn("Anomaly detected: {}", result.getAnomalies());
        }

        // Process payment...
    }
}
```

## Aspect Oriented Programming (AOP)

Annotate methods with `@DriftlockMonitor` to automatically detect anomalies in method arguments or return values.

```java
import io.driftlock.annotation.DriftlockMonitor;

@Service
public class OrderService {

    @DriftlockMonitor(streamId = "orders")
    public Order createOrder(OrderRequest request) {
        // ...
        return order;
    }
}
```

## Async Processing

By default, the starter uses an asynchronous executor to send events to Driftlock to avoid blocking the main thread. You can configure the thread pool:

```yaml
driftlock:
  executor:
    core-pool-size: 2
    max-pool-size: 10
    queue-capacity: 500
```

## Error Handling

Define a custom `DriftlockErrorHandler` to handle API errors.

```java
@Component
public class CustomErrorHandler implements DriftlockErrorHandler {
    @Override
    public void handleError(Exception e) {
        log.error("Driftlock error", e);
    }
}
```

## Actuator Integration

Driftlock exposes metrics via Spring Boot Actuator at `/actuator/metrics/driftlock.*`.

- `driftlock.requests.total`
- `driftlock.anomalies.detected`
- `driftlock.latency`
