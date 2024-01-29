package main

import (
	"fmt"
	"math/rand"
	"time"
)

func quickSort(a []int, l, r int) {
	if l >= r {
		return
	}
	p := partition(a, l, r)
	quickSort(a, l, p-1)
	quickSort(a, p+1, r)
}

func quickSortGo(a []int, l, r int, done chan struct{}) {
	if l >= r {
		done <- struct{}{}
		return
	}
	p := partition(a, l, r)
	childDone := make(chan struct{}, 2)
	go quickSortGo(a, l, p-1, childDone)
	go quickSortGo(a, p+1, r, childDone)
	<-childDone
	<-childDone
	done <- struct{}{}
}

func quickSortGoWithDepth(a []int, l, r int, done chan struct{}, depth int) {
	if l >= r {
		done <- struct{}{}
		return
	}
	depth--
	p := partition(a, l, r)
	if depth > 0 {
		childDone := make(chan struct{}, 2)
		go quickSortGoWithDepth(a, l, p-1, childDone, depth)
		go quickSortGoWithDepth(a, p+1, r, childDone, depth)
		<-childDone
		<-childDone
	} else {
		quickSort(a, l, p-1)
		quickSort(a, p+1, r)
	}
	done <- struct{}{}

}

// 将 a 分成左右两部分 返回中间下标
func partition(a []int, l, r int) int {
	p := a[r] // 最后一个数作为分界值
	i := l - 1
	for j := l; j < r; j++ { // 最后一个数已经作为分界值
		if a[j] < p {
			i++
			a[i], a[j] = a[j], a[i]
		}
	}
	a[i+1], a[r] = a[r], a[i+1]
	return i + 1
}

func benchQuickSort() {
	//	随机生成测试数据
	rand.Seed(time.Now().UnixNano())
	n := 10000000
	testData1, testData2, testData3 := make([]int, 0, n), make([]int, 0, n), make([]int, 0, n)
	for i := 0; i < n; i++ {
		val := rand.Intn(n * 100)
		testData1 = append(testData1, val)
		testData2 = append(testData2, val)
		testData3 = append(testData3, val)
	}
	//  串行执行
	start := time.Now()
	quickSort(testData1, 0, len(testData1)-1)
	fmt.Println("串行执行：", time.Since(start))
	//  并发执行
	done := make(chan struct{})
	start = time.Now()
	go quickSortGo(testData2, 0, len(testData2)-1, done)
	<-done
	fmt.Println("并发执行：", time.Since(start))
	//	并发优化
	done_ := make(chan struct{})
	start = time.Now()

	go quickSortGoWithDepth(testData3, 0, len(testData3)-1, done_, 16)
	<-done_
	fmt.Println("并发执行：", time.Since(start))
}

func main() {
	benchQuickSort()
	//串行执行： 751.744427ms
	//并发执行： 2.306186556s
	//并发执行： 168.963996ms
}
