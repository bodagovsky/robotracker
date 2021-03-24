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
	usersPool []*user
	mu        sync.Mutex
}

type user struct {
	id int
	ts int64
}

var pool = &Pool{
	usersPool: make([]*user, 0),
	mu:        sync.Mutex{},
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
		u.id, err = strconv.Atoi(strings.Join(userID, ""))
		if err != nil {
			log.Println(err)
			return
		}
		u.ts = time.Now().Unix()
		pool.mu.Lock()
		pool.usersPool = append(pool.usersPool, u)
		pool.mu.Unlock()
	}
}

func count(w http.ResponseWriter, r *http.Request) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	length := len(pool.usersPool)
	if length == 0 {
		w.Write([]byte("0"))
		return
	}
	minuteAgo := time.Now().Unix() - 60
	if pool.usersPool[length-1].ts < minuteAgo {
		w.Write([]byte("0"))
		return
	}
	n := 0
	if pool.usersPool[0].ts < minuteAgo {
		for {
			if (pool.usersPool[n].ts == minuteAgo) || (minuteAgo < pool.usersPool[n].ts && minuteAgo > pool.usersPool[n-1].ts) {
				break
			} else if minuteAgo < pool.usersPool[n].ts {
				n /= 2
			} else if minuteAgo > pool.usersPool[n].ts {
				n += n / 2
			}
		}
	}
	lastMinuteUsers := pool.usersPool[n:]
	var usersCount = make(map[int]int)
	var robots int
	for _, user := range lastMinuteUsers {
		usersCount[user.id]++
		if usersCount[user.id] == 100 {
			robots += 1
		}
	}

	w.Write([]byte(strconv.Itoa(robots)))
}
