package cryptox

import (
	"bytes"
	"crypto/des"
	"encoding/hex"
	"errors"
)

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
