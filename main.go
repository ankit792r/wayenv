package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) pageSize() int {
	size := m.height - 2 // footer/help

	if size < 1 {
		size = 1
	}

	return size
}

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

func getAllOsEnvs() []Env {

	var envSlice []Env

	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		env := Env{
			Name:  pair[0],
			Value: pair[1],
		}
		envSlice = append(envSlice, env)
	}
	return envSlice
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

func (m *model) updateScroll() {
	pageSize := m.pageSize()
	if m.selected < m.offset {
		m.offset = m.selected
	}

	if m.selected >= m.offset+pageSize {
		m.offset = m.selected - pageSize + 1
	}

	if m.selected >= m.offset+pageSize {
		m.offset = m.selected - pageSize + 1
	}
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

func (m model) renderList() string {
	var b strings.Builder
	nameWidth := m.nameColumnWidth()
	pageSize := m.pageSize()

	if m.selected >= m.offset+pageSize {
		m.offset = m.selected - pageSize + 1
	}

	end := min(m.offset+pageSize, len(m.envs))

	for i := m.offset; i < end; i++ {
		env := m.envs[i]

		prefix := " "
		if i == m.selected {
			prefix = ">"
		}

		line := fmt.Sprintf(
			"%s %-*s%s",
			prefix,
			nameWidth,
			env.Name,
			env.Value,
		)
		b.WriteString(line)
		b.WriteByte('\n')
	}

	b.WriteString("\n")
	b.WriteString("↑↓ navigate  a add  e edit  d delete  q quit")

	return b.String()
}

func (m model) renderForm() string {
	title := "Add Environment Variable"

	if m.editing {
		title = "Edit Environment Variable"
	}

	return fmt.Sprintf(
		"%s\n\nName : %s\nValue: %s\n\n"+
			"tab switch field\n"+
			"enter save\n"+
			"esc cancel",
		title,
		m.nameInput.View(),
		m.valueInput.View(),
	)
}

func (m model) renderDelete() string {
	if len(m.envs) == 0 {
		return ""
	}

	env := m.envs[m.selected]

	return fmt.Sprintf(
		"Delete %s?\n\n"+
			"y confirm\n"+
			"n cancel",
		env.Name,
	)
}

func (m model) nameColumnWidth() int {
	maxLen := 0

	for _, env := range m.envs {
		if len(env.Name) > maxLen {
			maxLen = len(env.Name)
		}
	}

	// longest name + 2 tabs (assuming tab = 4 spaces)
	return maxLen + 8
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
