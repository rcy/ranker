package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Node struct {
	label    string
	parents  []*Node
	children []*Node
}

func (n Node) String() string {
	return fmt.Sprintf("%s", n.label)
}

var nodeList []*Node

func createNode(label string, parent *Node) *Node {
	node := Node{label: label}
	nodeList = append(nodeList, &node)

	if parent != nil {
		prefer(parent, &node)
	}

	return &node
}

func prefer(parent, child *Node) {
	parent.children = append(parent.children, child)

	// walkthrough entire node list, and find nodes that include both parent and child as children and remove child
	for _, n := range nodeList {
		preferSibling(n, parent, child)
	}
}

func preferSibling(node, keep, drop *Node) {
	var hasKeep bool
	var hasDrop bool

	// if node.children includes both keep and drop, update it to remove drop
	for _, n := range node.children {
		if n == keep {
			hasKeep = true
		}
		if n == drop {
			hasDrop = true
		}
	}
	if hasKeep && hasDrop {
		children := []*Node{keep}
		for _, n := range node.children {
			if n != drop && n != keep {
				children = append(children, n)
			}
		}
		node.children = children
	}
}

func printNodes(root *Node) {
	fmt.Println("---")
	for _, n := range nodeList {
		// if n == root {
		// 	continue
		// }
		fmt.Printf("%s", n)
		if len(n.children) > 0 {
			fmt.Print(" > ")
			for _, c := range n.children {
				fmt.Printf("%s ", c)
			}
		}
		fmt.Println()
	}
	fmt.Println("---")
}

type Matchup struct {
	A *Node
	B *Node
}

func findMatchups() []Matchup {
	// go through each node, if there are multiple children, take the last 2
	matchups := []Matchup{}
	for _, n := range nodeList {
		cs := n.children
		if len(cs) >= 2 {
			matchups = append(matchups, Matchup{cs[len(cs)-1], cs[len(cs)-2]})
		}
	}
	return matchups
}

func faceoff(root *Node, matchup Matchup) (winner, loser *Node) {
	fmt.Printf("a: %s\nb: %s\n", matchup.A, matchup.B)

	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		log.Fatal(err)
	}

	switch input {
	case "a", "A":
		return matchup.A, matchup.B
	case "b", "B":
		return matchup.B, matchup.A
	case "?":
		printNodes(root)
	}
	return faceoff(root, matchup)
}

func readOptions(root *Node) {
	fmt.Println("Enter one option per line. Enter a blank line when done.")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			log.Fatal(err)
		}
		input := scanner.Text()
		if input == "" {
			break
		}
		createNode(input, root)
	}
}

func runTournament(root *Node) {
	fmt.Println("Enter a or b to indicate your preference for the following items:")
	var matchups []Matchup
	for {
		matchups = findMatchups()
		if len(matchups) == 0 {
			break
		}
		winner, loser := faceoff(root, matchups[0])
		prefer(winner, loser)
	}
}

func showResults(root *Node) {
	fmt.Println("Here are the results:")
	node := root
	order := 1
	for {
		if len(node.children) == 0 {
			break
		}
		fmt.Printf("%d\t%s\n", order, node.children[0])
		order += 1
		node = node.children[0]
	}
}

func main() {
	root := createNode("0", nil)

	readOptions(root)
	runTournament(root)
	showResults(root)
}
