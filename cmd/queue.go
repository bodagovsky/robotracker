package main

import "time"

type node struct {
	u    *user
	next *node
}

type userQueue struct {
	robots   int
	usersMap map[string]int
	head     *node
	tail     *node
}

func (u *userQueue) enqueue(usr *user) {
	if u.head == nil || u.tail == nil {
		u.head = &node{u: usr}
		u.tail = u.head
		return
	}
	u.usersMap[usr.id]++
	if u.usersMap[usr.id] > 100 {
		u.robots++
	}
	u.tail.next = &node{u: usr}
	u.tail = u.tail.next
	minuteAgo := usr.ts - 60
	if u.head.u.ts < minuteAgo {
		u.usersMap[u.head.u.id]--
		if u.usersMap[u.head.u.id] <= 100 {
			u.robots--
		}
		u.head = u.head.next
	}
	return

}

func (u *userQueue) count() int {
	ts := time.Now().Unix()
	for u.head != nil && u.head.u.ts < ts {
		u.usersMap[u.head.u.id]--
		if u.usersMap[u.head.u.id] < 100 {
			u.robots--
		}
		u.head = u.head.next
	}
	return u.robots
}
