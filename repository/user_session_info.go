package repository

import "time"

const (
	TBL_USER_SESS_INFO = "user_sess_info"
)

type UserSessionInfo struct {
	Id        int64
	UserId    string `gorm:"type:varchar(255)"`
	SessId    string `gorm:"type:varchar(255)"`
	CreatedAt int64
	UpdatedAt int64
}

func (UserSessionInfo UserSessionInfo) TableName() string {
	return TBL_USER_SESS_INFO
}

func (UserSessionInfo UserSessionInfo) Indexes() []*MySqlIndex {
	return []*MySqlIndex{}
}

func (mysql *Mysql) AddUserSessionInfo(userSessionInfo *UserSessionInfo) error {
	return mysql.DB.Create(userSessionInfo).Error
}

func (mysql *Mysql) CreateUserSessonInfo(userId string, sessId string) error {
	userSessionInfo := &UserSessionInfo{
		UserId:    userId,
		SessId:    sessId,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	return mysql.AddUserSessionInfo(userSessionInfo)
}
