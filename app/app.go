package app

import (
	"fmt"
	"log"
	"math"
	"sort"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Node struct {
	label    string
	stamp    int
	children []*Node
}

func (n Node) String() string {
	return fmt.Sprintf("%s", n.label)
}

func (m *model) createNode(label string) *Node {
	node := Node{label: label}
	m.nodeList = append(m.nodeList, &node)

	m.prefer(m.rootNode, &node)

	return &node
}

func (m *model) prefer(parent, child *Node) {
	m.stampSeq += 1
	parent.stamp = m.stampSeq
	child.stamp = m.stampSeq

	parent.children = append(parent.children, child)

	// order by time stamp, so we present less recent options first
	sort.Slice(parent.children, func(i, j int) bool {
		return parent.children[i].stamp < parent.children[j].stamp
	})

	// walkthrough entire node list, and find nodes that include both parent and child as children and remove child
	for _, n := range m.nodeList {
		m.preferSibling(n, parent, child)
	}
}

func (m *model) preferSibling(node, keep, drop *Node) {
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

func (m *model) printNodes(root *Node) {
	fmt.Println("digraph {")
	for _, n := range m.nodeList {
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

func (m *model) findMatchups() []Matchup {
	matchups := []Matchup{}
	for _, n := range m.nodeList {
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

func (m *model) faceoff(root *Node, matchup *Matchup) (winner, loser *Node) {
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
		m.printNodes(root)
	}
	return m.faceoff(root, matchup)
}

func (m *model) setMatchup() {
	matchups := m.findMatchups()
	if len(matchups) == 0 {
		m.matchup = nil
		m.selected = nil
	} else {
		m.matchup = &matchups[0]
		m.selected = m.matchup.A
	}
}

func (m model) orderedNodes() (result []*Node) {
	node := m.rootNode
	for {
		if len(node.children) == 0 {
			break
		}
		result = append(result, node.children[0])
		node = node.children[0]
	}
	return
}

type model struct {
	state     string
	rootNode  *Node
	nodeList  []*Node
	textInput textinput.Model
	matchup   *Matchup
	selected  *Node
	stampSeq  int
}

func InitialModel() model {
	ti := textinput.New()
	ti.Placeholder = "New item..."
	ti.Focus()

	rootNode := Node{label: "0"}

	return model{
		state:     "collect",
		textInput: ti,
		rootNode:  &rootNode,
		nodeList:  []*Node{&rootNode},
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) View() string {
	switch m.state {
	case "collect":
		switch len(m.nodeList) {
		case 1:
			m.textInput.Placeholder = "First item"
		case 2:
			m.textInput.Placeholder = "Second item"
		default:
			m.textInput.Placeholder = "Another item... hit enter when done"
		}
		return fmt.Sprintf(
			"Let's rank some stuff\n%s%s",
			m.viewNodes(),
			m.textInput.View(),
		) + "\n"
	case "rank":
		s := "Which do you prefer?\n\n"
		if m.selected == m.matchup.A {
			s += "> "
		} else {
			s += "  "
		}
		s += m.matchup.A.label
		s += "\n"
		if m.selected == m.matchup.B {
			s += "> "
		} else {
			s += "  "
		}
		s += m.matchup.B.label
		return s
	case "results":
		s := "Here are the results:\n\n"
		for i, node := range m.orderedNodes() {
			s += fmt.Sprintf("%d\t%s\n", i+1, node)
		}

		s += "\nPress r to redo with the same items\n"
		s += "Press s to start over with new items\n"
		s += "Press q to quit\n"

		return s
	}

	return fmt.Sprintf("unknown state: %s\n", m.state)
}

func (m model) viewNodes() string {
	s := ""
	for i, node := range m.nodeList {
		if node != m.rootNode {
			s += fmt.Sprintf("%d: %s\n", i, node)
		}
	}
	return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	switch m.state {
	case "collect":
		return m.UpdateCollect(msg)
	case "rank":
		return m.UpdateRank(msg)
	case "results":
		return m.UpdateResults(msg)
	}
	return m, nil
}

func (m model) UpdateRank(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			m.selected = m.matchup.A
		case tea.KeyDown:
			m.selected = m.matchup.B
		case tea.KeyEnter:
			if m.selected == m.matchup.A {
				m.prefer(m.selected, m.matchup.B)
			} else {
				m.prefer(m.selected, m.matchup.A)
			}
			m.setMatchup()
			if m.matchup == nil {
				m.state = "results"
			}
		}
	}
	return m, nil
}

func (m model) UpdateCollect(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			label := m.textInput.Value()
			if label == "" {
				m.state = "rank"
				m.setMatchup()
			} else {
				m.createNode(label)
				m.textInput.Reset()
			}
			return m, nil
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) UpdateResults(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "s":
				return InitialModel(), nil
			case "q":
				return m, tea.Quit
			}
		}
	}

	return m, nil
}
