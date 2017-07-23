package handler

import (
	"fmt"
	"strings"
	"log"
)


func ListenCommand(sess Sess_info)  {

	ws := sess.Ws

	for {

		online := SessionGet(sess.Uid,"online")
		if online,ok := online.(int); !ok || online==0 {
			return
		}

		buffer := make([]byte,1024)
		rlen,err := ws.Read(buffer)
		if err != nil {
			Logout(sess,true)
			return
		}

		data := string(buffer[:rlen])
		parse_result := strings.SplitN(data,":",2)
		if len(parse_result)!=2 {
			log.Println(fmt.Sprintf("response invalid data: %v",parse_result))
			ws.Write([]byte(" data format invalid : "+data))
			continue
		}

		call_cmd :=parse_result[0]
		response :=parse_result[1]

		cmd_handler := Command{}
		switch call_cmd {
		case "send" : //发消息
			go cmd_handler.Send(sess,response)
		case "logout"://退出
			go Logout(sess,true)
		default:
			ws.Write([]byte("undefined cmd : "+call_cmd))
		}

	}


}

func ListenMessage(sess Sess_info)  {

	my_msg_chan := SessionGet(sess.Uid,"msg_chan")

	if my_msg_chan==nil {
		sess.Ws.Write([]byte(fmt.Sprintf("exception: you msg chan invalid ")))
		return
	}

	for {

		online := SessionGet(sess.Uid,"online")
		if online,ok := online.(int); !ok || online==0 {
			return
		}

		if my_msg_chan,ok := my_msg_chan.(chan string); ok {
			receive := <-my_msg_chan
			sess.Ws.Write([]byte(fmt.Sprintf("receive:%s",receive)))
		}

	}

}

func Logout(sess Sess_info , close_conn bool)  {
	if close_conn {
		sess.Ws.Close()
	}
	SessionSet(sess.Uid,"online",0)
}