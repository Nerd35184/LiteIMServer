package server

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"server/repository"
	"server/util"
	"time"
)

const (
	WS_MSG_TYPE_AUTH = 1
)

func MessageInfoToGetMsgResponseItem(msgInfo *repository.MessageInfo) (*GetMsgResponseItem, error) {
	msgContent, err := util.DecodeBase64Str(msgInfo.Content)
	if err != nil {
		return nil, err
	}
	item := &GetMsgResponseItem{
		SessId:    msgInfo.SessId,
		MsgId:     msgInfo.MsgId,
		SenderId:  msgInfo.SenderId,
		MsgType:   msgInfo.MsgType,
		CreatedAt: msgInfo.CreatedAt,
		Content:   msgContent,
	}
	return item, nil
}

type TextMsgContent struct {
	Text string `json:"text"`
}

type AddContactRequestContent struct {
	UserId string `json:"userId"`
}

type SendMsgRequest struct {
	MsgType     int             `json:"msgType"`
	To          string          `json:"to"`
	ContactType int             `json:"contactType"`
	Content     json.RawMessage `json:"content"`
}

type SendMsgResponse struct {
}

func (server *Server) SendMsgBySenderId(userId string, request *SendMsgRequest) (*SendMsgResponse, error) {
	_, err := server.Db.GetUserContact(userId, request.ContactType, request.To)
	if err != nil {
		log.Printf("SendMsgBySenderId GetUserContact err:%s", err)
		return nil, err
	}
	dialogSessInfo, err := server.GetOrCreateOneToOneDialogueSess(userId, request.To)
	if err != nil {
		log.Printf("SendMsgBySenderId GetOrCreateOneToOneDialogueSess err:%s", err)
		return nil, err
	}
	msgInfo := &repository.MessageInfo{
		SenderId:  userId,
		SessId:    dialogSessInfo.Id,
		MsgType:   request.MsgType,
		MsgId:     dialogSessInfo.GenerateMsgId(),
		Content:   util.EncodeBase64Str(request.Content),
		CreatedAt: time.Now().Unix(),
	}
	err = server.Db.AddMessageInfo(msgInfo)
	if err != nil {
		log.Printf("SendMsgBySenderId AddMessageInfo err:%s", err)
		return nil, err
	}

	msgItem, err := MessageInfoToGetMsgResponseItem(msgInfo)
	if err != nil {
		log.Printf("SendMsgBySenderId AddMessageInfo err:%s", err)
		return nil, err
	}
	for _, userId := range dialogSessInfo.GetUserIds() {
		log.Printf("SendMsgBySenderId GetUserIds %s", userId)
		connSess, ok := server.ConnSesses[userId]
		if !ok {
			continue
		}
		connSess.WriteJson(msgItem)
	}
	return &SendMsgResponse{}, nil
}

func (server *Server) SendMsg(ctx context.Context, request *SendMsgRequest) (*SendMsgResponse, error) {
	if request.ContactType != repository.TBL_CONTACT_INFO_CONTACT_TYPE_REGULAR_USER {
		return nil, errors.New("only support single for now")
	}
	userId := ctx.Value(CTX_USER_ID_KEY).(string)
	return server.SendMsgBySenderId(userId, request)
}

type GetMsgRequest struct {
	BeginMsgId     int64  `json:"begin_msg_id"`
	BeginCreatedAt int64  `json:"begin_created_at"`
	SessId         string `json:"sess_id"`
	Offset         int    `json:"offset"`
	Limit          int    `json:"limit"`
}

type GetMsgResponseItem struct {
	SessId    string          `json:"sess_id"`
	MsgId     int64           `json:"msg_id"`
	SenderId  string          `json:"sender_id"`
	MsgType   int             `json:"msg_type"`
	Content   json.RawMessage `json:"content"`
	CreatedAt int64           `json:"created_at"`
}

type GetMsgResponse struct {
	Items []*GetMsgResponseItem
}

func (server *Server) GetMsg(ctx context.Context, request *GetMsgRequest) (*GetMsgResponse, error) {
	//todo 检查是否有权限获取
	msgs, err := server.Db.GetMsgInfos(request.SessId, request.BeginMsgId, request.BeginCreatedAt, request.Offset, request.Limit)
	if err != nil {
		return nil, err
	}
	items := []*GetMsgResponseItem{}
	for _, msg := range msgs {
		item, err := MessageInfoToGetMsgResponseItem(msg)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return &GetMsgResponse{
		Items: items,
	}, nil
}

type AuthMsg struct {
	Token string `json:"token"`
}

func (server *Server) HandleWS(w http.ResponseWriter, r *http.Request) {
	log.Printf("Server HandleWS")
	conn, err := server.WSUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Server HandleWS Upgrade err:%s", err)
		return
	}
	defer conn.Close()
	msgType, body, err := conn.ReadMessage()
	if err != nil {
		log.Printf("Server HandleWS ReadMessage err:%s", err)
		return
	}
	if msgType != WS_MSG_TYPE_AUTH {
		log.Printf("Server HandleWS WS_MSG_TYPE_AUTH err:%s", err)
		return
	}
	authMsg := &AuthMsg{}
	err = json.Unmarshal(body, authMsg)
	if err != nil {
		log.Printf("Server HandleWS Unmarshal err:%s", err)
		return
	}

	userId, ok := server.Token2UserId[authMsg.Token]
	if !ok {
		log.Printf("HandleWS token not found")
		w.Write([]byte("token not found"))
		return
	}

	connSess, ok := server.ConnSesses[userId]
	if !ok {
		connSess = NewConnSession(userId, nil)
		server.ConnSesses[userId] = connSess
	}
	err = connSess.SetConn(conn)
	if err != nil {
		return
	}
	defer connSess.SetConn(nil)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
