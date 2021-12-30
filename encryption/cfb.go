package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

var (
	ErrKeyLength = errors.New(`Wrong key length`)
)

type CFB struct {
	IV, key []byte
	block cipher.Block
}

func (cfb *CFB) SetKey(key []byte) error {
	if key == nil { key = make([]byte, 0) }
	normalizedKey := WriteInBytes(key, DefaultKey)

	normalizedKeyLength := len(normalizedKey)
	if normalizedKeyLength != 16 && normalizedKeyLength != 24 && normalizedKeyLength != 32 {
		ErrorLog.Println(ErrKeyLength, normalizedKeyLength)
		return ErrKeyLength
	}

	block, err := aes.NewCipher(normalizedKey)
	if err != nil {
		ErrorLog.Println(err.Error())
		return err
	}

	cfb.key = normalizedKey
	cfb.block = block
	return nil
}


func (cfb *CFB) SetIV(iv []byte) error {
	//TODO: fail if no block
	if iv == nil || len(iv) == 0 {
		iv = make([]byte, cfb.block.BlockSize())
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			ErrorLog.Println(err.Error())
			return err
		}
	}
	cfb.IV = iv
	return nil
}

func (cfb *CFB) Encrypt(data []byte) []byte {
	if data == nil || len(data) == 0 { return []byte{} }
	// TODO: exit if no iv or block
	encrypted := make([]byte, len(data))
	encrypter := cipher.NewCFBEncrypter(cfb.block, cfb.IV)
	encrypter.XORKeyStream(encrypted, data)
	return encrypted
}

func (cfb *CFB) Decrypt(encrypted []byte) []byte {
	if encrypted == nil || len(encrypted) == 0 { return []byte{} }
	// TODO: exit if no iv or block
	data := make([]byte, len(encrypted))
	decrypter := cipher.NewCFBDecrypter(cfb.block, cfb.IV)
	decrypter.XORKeyStream(data, encrypted)
	return data
}


func NewCFB(iv, key []byte) (*CFB, error) {
	cfb := CFB{}
	if err := cfb.SetKey(key); err != nil { return nil, err }
	if err := cfb.SetIV(iv); err != nil { return nil, err }

	return &cfb, nil
}