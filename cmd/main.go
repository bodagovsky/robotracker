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
	usersPool *UserQueue
	mu        *sync.Mutex
}

type user struct {
	id string
	ts int64
}

type userRecord struct {
	c       int
	isRobot bool
}

var pool = &Pool{
	usersPool: &UserQueue{
		usersMap: make(map[string]*userRecord),
	},
	mu: &sync.Mutex{},
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
		pool.usersPool.Enqueue(u)
		pool.mu.Unlock()
	}
}

func count(w http.ResponseWriter, r *http.Request) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	robots := pool.usersPool.Count()
	w.Write([]byte(strconv.Itoa(robots)))
}

type node struct {
	u    *user
	next *node
}

type UserQueue struct {
	robots   int
	usersMap map[string]*userRecord
	head     *node
	tail     *node
}

func (u *UserQueue) Enqueue(usr *user) {
	if u.head == nil || u.tail == nil {
		u.head = &node{u: usr}
		u.tail = u.head
		return
	}
	if _, ok := u.usersMap[usr.id]; ok {
		u.usersMap[usr.id].c++
	} else {
		u.usersMap[usr.id] = &userRecord{
			c:       1,
			isRobot: false,
		}
	}

	if u.usersMap[usr.id].c > 100 {
		if !u.usersMap[usr.id].isRobot {
			u.usersMap[usr.id].isRobot = true
			u.robots++
		}
	}
	u.tail.next = &node{u: usr}
	u.tail = u.tail.next
	minuteAgo := usr.ts - 60
	for u.head != nil && u.head.u.ts < minuteAgo {
		u.usersMap[u.head.u.id].c--
		if u.usersMap[u.head.u.id].c <= 100 {
			if u.usersMap[usr.id].isRobot {
				u.usersMap[usr.id].isRobot = false
				u.robots--
			}
		}
		u.head = u.head.next
	}
	return

}

func (u *UserQueue) Count() int {
	minuteAgo := time.Now().Unix() - 60
	for u.head != nil && u.head.u.ts < minuteAgo {
		u.usersMap[u.head.u.id].c--
		if u.usersMap[u.head.u.id].c <= 100 {
			if u.usersMap[u.head.u.id].isRobot {
				u.usersMap[u.head.u.id].isRobot = false
				u.robots--
			}
		}
		u.head = u.head.next
	}
	return u.robots
}
