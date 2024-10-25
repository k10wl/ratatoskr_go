package mongo_db

import (
	"ratatoskr/internal/models"
	"reflect"
	"testing"
)

func TestFixSorting(t *testing.T) {
	group := []models.Group{
		{
			OriginalIndex: 1,
		},
		{
			OriginalIndex: 0,
		},
		{
			OriginalIndex: 5,
		},
		{
			OriginalIndex: 3,
		},
		{
			OriginalIndex: 4,
		},
		{
			OriginalIndex: 2,
		},
	}

	sorted := fixSorting(group)

	afterSorting := make([]int, len(sorted))
	for i, v := range sorted {
		afterSorting[i] = v.OriginalIndex
	}

	if !reflect.DeepEqual([]int{0, 1, 2, 3, 4, 5}, afterSorting) {
		t.Errorf("Failed to sort groups by original indexes\n")
	}
}
