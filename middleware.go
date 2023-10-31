package bongo

import (
	"net/http"
)

func (b *Bongo) LoadSession(next http.Handler) http.Handler {
	return b.Session.LoadAndSave(next)
}
