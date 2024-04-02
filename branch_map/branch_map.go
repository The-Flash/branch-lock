package branch_map

import (
	"crypto/sha1"
	"sync"
)

const MAX_BRANCHES = 65535

type branch struct {
	sync.RWMutex
	name     string
	isLocked bool
}

type branchMap struct {
	mu       sync.Mutex
	branches [][]*branch
}

func NewBranchMap() branchMap {
	b := make([][]*branch, MAX_BRANCHES)
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
	br := b.branches[index]
	for _, branch := range br {
		if branch.name == name {
			return branch
		}
	}
	return nil
}

// Add add a branch to the map
func (b *branchMap) Add(name string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	index := b.getBranchIndex(name)
	mapArray := b.branches[index]
	if mapArray == nil {
		b.branches[index] = make([]*branch, 0)
	}
	arr := b.branches[index]
	for _, branch := range arr {
		if branch.name == name {
			return
		}
	}
	newBranch := &branch{name: name}
	arr = append(arr, newBranch)
	b.branches[index] = arr
}

func (b *branchMap) Get(name string) *branch {
	branch := b.getBranch(name)
	return branch
}

// Lock lock a branch for push
func (b *branchMap) Lock(name string) {
	br := b.getBranch(name)
	if br == nil {
		return
	}
	br.Lock()
	defer br.Unlock()
	br.isLocked = true
}

// UnLock unlock a branch for push
func (b *branchMap) UnLock(name string) {
	br := b.getBranch(name)
	if br == nil {
		return
	}
	br.Lock()
	defer br.Unlock()
	br.isLocked = false
}

// IsLocked a push is currently happening on a branch or not
func (b *branchMap) IsLocked(name string) bool {
	br := b.getBranch(name)
	if br == nil {
		return false
	}

	br.RLock()
	defer br.RUnlock()
	return br.isLocked
}
