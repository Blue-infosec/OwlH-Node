package controllers

import (
    "owlhnode/models"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
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
    err := models.PlayCollector()
    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        logs.Error("COLLECTOR CREATE -> error: %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title StopCollector
// @Description Stop collector
// @Success 200 {object} models.Collector
// @Failure 403 body is empty
// @router /stop [get]
func (n *CollectorController) StopCollector() {
    err := models.StopCollector()
    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        logs.Error("COLLECTOR CREATE -> error: %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title ShowCollector
// @Description Show collector
// @Success 200 {object} models.Collector
// @Failure 403 body is empty
// @router /show [get]
func (n *CollectorController) ShowCollector() {
    data, err := models.ShowCollector()
    n.Data["json"] = data

    if err != nil {
        logs.Error("COLLECTOR CREATE -> error: %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}