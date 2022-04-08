package decoder_test

import (
	"github.com/golobby/dotenv/pkg/decoder"
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
	AppName  string   `env:"APP_NAME"`
	AppPort  int32    `env:"APP_PORT"`
	IPs      []string `env:"IPS"`
	IDs      []int64  `env:"IDS"`
	float    float64  `env:"FLOAT"`
	FlagBox  *FlagBox
	QuoteBox struct {
		Quote1 string `env:"QUOTE1"`
		Quote2 string `env:"QUOTE2"`
		Quote3 string `env:"QUOTE3"`
		Quote4 string `env:"QUOTE4"`
		Quote5 string `env:"QUOTE5"`
	}
	Multiline string `env:"MULTILINE"`
}

func TestLoad(t *testing.T) {
	f, err := os.Open("./../../assets/.env")
	assert.NoError(t, err)

	c := &Config{}
	c.FlagBox = &FlagBox{}

	err = decoder.Decoder{File: f}.Decode(c)
	assert.NoError(t, err)

	assert.Equal(t, "DotEnv", c.AppName)
	assert.Equal(t, int32(8585), c.AppPort)
	assert.Equal(t, []string{"192.168.0.1", "192.168.0.2", "192.168.0.3"}, c.IPs)
	assert.Equal(t, []int64{10, 11, 12, 13, 14}, c.IDs)
	assert.Equal(t, 3.14, c.float)
	assert.Equal(t, true, c.FlagBox.Bool1)
	assert.Equal(t, false, c.FlagBox.Bool2)
	assert.Equal(t, true, c.FlagBox.Bool3)
	assert.Equal(t, false, c.FlagBox.Bool4)
	assert.Equal(t, "OK1", c.QuoteBox.Quote1)
	assert.Equal(t, " OK 2 ", c.QuoteBox.Quote2)
	assert.Equal(t, " OK ' 3 ", c.QuoteBox.Quote3)
	assert.Equal(t, " OK \" 4 ", c.QuoteBox.Quote4)
	assert.Equal(t, " OK # 5 ", c.QuoteBox.Quote5)
	assert.Equal(t, "1\n2\n3", c.Multiline)

	err = f.Close()
	assert.NoError(t, err)
}

func TestLoad_With_Default_Value(t *testing.T) {
	f, err := os.Open("./../../assets/.env")
	assert.NoError(t, err)

	type Config struct {
		AppName string `env:"APP_NAME"`
		AppUrl  string `env:"APP_URL"`
	}

	c := &Config{}
	c.AppUrl = "https://example.com"

	err = decoder.Decoder{File: f}.Decode(c)
	assert.NoError(t, err)

	assert.Equal(t, "DotEnv", c.AppName)
	assert.Equal(t, "https://example.com", c.AppUrl)

	err = f.Close()
	assert.NoError(t, err)
}

func TestLoad_With_Invalid_File(t *testing.T) {
	f, err := os.Open("./../../assets/.env.buggy")
	assert.NoError(t, err)

	c := &Config{}
	err = decoder.Decoder{File: f}.Decode(c)
	assert.Empty(t, c.AppName)

	err = f.Close()
	assert.NoError(t, err)
}

func TestLoad_With_Non_Readable_File(t *testing.T) {
	f, _ := os.OpenFile("./../../assets/.env.invalid", os.O_APPEND, 0644)

	c := &Config{}
	err := decoder.Decoder{File: f}.Decode(c)
	assert.Error(t, err)
}

func TestLoad_With_Invalid_Structure(t *testing.T) {
	f, err := os.Open("./../../assets/.env")
	assert.NoError(t, err)

	var number int
	err = decoder.Decoder{File: f}.Decode(&number)
	assert.Errorf(t, err, "dotenv: invalid structure")

	err = f.Close()
	assert.NoError(t, err)
}

func TestLoad_With_Invalid_Field_It_Should_Fail(t *testing.T) {
	f, err := os.Open("./../../assets/.env")
	assert.NoError(t, err)

	sample := &struct {
		BOOL1 bool `env:"APP_NAME"`
	}{}
	err = decoder.Decoder{File: f}.Decode(sample)
	assert.Error(t, err)

	err = f.Close()
	assert.NoError(t, err)
}

func TestLoad_With_Invalid_Nested_Structure_Field_It_Should_Fail(t *testing.T) {
	f, err := os.Open("./../../assets/.env")
	assert.NoError(t, err)

	type Inner struct {
		Float float64 `env:"APP_NAME"`
	}

	outer := &struct {
		Inner Inner
	}{}
	outer.Inner = Inner{}

	err = decoder.Decoder{File: f}.Decode(outer)
	assert.Error(t, err)

	err = f.Close()
	assert.NoError(t, err)
}

func TestLoad_With_Invalid_Nested_Structure_Ptr_Field_It_Should_Fail(t *testing.T) {
	f, err := os.Open("./../../assets/.env")
	assert.NoError(t, err)

	type Inner struct {
		Float float64 `env:"APP_NAME"`
	}

	outer := &struct {
		Inner *Inner
	}{}
	outer.Inner = &Inner{}

	err = decoder.Decoder{File: f}.Decode(outer)
	assert.Error(t, err)

	err = f.Close()
	assert.NoError(t, err)
}
