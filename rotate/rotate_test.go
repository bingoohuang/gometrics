package rotate

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testWrite(t *testing.T, f *File, path string) {
	n, err := io.WriteString(f, "hello\n")
	assert.NoError(t, err)
	assert.Equal(t, n, 6)
	assert.Equal(t, path, f.Path)

	n, err = f.Write([]byte("bar\n"))
	assert.NoError(t, err)
	assert.Equal(t, n, 4)

	err = f.Close()
	assert.NoError(t, err)

	d, err := ioutil.ReadFile(path)
	assert.NoError(t, err)
	assert.Equal(t, string(d), "hello\nbar\n")

	n, err = f.Write([]byte("and more\n"))
	f.Flush()
	assert.NoError(t, err)

	assert.Equal(t, 9, n)
}

func TestBasic(t *testing.T) {
	os.RemoveAll("test_dir")
	defer os.RemoveAll("test_dir")

	path := filepath.Join("test_dir", "second.log")
	f, err := NewFile(path)
	assert.NoError(t, err)
	assert.NotNil(t, f)
	testWrite(t, f, path)

	err = f.Close()
	assert.NoError(t, err)

	d, err := ioutil.ReadFile(path)
	assert.NoError(t, err)
	assert.Equal(t, string(d), "hello\nbar\nand more\n")
}

// nolint gomnd
func TestBasic_Location(t *testing.T) {
	os.RemoveAll("test_dir")
	defer os.RemoveAll("test_dir")

	path := filepath.Join("test_dir", "third.log")
	f, err := NewFile(path)
	assert.NoError(t, err)

	n, err := io.WriteString(f, "hello\n")
	assert.NoError(t, err)
	assert.Equal(t, n, 6)
	assert.Equal(t, path, f.Path)
}
