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

func Client(dest, ca, nameCert string) error {
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

	client := &x509.Certificate{
		SerialNumber: big.NewInt(2023),
		Subject: pkix.Name{
			CommonName:   "localhost",
			Organization: []string{"SolteDev."},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 6, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		DNSNames: []string{
			"client.localhost",
			"localhost",
		},
		IPAddresses: []net.IP{
			net.ParseIP("127.0.0.1"),
		},
	}

	clientPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Printf("Unable to generate client private key: %v", err)
		return err
	}

	pubKey := &clientPrivKey.PublicKey
	// note that we're using the CA's private key and certificate here to sign the client certificate
	clientBytes, err := x509.CreateCertificate(rand.Reader, client, caCert, pubKey, caKey)
	if err != nil {
		log.Printf("Unable to create client certificate: %v", err)
		return err
	}

	// Public key
	certOut, err := os.Create(dest + nameCert + ".crt")
	if err != nil {
		log.Printf("Unable to create client.crt: %v", err)
		return err
	}

	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: clientBytes})
	if err != nil {
		log.Printf("Unable to encode client certificate to PEM: %v", err)
		return err
	}

	err = certOut.Close()
	if err != nil {
		log.Printf("Unable to close client.crt: %v", err)
	}

	// Private key
	keyOut, err := os.Create(dest + nameCert + ".key")
	if err != nil {
		log.Fatalf("Unable to create client.key: %v", err)
		return err
	}

	err = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(clientPrivKey)})
	if err != nil {
		log.Printf("Unable to encode client key to PEM: %v", err)
		return err
	}

	err = keyOut.Close()
	if err != nil {
		log.Printf("Unable to close client.key: %v", err)
		return err
	}

	return nil
}
