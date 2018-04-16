package cache

import (
	"container/list"
	"fmt"
	"sync"
)

type fifoItem struct {
	Key   interface{}
	Value interface{}
}

type FIFOCache struct {
	amu      sync.Mutex
	rmu      sync.Mutex
	list     *list.List
	table    map[interface{}]*list.Element
	capacity int
}

func NewFIFOCache(cap int) *FIFOCache {
	return &FIFOCache{
		capacity: cap,
		list:     list.New(),
		table:    make(map[interface{}]*list.Element),
	}
}

func (f *FIFOCache) Count() int {
	return f.list.Len()
}

func (f *FIFOCache) Contains(key interface{}) bool {
	_, ok := f.table[key]
	return ok
}

func (f *FIFOCache) Get(key interface{}) interface{} {
	if element, ok := f.table[key]; ok {
		return element.Value.(*fifoItem).Value
	}
	return nil
}

func (f *FIFOCache) Set(key, value interface{}) {
	f.amu.Lock()
	defer f.amu.Unlock()
	if f.capacity > 0 && f.Count() == f.capacity && !f.Contains(key) {
		f.Remove(f.list.Front().Value.(*fifoItem).Key)
	}
	if element, ok := f.table[key]; ok {
		element.Value.(*fifoItem).Value = value
	} else {
		element := f.list.PushBack(&fifoItem{Key: key, Value: value})
		f.table[key] = element
	}
}

func (f *FIFOCache) Remove(key interface{}) {
	f.rmu.Lock()
	defer f.rmu.Unlock()
	if element, ok := f.table[key]; ok {
		f.list.Remove(element)
		delete(f.table, key)
	}
}

func (f *FIFOCache) Traverse() {
	data := f.list.Front()
	for {
		fmt.Println(data.Value.(*fifoItem))
		if data.Next() != nil {
			data = data.Next()
		} else {
			return
		}
	}
}
