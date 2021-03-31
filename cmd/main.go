package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Pool struct {
	usersPool *userQueue
	mu        *sync.Mutex
}

type user struct {
	id string
	ts int64
}

var pool = &Pool{
	usersPool: &userQueue{
		usersMap: make(map[string]int),
	},
	mu:        &sync.Mutex{},
}

func main() {
	http.HandleFunc("/", enqueue)
	http.HandleFunc("/count", count)
	http.ListenAndServe("127.0.0.1:8080", nil)
}

func enqueue(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	if userID, ok := r.Form["user_id"]; ok {
		var u = new(user)
		u.id = strings.Join(userID, "")
		u.ts = time.Now().Unix()
		pool.mu.Lock()
		pool.usersPool.enqueue(u)
		pool.mu.Unlock()
	}
}

func count(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(strconv.Itoa(pool.usersPool.robots)))
}

