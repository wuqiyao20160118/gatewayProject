package dao

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"src/gatewayProject/dto"
	"src/gatewayProject/public"
	"time"
)

/*
	Demo Example:
	type Area struct {
		Id        int       `json:"id" gorm:"primary_key" description:"自增主键"`
		AreaName  string    `json:"area_name" gorm:"column:area_name" description:"区域名称"`
		CityId    int       `json:"city_id" gorm:"column:city_id" description:"城市id"`
		UserId    int64     `json:"user_id" gorm:"column:user_id" description:"操作人"`
		UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
		CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	}

	func (t *Area) TableName() string {
		return "area"
	}

	func (t *Area) Find(c *gin.Context, tx *gorm.DB, id string) (*Area, error) {
		area:=&Area{}
		err := tx.SetCtx(public.GetGinTraceContext(c)).Where("id = ?", id).Find(area).Error
		if err != nil {
			return nil, err
		}
		return area, nil
	}
*/

type Admin struct {
	Id        int       `json:"id" gorm:"primary_key" description:"primary key"`
	UserName  string    `json:"user_name" gorm:"column:user_name" description:"username"`
	Salt      string    `json:"salt" gorm:"column:salt" description:"salt"`
	Password  string    `json:"password" gorm:"column:password" description:"password"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"update time"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"create time"`
	IsDelete  int       `json:"is_delete" gorm:"column:is_delete" description:"is_delete"`
}

func (t *Admin) TableName() string {
	return "gateway_admin"
}

func (t *Admin) Find(c *gin.Context, tx *gorm.DB, search *Admin) (*Admin, error) {
	admin := &Admin{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(admin).Error
	if err != nil {
		return nil, err
	}
	return admin, nil
}

func (t *Admin) LoginAndCheck(c *gin.Context, tx *gorm.DB, param *dto.AdminLoginInput) (*Admin, error) {
	// Step 1: params.UserName 取得管理员信息 adminInfo
	adminInfo, err := t.Find(c, tx, &Admin{UserName: param.UserName, IsDelete: 0})
	if err != nil {
		return nil, errors.New("UserName does not exist.")
	}

	// Step 2: adminInfo.salt + params.Password sha256 => saltPassword
	saltPassword := public.GenSaltPassword(adminInfo.Salt, param.Password)
	if adminInfo.Password != saltPassword {
		return nil, errors.New("Password error, please retry.")
	}

	return adminInfo, nil
}

func (t *Admin) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}
