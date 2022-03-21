package main

import (
	"context"
	"fmt"
	"robot/spider"
	"strings"
	"time"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/openapi"
)

type Processor struct {
	api openapi.OpenAPI
}

const (
	defaultMessage    = "你好，欢迎来到疫小知。在这里你可以查询疫情的实时情况！输入“帮助”查看相关指令"
	helpMessage       = "你可以输入如“全国疫情”，“广东省疫情”，“”深圳市疫情“，”防护措施“来查询相关数据哦"
	unknownMessage    = "疫小知暂时无法解析你的消息哦，输入“帮助”查看相关指令"
	preventionMessage = `戴口罩，勤洗手 ，测体温，勤通风，少聚集，勤消毒`
	phoneMessage      = "中国疾病防控中心电话：010-58900001。各地的疾控中心电话可以登录当地疾控中心官网查询。希望我的回答能够帮助到你"
	chinaTotalTpl     = `你好：%s
截至：%s 
累计确诊：%d
现有确诊：%d
现有本土确诊：%d
境外输入：%d
无症状感染者：%d
累计死亡：%d
	`
	tpl = `你好：%s
%s：
累计确诊：%d
现有确诊：%d
无症状感染者：%d
累计治愈：%d
累计死亡：%d
	`
	timeTpl = `你好：%s
在子频道 %s 收到消息
收到的消息发送时时间为：%s
当前本地时间为：%s
消息来自：%s
`
)

func (p Processor) ProcessMessage(input string, data *dto.WSATMessageData) error {
	ctx := context.Background()
	cmd := message.ParseCommand(input)
	toCreate := &dto.MessageToCreate{
		Content: defaultMessage + message.Emoji(14),
		MessageReference: &dto.MessageReference{
			// 引用这条消息
			MessageID:             data.ID,
			IgnoreGetMessageError: true,
		},
	}

	switch {
	case cmd.Cmd == "hi" || cmd.Cmd == "":
		p.sendReply(ctx, data.ChannelID, toCreate)
	case strings.Contains(cmd.Cmd, "全国疫情"):
		toCreate.Content = getChinaTotalContent(data)
		p.sendReply(ctx, data.ChannelID, toCreate)
	case cmd.Cmd == "time":
		toCreate.Content = genReplyContent(data)
		p.sendReply(ctx, data.ChannelID, toCreate)
	case cmd.Cmd == "帮助":
		toCreate.Content = helpMessage
		p.sendReply(ctx, data.ChannelID, toCreate)
	case cmd.Cmd == "电话":
		toCreate.Content = phoneMessage
		p.sendReply(ctx, data.ChannelID, toCreate)
	//广东省
	case strings.Contains(cmd.Cmd, "省") && !strings.Contains(cmd.Cmd, "市"):
		toCreate.Content = getProvinceContent(cmd.Cmd, data)
		p.sendReply(ctx, data.ChannelID, toCreate)
	//广东省深圳市
	case strings.Contains(cmd.Cmd, "省") && strings.Contains(cmd.Cmd, "市"):
		toCreate.Content = getProvinceCityContent(cmd.Cmd, data)
		p.sendReply(ctx, data.ChannelID, toCreate)
	//深圳市
	case strings.Contains(cmd.Cmd, "市"):
		toCreate.Content = getCityContent(cmd.Cmd, data)
		p.sendReply(ctx, data.ChannelID, toCreate)
	case strings.Contains(cmd.Cmd, "防护"):
		toCreate.Content = preventionMessage
		p.sendReply(ctx, data.ChannelID, toCreate)
	default:
		toCreate.Content = unknownMessage
		p.sendReply(ctx, data.ChannelID, toCreate)
	}
	return nil
}

func genReplyContent(data *dto.WSATMessageData) string {

	msgTime, _ := data.Timestamp.Time()
	return fmt.Sprintf(
		timeTpl,
		message.MentionUser(data.Author.ID),
		message.MentionChannel(data.ChannelID),
		msgTime, time.Now().Format(time.RFC3339),
		getIP(),
	)
}

func getChinaTotalContent(data *dto.WSATMessageData) string {
	ChinaTotal, err := spider.GetChinaTotal()
	if err != nil {
		return "网络开小差了哦，请稍后再试~"
	}

	return fmt.Sprintf(chinaTotalTpl, message.MentionUser(data.Author.ID),
		ChinaTotal.LastUpdateTime, ChinaTotal.Confirm, ChinaTotal.NowConfirm, ChinaTotal.LocalConfirm,
		ChinaTotal.ImportedCase, ChinaTotal.NoInfect, ChinaTotal.Dead,
	)
}

func getProvinceContent(cmd string, data *dto.WSATMessageData) string {
	idx := strings.Index(cmd, "省")
	province := cmd[:idx]
	if verifyArea(cmd[:idx]) == -1 {
		return "您输入的省份不存在！请重新输入正确的省份哦"
	}

	ProvinceCityTotal, err := spider.GetProvinceTotal(province)
	if err != nil {
		return "网络开小差了哦，请稍后再试~"
	}

	return fmt.Sprintf(tpl, message.MentionUser(data.Author.ID),
		province+"省", ProvinceCityTotal.Confirm, ProvinceCityTotal.NowConfirm, ProvinceCityTotal.Wzz,
		ProvinceCityTotal.Heal, ProvinceCityTotal.Dead,
	)

}

func getCityContent(cmd string, data *dto.WSATMessageData) string {
	idx := strings.Index(cmd, "市")
	city := cmd[:idx]

	ProvinceCityTotal, err := spider.GetCityTotal(city)
	if err != nil {
		return "网络开小差了哦，请稍后再试~"
	}

	if !ProvinceCityTotal.Exists {
		return "你输入的城市不存在！请重新输入正确的城市"
	}

	return fmt.Sprintf(tpl, message.MentionUser(data.Author.ID),
		city+"市", ProvinceCityTotal.Confirm, ProvinceCityTotal.NowConfirm, ProvinceCityTotal.Wzz,
		ProvinceCityTotal.Heal, ProvinceCityTotal.Dead,
	)

}

func getProvinceCityContent(cmd string, data *dto.WSATMessageData) string {
	p_idx := strings.Index(cmd, "省")
	c_idx := strings.Index(cmd, "市")
	province := cmd[:p_idx]
	city := cmd[p_idx+3 : c_idx]
	if verifyArea(province) == -1 {
		return "您输入的省份不存在！请重新输入正确的省份"
	}

	ProvinceCityTotal, err := spider.GetProvinceCityTotal(province, city)
	if err != nil {
		return "网络开小差了哦，请稍后再试~"
	}

	if !ProvinceCityTotal.Exists {
		return "你输入的城市不存在！请重新输入正确的城市"
	}

	return fmt.Sprintf(tpl, message.MentionUser(data.Author.ID),
		strings.Join([]string{province, "省", city, "市"}, ""),
		ProvinceCityTotal.Confirm, ProvinceCityTotal.NowConfirm, ProvinceCityTotal.Wzz, ProvinceCityTotal.Heal,
		ProvinceCityTotal.Dead,
	)
}
