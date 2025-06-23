package set

import "sync"

// Set 基于map实现的线程安全集合
type Set[T comparable] struct {
	mu sync.RWMutex
	m  map[T]struct{}
}

// New 创建新集合
func New[T comparable](capacity ...int) *Set[T] {
	s := &Set[T]{m: make(map[T]struct{})}
	if len(capacity) > 0 {
		s.m = make(map[T]struct{}, capacity[0])
	}
	return s
}

// Add 添加元素
func (s *Set[T]) Add(items ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range items {
		s.m[item] = struct{}{}
	}
}

// Remove 删除元素
func (s *Set[T]) Remove(items ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range items {
		delete(s.m, item)
	}
}

// Contains 检查元素是否存在
func (s *Set[T]) Contains(item T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.m[item]
	return exists
}

// Size 返回集合大小
func (s *Set[T]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.m)
}

// Clear 清空集合
func (s *Set[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m = make(map[T]struct{})
}

// ToSlice 转换为切片
func (s *Set[T]) ToSlice() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]T, 0, len(s.m))
	for k := range s.m {
		result = append(result, k)
	}
	return result
}
