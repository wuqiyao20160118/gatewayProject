package dao

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"src/gatewayProject/dto"
	"src/gatewayProject/public"
	"time"
)

type App struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	AppID     string    `json:"app_id" gorm:"column:app_id" description:"租户id	"`
	Name      string    `json:"name" gorm:"column:name" description:"租户名称	"`
	Secret    string    `json:"secret" gorm:"column:secret" description:"密钥"`
	WhiteIPS  string    `json:"white_ips" gorm:"column:white_ips" description:"ip白名单，支持前缀匹配"`
	Qpd       int64     `json:"qpd" gorm:"column:qpd" description:"日请求量限制"`
	Qps       int64     `json:"qps" gorm:"column:qps" description:"每秒请求量限制"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"添加时间	"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	IsDelete  int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
}

func (t *App) TableName() string {
	return "gateway_app"
}

func (t *App) APPList(c *gin.Context, tx *gorm.DB, param *dto.APPListInput) ([]App, int64, error) {
	total := int64(0)
	var list []App
	offset := (param.PageNo - 1) * param.PageSize

	query := tx.SetCtx(public.GetGinTraceContext(c))  // 控制台可以打印数据库查询
	query = query.Table(t.TableName()).Where("is_delete=0")  // 这里需要知道表的字段，故需要query.Table()
	if param.Info != "" {
		query = query.Where("(app_id like ? or name like ?)", "%"+param.Info+"%", "%"+param.Info+"%")
	}

	// 设置分页条数以及偏移量
	err := query.Limit(param.PageSize).Offset(offset).Order("id desc").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	errCount := query.Count(&total).Error
	if errCount != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (t *App) Find(c *gin.Context, tx *gorm.DB, search *App) (*App, error) {
	out := &App{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *App) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}
