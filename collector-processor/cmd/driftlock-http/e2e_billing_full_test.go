package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestE2E_FullBillingLifecycle tests the complete checkout flow:
// Signup -> Verify -> Checkout -> Webhook (Simulated) -> Entitlement
func TestE2E_FullBillingLifecycle(t *testing.T) {
	// Ensure we have a webhook secret for signature generation
	if os.Getenv("STRIPE_WEBHOOK_SECRET") == "" {
		os.Setenv("STRIPE_WEBHOOK_SECRET", "whsec_test_secret_12345")
	}

	te := setupTestEnv(t)
	defer te.Close()

	testEmail := generateTestEmail("lifecycle")
	defer te.cleanupTestTenants("%@e2e-test.driftlock.net")

	// 1. Signup
	t.Log("Step 1: Signing up new tenant...")
	resp, body := te.post("/v1/onboard/signup", map[string]interface{}{
		"email":        testEmail,
		"company_name": "Lifecycle Test Co",
		"plan":         "trial",
	}, nil)
	te.assertStatus(resp, http.StatusCreated)

	// 2. Verify Email & Get API Key
	t.Log("Step 2: Verifying email...")
	token, err := te.getVerificationToken(testEmail)
	if err != nil {
		t.Fatalf("Failed to get verification token: %v", err)
	}

	resp, body = te.get("/v1/onboard/verify?token="+token, nil)
	te.assertStatus(resp, http.StatusOK)

	var verifyResp map[string]interface{}
	json.Unmarshal(body, &verifyResp)
	apiKey := verifyResp["api_key"].(string)

	// 3. Verify Initial Plan (Trial/Pilot)
	t.Log("Step 3: Verifying initial plan...")
	stripeStatus, verified, tenantID := getTenantInfo(t, te, testEmail)
	if !verified {
		t.Error("Expected tenant to be verified")
	}
	if stripeStatus != "" {
		t.Logf("Initial stripe status: %s", stripeStatus)
	}
	t.Logf("Tenant ID: %s", tenantID)

	// 4. Initiate Checkout (Optional - primarily to ensure endpoint works)
	// We won't actually complete it on Stripe, but we'll simulate the webhook it WOULD send
	t.Log("Step 4: Initiating checkout for Pro tier...")
	// Only run this if Stripe keys are present, otherwise mock the anticipation
	if os.Getenv("STRIPE_SECRET_KEY") != "" {
		resp, body = te.post("/v1/billing/checkout", map[string]interface{}{
			"plan": "tensor", // Pro tier
		}, map[string]string{"X-Api-Key": apiKey})

		if resp.StatusCode == http.StatusOK {
			t.Log("Checkout session created successfully")
		} else {
			// If missing price IDs, it might 500, which is fine for this test
			// if we are just testing the webhook handler logic next.
			t.Log("Checkout creation skipped/failed (likely missing Stripe keys), proceeding to webhook simulation")
		}
	} else {
		t.Log("Skipping checkout initiation (no STRIPE_SECRET_KEY), proceeding to webhook simulation")
	}

	// 5. Simulate Webhook: checkout.session.completed
	t.Log("Step 5: Simulating Stripe Webhook...")

	// Construct the event payload
	// We need the tenant ID as client_reference_id
	subID := "sub_test_" + uuid.New().String()[:8]
	custID := "cus_test_" + uuid.New().String()[:8]

	webhookPayload := map[string]interface{}{
		"id":          "evt_test_" + uuid.New().String(),
		"object":      "event",
		"api_version": "2023-10-16", // Required by stripe-go library
		"type":        "checkout.session.completed",
		"data": map[string]interface{}{
			"object": map[string]interface{}{
				"id":                  "cs_test_" + uuid.New().String(),
				"object":              "checkout.session",
				"client_reference_id": tenantID.String(),
				"customer": map[string]interface{}{
					"id": custID,
				},
				"subscription": map[string]interface{}{
					"id":        subID,
					"trial_end": nil, // Immediate active
				},
				"metadata": map[string]interface{}{
					"tenant_id": tenantID.String(),
					"plan":      "tensor",
				},
			},
		},
	}

	payloadBytes, _ := json.Marshal(webhookPayload)

	// Generate Signature
	timestamp := time.Now().Unix()
	signedPayload := fmt.Sprintf("%d.%s", timestamp, string(payloadBytes))
	secret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedPayload))
	signature := hex.EncodeToString(mac.Sum(nil))

	stripeHeader := fmt.Sprintf("t=%d,v1=%s", timestamp, signature)

	// Send Webhook
	resp, _ = te.post("/v1/billing/webhook", webhookPayload, map[string]string{
		"Stripe-Signature": stripeHeader,
	})

	// Should be accepted (200 OK)
	te.assertStatus(resp, http.StatusOK)

	// 6. Verify Async Processing
	// The webhook handler writes to DB, then a worker processes it.
	// Wait for processing to complete.
	t.Log("Step 6: Waiting for webhook processing...")

	processed := false
	for i := 0; i < 20; i++ { // Wait up to 2 seconds
		var plan string
		var stripeStatus *string

		// Check tenant record directly to see if it was updated
		err := te.store.pool.QueryRow(context.Background(),
			"SELECT plan, stripe_status FROM tenants WHERE id = $1",
			tenantID).Scan(&plan, &stripeStatus)

		if err == nil {
			statusStr := ""
			if stripeStatus != nil {
				statusStr = *stripeStatus
			}
			if plan == "tensor" {
				t.Logf("Tenant upgraded to %s (Status: %s)", plan, statusStr)
				processed = true
				break
			}
		}
		time.Sleep(100 * time.Millisecond)
	}

	if !processed {
		// Check if the event failed in the webhook table
		var status string
		var lastError *string
		_ = te.store.pool.QueryRow(context.Background(),
			"SELECT status, last_error FROM stripe_webhook_events ORDER BY created_at DESC LIMIT 1").
			Scan(&status, &lastError)
		t.Fatalf("Webhook processing failed or timed out. Event status: %s, Error: %v", status, lastError)
	}

	// 7. Verify Entitlements
	// (Simulate checking if limits are increased, though limits are hardcoded by plan in code for now)
	t.Log("Step 7: Verifying successful upgrade...")

	// Verify database record has correct stripe IDs
	var dbCustID, dbSubID string
	err = te.store.pool.QueryRow(context.Background(),
		"SELECT stripe_customer_id, stripe_subscription_id FROM tenants WHERE id = $1",
		tenantID).Scan(&dbCustID, &dbSubID)

	if err != nil {
		t.Fatalf("Failed to query tenant: %v", err)
	}

	if dbCustID != custID {
		t.Errorf("Expected customer ID %s, got %s", custID, dbCustID)
	}
	if dbSubID != subID {
		t.Errorf("Expected subscription ID %s, got %s", subID, dbSubID)
	}

	t.Log("Test E2E_FullBillingLifecycle Passed!")
}

func getTenantInfo(t *testing.T, te *testEnv, email string) (string, bool, uuid.UUID) {
	var status *string
	var verifiedAt *time.Time
	var id uuid.UUID
	err := te.store.pool.QueryRow(context.Background(),
		"SELECT stripe_status, verified_at, id FROM tenants WHERE email = $1", email).
		Scan(&status, &verifiedAt, &id)
	if err != nil {
		t.Fatalf("Failed to get tenant info: %v", err)
	}
	s := ""
	if status != nil {
		s = *status
	}
	return s, verifiedAt != nil, id
}
