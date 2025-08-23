package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	textInput textinput.Model
	time      int
}

func initalModel() model {
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		time:      10,
	}
}

func (m model) View() string {
	return fmt.Sprintf(
		"What’s your favorite Pokémon %v?\n\n%s\n\n%s",
		m.time,
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", tea.KeyCtrlC.String():
			return m, tea.Quit
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func main() {
	p := tea.NewProgram(model{time: 10})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
