/*
 * @Description:
 * @Author: Moqi
 * @Date: 2018-12-12 10:34:19
 * @Email: str@li.cm
 * @Github: https://github.com/strugglerx
 * @LastEditors: Moqi
 * @LastEditTime: 2019-03-09 16:00:15
 */

package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/asmcos/requests"
	"github.com/astaxie/beego"
	"github.com/tidwall/gjson" // "github.com/PuerkitoBio/goquery"
	// // "reflect"
)

//自定义Headers里的Cookie
func normalHeader(session string) requests.Header {
	//弃用
	/* 	{
		"Host":       "eip.imnu.edu.cn",
		"Origin":     "http://eip.imnu.edu.cn",
		"Connection": "keep-alive",
		"User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_3 like Mac OS X) AppleWebKit/603.3.8 (KHTML, like Gecko) Mobile/14G60 MicroMessenger/6.6.7 NetType/WIFI Language/zh_CN",
		"Referer":    "http://eip.imnu.edu.cn/EIP/weixin/weui/chengjichaxun.html",
		"Cookie":     "",
	} */
	return DefaultHeader
}

//网费
func netlist(cookie *http.Cookie) string {
	proxy := beego.AppConfig.String("proxy::url")
	var reqs = requests.Requests()
	reqs.SetCookie(cookie)
	reqs.SetTimeout(3)
	//获取配置里的代理链接
	reqs.Proxy(proxy)
	req, err := reqs.Post(EipDomain+"/EIP/edu/wangfei/queryUsrBindProduct.htm", DefaultHeader)
	if err != nil {
		return "-1"
	}
	status := gjson.Get(req.Text(), "#").Bool()
	//fmt.Printf("%+v",status)
	if status {
		return gjson.Get(req.Text(), "0.otherData").String()
	} else {
		return "-1"
	}

}

//网费调用third接口

func netThird(user string) string {
	var reqs = requests.Requests()
	reqs.SetTimeout(3)
	headers := requests.Header{
		"Origin":     "https://servicewechat.com",
		"User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_3 like Mac OS X) AppleWebKit/603.3.8 (KHTML, like Gecko) Mobile/14G60 MicroMessenger/6.6.7 NetType/WIFI Language/zh_CN",
		"Referer":    "https://servicewechat.com/wx2bd97997e95bc55e/51/page-frame.html",
	}
	url := fmt.Sprintf("https://www.enjfun.com/weimnu/net?xh=%s", user)
	req, err := reqs.Get(url, headers)
	if err != nil {
		return "-1"
	}
	status := gjson.Get(req.Text(), "success").Bool()
	if status {
		return gjson.Get(req.Text(), "result.0.otherData").String()
	} else {
		return "-1"
	}
}

//校园卡余额 慢且经常不能用
func cardremain(cookie *http.Cookie) {
	var reqs = requests.Requests()
	reqs.SetCookie(cookie)
	reqs.Debug = 1
	reqs.Proxy("http://140.143.96.216:80")
	req, _ := reqs.Post(EipDomain+"/EIP/queryservice/query.htm?snumber=QRY_BAL&xh=20151105822", DefaultHeader)
	fmt.Println(req.Text())
}

//校园卡消费列表
func cardlist(cookie *http.Cookie) string {
	var reqs = requests.Requests()
	reqs.SetCookie(cookie)
	reqs.SetTimeout(3)
	req, err := reqs.Post(EipDomain+"/EIP/edu/ykt_tongji.htm", DefaultHeader)
	if err != nil {
		return "-1"
	}
	// fmt.Println(req.Text())
	return gjson.Get(req.Text(), "tongji.0.CDATE").String()
}

//校园卡消费明细
func carddetail(date string, cookie *http.Cookie) string {
	var reqs = requests.Requests()
	reqs.SetCookie(cookie)
	reqs.SetTimeout(3)
	data := requests.Datas{
		"date": date,
	}
	req, err := reqs.Post(EipDomain+"/EIP/edu/ykt_mingxi.htm", DefaultHeader, data)
	if err != nil {
		return "-1"
	}
	// fmt.Println(req.Text())
	pop := fmt.Sprintf("%d", gjson.Get(req.Text(), "mingxi.#").Int()-1)
	// s := strconv.Itoa(i) int转string
	remain := gjson.Get(req.Text(), "mingxi."+string(pop)+".ACCOST").String()
	RemainJson := make(map[string]string)
	RemainJson["CARDBAL"] = remain
	r, _ := json.Marshal(RemainJson)
	return string(r)
}

