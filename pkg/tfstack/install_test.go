package tfstack

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_detectTFVersion(t *testing.T) {
	// test-stack tf file location
	stackPath := "testdata/test-stack/"

	want := "1.2.6"
	got, err := detectTFVersion(stackPath)

	assert.NoError(t, err, "Unexpected error from detectTFVersion")
	assert.Equal(t, want, got, "Terraform version does not match expected output")
}

func Test_downloadBinary(t *testing.T) {

	tfver := "1.2.6"
	want := os.TempDir() + tfver + "/terraform"
	got, err := downloadBinary(os.TempDir()+tfver, tfver)

	assert.NoError(t, err, "Unexpected error from downloadBinary")
	assert.Equal(t, want, got, "Terraform binary path does not match expected output")
}
