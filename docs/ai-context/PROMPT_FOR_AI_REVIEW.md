# CRITICAL: Review and Fix Driftlock AI Integration for Launch

## Mission
Complete the Driftlock AI integration with Claude SDK and fix all issues to prepare for launch. Previous work has critical errors that need correction.

## IMMEDIATE FIXES NEEDED

### 1. Claude Model IDs are WRONG
- Current IDs: `claude-3-haiku-20240307`, `claude-3-sonnet-20240229`, `claude-3-opus-20240229`
- These are Claude 3.x models, NOT 4.5!
- Need to find and update ALL references to use Claude 4.5 model IDs
- Search Anthropic documentation for correct Claude 4.5 model IDs

### 2. Pricing Verification Required
The pricing in the plan file may be incorrect. You MUST:
1. Check the current Claude 4.5 pricing from official documentation
2. Verify all cost calculations in the codebase
3. Ensure the "Costco of anomaly detection" philosophy is maintained
4. Confirm 15% margin is correctly applied

### 3. Review Against Requirements
Go to `/Users/hunterbown/.claude/plans/sleepy-spinning-lynx.md` and verify:
- All implementation requirements are met
- Pricing aligns with the plan
- No deviations from the specified architecture

## TASKS TO COMPLETE

### Phase 1: Fix Critical Errors (Priority: CRITICAL)
1. Find correct Claude 4.5 model IDs from official Anthropic documentation
2. Update ALL model references in:
   - Database migrations (`api/migrations/20251201000000_ai_cost_control.sql`)
   - Smart router (`collector-processor/internal/ai/smart_router.go`)
   - Config system (`collector-processor/internal/ai/config.go`)
   - API handlers (`collector-processor/cmd/driftlock-http/ai_usage.go`)
   - Any other files with model references

### Phase 2: Verify and Fix Pricing
1. Check `/Users/hunterbown/.claude/plans/pricing-prompt.md` for original requirements
2. Verify pricing in `/docs/CLAUDE_4.5_PRICING.md` against current official pricing
3. Update all cost calculations to use correct Claude 4.5 pricing
4. Ensure batch processing discounts (50%) are correctly implemented
5. Verify the 15% margin is transparently shown to users

### Phase 3: Complete Integration
1. Add missing API endpoint registrations in main.go
2. Implement actual Claude SDK client (currently only structure exists)
3. Add webhook handlers for AI usage tracking
4. Connect AI routing to actual anomaly detection pipeline
5. Test end-to-end flow

### Phase 4: Deployment Preparation
1. Create Docker configurations for AI services
2. Add environment variables for Claude API keys
3. Update deployment scripts
4. Create monitoring and alerting for AI costs
5. Write integration tests

## FILES TO REVIEW AND POTENTIALLY FIX

### Core Implementation
- `collector-processor/internal/ai/` - All AI-related code
- `collector-processor/cmd/driftlock-http/ai_usage.go`
- `collector-processor/cmd/driftlock-http/main.go` - Add route registrations
- `api/migrations/20251201000000_ai_cost_control.sql`

### Pricing and Plans
- `/Users/hunterbown/.claude/plans/sleepy-spinning-lynx.md`
- `/Users/hunterbown/.claude/plans/pricing-prompt.md`
- `/docs/CLAUDE_4.5_PRICING.md`
- Update CLAUDE.md with new AI features

### Frontend
- `landing-page/src/components/dashboard/AIUsageWidget.vue`
- Add AI configuration to settings
- Update plan descriptions to reflect AI features

### Claude Code Plugin
- `.claude/plugins/driftlock/` - Complete implementation
- Add actual API integration
- Test plugin functionality

## CRITICAL REQUIREMENTS

1. **DO NOT DEVIATE FROM THE PLAN**: Check pricing-prompt.md and sleepy-spinning-lynx.md
2. **FIX MODEL IDs**: Must use Claude 4.5, not 3.x
3. **MAINTAIN MARGINS**: 15% fixed margin, transparent pricing
4. **COST CONTROLS**: Auto-deletion for free tier, configurable limits
5. **INTELLIGENT ROUTING**: Only analyze 1-3% of events
6. **BATCH PROCESSING**: 50% discount must be implemented

## TESTING REQUIRED

1. Unit tests for cost calculations
2. Integration tests for AI routing
3. Load tests with cost tracking
4. End-to-end tests from event to AI analysis
5. Verify free tier auto-deletion works

## DELIVERABLES

1. All model IDs corrected to Claude 4.5
2. Accurate pricing implementation
3. Working Claude SDK integration
4. Complete deployment configuration
5. Test suite passing
6. Documentation updated

## NEXT STEPS AFTER COMPLETION

1. Run the database migration
2. Test AI features in staging
3. Monitor costs carefully
4. Prepare launch checklist
5. Deploy to production

## WARNING

The current implementation has WRONG model IDs which will cause API failures. This must be fixed FIRST before any testing or deployment.

Also verify that the pricing maintains accessibility ($5 entry point) while ensuring profitability. The "Costco of anomaly detection" model must be preserved - tiny margins on infrastructure + AI costs passed through transparently.

Please be thorough and check ALL references. The future of Driftlock depends on getting this right!