package store

import (
	"os"
	"testing"
)

const (
	DB_PATH = "test.db"
)

var db *BoltStore

func init() {
	os.Setenv("VICTOR_STORAGE_PATH", DB_PATH)
}

func setup() {
	os.Create(DB_PATH)
	if db == nil {
		db = newBoltStore()
	}
}

func teardown() {
	os.Remove(DB_PATH)
}

func TestEmptyGet(t *testing.T) {
	setup()

	val, _ := db.Get("nothing")
	if val != "" {
		t.Error("Expected to get nothing before store has data, got: ", val)
	}

	teardown()
}

func TestSet(t *testing.T) {
	setup()

	db.Set("a", "b")
	val, _ := db.Get("a")

	if val != "b" {
		t.Error("Stored 'a': 'b', expected to get it back", val)
	}

	teardown()
}

func TestDelete(t *testing.T) {
	setup()

	db.Set("a", "b")
	db.Delete("a")
	val, _ := db.Get("a")
	if val != "" {
		t.Error("Expected to get nothing after deleting key, got: ", val)
	}

	teardown()
}
