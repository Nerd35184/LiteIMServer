package repository

import "fmt"

const (
	TBL_ADD_CONTACT_REQUEST = "add_contact_request"

	TBL_ADD_CONTACT_REQUEST_COL_RECEIVER = "receiver"
	TBL_ADD_CONTACT_REQUEST_COL_STATUS   = "status"
	TBL_ADD_CONTACT_REQUEST_COL_REQ_ID   = "req_id"

	ADD_CONTACT_REQUEST_STATUS_PENDING = 1
	ADD_CONTACT_REQUEST_STATUS_ACCEPT  = 2
	ADD_CONTACT_REQUEST_STATUS_REJECT  = 3
)

type AddContactRequest struct {
	Id        int64
	ReqId     string `gorm:"type:varchar(255)"`
	Sender    string `gorm:"type:varchar(255)"`
	Receiver  string `gorm:"type:varchar(255)"`
	Status    int
	CreatedAt int64
	UpdatedAt int64
}

func (AddContactRequest AddContactRequest) TableName() string {
	return TBL_ADD_CONTACT_REQUEST
}

func (AddContactRequest AddContactRequest) Indexes() []*MySqlIndex {
	return []*MySqlIndex{}
}

func (mysql *Mysql) AddAddContactRequest(addContactRequest *AddContactRequest) error {
	return mysql.DB.Create(addContactRequest).Error
}

func (mysql *Mysql) CountAddContactRequestByReceiverId(receiverId string) (int64, error) {
	result := &CountModel{}
	sql := fmt.Sprintf(`select count(1) as count from %s where %s = ?`, TBL_ADD_CONTACT_REQUEST, TBL_ADD_CONTACT_REQUEST_COL_RECEIVER)
	err := mysql.DB.Raw(sql, receiverId).Scan(&result).Error
	return result.Count, err
}

func (mysql *Mysql) GetAddContactRequestByReceiverId(receiverId string, offset int, limit int) ([]AddContactRequest, error) {
	sql := fmt.Sprintf(`select * from %s where %s = ? limit ?,?`, TBL_ADD_CONTACT_REQUEST, TBL_ADD_CONTACT_REQUEST_COL_RECEIVER)
	result := []AddContactRequest{}
	err := mysql.DB.Raw(sql, receiverId, offset, limit).Scan(&result).Error
	return result, err
}

func (mysql *Mysql) UpdateAddContactRequestStatus(ReqId string, receiverId string, Status int) error {
	sql := fmt.Sprintf(`update %s set %s = ? where %s = ? and %s =?`, TBL_ADD_CONTACT_REQUEST, TBL_ADD_CONTACT_REQUEST_COL_STATUS, TBL_ADD_CONTACT_REQUEST_COL_REQ_ID, TBL_ADD_CONTACT_REQUEST_COL_RECEIVER)
	return mysql.DB.Exec(sql, Status, ReqId, receiverId).Error
}
