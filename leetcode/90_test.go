package leetcode

import (
	"fmt"
	"slices"
	"testing"
)

func Test90(t *testing.T) {
	nums := []int{1, 2, 2}
	slices.Sort(nums)
	res := _subsetsWithDup(nums)
	res = append(res, []int{})
	fmt.Println(res)
}
func _subsetsWithDup(nums []int) [][]int {
	var res [][]int
	var leftRes [][]int
	var subRes [][]int
	for i := 0; i < len(nums); i++ {
		leftRes = append(leftRes, nums[:i+1])
		if i+1 < len(nums) && nums[i+1] != nums[i] {
			subRes = _subsetsWithDup(nums[i+1:])
			break
		}
	}
	res = subRes // should only copy once
	for _, item := range leftRes {
		res = append(res, item)
		for _, subItem := range subRes {
			newItem := append([]int{}, item...)
			res = append(res, append(newItem, subItem...))
		}
	}
	return res
}
