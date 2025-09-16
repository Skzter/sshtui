package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Page int

const (
	Auswahl Page = iota
	Connect
	Execute
)

type sshEntry struct {
	IP       string
	Port     int
	Username string
}

type keyMap struct {
	change  key.Binding
	execute key.Binding
	quit    key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.change, k.execute, k.quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.change, k.execute},
		{k.quit},
	}
}

var keys = keyMap{
	quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
	change: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "change"),
	),
	execute: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "execute"),
	),
}

type model struct {
	textInput    textinput.Model
	sshEntrys    []sshEntry
	cursor       int
	page         Page
	selectedIP   sshEntry
	typing       bool
	keys         keyMap
	help         help.Model
	shouldRunSSH bool
}

func main() {
	m := initialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	if m.shouldRunSSH {
		runSSH(m)
	}
}

func runSSH(m model) {

	args := [...]string{fmt.Sprintf("%s@%s", m.selectedIP.Username, m.selectedIP.IP), fmt.Sprintf("-P %d", m.selectedIP.Port)}
	cmd := exec.Command("ssh", args[0], args[1])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
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

	return model{
		sshEntrys:    sshEntrys,
		cursor:       0,
		page:         Auswahl,
		selectedIP:   sshEntry{},
		textInput:    ti,
		typing:       false,
		keys:         keys,
		help:         help.New(),
		shouldRunSSH: false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			if !m.typing {
				return m, tea.Quit
			}
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
				m.typing = true
				m.textInput.Focus()
			} else if m.page == Connect {
				m.typing = false
				m.selectedIP.Username = strings.TrimSpace(m.textInput.Value())
				m.textInput.Blur()
			}
		case "backspace":
			if m.page == Connect && !m.typing {
				m.textInput.SetValue("")
				m.selectedIP = sshEntry{}
				m.page = Auswahl
			}
		case "c":
			if m.page == Connect && !m.typing {
				m.textInput.Focus()
				m.typing = true
				return m, nil
			}
		case "e":
			if m.page == Connect && !m.typing {
				m.shouldRunSSH = true
				return m, tea.Quit
			}
		}
	}

	if m.page == Auswahl {
		m.keys.change.SetEnabled(false)
		m.keys.execute.SetEnabled(false)
	}
	if m.page == Connect {
		m.keys.change.SetEnabled(true)
		m.keys.execute.SetEnabled(true)
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}
	return m, cmd
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
		if !m.typing {
			if !(m.selectedIP.Username == "") {
				ssh := fmt.Sprintf("ssh %s@%s -P %d", m.selectedIP.Username, m.selectedIP.IP, m.selectedIP.Port)
				s += fmt.Sprintf("command to execute: %s\n", ssh)
			} else {
				s += fmt.Sprintf("no username given, please change\n")
			}
		}
	case Execute:
		return "executing ssh command"
	}
	helpview := m.help.View(m.keys)
	s += "\n" + helpview

	return s
}
