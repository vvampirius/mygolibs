package rwblocks

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"
)


func TestWriteReadEqual(t *testing.T) {
	x := []byte{1, 2, 3, 4, 5}
	buffer := bytes.NewBuffer(nil)
	if err := Write(buffer, x); err != nil { t.Fatal(err.Error()) }
	if len(buffer.Bytes()) != len(x) + 4 { t.Fail() }
	y, err := Read(buffer)
	if err != nil { t.Fatal(err.Error()) }
	if !bytes.Equal(x, y) { t.Fail() }
}

func TestReadBadBlock(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{5, 0, 0 ,0, 1, 2, 3, 4})
	if n, err := Read(buffer); err == nil { t.Fatal(n, err) }
}

func TestWriteFileReadToEOF(t *testing.T) {
	tmpDir, err := ioutil.TempDir(``, ``)
	if err != nil { t.Fatal(err.Error()) }
	defer os.RemoveAll(tmpDir)

	filePath := path.Join(tmpDir, `TestWriteFileReadToEOF`)
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil { t.Fatal(err.Error()) }
	defer f.Close()

	a := []byte{1, 2, 3, 4, 5}
	if err := Write(f, a); err != nil { t.Fatal(err.Error()) }
	if _, err := f.Seek(0, 0); err != nil { t.Fatal(err.Error()) }
	b, err := Read(f)
	if err != nil { t.Fatal(err.Error()) }
	if !bytes.Equal(a, b) { t.FailNow() }

	if _, err := Read(f); err == nil || err != io.EOF { t.Fatal(err) }
}
