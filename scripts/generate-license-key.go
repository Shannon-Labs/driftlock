package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"os"
	"time"
)

// This is the public key from license.go - you need the matching private key
const licensePublicKeyHex = "046b17d1f2e12c4247f8bce6e563a440f277037d812deb33a0f4a13945d898c2964fe342e2fe1a7f9b8ee7eb4a7c0f9e162bce33576b315ececbb6406837bf51f5"

func main() {
	tier := flag.String("tier", "EVAL", "License tier (e.g., EVAL, PRO, ENTERPRISE)")
	days := flag.Int("days", 365, "Number of days until expiry")
	privateKeyHex := flag.String("private-key", "", "ECDSA P-256 private key in hex format (D value, 64 hex chars)")
	flag.Parse()

	if *privateKeyHex == "" {
		fmt.Fprintf(os.Stderr, "Error: --private-key is required\n")
		fmt.Fprintf(os.Stderr, "\nUsage: go run generate-license-key.go --private-key <hex> [--tier EVAL] [--days 365]\n")
		fmt.Fprintf(os.Stderr, "\nTo generate a private key that matches the public key:\n")
		fmt.Fprintf(os.Stderr, "  You need the private key (D value) that corresponds to:\n")
		fmt.Fprintf(os.Stderr, "  Public key: %s\n", licensePublicKeyHex)
		fmt.Fprintf(os.Stderr, "\nIf you don't have the private key, you'll need to:\n")
		fmt.Fprintf(os.Stderr, "  1. Find where Shannon Labs stores the signing key\n")
		fmt.Fprintf(os.Stderr, "  2. Or generate a new keypair and update license.go with the new public key\n")
		os.Exit(1)
	}

	// Parse private key
	d, ok := new(big.Int).SetString(*privateKeyHex, 16)
	if !ok || d.Sign() <= 0 {
		fmt.Fprintf(os.Stderr, "Error: invalid private key hex\n")
		os.Exit(1)
	}

	curve := elliptic.P256()
	priv := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{Curve: curve},
		D:         d,
	}
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(d.Bytes())

	// Verify it matches the public key
	pubKeyBytes, _ := hex.DecodeString(licensePublicKeyHex)
	if len(pubKeyBytes) != 65 || pubKeyBytes[0] != 4 {
		fmt.Fprintf(os.Stderr, "Error: invalid public key format\n")
		os.Exit(1)
	}
	expectedX := new(big.Int).SetBytes(pubKeyBytes[1:33])
	expectedY := new(big.Int).SetBytes(pubKeyBytes[33:])
	if priv.PublicKey.X.Cmp(expectedX) != 0 || priv.PublicKey.Y.Cmp(expectedY) != 0 {
		fmt.Fprintf(os.Stderr, "Error: private key does not match the public key in license.go\n")
		fmt.Fprintf(os.Stderr, "  Expected public key: %s\n", licensePublicKeyHex)
		fmt.Fprintf(os.Stderr, "  Generated public key: %s\n", publicKeyToHex(&priv.PublicKey))
		os.Exit(1)
	}

	// Generate license key
	expiry := time.Now().Add(time.Duration(*days) * 24 * time.Hour).Unix()
	key := makeLicenseKey(priv, *tier, expiry)

	fmt.Printf("DRIFTLOCK_LICENSE_KEY=%s\n", key)
	fmt.Fprintf(os.Stderr, "\nLicense details:\n")
	fmt.Fprintf(os.Stderr, "  Tier: %s\n", *tier)
	fmt.Fprintf(os.Stderr, "  Expires: %s\n", time.Unix(expiry, 0).UTC().Format(time.RFC3339))
	fmt.Fprintf(os.Stderr, "\nExport it:\n")
	fmt.Fprintf(os.Stderr, "  export DRIFTLOCK_LICENSE_KEY=%s\n", key)
}

func makeLicenseKey(priv *ecdsa.PrivateKey, tier string, expiry int64) string {
	message := []byte(fmt.Sprintf("%s.%d", tier, expiry))
	digest := sha256.Sum256(message)
	r, s, err := ecdsa.Sign(rand.Reader, priv, digest[:])
	if err != nil {
		panic(fmt.Sprintf("signing failed: %v", err))
	}
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

func publicKeyToHex(pub *ecdsa.PublicKey) string {
	xBytes := pub.X.Bytes()
	yBytes := pub.Y.Bytes()
	// Pad to 32 bytes each
	xPadded := make([]byte, 32)
	yPadded := make([]byte, 32)
	copy(xPadded[32-len(xBytes):], xBytes)
	copy(yPadded[32-len(yBytes):], yBytes)
	return fmt.Sprintf("04%x%x", xPadded, yPadded)
}
