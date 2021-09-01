package tail

import (
	"errors"
	"io"
	"log"
	"os"
	"sync"
	"syscall"
	"time"
)

var (
	DebugLog = log.New(os.Stdout, ``, log.Lshortfile)
	ErrorLog = log.New(os.Stderr, ``, log.Lshortfile)
	ErrIsClosed = errors.New(`tail is closed`)
)

type Tail struct {
	Path string
	fd *os.File
	mu *sync.Mutex
	error error
	ReadBytes uint64
	Deleted uint64
	Moved uint64
	Truncated uint64
	isClosed bool
}

func (tail *Tail) FD() (*os.File, error) {
	if tail.isClosed { return nil, ErrIsClosed }
	if tail.fd != nil { return tail.fd, nil }
	fd, err := os.Open(tail.Path)
	if err != nil {
		ErrorLog.Printf("%s: %s\n", tail.Path, err.Error())
		return nil, err
	}
	fd.Seek(0, 2)
	tail.fd = fd
	return fd, nil
}

func (tail *Tail) Close() error {
	tail.isClosed = true
	if tail.fd != nil {
		fd := tail.fd
		tail.fd = nil
		return fd.Close()
	}
	return nil
}

func (tail *Tail) Read(p []byte) (int, error) {
	tail.mu.Lock()
	defer tail.mu.Unlock()
	for {
		fd, err := tail.FD()
		if err != nil { return 0, err }
		n, err := fd.Read(p)
		tail.ReadBytes = tail.ReadBytes + uint64(n)
		if err != nil && err != io.EOF {
			ErrorLog.Printf("%s: %s\n", tail.Path, err.Error())
			go fd.Close()
			tail.fd = nil
			return n, err
		}
		if n > 0 { return n, nil }
		tail.mu.Unlock()
		time.Sleep(100 * time.Millisecond)
		tail.mu.Lock()
	}
}

func (tail *Tail) IsDeleted() bool {
	fd, err := tail.FD()
	if err != nil { return false }
	fileInfo, err := fd.Stat()
	if err != nil {
		ErrorLog.Printf("%s: %s\n", tail.Path, err.Error())
		return false
	}
	if stat_t, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
		if stat_t.Nlink == 0 { return true }
	}
	return false
}

func (tail *Tail) getInode(sys interface{}) uint64 {
	if stat_t, ok := sys.(*syscall.Stat_t); ok {
		return stat_t.Ino
	}
	return 0
}

func (tail *Tail) IsMoved() bool {
	fd, err := tail.FD()
	if err != nil { return false }
	fileInfo, err := fd.Stat()
	if err != nil {
		ErrorLog.Printf("%s: %s\n", tail.Path, err.Error())
		return false
	}
	openedFileInode := tail.getInode(fileInfo.Sys())

	fileInfo, err = os.Stat(tail.Path)
	if err != nil {
		if !os.IsNotExist(err) { ErrorLog.Printf("%s: %s\n", tail.Path, err.Error()) }
		return true
	}
	currentFileInode := tail.getInode(fileInfo.Sys())

	if currentFileInode != openedFileInode { return true }
	return false
}

func (tail *Tail) IsTruncated() bool {
	fd, err := tail.FD()
	if err != nil { return false }
	fileInfo, err := fd.Stat()
	if err != nil {
		ErrorLog.Printf("%s: %s\n", tail.Path, err.Error())
		return false
	}
	currentPosition, err := fd.Seek(0, 1)
	if err != nil {
		ErrorLog.Printf("%s: %s\n", tail.Path, err.Error())
		return false
	}
	if currentPosition > fileInfo.Size() { return true }
	return false
}

func (tail *Tail) check() {
	for {
		time.Sleep(10 * time.Second)
		tail.mu.Lock()
		closeFD := false
		switch {
			case tail.IsDeleted():
				DebugLog.Printf("%s is deleted\n", tail.Path)
				tail.Deleted++
				closeFD = true
			case tail.IsMoved():
				DebugLog.Printf("%s is moved\n", tail.Path)
				tail.Moved++
				closeFD = true
			case tail.IsTruncated():
				DebugLog.Printf("%s is truncated\n", tail.Path)
				tail.Truncated++
				if tail.fd != nil { tail.fd.Seek(0, 2) }
		}
		if closeFD {
			go tail.fd.Close()
			tail.fd = nil
		}
		tail.mu.Unlock()
	}
}

func NewTail(filePath string) (*Tail, error) {
	tail := Tail{
		Path: filePath,
		mu: new(sync.Mutex),
	}
	if _, err := tail.FD(); err != nil { return nil, err }
	go tail.check()
	return &tail, nil
}