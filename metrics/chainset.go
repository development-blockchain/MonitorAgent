package metrics

import (
	"github.com/develope/MonitorAgent/common"
	"math/big"
)

var (
	subnamespace = "chain"
	chainHeight  = NewRegisteredGauge(subnamespace, "chain_height", []string{NodeIdentifier})
	peerCount    = NewRegisteredGauge(subnamespace, "peer_count", []string{NodeIdentifier})
	restHandler  = NewRegisteredCounterVec(subnamespace, "rest_handler", []string{NodeIdentifier, "method"})
	syncFlag     = NewRegisteredGauge(subnamespace, "sync_flag", []string{NodeIdentifier})
	syncPercent  = NewRegisteredGauge(subnamespace, "sync_percent", []string{NodeIdentifier})
)

// UpdateChainHeight 更新节点当前的链高度，每当高度变化时进行更新
func UpdateChainHeight(height *big.Int) {
	chainHeight.WithLabelValues(common.GetHostName()).Set(float64(height.Int64()))
}

// UpdatePeerCount 更新节点的peer数量，每当peer数量变化时进行更新
func UpdatePeerCount(count int) {
	peerCount.WithLabelValues(common.GetHostName()).Set(float64(count))
}

// UpdateRestHandler 记录节点接收到的rest请求, 每当有新的请求时进行更新，传入被调用的方法名
func UpdateRestHandler(method string) {
	restHandler.WithLabelValues(common.GetHostName(), method).Inc()
}

// UpdateSyncInfo 更新节点同步的情况，根据实际情况看是否设置.
// startBlock和endBlock 指的是当前同步开始的区块高度和要同步到的高度，
// curBlock 是当前已经同步到的高度.
func UpdateSyncInfo(syncing bool, startBlock, endBlock, curBlock *big.Int) {
	name := common.GetHostName()
	if syncing {
		syncFlag.WithLabelValues(name).Set(1)
		total := new(big.Int).Sub(endBlock, startBlock)
		current := new(big.Int).Sub(curBlock, startBlock)
		percent := new(big.Float).Quo(new(big.Float).SetInt(current), new(big.Float).SetInt(total))
		p, _ := percent.Float64()
		syncPercent.WithLabelValues(name).Set(p)
	} else {
		syncFlag.WithLabelValues(name).Set(0)
		syncPercent.WithLabelValues(name).Set(0)
	}
}
