package main

import (

    "github.com/astaxie/beego/logs"
    _ "owlhnode/routers"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/plugins/cors"
    "owlhnode/database"
    "owlhnode/stap"
    "owlhnode/analyzer"
    "owlhnode/plugin"
    "owlhnode/utils"
    "owlhnode/knownports"
    "owlhnode/geolocation"
    "owlhnode/monitor"
    "os"
	"bufio"
	"strings"
	"runtime"
)

func main() {
	//operative system values
	data:=OperativeSystemValues()
	for x := range data {
		if (x == "ID" || x == "ID_LIKE" || x == "VERSION_ID"){
			logs.Info(x +" -- "+data[x])
		}
	}
	
    var err error
    loadDataLogger := map[string]map[string]string{}
    loadDataLogger["logs"] = map[string]string{}
    loadDataLogger["logs"]["filename"] = ""
    loadDataLogger["logs"]["maxlines"] = ""
    loadDataLogger["logs"]["maxsize"] = ""
    loadDataLogger["logs"]["daily"] = ""
    loadDataLogger["logs"]["maxdays"] = ""
    loadDataLogger["logs"]["rotate"] = ""
    loadDataLogger["logs"]["level"] = ""
    loadDataLogger, err = utils.GetConf(loadDataLogger)    
    filename := loadDataLogger["logs"]["filename"]
    maxlines := loadDataLogger["logs"]["maxlines"]
    maxsize := loadDataLogger["logs"]["maxsize"]
    daily := loadDataLogger["logs"]["daily"]
    maxdays := loadDataLogger["logs"]["maxdays"]
    rotate := loadDataLogger["logs"]["rotate"]
    level := loadDataLogger["logs"]["level"]
    if err != nil {
        logs.Error("Main Error getting data from main.conf for load Logger data: "+err.Error())
    }
    logs.NewLogger(10000)
    logs.SetLogger(logs.AdapterFile,`{"filename":"`+filename+`", "maxlines":`+maxlines+` ,"maxsize":`+maxsize+`, "daily":`+daily+`, "maxdays":`+maxdays+`, "rotate":`+rotate+`, "level":`+level+`}`)

    ndb.SConn()
    ndb.PConn()
    ndb.NConn()

    // logs.Error("Version: 0.5.190415.0922")

    //Launch StapInit for 1st time for check status and go concurrency if status==true
    plugin.CheckServicesStatus()
    stap.StapInit()
    knownports.Init()
    analyzer.Init()
    geolocation.Init()
    monitor.Init()
    
    if beego.BConfig.RunMode == "dev" {
        beego.BConfig.WebConfig.DirectoryIndex = true
        beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
    }
    
    beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin"},
        ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
        AllowCredentials: true,
    }))

    beego.Run()
}

func OperativeSystemValues()(values map[string]string){
	if (runtime.GOOS == "linux"){
		logs.Info("============"+runtime.GOOS+"============")
		var OSmap = make(map[string]string)
		file, err := os.Open("/etc/os-release")
		if err != nil {logs.Error("No os-release file")}
		defer file.Close()
		
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if (scanner.Text() != ""){
				sidsSplit := strings.Split(scanner.Text(), "=")
				str := strings.Replace(sidsSplit[1], "\"", "", -1)
				OSmap[sidsSplit[0]] = str
			}			
		}
		return OSmap
	}else{
		return nil
	}
}