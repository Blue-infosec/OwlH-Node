package controllers

import (
    "owlhnode/models"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
    "owlhnode/validation"
)

type CollectorController struct {
    beego.Controller
}

// @Title PlayCollector
// @Description Play collector
// @Success 200 {object} models.Collector
// @Failure 403 body is empty
// @router /play [get]
func (n *CollectorController) PlayCollector() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "get")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{    
        err := models.PlayCollector()
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            logs.Error("COLLECTOR CREATE -> error: %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title StopCollector
// @Description Stop collector
// @Success 200 {object} models.Collector
// @Failure 403 body is empty
// @router /stop [get]
func (n *CollectorController) StopCollector() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "get")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{    
        err := models.StopCollector()
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            logs.Error("COLLECTOR CREATE -> error: %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title ShowCollector
// @Description Show collector
// @Success 200 {object} models.Collector
// @Failure 403 body is empty
// @router /show [get]
func (n *CollectorController) ShowCollector() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "get")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{    
        data, err := models.ShowCollector()
        n.Data["json"] = data
    
        if err != nil {
            logs.Error("COLLECTOR CREATE -> error: %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}