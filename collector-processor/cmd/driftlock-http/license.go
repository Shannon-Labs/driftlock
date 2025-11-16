package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
)

const licensePublicKeyHex = "046b17d1f2e12c4247f8bce6e563a440f277037d812deb33a0f4a13945d898c2964fe342e2fe1a7f9b8ee7eb4a7c0f9e162bce33576b315ececbb6406837bf51f5"

var licensePublicKey *ecdsa.PublicKey

func init() {
	pk, err := parseECDSAPublicKey(licensePublicKeyHex)
	if err != nil {
		panic(fmt.Sprintf("invalid license public key: %v", err))
	}
	licensePublicKey = pk
}

type licenseMetadata struct {
	Tier        string
	ExpiresAt   time.Time
	Fingerprint string
}

type licenseStatus struct {
	Status      string    `json:"status"`
	Tier        string    `json:"tier"`
	ExpiresAt   time.Time `json:"expires_at"`
	Fingerprint string    `json:"fingerprint"`
	Message     string    `json:"message,omitempty"`
}

var licenseInfo licenseMetadata

func loadLicense(now time.Time) (licenseMetadata, error) {
	// Development mode bypass - for local testing only
	if os.Getenv("DRIFTLOCK_DEV_MODE") == "true" {
		return licenseMetadata{
			Tier:        "DEV",
			ExpiresAt:   now.Add(365 * 24 * time.Hour), // 1 year from now
			Fingerprint: "dev-mode",
		}, nil
	}

	raw := strings.TrimSpace(os.Getenv("DRIFTLOCK_LICENSE_KEY"))
	if raw == "" {
		return licenseMetadata{}, fmt.Errorf("DRIFTLOCK_LICENSE_KEY not set (set DRIFTLOCK_DEV_MODE=true for development)")
	}
	parts := strings.Split(raw, ".")
	if len(parts) != 3 {
		return licenseMetadata{}, fmt.Errorf("license key malformed")
	}
	tier := parts[0]
	expiryUnix, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return licenseMetadata{}, fmt.Errorf("invalid expiry: %w", err)
	}
	expiry := time.Unix(expiryUnix, 0).UTC()
	message := strings.Join(parts[:2], ".")
	sig, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return licenseMetadata{}, fmt.Errorf("invalid signature encoding")
	}
	if !verifySignature(message, sig) {
		return licenseMetadata{}, fmt.Errorf("license signature invalid")
	}
	if now.After(expiry) {
		return licenseMetadata{}, fmt.Errorf("license expired at %s", expiry.UTC().Format(time.RFC3339))
	}
	return licenseMetadata{
		Tier:        tier,
		ExpiresAt:   expiry,
		Fingerprint: fingerprint(raw),
	}, nil
}

func verifySignature(message string, sig []byte) bool {
	if len(sig) != 64 {
		return false
	}
	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:])
	digest := sha256.Sum256([]byte(message))
	return ecdsa.Verify(licensePublicKey, digest[:], r, s)
}

func parseECDSAPublicKey(hexKey string) (*ecdsa.PublicKey, error) {
	raw, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, err
	}
	if len(raw) != 65 || raw[0] != 4 {
		return nil, fmt.Errorf("public key must be uncompressed P-256")
	}
	x := new(big.Int).SetBytes(raw[1:33])
	y := new(big.Int).SetBytes(raw[33:])
	pk := &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}
	if !pk.Curve.IsOnCurve(pk.X, pk.Y) {
		return nil, fmt.Errorf("point not on curve")
	}
	return pk, nil
}

func fingerprint(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])[:12]
}

func currentLicenseStatus(now time.Time) licenseStatus {
	status := licenseStatus{
		Tier:        licenseInfo.Tier,
		ExpiresAt:   licenseInfo.ExpiresAt,
		Fingerprint: licenseInfo.Fingerprint,
	}
	if licenseInfo.Fingerprint == "" {
		status.Status = "invalid"
		status.Message = "license missing"
		return status
	}
	remaining := licenseInfo.ExpiresAt.Sub(now)
	if remaining <= 0 {
		status.Status = "expired"
		status.Message = "license expired"
		return status
	}
	status.Status = "valid"
	if remaining < 14*24*time.Hour {
		status.Message = "license expiring soon"
	}
	return status
}
