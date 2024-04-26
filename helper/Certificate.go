package helper

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
	"math/big"
	"os"
	"time"
)

const (
	rootKeyFileName    = "external/certificate/root.key"
	privateKeyFileName = "external/certificate/private.key"
	publicKeyFileName  = "external/certificate/public.key"
	certFileName       = "external/certificate/ca.crt"
)

func rootKeyFile() *os.File {
	return OpenFile(rootKeyFileName)
}
func privateKeyFile() *os.File {
	return OpenFile(privateKeyFileName)
}
func publicKeyFile() *os.File {
	return OpenFile(publicKeyFileName)
}
func certFile() *os.File {
	return OpenFile(certFileName)
}

func InitCertificate() {
	log.Println("Certificate context init loading...")
	if !FileExists(privateKeyFileName) || !FileExists(publicKeyFileName) || !FileExists(certFileName) {
		log.Println("Certificate context generate loading...")
		generateCertificate()
		log.Println("Certificate context generate success!")
	}
	log.Println("Certificate context init success!")
}

func generateCertificate() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatalf("Failed to generate RSA key pair: %v", err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatalf("Failed to marshal public key: %v", err)
	}
	writePemFile(publicKeyFileName, "PUBLIC KEY", publicKeyBytes)

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	writePemFile(privateKeyFileName, "PRIVATE KEY", privateKeyBytes)

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
	}

	certTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: "JetProfile CA",
		},
		NotBefore: time.Now().Add(-24 * time.Hour),
		NotAfter:  time.Now().AddDate(100, 0, 0),
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}
	writePemFile(certFileName, "CERTIFICATE", certBytes)
}
