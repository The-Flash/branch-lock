package branch_map

import (
	"fmt"
	"testing"
)

func TestBranchMapAdd(t *testing.T) {
	branchMap := NewBranchMap()
	for i := range 10 {
		go func(name string) {
			branchMap.Add(name + fmt.Sprint(i))
		}("feat/new-branch")
		go func(name string) {
			branchMap.Add(name)
		}("feat/new-branch")
	}
}
func TestBranchMapLock(t *testing.T) {
	branchMap := NewBranchMap()
	name := "feat/new-branch"
	branchMap.Add(name)
	for _ = range 10 {
		go func() {
			branchMap.Lock(name)
		}()
	}
}
func TestBranchMapUnLock(t *testing.T) {
	branchMap := NewBranchMap()
	name := "feat/new-branch"
	branchMap.Add(name)
	for _ = range 10 {
		go func() {
			branchMap.UnLock(name)
		}()
	}
}
func TestBranchMapIsLocked(t *testing.T) {

	branchMap := NewBranchMap()
	name := "feat/new-branch"
	branchMap.Add(name)
	for _ = range 10 {
		go func() {
			branchMap.Lock(name)
			branch := branchMap.Get(name)
			if branch.isLocked == false {
				t.Fatal("branch is expected to be locked")
			}
			branchMap.UnLock(name)
			if branch.isLocked == false {
				t.Fatal("branch is expected to be unlocked")
			}
		}()
	}
}
