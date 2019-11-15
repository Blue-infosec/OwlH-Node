package analyzer

import (
    "encoding/json"
    "os"
    "io/ioutil"
    "strings"
    "github.com/hpcloud/tail"
    "github.com/astaxie/beego/logs"
    "bufio"
    "time"
    "strconv"
    "owlhnode/utils"
    "owlhnode/database"
    "owlhnode/geolocation"
    "regexp"
)

type iocAlert struct {
    Data            Data        `json:"data"`
    Full_log        string      `json:"full_log"`
}

type Data struct {
    Dstport         string      `json:"dstport"`
    Srcport         string      `json:"srcport"`
    Dstip           string      `json:"dstip"`
    Srcip           string      `json:"srcip"`
    IoC             string      `json:"ioc"`
    IoCsource       string      `json:"iocsource"`
    Signature       Signature   `json:"alert"`
}

type Signature struct {
    Signature       string      `json:"signature"`
    Signature_id    string      `json:"signature_id"`
}


type Analyzer struct {
    Enable          bool        `json:"enable"`
    Srcfiles        []string    `json:"srcfiles"`
    Feedfiles       []Feedfile
}

type Feedfile struct {
    File            string      `json:"feedfile"`
    Workers         int         `json:"workers"`
}

var Dispatcher = make(map[string]chan string)
var Writer = make(map[string]chan string)

var config Analyzer

func readconf()(err error) {

    cfg := map[string]map[string]string{}
    cfg["analyzer"] = map[string]string{}
    cfg["analyzer"]["analyzerconf"] = ""
    cfg,err = utils.GetConf(cfg)
    analyzerCFG := cfg["analyzer"]["analyzerconf"]
    if err != nil {
        logs.Error("AlertLog Error getting data from main.conf: "+err.Error())
        return
    }

    confFile, err := os.Open(analyzerCFG)
    if err != nil {
        logs.Error("Error openning analyzer CFG: "+err.Error())
        return err
    }
    defer confFile.Close()
    byteValue, _ := ioutil.ReadAll(confFile)
    err = json.Unmarshal(byteValue, &config)
    if err != nil {
        logs.Error(err.Error())
        return err
    }
    return nil
}

func readLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

func Registerchannel(uuid string) {
    Dispatcher[uuid] = make(chan string)
}

func RegisterWriter(uuid string) {
    Writer[uuid] = make(chan string)
}

func Domystuff(IoCs []string, uuid string, wkrid int, iocsrc string) {
    for {
        line := <- Dispatcher[uuid] 
        for ioc := range IoCs {
            if strings.Contains(line, IoCs[ioc]) {
                logs.Info("Match -> "+ line +" IoC found -> " + IoCs[ioc]  + " wkrid -> " + strconv.Itoa(wkrid))
                //ioc
                IoCtoAlert(line, IoCs[ioc], iocsrc)
            }
        }
    }
}

