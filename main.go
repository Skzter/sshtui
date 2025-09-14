package main

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type sshEntry struct {
	IP       string
	Port     int
	Username string
}

type model struct {
	sshEntrys []sshEntry
	cursor    int
	selected  map[int]sshEntry
}

func initialModel() model {

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	sshDir := home + "/.ssh/known_hosts"
	data, err := os.Open(sshDir)
	if err != nil {
		panic(err)
	}
	defer data.Close()

	scanner := bufio.NewScanner(data)

	sshEntrys := make(map[string]sshEntry)

	for scanner.Scan() {
		item := scanner.Text()
		startIP := strings.Index(item, "[")
		if startIP == -1 {
			continue
		}
		endIP := strings.Index(item, "]")
		IP := item[startIP+1 : endIP]
		startPort := strings.Index(item, ":")
		endPort := strings.Index(item, " ")
		PORT, err := strconv.Atoi(item[startPort+1 : endPort])
		if err != nil {
			panic(err)
		}
		sshEntrys[IP] = sshEntry{IP: IP, Port: PORT}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	var entryArr []sshEntry
	for _, value := range sshEntrys {
		entryArr = append(entryArr, value)
	}
	return model{
		sshEntrys: entryArr,
		cursor:    0,
		selected:  make(map[int]sshEntry),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
			// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.sshEntrys)-1 {
				m.cursor++
			}
		case "enter", "space":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = m.sshEntrys[m.cursor]
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "SSHTUI\n\n"
	// iterate over IPs
	for i, entry := range m.sshEntrys {
		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, entry.IP)
	}
	s += "\nPress q to quit"

	return s
}
