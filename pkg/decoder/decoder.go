package decoder

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/golobby/cast"
	"os"
	"reflect"
	"regexp"
	"strings"
	"unsafe"
)

type Decoder struct {
	File *os.File
}

// Decode reads a dot env (.env) file and fills the given struct fields.
func (d Decoder) Decode(structure interface{}) error {
	kvs, err := d.read(d.File)
	if err != nil {
		return err
	}

	if err = d.feed(structure, kvs); err != nil {
		return err
	}

	return nil
}

// read scans a dot env (.env) file and extracts its key/value pairs.
func (d Decoder) read(file *os.File) (map[string]string, error) {
	scanner := bufio.NewScanner(file)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	kvs := d.parse(strings.Join(lines, "\n"))

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("dotenv: error when scanning file; err: %v", err)
	}

	return kvs, nil
}

func (d Decoder) parse(lines string) map[string]string {
	kvs := make(map[string]string)

	re := regexp.MustCompile("(?m)^\\s*(?:export\\s+)?([\\w.-]+)(?:\\s*=\\s*?|:\\s+?)(\\s*'(?:\\\\'|[^'])*'|\\s*\"(?:\\\\\"|[^\"])*\"|\\s*`(?:\\\\`|[^`])*`|[^#\\r\\n]+)?\\s*(?:#.*)?$")
	matches := re.FindAllStringSubmatch(lines, -1)

	for _, match := range matches {
		var key string
		var value string

		key = match[1]
		if len(match) > 1 {
			value = match[2]
		}

		value = strings.TrimSpace(value)

		if value == "" {
			kvs[key] = value
			continue
		}

		maybeQuote := string(value[0])

		if maybeQuote == "'" || maybeQuote == "`" || maybeQuote == `"` {
			pattern := fmt.Sprintf(`^%s([\s\S]*)%s[^%s]*$`, maybeQuote, maybeQuote, maybeQuote)
			res := regexp.MustCompile(pattern).FindStringSubmatch(value)
			if res != nil {
				value = res[1]
			}
		}

		kvs[key] = value
	}

	return kvs
}

// feed sets struct fields with the given key/value pairs.
func (d Decoder) feed(structure interface{}, kvs map[string]string) error {
	inputType := reflect.TypeOf(structure)
	if inputType != nil {
		if inputType.Kind() == reflect.Ptr {
			if inputType.Elem().Kind() == reflect.Struct {
				return d.feedStruct(reflect.ValueOf(structure).Elem(), kvs)
			}
		}
	}

	return errors.New("dotenv: invalid structure")
}

// feedStruct sets reflected struct fields with the given key/value pairs.
func (d Decoder) feedStruct(s reflect.Value, vars map[string]string) error {
	for i := 0; i < s.NumField(); i++ {
		if t, exist := s.Type().Field(i).Tag.Lookup("env"); exist {
			if val, exist := vars[t]; exist {
				v, err := cast.FromType(val, s.Type().Field(i).Type)
				if err != nil {
					return fmt.Errorf("dotenv: cannot set `%v` field; err: %v", s.Type().Field(i).Name, err)
				}

				ptr := reflect.NewAt(s.Field(i).Type(), unsafe.Pointer(s.Field(i).UnsafeAddr())).Elem()
				ptr.Set(reflect.ValueOf(v))
			}
		} else if s.Type().Field(i).Type.Kind() == reflect.Struct {
			if err := d.feedStruct(s.Field(i), vars); err != nil {
				return err
			}
		} else if s.Type().Field(i).Type.Kind() == reflect.Ptr {
			if s.Field(i).IsZero() == false && s.Field(i).Elem().Type().Kind() == reflect.Struct {
				if err := d.feedStruct(s.Field(i).Elem(), vars); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
