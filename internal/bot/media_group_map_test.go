package bot

import (
	"reflect"
	"testing"
)

func TestMediaGroupMap(t *testing.T) {
	type tc struct {
		name      string
		subject   *mediaGroupMap
		operation func(*mediaGroupMap)
		expected  map[string][]item
	}

	table := []tc{
		{
			name:      "should initialize empty map",
			subject:   newMediaGroupMap(),
			operation: func(mgm *mediaGroupMap) {},
			expected:  map[string][]item{},
		},

		{
			name:    "should add items to the map",
			subject: newMediaGroupMap(),
			operation: func(mgm *mediaGroupMap) {
				mgm.add("1", item{
					messageID: 1,
					fileID:    "file id 1",
					mediaType: "photo",
				})
				mgm.add("1", item{
					messageID: 2,
					fileID:    "file id 2",
					mediaType: "photo",
				})
				mgm.add("2", item{
					messageID: 3,
					fileID:    "file id 3",
					mediaType: "photo",
				})
			},
			expected: map[string][]item{
				"1": {
					{
						messageID: 1,
						fileID:    "file id 1",
						mediaType: "photo",
					},
					{
						messageID: 2,
						fileID:    "file id 2",
						mediaType: "photo",
					},
				},
				"2": {
					{
						messageID: 3,
						fileID:    "file id 3",
						mediaType: "photo",
					},
				},
			},
		},

		{
			name: "should remove items from the map",
			subject: &mediaGroupMap{
				hashMap: map[string][]item{
					"1": {
						{
							messageID: 1,
							fileID:    "file id 1",
							mediaType: "photo",
						},
						{
							messageID: 2,
							fileID:    "file id 2",
							mediaType: "photo",
						},
					},
					"2": {
						{
							messageID: 3,
							fileID:    "file id 3",
							mediaType: "photo",
						},
						{
							messageID: 4,
							fileID:    "file id 4",
							mediaType: "photo",
						},
					},
				},
			},
			operation: func(mgm *mediaGroupMap) {
				mgm.remove("1")
			},
			expected: map[string][]item{
				"2": {
					{
						messageID: 3,
						fileID:    "file id 3",
						mediaType: "photo",
					},
					{
						messageID: 4,
						fileID:    "file id 4",
						mediaType: "photo",
					},
				},
			},
		},
	}

	for _, test := range table {
		test.operation(test.subject)
		if !reflect.DeepEqual(test.expected, test.subject.hashMap) {
			t.Errorf("%s - wrong output\nexpected: %+v\nactual:   %+v",
				test.name,
				test.expected,
				test.subject.hashMap,
			)
		}
	}
}
