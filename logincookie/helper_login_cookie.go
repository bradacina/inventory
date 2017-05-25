package logincookie

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	loginCookieName = "Inventory987-Login"
	cookieLifetime  = 3600 // 1 hour
	cookieExpire    = -1

	padding byte = 255
)

var (
	ErrorNoLoginCookie = errors.New("No Login Cookie")
)

type CookieAuthentication struct {
	domain  string
	hmacKey []byte
	aesKey  []byte
}

type LoginInfo struct {
	ID       int
	Username string
	IsAdmin  bool
}

func NewCookieAuthentication(
	domain string, hmacKey []byte, aesKey []byte) *CookieAuthentication {
	return &CookieAuthentication{domain, hmacKey, aesKey}
}

func serializeLoginInfo(info *LoginInfo) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(info)
	if err != nil {
		log.Panic("Cannot serialize logininfo")
	}

	serialized := buf.Bytes()
	remainder := len(serialized) % aes.BlockSize
	if remainder == 0 {
		return serialized
	}

	needed := aes.BlockSize - remainder
	initialLen := len(serialized)

	serialized = append(serialized, make([]byte, needed)...)
	for i := 0; i < needed; i++ {
		serialized[initialLen+i] = padding
	}

	return serialized
}

func deserializeLoginInfo(serialized []byte) *LoginInfo {
	for i := len(serialized) - 1; i >= 0; i-- {
		if serialized[i] != padding {
			serialized = serialized[:i+1]
			break
		}
	}

	buf := bytes.NewBuffer(serialized)
	dec := gob.NewDecoder(buf)
	var loginInfo LoginInfo
	dec.Decode(&loginInfo)

	return &loginInfo
}

func (ch *CookieAuthentication) encrypt(content []byte) []byte {

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

func (ch *CookieAuthentication) decrypt(ciphertext []byte) []byte {
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

func (ch *CookieAuthentication) sign(ciphertext []byte) []byte {
	mac := hmac.New(sha256.New, ch.hmacKey)
	mac.Write(ciphertext)
	signature := mac.Sum(nil)

	return signature
}

func (ch *CookieAuthentication) verifySignature(ciphertext, expectedSig []byte) bool {
	mac := hmac.New(sha256.New, ch.hmacKey)
	mac.Write(ciphertext)
	signature := mac.Sum(nil)

	return hmac.Equal(signature, expectedSig)
}

func (ch *CookieAuthentication) encryptAndSign(content []byte) ([]byte, []byte) {
	encrypted := ch.encrypt(content)
	signed := ch.sign(encrypted)

	return encrypted, signed
}

func (ch *CookieAuthentication) verifySignatureAndDecrypt(ciphertext, signature []byte) []byte {
	if !ch.verifySignature(ciphertext, signature) {
		log.Panicln("Cookie signature verification failed")
	}

	decrypted := ch.decrypt(ciphertext)
	return decrypted
}

func (ch *CookieAuthentication) encode(ciphertext, signature []byte) string {
	encodedSignature := base64.StdEncoding.EncodeToString(signature)
	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)

	return encodedCiphertext + "." + encodedSignature
}

func (ch *CookieAuthentication) decode(cookie string) ([]byte, []byte) {
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

func (ch *CookieAuthentication) SetLoginCookie(w http.ResponseWriter, loginInfo *LoginInfo) {

	value := ch.encode(ch.encryptAndSign(serializeLoginInfo(loginInfo)))

	cookie := http.Cookie{
		Domain:   ch.domain,
		MaxAge:   cookieLifetime,
		HttpOnly: true,
		Name:     loginCookieName,
		Secure:   true,
		Value:    value}

	http.SetCookie(w, &cookie)
}

func (ch *CookieAuthentication) GetLoginCookie(r *http.Request) (li *LoginInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			li = nil
			err = ErrorNoLoginCookie
		}
	}()

	cookie, err := r.Cookie(loginCookieName)
	if err != nil {
		return nil, err
	}

	value := deserializeLoginInfo(ch.verifySignatureAndDecrypt(ch.decode(cookie.Value)))

	return value, nil
}

func (ch *CookieAuthentication) DeleteLoginCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Domain:   ch.domain,
		MaxAge:   cookieExpire,
		HttpOnly: true,
		Name:     loginCookieName,
		Secure:   true}

	http.SetCookie(w, &cookie)
}
