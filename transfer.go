// implement of transfer router table based on Dijkstra
package main

import ()

/*
  表头结构
  ---------------
 |  车次码 | 指针 |
  ---------------

  头节点
  --------------------
 |站码 |相关信息 |指针  |
  --------------------

  表内节点
  --------------------------------------------------
 | 站码 | 里程 | 线码| 同站的下一车次 | 同一车次的下一个站 |
  --------------------------------------------------
*/

type TableHeader struct {
	trainCode string
	node      *Node
}

type Header struct {
	trainCode string
	info      string
	next      *Header
}

type Node struct {
	stationCode string
	mileage     float64
	lineCode    int
	nextTrain   string
	nextStation *Node
}
