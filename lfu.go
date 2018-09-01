package cache

import (
	"container/heap"
	"sync"
	"time"
)

type lfuItem struct {
	Key       interface{}
	Value     interface{}
	Frequence int
	Index     int
	Date      int64
}

type LFUCache struct {
	amu      sync.Mutex
	rmu      sync.Mutex
	table    map[interface{}]*lfuItem
	list     LS
	capacity int
}

func NewLFUCache(cap int) *LFUCache {

	return &LFUCache{
		table:    make(map[interface{}]*lfuItem, cap),
		list:     make([]*lfuItem, 0, cap),
		capacity: cap,
	}

}
func (l *LFUCache) Count() int {
	return len(l.list)
}

func (l *LFUCache) Contains(key interface{}) bool {
	_, ok := l.table[key]
	return ok
}

func (l *LFUCache) Get(key int) interface{} {
	l.rmu.Lock()
	defer l.rmu.Unlock()
	if element, ok := l.table[key]; ok {
		l.list.update(element)
		return element.Value
	}
	return nil
}

func (l *LFUCache) Put(key, value interface{}) {
	if l.capacity <= 0 {
		return
	}
	if element, ok := l.table[key]; ok {
		l.table[key].Value = value
		l.list.update(element)
		return
	}

	element := &lfuItem{Key: key, Value: value}
	if len(l.list) == l.capacity {
		temp := heap.Pop(&l.list).(*lfuItem)
		delete(l.table, temp.Key)
	}
	l.table[key] = element
	heap.Push(&l.list, element)
}

//sort

type LS []*lfuItem

func (ls LS) Len() int {
	return len(ls)
}

func (ls LS) Less(i, j int) bool {
	if ls[i].Frequence == ls[j].Frequence {
		return ls[i].Date < ls[j].Date
	}
	return ls[i].Frequence < ls[j].Frequence
}

func (ls LS) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
	ls[i].Index = i
	ls[j].Index = j
}

// heap

func (ls *LS) Push(i interface{}) {
	l := len(*ls)
	e := i.(*lfuItem)
	e.Index = l
	e.Date = time.Now().Unix()
	*ls = append(*ls, e)
}

func (ls *LS) Pop() interface{} {
	o := *ls
	l := len(o)
	e := o[l-1]
	e.Index = -1
	*ls = o[0 : l-1]
	return e
}

func (ls *LS) update(li *lfuItem) {
	li.Frequence++
	li.Date = time.Now().Unix()
	heap.Fix(ls, li.Index)
}