func Mapper(uuid string, wkrid int) {
    logs.Info("Mapper -> " + uuid + " -> Started")
    for {
        line := <- Dispatcher[uuid] 
        StatClue := regexp.MustCompile("event_type\":\"stats\"")
        isStat := StatClue.FindStringSubmatch(line)
        if isStat != nil {
            continue
        }
        FlowClue := regexp.MustCompile("event_type\":\"flow\"")
        isFlow := FlowClue.FindStringSubmatch(line)
        if isFlow != nil {
            continue
        }
        line = strings.Replace(line, "id.orig_h", "srcip", -1)
        line = strings.Replace(line, "id.orig_p", "srcport", -1)
        line = strings.Replace(line, "id.resp_h", "dstip", -1)
        line = strings.Replace(line, "id.resp_p", "dstport", -1)
        line = strings.Replace(line, "src_ip", "srcip", -1)
        line = strings.Replace(line, "src_port", "srcport", -1)
        line = strings.Replace(line, "dest_ip", "dstip", -1)
        line = strings.Replace(line, "dest_port", "dstport", -1)
        re := regexp.MustCompile("dstip\":\"([^\"]+)\"")
        match := re.FindStringSubmatch(line)
        if match != nil {
            geoinfo_dst := geolocation.GetGeoInfo(match[1])
            if geoinfo_dst != nil {
                geodstjson, _ := json.Marshal(geoinfo_dst)
                if last := len(line) - 1; last >= 0 && line[last] == '}' && string(geodstjson) != "{}" {
                    line = line[:last]
                    line = line + ",\"geolocation_dst\":"+string(geodstjson)+"}"
                }
            } 
            
        }
        re = regexp.MustCompile("srcip\":\"([^\"]+)\"")
        match = re.FindStringSubmatch(line)
        if match != nil {
            geoinfo_src := geolocation.GetGeoInfo(match[1])
            if geoinfo_src != nil {
                geosrcjson, _ := json.Marshal(geoinfo_src)
                if last := len(line) - 1; last >= 0 && line[last] == '}' && string(geosrcjson) != "{}"{
                    line = line[:last]
                    line = line + ",\"geolocation_src\":"+string(geosrcjson)+"}"
                }
            } 
        }
        writeline(line)
    }
}

func Writerproc(uuid string, wkrid int) {
    var err error
    AlertLog := map[string]map[string]string{}
    AlertLog["node"] = map[string]string{}
    AlertLog["node"]["alertLog"] = ""
    AlertLog,err = utils.GetConf(AlertLog)
    outputfile := AlertLog["node"]["alertLog"]
    if err != nil {
        logs.Error("AlertLog Error getting data from main.conf: " + err.Error())
        return
    }
    ofile, err := os.OpenFile(outputfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        logs.Error("Analyzer Writer: can't open output file: " + outputfile + " -> " + err.Error())
        return
    }
    logs.Info("Mapper -> writer -> Started -> " + outputfile)
    _, err = ofile.WriteString("started 2\n")
    defer ofile.Close()
    for {
        line := <- Writer[uuid] 
        _, err = ofile.WriteString(line+"\n")
        if err != nil {
            logs.Error("Analyzer Writer: can't write line to file: " + outputfile + " -> " + err.Error())
        }
    }
}

func Startanalyzer(file string, wkr int) {
    newuuid := utils.Generate()
    logs.Info(newuuid + ": starting analyzer with feed: "+file + " with " + strconv.Itoa(wkr) + " workers")
    Registerchannel(newuuid)
    IoCs, _ := readLines(file)
    for x:=0; x < wkr; x++ {
        go Domystuff(IoCs, newuuid, x, file)
    }
}

func StartMapper(wkr int) {
    newuuid := utils.Generate()
    logs.Info(newuuid + ": starting Mapper with " + strconv.Itoa(wkr) + " workers")
    Registerchannel(newuuid)
    for x:=0; x < wkr; x++ {
        go Mapper(newuuid, x)
    }
}

func StartWriter(wkr int) {
    newuuid := utils.Generate()
    logs.Info(newuuid + ": starting Writer with " + strconv.Itoa(wkr) + " workers")
    RegisterWriter(newuuid)
    for x:=0; x < wkr; x++ {
        go Writerproc(newuuid, x)
    }
}

func Starttail(file string) {
    logs.Info("Starting tail of file: "+file)
    var seekv tail.SeekInfo
    seekv.Offset = 0
    seekv.Whence = os.SEEK_END
    for {
        t, _ := tail.TailFile(file, tail.Config{Follow: true, Location: &seekv})
        for line := range t.Lines {
            dispatch(line.Text)
        }
    }
}

func LoadAnalyzers() {
    logs.Info("loading analyzers")
    for file := range config.Feedfiles {
        go Startanalyzer(config.Feedfiles[file].File, config.Feedfiles[file].Workers)
    }
}