//student information
func info(cookie *http.Cookie) string {
	var reqs = requests.Requests()
	reqs.SetCookie(cookie)
	reqs.SetTimeout(3)
	req, err := reqs.Post(EipDomain+"/EIP/edu/xueji.htm", DefaultHeader)
	if err != nil {
		return "-1"
	}
	return req.Text()
}

//class
func class_(date string, cookie *http.Cookie) string {
	var reqs = requests.Requests()
	reqs.SetCookie(cookie)
	reqs.SetTimeout(3)
	data := requests.Datas{
		"monday_": date,
	}
	req, err := reqs.Post(EipDomain+"/EIP/qiandao/kebiao/queryKebiaoByUserId.htm", DefaultHeader, data)
	if err != nil {
		return "-1"
	}
	if len(req.Text()) > 0 {
		return req.Text()
	}
	return "-1"

}

//scores
func score(cookie *http.Cookie) string {
	var reqs = requests.Requests()
	reqs.SetCookie(cookie)
	reqs.SetTimeout(3)
	req, err := reqs.Post(EipDomain+"/EIP/edu/chengji.htm", DefaultHeader)
	if err != nil {
		return "-1"
	}
	// fmt.Println(req.Text())
	change := []byte(req.Text())

	var eipformat []EipStr
	json.Unmarshal(change, &eipformat)
	// fmt.Printf("%+v", eipfor[0].CJ[0])
	var customformat []Custom
	for _, v := range eipformat {
		var tempitem []CustomItem
		for _, o := range v.CJ {
			item := CustomItem{o.KCM, o.XF, o.XKLX, o.CJ}
			tempitem = append(tempitem, item)
		}
		tempset := Custom{v.XQ, tempitem}
		customformat = append(customformat, tempset)
	}
	// fmt.Printf("%+v", customformat)

	customlen := len(customformat)        //获取长度
	recustom := make([]Custom, customlen) //开辟空间
	for i, v := range customformat {
		//反序列
		recustom[customlen-i-1] = v
	}
	result, _ := json.Marshal(recustom)
	return string(result)
}

//登录
func login(user string) (bool, *http.Cookie) {
	// reqs.Debug = 1
	headers := requests.Header{
		"Host":                      "eip.imnu.edu.cn",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Connection":                "keep-alive",
		"Upgrade-Insecure-Requests": "1",
		"Cookie":                    "EIPUserId=" + user,
		"User-Agent":                "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_3 like Mac OS X) AppleWebKit/603.3.8 (KHTML, like Gecko) Mobile/14G60 MicroMessenger/6.6.7 NetType/WIFI Language/zh_CN",
		"Accept-Language":           "zh-cn",
		"Accept-Encoding":           "gzip, deflate",
	}
	resp, err := requests.Get(EipDomain+"/EIP/weixinEnterprise/cookie.htm?url=/weixin/weui/jiugongge.htmlSPT623e8474576a41e5959ec47f0505109e", headers)
	cookie := &http.Cookie{}
	if err != nil || resp.R.StatusCode != 200 {
		return false, cookie
	}
	//fmt.Printf("%+v", resp.R)
	head := resp.R.Header["Content-Type"][0]
	reg := regexp.MustCompile(`gbk`)
	if reg.MatchString(head) {
		return false, resp.Cookies()[0]
	}
	return true, resp.Cookies()[0]
	// fmt.Println(resp.Cookies())
	//[]*http.Cookie
}

func EipEntry(user string, type_ string, date string) (string, error) {
	status, cookie := login(user)
	if status {
		switch type_ {
		case "score":
			return score(cookie), nil
		case "info":
			return info(cookie), nil
		case "card":
			date := cardlist(cookie)
			if len(date) == 8 {
				result := carddetail(date, cookie)
				return result, nil
			} else {
				return "", errors.New("fail")
			}
		case "net":
			netStatus, _ := beego.AppConfig.Bool("proxy::status")
			// 通过app.conf的proxy::status判断是否启动代理接口
			if netStatus {
				return netlist(cookie), nil
			} else {
				return netThird(user), nil
			}
		case "class_":
			return class_(date, cookie), nil
		default:
			return "", errors.New("fail")
		}
	}
	return "", errors.New("fail")
}
