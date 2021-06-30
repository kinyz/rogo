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

	n := rogo.NewNode("127.0.0.1", "22112", 33)

	sf := &stateFollow{n}
	err := n.StartStorageCluster(102, rogo.StorageTypeDisk, rogo.RoleObserver)
	if err != nil {
		panic(err)
	}

	cluster, err := n.RequestJoinCluster("127.0.0.1:22222", 102, "123", rogo.RoleObserver)
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

			//
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
			err := f.node.SendMessage(102, 222, []byte(parts[1]))
			if err != nil {
				log.Println(err)
			}

		}

	}
}
