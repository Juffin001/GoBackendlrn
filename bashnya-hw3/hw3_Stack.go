package main

import "fmt"

type Node struct {
	next *Node
	val  int
}

type Stack struct {
	last *Node
	size int
}

func IsEmpty(s Stack) bool {
	return s.last == nil
}

func Size(s Stack) int {
	return s.size
}

func push(s *Stack, v int) {
	if IsEmpty(*s) {
		s.size = 1
		s.last = &Node{val: v, next: nil}
	} else {
		var n = new(Node)
		n.val = v
		n.next = s.last
		s.last = n
		s.size += 1
	}
}

func pop(s *Stack) {
	s.size -= 1
	s.last = s.last.next
}

func clear(s *Stack) {
	s.last = nil
	s.size = 0
}

func main() {
	var n int
	fmt.Println("Введите n - число элементов стека:")
	fmt.Scan(&n)
	fmt.Printf("Введите %d чисел:", n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&nums[i])
	}
	s := Stack{last: nil, size: 0}
	for i := 0; i < n; i++ {
		push(&s, nums[i])
	}
	fmt.Println("Элементы стека:\n")
	ptr := s.last
	for i := 0; i < n; i++ {
		fmt.Println(ptr.val)
		ptr = ptr.next
	}
	fmt.Println("end")
}
