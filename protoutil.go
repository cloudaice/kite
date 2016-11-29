package kite

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"net"
	"strings"
)

const (
	UniqSize = 6
)

type Uniq [UniqSize]byte

// HardwareAddr return [UniqSize]byte
func HardwareAddr() Uniq {
	def := Uniq{}
	inters, err := net.Interfaces()
	if err != nil {
		return def
	}

	for _, inter := range inters {
		if strings.Index(inter.Name, "eth") == 0 && len(inter.HardwareAddr) >= UniqSize {
			copy(def[:], inter.HardwareAddr)
			return def
		}
	}

	return def
}

// RandBytes generate [lenn]byte
func RandBytes(lenn int) []byte {
	const alpha = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, lenn)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alpha[b%byte(len(alpha))]
	}
	return bytes
}

// Gzip
func Gzip(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	if err != nil {
		return nil, err
	}

	_, err = w.Write(data)
	if err != nil {
		return nil, err
	}
	w.Close()
	return b.Bytes(), nil
}

// UnGzip
func UnGzip(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return ioutil.ReadAll(r)
}

// Md5
func Md5(data []byte) []byte {
	sum := md5.Sum(data)
	return sum[:]
}

// AesEncrypt data block is 128 bit and key length is 128 bit too.
func AesEncrypt(originData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	originData = PKCSSPadding(originData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(originData))
	blockMode.CryptBlocks(crypted, originData)
	return crypted, nil
}

// AesEncrypt data block is 128 bit and key length is 128 bit too.
func AesDecrypt(crypted, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	originData := make([]byte, len(crypted))
	blockMode.CryptBlocks(originData, crypted)
	originData = PKCSSUnPadding(originData)
	return originData, nil
}

// PKCSSPadding adapt to aes encryption, doing like follow:
// block size is 4, text: "abcedd"
// len("abcedd") is 6, 6 % 4 = 2, let # = byte(2),
// the text after padding is "abcedd##"
// the most value can be show is 255, and aes largest block length is 256
func PKCSSPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCSSUnPadding used to depadding
func PKCSSUnPadding(originData []byte) []byte {
	length := len(originData)
	unpadding := int(originData[length-1])
	return originData[:(length - unpadding)]
}

// EqualSlice judge every byte in dst and src
func EqualSlice(dst, src []byte) bool {
	if len(dst) != len(src) {
		return false
	}
	for idx, b := range dst {
		if b != src[idx] {
			return false
		}
	}
	return true
}

// Generate Rsa Public Key and Private Key
func GenRsaKey(bits int) (pubKeyByte, privKeyByte []byte, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	privateKeyDer := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDer,
	}
	privKeyByte = pem.EncodeToMemory(privateKeyBlock)

	publicKey := &privateKey.PublicKey
	publicKeyDer, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, nil, err
	}
	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyDer,
	}
	pubKeyByte = pem.EncodeToMemory(publicKeyBlock)
	return
}

func RsaEncrypt(data, pubKey []byte) ([]byte, error) {
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return nil, errors.New("Public Key Error")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

func RsaDecrypt(data, privKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privKey)
	if block == nil {
		return nil, errors.New("Private Key Error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, data)
}
