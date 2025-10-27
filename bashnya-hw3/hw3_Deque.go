package main

import "fmt"

type Node struct {
	up   *Node
	down *Node
	val  int
}

type Deque struct {
	highest *Node
	lowest  *Node
	size    int
}

func Size(d Deque) int {
	return d.size
}

func IsEmpty(d Deque) bool {
	return d.size == 0
}

func PushFront(d *Deque, v int) {
	if IsEmpty(*d) {
		n := new(Node)
		n.val = v
		n.down = nil
		n.up = nil
		d.highest = n
		d.lowest = n
		d.size = 1
	} else {
		n := new(Node)
		n.val = v
		n.down = d.highest
		n.up = nil
		d.highest.up = n
		d.highest = n
		d.size += 1
	}
}

func PushBack(d *Deque, v int) {
	if IsEmpty(*d) {
		n := new(Node)
		n.val = v
		n.down = nil
		n.up = nil
		d.highest = n
		d.lowest = n
		d.size = 1
	} else {
		n := new(Node)
		n.val = v
		n.up = d.lowest
		n.down = nil
		d.lowest.down = n
		d.lowest = n
		d.size += 1
	}
}

func PopFront(d *Deque) {
	d.size -= 1
	if d.size == 0 {
		d.highest = nil
		d.lowest = nil
	} else {
		d.highest = d.highest.down
	}
}

func PopBack(d *Deque) {
	d.size -= 1
	if d.size == 0 {
		d.lowest = nil
		d.highest = nil
	} else {
		d.lowest = d.lowest.up
	}
}

func Clear(d *Deque) {
	d.lowest = nil
	d.highest = nil
	d.size = 0
}

func main() {
	var n int
	fmt.Println("Введите n - число элементов дека:")
	fmt.Scan(&n)
	fmt.Printf("Введите %d чисел:", n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&nums[i])
	}
	s := Deque{highest: nil, lowest: nil, size: 0}
	for i := 0; i < n; i++ {
		PushBack(&s, nums[i])
	}
	fmt.Println("Элементы дека:")
	ptr := s.highest
	for i := 0; i < n; i++ {
		fmt.Println(ptr.val)
		ptr = ptr.down
	}
	fmt.Println("end")
}