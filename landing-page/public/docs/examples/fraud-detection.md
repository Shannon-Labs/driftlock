# Fraud Detection

Learn how to use Driftlock to detect fraudulent transactions in real-time.

## Scenario

You are running an e-commerce platform. You want to flag transactions that deviate significantly from a user's normal purchasing behavior (e.g., unusually high amount, different currency, strange location).

## Implementation

### 1. Define the Event Structure

We will track `transaction` events.

```json
{
  "user_id": "u_123",
  "amount": 150.00,
  "currency": "USD",
  "merchant_category": "electronics",
  "ip_country": "US",
  "device_id": "d_xyz"
}
```

### 2. Integration Code (Node.js)

We'll use a separate stream for each user (`user-{id}`) to build a personalized baseline.

```javascript
const { DriftlockClient } = require('@driftlock/client');
const client = new DriftlockClient({ apiKey: process.env.API_KEY });

async function processTransaction(transaction) {
  // 1. Detect anomalies against the user's history
  const result = await client.detect({
    streamId: `user-${transaction.user_id}`,
    events: [{
      type: 'transaction',
      body: {
        amount: transaction.amount,
        currency: transaction.currency,
        merchant: transaction.merchant_category,
        country: transaction.ip_country
      }
    }]
  });

  // 2. Check for fraud
  const anomaly = result.anomalies[0];
  
  if (anomaly && anomaly.metrics.confidence > 0.9) {
    console.warn(`Fraud Alert for User ${transaction.user_id}:`, anomaly.why);
    return { status: 'blocked', reason: anomaly.why };
  }

  // 3. Proceed with payment
  return await paymentGateway.charge(transaction);
}
```

## Testing the Model

### Step 1: Train (Establish Baseline)

Send 10-20 "normal" transactions for a user.

```javascript
// Normal behavior: Small purchases in US
for (let i = 0; i < 20; i++) {
  await processTransaction({
    user_id: 'john_doe',
    amount: 20 + Math.random() * 30, // $20-$50
    currency: 'USD',
    merchant_category: 'retail',
    ip_country: 'US'
  });
}
```

### Step 2: Attack (Simulate Fraud)

Send a transaction that breaks the pattern.

```javascript
// Anomaly: Large purchase in different country
const result = await processTransaction({
  user_id: 'john_doe',
  amount: 2500.00, // High amount
  currency: 'EUR', // Different currency
  merchant_category: 'jewelry',
  ip_country: 'FR' // Different location
});

console.log(result); 
// Output: { status: 'blocked', reason: 'Significant deviation...' }
```

## Why Driftlock?

Traditional rule-based systems require you to write thousands of rules ("IF amount > 1000 AND country != US...").

Driftlock learns these patterns automatically for every single user. If a user normally spends $5000, a $2500 transaction won't be flagged. If they normally spend $10, it will be.
