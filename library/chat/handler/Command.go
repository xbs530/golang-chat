package handler

import (
	"strings"
	"strconv"
	"fmt"
)

type Command struct {

}

//个人消息
func (cmd Command) Send(sess Sess_info, data string)  {
	msg_info := strings.SplitN(data,":",2)
	if len(msg_info)!=2 {
		sess.Ws.Write([]byte("send fail: invalid format (0x01) "))
		return
	}
	to_uid,error := strconv.Atoi(msg_info[0])
	if error!=nil {
		sess.Ws.Write([]byte("send fail: invalid format (0x02) "))
		return
	}

	send(sess,to_uid,msg_info[1],"user")

	sess.Ws.Write([]byte(fmt.Sprintf("success:send to user-%d success ",to_uid)))
}

//广播消息
func (cmd Command) Broadcast(sess Sess_info, data string)  {
	user_count:=len(*sess.users)
	broadcast_chan := make(chan int)
	for _,to_uid := range *sess.users  {
		go func(to_uid int){
			send(sess,to_uid,data,"broadcast")
			broadcast_chan <- 1
		}(to_uid)
	}


	send_count:=0
	for {
		<-broadcast_chan
		send_count++
		if send_count>=user_count {
			close(broadcast_chan)
			goto endfor
		}
	}
	endfor:

	sess.Ws.Write([]byte(fmt.Sprintf("success:broadcast %d users ",user_count)))

}

func send(sess Sess_info,to_uid int,msg string,msg_type string)  {
	tmp_msg_chan := SessionGet(to_uid,"msg_chan")
	if tmp_msg_chan==nil {
		sess.Ws.Write([]byte(fmt.Sprintf("error:user-%d offline or not exists ",to_uid)))
		return
	}
	to_msg_chan,ok := tmp_msg_chan.(chan string)
	if !ok {
		sess.Ws.Write([]byte(fmt.Sprintf("error:user-%d msg chan exception ",to_uid)))
		return
	}

	to_msg := fmt.Sprintf("%d:%s:%s",sess.Uid,msg_type,msg)
	to_msg_chan <- to_msg
}

