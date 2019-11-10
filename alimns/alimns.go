package alimns

import (
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/qit-team/snow-core/config"
)

//依赖注入用的函数
func NewMnsClient(mnsConfig config.MnsConfig) (client ali_mns.MNSClient, err error) {
	//2.1初始化mns client
	defer func() {
		if e := recover(); e != nil {
			s := fmt.Sprintf("ali_mns client init panic: %s", fmt.Sprint(e))
			err = errors.New(s)
		}
	}()

	if mnsConfig.Url != "" {
		client = ali_mns.NewAliMNSClient(mnsConfig.Url,
			mnsConfig.AccessKeyId,
			mnsConfig.AccessKeySecret)
	}
	return
}

func GetMnsBasicQueue(client ali_mns.MNSClient, queueName string) ali_mns.AliMNSQueue {
	var defaultQueue ali_mns.AliMNSQueue

	//根据client创建manager
	queueManager := ali_mns.NewMNSQueueManager(client)

	// 暂时将visibilityTimeout 设置成120，后续将参数暴露给上层，可自行配置
	err := queueManager.CreateQueue(queueName, 0, 65536, 345600, 120, 0, 3)
	if err != nil && !ali_mns.ERR_MNS_QUEUE_ALREADY_EXIST_AND_HAVE_SAME_ATTR.IsEqual(err) {
		fmt.Println(err)
		return defaultQueue
	}
	//最终的最小执行单元queue
	return ali_mns.NewMNSQueue(queueName, client)
}
