// Весь код дерева написан мною (около 4 часов суммарно писал 0о0)
// Функцию print_tree и main написал дипсик, когда я попросил его наглядный пример для проверки работы дерева
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
}

func Insert(b *BST, v int, ptr *Node) {
	if b.cnt_nodes == 0 {
		n := new(Node)
		n.val = v
		b.top = n
		b.cnt_nodes = 1
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
		return
	}
	if v > ptr.val && ptr.right == nil {
		n := new(Node)
		n.val = v
		n.parent = ptr
		ptr.right = n
		b.cnt_nodes++
		return
	}
	if v < ptr.val {
		Insert(b, v, ptr.left)
	} else {
		Insert(b, v, ptr.right)
	}
}

func Find_min_in_right(ptr *Node) *Node { //указатель на правого чела
	if ptr.left == nil {
		return ptr
	}
	return Find_min_in_right(ptr.left)
}

func Remove(b *BST, v int, ptr *Node) {
	if ptr == nil {
		ptr = b.top
	}
	if ptr == nil {
		return
	}
	if ptr.val > v {
		if ptr.left != nil {
			Remove(b, v, ptr.left)
		}
		return
	}
	if ptr.val < v {
		if ptr.right != nil {
			Remove(b, v, ptr.right)
		}
		return
	}
	b.cnt_nodes--
	if ptr.left == nil && ptr.right == nil { //Лист
		if ptr.parent == nil {
			b.top = nil
		} else if ptr.parent.left == ptr {
			ptr.parent.left = nil
		} else if ptr.parent.right == ptr {
			ptr.parent.right = nil
		}
		return
	}
	// 1 ребенок
	if ptr.left == nil && ptr.right != nil {
		if ptr.parent == nil {
			b.top = ptr.right
		} else if ptr.parent.right == ptr {
			ptr.parent.right = ptr.right
		} else {
			ptr.parent.left = ptr.right
		}
		if ptr.right != nil {
			ptr.right.parent = ptr.parent
		}
		return
	}

	if ptr.right == nil && ptr.left != nil {
		if ptr.parent == nil {
			b.top = ptr.left
		} else if ptr.parent.right == ptr {
			ptr.parent.right = ptr.left
		} else {
			ptr.parent.left = ptr.left
		}
		if ptr.left != nil {
			ptr.left.parent = ptr.parent
		}
		return
	}

	//2 ребенка (0o0)
	min_in_right := Find_min_in_right(ptr.right)
	ptr.val = min_in_right.val
	b.cnt_nodes++ // рекурсия
	Remove(b, min_in_right.val, ptr.right)
}

func Find(v int, ptr *Node) *Node { //ссылку если нашел, nil - иначе
	if ptr == nil {
		return nil
	}
	if v == ptr.val {
		return ptr
	} else if v > ptr.val {
		return Find(v, ptr.right)
	} else {
		return Find(v, ptr.left)
	}
}

func Depth(b *BST, now int, ans int, ptr *Node) int { //b=b now=1 ans=0
	if ptr.left != nil {
		ans = max(ans, Depth(b, now+1, ans, ptr.left))
	}
	if ptr.right != nil {
		ans = max(ans, Depth(b, now+1, ans, ptr.right))
	}
	ans = max(ans, now)
	return ans
}

//Далее код дипсика

func printTree(ptr *Node, prefix string, isLeft bool) {
	if ptr == nil {
		return
	}

	fmt.Print(prefix)
	if isLeft {
		fmt.Print("├── ")
	} else {
		fmt.Print("└── ")
	}
	fmt.Println(ptr.val)

	if ptr.left != nil {
		newPrefix := prefix
		if isLeft {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		printTree(ptr.left, newPrefix, true)
	}

	if ptr.right != nil {
		newPrefix := prefix
		if isLeft {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		printTree(ptr.right, newPrefix, false)
	}
}

func main() {
	bst := &BST{}

	fmt.Println("=== Создаем дерево ===")
	values := []int{8, 3, 10, 1, 6, 14, 4, 7, 13}
	for _, v := range values {
		fmt.Printf("Вставляем %d:\n", v)
		Insert(bst, v, nil)
		printTree(bst.top, "", false)
		fmt.Printf("Узлов: %d, Глубина: %d\n\n", bst.cnt_nodes, Depth(bst, 1, 0, bst.top))
	}

	fmt.Println("=== Поиск значений ===")
	testValues := []int{6, 5, 14, 20}
	for _, v := range testValues {
		if found := Find(v, bst.top); found != nil {
			fmt.Printf("Значение %d найдено\n", v)
		} else {
			fmt.Printf("Значение %d НЕ найдено\n", v)
		}
	}
	fmt.Println()

	fmt.Println("=== Удаление листа (7) ===")
	Remove(bst, 7, nil)
	printTree(bst.top, "", false)
	fmt.Printf("Узлов: %d, Глубина: %d\n\n", bst.cnt_nodes, Depth(bst, 1, 0, bst.top))

	fmt.Println("=== Удаление узла с одним потомком (10) ===")
	Remove(bst, 10, nil)
	printTree(bst.top, "", false)
	fmt.Printf("Узлов: %d, Глубина: %d\n\n", bst.cnt_nodes, Depth(bst, 1, 0, bst.top))

	fmt.Println("=== Удаление узла с двумя потомками (3) ===")
	Remove(bst, 3, nil)
	printTree(bst.top, "", false)
	fmt.Printf("Узлов: %d, Глубина: %d\n\n", bst.cnt_nodes, Depth(bst, 1, 0, bst.top))

	fmt.Println("=== Удаление корня (8) ===")
	Remove(bst, 8, nil)
	printTree(bst.top, "", false)
	fmt.Printf("Узлов: %d, Глубина: %d\n", bst.cnt_nodes, Depth(bst, 1, 0, bst.top))
}
