package explorer

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nooooaaaaah/photoboard/internal/defs"
	"github.com/nooooaaaaah/photoboard/internal/model"
)

type Navigator struct{}

func (n Navigator) HandleNavigation(m model.Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	minColumnWidth := 30
	maxColumns := (m.WindowWidth - 2) / minColumnWidth
	if maxColumns < 1 {
		maxColumns = 1
	}

	switch msg.String() {
	case "right", "l", "enter":
		if m.ActiveColumn >= len(m.Columns) {
			return m, nil
		}

		if i, ok := m.Columns[m.ActiveColumn].List.SelectedItem().(defs.FileItem); ok {
			if i.IsDir {
				// Remove columns to the right of current
				m.Columns = m.Columns[:m.ActiveColumn+1]

				// If we're at max columns, remove leftmost column
				if len(m.Columns) >= maxColumns {
					m.Columns = m.Columns[1:]
					m.ActiveColumn--
				}

				// Add new column
				columnWidth := (m.WindowWidth - 2) / (len(m.Columns) + 1)
				if columnWidth < minColumnWidth {
					columnWidth = minColumnWidth
				}

				err := m.AddColumn(i.Path, columnWidth)
				if err == nil {
					m.ActiveColumn++
				}
			}
		}
		return m, nil

	case "left", "h", "backspace":
		if m.ActiveColumn > 0 {
			m.ActiveColumn--
			// Remove columns to the right
			m.Columns = m.Columns[:m.ActiveColumn+1]
		}
		return m, nil

	case "up", "k":
		if m.ActiveColumn < len(m.Columns) {
			m.Columns[m.ActiveColumn].List.CursorUp()
		}
		return m, nil

	case "down", "j":
		if m.ActiveColumn < len(m.Columns) {
			m.Columns[m.ActiveColumn].List.CursorDown()
		}
		return m, nil
	}
	return m, nil
}

func removeChildren(items []list.Item, parentIdx int, parentLevel int) []list.Item {
	result := make([]list.Item, 0)
	result = append(result, items[:parentIdx+1]...)

	for i := parentIdx + 1; i < len(items); i++ {
		if item, ok := items[i].(defs.TreeItem); ok {
			if item.Level <= parentLevel {
				result = append(result, items[i:]...)
				break
			}
		}
	}

	return result
}
