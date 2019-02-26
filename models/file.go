package models

import (
	"owlhnode/file"
	"github.com/astaxie/beego/logs"
)

func SendFile(filename string) (data map[string]string, err error) {
    logs.Info("SendFile into Node file")
	data,err = file.SendFile(filename)
    return data,err
}

func SaveFile(data map[string]string) (err error) {
    logs.Info("SaveFile into Node file")
	err = file.SaveFile(data)
    return err
}

func GetAllFiles() (data map[string]string, err error) {
    logs.Info("GetAllFiles into Node file")
	data,err = file.GetAllFiles()
    return data,err
}