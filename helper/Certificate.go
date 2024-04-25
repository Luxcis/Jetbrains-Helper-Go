package helper

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

const (
	rootKeyFileName    = "external/certificate/root.key"
	privateKeyFileName = "external/certificate/private.key"
	publicKeyFileName  = "external/certificate/public.key"
	certFileName       = "external/certificate/ca.crt"
)

var (
	rootKeyFile    *os.File
	privateKeyFile *os.File
	publicKeyFile  *os.File
	certFile       *os.File
)

func InitCertificate() {
	log.Println("Certificate context init loading...")
	rootKeyFile = OpenFile(rootKeyFileName)
	if !FileExists(privateKeyFileName) || !FileExists(publicKeyFileName) || !FileExists(certFileName) {
		log.Println("Certificate context generate loading...")
		generateCertificate()
		log.Println("Certificate context generate success!")
	} else {
		privateKeyFile = OpenFile(privateKeyFileName)
		publicKeyFile = OpenFile(publicKeyFileName)
		certFile = OpenFile(certFileName)
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

func writePemFile(fileName string, pemType string, bytes []byte) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Failed to close data to PEM file: %v", err)
		}
	}(file)

	pemBlock := &pem.Block{
		Type:  pemType,
		Bytes: bytes,
	}
	if err := pem.Encode(file, pemBlock); err != nil {
		log.Fatalf("Failed to write data to PEM file: %v", err)
	}
}
