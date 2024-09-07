package project

// 【项目表】
type Project struct {
	Id                 int64
	Cover              string
	Name               string
	Description        string
	AccessControlType  int
	WhiteList          string
	Sort               int
	Deleted            int
	TemplateCode       string
	Schedule           float64
	CreateTime         int64
	OrganizationCode   int64
	DeletedTime        string
	Private            int
	Prefix             string
	OpenPrefix         int
	Archive            int
	ArchiveTime        int64
	OpenBeginTime      int
	OpenTaskPrivate    int
	TaskBoardTheme     string
	BeginTime          int64
	EndTime            int64
	AutoUpdateSchedule int
}

func (*Project) TableName() string {
	return "ms_project"
}

type MemberProject struct {
	Id          int64
	ProjectCode int64
	MemberCode  int64
	JoinTime    int64
	IsOwner     int64
	Authorize   string
}

func (*MemberProject) TableName() string {
	return "ms_project_member"
}

type ProjectAndMenber struct {
	Project
	ProjectCode int64
	MemberCode  int64
	JoinTime    int64
	IsOwner     int64
	Authorize   string
}

func (m *ProjectAndMenber) GetAccessControlType() string {
	if m.AccessControlType == 0 {
		return "open"
	}
	if m.AccessControlType == 1 {
		return "private"
	}
	if m.AccessControlType == 2 {
		return "custom"
	}
	return ""
}

func ToMap(orgs []*ProjectAndMenber) map[int64]*ProjectAndMenber {
	m := make(map[int64]*ProjectAndMenber)
	for _, v := range orgs {
		m[v.Id] = v
	}
	return m
}
