package branch_map

import (
	"fmt"
	"testing"
)

func TestBranchMapLock(t *testing.T) {
	branchMap := NewBranchMap()
	name := "feat/new-branch"
	for range 10 {
		go func() {
			branchMap.Add(name)
			branchMap.Lock(name)
		}()
	}
	for i := range 10 {
		go func(i int) {
			name := "feat/new" + fmt.Sprint(i)
			branchMap.Add(name)
			branchMap.Lock("feat/new" + fmt.Sprint(i))
		}(i)
	}
}

func TestBranchMapUnLock(t *testing.T) {
	branchMap := NewBranchMap()
	name := "feat/new-branch"
	for range 10 {
		go func() {
			branchMap.Add(name)
			branchMap.UnLock(name)
		}()
	}
	for i := range 10 {
		go func(i int) {
			name := "feat/new" + fmt.Sprint(i)
			branchMap.Add(name)
			branchMap.UnLock("feat/new" + fmt.Sprint(i))
		}(i)
	}
}
func TestBranchMapIsLocked(t *testing.T) {
	branchMap := NewBranchMap()
	name := "feat/new-branch"
	for range 10 {
		branchMap.Add(name)
		branchMap.Lock(name)
		branch := branchMap.Get(name)
		if branch.isLocked == false {
			t.Fatal("branch is expected to be locked")
		}
		branchMap.UnLock(name)
		if branch.isLocked == true {
			t.Fatal("branch is expected to be unlocked")
		}
	}
}
