package defs

import "github.com/charmbracelet/lipgloss"

type Styler interface {
	ColumnStyle() lipgloss.Style
	ActiveColumnStyle() lipgloss.Style
	FilePreviewStyle() lipgloss.Style
	ImagePreviewStyle() lipgloss.Style
	GetFrameSize() (width, height int)
}
