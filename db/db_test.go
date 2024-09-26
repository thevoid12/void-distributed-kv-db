package db_test

import (
	"bytes"
	"main/db"
	"os"
	"testing"
)

func createTempDb(t *testing.T, readOnly bool) *db.Database {
	t.Helper()
	f, err := os.CreateTemp("void-kv-db", "")
	if err != nil {
		t.Fatalf("Could not create temp file: %v", err)
	}
	name := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(name) })

	db, closeFunc, err := db.NewDatabase(name, readOnly)
	if err != nil {
		t.Fatalf("Could not create a new database: %v", err)
	}
	t.Cleanup(func() { closeFunc() })
	return db
}
func TestGetSet(t *testing.T) {
	db := createTempDb(t, false)

	if err := db.SetKey("void", []byte("coder")); err != nil {
		t.Fatalf("Could not write key: %v", err)
	}

	value, err := db.GetKey("void")
	if err != nil {
		t.Fatalf(`Could not get the key "void": %v`, err)
	}

	if !bytes.Equal(value, []byte("coder")) {
		t.Errorf(`Unexpected value for key "void": got %q, want %q`, value, "coder")
	}

	k, v, err := db.GetNextKeyForReplication()
	if err != nil {
		t.Fatalf(`Unexpected error for GetNextKeyForReplication(): %v`, err)
	}
	if !bytes.Equal(k, []byte("void")) || !bytes.Equal(v, []byte("coder")) {
		t.Errorf(`GetNextKeyForReplication(): got %q, %q; want %q, %q`, k, v, "void", "coder")
	}
}
func TestDeleteReplicationKey(t *testing.T) {
	db := createTempDb(t, false)
	setKey(t, db, "void", "coder")
	k, v, err := db.GetNextKeyForReplication()
	if err != nil {
		t.Fatalf(`Unexpected error for GetNextKeyForReplication(): %v`, err)
	}

	if !bytes.Equal(k, []byte("void")) || !bytes.Equal(v, []byte("coder")) {
		t.Errorf(`GetNextKeyForReplication(): got %q, %q; want %q, %q`, k, v, "void", "coder")
	}
	if err := db.DeleteReplicationKey([]byte("void"), []byte("not-coder")); err == nil {
		t.Fatalf(`DeleteReplicationKey("void", "not-coder"): got nil error, want non-nil error`)
	}

	if err := db.DeleteReplicationKey([]byte("void"), []byte("coder")); err != nil {
		t.Fatalf(`DeleteReplicationKey("void", "coder"): got %q, want nil error`, err)
	}

	k, v, err = db.GetNextKeyForReplication()
	if err != nil {
		t.Fatalf(`Unexpected error for GetNextKeyForReplication(): %v`, err)
	}
	if k != nil || v != nil {
		t.Errorf(`GetNextKeyForReplication(): got %v, %v; want nil, nil`, k, v)
	}
}
func TestSetReadOnly(t *testing.T) {
	db := createTempDb(t, true)
	if err := db.SetKey("void", []byte("coder")); err == nil {
		t.Fatalf("SetKey(%q, %q): got nil error, want non-nil error", "void", []byte("coder"))
	}
}

func setKey(t *testing.T, d *db.Database, key, value string) {
	t.Helper()

	if err := d.SetKey(key, []byte(value)); err != nil {
		t.Fatalf("SetKey(%q, %q) failed: %v", key, value, err)
	}
}

func getKey(t *testing.T, d *db.Database, key string) string {
	t.Helper()

	value, err := d.GetKey(key)
	if err != nil {
		t.Fatalf("GetKey(%q) failed: %v", key, err)
	}

	return string(value)
}

func TestDeleteExtraKeys(t *testing.T) {
	db := createTempDb(t, false)

	setKey(t, db, "void", "coder")
	setKey(t, db, "thevoid12", "myusername")

	if err := db.DeleteExtraKeys(func(name string) bool { return name == "thevoid12" }); err != nil {
		t.Fatalf("Could not delete extra keys: %v", err)
	}

	if value := getKey(t, db, "void"); value != "coder" {
		t.Errorf(`Unexpected value for key "void": got %q, want %q`, value, "coder")
	}

	if value := getKey(t, db, "thevoid12"); value != "" {
		t.Errorf(`Unexpected value for key "thevoid12": got %q, want %q`, value, "")
	}
}
