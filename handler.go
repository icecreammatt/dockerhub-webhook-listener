package listener

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

type CallbackState string

const (
	Success CallbackState = "success"
	Failure CallbackState = "failure"
	Error   CallbackState = "error"
)

type CallbackPayload struct {
	State       CallbackState `json:"state"`
	Description string        `json:"description"`
	Context     string        `json:"context"`
	TargetUrl   string        `json:"target_url"`
}

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
	out, err := exec.Command("./reload.sh", msg.Repository.RepoName).Output()

	var callbackState CallbackState
	callbackPayload := CallbackPayload{
		State: callbackState,
		//TODO: Limit description to 255 characters
		Description: string(out),
		Context:     "Webhook Listener",
		TargetUrl:   "http://ci.acme.com/results/afd339c1c3d27",
	}
	if err != nil {
		log.Println("ERROR EXECUTING COMMAND IN RELOAD HANDLER!!")
		log.Println(err)
		callbackPayload.State = Error
		callbackPayload.Description = err.Error()
	} else {
		callbackPayload.State = Success
	}

	if msg.CallbackUrl != "" {
		res, err := json.Marshal(callbackPayload)
		if err != nil {
			fmt.Println("error marshaling callbackPayload")
		}
		_, err = http.Post(msg.CallbackUrl, "application/json", bytes.NewBuffer(res))
		log.Println("calling callback URL:", msg.CallbackUrl)
		if err != nil {
			fmt.Println("Error posting to callbackUrl", err.Error())
		}
	}
	log.Println("output of reload.sh is", string(out))
}

func MsgHandlers() Registry {
	var handlers Registry

	handlers.Add((&Logger{}).Call)
	handlers.Add(reloadHandler)

	return handlers
}
