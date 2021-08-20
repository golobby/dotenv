package dotenv_test

import (
	"github.com/golobby/dotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type FlagBox struct {
	Bool1 bool `dotenv:"BOOL1"`
	Bool2 bool `dotenv:"BOOL2"`
	Bool3 bool `dotenv:"BOOL3"`
	Bool4 bool `dotenv:"BOOL4"`
}

type Config struct {
	AppName  string  `dotenv:"APP_NAME"`
	AppPort  int32   `dotenv:"APP_PORT"`
	float    float64 `dotenv:"FLOAT"`
	FlagBox  *FlagBox
	QuoteBox struct {
		Quote1 string `dotenv:"QUOTE1"`
		Quote2 string `dotenv:"QUOTE2"`
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
	assert.Equal(t, int32(8585), c.AppPort)
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

func TestLoad_With_Default_Value(t *testing.T) {
	f, err := os.Open("resources/test/.env")
	assert.NoError(t, err)

	type Config struct {
		AppName string `dotenv:"APP_NAME"`
		AppUrl  string `dotenv:"APP_URL"`
	}

	c := &Config{}
	c.AppUrl = "https://example.com"

	err = dotenv.Load(f, c)
	assert.NoError(t, err)

	assert.Equal(t, "DotEnv", c.AppName)
	assert.Equal(t, "https://example.com", c.AppUrl)

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

func TestLoad_With_Non_Readable_File(t *testing.T) {
	f, _ := os.OpenFile("resources/test/.env.invalid", os.O_APPEND, 0644)

	c := &Config{}
	err := dotenv.Load(f, c)
	assert.Error(t, err)
}

func TestLoad_With_Invalid_Structure(t *testing.T) {
	f, err := os.Open("resources/test/.env")
	assert.NoError(t, err)

	var number int
	err = dotenv.Load(f, &number)
	assert.Errorf(t, err, "dotenv: invalid structure")

	err = f.Close()
	assert.NoError(t, err)
}

func TestLoad_With_Invalid_Field_It_Should_Fail(t *testing.T) {
	f, err := os.Open("resources/test/.env")
	assert.NoError(t, err)

	sample := struct {
		BOOL1 bool `dotenv:"APP_NAME"`
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
		Float float64 `dotenv:"APP_NAME"`
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
		Float float64 `dotenv:"APP_NAME"`
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
