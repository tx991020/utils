package utils

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
	"time"

	"strconv"

	"github.com/fatih/structs"
)

type JSONTime struct {
	T time.Time
}

func (t *JSONTime) MarshalJSON() ([]byte, error) {
	stamp := strconv.FormatInt(int64(t.T.UnixNano()/1e6), 10)
	return []byte(stamp), nil
}

func (t *JSONTime) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return nil
	}

	if b[0] == b[len(b)-1] && b[0] == '"' {
		i, err := strconv.ParseInt(string(b[1:len(b)-1]), 10, 64)
		t.T = time.Unix(i/1e3, (i%1e3)*1e6)
		return err
	} else {
		i, err := strconv.ParseInt(string(b), 10, 64)
		t.T = time.Unix(i/1e3, (i%1e3)*1e6)
		return err
	}
}

func (t *JSONTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	stamp := fmt.Sprintf("\"%d\"", t.T.Unix())
	e.EncodeToken(start)
	e.EncodeToken(stamp)
	e.EncodeToken(xml.EndElement{start.Name})
	return nil
}

// sql.Scanner
func (t *JSONTime) Scan(value interface{}) error {
	t.T = value.(time.Time)
	return nil
}

// sql.driver.Valuer, MUST BE (t JSONTime), NOT (t *JSONTime)
func (t *JSONTime) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	} else {
		return t.T, nil
	}
}

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

func FillStruct(m map[string]interface{}, o interface{}) error {
	for k, v := range m {
		err := SetField(o, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// copy a to b, by the same golang field name
func Wrap(a interface{}, b interface{}) error {
	m := structs.Map(a)
	return FillStruct(m, b)
}

func MustJson(o interface{},  escape bool) []byte {
	b, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	if !escape {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b
}
