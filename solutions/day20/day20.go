package main

import (
	"embed"
	"fmt"

	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
	"golang.org/x/exp/slices"
)

var log = logging.GetLogger()

//go:embed input.txt
var files embed.FS

type Node struct {
	Prev  *Node
	Next  *Node
	Value int
}

func getPartOneData() []*Node {
	var nodes = make([]*Node, 0)
	for _, line := range util.ReadProblemInput(files) {
		nodes = append(nodes, &Node{Value: util.AtoiNoError(line)})
	}

	for i := 1; i < len(nodes); i++ {
		nodes[i].Prev = nodes[i-1]
		nodes[i-1].Next = nodes[i]
	}
	nodes[0].Prev = nodes[len(nodes)-1]
	nodes[len(nodes)-1].Next = nodes[0]

	return nodes
}

func Swap(a, b *Node) {
	a.Prev.Next = b
	b.Next.Prev = a
	a.Next = b.Next
	b.Prev = a.Prev
	a.Prev = b
	b.Next = a
}

func RunDecryption(decryptKey, numMixes int) []int {
	data := getPartOneData()

	// Apply decryption key
	for _, n := range data {
		n.Value *= decryptKey
	}

	for i := 0; i < numMixes; i++ {
		for _, n := range data {
			shift := n.Value % (len(data) - 1) // modulo to ignore complete loops around the list
			for shift != 0 {
				if shift > 0 {
					Swap(n, n.Next)
					shift--
				} else {
					Swap(n.Prev, n)
					shift++
				}

			}
		}
	}

	idx := slices.IndexFunc(data, func(n *Node) bool { return n.Value == 0 })
	start := data[idx]
	encrypted := []int{start.Value}

	next := start.Next
	for next != start {
		encrypted = append(encrypted, next.Value)
		next = next.Next
	}

	return encrypted
}

func PartOne() any {
	encrypted := RunDecryption(1, 1)
	return encrypted[1000%len(encrypted)] + encrypted[2000%len(encrypted)] + encrypted[3000%len(encrypted)]
}

func PartTwo() any {
	encrypted := RunDecryption(811_589_153, 10)
	return encrypted[1000%len(encrypted)] + encrypted[2000%len(encrypted)] + encrypted[3000%len(encrypted)]
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
