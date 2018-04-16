package cache

import (
	"container/list"
	"fmt"
	"sync"
)

type lruItem struct {
	Key   interface{}
	Value interface{}
}

type LRUCache struct {
	amu      sync.Mutex
	rmu      sync.Mutex
	list     *list.List
	table    map[interface{}]*list.Element
	capacity int
}

func NewLRUCache(cap int) *LRUCache {
	return &LRUCache{
		capacity: cap,
		list:     list.New(),
		table:    make(map[interface{}]*list.Element),
	}
}

func (l *LRUCache) Count() int {
	return l.list.Len()
}

func (l *LRUCache) Contains(key interface{}) bool {
	_, ok := l.table[key]
	return ok
}

func (l *LRUCache) Get(key interface{}) interface{} {
	l.rmu.Lock()
	defer l.rmu.Unlock()
	if element, ok := l.table[key]; ok {
		l.list.MoveToBack(element)
		return element.Value.(*lruItem).Value
	}
	return nil
}

func (l *LRUCache) Set(key, value interface{}) {
	l.amu.Lock()
	defer l.amu.Unlock()
	if l.capacity > 0 && l.Count() == l.capacity && !l.Contains(key) {
		l.Remove(l.list.Front().Value.(*lruItem).Key)
	}
	if element, ok := l.table[key]; ok {
		element.Value.(*lruItem).Value = value
	} else {
		element := l.list.PushBack(&lruItem{Key: key, Value: value})
		l.table[key] = element
	}
}

func (l *LRUCache) Remove(key interface{}) {
	l.rmu.Lock()
	defer l.rmu.Unlock()
	if element, ok := l.table[key]; ok {
		l.list.Remove(element)
		delete(l.table, key)
	}
}

func (l *LRUCache) Traverse() {
	data := l.list.Front()
	for {
		fmt.Println(data.Value.(*lruItem))
		if data.Next() != nil {
			data = data.Next()
		} else {
			return
		}
	}
}
