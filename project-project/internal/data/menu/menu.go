package menu

type ProjectMenu struct {
	Id         int64
	Pid        int64
	Title      string
	Icon       string
	Url        string
	FilePath   string
	Params     string
	Node       string
	Sort       int
	Status     int
	CreateBy   int64
	IsInner    int
	Values     string
	ShowSlider int
}

func (*ProjectMenu) TableName() string {
	return "ms_project_menu"
}

type ProjectMenuChild struct {
	ProjectMenu
	StatusText string
	InnerText  string
	FullUrl    string
	Children   []*ProjectMenuChild
}

func GetChild(pms []*ProjectMenu) []*ProjectMenuChild { //填入子项目的信息。如何判断是否有子项目

	projectMap := make(map[int64]*ProjectMenuChild)
	for _, p := range pms {
		statusText := getStatus(p.Status)
		innerText := getInnerText(p.IsInner)
		fullUrl := getFullUrl(p.Url, p.Params, p.Values)
		projectMap[p.Id] = &ProjectMenuChild{*p, statusText, innerText, fullUrl, []*ProjectMenuChild{}} //为每一项初始化一个空的child列表
	}

	for _, p := range pms {
		if p.Pid > 0 { //如果说有父级：填充一下父级的child列表
			statusText := getStatus(p.Status)
			innerText := getInnerText(p.IsInner)
			fullUrl := getFullUrl(p.Url, p.Params, p.Values)
			projectMap[p.Pid].Children = append(projectMap[p.Pid].Children, &ProjectMenuChild{*p, statusText, innerText, fullUrl, []*ProjectMenuChild{}}) //想一想是否有问题，是否存在空数据的覆盖？
		}
	}
	//ProjectMap确实记录了完整的信息，但是只平铺记录了 某一个projectMenu和它下一层的child信息。
	//还差一步，就是如何构造层层嵌套的[]ProjectMenuChild结果？

	result := []*ProjectMenuChild{}
	for _, p := range pms {
		if p.Pid == 0 { //0级项目
			result = append(result, projectMap[p.Id])
		}
	}

	for i := 0; i < len(result); i++ {
		if result[i].Children != nil {
			for j := 0; j < len(result[i].Children); j++ {
				checkAndFill_Child(result[i].Children[j], projectMap)
			}
		}
	}

	return result
}

func checkAndFill_Child(cur *ProjectMenuChild, projectmap map[int64]*ProjectMenuChild) { //填充单个ProjectMenuChild
	if projectmap[cur.Id].Children == nil { //最后一级
		return
	}
	//通过所有数据的配置中心projectMap来判断是否存在下级，而不是正在填充的result本身。
	//如果存在下级，先填充下级。然后在判断下级是否存在下级
	cur.Children = projectmap[cur.Id].Children
	for i := 0; i < len(cur.Children); i++ {
		checkAndFill_Child(cur.Children[i], projectmap)
	}
}

func getFullUrl(url string, params string, values string) string {
	if (params != "" && values != "") || values != "" {
		return url + "/" + values
	}
	return url
}

func getInnerText(inner int) string {
	if inner == 0 {
		return "导航"
	}
	if inner == 1 {
		return "内页"
	}
	return ""
}

func getStatus(status int) string {
	if status == 0 {
		return "禁用"
	}
	if status == 1 {
		return "使用中"
	}
	return ""
}
