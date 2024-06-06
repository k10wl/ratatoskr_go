package bot

import (
	"reflect"
	"testing"
)

func TestMediaQueue(t *testing.T) {
	q := newMediaQueue()

	q.pop() // should not explode if list is empty

	peek := *q.peek()
	if !reflect.DeepEqual(peek, []mediaInfo{}) {
		t.Errorf("initial queue is not empty: %+v", peek)
	}

	q.push(&mediaInfo{messageID: "1", mediaFileID: "1"})
	peek = *q.peek()
	if !reflect.DeepEqual(peek, []mediaInfo{{messageID: "1", mediaFileID: "1"}}) {
		t.Errorf("did not push element into queue: %+v", peek)
	}

	q.push(&mediaInfo{messageID: "2", mediaFileID: "2"})
	peek = *q.peek()
	if !reflect.DeepEqual(
		peek,
		[]mediaInfo{
			{messageID: "1", mediaFileID: "1"},
			{messageID: "2", mediaFileID: "2"},
		},
	) {
		t.Errorf("did not push element into queue: %+v", peek)
	}

	head := q.pop()
	expected := mediaInfo{messageID: "1", mediaFileID: "1"}
	if !reflect.DeepEqual(*head, expected) {
		t.Errorf("popped unexpected element\nexpected: %+v\nactual:   %+v",
			expected,
			head,
		)
	}

	peek = *q.peek()
	if !reflect.DeepEqual(
		peek,
		[]mediaInfo{
			{messageID: "2", mediaFileID: "2"},
		},
	) {
		t.Errorf("did not popped element from queue: %+v", peek)
	}

	head = q.pop()
	expected = mediaInfo{messageID: "2", mediaFileID: "2"}
	if !reflect.DeepEqual(*head, expected) {
		t.Errorf("popped unexpected element\nexpected: %+v\nactual:   %+v",
			expected,
			head,
		)
	}

	peek = *q.peek()
	if len(peek) != 0 {
		t.Errorf("did not pop last element, actual: %+v", peek)
	}

	q.pop() // should not explode if list is empty
}
