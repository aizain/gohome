package leetcode

/*
35. 搜索插入位置 简单

给定一个排序数组和一个目标值，在数组中找到目标值，并返回其索引。如果目标值不存在于数组中，返回它将会被按顺序插入的位置。
请必须使用时间复杂度为 O(log n) 的算法。

示例 1:
输入: nums = [1,3,5,6], target = 5
输出: 2

示例 2:
输入: nums = [1,3,5,6], target = 2
输出: 1

示例 3:
输入: nums = [1,3,5,6], target = 7
输出: 4

提示:

1 <= nums.length <= 104
-104 <= nums[i] <= 104
nums 为 无重复元素 的 升序 排列数组
-104 <= target <= 104

来源：力扣（LeetCode）
链接：https://leetcode.cn/problems/search-insert-position
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
*/
func searchInsert(nums []int, target int) int {
	if nums == nil {
		return 0
	}
	ln := len(nums)
	if ln <= 0 || target <= nums[0] {
		return 0
	}
	if target > nums[ln-1] {
		return ln
	}
	l, r := 0, ln
	for l <= r {
		mid := (l + r) >> 1
		if lv := nums[mid]; target > lv && target <= nums[mid+1] {
			return mid + 1
		} else if target > lv {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}

	return l + 1
}
