package utils

import (
	"os"
	"testing"
)

func Test_replaceEnvironment(t *testing.T) {

	os.Setenv("ENV1", "myenvone")
	os.Setenv("ENV2", "myenvtwo")
	os.Setenv("ENV3", "myenvthree")

	rs := "test${ENV1}plus${ENV2}${ENV3}tail"
	want := "test" + os.Getenv("ENV1") + "plus" + os.Getenv("ENV2") + os.Getenv("ENV3") + "tail"

	if got := replaceEnvironment(rs); got != want {
		t.Errorf("replaceEnvironment() = %v, want %v", got, want)
	}

}

func Test_replaceEnvironment1(t *testing.T) {

	os.Setenv("ENV1", "myenvone")

	rs := "${ENV1}"
	want := os.Getenv("ENV1")

	if got := replaceEnvironment(rs); got != want {
		t.Errorf("replaceEnvironment() = %v, want %v", got, want)
	}

}
func Test_replaceEnvironment2(t *testing.T) {

	os.Setenv("ENV1", "myenvone")

	rs := "a${ENV1}b"
	want := "a" + os.Getenv("ENV1") + "b"

	if got := replaceEnvironment(rs); got != want {
		t.Errorf("replaceEnvironment() = %v, want %v", got, want)
	}

}
