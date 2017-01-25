#Querys train schedule use the origin command line tools


###output(train schedule)
![](http://p1.bqimg.com/567571/21b3d09e27e01ec1.gif)


###output(left tricket)
![](http://p1.bpimg.com/567571/bd4a89e17aa0bde0.gif)

###dependence
ASCLL TABLE Writer it is for generate ascii table on termial
and use below command to install

```
go get  github.com/olekukonko/tablewriter
```

###How to get started?
you should download project to path of yourself and then go build or go install, after that just run it!

```
go build 
go-home your train number (em: G4474)
```

###log

1. Add depart date parameter when query train schedule.
If you use a command like **go-home k502**, default value is tody,
so above command equal **go-home k502 2017-01-22**
In fact you should use a command like **go-home k502 2017-01-23** to query someday's plan that you want to know. 
2. Support query left tricket go througth api of 12306 
###Contact

Wechat: convertxy

QQ: 2698380951

Email: aliasliyu4@gamil.com
