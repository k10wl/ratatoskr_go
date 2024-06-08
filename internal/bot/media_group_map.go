package bot

import "sync"

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
	mgm.hashMap[key] = append(mgm.hashMap[key], value)
}

func (mgm *mediaGroupMap) remove(key string) {
	mgm.mu.Lock()
	defer mgm.mu.Unlock()
	delete(mgm.hashMap, key)
}
