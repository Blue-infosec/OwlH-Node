package stap

import (
    "github.com/astaxie/beego/logs"
    // "godoc.org/golang.org/x/crypto/ssh"
    "os"
    "os/exec"
    // "strings"
    // "regexp"
  	"owlhnode/utils"
  	// "owlhnode/database"
	  "io/ioutil"
      //"encoding/json"
      "time"
    //   "strconv"
    //   "errors"
      //"ssh.CleintConfig"
    //   "code.google.com/p/go.crypto/ssh"
    //   "sync"
    // "runtime"
	// "math/rand"
	// "golang.org/x/crypto/ssh"  
)

func Pcap_replay()() {
	loadStap := map[string]map[string]string{}
	loadStap["stap"] = map[string]string{}
    loadStap["stap"]["in_queue"] = ""
	loadStap["stap"]["out_queue"] = ""
	loadStap["stap"]["interface"] = ""
    loadStap = utils.GetConf(loadStap)
    inQueue := loadStap["stap"]["in_queue"]
	outQueue := loadStap["stap"]["out_queue"]
	stapInterface := loadStap["stap"]["interface"]
	// inQueue := "/usr/share/owlh/in_queue/"
	// outQueue := "/usr/share/owlh/out_queue/"
	// interface := enp0s3




	logs.Debug("Inside the PcapReplay, just before the loop")
	for{
		files, _ := ioutil.ReadDir(inQueue)
		if len(files) == 0 {
			logs.Error("Error Pcap_replay reading files: No files")
			time.Sleep(time.Second * 10)
		}
		x := 0
		// logs.Info(len(files))
		for _, f := range files{
			x += 1
			logs.Debug("Pcap_Replay-->"+f.Name())
			cmd := "tcpreplay -i "+stapInterface+" -t -l 1 "+inQueue+f.Name()
			logs.Debug(cmd)
			output, err := exec.Command("bash", "-c", cmd).Output()
			logs.Info(string(output))
			if err != nil{
				logs.Error("Error exec cmd command "+err.Error())
			}
			err = os.Rename(inQueue+f.Name(), outQueue+f.Name())
			if err != nil{
				logs.Error("Error moving file "+err.Error())
			}
		}
	}
}

// func loadQueuePath()(inQueue string, outQueue string){
// 	uuidParams, err := Sdb.Query("select server_param,server_value from servers where server_uniqueid = \""+uuid+"\";")
// 	defer uuidParams.Close()
// 	for uuidParams.Next(){
// 		if err = uuidParams.Scan(&param, &value); err!=nil {
// 			logs.Error("Error creating data Map: "+err.Error())
// 			return nil
// 		}
// 		stapServer[param]=value
// 	}
// }