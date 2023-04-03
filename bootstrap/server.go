package bootstrap

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
	"net/http"
	"os"
	"sync"
)

func NewServer(bot *telebot.Bot) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		update := &telebot.Update{}
		if err := json.NewDecoder(r.Body).Decode(update); err != nil {
			logrus.Errorf("Error decoding request body: %v", err)
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}

		mu := sync.Mutex{}
		mu.Lock()
		defer logrus.Info("unlock bot update")
		defer mu.Unlock()

		logrus.Info("lock bot update")
		bot.ProcessUpdate(*update)

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			logrus.Error(err)
			return
		}
	}
	go listenAndServe(handler)
}

func listenAndServe(handler func(w http.ResponseWriter, r *http.Request)) {
	port := os.Getenv("PORT")
	if port == "" {
		logrus.Fatal("port didn't set in your env")
	}

	logrus.Infof("Starting http server on port %s", port)
	if err := http.ListenAndServe(":"+port, http.HandlerFunc(handler)); err != nil {
		logrus.Fatalf("Error while starting or listening server: %v", err)
	}
}
