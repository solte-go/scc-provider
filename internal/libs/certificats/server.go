package certificats

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

func Server(dest, ca, nameCert string) error {
	caCertPEM, err := os.ReadFile(dest + ca + ".crt")
	if err != nil {
		log.Printf("Unable to read CA cert: %v", err)
		return err
	}

	block, _ := pem.Decode(caCertPEM)
	if block == nil {
		log.Printf("Unable to decode CA certificate PEM")
		return err
	}

	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Printf("Unable to parse CA certificate: %v", err)
		return err
	}

	// read and parse the CA's private key
	caKeyPEM, err := os.ReadFile(dest + ca + ".key")
	if err != nil {
		log.Printf("Unable to read CA key: %v", err)
		return err
	}

	block, _ = pem.Decode(caKeyPEM)
	if block == nil {
		log.Printf("Unable to decode CA key PEM")
		return err
	}

	caKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Printf("Unable to parse CA private key: %v", err)
		return err
	}

	server := &x509.Certificate{
		SerialNumber: big.NewInt(2022),
		Subject: pkix.Name{
			CommonName:   "localhost",
			Organization: []string{"Solte Dev."},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		DNSNames: []string{
			"server.localhost",
			"localhost",
		},
		IPAddresses: []net.IP{
			net.ParseIP("127.0.0.1"),
		},
	}

	serverPrivKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	pubKey := &serverPrivKey.PublicKey
	serverBytes, _ := x509.CreateCertificate(rand.Reader, server, caCert, pubKey, caKey)

	// Public key
	certOut, _ := os.Create(dest + nameCert + ".crt")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: serverBytes})
	certOut.Close()

	// Private key
	keyOut, _ := os.Create(dest + nameCert + ".key")
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(serverPrivKey)})
	keyOut.Close()

	return nil
}
