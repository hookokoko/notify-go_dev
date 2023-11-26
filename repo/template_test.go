package repo

import (
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func Test_Migrate(t *testing.T) {
	dsn := "root:@tcp(127.0.0.1:3306)/message_go?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	err = db.AutoMigrate(&MessageTemplate{})
	if err != nil {
		t.Error(err)
	}
}

func Test_CURD(t *testing.T) {
	dsn := "root:@tcp(127.0.0.1:3306)/message_go?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	dao := NewTemplateDao(db)

	tests := []struct {
		name string
		exec func()
	}{
		{
			name: "insert",
			exec: func() {
				err = dao.Insert(context.TODO(), &MessageTemplate{
					Name:           "测试邮件模版",
					AuditStatus:    20,
					IDType:         50,
					SendChannel:    40,
					TemplateType:   10,
					MsgType:        10,
					ExpectPushTime: "0",
					MsgContent:     "${name}先生/女士,你好，这是一条短信内容哦",
					Creator:        "chenhaokun",
					Updater:        "chenhaokun",
					Auditor:        "chenhaokun",
				})
				if err != nil {
					t.Error(err)
				}
			},
		},
		{
			name: "update",
			exec: func() {
				err := dao.Update(context.TODO(), &MessageTemplate{ID: 2, MsgContent: "${name}先生/女士,你好，这是一条邮件内容哦"})
				if err != nil {
					t.Error(err)
				}
			},
		},
		{
			name: "delete",
			exec: func() {
				err := dao.Delete(context.TODO(), 1)
				if err != nil {
					t.Error(err)
				}
			},
		},
		{
			name: "find",
			exec: func() {
				res, err := dao.FindById(context.TODO(), 2)
				if err != nil {
					t.Error(err)
				}
				t.Log(res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.exec()
		})
	}
}
