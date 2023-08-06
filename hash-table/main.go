package main

import "fmt"

const ArraySize = 7

// HashTable structure
type HashTable struct {
	array [ArraySize]*bucket
}

// bucket structure (will be linked list, in each slot/index of the HashTable)
type bucket struct {
	head *bucketNode
}

// bucketNode structure (node is each key/value)
type bucketNode struct {
	key  string
	next *bucketNode
}

// for HashTable
func hash(key string) int {
	//get ascii code for each character, sum it up and divide it by the array size and get the remainder
	sum := 0 //initialize sum
	//we need a for loop to loop through each character of the key
	for _, v := range key {
		sum += int(v) //so each letter is getting changed to an integer and getting added up
	}
	//we get the index, or the numerical value that corresponds to a poisition in the hash table
	return sum % ArraySize
}

// // insert will take in a key and add it to the hash table array
func (h *HashTable) Insert(key string) {
	index := hash(key)
	h.array[index].insert(key)
}

// // search will take in a key and return true if that key is stored in the hash table
func (h *HashTable) Search(key string) bool {
	index := hash(key)
	return h.array[index].search(key)
}

// // delete will take in a key and delete it from the hash table
func (h *HashTable) Delete(key string) {
	index := hash(key)
	h.array[index].delete(key)
}

// for bucket
// insert will take in a key, create a node with the key and insert the node in the bucket
func (b *bucket) insert(k string) {
	if !b.search(k) {
		newNode := &bucketNode{key: k} //initializes and sets the newNode to the address of a bucketNode with they key: k
		newNode.next = b.head          //sets the next attribute to the head node, the new node becomes the first node
		b.head = newNode               //makes the new node the first element in the linked list, or new head of the bucket
	} else {
		fmt.Println("The bucket node already exists")
	}
}

// search will take in a key and return true if the key is found
func (b *bucket) search(k string) bool {
	currentNode := b.head
	//going to keep on looping until we find a match, until the current node is empty
	for currentNode != nil {
		if currentNode.key == k {
			return true
		}
		currentNode = currentNode.next
	}
	return false

}

// delete
func (b *bucket) delete(k string) {
	//we don't want to miss if the matching key is the head
	if b.head.key == k { //if the head node is they key we want to delete we reset the head of this bucket to the second node
		b.head = b.head.next
		return
	}

	previousNode := b.head
	//the previousNode.next is the current node
	for previousNode.next != nil {
		if previousNode.next.key == k {
			//delete
			previousNode.next = previousNode.next.next
		}
		previousNode = previousNode.next
	}
}

//define hash function

//define init function that initializes the hash table

func main() {
	hashTable := Init() //creates a hashtable that has a bucket at each index, from the Init function we defined
	list := []string{
		"ERIC",
		"KENNY",
		"KYLE",
		"STAN",
		"RANDY",
		"BUTTERS",
		"TOKEN",
	}

	for _, v := range list {
		hashTable.Insert(v)
	}
	//hashTable.Delete("STAN")
	fmt.Println(hashTable.Search("STAN"))
	fmt.Println(hashTable.Search("KENNY"))
	// fmt.Println(testHashTable)
	// fmt.Println(hash("RANDY"))

	// testBucket := &bucket{}
	// testBucket.insert("RANDY")
	// testBucket.delete("RANDY")

	// fmt.Println(testBucket.search("RANDY"))
	// fmt.Println(testBucket.search("ERIC"))
}

// emils
// init will create a bucket in each slot of the hash table
func Init() *HashTable {
	result := &HashTable{}
	//for loop that goes through each i of the hash table
	for i := range result.array {
		//creates a bucket at each index
		result.array[i] = &bucket{}
	}
	return result //returns the hashtable with the buckets at each slot
}
