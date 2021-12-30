package rwblocks

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"os"
)

var (
	ErrorLog = log.New(os.Stderr, `error#`, log.Lshortfile)
	ErrShortBlock = errors.New(`Block is too short`)
)

func Write(w io.Writer, block []byte) error {
	blockLength := len(block)
	if blockLength == 0 { return nil }
	blockLengthBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(blockLengthBytes,  uint32(blockLength))
	if _, err := w.Write(blockLengthBytes); err != nil {
		ErrorLog.Println(err.Error())
		return err
	}
	if _, err := w.Write(block); err != nil {
		ErrorLog.Println(err.Error())
		return err
	}
	return nil
}


func Read(r io.Reader) ([]byte, error) {
	blockLengthBytes := make([]byte, 4)
	if _, err := r.Read(blockLengthBytes); err != nil {
		if err != io.EOF { ErrorLog.Println(err.Error()) }
		return []byte{}, err
	}
	blockLength := binary.LittleEndian.Uint32(blockLengthBytes)
	block := make([]byte, blockLength)
	n, err := r.Read(block)
	if err != nil {
		ErrorLog.Println(err.Error(), n)
		return []byte{}, err
	}
	if n < int(blockLength) {
		ErrorLog.Println(ErrShortBlock, blockLength, n)
		return []byte{}, ErrShortBlock
	}
	return block, nil
}