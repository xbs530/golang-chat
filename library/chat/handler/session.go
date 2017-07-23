package handler

import (
	"golang.org/x/net/websocket"
	"log"
	"fmt"
	"time"
	"strings"
	"strconv"
)

//会话结构
type Sess_info struct {
	Ws *websocket.Conn
	Uid int
}

//在线人数
var connection_count = 0
//用户编号
var user_no=0

//全局容器
var di=make(map[string]interface{})


func init()  {
	fmt.Println("server init ...")
}


//开启会话
func Session(ws *websocket.Conn)  {

	connection_count+=1
	user_no+=1

	sess := Sess_info{ws,user_no}

	SessionSet(sess.Uid,"online",1)

	//初始化消息通道
	SessionSet(sess.Uid,"msg_chan",make(chan string))

	go log.Println(fmt.Sprintf("the user-%d connect success(online:%d) ... ",sess.Uid,connection_count))

	//监听命令消息
	go ListenCommand(sess)

	//监听聊天消息
	go ListenMessage(sess)

	//心跳检测
	go Heartbeat(sess)

	//检查在线状态
	for {
		online := SessionGet(sess.Uid,"online")
		if online,ok := online.(int); !ok || online==0 {
			connection_count-=1
			log.Println(fmt.Sprintf("user=%d logout",sess.Uid))
			SessionDestory(sess.Uid)
			return
		}
		time.Sleep(1*time.Second)
	}

}

//会话变量获取
func SessionGet(uid int, field string) interface{}  {
	tmp_key := fmt.Sprintf("%d-%s",uid,field)
	if _,ok := di[tmp_key]; ok {
		return di[tmp_key]
	}
	return nil
}

//会话变量修改
func SessionSet(uid int, field string,value interface{})  {
	tmp_key := fmt.Sprintf("%d-%s",uid,field)
	di[tmp_key]=value
}

//会话销毁
func SessionDestory(uid int)  {

	for k,_ := range di {
		tmp :=strings.SplitN(k,"-",2)
		if prefix,_ := strconv.Atoi(tmp[0]);prefix==uid {
			delete(di,k)
		}
	}
}

//心跳检测
func Heartbeat(sess Sess_info)  {

	for  {

		online := SessionGet(sess.Uid,"online")
		if online,ok := online.(int); !ok || online==0 {
			return
		}

		ping_data := fmt.Sprintf("ping:%d",time.Now().Second())
		_,err := sess.Ws.Write([]byte(ping_data))
		if err != nil {
			Logout(sess,false)
			return
		}
		time.Sleep(5*time.Second)
	}

}
