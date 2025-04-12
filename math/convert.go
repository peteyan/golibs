package math

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/emmansun/gmsm/sm3"
	"github.com/emmansun/gmsm/sm4"
	"github.com/peteyan/golibs/strings"
	"github.com/zeromicro/go-zero/core/logx"
	"regexp"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type AlgorithmType int

const (
	SM4_ECB_PKCS5Padding AlgorithmType = iota
	SM4_ECB_NoPadding
)

// GenerateRandomKey 生成指定长度的随机密钥
// 推荐的密钥长度应至少为16字节（128位），更长的密钥提供更高的安全性。推荐使用32字节（256位）的密钥。
func GenerateRandomKey(length int) ([]byte, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// GenerateRandomKeyStr 生成指定长度的随机密钥可读字符串
func GenerateRandomKeyStr(length int) (string, error) {
	byteArray := make([]byte, length)
	_, err := rand.Read(byteArray)
	if err != nil {
		return "", err
	}
	for i, b := range byteArray {
		byteArray[i] = charset[b%byte(len(charset))]
	}
	return string(byteArray), nil
}

// HashSM3HMAC 计算HMAC值，使用SM3作为哈希函数
func HashSM3HMAC(key, plaintext []byte) []byte {
	h := hmac.New(sm3.New, key)
	h.Write(plaintext)
	return h.Sum(nil)
}

// pkcs5Padding 实现 PKCS5Padding (与 PKCS7Padding 相同)
func pkcs5Padding(src []byte) []byte {
	padding := sm4.BlockSize - len(src)%sm4.BlockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

// pkcs5UnPadding 去除 PKCS5Padding
func pkcs5UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, fmt.Errorf("src length is zero")
	}
	unPadding := int(src[length-1])
	if unPadding > sm4.BlockSize || unPadding == 0 {
		return nil, fmt.Errorf("invalid padding")
	}
	return src[:(length - unPadding)], nil
}

// EncryptSM4 使用SM4进行加密
func EncryptSM4(algorithm AlgorithmType, key, plaintext []byte) ([]byte, error) {
	cipher, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	switch algorithm {
	case SM4_ECB_PKCS5Padding:
		plaintext = pkcs5Padding(plaintext)
	case SM4_ECB_NoPadding:
		// 确保输入的明文长度是16字节的倍数
		if len(plaintext)%sm4.BlockSize != 0 {
			return nil, fmt.Errorf("plaintext length must be a multiple of %d bytes", sm4.BlockSize)
		}
	}
	ciphertext := make([]byte, len(plaintext))
	for i := 0; i < len(plaintext); i += sm4.BlockSize {
		cipher.Encrypt(ciphertext[i:i+sm4.BlockSize], plaintext[i:i+sm4.BlockSize])
	}
	return ciphertext, nil
}

// DecryptSM4 使用SM4进行解密
func DecryptSM4(algorithm AlgorithmType, key, ciphertext []byte) ([]byte, error) {
	cipher, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 确保输入的密文长度是16字节的倍数
	if len(ciphertext)%sm4.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext length must be a multiple of %d bytes", sm4.BlockSize)
	}
	plaintext := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += sm4.BlockSize {
		cipher.Decrypt(plaintext[i:i+sm4.BlockSize], ciphertext[i:i+sm4.BlockSize])
	}
	switch algorithm {
	case SM4_ECB_PKCS5Padding:
		plaintext, err = pkcs5UnPadding(plaintext)
		if err != nil {
			return nil, err
		}
	case SM4_ECB_NoPadding:
	}
	return plaintext, nil
}

// SignRSA 私钥签名
func SignRSA(signContent string, privateKey string, hash crypto.Hash) string {
	privateKey, _ = cleanKeys(privateKey)
	h2 := sha256.New()
	h2.Write([]byte(signContent))
	hashed := h2.Sum(nil)
	priKey, err := parsePrivateKey(privateKey)
	if err != nil {
		logx.Errorf("parse private key error: %v", err)
		return ""
	}
	signature, err := rsa.SignPKCS1v15(rand.Reader, priKey, hash, hashed)
	if err != nil {
		logx.Errorf("sign error: %v", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(signature)
}

// parsePrivateKey 解析私钥
func parsePrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	block, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}
	priKey, err := x509.ParsePKCS8PrivateKey(block)
	if err != nil {
		return nil, err
	}
	rsaPriKey, ok := priKey.(*rsa.PrivateKey)
	if !ok {
		return nil, err
	}
	return rsaPriKey, nil
}

// VerifySignRSA 公钥验签
func VerifySignRSA(signContent, sign, publicKey string, hash crypto.Hash) bool {
	h2 := sha256.New()
	h2.Write([]byte(signContent))
	hashed := h2.Sum(nil)
	pubKey, err := parsePublicKey(publicKey)
	if err != nil {
		logx.Errorf("parse public key error: %v, key is: %s", err, strings.DesensitizeCommon(publicKey))
		return false
	}
	sig, _ := base64.StdEncoding.DecodeString(sign)
	err = rsa.VerifyPKCS1v15(pubKey, hash, hashed, sig)
	if err != nil {
		logx.Errorf("verify sign error: %v, signContent is: %s, pubKey is: %s",
			err, signContent, strings.DesensitizeCommon(publicKey))
		return false
	}
	return true
}

// parsePublicKey 解析公钥
func parsePublicKey(publicKey string) (*rsa.PublicKey, error) {
	/*block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("公钥信息错误！")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)*/
	//key, _ := hex.DecodeString(publicKey)
	decodeKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, err
	}
	pubKey, err := x509.ParsePKIXPublicKey(decodeKey)
	if err != nil {
		return nil, err
	}
	return pubKey.(*rsa.PublicKey), nil
}

func cleanKeys(key string) (string, error) {
	//key = strings.ReplaceAll(key, "-----BEGIN PRIVATE KEY-----", "")
	//key = strings.ReplaceAll(key, "-----END PRIVATE KEY-----", "")
	re, err := regexp.Compile("\\s*|\t|\r|\n|")
	if err != nil {
		return "", err
	}
	key = re.ReplaceAllString(key, "")
	return key, nil
}

// Encrypt 公钥加密
//func Encrypt(content, publicKey string) string {
//	rsa.EncryptPKCS1v15(rand io.Reader, pub *PublicKey, []byte(content)) ([]byte, error)
//
//}

// DecryptRSA 私钥解密
func DecryptRSA(cipherText, privateKey string) string {
	privateKey, _ = cleanKeys(privateKey)
	priKey, err := parsePrivateKey(privateKey)
	if err != nil {
		return ""
	}
	t, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return ""
	}
	b, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, t)
	if err != nil {
		return ""
	}
	return string(b)
}

// EncryptRSA 私钥解密
func EncryptRSA(plainText, publicKey string) string {
	publicKey, _ = cleanKeys(publicKey)
	pubKey, err := parsePublicKey(publicKey)
	if err != nil {
		return ""
	}
	b, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, []byte(plainText))
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}
