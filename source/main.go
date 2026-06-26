package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type mode int

const (
	listMode mode = iota
	formMode
	deleteMode
)

type Env struct {
	Name  string
	Value string
}

type model struct {
	mode mode

	envs []Env

	selected int
	offset   int

	width  int
	height int

	nameInput  textinput.Model
	valueInput textinput.Model

	editing bool
	editIdx int
}

func initialModel() model {
	nameInput := textinput.New()
	nameInput.Placeholder = "ENV_NAME"
	nameInput.Focus()

	valueInput := textinput.New()
	valueInput.Placeholder = "value"

	envs := getAllOsEnvs()

	return model{
		mode:       listMode,
		envs:       envs,
		nameInput:  nameInput,
		valueInput: valueInput,
		editIdx:    -1,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:

		switch m.mode {

		case listMode:

			switch msg.String() {

			case "q", "ctrl+c":
				return m, tea.Quit

			case "up", "k":
				if m.selected > 0 {
					m.selected--
					m.updateScroll()
				}

			case "down", "j":
				if m.selected < len(m.envs)-1 {
					m.selected++
					m.updateScroll()
				}

			case "a":
				m.mode = formMode
				m.editing = false
				m.editIdx = -1

				m.nameInput.SetValue("")
				m.valueInput.SetValue("")

				m.nameInput.Focus()
				m.valueInput.Blur()

			case "e":
				if len(m.envs) == 0 {
					return m, nil
				}

				env := m.envs[m.selected]

				m.mode = formMode
				m.editing = true
				m.editIdx = m.selected

				m.nameInput.SetValue(env.Name)
				m.valueInput.SetValue(env.Value)

				m.nameInput.Focus()
				m.valueInput.Blur()

			case "d":
				if len(m.envs) > 0 {
					m.mode = deleteMode
				}
			}

		case deleteMode:

			switch msg.String() {

			case "y", "enter":
				if len(m.envs) > 0 {
					m.envs = append(
						m.envs[:m.selected],
						m.envs[m.selected+1:]...,
					)

					if m.selected >= len(m.envs) && m.selected > 0 {
						m.selected--
					}

					m.updateScroll()
				}

				m.mode = listMode

			case "n", "esc":
				m.mode = listMode
			}

		case formMode:

			switch msg.String() {

			case "esc":
				m.mode = listMode
				return m, nil

			case "tab":
				if m.nameInput.Focused() {
					m.nameInput.Blur()
					m.valueInput.Focus()
				} else {
					m.valueInput.Blur()
					m.nameInput.Focus()
				}
				return m, nil

			case "enter":

				name := strings.TrimSpace(m.nameInput.Value())
				value := strings.TrimSpace(m.valueInput.Value())

				if name == "" {
					return m, nil
				}

				env := Env{
					Name:  name,
					Value: value,
				}

				if m.editing {
					m.envs[m.editIdx] = env
				} else {
					m.envs = append(m.envs, env)
					m.selected = len(m.envs) - 1
					m.updateScroll()
				}

				m.mode = listMode
				return m, nil
			}

			var cmd tea.Cmd

			if m.nameInput.Focused() {
				m.nameInput, cmd = m.nameInput.Update(msg)
			} else {
				m.valueInput, cmd = m.valueInput.Update(msg)
			}

			return m, cmd
		}
	}

	return m, nil
}

func (m model) View() string {
	switch m.mode {

	case formMode:
		return m.renderForm()

	case deleteMode:
		return m.renderDelete()

	default:
		return m.renderList()
	}
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
