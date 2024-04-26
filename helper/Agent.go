package helper

import (
	"archive/zip"
	"fmt"
	"github.com/mholt/archiver/v3"
	"io"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

const (
	JaNetfilterFilePath = "external/agent/ja-netfilter"
	JaNetfilterZipFile  = JaNetfilterFilePath + ".zip"
	powerConfFileName   = JaNetfilterFilePath + "/config/power.conf"
)

var jaNetfilterZipFile *os.File

func InitAgent() {
	log.Println("Agent context init loading...")
	jaNetfilterZipFile, _ = getFileOrCreate(JaNetfilterZipFile)

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
	crt := readX509Certificate(certFile()) // Assuming certificate package provides these
	publicKey := readRSAPublicKey(publicKeyFile())
	rootPublicKey := readRSAPublicKey(rootKeyFile())

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
	// 检查是否已存在同名的 zip 文件
	if _, err := os.Stat(JaNetfilterZipFile); err == nil {
		// 如果存在，重命名之前的文件
		newName := JaNetfilterFilePath + "_old.zip"
		err := os.Rename(JaNetfilterZipFile, newName)
		if err != nil {
			fmt.Println("Failed to rename existing zip file:", err)
			return
		}
		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				log.Fatalf("Failed to remove existing zip file: %v", err)
			}
		}(newName) // 在压缩完成后删除旧的 zip 文件
		fmt.Println("Existing zip file renamed to:", newName)
	}

	// 创建一个新的 zip 文件
	zipFile, err := os.Create(JaNetfilterZipFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(zipFile *os.File) {
		err := zipFile.Close()
		if err != nil {
			if err != nil {
				log.Fatalf("Failed to close data to zip file: %v", err)
			}
		}
	}(zipFile)

	// 创建一个 zip.Writer
	zipWriter := zip.NewWriter(zipFile)
	defer func(zipWriter *zip.Writer) {
		err := zipWriter.Close()
		if err != nil {
			log.Fatalf("Failed to close data to zip writer: %v", err)
		}
	}(zipWriter)

	// 递归地压缩目录中的文件
	err = filepath.Walk(JaNetfilterFilePath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过.DS_Store文件
		if strings.HasSuffix(info.Name(), ".DS_Store") {
			return nil
		}

		// 获取文件头信息
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// 设置文件名
		header.Name, err = filepath.Rel(JaNetfilterFilePath, filePath)
		if err != nil {
			return err
		}

		// 在 Windows 下修复 zip 文件中的路径分隔符
		header.Name = strings.ReplaceAll(header.Name, "\\", "/")

		// 如果是目录，只添加目录名到 zip 文件中
		if info.IsDir() {
			header.Name += "/"
		}

		// 创建一个 zip 文件中的新文件
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// 如果是文件，将文件内容拷贝到 zip 文件中
		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					log.Fatalf("Failed to close data to zip file: %v", err)
				}
			}(file)

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
		return
	}
}
