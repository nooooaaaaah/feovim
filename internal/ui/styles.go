package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/nooooaaaaah/photoboard/internal/defs"
)

var DocStyle = lipgloss.NewStyle().Margin(1, 2)
var ImgPreviewStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1)
var FileViewStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1)

// DefaultStyler implements the Styler interface
type DefaultStyler struct {
	docStyle          lipgloss.Style
	filePreview       lipgloss.Style
	imagePreview      lipgloss.Style
	hoverStyle        lipgloss.Style
	columnStyle       lipgloss.Style
	activeColumnStyle lipgloss.Style
}

var _ defs.Styler = (*DefaultStyler)(nil) // Ensure interface implementation

func NewDefaultStyler() *DefaultStyler {
	return &DefaultStyler{
		docStyle: lipgloss.NewStyle().Margin(0, 0),
		filePreview: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1),
		imagePreview: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1),
		columnStyle: lipgloss.NewStyle().
			BorderRight(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Margin(0, 0).
			Padding(0, 0),
		activeColumnStyle: lipgloss.NewStyle().
			BorderRight(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("205")).
			Margin(0, 0).
			Padding(0, 0),
	}
}

func (s *DefaultStyler) ListStyle() lipgloss.Style {
	return s.docStyle
}

func (s *DefaultStyler) FilePreviewStyle() lipgloss.Style {
	return s.filePreview
}

func (s *DefaultStyler) ImagePreviewStyle() lipgloss.Style {
	return s.imagePreview
}

func (s *DefaultStyler) GetFrameSize() (width, height int) {
	return s.docStyle.GetFrameSize()
}

func (s *DefaultStyler) ColumnStyle() lipgloss.Style {
	return s.columnStyle
}

func (s *DefaultStyler) ActiveColumnStyle() lipgloss.Style {
	return s.activeColumnStyle
}
