package controllers

import (
    "owlhnode/models"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
    "encoding/json"
)

type WazuhController struct {
    beego.Controller
}

// @Title GetWazuh
// @Description get Wazuh status
// @Success 200 {object} models.wazuh
// @router / [get]
func (m *WazuhController) Get() {
    logs.Info ("Wazuh controller -> GET")
    mstatus, err := models.GetWazuh()
    m.Data["json"] = mstatus
    if err != nil {
        logs.Info("GetWazuh OUT -- ERROR : %s", err.Error())
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    m.ServeJSON()
}

// @Title RunWazuh
// @Description Run wazuh system
// @Success 200 {object} models.wazuh
// @Failure 403 body is empty
// @router /RunWazuh [put]
func (n *WazuhController) RunWazuh() {
    logs.Info("RunWazuh -> In")
    var anode map[string]string
    anode["action"] = "PUT"
    anode["controller"] = "WAZUH"
    anode["router"] = "@router /RunWazuh [put]"
    logs.Info("============")
    logs.Info("WAZUH - RunWazuh")
    for key :=range anode {
        logs.Info(key +" -> "+anode[key])
    }
    data,err := models.RunWazuh()
    n.Data["json"] = data
    if err != nil {
        logs.Info("RunWazuh OUT -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    logs.Info("RunWazuh -> OUT -> %s", n.Data["json"])
    n.ServeJSON()
}

// @Title StopWazuh
// @Description Run wazuh system
// @Success 200 {object} models.wazuh
// @Failure 403 body is empty
// @router /StopWazuh [put]
func (n *WazuhController) StopWazuh() {
    logs.Info("StopWazuh -> In")
    var anode map[string]string
    anode["action"] = "PUT"
    anode["controller"] = "WAZUH"
    anode["router"] = "@router /StopWazuh [put]"
    logs.Info("============")
    logs.Info("WAZUH - StopWazuh")
    for key :=range anode {
        logs.Info(key +" -> "+anode[key])
    }
    data,err := models.StopWazuh()
    n.Data["json"] = data
    if err != nil {
        logs.Info("StopWazuh OUT -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    logs.Info("StopWazuh -> OUT -> %s", n.Data["json"])
    n.ServeJSON()
}

// @Title PingWazuhFiles
// @Description get Wazuh status
// @Success 200 {object} models.wazuh
// @router /pingWazuhFiles [get]
func (m *WazuhController) PingWazuhFiles() {
    files, err := models.PingWazuhFiles()
    m.Data["json"] = files
    if err != nil {
        m.Data["json"] = map[int]map[string]string{0:{"ack": "false", "error": err.Error()}}
    }
    m.ServeJSON()
}

// @Title DeleteWazuhFile
// @Description Run wazuh system
// @Success 200 {object} models.wazuh
// @Failure 403 body is empty
// @router /deleteWazuhFile [delete]
func (n *WazuhController) DeleteWazuhFile() {
    var anode map[string]interface{}
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    anode["action"] = "PUT"
    anode["controller"] = "SURICATA"
    anode["router"] = "@router /StopSuricata [put]"
    err := models.DeleteWazuhFile(anode)
    n.Data["json"] = map[string]string{"ack": "true"}
    
    if err != nil {n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}}
    n.ServeJSON()
}

// @Title AddWazuhFile
// @Description Run wazuh system
// @Success 200 {object} models.wazuh
// @Failure 403 body is empty
// @router /addWazuhFile [put]
func (n *WazuhController) AddWazuhFile() {
    var anode map[string]interface{}
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    anode["action"] = "PUT"
    anode["controller"] = "SURICATA"
    anode["router"] = "@router /StopSuricata [put]"
    err := models.AddWazuhFile(anode)
    n.Data["json"] = map[string]string{"ack": "true"}
    
    if err != nil {n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}}
    n.ServeJSON()
}

// @Title LoadFileLastLines
// @Description Run wazuh system
// @Success 200 {object} models.wazuh
// @Failure 403 body is empty
// @router /loadFileLastLines [put]
func (n *WazuhController) LoadFileLastLines() {
    var anode map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    anode["action"] = "PUT"
    anode["controller"] = "SURICATA"
    anode["router"] = "@router /StopSuricata [put]"
    data,err := models.LoadFileLastLines(anode)
    n.Data["json"] = data
    
    if err != nil {n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}}
    n.ServeJSON()
}

// @Title SaveFileContentWazuh
// @Description Run wazuh system
// @Success 200 {object} models.wazuh
// @Failure 403 body is empty
// @router /saveFileContentWazuh [put]
func (n *WazuhController) SaveFileContentWazuh() {
    var anode map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    anode["action"] = "PUT"
    anode["controller"] = "SURICATA"
    anode["router"] = "@router /StopSuricata [put]"

    err := models.SaveFileContentWazuh(anode)
    n.Data["json"] = map[string]string{"ack": "true"}
    
    if err != nil {n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}}
    n.ServeJSON()
}