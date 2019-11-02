package gitdiff

import (
	"fmt"
	"os"
	"strings"
	"log"
)

// File describes changes to a single file. It can be either a text file or a
// binary file.
type File struct {
	OldName string
	NewName string

	IsNew    bool
	IsDelete bool
	IsCopy   bool
	IsRename bool

	OldMode os.FileMode
	NewMode os.FileMode

	OldOIDPrefix string
	NewOIDPrefix string
	Score        int

	// TextFragments contains the fragments describing changes to a text file. It
	// may be empty if the file is empty or if only the mode changes.
	TextFragments []*TextFragment

	// IsBinary is true if the file is a binary file. If the patch includes
	// binary data, BinaryFragment will be non-nil and describe the changes to
	// the data. If the patch is reversible, ReverseBinaryFragment will also be
	// non-nil and describe the changes needed to restore the original file
	// after applying the changes in BinaryFragment.
	IsBinary              bool
	BinaryFragment        *BinaryFragment
	ReverseBinaryFragment *BinaryFragment
}

func (f *File) HasFunc(match func(string)bool) bool {
	for _, fragment := range f.TextFragments {
		funcnames := fragment.FuncNames()
		for _, fname := range funcnames {
			if match(fname) {
				return true
			}
		}
	}
	return false
}

// TextFragment describes changed lines starting at a specific line in a text file.
type TextFragment struct {
	Comment string

	OldPosition int64
	OldLines    int64

	NewPosition int64
	NewLines    int64

	LinesAdded   int64
	LinesDeleted int64

	LeadingContext  int64
	TrailingContext int64

	Lines []Line
}

// Header returns the canonical header of this fragment.
func (f *TextFragment) Header() string {
	return fmt.Sprintf("@@ -%d,%d +%d,%d @@ %s", f.OldPosition, f.OldLines, f.NewPosition, f.NewLines, f.Comment)
}

func (f *TextFragment) FuncNames() []string {
	funcnames := []string{}
	addname := func(line string) {
		if strings.HasPrefix(line, "func") == false || strings.Contains(def, "(") == false {
			return
		}
		def := strings.TrimRight(line, " {")
		// func NAME(arg1, arg2, ...)
		if strings.Count(def, "(") == 1 && len(def) > 5 {
			funcnames = append(funcnames, def[5:])
			return
		}
		// func (this *Object) NAME(arg1, arg2, ...)
		low := strings.Index(def, ")") + 2
		if low >= len(def) {
			log.Printf("Can't handle '%s'")
			return
		}
		hig := strings.Index(def[low:], "(") + low
		if low < hig && low >= 0 && hig < len(def) {
			funcnames = append(funcnames, def[low:hig])
		} else {
			log.Printf("Can't handle '%s'")
		}
	}
	addname(f.Comment)
	for _, line := range f.Lines {
		addname(line.Line)
	}
	return funcnames
}

// Line is a line in a text fragment.
type Line struct {
	Op   LineOp
	Line string
}

func (fl Line) String() string {
	return fl.Op.String() + fl.Line
}

// LineOp describes the type of a text fragment line: context, added, or removed.
type LineOp int

const (
	// OpContext indicates a context line
	OpContext LineOp = iota
	// OpDelete indicates a deleted line
	OpDelete
	// OpAdd indicates an added line
	OpAdd
)

func (op LineOp) String() string {
	switch op {
	case OpContext:
		return " "
	case OpDelete:
		return "-"
	case OpAdd:
		return "+"
	}
	return "?"
}

// BinaryFragment describes changes to a binary file.
type BinaryFragment struct {
	Method BinaryPatchMethod
	Size   int64
	Data   []byte
}

// BinaryPatchMethod is the method used to create and apply the binary patch.
type BinaryPatchMethod int

const (
	// BinaryPatchDelta indicates the data uses Git's packfile encoding
	BinaryPatchDelta BinaryPatchMethod = iota
	// BinaryPatchLiteral indicates the data is the exact file content
	BinaryPatchLiteral
)
