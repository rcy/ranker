package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var debug bool

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
	if debug {
		fmt.Printf("prefer %s over %s\n", parent, child)
	}
	parent.children = append(parent.children, child)

	// walkthrough entire node list, and find nodes that include both parent and child as children and remove child
	for _, n := range nodeList {
		preferSibling(n, parent, child)
	}

	printNodes()
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

func printNodes() {
	if !debug {
		return
	}
	for _, n := range nodeList {
		fmt.Printf("%s: ", n)
		for _, c := range n.children {
			fmt.Printf("%s ", c)
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
	// go through each node
	// if there are multiple children, take the first 2
	matchups := []Matchup{}
	for _, n := range nodeList {
		cs := n.children
		if len(cs) >= 2 {
			matchups = append(matchups, Matchup{cs[len(cs)-1], cs[len(cs)-2]})
		}
	}
	if debug {
		fmt.Printf("matchups: %v\n", matchups)
	}
	return matchups
}

func faceoff(matchup Matchup) (winner, loser *Node) {
	fmt.Printf("a: %s\nb: %s\n", matchup.A, matchup.B)

	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		log.Fatal(err)
	}
	if input == "a" || input == "A" {
		winner = matchup.A
		loser = matchup.B
	} else {
		winner = matchup.B
		loser = matchup.A
	}
	//fmt.Printf("%s > %s\n\n", winner, loser)
	return
}

func readOptions(root *Node) {
	fmt.Println("Enter one option per line. Enter a blank line when done.")

	scanner := bufio.NewScanner(os.Stdin)

	//count := 0
	for {
		// count += 1
		// fmt.Printf("Option %d: ", count)
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			log.Fatal(err)
		}
		input := scanner.Text()
		if input == "" {
			break
		}
		//fmt.Printf("added: %s-\n", input)
		createNode(input, root)
	}
}

func runTournament() {
	fmt.Println("Enter a or b to indicate your preference for the following items:")
	var matchups []Matchup
	for {
		matchups = findMatchups()
		if len(matchups) == 0 {
			break
		}
		winner, loser := faceoff(matchups[0])
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
	runTournament()
	showResults(root)
}
