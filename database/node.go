package ndb

import (
    "github.com/astaxie/beego/logs"
    "database/sql"
	"os"
	"owlhnode/utils"
	_ "github.com/mattn/go-sqlite3"
)

var (
    Nodedb *sql.DB
)

func NConn() {
    var err error
	loadDataSQL := map[string]map[string]string{}
	loadDataSQL["nodeConn"] = map[string]string{}
	loadDataSQL["nodeConn"]["path"] = ""
	loadDataSQL["nodeConn"]["cmd"] = "" 
	loadDataSQL, err = utils.GetConf(loadDataSQL)    
    path := loadDataSQL["nodeConn"]["path"]
	cmd := loadDataSQL["nodeConn"]["cmd"]
	if err != nil {
		logs.Error("NConn Error getting data from main.conf")
	}
	_, err = os.Stat(path) 
	if err != nil {
		panic("Fail opening servers.db from path: "+path+"  --  "+err.Error())
	}	
	Nodedb, err = sql.Open(cmd,path)
    if err != nil {
        logs.Error("Nodedb/stap -- servers.db Open Failed: "+err.Error())
    }else {
		logs.Info("Nodedb/stap -- servers.db -> sql.Open, servers.db Ready") 
	}
}

func LoadDataflowValues()(path map[string]map[string]string, err error){
	var pingData = map[string]map[string]string{}
    var uniqid string
    var param string
	var value string

	sql := "select flow_uniqueid, flow_param, flow_value from dataflow;";
	
	rows, err := Nodedb.Query(sql)
	if err != nil {
		logs.Error("LoadDataflowValues Nodedb.Query Error : %s", err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&uniqid, &param, &value); err != nil {
            logs.Error("LoadDataflowValues -- Query return error: %s", err.Error())
            return nil, err
		}
        if pingData[uniqid] == nil { pingData[uniqid] = map[string]string{}}
        pingData[uniqid][param]=value
	} 
	return pingData,nil
}

func ChangeDataflowValues(uuid string, param string, value string) (err error) {
	updateDataflowNode, err := Nodedb.Prepare("update dataflow set flow_value = ? where flow_uniqueid = ? and flow_param = ?;")
	if (err != nil){
		logs.Error("ChangeDataflowValues UPDATE prepare error: "+err.Error())
		return err
	}
	_, err = updateDataflowNode.Exec(&value, &uuid, &param)
	defer updateDataflowNode.Close()
	if (err != nil){
		logs.Error("ChangeDataflowValues UPDATE error: "+err.Error())
		return err
	}
	return nil
}

func ChangeNodeconfigValues(uuid string, param string, value string)(err error){
	updateNodeconfig, err := Nodedb.Prepare("update nodeconfig set config_value = ? where config_uniqueid = ? and config_param = ?;")
	if (err != nil){
		logs.Error("Change Nodeconfig Values UPDATE prepare error: "+err.Error())
		return err
	}
	_, err = updateNodeconfig.Exec(&value, &uuid, &param)
	defer updateNodeconfig.Close()
	if (err != nil){
		logs.Error("Change Nodeconfig Values UPDATE error: "+err.Error())
		return err
	}
	return nil
}

func LoadNodeconfigValues()(path map[string]map[string]string, err error){
	var configValues = map[string]map[string]string{}
    var uniqid string
    var param string
	var value string

	sql := "select config_uniqueid, config_param, config_value from nodeconfig;";
	
	rows, err := Nodedb.Query(sql)
	if err != nil {
		logs.Error("LoadNodeconfigValues Nodedb.Query Error : %s", err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&uniqid, &param, &value); err != nil {
            logs.Error("LoadNodeconfigValues -- Query return error: %s", err.Error())
            return nil, err
		}
        if configValues[uniqid] == nil { configValues[uniqid] = map[string]string{}}
        configValues[uniqid][param]=value
	} 
	return configValues,nil
}

func GetNodeconfigValue(uuid string, param string)(val string, err error){
	var value string

	sql := "select config_value from nodeconfig where config_param=\""+param+"\" and config_uniqueid=\""+uuid+"\";";
	rows, err := Nodedb.Query(sql)
	if err != nil { logs.Error("GetNodeconfigValue Nodedb.Query Error : %s", err.Error()); return "", err}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&value); err != nil { logs.Error("GetNodeconfigValue -- Query return error: %s", err.Error()); return "", err}
	} 
	return value,nil
}