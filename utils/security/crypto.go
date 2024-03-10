package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// S256Challenge creates base64 encoded sha256 challenge string derived from code.
// The padding of the result base64 string is stripped per [RFC 7636].
//
// [RFC 7636]: https://datatracker.ietf.org/doc/html/rfc7636#section-4.2
func S256Challenge(code string) string {
	h := sha256.New()
	h.Write([]byte(code))
	return strings.TrimRight(base64.URLEncoding.EncodeToString(h.Sum(nil)), "=")
}

// MD5 creates md5 hash from the provided plain text.
func MD5(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Md5ByBytes(b []byte) string {
	return fmt.Sprintf("%x", md5.Sum(b))
}

// GetMD5Encode 返回一个32位md5加密后的字符串
func GetMD5Encode(data string, salt *string) string {
	h := md5.New()
	h.Write([]byte(data))
	if salt != nil {
		h.Write([]byte(*salt))
	}
	return hex.EncodeToString(h.Sum(nil))
}

// Get16MD5Encode 返回一个16位md5加密后的字符串
func Get16MD5Encode(data string) string {
	return GetMD5Encode(data, nil)[8:24]
}

// SHA256 creates sha256 hash as defined in FIPS 180-4 from the provided text.
func SHA256(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// SHA512 creates sha512 hash as defined in FIPS 180-4 from the provided text.
func SHA512(text string) string {
	h := sha512.New()
	h.Write([]byte(text))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// HS256 creates a HMAC hash with sha256 digest algorithm.
func HS256(text string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(text))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// HS512 creates a HMAC hash with sha512 digest algorithm.
func HS512(text string, secret string) string {
	h := hmac.New(sha512.New, []byte(secret))
	h.Write([]byte(text))
	return fmt.Sprintf("%x", h.Sum(nil))
}

/**
概念
CBC（密文分组链接方式）有向量的概念, 它的实现机制使加密的各段数据之间有了联系。
加密步骤：
    首先将数据按照8个字节一组进行分组得到D1D2......Dn（若数据不是8的整数倍，用指定的PADDING数据补位）
    第一组数据D1与初始化向量I异或后的结果进行DES加密得到第一组密文C1（初始化向量I为全零）
    第二组数据D2与第一组的加密结果C1异或以后的结果进行DES加密，得到第二组密文C2
    之后的数据以此类推，得到Cn
    按顺序连为C1C2C3......Cn即为加密结果。
解密是加密的逆过程：
    首先将数据按照8个字节一组进行分组得到C1C2C3......Cn
    将第一组数据进行解密后与初始化向量I进行异或得到第一组明文D1（注意：一定是先解密再异或）
    将第二组数据C2进行解密后与第一组密文数据进行异或得到第二组数据D2
    之后依此类推，得到Dn
    按顺序连为D1D2D3......Dn即为解密结果。
特点
    不容易主动攻击,安全性好于ECB,适合传输长度长的报文,是SSL、IPSec的标准。
    每个密文块依赖于所有的信息明文消息中一个改变会影响所有密文块
    发送方和接收方都需要知道初始化向量
    加密过程是串行的，无法被并行化(在解密时，从两个邻接的密文块中即可得到一个平文块。因此，解密过程可以被并行化
*/

// AesEncrypt CBC加密 key
// iv必须是16位
func AesEncrypt(orig string, key string) string {

	origData := []byte(orig)
	k := []byte(key)

	// new Cipher
	block, _ := aes.NewCipher(k)
	// fetch block size
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)

	return base64.StdEncoding.EncodeToString(cryted)
}

// AesDecrypt CBC解密key
// iv必须是16位
func AesDecrypt(cryted string, key string) string {
	// 先解密base64
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	k := []byte(key)

	block, _ := aes.NewCipher(k)

	blockSize := block.BlockSize()

	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])

	orig := make([]byte, len(crytedByte))

	blockMode.CryptBlocks(orig, crytedByte)
	orig = PKCS7UnPadding(orig)
	return string(orig)
}

// PKCS7Padding 明文补码算法
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding 明文减码算法
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// Equal compares two hash strings for equality without leaking timing information.
func Equal(hash1 string, hash2 string) bool {
	return subtle.ConstantTimeCompare([]byte(hash1), []byte(hash2)) == 1
}

/**
ECB（电子密本方式）就是将数据按照8个字节一段进行DES加密或解密得到一段8个字节的密文或者明文，
最后一段不足8个字节，按照需求补足8个字节进行计算，之后按照顺序将计算所得的数据连在一起即可，
各段数据之间互不影响。
特点：
简单，有利于并行计算，误差不会被传送；
不能隐藏明文的模式；在密文中出现明文消息的重复
可能对明文进行主动攻击；加密消息块相互独立成为被攻击的弱点
*/

// DESEncrypt ECB加密
func DESEncrypt(text string, key []byte) (string, error) {
	src := []byte(text)
	block, err := des.NewCipher(key) // max size 8
	//block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}
	bs := block.BlockSize()
	src = ZeroPadding(src, bs)
	if len(src)%bs != 0 {
		return "", errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	return hex.EncodeToString(out), nil
}

// DESDecrypt ECB解密
func DESDecrypt(decrypted string, key []byte) (string, error) {
	src, err := hex.DecodeString(decrypted)
	if err != nil {
		return "", err
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	out := make([]byte, len(src))
	dst := out
	bs := block.BlockSize()
	if len(src)%bs != 0 {
		return "", errors.New("crypto/cipher: input not full blocks")
	}
	for len(src) > 0 {
		block.Decrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	out = ZeroUnPadding(out)
	return string(out), nil
}

func DESCheck(password, encrypted, salt string) (bool, error) {
	pwd, err := DESDecrypt(encrypted, []byte(salt))
	if err != nil {
		return false, err
	}
	return password == pwd, nil
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}
