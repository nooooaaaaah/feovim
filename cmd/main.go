package main

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	zone "github.com/lrstanley/bubblezone"
	"github.com/nooooaaaaah/photoboard/internal/explorer"
	"github.com/nooooaaaaah/photoboard/internal/model"
	"github.com/nooooaaaaah/photoboard/internal/ui"
	"github.com/nooooaaaaah/photoboard/internal/utils"
)

type ModelWrapper struct {
	model model.Model
}

func (m ModelWrapper) Init() tea.Cmd {
	return m.model.Init()
}

func (m ModelWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if msg.Action != tea.MouseActionRelease {
			return m, nil
		}

		newModel, cmd := m.model.Update(msg)
		if model, ok := newModel.(model.Model); ok {
			m.model = model
			return m, cmd
		}

	case tea.KeyMsg:
		newModel, cmd := m.model.Update(msg)
		if model, ok := newModel.(model.Model); ok {
			m.model = model
			return m, cmd
		}
	}

	newModel, cmd := m.model.Update(msg)
	if model, ok := newModel.(model.Model); ok {
		m.model = model
	}
	return m, cmd
}

func (m ModelWrapper) View() string {
	return zone.Scan(m.model.View())
}

func main() {
	log.SetReportCaller(true)
	log.SetTimeFormat(time.Kitchen)

	// Initialize global zone manager
	zone.NewGlobal()
	defer zone.Close()

	dir := utils.GetRootPath()

	styler := ui.NewDefaultStyler()
	nav := explorer.Navigator{}
	prev := explorer.Previewer{}
	uiHandler := ui.NewWindowHandler(styler)

	// Create model with all dependencies
	m := model.NewModel(dir, styler, nav, prev, uiHandler)

	// Configure initial delegate for the first column
	delegate := list.NewDefaultDelegate()
	delegate.SetSpacing(0)
	delegate.ShowDescription = false
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("205"))

	// Initialize the first column with proper width
	initialWidth := 30 // This will be adjusted by window resize
	err := m.AddColumn(dir, initialWidth)
	if err != nil {
		log.Error("Failed to create initial column", "error", err)
		return
	}

	// Create wrapper
	wrapper := ModelWrapper{
		model: m,
	}

	p := tea.NewProgram(
		&wrapper,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		return
	}
}
