package main

import (
	"github.com/wxnacy/wgo/arrays"
	"net"
)

// getIP 一个简易的获取本机 IP 方法
func getIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "0.0.0.0"
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "0.0.0.1"
}

func verifyArea(s string) int {
	provinces := []string{"香港", "吉林", "台湾", "广东", "山东", "福建", "上海", "浙江", "天津",
		"辽宁", "陕西", "广西", "河北", "甘肃", "北京", "江苏", "黑龙江", "四川", "云南", "重庆", "河南", "湖南",
		"安徽", "湖北", "内蒙古", "山西", "贵州", "江西", "海南", "澳门", "青海", "西藏", "宁夏", "新疆"}
	index := arrays.ContainsString(provinces, s)
	return index
}
