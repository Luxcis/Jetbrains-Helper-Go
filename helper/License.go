package helper

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
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
	privateKey := readRSAPrivateKey(privateKeyFile())
	// publicKey = readRSAPublicKey(publicKeyFile)
	signatureBase64 := signWithRSA(privateKey, licensePartJSON)
	cert := readX509Certificate(certFile())
	certEncoded, err := x509.MarshalPKIXPublicKey(cert.PublicKey)
	if err != nil {
		return "", err
	}
	certBase64 := base64.StdEncoding.EncodeToString(certEncoded)

	return licenseID + "-" + licensePartBase64 + "-" + signatureBase64 + "-" + certBase64, nil
}
