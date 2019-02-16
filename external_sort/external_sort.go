package external_sort

import (
	"sort"
)

// 对两个有序序列（用channel来存储）进行归并
func Merge(in1, in2 <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		v1, ok1 := <-in1
		v2, ok2 := <-in2

		for ok1 || ok2 {
			if !ok2 || (ok1 && v1 <= v2) {
				out <- v1
				v1, ok1 = <-in1
			} else {
				out <- v2
				v2, ok2 = <-in2
			}
		}
		close(out)
	}()
	return out
}

// 若干个有序序列合并
func MergeN(inputs ...<-chan int) <-chan int {
	n := len(inputs)
	if n == 1 {
		return inputs[0]
	}

	mid := len(inputs) / 2
	return Merge(
		MergeN(inputs[:mid]...),
		MergeN(inputs[mid:]...))
}

func MemorySort(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		a := make([]int, 0)
		for v := range in {
			a = append(a, v)
		}
		sort.Ints(a)

		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return out
}

// 从任意数据源读取数据, 转成-chan int的接口
type Reader interface {
	Read() (<-chan int, error)
}

func ReaderSource(reader Reader) (<-chan int, error) {
	return reader.Read()
}
