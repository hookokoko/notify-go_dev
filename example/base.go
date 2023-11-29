package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	////go consumer.Consume()
	////<-time.Tick(3 * time.Second)
	//params := []notify.MessageParam{
	//	{
	//		BizId:     "test biz",
	//		Receiver:  "648646891@qq.com",
	//		Variables: map[string]string{"name": "chenxx"},
	//	},
	//}
	////d := experience.NewDelivery(params, 123,
	////	experience.WithSendService(experience.NewKafkaSendService(kafka.Kafka{Hosts: []string{"127.0.0.1:9092"}})))
	//
	//d := notify.NewNotification(params, 123, storage.MysqlDB())
	//
	//err := d.Send(context.Background())
	//if err != nil {
	//	panic(err)
	//}
	//<-time.Tick(3 * time.Second)
	Test_toml()
}

func Test_toml() {
	var res []map[string][]map[string]any
	str := `[
{"message_type_10":[{"weights":99,"scriptName":"TencentSmsScript"},{"weights":1,"scriptName":"YunPianSmsScript"}]},
{"message_type_20":[{"weights":99,"scriptName":"TencentSmsScript"},{"weights":1,"scriptName":"YunPianSmsScript"}]},
{"message_type_30":[{"weights":20,"scriptName":"TencentSmsScript"}]},
{"message_type_40":[{"weights":20,"scriptName":"TencentSmsScript"}]},
{"message_type_50":[{"weights":20,"scriptName":"LinTongSmsScript"}]}]`
	err := json.Unmarshal([]byte(str), &res)
	fmt.Println(err)
	fmt.Println(res)
	var smsConfig []MessageTypeSmsConfig
	b, _ := json.Marshal(res[0]["message_type_10"])
	err = json.Unmarshal(b, &smsConfig)
	fmt.Println(err)
	fmt.Println(smsConfig)
}

type MessageTypeSmsConfig struct {
	// 权重(决定着流量的占比)
	Weights int `json:"weights"`
	//短信模板若指定了账号，则该字段有值
	SendAccount int `json:"sendAccount"`
	//script名称
	ScriptName string `json:"scriptName"`
}
