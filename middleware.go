package bongo

import (
	"fmt"
	"net/http"
)

func (b *Bongo) LoadSession(next http.Handler) http.Handler {
	fmt.Println("Session loaded")
	return b.Session.LoadAndSave(next)
}
