package p2p

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestDeleteStoreKey(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	store := NewStore(opts)

	data := "Hello there"
	key := "myDataKey"
	if err := store.WriteStream(key, bytes.NewBufferString(data)); err != nil {
		t.Error(err)
	}

	if err := store.Delete(key); err != nil {
		t.Error(err)
	}
}

func TestTransformFunc(t *testing.T) {
	key := "mykey"
	pathKey := CASPathTransformFunc(key)
	expectedOriginalKey := "816cc20437d859538736e1ef46558b7bda486c06"
	expectedPathName := "816cc/20437/d8595/38736/e1ef4/6558b/7bda4/86c06"
	if pathKey.PathName != expectedPathName {
		t.Errorf("Expected: %s\nBut was: %s\n", expectedPathName, pathKey.PathName)
	}
	if pathKey.FileName != expectedOriginalKey {
		t.Errorf("Expected: %s\nBut was: %s\n", expectedOriginalKey, pathKey.FileName)
	}
	// assert.Equal(t, expectedPathName, pathName)
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	store := NewStore(opts)

	data := "Hello there"
	key := "myDataKey"
	if err := store.WriteStream(key, bytes.NewBufferString(data)); err != nil {
		t.Error(err)
	}

	r, err := store.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)
	retrivedData := string(b)
	fmt.Println(retrivedData)

	if data != retrivedData {
		t.Errorf("Expected: %s\nBut was: %s\n", data, retrivedData)
	}

	if err := store.Delete(key); err != nil {
		t.Error(err)
	}
}
