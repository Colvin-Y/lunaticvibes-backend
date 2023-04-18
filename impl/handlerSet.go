package processor

import (
	"fmt"
	"net/http"
)

type ScoreProcessorHandlerFunc func(w http.ResponseWriter, r *http.Request, sp *ScoreProcessor) error
type SignUpProcessorHandlerFunc func(w http.ResponseWriter, r *http.Request, sp *SignUpProcessor) error

func (sp *ScoreProcessor) Sync(w http.ResponseWriter, r *http.Request, handlers ...ScoreProcessorHandlerFunc) {
	sp.Logger.Info(fmt.Sprintf("start to sync req[%v]", *r))
	for _, handler := range handlers {
		err := handler(w, r, sp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			sp.Logger.Error(err.Error())
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	sp.Logger.Info(fmt.Sprintf("process req[%v] success", *r))
}

func (sup *SignUpProcessor) Sync(w http.ResponseWriter, r *http.Request, handlers ...SignUpProcessorHandlerFunc) {
	sup.Logger.Info(fmt.Sprintf("start to sync req[%v]", *r))
	for _, handler := range handlers {
		err := handler(w, r, sup)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			sup.Logger.Error(err.Error())
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	sup.Logger.Info(fmt.Sprintf("process req[%v] success", *r))
}
