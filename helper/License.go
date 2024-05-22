package helper

import (
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"strings"
)

func GenerateLicense(licensesName, assigneeName, expiryDate string, productCodeSet []string) (string, error) {
	licenseID := uuid.New().String()
	licenseID = strings.Replace(licenseID, "-", "", -1)
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
	privateKey := readRSAPrivateKey(privateKeyFileName)
	// publicKey := readRSAPublicKey(publicKeyFileName)
	signatureBase64 := signWithRSA(privateKey, licensePartJSON)
	cert := readX509Certificate(certFileName)
	certBase64 := base64.StdEncoding.EncodeToString(cert.Raw)
	println(signatureBase64)
	return licenseID + "-" + licensePartBase64 + "-" + signatureBase64 + "-" + certBase64, nil
}
