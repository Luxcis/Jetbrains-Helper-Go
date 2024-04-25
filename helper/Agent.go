package helper

import (
	"github.com/mholt/archiver/v3"
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
	ruleValue := generatePowerConfigRule()
	configStr := generatePowerConfigStr(ruleValue)
	overridePowerConfFileContent(configStr)
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
	// 删除已存在的压缩文件
	if _, err := os.Stat(jaNetfilterZipFile.Name()); !os.IsNotExist(err) {
		if err := os.Remove(jaNetfilterZipFile.Name()); err != nil {
			log.Fatalf("Failed to remove existing zip file: %v", err)
		}
	}

	// 压缩文件
	if err := archiver.Archive([]string{JaNetfilterFilePath}, jaNetfilterZipFile.Name()); err != nil {
		log.Fatalf("Failed to zip folder: %v", err)
	}
}
