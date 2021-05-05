package timeslice

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestSaveLoad(t *testing.T) {
	file, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	initial := []int64{34853452, 234582934, 234589234059, time.Now().Unix()}
	if err := Save(initial, file.Name()); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, err := Load(file.Name())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !reflect.DeepEqual(got, initial) {
		t.Errorf("Loaded() = %v, want %v", got, initial)
	}
}
