package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

//使用PKCS7进行填充，IOS也是7
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unpadding := int(origData[length-1])
	idx := length - unpadding
	if len(origData) <= idx || idx < 0 {
		return nil, errors.New(fmt.Sprintf("slice bounds out of range %v", idx))
	}
	return origData[:(idx)], nil
}

//aes加密，填充秘钥key的16位，24,32分别对应AES-128, AES-192, or AES-256.
func AesCBCEncrypt(rawData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//填充原文
	blockSize := block.BlockSize()
	rawData = PKCS7Padding(rawData, blockSize)
	//初始向量IV必须是唯一，但不需要保密
	cipherText := make([]byte, blockSize+len(rawData))
	//block大小 16
	iv := cipherText[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	//block大小和初始向量大小一定要一致
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[blockSize:], rawData)

	return cipherText, nil
}

func AesCBCDecrypt(encryptData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()

	if len(encryptData) < blockSize {
		return nil, errors.New("cipher text too short")
	}
	iv := encryptData[:blockSize]
	encryptData = encryptData[blockSize:]

	// CBC mode always works in whole blocks.
	if len(encryptData)%blockSize != 0 {
		return nil, errors.New("cipher text is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(encryptData, encryptData)
	//解填充
	encryptData, err = PKCS7UnPadding(encryptData)
	return encryptData, err
}

func EncryptByte(rawData []byte, key string) (string, error) {
	data, err := AesCBCEncrypt(rawData, []byte(key))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func DecryptByte(rawData []byte, key string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(string(rawData))
	if err != nil {
		return nil, err
	}
	dnData, err := AesCBCDecrypt(data, []byte(key))
	if err != nil {
		return nil, err
	}
	return dnData, nil
}
func Encrypt(rawData, key string) (string, error) {
	return EncryptByte([]byte(rawData), key)
}

func Decrypt(rawData string, key string) (string, error) {
	byteData, err := DecryptByte([]byte(rawData), key)
	if err != nil {
		return "", err
	}
	return string(byteData), nil
}
