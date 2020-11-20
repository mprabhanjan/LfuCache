package lfuCache
import (
    "errors"
    "container/heap"
    "log"
    "time"
)

type Node struct {
    key int
    data interface{}
    frequency int
    index int
    tstamp time.Time
}

type PrioQ []*Node

type LfuCache struct {
    keyMap  map[int]*Node
    max_size int
    prioQ PrioQ
}

func (prq PrioQ) Len() int {
    return len(prq)
}

func (prq PrioQ) Less(i, j int) bool {
    if prq[i].frequency < prq[j].frequency {
        return true
    }
    if prq[i].frequency == prq[j].frequency {
        return prq[i].tstamp.Before(prq[j].tstamp)
    }
    return false
}

func (prq PrioQ) Swap(i, j int) {
    prq[i], prq[j] = prq[j], prq[i]
    prq[i].index = i
    prq[j].index = j
}

func (prq *PrioQ) Push(x interface{}) {
    n := len(*prq)
    item := x.(*Node)
    item.index = n
    *prq = append(*prq, item)
}

func (prq *PrioQ) Pop() interface{} {
    old := *prq
    n := len(old)
    item := old[n-1]
    old[n-1] = nil  // avoid memory leak
    item.index = -1 // for safety
    *prq = old[0 : n-1]
    return item
}


func NewLfuCache(max_size int) *LfuCache {
    lfuCache := &LfuCache{
        keyMap: make(map[int]*Node),
        prioQ: make([]*Node, 0),
        max_size: max_size,
    }
    heap.Init(&lfuCache.prioQ)
    return lfuCache
}

func (lfuc *LfuCache) Add(key int, data interface{}) error {

    // If key already exists, update the data and tstamp
    tmp_node, ok := lfuc.keyMap[key]
    if ok {
        tmp_node.data = data
        tmp_node.tstamp = time.Now()
        tmp_node.frequency += 1
        heap.Fix(&lfuc.prioQ, tmp_node.index)
        return nil
    }

    new_node := &Node {
        key: key,
        data: data,
        frequency: 1,
        tstamp: time.Now(),
    }

    if len(lfuc.keyMap) >= lfuc.max_size {
        pop_elem := heap.Pop(&lfuc.prioQ)
        lfu_node, ok := pop_elem.(*Node)
        if !ok {
            log.Fatalf("LfuCache.Add: lfu_node is not of *Node Type!!")
        }
        log.Printf("Cache size at the limit; Dropping key %d with frequency %d",
            lfu_node.key, lfu_node.frequency)
        delete(lfuc.keyMap, lfu_node.key)
    }

    if !(len(lfuc.keyMap) < lfuc.max_size) {
        log.Fatalf("LfuCache.Add: Unexpected - keyMap len is still larger than max_size!!")
    }

    lfuc.keyMap[key] = new_node
    heap.Push(&lfuc.prioQ, new_node)
    return nil
}

func (lfuc *LfuCache) Get(key int) (interface{}, error) {

    key_node, ok := lfuc.keyMap[key]
    if !ok {
        return nil, errors.New("Not found!")
    }
    key_node.frequency += 1
    key_node.tstamp = time.Now()
    heap.Fix(&lfuc.prioQ, key_node.index)
    return key_node.data, nil
}

func (lfuc *LfuCache) Delete(key int) (interface{}, error) {

    key_node, ok := lfuc.keyMap[key]
    if !ok {
        return nil, errors.New("Not found!")
    }
    heap.Remove(&lfuc.prioQ, key_node.index)
    delete(lfuc.keyMap, key)
    return key_node.data, nil
}