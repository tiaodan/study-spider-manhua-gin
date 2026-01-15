/*
功能: 通用工具类
*/
package util

import "study-spider-manhua-gin/src/log"

/* 判断数组内容是否有重复项

实现: 通过 map实现
优点：

	时间复杂度：O(n)
	空间复杂度：O(n)（最坏情况）
	代码极其简洁，可读性强
	在实际场景中性能最好（尤其是数据量大时）
*/
func HasDuplicate(nums []int) bool {
	if len(nums) == 0 {
		log.Warn("判断重复, 数组为空, 不遍历了")
		return false
	}

	// v0.3写法: 找到所有重复项再return
	seen := make(map[int]bool)
	log.Debug("开始检测重复...")

	var repeateArr []int
	for _, num := range nums {
		// log.Debugf("第 %d 次循环，当前数字: %d\n", i+1, num) // 调试用，注释，占日志

		if seen[num] { // 如果这个数字之前已经见过
			// log.Debugf("!!! 发现重复: %d 已经出现过了！\n", num) // 调试用，注释，占日志
			repeateArr = append(repeateArr, num)
			continue
		}

		// 第一次见到这个数字，标记为已见过
		seen[num] = true
		// log.Debugf("   第一次见, 标记 %d 为已见过\n", num) // 调试用，注释，占日志
		// log.Debug("   当前 seen map:", seen) // 调试用，注释，占日志
	}

	if len(repeateArr) > 0 {
		log.Warn("发现所有重复项: ", repeateArr)
		return true
	}

	log.Debugf("遍历结束，没有发现重复")
	return false

	/* v0.2写法 遇到1个重复就返回
	log.Debug("传参 = ", nums)
	seen := make(map[int]bool)
	log.Debug("开始检测重复...")

	for _, num := range nums { // 这样写 num才是实际值，如果少_ 那num就是index，不对
		// log.Debugf("第 %d 次循环，当前数字: %d\n", i+1, num) // 调试用，注释，占日志

		if seen[num] { // 如果这个数字之前已经见过
			log.Debugf("!!! 发现重复: %d 已经出现过了！\n", num)
			return true
		}

		// 第一次见到这个数字，标记为已见过
		seen[num] = true
		// log.Debugf("   第一次见, 标记 %d 为已见过\n", num) // 调试用，注释，占日志
		// log.Debug("   当前 seen map:", seen) // 调试用，注释，占日志
	}

	log.Debugf("遍历结束，没有发现重复")
	return false
	*/

	/* v0.1 简洁写法
	   seen := make(map[int]bool)
	   for _, num := range nums {
	       if seen[num] {
	           return true
	       }
	       seen[num] = true
	   }
	   return false
	*/
}

// 拆分数组 (int类型),分成N份
/*
参数：arr []int,
	max int 拆分的数组，最大个数
*/
func SplitIntArr(arr []int, splitArrMax int) [][]int {
	var result [][]int

	for i := 0; i < len(arr); i += splitArrMax {
		end := i + splitArrMax
		if end > len(arr) { // 结束for 标志
			end = len(arr)
		}
		result = append(result, arr[i:end])
	}

	return result
}
