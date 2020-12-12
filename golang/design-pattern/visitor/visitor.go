package visitor

import (
	"fmt"
	"path"
)

type Visitor interface {
	Visit(IResourceFile) error
}

type IResourceFile interface {
	Accept(Visitor) error
}

func NewResourceFile(filepath string) (IResourceFile, error) {
	switch path.Ext(filepath) {
	case ".ppt":
		return &PPTFile{path: filepath}, nil
	case ".pdf":
		return &PDFFile{path: filepath}, nil
	default:
		return nil, fmt.Errorf("not found file type: %s", filepath)
	}
}

type PPTFile struct {
	path string
}

func (f *PPTFile) Accept(visitor Visitor) error {
	return visitor.Visit(f)
}

type PDFFile struct {
	path string
}

func (f *PDFFile) Accept(visitor Visitor) error {
	return visitor.Visit(f)
}

type Compressor struct{}

func (c Compressor) Visit(r IResourceFile) error {
	switch f := r.(type) {
	case *PPTFile:
		return c.VisitPPTFile(f)
	case *PDFFile:
		return c.VisitPDFFile(f)
	default:
		return fmt.Errorf("not found resource typr: %#v", r)
	}	
}

func (c *Compressor) VisitPPTFile(f *PPTFile) error {
	fmt.Println("this is ppt file")
	return nil
}

// VisitPDFFile VisitPDFFile
func (c *Compressor) VisitPDFFile(f *PdfFile) error {
	fmt.Println("this is pdf file")
	return nil
}

