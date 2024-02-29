package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"server/conf"
	"server/repository"
	"server/util"
	"sync"

	"github.com/gorilla/websocket"
)

type ConnSession struct {
	UserId string
	Conn   *websocket.Conn
}

func (connSession *ConnSession) WriteJson(body interface{}) error {
	log.Printf("connSession WriteJson %s %s", connSession.UserId, util.ToJsonStr(body))
	conn := connSession.Conn
	if conn == nil {
		return errors.New("not connect")
	}
	return connSession.Conn.WriteJSON(body)
}

func (connSession *ConnSession) SetConn(conn *websocket.Conn) error {
	connSession.Conn = conn
	return nil
}

func NewConnSession(
	UserId string,
	Conn *websocket.Conn,
) *ConnSession {
	return &ConnSession{
		UserId: UserId,
		Conn:   Conn,
	}
}

type DialogueSessLiteInfo struct {
	Id      string
	Seq     int64
	L       sync.Locker
	UserIds []string
}

func (dialogueSessLiteInfo *DialogueSessLiteInfo) GetId() string {
	return dialogueSessLiteInfo.Id
}

//todo，做群聊的时候，这里应该有并发问题
func (dialogueSessLiteInfo *DialogueSessLiteInfo) GetUserIds() []string {
	return dialogueSessLiteInfo.UserIds
}

func (dialogueSessLiteInfo *DialogueSessLiteInfo) GenerateMsgId() int64 {
	dialogueSessLiteInfo.L.Lock()
	defer dialogueSessLiteInfo.L.Unlock()
	dialogueSessLiteInfo.Seq++
	return dialogueSessLiteInfo.Seq
}

type Server struct {
	HttpHostPort             string
	TcpHostPort              string
	Token2UserId             map[string]string
	ConnSesses               map[string]*ConnSession
	DialogueSessLiteInfos    map[string]*DialogueSessLiteInfo
	Db                       *repository.Mysql
	StaticFileSystemPathRoot string
	StaticHttpRootPath       string
	L                        sync.Mutex
	WSUpgrader               *websocket.Upgrader
}

func GenereteOneToOneTypeSessionId(userA string, userB string) string {
	if userA < userB {
		return fmt.Sprintf("%s_%s", userA, userB)
	}
	return fmt.Sprintf("%s_%s", userB, userA)
}

func (server *Server) GetOrCreateOneToOneDialogueSess(userA string, userB string) (*DialogueSessLiteInfo, error) {
	var id string
	if userA < userB {
		id = fmt.Sprintf("%s_%s", userA, userB)
	} else {
		id = fmt.Sprintf("%s_%s", userB, userA)
	}
	server.L.Lock()
	defer server.L.Unlock()
	dialogueSessLiteInfo, ok := server.DialogueSessLiteInfos[id]
	if !ok {
		max, err := server.Db.GetMaxMsgId(id)
		if err != nil {
			return nil, err
		}
		server.DialogueSessLiteInfos[id] = &DialogueSessLiteInfo{
			Id:      id,
			UserIds: []string{userA, userB},
			Seq:     max,
			L:       &sync.Mutex{},
		}
	}
	return dialogueSessLiteInfo, nil

}

func NewServer(config *conf.Config) (*Server, error) {
	db, err := repository.InitMysqlDb(config.MysqlConfig)
	if err != nil {
		return nil, err
	}

	server := &Server{
		TcpHostPort:              config.TcpHostPort,
		HttpHostPort:             config.HttpHostPort,
		ConnSesses:               make(map[string]*ConnSession),
		Db:                       db,
		StaticFileSystemPathRoot: config.StaticFileSystemPathRoot,
		DialogueSessLiteInfos:    map[string]*DialogueSessLiteInfo{},
		WSUpgrader:               &websocket.Upgrader{},
		Token2UserId:             make(map[string]string),
		StaticHttpRootPath:       config.StaticHttpRootPath,
	}
	RegisterDataHttpHandleFunc(server, "/static/download/", false, server.Download)
	RegisterDataHttpHandleFunc(server, "/static/upload/", false, server.Upload)

	RegisterJsonHttpHandleFunc(server, "/api/register", false, server.Register)
	RegisterJsonHttpHandleFunc(server, "/api/login", false, server.Login)
	RegisterJsonHttpHandleFunc(server, "/api/get_contact_list", true, server.GetContactList)
	RegisterJsonHttpHandleFunc(server, "/api/get_user_info", true, server.GetUserInfo)
	RegisterJsonHttpHandleFunc(server, "/api/get_user_info_by_nickname", true, server.GetUserInfoByNickname)
	RegisterJsonHttpHandleFunc(server, "/api/set_user_info", true, server.SetUserInfo)

	RegisterJsonHttpHandleFunc(server, "/api/add_contact", true, server.AddContact)
	RegisterJsonHttpHandleFunc(server, "/api/remove_contact", true, server.RemoveContact)
	RegisterJsonHttpHandleFunc(server, "/api/get_add_contact_req", true, server.GetAddContactReq)
	RegisterJsonHttpHandleFunc(server, "/api/confirm_add_contact_req", true, server.ConfirmAddContantReq)

	RegisterJsonHttpHandleFunc(server, "/api/send_msg", true, server.SendMsg)
	RegisterJsonHttpHandleFunc(server, "/api/get_msg", true, server.GetMsg)

	http.HandleFunc("/ws", server.HandleWS)
	return server, nil
}

func (server *Server) Start() error {
	log.Printf("http start %s", server.HttpHostPort)
	return http.ListenAndServe(server.HttpHostPort, nil)
}
