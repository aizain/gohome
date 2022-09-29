package leetcode

import "time"

/**
firstBadVersion
278. 第一个错误的版本 简单

你是产品经理，目前正在带领一个团队开发新的产品。不幸的是，你的产品的最新版本没有通过质量检测。
由于每个版本都是基于之前的版本开发的，所以错误的版本之后的所有版本都是错的。
假设你有 n 个版本 [1, 2, ..., n]，你想找出导致之后所有版本出错的第一个错误的版本。
你可以通过调用 bool isBadVersion(version) 接口来判断版本号 version 是否在单元测试中出错。
实现一个函数来查找第一个错误的版本。你应该尽量减少对调用 API 的次数。

提示：
1 <= bad <= n <= 231 - 1

来源：力扣（LeetCode）
链接：https://leetcode.cn/problems/first-bad-version
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
*/
/**
 * Forward declaration of isBadVersion API.
 * @param   version   your guess about first bad version
 * @return 	 	      true if current version is bad
 *			          false if current version is good
 * func isBadVersion(version int) bool;
 */
func firstBadVersion(n int) int {
	return directlyFindVersion(n)
}
func incrFindVersion(n int) int {
	last := 0
	llast := 0
	for {
		ok := isBadVersion(n)
		if ok {
			if last == 1 {
				return n
			}
			if last > 0 {
				last = 0
				llast = 0
			} else if last <= -1 {
				tmp := last
				last += llast
				llast = tmp
			} else {
				last = -1
				llast = 0
			}

		} else {
			if last == -1 {
				return n + 1
			}
			if last < 0 {
				last = 0
				llast = 0
			} else if last >= 1 {
				tmp := last
				last += llast
				llast = tmp
			} else {
				last = 1
				llast = 0
			}
		}
		n += last
	}
}

func directlyFindVersion(n int) int {
	// 0 none 1 good -1 bad
	last := 0
	for {
		ok := isBadVersion(n)
		if ok {
			if last == 1 {
				return n
			}
			last = -1
		} else {
			if last == -1 {
				return n + 1
			}
			last = 1
		}
		n += last
	}

}

func isBadVersion(n int) bool {
	time.Sleep(100 * time.Microsecond)
	if n >= 4 {
		return true
	}
	return false
}
