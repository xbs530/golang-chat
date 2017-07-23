package handler

import (
	"strings"
	"strconv"
	"fmt"
)

type Command struct {

}

func (cmd Command) Send(sess Sess_info, data string)  {
	msg_info := strings.SplitN(data,":",2)
	if len(msg_info)!=2 {
		sess.Ws.Write([]byte("send fail: invalid format (0x01) "))
		return
	}
	to_uid,error := strconv.Atoi(msg_info[0]);
	if error!=nil {
		sess.Ws.Write([]byte("send fail: invalid format (0x02) "))
		return
	}

	tmp_msg_chan := SessionGet(to_uid,"msg_chan")
	if tmp_msg_chan==nil {
		sess.Ws.Write([]byte(fmt.Sprintf("send fail: user-%d offline or not exists ",to_uid)))
		return
	}
	to_msg_chan,ok := tmp_msg_chan.(chan string)
	if !ok {
		sess.Ws.Write([]byte(fmt.Sprintf("send fail: user-%d msg chan exception ",to_uid)))
		return
	}

	to_msg := fmt.Sprintf("%d:%s",sess.Uid,msg_info[1])
	to_msg_chan <- to_msg

	sess.Ws.Write([]byte(fmt.Sprintf("send to user-%d success ",to_uid)))
}

