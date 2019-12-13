package controllers

import (
    "owlhnode/models"
    "encoding/json"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
)

type GroupController struct {
    beego.Controller
}

// @Title SyncSuricataGroupValues
// @Description get Suricata group values
// @Success 200 {object} models.suricata
// @router /sync [put]
func (n *GroupController) SyncSuricataGroupValues() {
	var anode map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
	err := models.SyncSuricataGroupValues(anode)
    
    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        logs.Error("SyncSuricataGroupValues controller -> GET -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}