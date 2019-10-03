package wazuh

import (
    "github.com/astaxie/beego/logs"
    "os"
    "os/exec"
    "regexp"
    "strings"
    "errors"
    "bufio"
    // "io/ioutil"
    "encoding/json"
    // "bytes"
    "owlhnode/utils"
)

func WazuhPath() (exists bool) {
	var err error
    //Retrieve path for wazuh.
	loadDataWazuhPath := map[string]map[string]string{}
	loadDataWazuhPath["loadDataWazuhPath"] = map[string]string{}
	loadDataWazuhPath["loadDataWazuhPath"]["path"] = ""
    loadDataWazuhPath,err = utils.GetConf(loadDataWazuhPath)    
    path := loadDataWazuhPath["loadDataWazuhPath"]["path"]
	if err != nil {
		logs.Error("WazuhPath Error getting data from main.conf")
	}
	
    if _, err := os.Stat(path); os.IsNotExist(err) {
        logs.Error("Wazuh is not installed, at least at /var/ossec folder does not exist")
        return false
    }
    return true
}

func WazuhBin() (exists bool) {
	var err error
    //Retrieve bin for wazuh.
	loadDataWazuhBin := map[string]map[string]string{}
	loadDataWazuhBin["loadDataWazuhBin"] = map[string]string{}
    loadDataWazuhBin["loadDataWazuhBin"]["bin"] = ""
    loadDataWazuhBin,err = utils.GetConf(loadDataWazuhBin)    
	bin := loadDataWazuhBin["loadDataWazuhBin"]["bin"]
	if err != nil {
		logs.Error("WazuhBin Error getting data from main.conf")
	}
    if _, err := os.Stat(bin); os.IsNotExist(err) {
        logs.Error("Wazuh bin does not exist")
        return false
    }
    logs.Info("Wazuh bin exist")
    return true
}

func WazuhRunning() (running bool) {
	var err error
    //Retrieve running for wazuh.
	loadDataWazuhRunning := map[string]map[string]string{}
	loadDataWazuhRunning["loadDataWazuhRunning"] = map[string]string{}
    loadDataWazuhRunning["loadDataWazuhRunning"]["cmd"] = ""
    loadDataWazuhRunning["loadDataWazuhRunning"]["param"] = ""
    loadDataWazuhRunning["loadDataWazuhRunning"]["command"] = ""
    loadDataWazuhRunning,err = utils.GetConf(loadDataWazuhRunning)    
    cmd := loadDataWazuhRunning["loadDataWazuhRunning"]["cmd"]
    param := loadDataWazuhRunning["loadDataWazuhRunning"]["param"]
    command := loadDataWazuhRunning["loadDataWazuhRunning"]["command"]
	if err != nil {
		logs.Error("WazuhRunning Error getting data from main.conf")
	}
    out, err := exec.Command(command, param, cmd).Output()
    if err == nil {
        if strings.Contains(string(out), "is running") {
            logs.Info("Wazuh is running ->"+string(out))
            return true
        }
    }
    logs.Error("Wazuh is NOT running -> " + string(out))
    return false
}

func Installed() (isIt map[string]bool, err error){
    wazuh := make(map[string]bool)
    wazuh["path"] = WazuhPath()
    wazuh["bin"] = WazuhBin()
    wazuh["running"] = WazuhRunning()

    if wazuh["path"] || wazuh["bin"] || wazuh["running"]  {
        logs.Info("Wazuh installed and running")
        return wazuh, nil
    } else {
        logs.Error("Wazuh isn't present or not running")
        return wazuh, errors.New("Wazuh isn't present or not running")
    }
}

//Run wazuh
func RunWazuh()(data string, err error){

    // //Retrieve path for wazuh.
    StartWazuh := map[string]map[string]string{}
    StartWazuh["wazuhStart"] = map[string]string{}
    StartWazuh["wazuhStart"]["start"] = ""
    StartWazuh["wazuhStart"]["param"] = ""
    StartWazuh["wazuhStart"]["command"] = ""
    StartWazuh,err = utils.GetConf(StartWazuh)    
    cmd := StartWazuh["wazuhStart"]["start"]
    param := StartWazuh["wazuhStart"]["param"]
    command := StartWazuh["wazuhStart"]["command"]
	if err != nil {
		logs.Error("RunWazuh Error getting data from main.conf")
	}

    _,err = exec.Command(command, param, cmd).Output()
    if err != nil {
        logs.Error("Error launching wazuh: "+err.Error())
        return "",err
    }
    return "wazuh system is on",nil
}

//Stop wazuh
func StopWazuh()(data string, err error){

    // //Retrieve path for wazuh.
    StopWazuh := map[string]map[string]string{}
	StopWazuh["wazuhStop"] = map[string]string{}
    StopWazuh["wazuhStop"]["stop"] = ""
    StopWazuh["wazuhStop"]["param"] = ""
    StopWazuh["wazuhStop"]["command"] = ""
    StopWazuh,err = utils.GetConf(StopWazuh)    
    cmd := StopWazuh["wazuhStop"]["stop"]
    param := StopWazuh["wazuhStop"]["param"]
    command := StopWazuh["wazuhStop"]["command"]
    if err != nil {logs.Error("RunWazuh Error getting data from main.conf")}
    
    _,err = exec.Command(command, param, cmd).Output()
    if err != nil {logs.Error("Error stopping Wazuh: "+err.Error()); return "",err}

    return "Wazuh stopped ",nil
}

