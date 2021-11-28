package main

import (
	"fmt"
	"os"
	"sort"
	"time"
)

func powmod(a, d, n int) int {
	r := a
	for d > 1 {
		r *= a
		r %= n
		d--
	}
	return r
}

func single_test(n, a int) bool {
	s := 0
	d := n - 1
	for d&1 == 0 {
		d /= 2
		s++
	}

	if powmod(a, d, n) == 1 {
		return true
	}

	for r := 0; r < s; r++ {
		if powmod(a, d, n) == n-1 {
			return true
		}
	}

	return false
}

func miller_rabin(n int, as ...int) bool {
	for _, a := range as {
		if !single_test(n, a) {
			return false
		}
	}
	return true
}

type tuple struct {
	a int
	b int
}

func newTuple(a, b int) tuple {
	return tuple{a, b}
}

type sortRight []tuple

func (a sortRight) Len() int           { return len(a) }
func (a sortRight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortRight) Less(i, j int) bool { return a[i].b < a[j].b }

func main() {
	teams := []struct {
		max int
		as  []int
	}{
		{2047, []int{2}},
		{1373653, []int{2, 3}},
		{9080191, []int{31, 73}},
		{25326001, []int{2, 3, 5}},
		{3215031751, []int{2, 3, 5, 7}},
		{4759123141, []int{2, 7, 61}},
		{1122004669633, []int{2, 13, 23, 1662803}},
		{2152302898747, []int{2, 3, 5, 7, 11}},
		{3474749660383, []int{2, 3, 5, 7, 11, 13}},
		{341550071728321, []int{2, 3, 5, 7, 11, 13, 17}},
	}

	ticker := time.NewTicker(time.Millisecond * 100)

	lies := make(map[int]int)
	last := int(2048)

	n := int(3)
	for i := 0; i < len(teams) && n <= last; i++ {
		team := teams[i]
		for n < team.max && n <= last {
			if !miller_rabin(n, team.as...) {
				for a := int(2); a < n; a++ {
					if single_test(n, a) {
						lies[a]++
					}
				}
			}

			select {
			case <-ticker.C:
				fmt.Fprintf(os.Stderr, "\r%.3f%%", float64(n)/float64(last)*100)
			default:
			}

			n += 2
		}
		fmt.Fprint(os.Stderr, "\r100.000%")
	}
	fmt.Fprintln(os.Stderr)

	tt := make([]tuple, 0, len(lies))
	for a, b := range lies {
		tt = append(tt, tuple{a, b})
	}
	sort.Sort(sort.Reverse(sortRight(tt)))
	for _, t := range tt {
		fmt.Println(t.a, t.b)
	}
}
