package structures

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

type Item interface{}

type CircularQueue struct {
	items []Item
	lock  sync.RWMutex
}

type Set struct {
	items map[string]bool
}

type SyncSet struct {
	items map[Item]bool
	lock  sync.RWMutex
}

type SyncMap struct {
	data map[Item]Item
	lock sync.RWMutex
}

func NewQueue() *CircularQueue {
	return &CircularQueue{items: make([]Item, 0)}
}

func NewSet() *Set {
	s := &Set{}
	s.items = make(map[string]bool)
	return s
}

func (s *Set) Insert(value string) {
	s.items[value] = true
}

func (s *Set) Remove(value string) {
	delete(s.items, value)
}

func (s *Set) Contains(value string) bool {
	_, ok := s.items[value]
	return ok
}

func (queue *CircularQueue) Push(item interface{}) {
	queue.lock.Lock()
	queue.items = append(queue.items, item)
	queue.lock.Unlock()
}

func (queue *CircularQueue) InitFromFile(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		queue.Push(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (queue *CircularQueue) PopString() string {
	queue.lock.Lock()
	item := queue.items[0]
	if len(queue.items) > 1 {
		queue.items = queue.items[1:]
		queue.items = append(queue.items, item)
	}
	queue.lock.Unlock()
	return fmt.Sprintf("%v", item)

}

func (queue *CircularQueue) PrintStatus() {
	for _, item := range queue.items[0:] {
		fmt.Println(item)
	}
}

func NewSyncSet() *SyncSet {
	return &SyncSet{items: make(map[Item]bool)}
}

func (set *SyncSet) Insert(item Item) bool {
	set.lock.Lock()
	if set.Contains(item) {
		set.lock.Unlock()
		return false
	}

	set.items[item] = true
	set.lock.Unlock()
	return true
}

func (set *SyncSet) Contains(item Item) bool {
	_, ok := set.items[item]
	return ok
}

// func NewSyncMap() *SyncMap {

// }
