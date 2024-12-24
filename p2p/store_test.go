package p2p

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestTransformFunc(t *testing.T) {
	key := "mykey"
	pathKey := CASPathTransformFunc(key)
	expectedFileName := "816cc20437d859538736e1ef46558b7bda486c06"
	expectedPathName := "816cc/20437/d8595/38736/e1ef4/6558b/7bda4/86c06"

	if pathKey.PathName != expectedPathName {
		t.Errorf("Expected: %s\nBut was: %s\n", expectedPathName, pathKey.PathName)
	}

	if pathKey.FileName != expectedFileName {
		t.Errorf("Expected: %s\nBut was: %s\n", expectedFileName, pathKey.FileName)
	}
}

func TestStore(t *testing.T) {
	store := newStore()
	defer teardown(t, store)
	for i := 0; i < 50; i++ {
		data := "Random Data"
		key := fmt.Sprintf("myKey_%d", i)
		if err := store.Write(key, bytes.NewBufferString(data)); err != nil {
			t.Error(err)
		}

		if ok := store.Has(key); !ok {
			t.Error(fmt.Errorf("expected to have key: %s", key))
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

		if ok := store.Has(key); ok {
			t.Error(fmt.Errorf("expected NOT to have key: %s", key))
		}
	}
}

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	return NewStore(opts)
}

func teardown(t *testing.T, store *Store) {
	if err := store.Clear(); err != nil {
		t.Error(err)
	}
}
