[![Sourcegraph](https://sourcegraph.com/github.com/JingDa-open-source-community/go-home/-/badge.svg)](https://sourcegraph.com/github.com/JingDa-open-source-community/go-home?badge)

#Querys train schedule use the origin command line tools


###output(train schedule)
![](http://p1.bqimg.com/567571/21b3d09e27e01ec1.gif)


###output(left tricket)
![](http://p1.bpimg.com/567571/bd4a89e17aa0bde0.gif)

###output(update)
![](http://i1.piimg.com/567571/ad64c6ff02bbca8b.gif)

###output(transfer query)
![Markdown](http://i1.piimg.com/1949/744687e65ea09b88.gif)

###dependence
1. ASCLL TABLE Writer it is for generate ascii table on termial
and use below command to install

###Create table
```
CREATE TABLE `station_lat_lgt` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `station` varchar(10) DEFAULT NULL,
  `latitude` varchar(30) DEFAULT NULL,
  `longitude` varchar(30) DEFAULT NULL,
  `insert_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;

CREATE TABLE `train_list` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `code` varchar(10) DEFAULT NULL,
  `train_no` varchar(30) DEFAULT NULL,
  `there` varchar(30) DEFAULT NULL,
  `home` varchar(30) DEFAULT NULL,
  `depart_date` date NOT NULL,
  `insert_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=265205 DEFAULT CHARSET=utf8;

```

###Databse
database name is 12306 go-home will auto create that.

###Config file

```
###this file is for config connect of mysql

type = "mysql"      // maybe you like use other database just change this type
host = "127.0.0.1"  // replace your host
port = 3306         // replace your port
user = root         // replace your user name
pass = digitalx168  //replace your password
databaseName = "12306" 
isCreate = false

```

###How to get started?
you should download project to path of yourself and then go build or go install, after that just run it!

```
go get github.com/JingDa-open-source-community/go-home
or you can use Makefile
cd yourpath/github.com/JingDa-open-source-community/go-home
make

go-home train d332 2017-02-25

go-home left 上海 永修 2017-03-03

// update mysql data
go-home update
```

###log

1. Add depart date parameter when query train schedule.
If you use a command like **go-home train k502 2017-01-22**, default value is tody,
so above command equal **go-home k502 2017-01-22**
In fact you should use a command like **go-home train k502 2017-01-27** to query someday's plan that you want to know. 
2. Support query left tricket go througth api of 12306 
3. Support user update local data by "update" option
4. Support go get
5. Support both ansi in window console and cygwinn&&mysys terminal 
6. Support transfer query

###Contact

Wechat: convertxy

QQ: 2698380951

Email: aliasliyu4@gamil.com

