package lfuCache

import (
    "fmt"
    "log"
    "math/rand"
    "testing"
)

var lfuCache *LfuCache

func Initialize(size int) *LfuCache {
    if lfuCache == nil {
        lfuCache = NewLfuCache(size)
    }
    return lfuCache
}

func TestAdd(t *testing.T) {

    Initialize(10)
    lfuCache.Add(99, "AddEntryData")
    data, err := lfuCache.Get(99)
    if err != nil {
       t.Errorf("TestAdd: key %d not found after Add()!!", 99)
    }

    data, ok := data.(string)
    if !ok || data != "AddEntryData" {
        t.Errorf("TestAdd: key %d contains invalid data %s !!", 99, data)
    }
}
func TestDelete(t *testing.T) {

    Initialize(10)
    lfuCache.Add(99, "AddEntryData")
    data, err := lfuCache.Get(99)
    if err != nil {
        t.Errorf("TestAdd: key %d not found after Add()!!", 99)
    }

    data, ok := data.(string)
    if !ok || data != "AddEntryData" {
        t.Errorf("TestAdd: key %d contains invalid data %s !!", 99, data)
    }

    data, err = lfuCache.Delete(99)
    if err != nil {
        t.Errorf("TestAdd: key %d not found after Add()!!", 99)
    }

    data, ok = data.(string)
    if !ok || data != "AddEntryData" {
        t.Errorf("TestAdd: key %d contains invalid data %s !!", 99, data)
    }

    data, err = lfuCache.Get(99)
    if err == nil {
        t.Errorf("TestAdd: key %d still found after Delete() with data %s !!",
            99, data.(string))
    }
}


func TestAddAfterCapacity(t *testing.T) {
    size := 10
    Initialize(size)
    for i := 0; i < size; i++ {
        data := fmt.Sprintf("Data-%d", i+100)
        lfuCache.Add(i+100, data)
    }

    for i := 0; i < size-1; i++ {
        _, _ = lfuCache.Get(i+100)
    }

    new_key := size+100
    new_data := fmt.Sprintf("Data-%d", new_key)
    lfuCache.Add(new_key, new_data)

    d, err := lfuCache.Get(new_key-1)
    if err == nil {
        t.Errorf("TestAddAfterCapacity: key %d still found after adding a newer key!! with data %s !!",
            new_key-1, d.(string))
    }

    d, err = lfuCache.Get(new_key)
    if err != nil {
        t.Errorf("TestAddAfterCapacity: Newer key %d not found after adding it!!", new_key)
    }
    data, ok := d.(string)
    if !ok || data != new_data {
        t.Errorf("TestAdd: key %d contains invalid data %s !!", new_key, data)
    }
}
func TestScaleAddAfterCapacity(t *testing.T) {
    size := 100 * 1024 * 1024
    lfuCache = nil
    Initialize(size)
    offset := 0
    for i := 0; i < size; i++ {
        data := fmt.Sprintf("Data-%d", i+offset)
        lfuCache.Add(i+offset, data)
    }

    rand_int := rand.Intn(size)
    for i := 0; i < rand_int; i++ {
        _, _ = lfuCache.Get(i+offset)
    }
    for i := rand_int+1; i < size; i++ {
        _, _ = lfuCache.Get(i+offset)
    }

    log.Printf("TestScaleAddAfterCapacity: rand_int=%d\n", rand_int)
    new_key := size+offset
    new_data := fmt.Sprintf("Data-%d", new_key+offset)
    lfuCache.Add(new_key, new_data)

    d, err := lfuCache.Get(rand_int+offset)
    if err == nil {
        t.Errorf("TestAddAfterCapacity: key %d still found after adding a newer key!! with data %s !!",
            rand_int+offset, d.(string))
    }

    d, err = lfuCache.Get(new_key)
    if err != nil {
        t.Errorf("TestAddAfterCapacity: Newer key %d not found after adding it!!", new_key)
    }
    data, ok := d.(string)
    if !ok || data != new_data {
        t.Errorf("TestAdd: key %d contains invalid data %s !!", new_key, data)
    }
}
