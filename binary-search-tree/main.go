package main

import "fmt"

var count int

// Node reperesents the components of a binary search tree
type Node struct {
	Key   int
	Left  *Node
	Right *Node
}

// insert will add a node to the tree
func (n *Node) Insert(k int) {
	if n.Key < k {
		//move right
		if n.Right == nil {
			//if the right node is empty great, place it
			n.Right = &Node{Key: k}
		} else {
			//if it's not empty then make recursive call which will do the check again until there is an empty slot
			n.Right.Insert(k)
		}
	} else if n.Key > k {
		//move left
		if n.Left == nil {
			n.Left = &Node{Key: k}
		} else {
			n.Left.Insert(k)
		}
	}
}

// search will take in a key value
// and return true if there is a node with that value
func (n *Node) Search(k int) bool {
	count++
	//if a match is found this conditional will be ignored, if neither n.Key < k and n.Key > k are true and the match isn't found it will execute
	if n == nil {
		return false
	}
	//we move the key in the appropriate direction
	if n.Key < k {
		return n.Right.Search(k)
	} else if n.Key > k {
		return n.Left.Search(k)
	}
	//if n is not nil nor n.Key < k and n.Key > k are all false it will return true, meaning a match is found and it doesn't need to look for it anymore nor return nil because it's found
	return true
}

func main() {
	tree := &Node{Key: 100}
	tree.Insert(52)
	tree.Insert(203)
	tree.Insert(19)
	tree.Insert(76)
	tree.Insert(150)
	tree.Insert(310)
	tree.Insert(7)
	tree.Insert(24)
	tree.Insert(88)
	tree.Insert(276)
	fmt.Println(tree.Search(310))
	fmt.Println(count)
}
