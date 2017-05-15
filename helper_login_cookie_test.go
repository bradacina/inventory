package main

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

type testResponseWriter struct {
	body   []byte
	header map[string][]string
}

func TestSerializeDeserialize(t *testing.T) {

	for i := 1; i <= 20; i++ {
		var s string
		for j := 0; j < i; j++ {
			s = s + "x"
		}
		l := loginInfo{s}
		serialized := serializeLoginInfo(&l)
		result := deserializeLoginInfo(serialized)

		if result.Username != l.Username {
			t.Error("The logininfo " + strconv.Itoa(i) + " did not serialize/deserialize correctly")
		}
	}
}

func TestEncryptDecrypt(t *testing.T) {

	ch := cookieHelper{"test.com", []byte("1234567890123456"), []byte("1234567890123456")}

	var testString string

	for i := 16; i < 300; i += 16 {
		testString = testString + "AxxxxxxxxxxxxxxA"
		encrypted := ch.encrypt([]byte(testString))
		decrypted := ch.decrypt(encrypted)

		if string(decrypted) != testString {
			t.Error("Encrypt/Decrypt failed for ", testString)
		}
	}
}

func TestSignVerify(t *testing.T) {
	ch := cookieHelper{"test.com", []byte("1234567890123456"), []byte("1234567890123456")}

	var testString string

	for i := 16; i < 300; i += 16 {
		testString = testString + "AxxxxxxxxxxxxxxA"
		signature := ch.sign([]byte(testString))
		isAuthentic := ch.verifySignature([]byte(testString), signature)

		if !isAuthentic {
			t.Error("Sign/Verify failed for ", testString)
		}
	}
}

func TestEncodeDecode(t *testing.T) {
	ch := cookieHelper{}

	cipher := "test cypher"
	sig := "test signature"
	encoded := ch.encode([]byte(cipher), []byte(sig))
	decCipher, decSig := ch.decode(encoded)

	if cipher != string(decCipher) || sig != string(decSig) {
		t.Error("Encode/Decode failed for", cipher, sig)
	}
}

func TestSetCookie(t *testing.T) {
	ch := cookieHelper{"test.com", []byte("1234567890123456"), []byte("1234567890123456")}

	w := &testResponseWriter{header: make(map[string][]string)}
	login := &loginInfo{"x"}
	ch.setLoginCookie(w, login)

	if val, ok := w.header["Set-Cookie"]; ok {
		if !strings.Contains(val[0], "Inventory987-Login") {
			t.Error("Cannot find name of cookie in Set-Cookie")
		}
	} else {
		t.Error("Did not Set-Cookie")
	}
}

func (w *testResponseWriter) Header() http.Header {
	return w.header
}

func (w *testResponseWriter) Write(content []byte) (int, error) {
	w.body = append(w.body, content...)
	return len(content), nil
}

func (w *testResponseWriter) WriteHeader(statusCode int) {

}
