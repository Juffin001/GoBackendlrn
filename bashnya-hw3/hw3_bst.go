// КОД НЕ ДОПИСАН. Еще ничего не работает
package main

import "fmt"

type Node struct {
	val    int
	left   *Node
	right  *Node
	parent *Node
}

type BST struct {
	top       *Node
	cnt_nodes int
	depth     int
}

func Insert(b *BST, v int, ptr *Node, cur_depth int) {
	if b.cnt_nodes == 0 {
		n := new(Node)
		n.val = v
		b.top = n
		b.cnt_nodes = 1
		b.depth = 1
		return
	}
	if ptr == nil {
		ptr = b.top
	}
	if v < ptr.val && ptr.left == nil {
		n := new(Node)
		n.val = v
		n.parent = ptr
		ptr.left = n
		b.cnt_nodes++
		b.depth = max(b.depth, cur_depth)
		return
	}
	if v > ptr.val && ptr.right == nil {
		n := new(Node)
		n.val = v
		n.parent = ptr
		ptr.right = n
		b.cnt_nodes++
		b.depth = max(b.depth, cur_depth)
		return
	}
	if v < b.top.val {
		Insert(b, v, ptr.left, cur_depth+1)
	} else {
		Insert(b, v, ptr.right, cur_depth+1)
	}
}

func Remove(b *BST, v int, ptr *Node) {
	if ptr == nil {
		ptr = b.top
	}
	if b.top.val == v && b.cnt_nodes == 1 {
		b.top = nil
		return
	}
	if ptr.val == v && ptr.left == nil && ptr.right == nil {
		if ptr.parent.left.val == v {
			ptr.parent.left = nil
		} else {
			ptr.parent.right = nil
		}
		return
	}
	if ptr.val == v && ptr.right == nil {
		if ptr.parent.left.val == v {
			ptr.parent.left = ptr.left
		} else {
			ptr.parent.right = ptr.left
		}
		return
	}
	if ptr.val == v && ptr.left == nil {
		if ptr.parent.left.val == v {
			ptr.parent.left = ptr.right
		} else {
			ptr.parent.right = ptr.right
		}
		return
	}
	if ptr.val == v { //всегда заменяем левым

	}
	if ptr.val > v {
		Remove(b, v, ptr.left)
	} else if ptr.val < v {
		Remove(b, v, ptr.right)
	}
}

func main() {
	var n int
	fmt.Println("Введите n - число элементов дерева:")
	fmt.Scan(&n)
	fmt.Printf("Введите %d чисел:\n", n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&nums[i])
	}
	s := BST{top: nil, cnt_nodes: 0, depth: 0}
	for i := 0; i < n; i++ {
		Insert(&s, nums[i], nil, 1)
	}
	fmt.Println("Элементы дерева:")
	ptr := s.top
	for i := 0; i < 4; i++ {
		fmt.Println(ptr.val)
		ptr = ptr.left
	}
	fmt.Println("end")
}
