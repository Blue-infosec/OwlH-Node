// @APIVersion 0.1.0
// @Title OwlH Node API
// @Description OwlH node API
// @Contact support@owlh.net
// @TermsOfServiceUrl http://www.owlh.net
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"owlhnode/controllers"
	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/node",
		beego.NSNamespace("/suricata",
			beego.NSInclude(
				&controllers.SuricataController{},
			),
		),
		beego.NSNamespace("/zeek",
			beego.NSInclude(
				&controllers.ZeekController{},
			),
		),
		beego.NSNamespace("/wazuh",
			beego.NSInclude(
				&controllers.WazuhController{},
			),
		),
		beego.NSNamespace("/file",
			beego.NSInclude(
				&controllers.FileController{},
			),
		),
		beego.NSNamespace("/stap",
			beego.NSInclude(
				&controllers.StapController{},
			),
		),
		beego.NSNamespace("/ping",
			beego.NSInclude(
				&controllers.PingController{},
			),
		),
		beego.NSNamespace("/collector",
			beego.NSInclude(
				&controllers.CollectorController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
