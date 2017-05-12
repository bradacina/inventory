package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
)

const (
	LoginCookieName = "Inventory987-Login"
	CookieLifetime  = 3600 // 1 hour
)

type cookieHelper struct {
	domain string
	key    []byte
}

func (ch *cookieHelper) setLoginCookie(w http.ResponseWriter, username string) {

	mac := hmac.New(sha256.New, ch.key)
	mac.Write([]byte(username))
	encrypted := mac.Sum(nil)
	encUsername := base64.StdEncoding.EncodeToString(encrypted)

	cookie := http.Cookie{
		Domain:   ch.domain,
		MaxAge:   CookieLifetime,
		HttpOnly: true,
		Name:     LoginCookieName,
		Path:     "/",
		Secure:   true,
		Value:    encUsername}

	http.SetCookie(w, &cookie)
}

func (ch *cookieHelper) getLoginCookie(r *http.Request) (string, err) {
	cookie, err := r.Cookie(LoginCookieName)
	if err != nil {
		return nil, err
	}

	if cookie.Domain != ch.domain ||
	cookie.Path != "/" {
		return nil, http.ErrNoCookie
	}

	encrypted, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, err
	}

	mac := hmac.New(sha256.New, ch.key)
	mac.
}

func (ch *cookieHelper) deleteLoginCookie(w http.ResponseWriter) {

}
