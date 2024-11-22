package utils

import (
	"path/filepath"
	"strings"

	"github.com/qeesung/image2ascii/convert"
)

func IsImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif"
}

func ImageToAscii(path string) string {
	converter := convert.NewImageConverter()
	options := convert.DefaultOptions
	options.FixedWidth = 80
	options.FixedHeight = 40

	res := converter.ImageFile2ASCIIString(path, &options)
	return res
}
