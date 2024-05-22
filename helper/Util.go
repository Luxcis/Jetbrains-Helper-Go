package helper

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
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

func readRSAPublicKey(file string) *rsa.PublicKey {
	key, err := readPemKey(file)
	if err != nil {
		log.Fatalf("无法解析公钥: %v", err)
	}
	pk := key.(*rsa.PublicKey)
	return pk
}

func readRSAPrivateKey(file string) *rsa.PrivateKey {
	key, err := readPemKey(file)
	if err != nil {
		log.Fatalf("无法解析私钥: %v", err)
	}
	pk := key.(*rsa.PrivateKey)
	return pk
}

func readX509Certificate(file string) *x509.Certificate {
	key, err := readPemKey(file)
	if err != nil {
		log.Fatalf("无法解析私钥: %v", err)
	}
	ck := key.(*x509.Certificate)
	return ck
}

func readPemKey(file string) (interface{}, error) {
	keyData, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing key")
	}

	switch block.Type {
	case "EC PRIVATE KEY":
		key, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			switch key := key.(type) {
			case *ecdsa.PrivateKey:
				return key, nil
			default:
				return nil, errors.New("unsupported private key type")
			}
		}
		return key, nil
	case "RSA PRIVATE KEY":
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return key, nil
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		switch key := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, errors.New("unknown private key type")
		}
	case "EC PUBLIC KEY", "PUBLIC KEY":
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		switch key := key.(type) {
		case *ecdsa.PublicKey, *rsa.PublicKey:
			return key, nil
		default:
			return nil, errors.New("unknown public key type")
		}
	case "CERTIFICATE":
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		return cert, nil
	default:
		return nil, errors.New("unrecognized key type: " + block.Type)
	}
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
	// 对数据进行SHA1哈希
	hash := sha1.New()
	_, err := hash.Write(data)
	if err != nil {
		log.Fatalf("哈希计算失败: %v", err)
	}
	hashed := hash.Sum(nil)
	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hashed)
	if err != nil {
		log.Fatalf("Failed to sign: %v", err)
	}
	signature := base64.StdEncoding.EncodeToString(sign)
	return signature
}
