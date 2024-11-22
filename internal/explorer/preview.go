package explorer

import (
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/nooooaaaaah/photoboard/internal/defs"
	"github.com/nooooaaaaah/photoboard/internal/model"
	"github.com/nooooaaaaah/photoboard/internal/utils"
	"github.com/nooooaaaaah/photoboard/internal/utils/highlight"
)

type Previewer struct{}

func (p Previewer) HandlePreviewUpdate(m model.Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.ShowPreview = false
		return m, nil
	case "up", "k":
		m.Viewport.LineUp(1)
	case "down", "j":
		m.Viewport.LineDown(1)
	case "pgup":
		m.Viewport.HalfViewUp()
	case "pgdown":
		m.Viewport.HalfViewDown()
	}

	var cmd tea.Cmd
	m.Viewport, cmd = m.Viewport.Update(msg)
	return m, cmd
}

func (p Previewer) StartPreview(m model.Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if len(m.Columns) == 0 || m.ActiveColumn >= len(m.Columns) {
		return m, nil
	}

	activeList := m.Columns[m.ActiveColumn].List
	if i, ok := activeList.SelectedItem().(defs.FileItem); ok && !i.IsDir {
		if utils.IsImageFile(i.Path) {
			return handleImagePreview(m, i)
		}
		return handleFilePreview(m, i)
	}
	return m, nil
}

func handleImagePreview(m model.Model, item defs.FileItem) (tea.Model, tea.Cmd) {
	asciiArt := utils.ImageToAscii(item.Path)
	m.ShowPreview = true
	m.PreviewIsImage = true
	m.Viewport = viewport.New(80, 40)
	m.Viewport.SetContent(asciiArt)
	return m, nil
}

func handleFilePreview(m model.Model, item defs.FileItem) (tea.Model, tea.Cmd) {
	content, err := os.ReadFile(item.Path)
	if err != nil {
		log.Error("Failed to read file", "error", err)
		return m, nil
	}

	m.ShowPreview = true
	m.PreviewIsImage = false
	m.PreviewContent = highlight.GetSyntaxHighlightedContent(content, item.Path)
	m.Viewport = viewport.New(80, 40)
	m.Viewport.SetContent(m.PreviewContent)
	return m, nil
}
