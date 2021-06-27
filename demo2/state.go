package main

import (
	"bufio"
	"github.com/lni/goutils/syncutil"
	"log"
	"os"
	"rogo/node"
	"rogo/pb"
	"strconv"
	"strings"
)

const sId = 128
func main(){


	n:=node.NewNode("127.0.0.1","12212",2)
	s:=&stateHandle{n: n}
	err := n.JoinKvCluster(sId,s )
	if err != nil {
		return
	}

	raftStopper := syncutil.NewStopper()
	consoleStopper := syncutil.NewStopper()
	ch := make(chan string, 16)

	consoleStopper.RunWorker(func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			s, err := reader.ReadString('\n')
			if err != nil {
				close(ch)
				return
			}
			if s == "exit\n" {
				raftStopper.Stop()
				n.LogoutCluster(sId)
				return
			}
			//	log.Println("哈哈",s)
			ch <- s
		}
	})
	//printUsage()

	raftStopper.RunWorker(func() {
		//	cs := n.GetSession(sId)
		for {
			select {
			case v, ok := <-ch:
				if !ok {
					return
				}
				msg := strings.Replace(v, "\n", "", 1)

				s.parseCommand(msg)
			case <-raftStopper.ShouldStop():
				return
			}
		}
	})
	raftStopper.Wait()


}

type stateHandle struct {
	n node.Node

}

func (s *stateHandle) SyncMessage(message *pb.Message) {
	log.Println(message)
}


func (s *stateHandle)parseCommand(msg string)  {
	parts := strings.Split(strings.TrimSpace(msg), " ")
	if len(parts) == 0 {
		log.Println("键入命令错误")
		return
	}


	switch parts[0] {
	case "put":
		if len(parts) != 3 {
			log.Println("put 键入命令错误")
			return
		}
		err := s.n.SyncData(sId, parts[1], []byte(parts[2]))
		if err != nil {
			log.Println(err)
			return
		}

	case "get":
		if len(parts) != 2 {
			log.Println("get 键入命令错误")
			return
		}
		data, err := s.n.GetSyncData(sId, parts[1])
		if err != nil {
			log.Println("get 解析错误:",err)
			return
		}
		log.Println("key:",parts[1]," value:",string(data))
	case "del":
		if len(parts) != 2 {
			log.Println("del 键入命令错误")
			return
		}
		err := s.n.RemoveSyncData(sId, parts[1])
		if err != nil {
			log.Println("del 解析错误:",err)
			return
		}
	case "leader":
		log.Println("未启用")
		return
	case "addNode":
		//log.Println(len(parts))
		if len(parts) != 3 {
			log.Println("addNode 键入命令错误")
			return
		}
		keyid, _ := strconv.ParseInt(parts[1], 10, 64)
		err := s.n.AddNode(sId, uint64(keyid), parts[2])
		if err != nil {
			log.Println("addNode 解析错误:",err)
			return
		}
		log.Println("添加成功")
	case "delNode":
		log.Println("未启用")
		return
	case "msg":
		if len(parts) < 2 {
			log.Println("msg 键入命令错误")
			return
		}
		err := s.n.SendMsg(sId, 11, []byte(parts[1]))
		if err != nil {
			log.Println("msg 解析错误:",err)
			return
		}


	}


}