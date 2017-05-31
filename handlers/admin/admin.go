package admin

import "github.com/bradacina/inventory/deps"

type AdminHandler struct {
	*deps.Deps
}

func NewHandler(deps *deps.Deps) *AdminHandler {
	return &AdminHandler{deps}
}
