package export

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shannon-labs/driftlock/api-server/internal/models"
	"github.com/shannon-labs/driftlock/pkg/version"
)

// Exporter handles evidence bundle generation
type Exporter struct {
	signBundles bool
}

// NewExporter creates a new evidence exporter
func NewExporter(signBundles bool) *Exporter {
	return &Exporter{signBundles: signBundles}
}

// ExportAnomaly generates an evidence bundle for an anomaly
func (e *Exporter) ExportAnomaly(anomaly *models.Anomaly, exportedBy string) (*models.EvidenceBundle, error) {
	bundle := &models.EvidenceBundle{
		Anomaly:    *anomaly,
		ExportedAt: time.Now(),
		ExportedBy: exportedBy,
		Version:    version.Version(),
		AdditionalMetadata: map[string]interface{}{
			"export_format": "driftlock-evidence-v1",
			"compliance_frameworks": []string{"DORA", "NIS2", "AI Act"},
		},
	}

	// Generate cryptographic signature if enabled
	if e.signBundles {
		sig, err := e.signBundle(bundle)
		if err != nil {
			return nil, fmt.Errorf("failed to sign bundle: %w", err)
		}
		bundle.Signature = &sig
	}

	return bundle, nil
}

// ExportJSON exports the evidence bundle as JSON
func (e *Exporter) ExportJSON(bundle *models.EvidenceBundle) ([]byte, error) {
	data, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal bundle: %w", err)
	}
	return data, nil
}

// signBundle generates a SHA-256 signature for tamper-evidence
func (e *Exporter) signBundle(bundle *models.EvidenceBundle) (string, error) {
	// Create canonical representation for signing
	signData := struct {
		AnomalyID  string    `json:"anomaly_id"`
		Timestamp  time.Time `json:"timestamp"`
		NCDScore   float64   `json:"ncd_score"`
		PValue     float64   `json:"p_value"`
		ExportedAt time.Time `json:"exported_at"`
		ExportedBy string    `json:"exported_by"`
	}{
		AnomalyID:  bundle.Anomaly.ID.String(),
		Timestamp:  bundle.Anomaly.Timestamp,
		NCDScore:   bundle.Anomaly.NCDScore,
		PValue:     bundle.Anomaly.PValue,
		ExportedAt: bundle.ExportedAt,
		ExportedBy: bundle.ExportedBy,
	}

	canonical, err := json.Marshal(signData)
	if err != nil {
		return "", err
	}

	// Compute SHA-256 hash
	hash := sha256.Sum256(canonical)
	signature := hex.EncodeToString(hash[:])

	return signature, nil
}

// VerifySignature verifies the bundle signature
func (e *Exporter) VerifySignature(bundle *models.EvidenceBundle) (bool, error) {
	if bundle.Signature == nil {
		return false, fmt.Errorf("bundle has no signature")
	}

	// Regenerate signature
	sig, err := e.signBundle(bundle)
	if err != nil {
		return false, err
	}

	// Compare signatures
	return sig == *bundle.Signature, nil
}
