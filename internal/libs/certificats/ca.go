package certificats

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"time"
)

func RootCA(dest, caName string) error {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2021),
		Subject: pkix.Name{
			Organization:  []string{"Solte Dev."},
			Country:       []string{"FI"},
			Province:      []string{""},
			Locality:      []string{"Helsinki"},
			StreetAddress: []string{"Alesy Katty"},
			PostalCode:    []string{"13220"},
			CommonName:    "soltedev.pro",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 6, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		Version:               1,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatalf("Unable to generate CA private key: %v", err)
		return err
	}

	pubKey := &caPrivKey.PublicKey
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, pubKey, caPrivKey)
	if err != nil {
		log.Fatalf("Unable to create CA certificate: %v", err)
		return err
	}

	// Public key
	certOut, err := os.Create(dest + caName + ".crt")
	if err != nil {
		log.Fatalf("Unable to create CA crt: %v", err)
		return err
	}

	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: caBytes})
	certOut.Close()

	// Private key
	keyOut, err := os.Create(dest + caName + ".key")
	if err != nil {
		log.Fatalf("Unable to create CA key: %v", err)
		return err
	}

	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey)})
	keyOut.Close()

	return nil
}
