package dotenv_test

import (
	"github.com/golobby/dotenv/v2"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type FlagBox struct {
	Bool1 bool `env:"BOOL1"`
	Bool2 bool `env:"BOOL2"`
	Bool3 bool `env:"BOOL3"`
	Bool4 bool `env:"BOOL4"`
}

type Config struct {
	AppName  string  `env:"APP_NAME"`
	AppPort  string  `env:"APP_PORT"`
	float    float64 `env:"FLOAT"`
	FlagBox  *FlagBox
	QuoteBox struct {
		Quote1 string `env:"QUOTE1"`
		Quote2 string `env:"QUOTE2"`
	}
}

func TestLoad(t *testing.T) {
	f, err := os.Open("resources/test/.env")
	assert.NoError(t, err)

	c := &Config{}
	c.FlagBox = &FlagBox{}

	err = dotenv.Load(f, c)
	assert.NoError(t, err)

	assert.Equal(t, "DotEnv", c.AppName)
	assert.Equal(t, "8585", c.AppPort)
	assert.Equal(t, 3.14, c.float)
	assert.Equal(t, true, c.FlagBox.Bool1)
	assert.Equal(t, false, c.FlagBox.Bool2)
	assert.Equal(t, true, c.FlagBox.Bool3)
	assert.Equal(t, false, c.FlagBox.Bool4)
	assert.Equal(t, "OK1", c.QuoteBox.Quote1)
	assert.Equal(t, " OK 2 ", c.QuoteBox.Quote2)

	err = f.Close()
	assert.NoError(t, err)
}

func TestLoad_With_Invalid_File(t *testing.T) {
	f, err := os.Open("resources/test/.env.buggy1")
	assert.NoError(t, err)

	c := &Config{}
	err = dotenv.Load(f, c)
	assert.Errorf(t, err, "dotenv: invalid syntax in line 1")

	err = f.Close()
	assert.NoError(t, err)
}

func TestLoad_With_Invalid_Field_It_Should_Fail(t *testing.T) {
	f, err := os.Open("resources/test/.env")
	assert.NoError(t, err)

	sample := struct {
		BOOL1 bool `env:"APP_NAME"`
	}{}
	err = dotenv.Load(f, &sample)
	assert.Error(t, err)

	err = f.Close()
	assert.NoError(t, err)
}

func TestLoad_With_Invalid_Nested_Structure_Field_It_Should_Fail(t *testing.T) {
	f, err := os.Open("resources/test/.env")
	assert.NoError(t, err)

	type Inner struct {
		Float float64 `env:"APP_NAME"`
	}

	outer := struct {
		Inner Inner
	}{}
	outer.Inner = Inner{}

	err = dotenv.Load(f, &outer)
	assert.Error(t, err)

	err = f.Close()
	assert.NoError(t, err)
}

func TestLoad_With_Invalid_Nested_Structure_Ptr_Field_It_Should_Fail(t *testing.T) {
	f, err := os.Open("resources/test/.env")
	assert.NoError(t, err)

	type Inner struct {
		Float float64 `env:"APP_NAME"`
	}

	outer := struct {
		Inner *Inner
	}{}
	outer.Inner = &Inner{}

	err = dotenv.Load(f, &outer)
	assert.Error(t, err)

	err = f.Close()
	assert.NoError(t, err)
}
