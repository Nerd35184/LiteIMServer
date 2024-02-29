package repository

import (
	"fmt"
	"server/conf"
	"testing"
)

func TestInitData(t *testing.T) {
	MysqlConfig := &conf.MysqlConfig{
		Username: "root",
		Password: "123456",
		Host:     "127.0.0.1",
		Port:     3306,
		DbName:   "lite_im",
	}
	mysql, err := InitMysqlDb(MysqlConfig)
	if err != nil {
		panic(err)
	}
	userInfos := []*UserInfo{}

	for i := 0; i < 10; i++ {
		userInfo := &UserInfo{
			UserId:    fmt.Sprintf("%d", i),
			Username:  fmt.Sprintf("username%d", i),
			Password:  fmt.Sprintf("password%d", i),
			Nickname:  fmt.Sprintf("nickname%d", i),
			Signature: fmt.Sprintf("signature%d", i),
			Avatar:    fmt.Sprintf("http://127.0.0.1:8080/static/download/%d.jpg", i%2+1),
		}
		userInfos = append(userInfos, userInfo)
	}
	for _, useruserInfo := range userInfos {
		mysql.AddUserInfo(useruserInfo)
	}

	contactInfos := []*ContactInfo{}

	for i := 2; i < 10; i++ {
		contactInfo := &ContactInfo{
			UserId:      "1",
			ContactId:   fmt.Sprintf("%d", i),
			ContactType: TBL_CONTACT_INFO_CONTACT_TYPE_REGULAR_USER,
		}
		contactInfos = append(contactInfos, contactInfo)
	}

	for _, contactInfo := range contactInfos {
		mysql.AddContactInfo(contactInfo)
	}
}
