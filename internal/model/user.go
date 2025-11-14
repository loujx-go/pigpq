package model

import (
	. "pigpq/internal/model/base"
	utils "pigpq/internal/pkg/untils"
)

type User struct {
	BaseModel
	Username               string           `gorm:"column:username;type:varchar(255);not null" json:"username"`                                                          // 用户名
	Nickname               string           `gorm:"column:nickname;type:varchar(255);not null" json:"nickname"`                                                          // 昵称
	Phone                  string           `gorm:"column:phone;type:varchar(15);not null;uniqueIndex:uniq_phone,where:deleted_at is null;index:idx_phone" json:"phone"` // 手机号
	Avatar                 string           `gorm:"column:avatar;type:varchar(255);not null" json:"avatar"`                                                              // 头像
	BgImage                string           `gorm:"column:bg_image;type:varchar(255);default:''" json:"bg_image"`                                                        // 背景图片
	Gender                 int8             `gorm:"column:gender;type:tinyint;default:1;comment:性别|radio|1:男;2:女;3:未知" json:"gender"`                                    // 性别|radio|1:男;2:女;3:未知
	Birthday               utils.FormatDate `gorm:"column:birthday;type:date" json:"birthday"`                                                                           // 生日
	Intro                  string           `gorm:"column:intro;type:varchar(255);default:''" json:"intro"`                                                              // 简介
	Status                 int8             `gorm:"column:status;type:tinyint;default:1;comment:状态|radio|1:启用;2:禁用;3:已注销;" json:"status"`                                // 状态|radio|1:启用;2:禁用;3:已注销;
	IsRobot                int8             `gorm:"column:is_robot;type:tinyint;default:0;comment:是否是机器人|radio|0:否;1:是;" json:"is_robot"`                                // 是否是机器人|radio|0:否;1:是;
	WriteOffAt             utils.FormatDate `gorm:"column:write_off_at" json:"write_off_at"`                                                                             // 注销时间
	Email                  string           `gorm:"column:email;type:varchar(255);not null" json:"email"`                                                                // 邮箱
	Province               string           `gorm:"column:province;type:varchar(100);default:''" json:"province"`                                                        // 省
	ProvinceAdCode         string           `gorm:"column:province_adcode;type:varchar(10);default:''" json:"province_adcode"`                                           // 省份编码（省市区表adcode）
	City                   string           `gorm:"column:city;type:varchar(100);default:''" json:"city"`                                                                // 市
	CityAdCode             string           `gorm:"column:city_adcode;type:varchar(10);default:''" json:"city_adcode"`                                                   // 市编码（省市区表adcode）
	Area                   string           `gorm:"column:area;type:varchar(100);default:''" json:"area"`                                                                // 区
	AreaAdCode             string           `gorm:"column:area_adcode;type:varchar(10);default:''" json:"area_adcode"`                                                   // 区编码（省市区表adcode）
	Address                string           `gorm:"column:address;type:varchar(255);default:''" json:"address"`                                                          // 详细地址
	Password               string           `gorm:"column:password;type:varchar(255);not null" json:"-"`                                                                 // 密码（敏感字段，JSON 不输出）
	LastEditNicknameTime   uint             `gorm:"column:last_edit_nickname_time;type:int unsigned;default:0" json:"last_edit_nickname_time"`                           // 最后一次修改nickname的时间（时间戳）
	FollowNumber           uint             `gorm:"column:follow_number;type:int unsigned;default:0" json:"follow_number"`                                               // 关注人数
	FansNumber             uint             `gorm:"column:fans_number;type:int unsigned;default:0" json:"fans_number"`                                                   // 粉丝人数
	NoteNumber             uint             `gorm:"column:note_number;type:int unsigned;default:0" json:"note_number"`                                                   // 笔记数
	InterestTagIDs         string           `gorm:"column:interest_tag_ids;type:varchar(255);default:''" json:"interest_tag_ids"`                                        // 兴趣标签ID，多个以英文逗号分割
	TagIDs                 string           `gorm:"column:tag_ids;type:varchar(255);default:''" json:"tag_ids"`                                                          // 用户标签id，多个以英文逗号分割
	FromSource             uint8            `gorm:"column:from_source;type:tinyint unsigned;default:1;comment:来源|radio|1:自然注册2：推荐分享" json:"from_source"`                 // 来源|radio|1:自然注册2：推荐分享
	RegisterSource         uint8            `gorm:"column:register_source;type:tinyint unsigned;default:1;comment:注册来源|radio|1：App" json:"register_source"`              // 注册来源|radio|1：App
	ParentID               int64            `gorm:"column:parent_id;default:0" json:"parent_id"`                                                                         // 推荐人id|input
	ShareBindingTime       utils.FormatDate `gorm:"column:share_binding_time" json:"share_binding_time"`                                                                 // 绑定推荐人时间|input
	LastLoginAt            utils.FormatDate `gorm:"column:last_login_at" json:"last_login_at"`                                                                           // 最后登录时间|input
	LastLoginAppVersion    string           `gorm:"column:last_login_app_version;type:varchar(20);default:''" json:"last_login_app_version"`                             // 最后登录的App版本|input
	LastLoginDeviceModelOS string           `gorm:"column:last_login_device_model_os;type:varchar(50);default:''" json:"last_login_device_model_os"`                     // 最后登录的设备型号和操作系统|input
	ShareNo                string           `gorm:"column:share_no;type:varchar(24);default:'';index:idx_share_no" json:"share_no"`                                      // 用户邀请码|input
}

func NewUser() *User {
	return &User{}
}

func (u *User) GetUserByPhone(phone string) *User {
	if u.DB().Where("phone = ?", phone).Find(u).Error != nil {
		return nil
	}
	return u
}

func (u *User) GetUserById(id uint) *User {
	if u.DB().Where("id = ?", id).Find(u).Error != nil {
		return nil
	}
	return u
}
