package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"rogo"
	"strings"
	"sync/atomic"
)

func main() {
	n := rogo.NewNode("127.0.0.1", "13202", 1)
	ma := &mainCluster{nodeId: 2, key: "123", node: n}

	err := n.CreatStorageCluster(102, rogo.StorageTypeDisk, false)
	if err != nil {
		panic(err)
	}

	err = n.StartOauthService("0.0.0.0:22222", ma)
	if err != nil {
		panic(err)
	}

	ma.listen()
}

type mainCluster struct {
	nodeId uint64

	key  string
	node rogo.Node
}

func (m *mainCluster) Oauth(req rogo.RequestJoinCluster) error {
	if req.GetKey() != m.key {
		return errors.New("key错误")
	}
	atomic.AddUint64(&m.nodeId, +1)
	log.Println("请求连接:", req.NodeId)
	return nil

}

func (m *mainCluster) SyncMessage(message rogo.Message) {
	log.Println(message)
}

func (m *mainCluster) listen() {
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
			err := m.node.SyncData(102, parts[1], []byte(parts[2]))
			if err != nil {
				log.Println(err)
				continue
			}
		case "get":
			if len(parts) != 2 {
				log.Println("不等于2")

				continue
			}
			data, err := m.node.GetSyncData(102, parts[1])
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println(string(data))
		}
	}
}
