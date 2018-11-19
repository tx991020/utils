// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	// DefaultEvery means the clock time of recycling the expired cache items in memory.
	DefaultEvery = 60 // 1 minute
)

// MemoryItem store memory cache item.
type MemoryItem struct {
	sync.RWMutex
	val         interface{}
	createdTime time.Time
	lifespan    time.Duration
}

func (mi *MemoryItem) isExpire() bool {
	// 0 means forever
	if mi.lifespan == 0 {
		return false
	}
	return time.Now().Sub(mi.createdTime) > mi.lifespan
}

// MemoryCache is Memory cache adapter.
// it contains a RW locker for safe map storage.
type MemoryCache struct {
	sync.RWMutex
	l     Loader
	dur   time.Duration
	items map[string]*MemoryItem
	Every int // run an expiration check Every clock time
}

// NewMemoryCache returns a new MemoryCache.
func NewMemoryCache() Cache {
	cache := MemoryCache{items: make(map[string]*MemoryItem)}
	return &cache
}

// Get cache from memory.
// if non-existed or expired, return nil.
func (bc *MemoryCache) Get(name string) (interface{}, bool) {
	bc.RLock()
	defer bc.RUnlock()

	var ret interface{} = nil
	itm, ok := bc.items[name]
	if !ok || itm.isExpire() {
		ok = false
		bc.RUnlock()
		bc.Lock()
		itm, ok = bc.items[name] // check twice, in case of flood redis.get
		if !ok || itm.isExpire() {
			ok = false
			o, lspan := bc.l.Load(name)
			if o != nil { // refresh
				itm = &MemoryItem{
					val:         o,
					createdTime: time.Now(),
					lifespan:    lspan,
				}
				bc.items[name] = itm
			}
			ret = o
		} else {
			ret = itm.val
		}
		bc.Unlock()
		bc.RLock()
	} else {
		ret = itm.val
	}

	if ret == nil && ok == true {
		fmt.Println("WARN: Cache Get(%v) return with nil but exist!!!! item: %v", itm)
	}
	return ret, ok
}

// GetMulti gets caches from memory.
// if non-existed or expired, return nil.
func (bc *MemoryCache) GetMulti(names []string) ([]interface{}, []bool) {
	var rc []interface{}
	var ec []bool
	for _, name := range names {
		o, e := bc.Get(name)
		rc = append(rc, o)
		ec = append(ec, e)
	}
	return rc, ec
}

// Put cache to memory.
// if lifespan is 0, it will be forever till restart.
func (bc *MemoryCache) Put(name string, value interface{}, lifespan time.Duration) error {
	if value == nil {
		fmt.Println("WARN: Cache Put(%v, nil)", name)
	}
	bc.Lock()
	defer bc.Unlock()
	bc.items[name] = &MemoryItem{
		val:         value,
		createdTime: time.Now(),
		lifespan:    lifespan,
	}
	if bc.l != nil {
		if err := bc.l.Put(name, value); err != nil {
			return err
		}
	}
	return nil
}

// Invalid cache in memory.
func (bc *MemoryCache) Invalid(name string) error {
	bc.Lock()
	defer bc.Unlock()
	if _, ok := bc.items[name]; !ok {
		return errors.New("key not exist")
	}
	delete(bc.items, name)
	if _, ok := bc.items[name]; ok {
		return errors.New("delete key error")
	}
	return nil
}

// Delete cache in memory.
func (bc *MemoryCache) Delete(name string) error {
	bc.Lock()
	defer bc.Unlock()
	if bc.l != nil {
		if err := bc.l.Delete(name); err != nil {
			return err
		}
	}
	if _, ok := bc.items[name]; !ok {
		return errors.New("key not exist")
	}
	delete(bc.items, name)
	if _, ok := bc.items[name]; ok {
		return errors.New("delete key error")
	}
	return nil
}

// Incr increase cache counter in memory.
// it supports int,int32,int64,uint,uint32,uint64.
func (bc *MemoryCache) Incr(key string) error {
	bc.RLock()
	defer bc.RUnlock()
	itm, ok := bc.items[key]
	if !ok {
		return errors.New("key not exist")
	}
	switch itm.val.(type) {
	case int:
		itm.val = itm.val.(int) + 1
	case int32:
		itm.val = itm.val.(int32) + 1
	case int64:
		itm.val = itm.val.(int64) + 1
	case uint:
		itm.val = itm.val.(uint) + 1
	case uint32:
		itm.val = itm.val.(uint32) + 1
	case uint64:
		itm.val = itm.val.(uint64) + 1
	default:
		return errors.New("item val is not (u)int (u)int32 (u)int64")
	}
	return nil
}

// Decr decrease counter in memory.
func (bc *MemoryCache) Decr(key string) error {
	bc.RLock()
	defer bc.RUnlock()
	itm, ok := bc.items[key]
	if !ok {
		return errors.New("key not exist")
	}
	switch itm.val.(type) {
	case int:
		itm.val = itm.val.(int) - 1
	case int64:
		itm.val = itm.val.(int64) - 1
	case int32:
		itm.val = itm.val.(int32) - 1
	case uint:
		if itm.val.(uint) > 0 {
			itm.val = itm.val.(uint) - 1
		} else {
			return errors.New("item val is less than 0")
		}
	case uint32:
		if itm.val.(uint32) > 0 {
			itm.val = itm.val.(uint32) - 1
		} else {
			return errors.New("item val is less than 0")
		}
	case uint64:
		if itm.val.(uint64) > 0 {
			itm.val = itm.val.(uint64) - 1
		} else {
			return errors.New("item val is less than 0")
		}
	default:
		return errors.New("item val is not int int64 int32")
	}
	return nil
}

// IsExist check cache exist in memory.
func (bc *MemoryCache) IsExist(name string) bool {
	bc.RLock()
	defer bc.RUnlock()
	if v, ok := bc.items[name]; ok {
		return !v.isExpire()
	}
	return false
}

// ClearAll will delete all cache in memory.
func (bc *MemoryCache) ClearAll() error {
	bc.Lock()
	defer bc.Unlock()
	bc.items = make(map[string]*MemoryItem)
	return nil
}

// StartAndGC start memory cache. it will check expiration in every clock time.
func (bc *MemoryCache) StartAndGC(config string, l Loader) error {
	var cf map[string]int
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["interval"]; !ok {
		cf = make(map[string]int)
		cf["interval"] = DefaultEvery
	}
	dur := time.Duration(cf["interval"]) * time.Second
	bc.Every = cf["interval"]
	bc.dur = dur
	bc.l = l
	go bc.vaccuum()
	return nil
}

// check expiration.
func (bc *MemoryCache) vaccuum() {
	if bc.Every < 1 {
		return
	}
	for {
		<-time.After(bc.dur)
		if bc.items == nil {
			return
		}
		for name := range bc.items {
			bc.itemExpired(name)
		}
	}
}

// itemExpired returns true if an item is expired.
func (bc *MemoryCache) itemExpired(name string) bool {
	bc.Lock()
	defer bc.Unlock()

	itm, ok := bc.items[name]
	if !ok {
		return true
	}
	if itm.isExpire() {
		delete(bc.items, name)
		return true
	}
	return false
}

func init() {
	Register("memory", NewMemoryCache)
}
