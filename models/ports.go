package models

import (
    "owlhnode/knownports"
   "owlhnode/changeControl"
    "github.com/astaxie/beego/logs")

func ShowPorts() (data map[string]map[string]string, err error) {
    data,err = knownports.ShowPorts()
    changecontrol.ChangeControlInsertData(err, "ShowPorts")    
    return data,err
}

func ChangeMode(anode map[string]string) (err error) {
    cc := anode
    logs.Info("============")
    logs.Info("PORTS - ChangeMode")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    err = knownports.ChangeMode(anode)
    changecontrol.ChangeControlInsertData(err, "ChangeMode")    
    return err
}

func ChangeStatus(anode map[string]string) (err error) {
    cc := anode
    logs.Info("============")
    logs.Info("PORTS - ChangeStatus")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    err = knownports.ChangeStatus(anode)
    changecontrol.ChangeControlInsertData(err, "ChangeStatus")    
    return err
}

func DeletePorts(anode map[string]string) (err error) {
    cc := anode
    logs.Info("============")
    logs.Info("PORTS - DeletePorts")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    err = knownports.DeletePorts(anode)
    changecontrol.ChangeControlInsertData(err, "DeletePorts")    
    return err
}

func DeleteAllPorts(anode map[string]string) (err error) {
    cc := anode
    logs.Info("============")
    logs.Info("PORTS - DeleteAllPorts")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    err = knownports.DeleteAllPorts()
    changecontrol.ChangeControlInsertData(err, "DeleteAllPorts")    
    return err
}

func PingPorts() (data map[string]map[string]string ,err error) {
    data, err = knownports.PingPorts()
    changecontrol.ChangeControlInsertData(err, "PingPorts")    
    return data, err
}