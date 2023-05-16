package netutil

import (
	"log"
	"regexp"
	"strings"
)

func CheckExistSNat() map[string]bool {
	s := listIptablesNatTable("POSTROUTING")
	lines := strings.Split(s, "\n")
	mp := make(map[string]bool)
	for i := range lines {
		if !strings.Contains(lines[i], "MASQUERADE") {
			continue
		}
		snatstr := strings.Trim(lines[i], " ")
		s := delete_extra_space(snatstr)
		split := strings.Split(s, " ")
		if len(split) < 10 {
			continue
		}
		mp[split[9]] = true
	}
	return mp
}

func listIptablesNatTable(s string) string {
	//iptables -t nat -nvL xxx --line
	lines := ExecCmd("iptables", "-t", "nat", "-nvL", s, "--line")
	log.Printf("iptables nat table[%s]:\t %s", s, lines)
	return lines
}

/*
函数名：delete_extra_space(s string) string
功  能:删除字符串中多余的空格(含tab)，有多个空格时，仅保留一个空格，同时将字符串中的tab换为空格
参  数:s string:原始字符串
返回值:string:删除多余空格后的字符串
创建时间:2018年12月3日
修订信息:
*/
func delete_extra_space(s string) string {
	//删除字符串中的多余空格，有多个空格时，仅保留一个空格
	s1 := strings.Replace(s, "	", " ", -1)       //替换tab为空格
	regstr := "\\s{2,}"                          //两个及两个以上空格的正则表达式
	reg, _ := regexp.Compile(regstr)             //编译正则表达式
	s2 := make([]byte, len(s1))                  //定义字符数组切片
	copy(s2, s1)                                 //将字符串复制到切片
	spc_index := reg.FindStringIndex(string(s2)) //在字符串中搜索
	for len(spc_index) > 0 {                     //找到适配项
		s2 = append(s2[:spc_index[0]+1], s2[spc_index[1]:]...) //删除多余空格
		spc_index = reg.FindStringIndex(string(s2))            //继续在字符串中搜索
	}
	return string(s2)
}
