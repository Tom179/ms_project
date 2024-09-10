package project

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"strconv"
	"test.com/project-api/api/rpc"
	common "test.com/project-common"
	"test.com/project-common/errs"
	project "test.com/project-grpc/project"
	"time"
)

type ProjectHandler struct {
}

func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{}
}

func (p *ProjectHandler) index(c *gin.Context) {
	idStr, _ := c.Get("memberId")
	//fmt.Println(reflect.TypeOf(idStr))
	id := strconv.FormatInt(idStr.(int64), 10)
	result := &common.Result{}

	rsp, err := rpc.ProjectServiceClient.Index(context.Background(), &project.IndexMessage{Token: id})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	req := []*Menu{}

	err = copier.Copy(&req, rsp.Menus)
	if err != nil {
		fmt.Println("projectIndex拷贝失败:", err)
		return
	}
	fmt.Println("拷贝成功")

	c.JSON(http.StatusOK, result.Success(req))
}

type Menu struct {
	Id         int64   `json:"id"`
	Pid        int64   `json:"pid"`
	Title      string  `json:"title"`
	Icon       string  `json:"icon"`
	Url        string  `json:"url"`
	FilePath   string  `json:"file_path"`
	Params     string  `json:"params"`
	Node       string  `json:"node"`
	Sort       int32   `json:"sort"`
	Status     int32   `json:"status"`
	CreateBy   int64   `json:"create_by"`
	IsInner    int32   `json:"is_inner"`
	Values     string  `json:"values"`
	ShowSlider int32   `json:"show_slider"`
	StatusText string  `json:"statusText"`
	InnerText  string  `json:"innerText"`
	FullUrl    string  `json:"fullUrl"`
	Children   []*Menu `json:"children"`
}

type Page struct {
	Page     int64  `json:"page" form:"page"`
	PageSize int64  `json:"pageSize" form:"pageSize"`
	SelectBy string `json:"SelectBy" form:"SelectBy"`
}

func InitPageForm(c *gin.Context) (int64, int64) {
	page, _ := strconv.ParseInt(c.PostForm("page"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.PostForm("pageSize"), 10, 64)
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (p *ProjectHandler) MyProject(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	page, pageSize := InitPageForm(c)

	msg := &project.ProjectRpcMessage{
		MemberId:   c.GetInt64("memberId"),
		MemberName: c.GetString("memberName"),
		Page:       page,
		PageSize:   pageSize,
		SelectBy:   c.PostForm("selectBy"),
	}

	fmt.Println("查询类型selectedBy为", msg.SelectBy)
	rsp, err := rpc.ProjectServiceClient.FindProjectByMemId(ctx, msg)

	result := common.Result{}
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	pms := []*ProAndMember{}
	copier.Copy(&pms, rsp.Pm)
	/*
		fmt.Println("grpc返回的响应为", rsp.Pm)
		fmt.Println("api返回的响应为", myProjectList)*/
	if pms == nil {
		pms = []*ProAndMember{}
	}

	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  pms,
		"total": rsp.Total,
	}))
}

func (p *ProjectHandler) ProjectTemplate(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	page, pageSize := InitPageForm(c)
	viewType, _ := strconv.ParseInt(c.PostForm("viewType"), 10, 32)

	msg := &project.ProjectRpcMessage{
		MemberId:   c.GetInt64("memberId"),
		MemberName: c.GetString("memberName"),
		Page:       page,
		PageSize:   pageSize,
		SelectBy:   c.PostForm("selectBy"),
		ViewType:   int32(viewType),
	}
	rpcRsp, err := rpc.ProjectServiceClient.FindProjectTemplate(ctx, msg)
	fmt.Println("响应为", rpcRsp)

	result := common.Result{}
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	myPeojrctList := []*ProjectTemplate{}
	//rpcRsp.Ptm
	copier.Copy(&myPeojrctList, rpcRsp.Ptm)

	fmt.Println("grpc返回的响应为", rpcRsp.Ptm) //空指针了
	//fmt.Println("api返回的响应为", myProjectList)

	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  myPeojrctList,
		"total": rpcRsp.Total,
	}))
}

type Project struct {
	Id                 int64   `json:"id"`
	Cover              string  `json:"cover"`
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	AccessControlType  string  `json:"access_control_type"`
	WhiteList          string  `json:"white_list"`
	Order              int     `json:"order"`
	Deleted            int     `json:"deleted"`
	TemplateCode       string  `json:"template_code"`
	Schedule           float64 `json:"schedule"`
	CreateTime         string  `json:"create_time"`
	OrganizationCode   string  `json:"organization_code"`
	DeletedTime        string  `json:"deleted_time"`
	Private            int     `json:"private"`
	Prefix             string  `json:"prefix"`
	OpenPrefix         int     `json:"open_prefix"`
	Archive            int     `json:"archive"`
	ArchiveTime        int64   `json:"archive_time"`
	OpenBeginTime      int     `json:"open_begin_time"`
	OpenTaskPrivate    int     `json:"open_task_private"`
	TaskBoardTheme     string  `json:"task_board_theme"`
	BeginTime          int64   `json:"begin_time"`
	EndTime            int64   `json:"end_time"`
	AutoUpdateSchedule int     `json:"auto_update_schedule"`
	Code               string  `json:"code"`
}

type MemberProject struct {
	Id          int64  `json:"id"`
	ProjectCode int64  `json:"project_code"`
	MemberCode  int64  `json:"member_code"`
	JoinTime    string `json:"join_time"`
	IsOwner     int64  `json:"is_owner"`
	Authorize   string `json:"authorize"`
}
type ProAndMember struct {
	Project
	ProjectCode int64  `json:"project_code"`
	MemberCode  int64  `json:"member_code"`
	JoinTime    int64  `json:"join_time"`
	IsOwner     int64  `json:"is_owner"`
	Authorize   string `json:"authorize"`
	OwnerName   string `json:"owner_name"`
	Collected   int    `json:"collected"`
}
type ProjectTemplate struct {
	Id               int                   `json:"id"`
	Name             string                `json:"name"`
	Description      string                `json:"description"`
	Sort             int                   `json:"sort"`
	CreateTime       string                `json:"create_time"`
	OrganizationCode string                `json:"organization_code"`
	Cover            string                `json:"cover"`
	MemberCode       string                `json:"member_code"`
	IsSystem         int                   `json:"is_system"`
	TaskStages       []*TaskStagesOnlyName `json:"task_stages"`
	Code             string                `json:"code"`
}

type TaskStagesOnlyName struct {
	Name string `json:"name"`
}
