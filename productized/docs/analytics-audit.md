# DriftLock Analytics and Audit Documentation

This document details the analytics and audit logging systems in DriftLock.

## Overview

DriftLock implements comprehensive analytics and audit logging to track user behavior, system performance, and compliance requirements.

## Analytics System

### Google Analytics 4 (GA4) Integration

The system integrates with GA4 for product analytics:

#### Events Tracked
- User actions (login, logout, dashboard views)
- API usage (endpoint calls, response times)
- Anomaly detection events (types, severities)
- Billing events (subscriptions, payments)
- Onboarding progress

#### Implementation
- Measurement Protocol for server-side events
- Client-side tracking for UI interactions
- Privacy-compliant data collection

### Custom Analytics Backend

In addition to GA4, the system maintains custom analytics:
- API usage metrics
- Anomaly detection performance
- User engagement metrics
- System health indicators

## Audit Logging System

### Purpose
- Compliance with regulations (SOC2, GDPR, etc.)
- Security monitoring
- Forensic analysis
- Change tracking

### Audit Log Schema

Each audit log contains:

```go
type AuditLog struct {
    ID          uint      // Unique identifier
    UserID      uint      // User who performed the action
    Action      string    // Type of action (login, create, update, delete)
    Resource    string    // Resource type (anomaly, user, billing, etc.)
    ResourceID  string    // Specific resource identifier
    IPAddress   string    // User's IP address
    UserAgent   string    // User agent string
    Details     string    // Additional details about the action
    CreatedAt   time.Time // Timestamp of the action
}
```

### Audited Actions

#### User Actions
- Login / Logout
- Profile updates
- Password changes

#### Anomaly Actions
- Anomaly creation
- Anomaly resolution
- Anomaly deletion
- Anomaly viewing

#### Billing Actions
- Subscription changes
- Payment processing
- Invoice generation
- Plan upgrades/downgrades

#### System Actions
- Configuration changes
- Administrative actions
- Data exports
- API key management

### Compliance Features

#### GDPR Compliance
- Right to erasure (for audit logs)
- Data retention policies (30 days default)
- Data portability
- Purpose limitation

#### SOC2 Compliance
- Detailed access logs
- System event logs
- Configuration change logs
- Regular audit reports

## Implementation

### Middleware Integration
Analytics and audit logging are implemented as middleware that:
- Runs on protected endpoints
- Captures user context
- Logs actions without impacting performance
- Handles failures gracefully

### Database Storage
Audit logs are stored in PostgreSQL with:
- Indexed fields for efficient querying
- Partitioning by date for large datasets
- Archival policies for old data
- Secure access controls

### Privacy Considerations
- Only necessary information is collected
- PII is handled according to privacy policies
- Data is encrypted at rest
- Secure deletion procedures

## API Endpoints

### Analytics Endpoints
- GET `/analytics/user-activity` - User activity metrics
- GET `/analytics/system-performance` - System performance metrics

### Audit Endpoints
- GET `/audit/logs` - Retrieve audit logs
- GET `/audit/user/:id` - User-specific audit logs
- GET `/audit/resource/:type/:id` - Resource-specific audit logs

## Monitoring and Alerting

### Analytics Monitoring
- Daily usage reports
- Anomaly detection metrics
- User engagement metrics
- API performance metrics

### Audit Monitoring
- Suspicious activity alerts
- Failed authentication attempts
- Unusual access patterns
- Compliance violations

## Data Retention

### Analytics Data
- Raw event data: 30 days
- Aggregated metrics: 1 year
- Anonymized data: 7 years

### Audit Logs
- Full audit logs: 1 year (compliance)
- Summary logs: 7 years (forensics)
- Automated deletion policies

## Performance Considerations

### Impact Minimization
- Asynchronous logging to avoid blocking requests
- Batch processing for high-volume events
- Optimized database writes
- Caching for frequently accessed data

### Scalability
- Horizontal scaling support
- Database connection pooling
- Queue-based processing
- Monitoring and alerting for performance issues