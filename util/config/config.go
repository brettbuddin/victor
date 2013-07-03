package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const ErrorsKey = "_errors"

type Object map[string]interface{}

// Objects

func (o Object) RequiredObject(key string) Object {
	return o.object(key, true)
}

func (o Object) OptionalObject(key string) Object {
	return o.object(key, false)
}

func (o Object) object(key string, required bool) Object {
	el, ok := o[key]
	if !ok {
		if required {
			o.pushError(fmt.Errorf("Missing required key %q (object)", key))
			return make(Object)
		}

		return make(Object)
	}

	m, ok := el.(map[string]interface{})
	if !ok {
		o.pushError(fmt.Errorf("Expected key %q to be an object, not %T", key, el))
		return make(Object)
	}

	return Object(m)
}

// Lists of Strings

func (o Object) RequiredStringList(key string) []string {
	return o.stringList(key, true)
}

func (o Object) OptionalStringList(key string) []string {
	return o.stringList(key, false)
}

func (o Object) stringList(key string, required bool) []string {
	el, ok := o[key]
	if !ok {
		if required {
			o.pushError(fmt.Errorf("Missing required key %q (list of strings)", key))
		}

		return nil
	}

	l, ok := el.([]interface{})
	if !ok {
		o.pushError(fmt.Errorf("Expected key %q to be a list of strings, not %T", key, el))
		return nil
	}

	list := make([]string, len(l))
	for i, el := range l {
		s, ok := el.(string)
		if !ok {
			o.pushError(fmt.Errorf("Expected key %q index %d to be a string, not %T", key, i, el))
			return nil
		}

		list[i] = s
	}

	return list
}

// Lists of Integers

func (o Object) RequiredIntList(key string) []int {
	return o.intList(key, true)
}

func (o Object) OptionalIntList(key string) []int {
	return o.intList(key, false)
}

func (o Object) intList(key string, required bool) []int {
	el, ok := o[key]
	if !ok {
		if required {
			o.pushError(fmt.Errorf("Missing required key %q (list of strings)", key))
		}

		return nil
	}

	l, ok := el.([]interface{})
	if !ok {
		o.pushError(fmt.Errorf("Expected key %q to be a list of strings, not %T", key, el))
		return nil
	}

	list := make([]int, len(l))
	for i, el := range l {
		v, ok := el.(int)
		if !ok {
			o.pushError(fmt.Errorf("Expected key %q index %d to be a string, not %T", key, i, el))
			return nil
		}

		list[i] = v
	}

	return list
}

// Strings

func (o Object) RequiredString(key string) string {
	return o.string(key, nil)
}

func (o Object) OptionalString(key string, def string) string {
	return o.string(key, &def)
}

func (o Object) string(key string, def *string) string {
	el, ok := o[key]
	if !ok {
		if def == nil {
			o.pushError(fmt.Errorf("Missing required key %q (string)", key))
			return ""
		}

		return *def
	}

	s, ok := el.(string)
	if !ok {
		o.pushError(fmt.Errorf("Expected key %q to be an string, not %T", key, el))
		return ""
	}

	return s
}

// Booleans

func (o Object) RequiredBool(key string) bool {
	return o.bool(key, nil)
}

func (o Object) OptionalBool(key string, def bool) bool {
	return o.bool(key, &def)
}

func (o Object) bool(key string, def *bool) bool {
	el, ok := o[key]
	if !ok {
		if def == nil {
			o.pushError(fmt.Errorf("Missing required key %q (bool)", key))
			return false
		}

		return *def
	}

	b, ok := el.(bool)
	if !ok {
		o.pushError(fmt.Errorf("Expected key %q to be an bool, not %T", key, el))
		return false
	}

	return b
}

// Integers

func (o Object) RequiredInt(key string) int {
	return o.int(key, nil)
}

func (o Object) OptionalInt(key string, def int) int {
	return o.int(key, &def)
}

func (o Object) int(key string, def *int) int {
	el, ok := o[key]
	if !ok {
		if def == nil {
			o.pushError(fmt.Errorf("Missing required key %q (int)", key))
			return 0
		}

		return *def
	}

	i, ok := el.(int)
	if !ok {
		o.pushError(fmt.Errorf("Expected key %q to be an int, not %T", key, el))
		return 0
	}

	return i
}

func (o Object) pushError(err error) {
	_, ok := o[ErrorsKey]
	if ok {
		o[ErrorsKey] = append(o[ErrorsKey].([]error), err)
		return
	}

	o[ErrorsKey] = []error{err}
}

func File(path string) (map[string]interface{}, error) {
	var (
		file *os.File
		err  error
	)

	if file, err = os.Open(path); err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := decode(bufio.NewReader(file))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func decode(reader io.Reader) (map[string]interface{}, error) {
	var payload map[string]interface{}

	decoder := json.NewDecoder(reader)

	if err := decoder.Decode(&payload); err != nil {
		return payload, err
	}

	return payload, nil
}
