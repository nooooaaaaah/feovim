package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/lucasb-eyer/go-colorful"
)

type fileItem struct {
	filename string
	path     string
	modified string
	isDir    bool
}

func (f fileItem) Title() string {
	if f.isDir {
		return f.filename + "/"
	}
	return f.filename
}
func (f fileItem) Description() string { return "Modified: " + f.modified }
func (f fileItem) FilterValue() string { return f.filename }

type model struct {
	list           list.Model
	currentPath    string
	previousPath   string
	showPreview    bool
	previewContent string
	viewport       viewport.Model
	previewIsImage bool
	imageContent   string
}

func (m model) Init() tea.Cmd {
	return nil
}

func openFile(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Run()
}

func getRootPath() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}

func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif"
}

func imageToAscii(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var sb strings.Builder
	for y := 0; y < height; y += 2 {
		for x := 0; x < width; x += 2 {
			r, g, b, _ := img.At(x, y).RGBA()
			c := colorful.Color{R: float64(r) / 65535, G: float64(g) / 65535, B: float64(b) / 65535}
			_, _, v := c.Hsv()
			char := "█"
			if v < 0.3 {
				char = " "
			} else if v < 0.6 {
				char = "▒"
			} else if v < 0.9 {
				char = "█"
			}
			sb.WriteString(char)
		}
		sb.WriteString("\n")
	}
	return sb.String(), nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "p" {
			if i, ok := m.list.SelectedItem().(fileItem); ok && !i.isDir {
				if isImageFile(i.path) {
					asciiArt, err := imageToAscii(i.path)
					if err != nil {
						log.Error("Failed to convert image", "error", err)
						return m, nil
					}
					m.showPreview = true
					m.previewIsImage = true
					m.imageContent = asciiArt
					m.viewport = viewport.New(80, 40)
					m.viewport.SetContent(asciiArt)
				} else {
					content, err := os.ReadFile(i.path)
					if err != nil {
						log.Error("Failed to read file", "error", err)
						return m, nil
					}
					m.showPreview = true
					m.previewIsImage = false
					m.previewContent = string(content)
					m.viewport = viewport.New(80, 40)
					m.viewport.SetContent(string(content))
				}
				return m, nil
			}
		}
		if m.showPreview {
			switch msg.String() {
			case "esc":
				m.showPreview = false
				return m, nil
			case "up", "k":
				m.viewport.LineUp(1)
			case "down", "j":
				m.viewport.LineDown(1)
			case "pgup":
				m.viewport.HalfViewUp()
			case "pgdown":
				m.viewport.HalfViewDown()
			}
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}
		if msg.String() == "enter" || msg.String() == "l" {
			if i, ok := m.list.SelectedItem().(fileItem); ok {
				if i.isDir {
					m.previousPath = m.currentPath
					m.currentPath = i.path
					items, err := getFiles(i.path)
					if err != nil {
						log.Error("Failed to get files", "error", err)
						return m, nil
					}
					m.list.SetItems(items)
				} else {
					if err := openFile(i.path); err != nil {
						log.Error("Failed to open file", "error", err)
					}
				}
			}
		}
		if msg.String() == "backspace" || msg.String() == "h" {
			parentDir := filepath.Dir(m.currentPath)
			if parentDir != m.currentPath {
				items, err := getFiles(parentDir)
				if err != nil {
					log.Error("Failed to get files", "error", err)
					return m, nil
				}
				m.previousPath = m.currentPath
				m.currentPath = parentDir
				m.list.SetItems(items)
			}
		}
		if msg.String() == "home" {
			rootPath := getRootPath()
			items, err := getFiles(rootPath)
			if err != nil {
				log.Error("Failed to get files", "error", err)
				return m, nil
			}
			m.previousPath = m.currentPath
			m.currentPath = rootPath
			m.list.SetItems(items)
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		if m.showPreview {
			m.viewport.Width = msg.Width - h - 2
			m.viewport.Height = msg.Height - v - 2
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.showPreview {
		if m.previewIsImage {
			return lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Padding(1).
				Render("Image Preview:\n\n" + m.viewport.View())
		}
		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1).
			Render("File Preview:\n\n" + m.viewport.View())
	}
	return docStyle.Render(m.list.View())
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func getFiles(dir string) ([]list.Item, error) {
	var items []list.Item

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	if filepath.Dir(dir) != dir {
		parentItem := fileItem{
			filename: "..",
			path:     filepath.Dir(dir),
			modified: "",
			isDir:    true,
		}
		items = append(items, parentItem)
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		file := fileItem{
			filename: info.Name(),
			path:     filepath.Join(dir, info.Name()),
			modified: info.ModTime().Format("2006-01-02 15:04"),
			isDir:    entry.IsDir(),
		}
		items = append(items, file)
	}

	return items, nil
}

func main() {
	log.SetReportCaller(true)
	log.SetTimeFormat(time.Kitchen)

	dir := getRootPath()

	items, err := getFiles(dir)
	if err != nil {
		log.Error("Failed to get files", "error", err)
		return
	}

	m := model{
		list:        list.New(items, list.NewDefaultDelegate(), 0, 0),
		currentPath: dir,
	}
	m.list.Title = "File Explorer"

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		return
	}
}
