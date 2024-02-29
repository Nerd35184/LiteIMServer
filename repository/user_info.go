package repository

import "fmt"

const (
	TABLE_USER_INFO_NAME          = "user_info"
	TABLE_USER_INFO_COL_USERID    = "user_id"
	TABLE_USER_INFO_COL_USERNAME  = "username"
	TABLE_USER_INFO_COL_PASSWORD  = "password"
	TABLE_USER_INFO_COL_NICKNAME  = "nickname"
	TABLE_USER_INFO_COL_SIGNATURE = "signature"
	TABLE_USER_INFO_COL_AVATAR    = "avatar"
)

type UserInfo struct {
	Id        int64
	UserId    string `gorm:"type:varchar(255)"`
	Username  string `gorm:"type:varchar(255)"`
	Password  string `gorm:"type:varchar(255)"`
	Nickname  string `gorm:"type:varchar(255)"`
	Signature string `gorm:"type:varchar(255)"`
	Avatar    string `gorm:"type:varchar(255)"`
}

func (UserInfo UserInfo) TableName() string {
	return TABLE_USER_INFO_NAME
}

func (UserInfo UserInfo) Indexes() []*MySqlIndex {
	return []*MySqlIndex{
		{
			Unique: true,
			Cols: []string{
				TABLE_USER_INFO_COL_USERID,
			},
		},
		{
			Unique: true,
			Cols: []string{
				TABLE_USER_INFO_COL_USERNAME,
			},
		},
	}
}

func (mysql *Mysql) AddUserInfo(userInfo *UserInfo) error {
	return mysql.DB.Create(userInfo).Error
}

func (mysql *Mysql) SetUserInfo(userId string, nickname string, signature string, avatar string) error {
	sql := fmt.Sprintf(`update %s set %s = ?,%s = ?,%s = ? where %s = ?`, TABLE_USER_INFO_NAME, TABLE_USER_INFO_COL_NICKNAME, TABLE_USER_INFO_COL_SIGNATURE, TABLE_USER_INFO_COL_AVATAR, TABLE_USER_INFO_COL_USERID)
	return mysql.DB.Exec(sql, nickname, signature, avatar, userId).Error
}

func (mysql *Mysql) GetUserInfo(userId string) (*UserInfo, error) {
	sql := fmt.Sprintf(`select * from %s where %s = ?`, TABLE_USER_INFO_NAME, TABLE_USER_INFO_COL_USERID)
	userInfo := &UserInfo{}
	err := mysql.DB.Raw(sql, userId).First(userInfo).Error
	return userInfo, err
}

func (mysql *Mysql) CountUserInfoByNickname(nickname string) (int64, error) {
	sql := fmt.Sprintf(`select %s from %s where %s like ?`, COUNT_SELECT, TABLE_USER_INFO_NAME, TABLE_USER_INFO_COL_NICKNAME)
	c := &CountModel{}
	err := mysql.DB.Raw(sql, "%"+nickname+"%").Scan(c).Error
	return c.Count, err
}

func (mysql *Mysql) GetUserInfoByNickname(nickname string, offset int, limit int) ([]*UserInfo, error) {
	sql := fmt.Sprintf(`select * from %s where %s like ? limit ?,?`, TABLE_USER_INFO_NAME, TABLE_USER_INFO_COL_NICKNAME)
	userInfo := []*UserInfo{}
	err := mysql.DB.Raw(sql, "%"+nickname+"%", offset, limit).Scan(&userInfo).Error
	return userInfo, err
}

func (mysql *Mysql) GetUserInfoByUsernamePassword(username string, password string) (*UserInfo, error) {
	sql := fmt.Sprintf(`select * from %s where %s = ? and %s = ?`, TABLE_USER_INFO_NAME, TABLE_USER_INFO_COL_USERNAME, TABLE_USER_INFO_COL_PASSWORD)
	userInfo := &UserInfo{}
	err := mysql.DB.Raw(
		sql,
		username,
		password).First(userInfo).Error
	return userInfo, err
}

func (mysql *Mysql) GetUserInfoByIds(ids []string) ([]*UserInfo, error) {
	sql := fmt.Sprintf(`select * from %s where %s in ?`, TABLE_USER_INFO_NAME, TABLE_USER_INFO_COL_USERID)
	userInfos := []*UserInfo{}
	err := mysql.DB.Raw(
		sql, ids).Scan(&userInfos).Error
	return userInfos, err
}
