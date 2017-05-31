package login

import (
	"net/http"

	"github.com/bradacina/inventory/routing"
)

func (lh *LoginHandler) Logout(w http.ResponseWriter, r *http.Request) {
	lh.CookieAuth.DeleteLoginCookie(w)
	http.Redirect(w, r, routing.RouteLogin, http.StatusSeeOther)
}
