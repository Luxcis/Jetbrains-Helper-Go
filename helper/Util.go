package helper

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"io"
	"log"
	"os"
)

func ReadJson(path string, payload interface{}) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return err
	}
	return nil
}

func OpenFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("Failed to open or create file: %v", err)
	}
	return file
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func readRSAPublicKey(file *os.File) *rsa.PublicKey {
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading public key file: %v", err)
	}
	block, _ := pem.Decode(data)
	if block.Type == "CERTIFICATE" {
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			log.Fatalf("Failed to parse certificate public key: %v", err)
		}
		return cert.PublicKey.(*rsa.PublicKey)
	} else {
		pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			log.Fatalf("Failed to parse public key: %v", err)
		}
		return pubKey.(*rsa.PublicKey)
	}
}

func readRSAPrivateKey(file *os.File) *rsa.PrivateKey {
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading private key file: %v", err)
	}
	block, _ := pem.Decode(data)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	return privateKey
}

func readX509Certificate(file *os.File) *x509.Certificate {
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading certificate file: %v", err)
	}
	block, _ := pem.Decode(data)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse certificate: %v", err)
	}
	return cert
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

func signWithRSA(privateKey *rsa.PrivateKey, data []byte) string {
	hashed := sha256.Sum256(data)
	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		log.Fatalf("Failed to sign: %v", err)
	}
	signature := base64.StdEncoding.EncodeToString(sign)
	return signature
}
