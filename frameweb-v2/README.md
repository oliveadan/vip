# 基于beego的web框架

项目依赖：
go get github.com/astaxie/beego
go get github.com/go-sql-driver/mysql
go get github.com/skip2/go-qrcode
go get github.com/garyburd/redigo/redis
go get gopkg.in/gomail.v2

ALTER TABLE `ph_site_config`
CHANGE COLUMN `value` `value` VARCHAR(1024) NOT NULL DEFAULT '' ;
