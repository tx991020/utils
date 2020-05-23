package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"

	"github.com/tx991020/utils/logs"
)

func Render(path string, config interface{}) ([]byte, error) {
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

func RenderAppend(old, new, data string) error {
	readAll, err := ioutil.ReadFile(old)
	if err != nil {
		return err
	}
	replaceString, err := Replace(`\/\/{{[^}]+}}`, []byte(string(fmt.Sprintf(data+
		"\n    //{{.data}}`"))), readAll)
	if err != nil {
		return err
	}
	if err := PutContents(new, string(replaceString)); err != nil {
		return err
	}
	return nil
}
