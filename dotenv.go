// Package dotenv is a lightweight package for loading dot env (.env) files into structs in Go projects.
package dotenv

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

// Load read a dot env (.env) file and fills the given struct.
func Load(file *os.File, structure interface{}) error {
	kvs, err := read(file)
	if err != nil {
		return err
	}

	if err = fill(structure, kvs); err != nil {
		return err
	}

	return nil
}

// read scans given dot env (.env) file and extract its key/value pairs
func read(file *os.File) (map[string]string, error) {
	kvs := map[string]string{}
	scanner := bufio.NewScanner(file)

	for i := 1; scanner.Scan(); i++ {
		if k, v, err := parse(scanner.Text()); err != nil {
			return nil, fmt.Errorf("dotenv: error in line %v; err: %v", i, err)
		} else if k != "" {
			kvs[k] = v
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("dotenv: error when scanning file; err: %v", err)
	}

	return kvs, nil
}

// parse extracts a key/value pair from the given dot env (.env) single line
func parse(line string) (string, string, error) {
	ln := strings.TrimSpace(line)
	kv := []string{"", ""}
	ci := 0
	qt := false

	for i := 0; i < len(ln); i++ {
		if string(ln[i]) == "#" {
			break
		}

		if string(ln[i]) == "=" && ci == 0 {
			ci = 1
			continue
		}

		if string(ln[i]) == " " {
			if kv[ci] == "" {
				if qt == false {
					continue
				}
			} else {
				if ci == 1 && qt == false {
					break
				}
			}
		}

		if string(ln[i]) == "\"" && ci == 1 {
			if kv[ci] == "" {
				qt = true
				continue
			} else if qt == true {
				break
			}
		}

		kv[ci] += string(ln[i])
	}

	if (ci == 0 && kv[0] != "") || (ci == 1 && kv[0] == "") {
		return "", "", fmt.Errorf("dotenv: invalid syntax")
	}

	return strings.TrimSpace(kv[0]), kv[1], nil
}

// fill sets a struct fields with the given map of key/value pairs
func fill(structure interface{}, kvs map[string]string) error {
	inputType := reflect.TypeOf(structure)
	if inputType != nil {
		if inputType.Kind() == reflect.Ptr {
			if inputType.Elem().Kind() == reflect.Struct {
				return fillStruct(reflect.ValueOf(structure).Elem(), kvs)
			}
		}
	}

	return errors.New("dotenv: invalid structure")
}

// fillStruct sets a reflected struct fields with the given map of key/value pairs
func fillStruct(s reflect.Value, vars map[string]string) error {
	for i := 0; i < s.NumField(); i++ {
		if t, exist := s.Type().Field(i).Tag.Lookup("env"); exist {
			if val, exist := vars[t]; exist {
				v, err := cast.FromString(val, s.Type().Field(i).Type.Name())
				if err != nil {
					return fmt.Errorf("dotenv: cannot set `%v` field; err: %v", s.Type().Field(i).Name, err)
				}

				ptr := reflect.NewAt(s.Field(i).Type(), unsafe.Pointer(s.Field(i).UnsafeAddr())).Elem()
				ptr.Set(reflect.ValueOf(v))
			}
		} else if s.Type().Field(i).Type.Kind() == reflect.Struct {
			if err := fillStruct(s.Field(i), vars); err != nil {
				return err
			}
		} else if s.Type().Field(i).Type.Kind() == reflect.Ptr {
			if s.Field(i).IsZero() == false && s.Field(i).Elem().Type().Kind() == reflect.Struct {
				if err := fillStruct(s.Field(i).Elem(), vars); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