func PingWazuhFiles() (files map[string]string, err error) {    
    file, err := os.Open("/var/ossec/etc/ossec.conf")
    if err != nil {logs.Error(err)}
    defer file.Close()

    scanner := bufio.NewScanner(file)
    isInit := false
    isEnd := false
    filesPath := make(map[string]string)
    for scanner.Scan() {
        var init = regexp.MustCompile(`<!-- OWLH INIT -->`)
        var end = regexp.MustCompile(`<!-- OWLH END -->`)
        owlhInit := init.FindStringSubmatch(scanner.Text())
        owlhEnd := end.FindStringSubmatch(scanner.Text())
        if owlhInit != nil{ isInit = true; continue}
        if owlhEnd != nil{ isEnd = true}
        if isInit && !isEnd {
            var locationPath = regexp.MustCompile(`<location>([^"]+)<\/location>`)
            locationFound := locationPath.FindStringSubmatch(scanner.Text())
            if locationFound != nil {
                filesPath[locationFound[1]] = locationFound[1]
            }
        }
    }
    if err := scanner.Err(); err != nil {logs.Error(err)}

    return filesPath ,err
}

type AllFiles struct {
    UUID string `json:"uuid"`
    Paths []string `json:"paths"`
}


func DeleteWazuhFile(anode map[string]interface{})(err error) {
    receivedWazuhFiles := AllFiles{}
    byteData, _ := json.Marshal(anode)
    json.Unmarshal(byteData, &receivedWazuhFiles)

    file, err := os.Open("/var/ossec/etc/ossec.conf")
    if err != nil {logs.Error("Error DeleteWazuhFile readding file: "+err.Error()); return err}
    defer file.Close()

    // var buf bytes.Buffer
    isInit := false
    isEnd := false
    isSecondEnd := false
    scanner := bufio.NewScanner(file)
    var h int 
    h = 0
    fileContent := make(map[int]string)
    for scanner.Scan() {
        // var init = regexp.MustCompile(`<!-- OWLH INIT -->\n<ossec_config>`)
        // var end = regexp.MustCompile(`</ossec_config>\n<!-- OWLH END -->`)
        var init = regexp.MustCompile(`<!-- OWLH INIT -->`)
        var end = regexp.MustCompile(`<!-- OWLH END -->`)
        owlhInit := init.FindStringSubmatch(scanner.Text())
        owlhEnd := end.FindStringSubmatch(scanner.Text())
           
        if owlhInit != nil{
            isInit = true
            fileContent[h] = "<!-- OWLH INIT -->"
            h++
            fileContent[h] = "<ossec_config>"
            h++ 
        }
        if owlhEnd != nil{isEnd = true}
        if isInit && !isEnd {
            
            for x := range receivedWazuhFiles.Paths{
                logs.Emergency(receivedWazuhFiles.Paths[x])
                fileContent[h] = "\t<localfile>"; h++
                fileContent[h] = "\t\t<log_format>syslog</log_format>"; h++
                fileContent[h] = "\t\t<location>"+receivedWazuhFiles.Paths[x]+"</location>  "; h++
                fileContent[h] = "\t</localfile>"; h++
            }
            isEnd = true
            fileContent[h] = "<ossec_config>"
            h++
            fileContent[h] = "<!-- OWLH END -->"
        }else if isInit && isEnd && !isSecondEnd{
            var secondEnd = regexp.MustCompile(`<!-- OWLH END -->`)
            secondOwlhEnd := secondEnd.FindStringSubmatch(scanner.Text())
            if secondOwlhEnd == nil{
                continue
            }else if secondOwlhEnd != nil && !isSecondEnd {
                isSecondEnd = true
                // h++
                continue
            }
            
        }else{
            fileContent[h] = scanner.Text()
        }
        h++
    }
    if err := scanner.Err(); err != nil {logs.Error("DeleteWazuhFile. Scanner file error: "+err.Error()); return err}

    // saveIntoFile, err := os.Open("/var/ossec/etc/ossec.conf")
    saveIntoFile, err := os.OpenFile("/var/ossec/etc/ossec.conf", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
    if err != nil {logs.Error("Error DeleteWazuhFile readding file: "+err.Error()); return err}
    defer saveIntoFile.Close()
    saveIntoFile.Truncate(0)
    saveIntoFile.Seek(0,0)
    for x:=0 ; x < h ; x++{
        // _, err := file.WriteString(fileContent[x])
        _, err = saveIntoFile.WriteAt([]byte(fileContent[x]+"\n"), 0) // Write at 0 beginning
        if err != nil {logs.Error("DeleteWazuhFile failed writing to file: %s", err); return err}
    }

    return err
}   