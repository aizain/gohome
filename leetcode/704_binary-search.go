package leetcode

// search
// 704. 二分查找 简单
//
// 给定一个 n 个元素有序的（升序）整型数组 nums 和一个目标值 target  ，写一个函数搜索 nums 中的 target，如果目标值存在返回下标，否则返回 -1。
//
// 提示：
// 1、你可以假设 nums 中的所有元素是不重复的。
// 2、n 将在 [1, 10000]之间。
// 3、nums 的每个元素都将在 [-9999, 9999]之间。
// 来源：力扣（LeetCode）
// 链接：https://leetcode.cn/problems/binary-search
// 著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
func search(nums []int, target int) int {
	return directlySearch(nums, target)
}

// directlySearch 直接搜索
func directlySearch(nums []int, target int) int {
	if nums == nil {
		return -1
	}
	for i, num := range nums {
		if num == target {
			return i
		}
	}
	return -1
}

// binarySearch 二分搜索
func binarySearch(nums []int, target int) int {
	if nums == nil {
		return -1
	}
	if len(nums) < 2 {
		return directlySearch(nums, target)
	}
	lPtr := 0
	rPtr := len(nums) - 1
	for rPtr-lPtr > 1 {
		cur := (rPtr + lPtr) / 2
		if v := nums[cur]; target == v {
			return cur
		} else if target < v {
			rPtr = cur
		} else {
			lPtr = cur
		}
	}

	if target == nums[rPtr] {
		return rPtr
	}
	if target == nums[lPtr] {
		return lPtr
	}
	return -1
}

// search20221115 搜索练习
//func search20221115(nums []int, target int) int {
//	i, j := 0, len(nums)-1
//	for i <= j {
//		mid := (i + j) >> 1
//		if nums[mid] == target {
//			return mid
//		}
//		if target < nums[mid] {
//			j--
//		} else {
//			i++
//		}
//	}
//
//	return -1
//}

// search20221115 搜索练习
func search20221115(nums []int, target int) int {
	i, j := 0, len(nums)-1
	for i <= j {
		mid := (i + j) >> 1
		if nums[mid] == target {
			return mid
		}
		if target < nums[mid] {
			j = mid - 1
		} else {
			i = mid + 1
		}
	}

	return -1
}
