package initial

import (
	. "phagego/frameweb-v2/initial"
	// . "phagego/frameweb-v2/models"
)

func init() {
	InitLog()
	InitSql()
	InitBeeCache()
	InitFilter()
	InitMailConf()
	InitSysTemplateFunc()
	initTemplateFunc()

	InitLoginFilter()
}