func LoadSources() {
    logs.Info("loading sources")
    for file := range config.Srcfiles {
        go Starttail(config.Srcfiles[file])
    }
}

func LoadMapper() {
    logs.Info("loading Mappers")
    go StartMapper(4)
}

func dispatch(line string) {
    for channel := range Dispatcher {
        Dispatcher[channel] <- line
    }
}

func writeline(line string) {
    for channel := range Writer {
        Writer[channel] <- line
    }
}

func IoCtoAlert(line, ioc, iocsrc string) {
    var err error
    AlertLog := map[string]map[string]string{}
    AlertLog["Node"] = map[string]string{}
    AlertLog["Node"]["AlertLog"] = ""
    AlertLog,err = utils.GetConf(AlertLog)
    AlertLogJson := AlertLog["Node"]["AlertLog"]
    if err != nil {
        logs.Error("AlertLog Error getting data from main.conf: "+err.Error())
        return
    }

    alert     := iocAlert{}
    data      := Data{}
    signature := Signature{}

    signature.Signature = "OwlH IoC found - "+ioc
    signature.Signature_id = "8000101"

    // data.Dstport = dstport
    // data.Dstip = dstip
    // data.Srcip = srcip
    // data.Srcport = srcport
    data.Signature = signature
    data.IoC = ioc
    data.IoCsource = iocsrc

    alert.Data = data
    alert.Full_log = line
    alertOutput, _ := json.Marshal(alert)

    err = utils.WriteNewDataOnFile(AlertLogJson, alertOutput)
    if err != nil {
        logs.Error("Error saving data IoCtoAlert: %s", err.Error())
    }

}

func InitAnalizer() {
    logs.Info("starting analyzer")
    analyzer,_ := PingAnalyzer()
    if analyzer["status"] == "Disabled"{
        return
    }
    readconf()
    StartWriter(1)
    LoadMapper()
    LoadAnalyzers()
    LoadSources()
    for {
        analyzer,_ = PingAnalyzer()
        if analyzer["status"] == "Disabled"{
            break
        }
        time.Sleep(time.Second * 3)
    }
}

func Init(){
    go InitAnalizer()
}

func PingAnalyzer()(data map[string]string ,err error) {
    wazuhFile := map[string]map[string]string{}
    wazuhFile["node"] = map[string]string{}
    wazuhFile["node"]["alertLog"] = ""
    wazuhFile,err = utils.GetConf(wazuhFile)    
    filePath := wazuhFile["node"]["alertLog"]
    if err != nil {logs.Error("PingAnalyzer Error getting data from main.conf")}

    analyzerData := make(map[string]string)
    analyzerData["status"] = "Disabled"


    analyzerStatus,err := ndb.GetStatusAnalyzer()
    if err != nil { logs.Error("Error getting Analyzer data: "+err.Error()); return analyzerData,err}

    analyzerData["status"] = analyzerStatus
    analyzerData["path"] = filePath

    fi, err := os.Stat(filePath);
    if err != nil { logs.Error("Can't access Analyzer ouput file data: "+err.Error()); return analyzerData,err}
    size := fi.Size()

    analyzerData["size"] = strconv.FormatInt(size, 10)

    return analyzerData, nil
}

func ChangeAnalyzerStatus(anode map[string]string) (err error) {
    logs.Emergency("ANALYZER STATUS - NEW STATUS - "+anode["status"])
    err = ndb.UpdateAnalyzer("analyzer", "status", anode["status"])
    if err != nil { logs.Error("Error updating Analyzer status: "+err.Error()); return err}
    
    return nil
}

func SyncAnalyzer(file map[string][]byte) (err error) {
    err = utils.WriteNewDataOnFile("conf/analyzer.json", file["data"])
    if err != nil { logs.Error("Analyzer/SyncAnalyzer Error updating Analyzer file: "+err.Error()); return err}
    
    return err
}