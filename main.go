package main

import (
	"fmt"
	"net/http"
	"sync"
	"flag"
)

type MessageQueue struct {
	mu sync.Mutex
	messages []string
}

func (q *MessageQueue) Enqueue(message string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.messages = append(q.messages, message)
}

func (q *MessageQueue) Dequeue() string {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.messages) == 0 {
		return ""
	}
	message := q.messages[0]
	q.messages = q.messages[1:]
	return message
}

func (q *MessageQueue) size() string {
	q.mu.Lock()
	defer q.mu.Unlock()
	return fmt.Sprintf("%d", len(q.messages))
}

func main() {
	port := flag.Int("p", 8080, "port to listen on")
	flag.Parse()

	queue := &MessageQueue{}

	http.HandleFunc("/enqueue", func(w http.ResponseWriter, r *http.Request) {
		message := r.FormValue("message")
		queue.Enqueue(message)
		fmt.Fprintf(w, "OK")
	})

	http.HandleFunc("/dequeue", func(w http.ResponseWriter, r *http.Request) {
		message := queue.Dequeue()
		fmt.Fprintf(w, message)
	})

	http.HandleFunc("/size", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, queue.size())
	})

	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)

}