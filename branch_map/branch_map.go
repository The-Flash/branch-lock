package branch_map

import (
	"crypto/sha1"
	"sync"
)

const MAX_BRANCHES = 65535

type branch struct {
	sync.RWMutex
	isLocked bool
}

type branchMap struct {
	mu       sync.Mutex
	branches []*branch
}

func NewBranchMap() branchMap {
	b := make([]*branch, MAX_BRANCHES, MAX_BRANCHES)
	return branchMap{
		branches: b,
	}
}

func (b *branchMap) getBranchIndex(name string) int {
	sum := sha1.Sum([]byte(name))
	hash := int(sum[13])<<8 | int(sum[17])
	return hash % len(b.branches)
}

func (b *branchMap) getBranch(name string) *branch {
	index := b.getBranchIndex(name)
	b.mu.Lock()
	defer b.mu.Unlock()
	if br := b.branches[index]; br == nil {
		b.add(name)
	}
	return b.branches[index]
}

func (b *branchMap) add(name string) {
	index := b.getBranchIndex(name)
	b.branches[index] = &branch{}
}

func (b *branchMap) Get(name string) *branch {
	branch := b.getBranch(name)
	return branch
}

// Lock lock a branch for push
func (b *branchMap) Lock(name string) {
	b.mu.Lock()
	br := b.getBranch(name)
	defer br.Unlock()
	br.isLocked = true
}

// UnLock unlock a branch for push
func (b *branchMap) UnLock(name string) {
	br := b.getBranch(name)
	br.Lock()
	defer br.Unlock()
	br.isLocked = false
}

// IsLocked a push is currently happening on a branch or not
func (b *branchMap) IsLocked(name string) bool {
	br := b.getBranch(name)
	br.RLock()
	defer br.RUnlock()
	return br.isLocked
}
