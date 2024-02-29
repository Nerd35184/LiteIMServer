package server

import (
	"context"
	"encoding/json"
	"log"
	"server/conf"
	"server/repository"
	"server/util"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (server *Server) Register(ctx context.Context, request *RegisterRequest) (*EmptyResponseBody, error) {
	userInfo := &repository.UserInfo{
		UserId:   util.RandomNumStr(8),
		Username: request.Username,
		Password: request.Password,
	}
	err := server.Db.AddUserInfo(userInfo)
	return nil, err
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserId    string `json:"userId"`
	Token     string `json:"token"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Signature string `json:"signature"`
}

func (server *Server) Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error) {
	userInfo, err := server.Db.GetUserInfoByUsernamePassword(request.Username, request.Password)
	if err != nil {
		return nil, err
	}
	token := util.RandomNumLetterStr(16)
	server.Token2UserId[token] = userInfo.UserId
	log.Printf("Login %s", util.ToJsonStr(server.Token2UserId))
	return &LoginResponse{
		Token:     token,
		UserId:    userInfo.UserId,
		Nickname:  userInfo.Nickname,
		Avatar:    userInfo.Avatar,
		Signature: userInfo.Signature,
	}, nil
}

type GetContactListRequest struct {
	ContactType int `json:"contactType"`
	Offset      int `json:"offset"`
	Limit       int `json:"limit"`
}

type ContactListResponseItem struct {
	UserId      string `json:"userId"`
	ContactType int    `json:"contactType"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Signature   string `json:"signature"`
}
type GetContactListResponse struct {
	Count       int64                      `json:"count"`
	ContactList []*ContactListResponseItem `json:"items"`
}

func (server *Server) GetContactList(ctx context.Context, request *GetContactListRequest) (*GetContactListResponse, error) {
	count, err := server.Db.CountUserContactByType(ctx.Value(CTX_USER_ID_KEY).(string), request.ContactType)
	if err != nil {
		return nil, err
	}
	infos, err := server.Db.GetUserContactsByType(ctx.Value(CTX_USER_ID_KEY).(string), request.ContactType, request.Offset, request.Limit)
	if err != nil {
		return nil, err
	}
	contactIds := []string{}
	for _, info := range infos {
		contactIds = append(contactIds, info.ContactId)
	}

	userInfos, err := server.Db.GetUserInfoByIds(contactIds)
	if err != nil {
		return nil, err
	}
	contactList := make([]*ContactListResponseItem, 0, len(contactIds))
	for _, userInfo := range userInfos {
		contactListResponseItem := &ContactListResponseItem{
			UserId:      userInfo.UserId,
			Nickname:    userInfo.Nickname,
			Avatar:      userInfo.Avatar,
			ContactType: request.ContactType,
			Signature:   userInfo.Signature,
		}
		contactList = append(contactList, contactListResponseItem)
	}

	return &GetContactListResponse{
		Count:       count,
		ContactList: contactList,
	}, nil
}

type GetUserInfoRequest struct {
	UserId string `json:"userId"`
}

type GetUserInfoResponse struct {
	UserId    string `json:"userId"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Signature string `json:"signature"`
}

func UserInfoToGetUserInfoResponse(info *repository.UserInfo) *GetUserInfoResponse {
	return &GetUserInfoResponse{
		UserId:    info.UserId,
		Nickname:  info.Nickname,
		Avatar:    info.Avatar,
		Signature: info.Signature,
	}
}

func (server *Server) GetUserInfo(ctx context.Context, request *GetUserInfoRequest) (*GetUserInfoResponse, error) {
	info, err := server.Db.GetUserInfo(request.UserId)
	if err != nil {
		return nil, err
	}
	return UserInfoToGetUserInfoResponse(info), nil
}

type GetUserInfoByNicknameRequest struct {
	Nickname string `json:"nickname"`
	Offset   int    `json:"offset"`
	Limit    int    `json:"limit"`
}

type GetUserInfoByNicknameResponse struct {
	Count int64                  `json:"count"`
	Items []*GetUserInfoResponse `json:"items"`
}

func (server *Server) GetUserInfoByNickname(ctx context.Context, request *GetUserInfoByNicknameRequest) (*GetUserInfoByNicknameResponse, error) {
	count, err := server.Db.CountUserInfoByNickname(request.Nickname)
	if err != nil {
		return nil, err
	}
	infos, err := server.Db.GetUserInfoByNickname(request.Nickname, request.Offset, request.Limit)
	if err != nil {
		return nil, err
	}
	items := []*GetUserInfoResponse{}
	for _, info := range infos {
		item := UserInfoToGetUserInfoResponse(info)
		items = append(items, item)
	}
	return &GetUserInfoByNicknameResponse{
		Count: count,
		Items: items,
	}, nil
}

type SetUserInfoRequest struct {
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Signature string `json:"signature"`
}

type SetUserInfoResponse struct {
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Signature string `json:"signature"`
}

func (server *Server) SetUserInfo(ctx context.Context, request *SetUserInfoRequest) (*SetUserInfoResponse, error) {
	err := server.Db.SetUserInfo(
		ctx.Value(CTX_USER_ID_KEY).(string),
		request.Nickname,
		request.Signature,
		request.Avatar,
	)
	if err != nil {
		return nil, err
	}
	return &SetUserInfoResponse{
		Nickname:  request.Nickname,
		Signature: request.Signature,
		Avatar:    request.Avatar,
	}, nil
}

type AddContactRequest struct {
	UserId string `json:"userId"`
}

type AddContactResponse struct {
}

func (server *Server) AddContact(ctx context.Context, request *AddContactRequest) (*AddContactResponse, error) {
	addAddContactRequest := &repository.AddContactRequest{
		Sender:   ctx.Value(CTX_USER_ID_KEY).(string),
		Receiver: request.UserId,
		Status:   repository.ADD_CONTACT_REQUEST_STATUS_PENDING,
	}
	err := server.Db.AddAddContactRequest(addAddContactRequest)
	if err != nil {
		return nil, err
	}
	config := conf.GetConfig()

	content := &AddContactRequestContent{}

	contentByte, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}
	sendMsgRequest := &SendMsgRequest{
		MsgType:     repository.MESSAGE_INFO_TYPE_ADD_CONTANT_REQUEST,
		To:          request.UserId,
		ContactType: repository.TBL_CONTACT_INFO_CONTACT_TYPE_SYSTEM_USER,
		Content:     json.RawMessage(contentByte),
	}
	server.SendMsgBySenderId(config.SystemUserIds.AddContactUserId, sendMsgRequest)
	return &AddContactResponse{}, nil
}

type RemoveContantRequest struct {
	UserId string `json:"userId"`
}

type RemoveContantResponse struct {
}

func (server *Server) RemoveContact(ctx context.Context, request *RemoveContantRequest) (*RemoveContantResponse, error) {
	return nil, server.Db.RemoveContact(
		ctx.Value(CTX_USER_ID_KEY).(string), request.UserId,
	)
}

type GetAddContactReqRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type AddContactReqItem struct {
	ReqId     string `json:"req_id"`
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Status    int    `json:"status"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"update_at"`
}

type GetAddContactReqResponse struct {
	Count int64                `json:"count"`
	Items []*AddContactReqItem `json:"items"`
}

func (server *Server) GetAddContactReq(ctx context.Context, request *GetAddContactReqRequest) (*GetAddContactReqResponse, error) {
	userId := ctx.Value(CTX_USER_ID_KEY).(string)
	count, err := server.Db.CountAddContactRequestByReceiverId(userId)
	if err != nil {
		return nil, err
	}

	addContactRequests, err := server.Db.GetAddContactRequestByReceiverId(
		userId,
		request.Offset,
		request.Limit,
	)

	items := []*AddContactReqItem{}
	for _, addContactRequest := range addContactRequests {
		item := &AddContactReqItem{
			Sender:    addContactRequest.Sender,
			Receiver:  addContactRequest.Receiver,
			Status:    addContactRequest.Status,
			CreatedAt: addContactRequest.CreatedAt,
			UpdatedAt: addContactRequest.UpdatedAt,
		}
		items = append(items, item)
	}

	return &GetAddContactReqResponse{
		Items: items,
		Count: count,
	}, nil
}

type ConfirmAddContantReqRequest struct {
	Status int
	ReqId  string
}

type ConfirmAddContantReqResponse struct {
}

func (server *Server) ConfirmAddContantReq(ctx context.Context, request *ConfirmAddContantReqRequest) (*ConfirmAddContantReqResponse, error) {
	userId := ctx.Value(CTX_USER_ID_KEY).(string)
	err := server.Db.UpdateAddContactRequestStatus(request.ReqId, userId, request.Status)
	return &ConfirmAddContantReqResponse{}, err
}
