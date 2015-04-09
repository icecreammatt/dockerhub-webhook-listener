package listener

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

type Handler interface {
	Call(HubMessage)
}

type Logger struct{}

func (l *Logger) Call(msg HubMessage) {
	log.Print(msg)
}

type Registry struct {
	entries []func(HubMessage)
}

func (r *Registry) Add(h func(msg HubMessage)) {
	r.entries = append(r.entries, h)
	return
}

func (r *Registry) Call(msg HubMessage) {
	for _, h := range r.entries {
		go h(msg)
	}
}

func reloadHandler(msg HubMessage) {
	fmt.Println("Reload handler called")
	log.Println("received message to reload ...")
	log.Printf("Callback URL: %s\n", msg.CallbackUrl)
	if msg.CallbackUrl != "" {
		jsonStr := []byte(`{"state":"success"}`)
		http.Post(msg.CallbackUrl, "application/json", bytes.NewBuffer(jsonStr))
		log.Println("calling callback URL:", msg.CallbackUrl)
	}
	out, err := exec.Command("./reload.sh", msg.Repository.RepoName).Output()
	if err != nil {
		log.Println("ERROR EXECUTING COMMAND IN RELOAD HANDLER!!")
		log.Println(err)
		return
	}
	log.Println("output of reload.sh is", string(out))
}

func MsgHandlers() Registry {
	var handlers Registry

	handlers.Add((&Logger{}).Call)
	handlers.Add(reloadHandler)

	return handlers
}
