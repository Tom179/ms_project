package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/resolver"
)

type Server struct { //要注册的服务
	Name    string `json:"name"`
	Addr    string `json:"addr"`    //服务地址
	Version string `json:"version"` //服务版本
	Weight  int64  `json:"weight"`  //服务权重
}

func BuildPrefix(info Server) string { //返回:服务名-版本
	if info.Version == "" {
		return fmt.Sprintf("%s", info.Name)
	}
	return fmt.Sprintf("%s%s", info.Name, info.Version)
}

func BuildRegPath(info Server) string { //返回:服务名-版本 服务地址
	return fmt.Sprintf("%s%s", BuildPrefix(info), info.Addr)
}

func ParseValue(value []byte) (Server, error) {
	info := Server{}
	if err := json.Unmarshal(value, &info); err != nil {
		return info, err
	}
	return info, nil
} //从字节切片中解析出服务实例

func SplitPath(path string) (Server, error) {
	info := Server{}
	strs := strings.Split(path, "/")
	if len(strs) == 0 {
		return info, errors.New("invalid path")
	}
	info.Addr = strs[len(strs)-1] //拆分出服务地址返回一个新的服务实例（只实例化了地址）
	return info, nil
}

// Exist helper function
func Exist(l []resolver.Address, addr resolver.Address) bool {
	for i := range l {
		if l[i].Addr == addr.Addr {
			return true
		}
	}
	return false
}

// Remove helper function
func Remove(s []resolver.Address, addr resolver.Address) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr.Addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}

func BuildResolverUrl(app string) string {
	return schema + ":///" + app
}
