package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	"github.com/nooooaaaaah/photoboard/internal/defs"
	// To fix circular imports:
	// 1. Move FileItem struct definition to a separate package (e.g. "types")
	// 2. Update imports in both packages to use "types" package instead
	// 3. Remove direct dependencies between utils and model packages
)

func GetFiles(dir string) ([]list.Item, error) {
	var items []list.Item

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	if filepath.Dir(dir) != dir {
		parentItem := defs.FileItem{
			Filename: "..",
			Path:     filepath.Dir(dir),
			Modified: "",
			IsDir:    true,
		}
		items = append(items, parentItem)
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		file := defs.FileItem{
			Filename: info.Name(),
			Path:     filepath.Join(dir, info.Name()),
			Modified: info.ModTime().Format("2006-01-02 15:04"),
			IsDir:    entry.IsDir(),
		}
		items = append(items, file)
	}

	sort.Slice(items, func(i, j int) bool {
		itemI := items[i].(defs.FileItem)
		itemJ := items[j].(defs.FileItem)

		// Always keep parent directory (..) at the top
		if itemI.Filename == ".." {
			return true
		}
		if itemJ.Filename == ".." {
			return false
		}

		// Directories before files
		if itemI.IsDir != itemJ.IsDir {
			return itemI.IsDir
		}

		// Alphabetical order within same type
		return itemI.Filename < itemJ.Filename
	})

	return items, nil
}

func OpenFile(path string) error {
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

func GetRootPath() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}

func GetFilesTree(dir string, level int) ([]list.Item, error) {
	var items []list.Item

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// Add parent directory entry if not at root
	if level == 0 && filepath.Dir(dir) != dir {
		parentItem := defs.TreeItem{
			Filename: "..",
			Path:     filepath.Dir(dir),
			Modified: "",
			IsDir:    true,
			Level:    level,
		}
		items = append(items, parentItem)
	}

	// First add directories
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		item := defs.TreeItem{
			Filename: info.Name(),
			Path:     filepath.Join(dir, info.Name()),
			Modified: info.ModTime().Format("2006-01-02 15:04"),
			IsDir:    true,
			Level:    level,
			IsOpen:   false,
		}
		items = append(items, item)
	}

	// Then add files
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		item := defs.TreeItem{
			Filename: info.Name(),
			Path:     filepath.Join(dir, info.Name()),
			Modified: info.ModTime().Format("2006-01-02 15:04"),
			IsDir:    false,
			Level:    level,
		}
		items = append(items, item)
	}

	return items, nil
}
