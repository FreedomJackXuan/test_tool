package main

import (
	"./protos"
	"bufio"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

type Army struct {
	ArmyType string
	ArmyId   int32
	Point    int32
	Rarity   int32
}


func sendJoin(c *gosocketio.Client, method string) {
	armys := &protos.Armys{}
	army := &protos.Army{}
	army.ArmyId = 10203
	army.ArmyNum = 10
	armys.Army = append(armys.Army, army)
	common := &protos.Common{}
	common.Msg = "hahah"
	common.Code = 0
	armys.Common = common
	army = &protos.Army{}
	army.ArmyId = 10103
	army.ArmyNum = 20
	armys.Army = append(armys.Army, army)

	bytes, _ :=  proto.Marshal(armys)
	fmt.Println(bytes)
	//c.Emit("cure.hurt_armys",bytes)

	result, err := c.Ack(method, bytes, time.Second*10 )
	if err != nil {
		log.Fatal("ack result error: ",err)
	} else {
		log.Println("Ack result to cure.hurt_armys: ", result)
	}
}

func event (c *gosocketio.Client) {

}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	trans := transport.GetDefaultWebsocketTransport()
	header := http.Header{}
	header.Add("userId", "123")
	trans.RequestHeader = header

	c, err := gosocketio.Dial(
		gosocketio.GetUrl("127.0.0.1", 7999, false),
		trans)
	if err != nil {
		log.Fatal(err)
	}
	//time.Sleep(1*time.Second)

	//if err!=nil {
	//	log.Fatal("emit error")
	//}
	err = c.On("cure.begining", func(h *gosocketio.Channel) {
		fmt.Println("aaaaaaaaaaaaaaa")
		//fmt.Println(msg)
	})
	if err != nil {
		log.Fatal(err)
	}
	err = c.On("reply", func(h *gosocketio.Channel, msg string) {
		fmt.Println("msg", msg)
	})
	err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		log.Fatal("Disconnected")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("Connected")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On("cure.complete",  func(h *gosocketio.Channel) {
		log.Println("cure.complete")
	})

	reader := bufio.NewReader(os.Stdin)
	for {
		line, _, _ := reader.ReadLine()
		lin := line[0]
		if lin == 49{
			sendJoin(c, "cure.hurting")
		}else if lin == 50 {
			sendJoin(c, "cure.hurt_armys")
		}else if lin == 51 {
			sendJoin(c, "cure.cure_armys")
		}else if lin == 52 {
			sendJoin(c, "cure.begining")
		}else if lin == 53 {
			sendJoin(c, "cure.imediately")
		}
	}
}