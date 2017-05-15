package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	LoginCookieName = "Inventory987-Login"
	CookieLifetime  = 3600 // 1 hour
	CookieExpire    = -1

	Padding byte = 255
)

type cookieHelper struct {
	domain  string
	hmacKey []byte
	aesKey  []byte
}

type loginInfo struct {
	Username string
}

func serializeLoginInfo(info *loginInfo) []byte {
	serialized := []byte(info.Username)

	remainder := len(serialized) % aes.BlockSize
	if remainder == 0 {
		return serialized
	}

	needed := aes.BlockSize - remainder
	initialLen := len(serialized)

	serialized = append(serialized, make([]byte, needed)...)
	for i := 0; i < needed; i++ {
		serialized[initialLen+i] = Padding
	}

	return serialized
}

func deserializeLoginInfo(serialized []byte) *loginInfo {
	for i := len(serialized) - 1; i >= 0; i-- {
		if serialized[i] != Padding {
			serialized = serialized[:i+1]
			break
		}
	}

	return &loginInfo{Username: string(serialized)}
}

func (ch *cookieHelper) encrypt(content []byte) []byte {

	if len(content)%aes.BlockSize != 0 {
		log.Panic("Content to encrypt is not a multiple of aes.BlockSize")
	}

	block, err := aes.NewCipher(ch.aesKey)
	if err != nil {
		log.Panicln("Could not create aes cipher", err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(content))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Panicln("Could not create iv", err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], content)

	return ciphertext
}

func (ch *cookieHelper) decrypt(ciphertext []byte) []byte {
	block, err := aes.NewCipher(ch.aesKey)
	if err != nil {
		log.Panicln("Could not create aes cipher", err)
	}

	if len(ciphertext) < aes.BlockSize {
		log.Panicln("Ciphertext is less than aes.BlockSize")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		log.Panicln("Ciphertext length is not a multiple of aes.BlockSize")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return ciphertext
}

func (ch *cookieHelper) sign(ciphertext []byte) []byte {
	mac := hmac.New(sha256.New, ch.hmacKey)
	mac.Write(ciphertext)
	signature := mac.Sum(nil)

	return signature
}

func (ch *cookieHelper) verifySignature(ciphertext, expectedSig []byte) bool {
	mac := hmac.New(sha256.New, ch.hmacKey)
	mac.Write(ciphertext)
	signature := mac.Sum(nil)

	return hmac.Equal(signature, expectedSig)
}

func (ch *cookieHelper) encryptAndSign(content []byte) ([]byte, []byte) {
	encrypted := ch.encrypt(content)
	signed := ch.sign(encrypted)

	return encrypted, signed
}

func (ch *cookieHelper) verifySignatureAndDecrypt(ciphertext, signature []byte) []byte {
	if !ch.verifySignature(ciphertext, signature) {
		log.Panicln("Cookie signature verification failed")
	}

	decrypted := ch.decrypt(ciphertext)
	return decrypted
}

func (ch *cookieHelper) encode(ciphertext, signature []byte) string {
	encodedSignature := base64.StdEncoding.EncodeToString(signature)
	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)

	return encodedCiphertext + "." + encodedSignature
}

func (ch *cookieHelper) decode(cookie string) ([]byte, []byte) {
	idx := strings.IndexByte(cookie, '.')
	if idx == -1 {
		log.Panicln("Cookie was not in the correct format")
	}

	encodedCipher := cookie[:idx]
	encodedSig := cookie[idx+1:]

	signature, err := base64.StdEncoding.DecodeString(encodedSig)
	if err != nil {
		log.Panicln(err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encodedCipher)
	if err != nil {
		log.Panicln(err)
	}

	return ciphertext, signature
}

func (ch *cookieHelper) setLoginCookie(w http.ResponseWriter, loginInfo *loginInfo) {

	value := ch.encode(ch.encryptAndSign(serializeLoginInfo(loginInfo)))

	cookie := http.Cookie{
		Domain:   ch.domain,
		MaxAge:   CookieLifetime,
		HttpOnly: true,
		Name:     LoginCookieName,
		Secure:   true,
		Value:    value}

	http.SetCookie(w, &cookie)
}

func (ch *cookieHelper) getLoginCookie(r *http.Request) *loginInfo {
	cookie, err := r.Cookie(LoginCookieName)
	if err != nil {
		log.Panicln("Request did not contain the login cookie")
	}

	value := deserializeLoginInfo(ch.verifySignatureAndDecrypt(ch.decode(cookie.Value)))

	return value
}

func (ch *cookieHelper) deleteLoginCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Domain:   ch.domain,
		MaxAge:   CookieExpire,
		HttpOnly: true,
		Name:     LoginCookieName,
		Secure:   true}

	http.SetCookie(w, &cookie)
}
