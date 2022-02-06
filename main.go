package main

import (
	"bytes"
	"fmt"
	suid "github.com/gofrs/uuid"
	"github.com/google/uuid"
	"io"
	"math/rand"
	"net"
	"time"
)

type Out struct {
	c   net.Conn
	err error
}

func main() {
	//network.CLIENT
	//fg := flag.Int("fg", 1, "")
	//flag.Parse()
	a := make([]byte, 0, 4)
	//a = []byte{1, 2, 3, 4}
	b := make([]byte, 4)
	b = []byte{5, 6, 7, 8}
	a = append(a, b...)
	fmt.Printf(">>>>>:::%+v\n", a)
	rand.Seed(12)
	fmt.Printf(">>>r:%v\n", rand.Uint32())
	fmt.Printf(">>>r:%v\n", rand.Uint32())
	rd := bytes.NewReader([]byte("aaaaaaaaaaa"))
	n := 3
	var err error
	buf := make([]byte, 3)
	for n != 0 {
		n, err = io.ReadFull(rd, buf)
		fmt.Printf("=======n:%v,err:%v\n", n, err)
	}
	//dur := time.NewTicker(time.Duration(1) * time.Second)
	dur := time.NewTimer(1 * time.Second)
	for i := 0; i < 1; i++ {
		select {
		case <-dur.C:
			fmt.Printf("===i:%+v\n", i)
		}
		dur.Reset(1 * time.Second)
	}
	//wg := sync.WaitGroup{}
	//wg.Add(2)
	//go func() {
	//	//defer wg.Done()
	//	<-dur.C
	//	fmt.Printf("1,结束\n")
	//}()

	//dur.Reset(6 * time.Second)
	//dur.
	//time.Sleep(time.Second * 2)
	//dur.Reset(2 * time.Second)
	//fmt.Printf("-----------\n")
	//<-dur.C
	////<-dur.C
	//fmt.Printf("2,结束\n")
	//dur.Stop()
	//fmt.Printf("after stop\n")
	//<-dur.C
	fmt.Printf("3,结束\n")
	uid := uuid.New().String()
	uid = uuid.NewString()
	fmt.Printf("uid :%v\n", uid)
	su, _ := suid.NewV4()
	fmt.Printf("suid:%v\n", su.String())
	//tm := time.NewTimer(2 * time.Second)

	//if *fg == 1 {
	//	lst, err := net.Listen("tcp", ":18888")
	//	if err != nil {
	//		fmt.Printf("listen error:%+v\n", err)
	//	}
	//	oc := make(chan Out, 1)
	//	go func(o chan Out) {
	//		defer func() {
	//			fmt.Printf("listener accept exit ok\n")
	//		}()
	//		c, err := lst.Accept()
	//		fmt.Printf("accept 之后:::%+v,%+v\n", c, err)
	//		o <- Out{c, err}
	//
	//	}(oc)
	//	dur := time.Tick(50 * time.Second)
	//
	//	//for {
	//	var o Out
	//	fmt.Printf("########%+v\n", o)
	//	select {
	//	case o = <-oc:
	//		fmt.Printf("o = <-oc:%+v,%+v\n", o.c, o.err)
	//		//lst.Close()
	//
	//	case <-dur:
	//		fmt.Printf("time out\n")
	//		lst.Close()
	//		return
	//	}
	//	//}
	//	//bbb := make([]byte, 100)
	//	//fmt.Printf("开始读数据\n")
	//	//num, err1 := o.c.Read(bbb)
	//	//fmt.Printf("读数据结束,num:%+v,err1:%+v\n", num, err1)
	//	time.Sleep(10 * time.Second)
	//	o.c.Close()
	//	fmt.Printf("关闭conn\n")
	//	time.Sleep(1000 * time.Second)
	//	//fmt.Printf("accept ok,c:%+v\n", c)
	//}
	//if *fg == 0 {
	//	dur := time.Tick(50 * time.Second)
	//	for {
	//		select {
	//		case <-dur:
	//			fmt.Printf("client timeout\n")
	//			return
	//		default:
	//			c, err := net.Dial("tcp", "127.0.0.1:18888")
	//			if err != nil {
	//				//fmt.Printf("Dial error:%+v\n", err)
	//				continue
	//			}
	//			fmt.Printf("dial ok,c:%+v\n", c)
	//			//c.Close()
	//			bbb := make([]byte, 100)
	//			num, err1 := c.Read(bbb)
	//			fmt.Printf("cli读数据结束,num:%+v,err1:%+v\n", num, err1)
	//			time.Sleep(500 * time.Second)
	//			return
	//		}
	//
	//	}
	//	//c, err := net.Dial("tcp", "127.0.0.1:18888")
	//	//if err != nil {
	//	//	fmt.Printf("Dial error:%+v\n", err)
	//	//	return
	//	//}
	//	//fmt.Printf("dial ok,c:%+v\n", c)
	//}
	fmt.Printf("==========\n\n")
}
