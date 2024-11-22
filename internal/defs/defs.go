package defs

import "strings"

type FileItem struct {
	Filename string
	Path     string
	Modified string
	IsDir    bool
}

func (f FileItem) Title() string {
	if f.IsDir {
		return f.Filename + "/"
	}
	return f.Filename
}
func (f FileItem) Description() string { return "Modified: " + f.Modified }
func (f FileItem) FilterValue() string { return f.Filename }

type TreeItem struct {
    Filename string
    Path     string
    Modified string
    IsDir    bool
    Level    int      // Indentation level
    Children []TreeItem
    IsOpen   bool     // Whether the folder is expanded
}

func (t TreeItem) Title() string {
    prefix := strings.Repeat("  ", t.Level)
    if t.IsDir {
        if t.IsOpen {
            return prefix + "▼ " + t.Filename + "/"
        }
        return prefix + "▶ " + t.Filename + "/"
    }
    return prefix + "  " + t.Filename
}

func (t TreeItem) Description() string { return "Modified: " + t.Modified }
func (t TreeItem) FilterValue() string { return t.Filename }
