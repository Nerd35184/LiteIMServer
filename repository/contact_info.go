package repository

import "fmt"

const (
	TBL_CONTACT_INFO = "contact_info"

	TBL_CONTACT_INFO_CONTACT_TYPE_SYSTEM_USER  = 1
	TBL_CONTACT_INFO_CONTACT_TYPE_REGULAR_USER = 2
	TBL_CONTACT_INFO_CONTACT_TYPE_GROUP        = 3
)

type ContactInfo struct {
	Id          int64
	UserId      string `gorm:"type:varchar(255)"`
	ContactId   string `gorm:"type:varchar(255)"`
	ContactType int
}

func (ContactInfo ContactInfo) TableName() string {
	return TBL_CONTACT_INFO
}

func (ContactInfo ContactInfo) Indexes() []*MySqlIndex {
	return []*MySqlIndex{}
}

func (mysql *Mysql) AddContactInfo(contactInfo *ContactInfo) error {
	return mysql.DB.Create(contactInfo).Error
}

func (mysql *Mysql) CountUserContactByType(UserId string, ContactType int) (int64, error) {
	result := &CountModel{}
	sql := fmt.Sprintf(`select count(1) as count from %s where user_id = ? and contact_type = ?`, TBL_CONTACT_INFO)
	err := mysql.DB.Raw(sql, UserId, ContactType).Scan(&result).Error
	return result.Count, err
}

func (mysql *Mysql) GetUserContactsByType(UserId string, ContactType int, offset int, limit int) ([]*ContactInfo, error) {
	result := []*ContactInfo{}
	sql := fmt.Sprintf(`select * from %s where user_id = ? and contact_type = ? limit ?,?`, TBL_CONTACT_INFO)
	err := mysql.DB.Raw(sql, UserId, ContactType, offset, limit).Scan(&result).Error
	return result, err
}

func (mysql *Mysql) GetUserContact(UserId string, ContactType int, ContactId string) (*ContactInfo, error) {
	result := &ContactInfo{}
	sql := fmt.Sprintf(`select * from %s where user_id = ? and contact_type = ? and contact_id = ?`, TBL_CONTACT_INFO)
	err := mysql.DB.Raw(sql, UserId, ContactType, ContactId).First(&result).Error
	return result, err
}

func (mysql *Mysql) RemoveContact(UserId string, ContactId string) error {
	sql := fmt.Sprintf(`delete from %s where user_id = ? and contact_id = ?`, TBL_CONTACT_INFO)
	err := mysql.DB.Exec(sql, UserId, ContactId).Error
	return err
}
