package repository

const (
	TBL_DIALOGUE_SESS = "dialogue_sess"
)

//todo 用作群聊，此处考虑可能不需要这个表，直接把群这个概念归并入user
type DialogueSess struct {
	Id     int64
	Type   int
	SessId string `gorm:"type:varchar(255)"`
}

func (DialogueSess DialogueSess) TableName() string {
	return TBL_DIALOGUE_SESS
}

func (DialogueSess DialogueSess) Indexes() []*MySqlIndex {
	return []*MySqlIndex{}
}

func (mysql *Mysql) AddDialogueSess(dialogueSess *DialogueSess) error {
	return mysql.DB.Create(dialogueSess).Error
}
