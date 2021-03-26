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
	id string
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
		u.id = strings.Join(userID, "")
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
	n := binarySearch(pool.usersPool, minuteAgo)

	pool.usersPool = pool.usersPool[n:]
	var usersCount = make(map[string]int)
	var robots int
	for _, user := range pool.usersPool {
		usersCount[user.id]++
		if usersCount[user.id] == 100 {
			robots += 1
		}
	}

	w.Write([]byte(strconv.Itoa(robots)))
}

func binarySearch(arr []*user, element int64) int {
	i := 0
	j := len(arr) - 1
	if element < arr[i].ts {
		return i
	}
	if element > arr[j].ts {
		return j
	}
	for i != j {
		n := (i+j)/2 + 1
		if arr[n].ts == element || (element > arr[n-1].ts && element < arr[n].ts) {
			//убедимся, что при наличии множества одинаковых значений мы выберем все имеющиеся
			for n > 0 && arr[n].ts == arr[n-1].ts {
				n--
			}
			return n
		}
		if element < arr[n].ts {
			j = n - 1
		} else {
			i = n + 1
		}
	}
	return i
}
