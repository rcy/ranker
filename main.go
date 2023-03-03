package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
)

type Node struct {
	label    string
	stamp    int
	children []*Node
}

var stampSeq int

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
	stampSeq += 1
	parent.stamp = stampSeq
	child.stamp = stampSeq

	parent.children = append(parent.children, child)

	// order by time stamp, so we present less recent options first
	sort.Slice(parent.children, func(i, j int) bool {
		return parent.children[i].stamp < parent.children[j].stamp
	})

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
		children := []*Node{}
		for _, n := range node.children {
			if n != drop {
				children = append(children, n)
			}
		}
		sort.Slice(children, func(i, j int) bool {
			return children[i].stamp < children[j].stamp
		})
		node.children = children
	}
}

func printNodes(root *Node) {
	fmt.Println("digraph {")
	for _, n := range nodeList {
		if len(n.children) > 0 {
			for _, c := range n.children {
				if n == root {
					fmt.Printf("  \"%s\"", c)
				} else {
					fmt.Printf("  \"%s\" -> \"%s\"", n, c)
				}
				fmt.Println()
			}
			fmt.Println()
		}
	}
	fmt.Println("}")
}

type Matchup struct {
	A *Node
	B *Node
}

func findMatchups() []Matchup {
	matchups := []Matchup{}
	for _, n := range nodeList {
		if len(n.children) >= 2 {
			// within this sibling group take the first 2, as they are sorted by stamp, this will be the oldest pair
			matchups = append(matchups, Matchup{n.children[0], n.children[1]})
		}
	}
	// sort the matchups, by oldest recent, then oldest secondary
	sort.Slice(matchups, func(i, j int) bool {
		imax := math.Max(float64(matchups[i].A.stamp), float64(matchups[i].B.stamp))
		jmax := math.Max(float64(matchups[j].A.stamp), float64(matchups[j].B.stamp))
		if imax == jmax {
			imin := math.Min(float64(matchups[i].A.stamp), float64(matchups[i].B.stamp))
			jmin := math.Min(float64(matchups[j].A.stamp), float64(matchups[j].B.stamp))
			return imin < jmin
		} else {
			return imax < jmax
		}
	})
	return matchups
}

func faceoff(root *Node, matchup *Matchup) (winner, loser *Node) {
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

func nextMatchup() *Matchup {
	matchups := findMatchups()
	//fmt.Printf("MATCHUPS: %v\n", matchups)
	if len(matchups) == 0 {
		return nil
	}
	return &matchups[0]
}

func runTournament(root *Node) {
	fmt.Println("Enter a or b to indicate your preference for the following items:")
	var matchup *Matchup
	for {
		matchup = nextMatchup()
		if matchup == nil {
			break
		}
		winner, loser := faceoff(root, matchup)
		fmt.Printf("%s > %s\n\n", winner, loser)
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
