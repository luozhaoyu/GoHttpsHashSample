package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Message struct {
	Message string
}

type MessageResponse struct {
	Digest  string `json:"digest,omitempty"`
	Message string `json:"message,omitempty"`
	ErrMsg  string `json:"err_msg,omitempty"`
}

type messageServer struct {
	digestedMessages map[string]string
}

func NewMessageServer() *messageServer {
	s := new(messageServer)
	s.digestedMessages = make(map[string]string)
	return s
}

func (s *messageServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	response := new(MessageResponse)

	if req.Method == "POST" {
		jsonDecoder := json.NewDecoder(req.Body)
		m := new(Message)
		if err := jsonDecoder.Decode(m); err != nil {
			s.writeErrorResponse(w, err, 500)
			return
		}

		response.Digest = fmt.Sprintf("%x", sha256.Sum256([]byte(m.Message)))
		s.digestedMessages[response.Digest] = m.Message
		w.WriteHeader(201)
	} else if req.Method == "GET" {
		digest, err := validateInput(req.URL.Path)
		if err != nil {
			s.writeErrorResponse(w, err, 400)
			return
		}
		message, ok := s.digestedMessages[digest]
		if !ok {
			err := errors.New(fmt.Sprintf("Message Not Found for: %v", s.digestedMessages))
			s.writeErrorResponse(w, err, 404)
			return
		}
		response.Message = message
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(response); err != nil {
		s.writeErrorResponse(w, err, 500)
	}
}

func (s *messageServer) writeErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(err.Error()))
}

// validateInput tries to convert input into digeste [32]byte
func validateInput(path string) (string, error) {
	splits := strings.Split(path, "/")
	if len(splits) != 3 {
		return "", errors.New(fmt.Sprintf("invalid input: %s", path))
	}
	return splits[2], nil
}

func main() {
	server := NewMessageServer()
	http.Handle("/messages", server)
	http.Handle("/messages/", server)
	log.Print("Listening...")
	if err := http.ListenAndServeTLS(":5000", "localhost.crt", "server.key", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
