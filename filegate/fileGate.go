package filegate

import (
	"os"
	"sync"
)

type FileGate struct {
	filePath string
	perm os.FileMode
	mu sync.Mutex
}

func (fileGate *FileGate) Open() (*ReadWriteCloser, error) {
	fileGate.mu.Lock()
	f, err := os.OpenFile(fileGate.filePath, os.O_CREATE|os.O_RDWR, fileGate.perm)
	if err != nil {
		fileGate.mu.Unlock()
		return nil, err
	}
	rwc := ReadWriteCloser{
		fileFd: f,
		fileGateUnlock: fileGate.mu.Unlock,
	}
	return &rwc, nil
}


func NewFileGate(filePath string, perm os.FileMode) *FileGate {
	fileGate := FileGate{
		filePath: filePath,
		perm: perm,
	}
	return &fileGate
}