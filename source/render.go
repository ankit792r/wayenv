package main

import (
	"fmt"
	"strings"
)

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
