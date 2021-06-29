package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"rogo"
	"rogo/pb"
	"strings"
	"time"
)

func main() {

	n := rogo.NewNode("127.0.0.1", "22112", 33)

	sf := &stateFollow{n}
	err := n.JoinStateCluster(102, sf)
	if err != nil {
		panic(err)
	}

	cluster, err := n.RequestJoinCluster("127.0.0.1:22222", 102, "123")
	if err != nil {
		panic(err)
	}
	log.Println(cluster)

	sf.listen()

}

type stateFollow struct {
	node rogo.Node
}

func (f *stateFollow) SyncMessage(message *pb.Message) {
	log.Println(message)
}

func (f *stateFollow) listen() {
	for {
		reader := bufio.NewReader(os.Stdin)
		s, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			continue
		}
		msg := strings.Replace(s, "\n", "", 1)

		parts := strings.Split(strings.TrimSpace(msg), " ")
		if len(parts) == 0 {
			continue
		}
		switch parts[0] {
		case "put":
			if len(parts) != 3 {
				log.Println("不等于3")
				continue
			}
			err := f.node.SyncData(102, parts[1], []byte(parts[2]))
			if err != nil {
				log.Println(err)
				continue
			}
		case "get":
			if len(parts) != 2 {
				log.Println("不等于2")

				continue
			}

			data, err := f.node.GetSyncData(102, parts[1])
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println(string(data))
			hs := f.node.GetNodeHost()

			index, err := hs.ReadIndex(102, 5*time.Second)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println("开始")
			//hs.RequestAddObserver()
			<-index.AppliedC()
			log.Println("完成")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			session, err := hs.SyncGetSession(ctx, 102)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println(session.GetClientID())
			cancel()

			node, err := hs.ReadLocalNode(index, []byte(parts[1]))
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println(string(node.([]byte)))

		case "msg":
			if len(parts) != 2 {
				log.Println("不等于2")

				continue
			}
			f.node.SendMsg(102, 222, []byte(parts[1]))

		}

	}
}
