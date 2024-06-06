package bot

import "container/list"

type mediaInfo struct {
	messageID   string
	mediaFileID string
}

type mediaQueue struct {
	list *list.List
}

func newMediaQueue() *mediaQueue {
	return &mediaQueue{
		list: list.New(),
	}
}

func (m *mediaQueue) push(node *mediaInfo) {
	m.list.PushBack(node)
}

func (m *mediaQueue) pop() *mediaInfo {
	if m.list.Len() == 0 {
		return nil
	}
	val := m.list.Remove(m.list.Front()).(*mediaInfo)
	return val
}

func (m *mediaQueue) peek() *[]mediaInfo {
	if m.list.Len() == 0 {
		return &[]mediaInfo{}
	}
	media := []mediaInfo{}
	node := m.list.Front()
	for node != nil {
		media = append(media, *node.Value.(*mediaInfo))
		node = node.Next()
	}
	return &media
}
