package filegate

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"testing"
	"time"
)

func TestNewFileGate(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	fileGate := NewFileGate(path.Join(os.TempDir(), `test_new_file_gate`), 0644)
	rwc, err := fileGate.Open()
	if err != nil { t.Fatal(err.Error()) }
	fmt.Fprint(rwc, `Hello World!`)
	b, err := io.ReadAll(rwc)
	if string(b) != `Hello World!` {
		t.Fatal(string(b))
	}
	c := make(chan string, 0)
	go func() {
		rwc1, err := fileGate.Open()
		defer rwc1.Close()
		if err != nil { t.Fatal(err.Error()) }
		p, err := io.ReadAll(rwc1)
		if err != nil { t.Fatal(err.Error()) }
		c <- string(p)
	}()
	time.Sleep(time.Second)
	fmt.Fprint(rwc, `Hello World again!`)
	rwc.Close()
	if s := <-c; s != `Hello World again!` {
		t.Fatal(s)
	}
}
