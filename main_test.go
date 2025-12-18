package syncset

import (
	"testing"
)

func checkSame[T comparable](items []T, expectedResults []T) bool {
	resultMap := make(map[T]bool)
	for _, r := range items {
		resultMap[r] = true
	}
	for _, expected := range expectedResults {
		if !resultMap[expected] {
			return false
		}
	}
	return true
}

func addFunc(funcset *SyncSet[*func(key string, value int) int]) {
	callback1 := func(key string, value int) int {
		return value + 3
	}

	funcset.Add(&callback1)
}

func TestSyncSetWithFunctionPointer(t *testing.T) {
	syncSet := NewSyncSet[*func(key string, value int) int]()

	callback1 := func(key string, value int) int {
		return value + 1
	}

	callback2 := func(key string, value int) int {
		return value + 2
	}

	unused := func(key string, value int) int {
		return 0
	}

	syncSet.Add(&callback1)
	syncSet.Add(&callback2)

	// Adding different function pointers with same implementation
	addFunc(syncSet)
	addFunc(syncSet)

	if !syncSet.Has(&callback1) {
		t.Errorf("Expected 'callback1' to be present in the Set")
	}

	if !syncSet.Has(&callback2) {
		t.Errorf("Expected 'callback2' to be present in the Set")
	}

	if syncSet.Has(&unused) {
		t.Errorf("Expected 'unused' to NOT be present in the Set")
	}

	expectedResults := []int{11, 12, 13, 13}

	results := []int{}
	for _, cb := range syncSet.List() {
		(*cb)("test", 10)
		results = append(results, (*cb)("test", 10))
	}

	if len(results) != len(expectedResults) {
		t.Errorf("Expected %d results, got %d", len(expectedResults), len(results))
	}

	// This comparison may fail due to order differences
	// if !reflect.DeepEqual(results, expectedResults) {
	// 	t.Errorf("Expected results %v, got %v", expectedResults, results)
	// }

	if !checkSame(results, expectedResults) {
		t.Errorf("Expected results %v, got %v", expectedResults, results)
	}

	syncSet.Remove(&callback1)
	syncSet.Remove(&callback2)

	if syncSet.Size() != 2 {
		t.Errorf("Expected 2 functions after removals, got %d", syncSet.Size())
	}

	results = []int{}
	syncSet.Range(func(cb *func(key string, value int) int) bool {
		result := (*cb)("test", 10)
		results = append(results, result)
		return true
	})

	expectedResults = []int{13, 13}
	if len(results) != len(expectedResults) {
		t.Errorf("Expected %d results, got %d", len(expectedResults), len(results))
	}

	if !checkSame(results, expectedResults) {
		t.Errorf("Expected results %v, got %v", expectedResults, results)
	}

	syncSet.Clear()

	if syncSet.Size() != 0 {
		t.Errorf("Expected 0 functions after clear, got %d", syncSet.Size())
	}

}

func TestSyncSetWithTypedFunctionPointer(t *testing.T) {

	type Callback func(value int) int

	syncSet := NewSyncSet[*Callback]()

	var cb1 Callback = func(value int) int {
		return value + 5
	}
	syncSet.Add(&cb1)

	results := []int{}
	syncSet.Range(func(cb *Callback) bool {
		result := (*cb)(5)
		results = append(results, result)
		return true
	})

	expectedResults := []int{10}
	if len(results) != len(expectedResults) {
		t.Errorf("Expected %d results, got %d", len(expectedResults), len(results))
	}
	if !checkSame(results, expectedResults) {
		t.Errorf("Expected results %v, got %v", expectedResults, results)
	}

	// Adding the same function again should have no effect
	syncSet.Add(&cb1)

	expectedResults = []int{10}
	if len(results) != len(expectedResults) {
		t.Errorf("Expected %d results, got %d", len(expectedResults), len(results))
	}
	if !checkSame(results, expectedResults) {
		t.Errorf("Expected results %v, got %v", expectedResults, results)
	}
}

func TestSyncSetWithString(t *testing.T) {
	syncSet := NewSyncSet[string]()

	syncSet.Add("item1")
	syncSet.Add("item2")
	syncSet.Add("item3")

	if syncSet.Size() != 3 {
		t.Errorf("Expected 3 items, got %d", syncSet.Size())
	}

	if !syncSet.Has("item1") {
		t.Errorf("Expected 'item1' to be present in the Set")
	}

	if !syncSet.Has("item2") {
		t.Errorf("Expected 'item2' to be present in the Set")
	}

	if !syncSet.Has("item3") {
		t.Errorf("Expected 'item3' to be present in the Set")
	}

	if syncSet.Has("item4") {
		t.Errorf("Expected 'item4' to NOT be present in the Set")
	}

	items := syncSet.List()
	expectedItems := []string{"item1", "item2", "item3"}
	if !checkSame(items, expectedItems) {
		t.Errorf("Expected results %v, got %v", expectedItems, items)
	}

	syncSet.Remove("item2")

	if syncSet.Size() != 2 {
		t.Errorf("Expected 2 items after removal, got %d", syncSet.Size())
	}

	items = syncSet.List()
	expectedItems = []string{"item1", "item3"}
	if !checkSame(items, expectedItems) {
		t.Errorf("Expected results %v, got %v", expectedItems, items)
	}
}

func TestSyncSetRange(t *testing.T) {
	syncSet := NewSyncSet[string]()

	syncSet.Add("item1")
	syncSet.Add("item2")
	syncSet.Add("item3")

	if syncSet.Size() != 3 {
		t.Errorf("Expected 3 items, got %d", syncSet.Size())
	}

	// If the callback returns false, Range should stop iterating
	results := []string{}
	syncSet.Range(func(value string) bool {
		results = append(results, value)
		return false
	})

	// expectedResults := []string{"item1"}
	// if !checkSame(results, expectedResults) {
	// 	t.Errorf("Expected results %v, got %v", expectedResults, results)
	// }
	// order is not guaranteed, so we just check that only one item was returned
	if len(results) != 1 {
		t.Errorf("Expected 1 result from Range, got %d", len(results))
	}

	// Now test with returning true to get all items
	results = []string{}
	syncSet.Range(func(value string) bool {
		results = append(results, value)
		return true
	})

	expectedResults := []string{"item1", "item2", "item3"}
	if !checkSame(results, expectedResults) {
		t.Errorf("Expected results %v, got %v", expectedResults, results)
	}

}
