package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Page int

const (
	Auswahl Page = iota
	Connect
)

type sshEntry struct {
	IP       string
	Port     int
	Username string
}

type model struct {
	textInput  textinput.Model
	sshEntrys  []sshEntry
	cursor     int
	page       Page
	selectedIP sshEntry
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
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

	var sshEntrys []sshEntry

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
		entry := sshEntry{IP: IP, Port: PORT}
		if !slices.Contains(sshEntrys, entry) {
			sshEntrys = append(sshEntrys, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	ti := textinput.New()
	ti.Focus()

	return model{
		sshEntrys:  sshEntrys,
		cursor:     0,
		page:       Auswahl,
		selectedIP: sshEntry{},
		textInput:  ti,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
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
		case "space", "enter":
			if m.page == Auswahl {
				m.selectedIP = m.sshEntrys[m.cursor]
				m.page = Connect
			}
		case "backspace":
			if m.page > 0 {
				m.page--
			}
		}
	}

	if m.page == Connect {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	s := "SSHTUI\n\n"

	switch m.page {
	case Auswahl:
		s += "Auswahl der IPs:\n"

		// iterate over IPs
		for i, entry := range m.sshEntrys {
			// Is the cursor pointing at this choice?
			cursor := " " // no cursor
			if m.cursor == i {
				cursor = ">" // cursor!
			}

			// Render the row
			s += fmt.Sprintf("%s %s\n", cursor, entry.IP)
		}

	case Connect:
		s += fmt.Sprintf("Connecten zu IP: %s Port: %d\n", m.selectedIP.IP, m.selectedIP.Port)
		s += fmt.Sprintf("Username %s\n", m.textInput.View())

	}
	s += "\nPress q to quit"

	return s
}
