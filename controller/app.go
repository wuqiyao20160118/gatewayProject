package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"src/gatewayProject/dao"
	"src/gatewayProject/dto"
	"src/gatewayProject/golang_common/lib"
	"src/gatewayProject/middleware"
	"src/gatewayProject/public"
	"time"
)

type APPController struct {
}

func APPRegister(group *gin.RouterGroup) {
	app := &APPController{}
	group.GET("/app_list", app.APPList)
	group.GET("/app_delete", app.APPDelete)
	group.GET("/app_detail", app.APPDetail)
	group.GET("/app_stat", app.APPStat)
	group.POST("/app_add", app.APPAdd)
	group.POST("/app_update", app.APPUpdate)
}

// APPList godoc
// @Summary APP List
// @Description APP List
// @Tags APP Interface
// @ID /app/app_list
// @Accept  json
// @Param info query string false "keyword"
// @Param page_no query int true "pageNum"
// @Param page_size query int false "pageSize"
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.APPListOutput} "success"
// @Router /app/app_list [get]
func (app *APPController) APPList(c *gin.Context) {
	params := &dto.APPListInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	appInfo := &dao.App{}
	tx, err := lib.GetGormPool("default") // 使用配置文件中default的数据库连接池
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	list, total, err := appInfo.APPList(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	var outList []dto.APPListItemOutput
	for _, listItem := range list {
		//appCounter, err := public.FlowCounterHandler.GetCounter(public.FlowAppPrefix + listItem.AppID)
		//if err != nil {
		//	middleware.ResponseError(c, 2003, err)
		//	c.Abort()
		//	return
		//}

		outItem := dto.APPListItemOutput{
			ID:       listItem.ID,
			AppID:    listItem.AppID,
			Name:     listItem.Name,
			Secret:   listItem.Secret,
			WhiteIPS: listItem.WhiteIPS,
			Qpd:      listItem.Qpd,
			Qps:      listItem.Qps,
			RealQpd:  0, //appCounter.TotalCount,
			RealQps:  0, //appCounter.QPS,
		}
		outList = append(outList, outItem)
	}

	out := &dto.APPListOutput{
		Total: total,
		List:  outList,
	}

	middleware.ResponseSuccess(c, out)
}

// APPDetail godoc
// @Summary APP Detail
// @Description APP Detail
// @Tags APP Interface
// @ID /app/app_detail
// @Accept  json
// @Produce  json
// @Param id query string true "APP ID"
// @Success 200 {object} middleware.Response{data=dao.App} "success"
// @Router /app/app_detail [get]
func (app *APPController) APPDetail(c *gin.Context) {
	params := &dto.APPDetailInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 读取App基本信息
	detail := &dao.App{ID: params.ID}
	detail, err = detail.Find(c, tx, detail)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	middleware.ResponseSuccess(c, detail)
}

// APPDelete godoc
// @Summary APP Delete
// @Description APP Delete
// @Tags APP Interface
// @ID /app/app_delete
// @Accept  json
// @Produce  json
// @Param id query string true "APP ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_delete [get]
func (app *APPController) APPDelete(c *gin.Context) {
	params := &dto.APPDetailInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	//读取基本信息
	detail := &dao.App{ID: params.ID}
	detail, err = detail.Find(c, tx, detail)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	detail.IsDelete = 1
	if err := detail.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")
}

// APPAdd godoc
// @Summary Add APP User Service
// @Description Add APP User Service
// @Tags APP Interface
// @ID /app/app_add
// @Accept  json
// @Param body body dto.APPAddHttpInput true "body"
// @Produce  json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_add [post]
func (app *APPController) APPAdd(c *gin.Context) {
	params := &dto.APPAddHttpInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 开启事务
	tx = tx.Begin()

	//验证app_id是否被占用
	search := &dao.App{AppID: params.AppID}
	if _, err := search.Find(c, tx, search); err == nil {
		// 错误回滚
		tx.Rollback()
		middleware.ResponseError(c, 2002, errors.New("租户ID被占用，请重新输入"))
		return
	}

	// 生成32密钥
	if params.Secret == "" {
		params.Secret = public.MD5(params.AppID)
	}

	// 插入表数据
	appModel := &dao.App{
		AppID:    params.AppID,
		Name:     params.Name,
		Secret:   params.Secret,
		WhiteIPS: params.WhiteIPS,
		Qps:      params.Qps,
		Qpd:      params.Qpd,
	}
	if err := appModel.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, err)
		return
	}

	tx.Commit()

	middleware.ResponseSuccess(c, "")
}

// APPUpdate godoc
// @Summary Update APP Service
// @Description Update APP Service
// @Tags APP Interface
// @ID /app/app_update
// @Accept  json
// @Param body body dto.APPUpdateHttpInput true "body"
// @Produce  json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_update [post]
func (app *APPController) APPUpdate(c *gin.Context) {
	params := &dto.APPUpdateHttpInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 开启事务
	tx = tx.Begin()

	//验证app_id是否存在
	search := &dao.App{AppID: params.AppID}
	info, err := search.Find(c, tx, search)
	if err != nil {
		// 错误回滚
		tx.Rollback()
		middleware.ResponseError(c, 2002, errors.New("租户ID不存在，请重新输入"))
		return
	}

	// 生成32密钥
	if params.Secret == "" {
		params.Secret = public.MD5(params.AppID)
	}

	info.Name = params.Name
	info.Secret = params.Secret
	info.WhiteIPS = params.WhiteIPS
	info.Qps = params.Qps
	info.Qpd = params.Qpd

	// 插入表数据
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, err)
		return
	}

	tx.Commit()

	middleware.ResponseSuccess(c, "")
}

// ServiceStat godoc
// @Summary ServiceStat
// @Description ServiceStat
// @Tags Service Interface
// @ID /service/service_stat
// @Accept  json
// @Produce  json
// @Param id query string true "service ID"
// @Success 200 {object} middleware.Response{data=dto.ServiceStatOutput} "success"
// @Router /service/service_stat [get]
func (app *APPController) APPStat(c *gin.Context) {
	params := &dto.APPDetailInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	//读取基本信息
	detail := &dao.App{ID: params.ID}
	detail, err = detail.Find(c, tx, detail)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	var todayList []int64
	currentTime := time.Now()
	for i := 0; i <= currentTime.Hour(); i++ {
		//dateTime := time.Date(currentTime.Year(),currentTime.Month(),currentTime.Day(),i,0,0,0,lib.TimeLocation)
		//hourData,_:=counter.GetHourData(dateTime)
		//todayList = append(todayList, hourData)
		todayList = append(todayList, 0)
	}

	var yesterdayList []int64
	//yesterTime:= currentTime.Add(-1*time.Duration(time.Hour*24))
	for i := 0; i <= 23; i++ {
		//dateTime := time.Date(yesterTime.Year(),yesterTime.Month(),yesterTime.Day(),i,0,0,0,lib.TimeLocation)
		//hourData,_:=counter.GetHourData(dateTime)
		//yesterdayList = append(yesterdayList, hourData)
		yesterdayList = append(yesterdayList, 0)
	}

	middleware.ResponseSuccess(c, &dto.StatisticsOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	})
}
