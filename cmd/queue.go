package main

import "time"

type node struct {
	u    *user
	next *node
}

type userQueue struct {
	head *node
	tail *node
}

func (u *userQueue) enqueue(usr *user) {
	if u.head == nil || u.tail == nil {
		u.head = &node{u: usr}
		u.tail = u.head
		return
	}
	u.tail.next = &node{u: usr}
	u.tail = u.tail.next
	minuteAgo := time.Now().Unix() - 60
	if u.head.u.ts < minuteAgo {
		u.head = u.head.next
	}
	return

}

func (u userQueue) count() int {
	var counter = make(map[string]int)
	var robots int
	for u.head != nil {
		counter[u.head.u.id]++
		if counter[u.head.u.id] == 100 {
			robots++
		}
		u.head = u.head.next
	}
	return robots
}
