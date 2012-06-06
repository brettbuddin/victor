package json

import (
    "log"
    "encoding/json"
    "errors"
)

type Json struct {
    Data interface{}
}

func Unmarshal(buf []byte) (*Json, error) {
    j   := new(Json)
    err := json.Unmarshal(buf, &j.Data)
    
    if err != nil {
        return nil, err
    }

    return j, nil
}

func (self *Json) Get(key string) *Json {
    m, err := self.Map()

    if err == nil {
        if val, exists := m[key]; exists {
            return &Json{val}
        }
    }

    return &Json{nil}
}

func (self *Json) Map() (map[string]interface{}, error) {
    if v, exists := (self.Data).(map[string]interface{}); exists {
        return v, nil
    }

    return nil, errors.New("could not assert value to map")
}

func (self *Json) Array() ([]interface{}, error) {
    if v, exists := (self.Data).([]interface{}); exists {
        return v, nil
    }

    return nil, errors.New("could not assert value to array")
}

func (self *Json) Int() (int, error) {
    if v, exists := (self.Data).(int); exists {
        return v, nil
    }

    return -1, errors.New("could not assert value to int")
}

func (self *Json) String() (string, error) {
    if v, exists := (self.Data).(string); exists {
        return v, nil
    }

    return "", errors.New("could not assert value to string")
}

func (self *Json) Bytes() ([]byte, error) {
    if v, exists := (self.Data).([]byte); exists {
        return v, nil
    }

    return nil, errors.New("could not assert value to bytes")
}

func (self *Json) MustString(args ...string) string {
    result, err := self.String()

    if err == nil {
        return result
    }

    if len(args) == 1 {
        result = args[0]
    } else {
        log.Panic("could not force json object into string") 
    }

    return result
}
