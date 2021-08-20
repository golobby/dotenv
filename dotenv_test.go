package dotenv_test

import (
	"github.com/golobby/dotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewDecoder(t *testing.T) {
	f, err := os.Open("./assets/.env")
	assert.NoError(t, err)

	d := dotenv.NewDecoder(f)

	assert.Same(t, f, d.File)
}
