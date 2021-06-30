package main

import (
	"bufio"
	"log"
	"os"
	"rogo"
	"rogo/pb"
	"strings"
)

func main() {

	n := rogo.NewNode("127.0.0.1", "13204", 3)
	f := &follower{node: n}

	err := n.StartStorageCluster(102, rogo.StorageTypeMemory, rogo.RoleFollower)
	if err != nil {
		panic(err)
	}

	//cluster, err := n.RequestJoinCluster("127.0.0.1:22222", 102, "123", rogo.RoleFollower)
	//if err != nil {
	//	panic(err)
	//}
	//log.Println(cluster)
	f.listen()
}

type follower struct {
	node rogo.Node
}

func (f *follower) SyncMessage(message *pb.Message) {
	log.Println(message)
}

func (f *follower) listen() {
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
			data, err := f.node.GetData(102, parts[1])
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println(string(data))
		case "msg":
			if len(parts) != 2 {
				log.Println("不等于2")

				continue
			}
			f.node.SendMessage(102, 222, []byte(parts[1]))
		case "del":
			if len(parts) != 2 {
				log.Println("不等于2")

				continue
			}
			err := f.node.RemoveData(102, parts[1])
			if err != nil {
				log.Println(err)

				continue

			}
			log.Println("del ok")

		}

	}
}
