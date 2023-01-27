package fileutil

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

const (
	// ContentTypeExcel excel
	ContentTypeExcel = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	// ContentTypePDF pdf
	ContentTypePDF = "application/pdf"
)

// File file
type File struct {
	filename    string
	basePath    string
	md5         string
	contentType string
}

// New new file
func New(mf *multipart.FileHeader) (*File, error) {
	src, err := mf.Open()
	if err != nil {
		return nil, err
	}
	defer func() { _ = src.Close() }()

	td, err := ioutil.TempDir("", "boots-")
	if err != nil {
		return nil, err
	}

	dst, err := os.Create(fmt.Sprintf("%s/%s", td, mf.Filename))
	if err != nil {
		return nil, err
	}
	defer func() { _ = dst.Close() }()

	if _, err = io.Copy(dst, src); err != nil {
		return nil, err
	}
	hash := md5.New()
	hashInBytes := hash.Sum(nil)[:16]
	md5 := hex.EncodeToString(hashInBytes)
	return &File{
		filename: mf.Filename,
		basePath: td,
		md5:      strings.ToUpper(md5),
	}, nil
}

// NewWithFilename new file
func NewWithFilename(name string) (*File, error) {
	td, err := ioutil.TempDir("", "boots-")
	if err != nil {
		return nil, err
	}
	var ct string
	switch filepath.Ext(name) {
	case ".xlsx":
		ct = ContentTypeExcel
	case ".pdf":
		ct = ContentTypePDF
	}
	return &File{
		filename:    name,
		basePath:    td,
		contentType: ct,
	}, nil
}

// Path file path
func (f *File) Path() string {
	return fmt.Sprintf("%s/%s", f.basePath, f.filename)
}

// Name file name
func (f *File) Name() string {
	return f.filename
}

// MD5 md5
func (f *File) MD5() string {
	return f.md5
}

// Close close
func (f *File) Close() error {
	return os.RemoveAll(f.basePath)
}

// ContentType content type
func (f *File) ContentType() string {
	return f.contentType
}

// Ext ext
func (f *File) Ext() string {
	return filepath.Ext(f.filename)
}
