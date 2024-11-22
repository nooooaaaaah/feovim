package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nooooaaaaah/photoboard/internal/defs"
	"github.com/nooooaaaaah/photoboard/internal/model"
)

type WindowHandler struct {
	styler defs.Styler
}

func NewWindowHandler(styler defs.Styler) *WindowHandler {
	return &WindowHandler{
		styler: styler,
	}
}

func (wh *WindowHandler) HandleWindowResize(m model.Model, msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.WindowWidth = msg.Width
	m.WindowHeight = msg.Height

	if len(m.Columns) == 0 {
		return m, nil
	}

	// Calculate available width
	availableWidth := msg.Width - 2
	minColumnWidth := 30
	maxColumns := availableWidth / minColumnWidth
	if maxColumns < 1 {
		maxColumns = 1
	}

	// Calculate width for visible columns
	visibleColumns := len(m.Columns)
	if visibleColumns > maxColumns {
		visibleColumns = maxColumns
	}
	columnWidth := availableWidth / visibleColumns

	// Update widths for all columns
	for i := range m.Columns {
		m.Columns[i].Width = columnWidth
		m.Columns[i].List.SetSize(columnWidth-2, msg.Height-2)
	}

	return m, nil
}
