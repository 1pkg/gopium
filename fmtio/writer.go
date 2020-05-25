package fmtio

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"1pkg/gopium"
)

// Stdout defines writer implementation
// which only returns os stdout all the time
type Stdout struct{}

// Generate stdout implementation
func (Stdout) Generate(string) (io.WriteCloser, error) {
	return stdout{}, nil
}

// File defines writer implementation
// which creates underlying single file
// with provided name and ext on provided loc
type File struct {
	Name string
	Ext  string
}

// Generate file implementation
func (f File) Generate(loc string) (io.WriteCloser, error) {
	path := filepath.Dir(loc)
	return os.Create(filepath.Join(path, fmt.Sprintf("%s.%s", f.Name, f.Ext)))
}

// Files defines writer implementation
// which creates underlying files list
// with provided ext on provided loc
type Files struct {
	Ext string
}

// Generate files implementation
func (f Files) Generate(loc string) (io.WriteCloser, error) {
	path := strings.Replace(loc, filepath.Ext(loc), fmt.Sprintf(".%s", f.Ext), 1)
	return os.Create(path)
}

// Origin defines category writer implementation
// which simply uses underlying writter
type Origin struct {
	Writter gopium.Writer
}

// Category origin implementation
func (o Origin) Category(cat string) error {
	return nil
}

// Generate origin implementation
func (o Origin) Generate(loc string) (io.WriteCloser, error) {
	return o.Writter.Generate(loc)
}

// Suffix defines category writer implementation
// which replaces category for writer
// with provided suffixed category
type Suffix struct {
	Writter gopium.Writer
	Suffix  string
	oldcat  string
	newcat  string
}

// Category suffix implementation
func (s *Suffix) Category(cat string) error {
	// add suffix to category
	scat := fmt.Sprintf("%s_%s", cat, s.Suffix)
	s.oldcat = cat
	s.newcat = scat
	// create dir for new suffixed cat
	return os.MkdirAll(scat, os.ModePerm)
}

// Generate suffix implementation
func (s Suffix) Generate(loc string) (io.WriteCloser, error) {
	loc = strings.Replace(loc, s.oldcat, s.newcat, 1)
	return s.Writter.Generate(loc)
}
