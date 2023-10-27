package fs

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/trakval/records"
)

var r records.Records

func TestMain(m *testing.M) {
	r = NewFsRecords("/tmp/fs_records")
	r.Connect()
	c := m.Run()
	r.Close()
	os.RemoveAll("/tmp/fs_records")
	os.Exit(c)
}

func TestCreate(t *testing.T) {
	key := "1"
	frontmatter := map[string]interface{}{
		"k1": "v1",
		"k2": "v2",
	}
	record := map[string]interface{}{
		"frontmatter": frontmatter,
		"body":        "test",
	}
	rk, err := r.CreateRecord(key, record)
	if err != nil {
		t.Errorf("error on create: %v", err)
	}
	if rk != key {
		t.Errorf("error on returned key!")
	}
	fi, err := os.Stat(fmt.Sprintf("/tmp/fs_records/vault/%s.md", key))
	if err != nil {
		t.Errorf("error checking created file: %v", err)
	}
	t.Logf("FileInfo: %v", fi)
}

func TestFetch(t *testing.T) {
	key := "2"
	frontmatter := map[string]interface{}{
		"k1": "v1",
		"k2": "v2",
	}
	record := map[string]interface{}{
		"frontmatter": frontmatter,
		"body":        "test",
	}
	_, err := r.CreateRecord(key, record)
	if err != nil {
		t.Errorf("error on create: %v", err)
	}

	rk, fr, err := r.FetchRecord(key)
	if err != nil {
		t.Errorf("error on fetch: %v", err)
	}
	if rk != key {
		t.Errorf("error on returned key!")
	}
	fm, ok := fr["frontmatter"].(map[string]interface{})
	if !ok {
		t.Errorf("fetched record did not contain frontmatter!")
	}
	if fm["k1"] != "v1" {
		t.Errorf("saved frontmatter value not retrieved properly!")
	}
	if fm["k2"] != "v2" {
		t.Errorf("saved frontmatter value not retrieved properly!")
	}
	body, ok := fr["body"].(string)
	if !ok {
		t.Errorf("fetched record did not contain body!")
	}
	if body != "test" {
		t.Errorf("saved body value not retrieved properly!")
	}
}

func TestUpdate(t *testing.T) {
	key := "3"
	frontmatter := map[string]interface{}{
		"k1": "v1",
		"k2": "v2",
	}
	record := map[string]interface{}{
		"frontmatter": frontmatter,
		"body":        "test",
	}
	_, err := r.CreateRecord(key, record)
	if err != nil {
		t.Errorf("error on create: %v", err)
	}

	rk, fr, err := r.FetchRecord(key)
	if err != nil {
		t.Errorf("error on fetch: %v", err)
	}
	if rk != key {
		t.Errorf("error on returned key!")
	}
	fm, ok := fr["frontmatter"].(map[string]interface{})
	if !ok {
		t.Errorf("fetched record did not contain frontmatter!")
	}
	_, ok = fr["body"].(string)
	if !ok {
		t.Errorf("fetched record did not contain body!")
	}

	fm["k1"] = "v3"
	fm["k2"] = "v4"
	record["frontmatter"] = fm
	record["body"] = "test-update"

	_, err = r.UpdateRecord(key, record)
	if err != nil {
		t.Errorf("error on update: %v", err)
	}

	rk, fr, err = r.FetchRecord(key)
	if err != nil {
		t.Errorf("error on fetch(after update): %v", err)
	}
	if rk != key {
		t.Errorf("error on returned key!")
	}
	fm, ok = fr["frontmatter"].(map[string]interface{})
	if !ok {
		t.Errorf("fetched record did not contain frontmatter(after update)!")
	}
	body, ok := fr["body"].(string)
	if !ok {
		t.Errorf("fetched record did not contain body(after update)!")
	}
	if fm["k1"] != "v3" {
		t.Errorf("saved frontmatter value not retrieved properly(after update)!")
	}
	if fm["k2"] != "v4" {
		t.Errorf("saved frontmatter value not retrieved properly(after update)!")
	}
	if body != "test-update" {
		t.Errorf("saved body value not retrieved properly(after update)!")
	}
}

func TestDelete(t *testing.T) {
	key := "4"
	frontmatter := map[string]interface{}{
		"k1": "v1",
		"k2": "v2",
	}
	record := map[string]interface{}{
		"frontmatter": frontmatter,
		"body":        "test",
	}
	rk, err := r.CreateRecord(key, record)
	if err != nil {
		t.Errorf("error on create: %v", err)
	}
	if rk != key {
		t.Errorf("error on returned key!")
	}
	rk, err = r.DeleteRecord(key)
	if err != nil {
		t.Errorf("error on delete: %v", err)
	}
	if rk != key {
		t.Errorf("error on returned key!")
	}
	_, err = os.Stat(fmt.Sprintf("/tmp/fs_records/vault/%s.md", key))
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("error on delete: %v", err)
	}
}
