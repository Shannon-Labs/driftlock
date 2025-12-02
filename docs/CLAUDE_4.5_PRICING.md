# Claude 4.5 Pricing (Updated December 2024)

## Standard Pricing

### Opus 4.5
Most intelligent model for building agents and coding
- **Input**: $5 / MTok
- **Output**: $25 / MTok

### Sonnet 4.5
Optimal balance of intelligence, cost, and speed
- **Input**:
  - Prompts ≤ 200K tokens: $3 / MTok
  - Prompts > 200K tokens: $6 / MTok
- **Output**:
  - Prompts ≤ 200K tokens: $15 / MTok
  - Prompts > 200K tokens: $22.50 / MTok

### Haiku 4.5
Fastest, most cost-efficient model
- **Input**: $1 / MTok
- **Output**: $5 / MTok

## Prompt Caching Pricing

### Opus 4.5
- **Write**: $6.25 / MTok
- **Read**: $0.50 / MTok

### Sonnet 4.5
- **≤ 200K tokens**:
  - Write: $3.75 / MTok
  - Read: $0.30 / MTok
- **> 200K tokens**:
  - Write: $7.50 / MTok
  - Read: $0.60 / MTok

### Haiku 4.5
- **Write**: $1.25 / MTok
- **Read**: $0.10 / MTok

## Batch Processing
Save 50% with batch processing (learn more)

## Cost Calculations for Driftlock

Based on typical anomaly analysis (500 input tokens, 200 output tokens):

### Standard API Pricing (with 15% margin)
- **Haiku 4.5**:
  - Base: (0.5 × $1 + 0.2 × $5) = $1.50
  - With margin: $1.73
- **Sonnet 4.5**:
  - Base: (0.5 × $3 + 0.2 × $15) = $4.50
  - With margin: $5.18
- **Opus 4.5**:
  - Base: (0.5 × $5 + 0.2 × $25) = $7.50
  - With margin: $8.63

### Batch API Pricing (with 15% margin)
- **Haiku 4.5**:
  - Base: $1.50 × 0.5 = $0.75
  - With margin: $0.86
- **Sonnet 4.5**:
  - Base: $4.50 × 0.5 = $2.25
  - With margin: $2.59
- **Opus 4.5**:
  - Base: $7.50 × 0.5 = $3.75
  - With margin: $4.31

## Model Selection Strategy

1. **Free Tier**: No AI access
2. **Radar ($5/mo)**: Haiku 4.5 with batch processing
3. **Lock ($15/mo)**: Sonnet 4.5 with batch processing
4. **Orbit ($50/mo)**: Opus 4.5 with batch processing

## Cost Optimization

1. **Batch Processing**: 50% discount on all API calls
2. **Prompt Caching**: Cache frequent patterns to reduce costs
3. **Smart Routing**: Only analyze high-confidence anomalies (1-3% of events)
4. **Token Optimization**: Minimize prompt size while maintaining accuracy

Last updated: December 1, 2024