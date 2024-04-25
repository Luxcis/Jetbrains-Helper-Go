package helper

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/mholt/archiver/v3"
	"io"
	"log"
	"math/big"
	"os"
	"strings"
)

const (
	JaNetfilterFilePath = "external/agent/ja-netfilter"
	powerConfFileName   = JaNetfilterFilePath + "/config/power.conf"
)

var (
	jaNetfilterZipFile *os.File
)

func InitAgent() {
	log.Println("Agent context init loading...")
	jaNetfilterZipFile, _ = getFileOrCreate(JaNetfilterFilePath + ".zip")

	if _, err := os.Stat(JaNetfilterFilePath); os.IsNotExist(err) {
		unzipJaNetfilter()
		if !powerConfHasInit() {
			log.Println("Agent config init loading...")
			loadPowerConf()
			zipJaNetfilter()
			log.Println("Agent config init success!")
		}
	}
	log.Println("Agent context init success!")
}

func getFileOrCreate(fileName string) (*os.File, error) {
	return os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0600)
}

func powerConfHasInit() bool {
	data, err := os.ReadFile(powerConfFileName)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	powerConfStr := string(data)
	return strings.Contains(powerConfStr, "[Result]") && strings.Contains(powerConfStr, "EQUAL,")
}

func loadPowerConf() {
	ch := make(chan string)
	go func() {
		ruleValue := generatePowerConfigRule()
		configStr := generatePowerConfigStr(ruleValue)
		overridePowerConfFileContent(configStr)
		close(ch)
	}()
	<-ch // wait for the goroutine to finish
}

func generatePowerConfigRule() string {
	crt := readX509Certificate(certFile) // Assuming certificate package provides these
	publicKey := readRSAPublicKey(publicKeyFile)
	rootPublicKey := readRSAPublicKey(rootKeyFile)

	x := new(big.Int).SetBytes(crt.Signature)
	y := big.NewInt(int64(publicKey.E))      // Convert the public exponent to *big.Int
	z := rootPublicKey.N                     // Modulus of the root public key
	r := new(big.Int).Exp(x, y, publicKey.N) // Use y which is *big.Int now
	return strings.Join([]string{"EQUAL", x.String(), y.String(), z.String(), "->", r.String()}, ",")
}

func generatePowerConfigStr(ruleValue string) string {
	return "[Result]\n" + ruleValue
}

func overridePowerConfFileContent(configStr string) {
	if err := os.WriteFile(powerConfFileName, []byte(configStr), 0644); err != nil {
		log.Fatalf("Error writing file: %v", err)
	}
}

func unzipJaNetfilter() {
	if err := archiver.Unarchive(jaNetfilterZipFile.Name(), JaNetfilterFilePath); err != nil {
		log.Fatalf("Failed to unzip file: %v", err)
	}
}

func zipJaNetfilter() {
	if err := archiver.Archive([]string{JaNetfilterFilePath}, jaNetfilterZipFile.Name()); err != nil {
		log.Fatalf("Failed to zip folder: %v", err)
	}
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

func readRSAPublicKey(file *os.File) *rsa.PublicKey {
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading public key file: %v", err)
	}
	block, _ := pem.Decode(data)
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse public key: %v", err)
	}
	return pubKey.(*rsa.PublicKey)
}
