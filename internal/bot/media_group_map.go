package bot

import (
	"slices"
	"sync"
)

type item struct {
	mediaType string
	fileID    string
	messageID int64
}

type mediaGroupMap struct {
	mu      sync.Mutex
	hashMap map[string][]item
}

func newMediaGroupMap() *mediaGroupMap {
	return &mediaGroupMap{
		hashMap: map[string][]item{},
	}
}

func (mgm *mediaGroupMap) add(key string, value item) {
	mgm.mu.Lock()
	defer mgm.mu.Unlock()
	items := mgm.hashMap[key]
	i := slices.IndexFunc(items, func(el item) bool {
		if el.messageID > value.messageID {
			return true
		}
		return false
	})
	if i == -1 {
		i = len(items)
	}
	mgm.hashMap[key] = slices.Insert(items, i, value)
}

func (mgm *mediaGroupMap) remove(key string) {
	mgm.mu.Lock()
	defer mgm.mu.Unlock()
	delete(mgm.hashMap, key)
}

func (mgm *mediaGroupMap) get(key string) []item {
	mgm.mu.Lock()
	defer mgm.mu.Unlock()
	return mgm.hashMap[key]
}
