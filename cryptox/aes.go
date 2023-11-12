package cryptox

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

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
