package db_test

import (
	"bytes"
	"main/pkg/db"
	"os"
	"testing"
)

func TestGetSet(t *testing.T) {
	f, err := os.CreateTemp("void-kv-db", "")
	if err != nil {
		t.Fatalf("Could not create temp file: %v", err)
	}
	name := f.Name()
	f.Close()
	defer os.Remove(name)
	db, closeFunc, err := db.NewDatabase(name)
	if err != nil {
		t.Fatalf("Could not create a new database: %v", err)
	}
	defer closeFunc()
	if err := db.SetKey("party", []byte("Great")); err != nil {
		t.Fatalf("Could not write key: %v", err)
	}
	value, err := db.GetKey("party")
	if err != nil {
		t.Fatalf(`Could not get the key "party": %v`, err)
	}
	if !bytes.Equal(value, []byte("Great")) {
		t.Errorf(`Unexpected value for key "party": got %q, want %q`, value, "Great")
	}
}
