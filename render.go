
package utils

import (
	"bytes"
	"github.com/tx991020/utils/logs"
	"html/template"
)

func Render(path string, config map[string]interface{}) ([]byte, error) {
	t, err := template.ParseFiles(path)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	bufer := bytes.NewBuffer(nil)
	err = t.Execute(bufer, config)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	return bufer.Bytes(), nil

}