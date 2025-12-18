package syncset

import (
	"sync"
)

// A Set implemented by a sync.Map.
type SyncSet[T comparable] struct {
	items sync.Map
}

// Create a new SyncSet instance.
// Specifying the type of items. They mmust be comparable. It can be pointers.
func NewSyncSet[T comparable]() *SyncSet[T] {
	return &SyncSet[T]{
		items: sync.Map{},
	}
}

// Add an item to the Set.
func (this *SyncSet[T]) Add(item T) {
	this.items.Store(item, struct{}{})
}

// Remove item from Set.
func (this *SyncSet[T]) Remove(item T) {
	this.items.Delete(item)
}

// Returns a slice of all items in the SyncSet.
// Warning: The order of items in the returned slice is not guaranteed.
func (this *SyncSet[T]) List() []T {
	list := []T{}
	this.items.Range(func(k, v any) bool {
		if item, ok := k.(T); ok {
			list = append(list, item)
		}
		return true
	})
	return list
}

// Range calls f sequentially for each item present in the Set.
// If f returns false, range stops the iteration.
// Uses underlying sync.Map's Range method.
func (this *SyncSet[T]) Range(f func(item T) bool) {
	this.items.Range(func(k, v any) bool {
		if item, ok := k.(T); ok {
			return f(item)
		}
		return true
	})
}

// Returns the number of items in the Set.
func (this *SyncSet[T]) Size() int {
	count := 0
	this.items.Range(func(k, v any) bool {
		count++
		return true
	})
	return count
}

// Clears every items from the Set.
func (this *SyncSet[T]) Clear() {
	this.items.Clear()
}
