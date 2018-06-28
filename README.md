# goTools
> create some tools use go lang.
1. mysql exporter
2. ip spider

## 1. mysql exporter
a tool can export mysql's table,data,views,function & stored procedure together or independent.

### characteristic
- can export table, data, views, funcs.
- multi grountinue to export many database together.
- filter \xfffd.
- solve the dependence of views.
- can recieve a cli arg, which must be in table, data, views & funcs, to export single content.

### configs.json
you must create configs.json as:

```
{
    "db_name1": {
        "db_host": "192.168.1.8",
        "db_port": 3306,
        "db_user": "root",
        "db_pass": "123456",
        "db_name": "name1",
        "db_charset": "utf8mb4",
        "file_alias": "file name1"
    },
    "db_name2": {
        "db_host": "localhost",
        "db_port": 3306,
        "db_user": "root",
        "db_pass": "123456",
        "db_name": "name2",
        "db_charset": "utf8mb4"
    },
    "database_dialect": "mysql",
    "workDir": "/home/zhoutk/gocodes/goTools/"
}
```
### Instructions
```
git clone https://github.com/zhoutk/goTools
cd goTools
go get
go run main.go

go buid main.go
./main                  #export all things of database
./main table            #export tables
./main data             #export tables & data
./main views            #export views
./main funcs            #export funcs & stored procedures
```

## 2. ip spider
a tool can spider ip address info from appointed web page.

### characteristic
- multi grountinue to spider web data.
- write mysql batch.
- update mysql batch.

### sql scripts
you can create table use it:

```
CREATE TABLE `ip_addr_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '索引，自动增长',
  `ip_addr_begin` varchar(32) NOT NULL DEFAULT '' COMMENT 'ip地址段开始',
  `ip_addr_end` varchar(32) DEFAULT '' COMMENT 'ip地址段结束',
  `province` varchar(32) DEFAULT '' COMMENT '所属省',
  `ip_comp` varchar(32) DEFAULT '' COMMENT '运营商',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ip_addr` (`ip_addr_begin`,`ip_addr_end`)
) ENGINE=InnoDB AUTO_INCREMENT=7268 DEFAULT CHARSET=utf8 COMMENT='表';
```
### Instructions
```
git clone https://github.com/zhoutk/goTools
cd goTools
go get
go run ip.go

go buid ip.go
./ip 
```

