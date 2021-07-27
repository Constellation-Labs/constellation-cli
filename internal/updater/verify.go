package updater

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
)

func readPublicKey(fileName string) (*rsa.PublicKey, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, info.Size())
	file.Read(buffer)
	block, _ := pem.Decode(buffer)
	publicInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return publicInterface.(*rsa.PublicKey), nil
}

func Verify(fileName string, signatureFileName string, publicKeyFileName string) error {
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	signature, err := ioutil.ReadFile(signatureFileName)
	if err != nil {
		return err
	}

	publicKey, err := readPublicKey(publicKeyFileName)
	if err != nil {
		return err
	}

	sha := sha512.New()
	sha.Write(fileContent)
	hashed := sha.Sum(nil)

	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA512, hashed, signature)
}
