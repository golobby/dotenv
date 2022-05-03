package decoder

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/golobby/cast"
	"os"
	"reflect"
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

	lines := make([]string, 0)
	for i := 1; scanner.Scan(); i++ {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("dotenv: error when scanning file; err: %v", err)
	}

	kvs, err := d.parseLines(lines)
	if err != nil {
		return nil, err
	}

	return kvs, nil
}

// parseLines extracts a key/value pairs from the given dot env (.env).
func (d Decoder) parseLines(lines []string) (map[string]string, error) {
	kv := map[string]string{}
	escape := "\\"
	quote := ""
	quoteWas := false
	quoteOpen := false
	key := ""
	val := ""
	valueMode := false

	for l := 0; l < len(lines); l++ {
		ln := strings.TrimSpace(lines[l])

		prev := ""

		for i := 0; i < len(ln); i++ {

			char := string(ln[i])

			// comment symbol
			if char == "#" && (!valueMode || (valueMode && !quoteOpen)) {
				break
			}

			// enable value mode
			if char == "=" && !valueMode {
				valueMode = true
				continue
			}

			// skip space symbol in value mode
			if char == " " && valueMode {
				if !quoteOpen && val == "" {
					continue
				}
			}

			// quote symbols in value mode
			if (char == "\"" || char == "'") && valueMode {
				// if value is empty
				if val == "" {
					quoteOpen = true
					quote = char
					continue
				}

				// if close quote occurred
				if quoteOpen && quote == char {
					if prev != escape {
						quoteOpen = false
						quoteWas = true
						break
					}
				}
			}

			if !valueMode {
				key += char
			}

			if valueMode {
				if i == 0 {
					val += "\n"
				}
				if prev == escape && char == quote {
					val = strings.TrimSuffix(val, escape)
				}
				val += char
			}
			prev = char
		}

		// end of line
		if !quoteOpen {
			key = strings.TrimSpace(key)
			if !quoteWas {
				val = strings.TrimSpace(val)
			}

			if valueMode && key == "" {
				return nil, fmt.Errorf("dotenv: invalid syntax in line %d", l)
			}

			if len(key) > 0 {
				kv[key] = val
				key = ""
				val = ""
			}
			valueMode = false
			quoteWas = false
			quote = ""
		}
	}

	if quoteOpen {
		return nil, fmt.Errorf("dotenv: invalid syntax")
	}

	return kv, nil
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
