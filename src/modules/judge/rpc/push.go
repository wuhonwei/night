package rpc

import (
	"nightingale-club203/src/common/dataobj"
	"nightingale-club203/src/modules/judge/cache"
	"nightingale-club203/src/modules/judge/judge"
	"nightingale-club203/src/toolkits/stats"

	"github.com/toolkits/pkg/logger"
)

type Judge int

func (j *Judge) Ping(req dataobj.NullRpcRequest, resp *dataobj.SimpleRpcResponse) error {
	return nil
}

func (j *Judge) Send(items []*dataobj.JudgeItem, resp *dataobj.SimpleRpcResponse) error {
	// 把当前时间的计算放在最外层，是为了减少获取时间时的系统调用开销

	for _, item := range items {
		now := item.Timestamp
		pk := item.MD5()
		logger.Debugf("recv-->%+v", item)
		stats.Counter.Set("push.in", 1)

		go judge.ToJudge(cache.HistoryBigMap[pk[0:2]], pk, item, now)
	}

	return nil
}