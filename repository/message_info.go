package repository

import "fmt"

const (
	TBL_MESSAGE_INFO                      = "message_info"
	TBL_MESSAGE_INFO_COL_MEG_ID           = "msg_id"
	TBL_MESSAGE_INFO_COL_SESS_ID          = "sess_id"
	TBL_MESSAGE_INFO_COL_CREATED_AT       = "created_at"
	MESSAGE_INFO_TYPE_TEXT                = 1
	MESSAGE_INFO_TYPE_ADD_CONTANT_REQUEST = 2
)

type MessageInfo struct {
	Id        int64
	SessId    string `gorm:"type:varchar(255)"`
	MsgId     int64
	SenderId  string `gorm:"type:varchar(255)"`
	MsgType   int
	Content   string `gorm:"type:varchar(1024)"`
	CreatedAt int64
}

func (MessageInfo MessageInfo) TableName() string {
	return TBL_MESSAGE_INFO
}
func (MessageInfo MessageInfo) Indexes() []*MySqlIndex {
	return []*MySqlIndex{}
}

func (mysql *Mysql) AddMessageInfo(messageInfo *MessageInfo) error {
	return mysql.DB.Create(messageInfo).Error
}

func (mysql *Mysql) GetMaxMsgId(sessId string) (int64, error) {
	sql := fmt.Sprintf(`select max(%s) as m from %s where %s = ?`, TBL_MESSAGE_INFO_COL_MEG_ID, TBL_MESSAGE_INFO, TBL_MESSAGE_INFO_COL_SESS_ID)
	result := MaxModel{}
	err := mysql.DB.Raw(sql, sessId).Scan(&result).Error
	return result.M, err
}

func (mysql *Mysql) GetMsgInfos(sessId string, beginMsgId int64, beginCreatedAt int64, offset int, limit int) ([]*MessageInfo, error) {
	sql := fmt.Sprintf(`select * from %s where %s = ? and %s > and %s > ? limit ?,? asc by %s`,
		TBL_MESSAGE_INFO,
		TBL_MESSAGE_INFO_COL_SESS_ID,
		TBL_MESSAGE_INFO_COL_MEG_ID,
		TBL_MESSAGE_INFO_COL_CREATED_AT, TBL_MESSAGE_INFO_COL_MEG_ID)
	result := []*MessageInfo{}
	err := mysql.DB.Raw(sql, sessId, beginMsgId, beginCreatedAt, offset, limit).Scan(&result).Error
	return result, err
}
