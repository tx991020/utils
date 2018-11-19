package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"

	"net/http"

	"errors"
	"strings"

	"gitee.com/johng/gf/g/os/gtime"
	"gitee.com/johng/gf/g/util/gregex"
	"github.com/GiterLab/aliyun-sms-go-sdk/dysms"
	"github.com/astaxie/beego/logs"

	"encoding/json"
	"github.com/golang/glog"
	"github.com/pborman/uuid"
	"github.com/tealeg/xlsx"
)

var (
	accessKeyId     = ""
	accessKeySecret = ""
	signName        = ""
	templateParam   = make(map[string]string)
	failList        = []string{}
)


func Upload(c *gin.Context) {

	templatecode := c.PostForm("templatecode")
	if templatecode == "" {
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}
	fmt.Println(file.Filename)
	if err := c.SaveUploadedFile(file, file.Filename); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	if err := OpenFile(file.Filename, templatecode); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("csv文件格式不对 err: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"errcode": 0, "fails": failList})

}
func OpenFile(name, templatecode string) error {
	if !strings.HasSuffix(name, ".xlsx") {
		return errors.New("文件名需要已.xlsx结尾")
	}
	xlFile, err := xlsx.OpenFile(name)
	if err != nil {
		logs.Error(err.Error())
		return err
	}

	for _, sheet := range xlFile.Sheets {

		for _, row := range sheet.Rows {
			phone := ""
			code := ""
			for _, cell := range row.Cells {
				text := cell.String()

				if gregex.IsMatchString(`\d{11}`, text) {
					phone = text
				}
				if gregex.IsMatchString(`[A-Za-z0-9]{12}`, text) {
					templateParam["code"] = text
					c, _ := json.Marshal(templateParam)
					code = string(c)
				}
			}
			fmt.Println(phone, code)
			//发
			//Send(phone, code, templatecode)

			//查
			//Check(phone)

		}
	}
	return nil
}

func Send(phoneNumbers, templateParam, templateCode string) {
	if phoneNumbers == "" {
		return
	}
	dysms.HTTPDebugEnable = true
	dysms.SetACLClient(accessKeyId, accessKeySecret)

	//短信发送
	respSendSms, err := dysms.SendSms(uuid.New(), phoneNumbers, signName, templateCode, templateParam).DoActionWithException()
	if err != nil {
		fmt.Println("send sms failed", err, respSendSms.Error())
		return

	}
	fmt.Println("send sms succeed", respSendSms.String())

}

func Check(phoneNum string) {
	if phoneNum == "" {
		return
	}
	// 查询短信
	dysms.HTTPDebugEnable = true
	dysms.SetACLClient(accessKeyId, accessKeySecret)
	date := gtime.Date()
	date1 := strings.Replace(date, "-", "", -1)

	respQuerySendDetails, err := dysms.QuerySendDetails("", phoneNum, "10", "1", date1).DoActionWithException()
	if err != nil {
		fmt.Println("query sms failed", err, respQuerySendDetails.Error())
		return

	}
	fmt.Println(*respQuerySendDetails.TotalCount)
	if *respQuerySendDetails.TotalCount < 1 {
		fmt.Println("发送失败%s", phoneNum)

		glog.Error(phoneNum)

	}

}
