package repo

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type ChannelAccount struct {
	ID            int                   `gorm:"primaryKey" json:"id"`
	Name          string                `gorm:"column:name;not null;comment:账号名称" json:"name"`                     // 账号名称
	SendChannel   string                `gorm:"column:send_channel;not null;comment:发送渠道" json:"send_channel"`     // 发送渠道
	AccountConfig string                `gorm:"column:account_config;not null;comment:账户配置" json:"account_config"` // 账户配置
	Creator       string                `gorm:"column:creator;not null;comment:账号拥有者" json:"creator"`              //账号拥有者
	Created       int64                 `gorm:"column:created;autoCreateTime;comment:创建时间" json:"created"`         // 创建时间,单位 s
	Updated       int64                 `gorm:"column:updated;autoCreateTime;comment:更新时间" json:"updated"`
	IsDeleted     soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
	DeletedAt     int64
}

// TableName SendAccount's table name
func (*ChannelAccount) TableName() string {
	return "t_channel_account"
}

type channelAccountDao struct {
	db *gorm.DB
}

type IChannelAccountDao interface {
	Insert(ctx context.Context, account *ChannelAccount) error
	Update(ctx context.Context, account *ChannelAccount) error
	Delete(ctx context.Context, id int) error
	FindById(ctx context.Context, id int) (*ChannelAccount, error)
	GetAccountConfigByScriptName(ctx context.Context, scriptName string) (*ChannelAccount, error)
}

func NewChannelAccountDao(db *gorm.DB) IChannelAccountDao {
	return &channelAccountDao{db: db}
}

func (c *channelAccountDao) Insert(ctx context.Context, account *ChannelAccount) error {
	//TODO implement me
	panic("implement me")
}

func (c *channelAccountDao) Update(ctx context.Context, account *ChannelAccount) error {
	//TODO implement me
	panic("implement me")
}

func (c *channelAccountDao) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func (c *channelAccountDao) FindById(ctx context.Context, id int) (*ChannelAccount, error) {
	//TODO implement me
	panic("implement me")
}

func (c *channelAccountDao) GetAccountConfigByScriptName(ctx context.Context, scriptName string) (*ChannelAccount, error) {
	//TODO implement me
	panic("implement me")
}
