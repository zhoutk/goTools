# goTools
create some tools use go lang.

# characteristic
- can export table, data, views, funcs.
- multi grountinue to export many database.
- filter \xfffd.
- solve the dependence of views.
- can recieve a cli arg, which must be in table, data, views & funcs, to export single content.

# configs.json
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

