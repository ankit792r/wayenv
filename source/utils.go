package main

import (
	"os"
	"strings"
)

func (m model) pageSize() int {
	size := m.height - 2 // footer/help

	if size < 1 {
		size = 1
	}

	return size
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
