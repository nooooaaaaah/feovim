package model

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/nooooaaaaah/photoboard/internal/defs"
	"github.com/nooooaaaaah/photoboard/internal/utils"
)

type Navigator interface {
	HandleNavigation(Model, tea.KeyMsg) (tea.Model, tea.Cmd)
}

type Previewer interface {
	HandlePreviewUpdate(Model, tea.KeyMsg) (tea.Model, tea.Cmd)
	StartPreview(Model, tea.KeyMsg) (tea.Model, tea.Cmd)
}

type UIHandler interface {
	HandleWindowResize(Model, tea.WindowSizeMsg) (tea.Model, tea.Cmd)
}

type ColumnView struct {
	List     list.Model
	Path     string
	Selected string
	Width    int
}

type Model struct {
	Columns        []ColumnView
	ActiveColumn   int
	ShowPreview    bool
	PreviewContent string
	Viewport       viewport.Model
	PreviewIsImage bool
	imageContent   string
	Styler         defs.Styler
	navigator      Navigator
	previewer      Previewer
	uiHandler      UIHandler
	WindowWidth    int
	WindowHeight   int
}

func NewModel(path string, styler defs.Styler, nav Navigator, prev Previewer, ui UIHandler) Model {
	return Model{
		Columns:      make([]ColumnView, 0),
		ActiveColumn: 0,
		Styler:       styler,
		navigator:    nav,
		previewer:    prev,
		uiHandler:    ui,
		WindowWidth:  80,
		WindowHeight: 24,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if m.ShowPreview {
			return m.previewer.HandlePreviewUpdate(m, msg)
		}

		switch msg.String() {
		case "p":
			return m.previewer.StartPreview(m, msg)
		case "enter", "l", "backspace", "h", "home":
			return m.navigator.HandleNavigation(m, msg)
		}

	case tea.WindowSizeMsg:
		return m.uiHandler.HandleWindowResize(m, msg)

	case tea.MouseMsg:
		if msg.Action != tea.MouseActionRelease {
			return m, nil
		}

		// Handle preview mode clicks
		if m.ShowPreview {
			if zone.Get("exit-preview").InBounds(msg) {
				m.ShowPreview = false
				return m, nil
			}
			return m, nil
		}

		// Handle list item clicks
		if m.ActiveColumn < len(m.Columns) {
			activeList := m.Columns[m.ActiveColumn].List
			for i := range activeList.Items() {
				if zone.Get(fmt.Sprintf("item-%d", i)).InBounds(msg) {
					activeList.Select(i)
					if msg.Button == tea.MouseButtonLeft {
						if i, ok := activeList.SelectedItem().(defs.FileItem); ok {
							if i.IsDir {
								return m.navigator.HandleNavigation(m, tea.KeyMsg{Type: tea.KeyEnter})
							} else {
								return m.previewer.StartPreview(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
							}
						}
					}
					return m, nil
				}
			}

			// Update active column's list
			var cmd tea.Cmd
			m.Columns[m.ActiveColumn].List, cmd = m.Columns[m.ActiveColumn].List.Update(msg)
			return m, cmd
		}
	}

	// If we have an active column, update its list
	if m.ActiveColumn < len(m.Columns) {
		var cmd tea.Cmd
		m.Columns[m.ActiveColumn].List, cmd = m.Columns[m.ActiveColumn].List.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if m.ShowPreview {
		if m.PreviewIsImage {
			return m.Styler.ImagePreviewStyle().Render(m.Viewport.View())
		}
		return m.Styler.FilePreviewStyle().Render(m.Viewport.View())
	}

	if len(m.Columns) == 0 {
		return "No columns to display"
	}

	// Calculate total available width
	availableWidth := m.WindowWidth - 2 // Account for margins
	minColumnWidth := 30
	maxColumns := availableWidth / minColumnWidth
	if maxColumns < 1 {
		maxColumns = 1
	}

	// Determine which columns to display
	startCol := m.ActiveColumn - (maxColumns - 1)
	if startCol < 0 {
		startCol = 0
	}
	endCol := startCol + maxColumns
	if endCol > len(m.Columns) {
		endCol = len(m.Columns)
		startCol = endCol - maxColumns
		if startCol < 0 {
			startCol = 0
		}
	}

	// Calculate width for each visible column
	visibleColumns := endCol - startCol
	columnWidth := availableWidth / visibleColumns

	var columns []string
	for i := startCol; i < endCol; i++ {
		col := m.Columns[i]
		col.Width = columnWidth // Update column width

		style := m.Styler.ColumnStyle()
		if i == m.ActiveColumn {
			style = m.Styler.ActiveColumnStyle()
		}

		// Create column header
		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			Width(columnWidth - 2).
			Background(lipgloss.Color("240"))

		header := headerStyle.Render(filepath.Base(col.Path))

		// Create list items
		var items []string
		for j, item := range col.List.Items() {
			if fileItem, ok := item.(defs.FileItem); ok {
				itemStyle := lipgloss.NewStyle().
					Width(columnWidth-2).
					Padding(0, 1)

				if j == col.List.Index() && i == m.ActiveColumn {
					itemStyle = itemStyle.
						Background(lipgloss.Color("205")).
						Foreground(lipgloss.Color("0"))
				}

				prefix := "  "
				if fileItem.IsDir {
					prefix = "▶ "
				}

				zoneID := fmt.Sprintf("item-%d", j)
				itemContent := zone.Mark(zoneID, itemStyle.Render(prefix+fileItem.Filename))
				items = append(items, itemContent)
			}
		}

		columnContent := lipgloss.JoinVertical(lipgloss.Left, append([]string{header}, items...)...)
		columns = append(columns, style.Render(columnContent))
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, columns...)
}

func (m *Model) AddColumn(path string, width int) error {
	// Calculate how many columns can fit
	minColumnWidth := 30
	maxColumns := (m.WindowWidth - 2) / minColumnWidth
	if maxColumns < 1 {
		maxColumns = 1
	}

	items, err := utils.GetFiles(path)
	if err != nil {
		return err
	}

	delegate := list.NewDefaultDelegate()
	delegate.SetSpacing(0)
	delegate.ShowDescription = false

	newList := list.New(items, delegate, width, 0)
	newList.SetShowTitle(false)
	newList.SetFilteringEnabled(false)
	newList.SetShowStatusBar(false)
	newList.SetShowHelp(false)

	column := ColumnView{
		List:  newList,
		Path:  filepath.Base(path),
		Width: width,
	}

	// If we're at max columns, remove leftmost column
	if len(m.Columns) >= maxColumns {
		m.Columns = m.Columns[1:]
		m.ActiveColumn--
	}

	m.Columns = append(m.Columns, column)
	return nil
}
