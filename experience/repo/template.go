package repo

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type MessageTemplate struct {
	ID             int                   `gorm:"column:id;primaryKey" json:"id"`
	Name           string                `gorm:"column:name;type:VARCHAR(255);comment:模板标题" json:"name"`                            // 模板标题
	AuditStatus    int                   `gorm:"column:audit_status;type:SMALLINT;comment:当前消息审核状态" json:"audit_status"`            // 当前消息审核状态: 10.待审核 20.审核成功 30.被拒绝
	FlowId         int                   `gorm:"column:flow_id;comment:工单ID（审核模板走工单）" json:"flow_id"`                               // 工单ID（审核模板走工单）
	MsgStatus      int                   `gorm:"column:msg_status;type:SMALLINT;comment:消息状态" json:"msg_status"`                    // 消息状态
	CronTaskId     int                   `gorm:"column:cron_task_id;comment:定时任务Id" json:"cron_task_id"`                            // 定时任务Id(由xxl-job返回)
	IDType         int                   `gorm:"column:id_type;type:SMALLINT;comment:消息的发送ID类型" json:"id_type"`                     // 消息的发送ID类型: 10. userId 20.did 30.手机号 40.openId 50.email 60.企业微信userId
	SendChannel    int                   `gorm:"column:send_channel;type:SMALLINT;comment:消息发送渠道" json:"send_channel"`              // 消息发送渠道: 10.IM 20.Push 30.短信 40.Email 50.公众号 60.小程序 70.企业微信
	TemplateType   int                   `gorm:"column:template_type;type:SMALLINT;comment:模板类型" json:"template_type"`              // 模板类型 10.运营类 20.技术类接口调用
	ShieldType     int                   `gorm:"column:shield_type;type:SMALLINT;comment:屏蔽类型" json:"shield_type"`                  // 屏蔽类型 10.夜间不屏蔽 20.夜间屏蔽 30.夜间屏蔽(次日早上9点发送)
	MsgType        int                   `gorm:"column:msg_type;type:SMALLINT;comment:消息类型" json:"msg_type"`                        // 消息类型 10.通知类消息 20.营销类消息 30.验证码类消息
	ExpectPushTime string                `gorm:"column:expect_push_time;type:VARCHAR(255);comment:推送消息的时间" json:"expect_push_time"` // 推送消息的时间, 0: 立即发送, else: crontab 表达式
	MsgContent     string                `gorm:"column:msg_content;type:LONGTEXT;comment:消息内容" json:"msg_content"`                  // 消息内容 占位符用${var}表示
	SendAccount    int                   `gorm:"column:send_account;comment:发送账号" json:"send_account"`                              // 发送账号 一个渠道下可存在多个账号
	Creator        string                `gorm:"column:creator;type:VARCHAR(255);comment:创建者" json:"creator"`                       // 创建者
	Updater        string                `gorm:"column:updater;type:VARCHAR(255);comment:更新者" json:"updater"`                       // 更新者
	Auditor        string                `gorm:"column:auditor;type:VARCHAR(255);comment:审核人" json:"auditor"`                       // 审核人
	Team           string                `gorm:"column:team;type:VARCHAR(255);comment:业务方团队" json:"team"`                           // 业务方团队
	Proposer       string                `gorm:"column:proposer;type:VARCHAR(255);comment:业务方" json:"proposer"`                     // 业务方
	Created        int64                 `gorm:"column:created;autoCreateTime;comment:创建时间" json:"created"`                         // 创建时间,单位 s
	Updated        int64                 `gorm:"column:updated;autoCreateTime;comment:更新时间" json:"updated"`
	IsDeleted      soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
	DeletedAt      int64
}

func (m MessageTemplate) TableName() string {
	return "t_message_template"
}

type templateDao struct {
	db *gorm.DB
}

type ITemplateDao interface {
	Insert(ctx context.Context, template *MessageTemplate) error
	Update(ctx context.Context, template *MessageTemplate) error
	Delete(ctx context.Context, id int) error
	FindById(ctx context.Context, id int) (*MessageTemplate, error)
}

func NewTemplateDao(db *gorm.DB) ITemplateDao {
	return &templateDao{db: db}
}

func (t *templateDao) Insert(ctx context.Context, template *MessageTemplate) error {
	return t.db.WithContext(ctx).Create(template).Error
}

func (t *templateDao) Update(ctx context.Context, template *MessageTemplate) error {
	return t.db.WithContext(ctx).Updates(template).Error
}

func (t *templateDao) Delete(ctx context.Context, id int) error {
	return t.db.WithContext(ctx).Delete(&MessageTemplate{}, id).Error
}

func (t *templateDao) FindById(ctx context.Context, id int) (*MessageTemplate, error) {
	var tpl MessageTemplate
	err := t.db.WithContext(ctx).Find(&tpl, "id =?", id).Error
	return &tpl, err
}
