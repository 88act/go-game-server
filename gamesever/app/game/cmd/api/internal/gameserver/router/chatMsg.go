package router

import (
	"time"

	"go-cms/app/game/cmd/api/internal/gameserver/core"
	"go-cms/app/game/cmd/api/internal/gameserver/pb"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/protobuf/proto"
)

type ChatMsgRouter struct {
	znet.BaseRouter
}

func (m *ChatMsgRouter) Handle(req ziface.IRequest) {
	ctx := req.GetConnection().Context()
	pID, err := req.GetConnection().GetProperty("pID")
	if err != nil {
		req.GetConnection().Stop()
		logMsg := getLogMsg(req) + err.Error()
		logc.Error(ctx, logMsg)
		return
	}
	player := core.WorldMgrObj.GetPlayerByPID(pID.(int64))
	msg := &pb.ChatMsg{}
	err = proto.Unmarshal(req.GetData(), msg)
	if err != nil {
		logMsg := getLogMsg(req) + err.Error()
		logc.Error(ctx, logMsg)
		return
	}
	msg.SendTime = time.Now().Unix()
	resp := &pb.ChatMsgResp{}
	resp.MsgList = append(resp.MsgList, msg)
	if msg.ChatType == 3 {
		//发给对方
		core.GetWM().SendOne(ctx, pb.S_ChatMsgResp, resp, msg.ObjId)
		//也发给自己
		player.SendMsg(pb.S_ChatMsgResp, resp)
	} else if msg.ChatType == 2 {
		core.GetWM().SendRoom(ctx, pb.S_ChatMsgResp, resp, player.RoomId)
	} else if msg.ChatType == 1 {
		core.GetWM().SendWorld(ctx, pb.S_ChatMsgResp, resp)
	}
}