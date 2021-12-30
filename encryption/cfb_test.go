package encryption

import (
	"bytes"
	"testing"
)

func TestCFBKey(t *testing.T) {
	cfb, err := NewCFB(nil, nil)
	if err != nil { t.Fatal(err.Error()) }
	if !bytes.Equal(cfb.key, DefaultKey) {
		t.Log(`Is not default key`, cfb.key, DefaultKey)
		t.Fail()
	}

	cfb, err = NewCFB(nil, []byte{})
	if err != nil {
		t.Log(err.Error())
		t.Fail()
		return
	}
	if !bytes.Equal(cfb.key, DefaultKey) {
		t.Log(`Is not default key`, cfb.key, DefaultKey)
		t.Fail()
	}

	cfb, err = NewCFB(nil, []byte{1, 2, 3, 4})
	if err != nil { t.Fatal(err.Error()) }
	if !bytes.Equal(cfb.key[:6], []byte{1, 2, 3, 4, 0, 0}) {
		t.Log(`Is not equal`, cfb.key, DefaultKey)
		t.Fail()
	}
}

func TestCFBEncryptDecrypt(t *testing.T) {
	cfb, err := NewCFB(nil, []byte{1, 2, 3, 4})
	if err != nil { t.Fatal(err.Error()) }

	a := []byte(`Hello world!`)
	b := cfb.Encrypt(a)
	t.Log(b)
	c := cfb.Decrypt(b)
	t.Log(c)
	t.Log(string(c))
	if !bytes.Equal(a, c) { t.Fail() }
}
