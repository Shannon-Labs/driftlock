package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"testing"
	"time"
)

func TestLoadLicenseValid(t *testing.T) {
	priv := testPrivateKey()
	expiry := time.Now().Add(24 * time.Hour).Unix()
	key := makeLicenseKey(priv, "EVAL", expiry)
	t.Setenv("DRIFTLOCK_LICENSE_KEY", key)
	meta, err := loadLicense(time.Now())
	if err != nil {
		t.Fatalf("expected valid license, got %v", err)
	}
	if meta.Tier != "EVAL" {
		t.Fatalf("unexpected tier %s", meta.Tier)
	}
}

func TestLoadLicenseInvalidSignature(t *testing.T) {
	priv := testPrivateKey()
	expiry := time.Now().Add(24 * time.Hour).Unix()
	key := makeLicenseKey(priv, "EVAL", expiry)
	key = key + "tamper"
	t.Setenv("DRIFTLOCK_LICENSE_KEY", key)
	if _, err := loadLicense(time.Now()); err == nil {
		t.Fatalf("expected invalid signature error")
	}
}

func TestLoadLicenseExpired(t *testing.T) {
	priv := testPrivateKey()
	expiry := time.Now().Add(-1 * time.Hour).Unix()
	key := makeLicenseKey(priv, "EVAL", expiry)
	t.Setenv("DRIFTLOCK_LICENSE_KEY", key)
	if _, err := loadLicense(time.Now()); err == nil {
		t.Fatalf("expected expiry error")
	}
}

func testPrivateKey() *ecdsa.PrivateKey {
	curve := elliptic.P256()
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = curve
	priv.D = big.NewInt(1)
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(priv.D.Bytes())
	return priv
}

func makeLicenseKey(priv *ecdsa.PrivateKey, tier string, expiry int64) string {
	message := []byte(fmt.Sprintf("%s.%d", tier, expiry))
	digest := sha256.Sum256(message)
	r, s, _ := ecdsa.Sign(rand.Reader, priv, digest[:])
	sig := append(padScalar(r.Bytes()), padScalar(s.Bytes())...)
	return fmt.Sprintf("%s.%d.%s", tier, expiry, base64.RawStdEncoding.EncodeToString(sig))
}

func padScalar(b []byte) []byte {
	if len(b) == 32 {
		return b
	}
	buf := make([]byte, 32)
	copy(buf[32-len(b):], b)
	return buf
}
