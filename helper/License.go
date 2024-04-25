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
	"errors"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
)

func GenerateLicense(licensesName, assigneeName, expiryDate string, productCodeSet []string) (string, error) {
	licenseID := uuid.NewString()
	var products []Product
	for _, code := range productCodeSet {
		products = append(products, Product{
			Code:         code,
			FallbackDate: expiryDate,
			PaidUpTo:     expiryDate,
		})
	}
	licensePart := LicensePart{
		LicenseID:    licenseID,
		LicenseeName: licensesName,
		AssigneeName: assigneeName,
		Products:     products,
		Metadata:     "0120230914PSAX000005",
	}

	licensePartJSON, err := json.Marshal(licensePart)
	if err != nil {
		return "", err
	}
	licensePartBase64 := base64.StdEncoding.EncodeToString(licensePartJSON)

	privateKey, err := readPrivateKey(privateKeyFileName)
	if err != nil {
		return "", err
	}
	_, err = readPublicKey(publicKeyFileName)
	if err != nil {
		return "", err
	}

	signature, err := signWithRSA(privateKey, licensePartJSON)
	if err != nil {
		return "", err
	}
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	cert, err := readCertificate(certFileName)
	if err != nil {
		return "", err
	}
	certEncoded, err := x509.MarshalPKIXPublicKey(cert.PublicKey)
	if err != nil {
		return "", err
	}
	certBase64 := base64.StdEncoding.EncodeToString(certEncoded)

	return licenseID + "-" + licensePartBase64 + "-" + signatureBase64 + "-" + certBase64, nil
}

func readPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func readPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}
	return rsaPub, nil
}

func signWithRSA(privateKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	hashed := sha256.Sum256(data)
	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
}

func readCertificate(path string) (*x509.Certificate, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	return x509.ParseCertificate(block.Bytes)
}
