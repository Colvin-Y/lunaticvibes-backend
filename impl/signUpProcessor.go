package processor

import (
	"net/http"
)

func (sup *SignUpProcessor) SignUpHandlerSet(w http.ResponseWriter, r *http.Request) {
	handlers := []SignUpProcessorHandlerFunc{}
	sup.Sync(w, r, handlers...)
}
