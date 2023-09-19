package wombat

import (
	"context"
	"encoding/json"
	"github.com/KyberNetwork/ethrpc"
	"github.com/KyberNetwork/kyberswap-dex-lib/pkg/entity"
	"github.com/KyberNetwork/logger"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
)

type PoolTracker struct {
	config       *Config
	ethrpcClient *ethrpc.Client
}

func NewPoolTracker(cfg *Config, ethrpcClient *ethrpc.Client) *PoolTracker {
	return &PoolTracker{
		config:       cfg,
		ethrpcClient: ethrpcClient,
	}
}

func (d *PoolTracker) GetNewPoolState(ctx context.Context, p entity.Pool) (entity.Pool, error) {
	logger.WithFields(logger.Fields{
		"address": p.Address,
	}).Infof("[%s] Start getting new states of pool", p.Type)

	var ampFactor, haircutRate, startCovRatio, endcovRatio *big.Int
	var assetData = make([]Asset, len(p.Tokens))

	calls := d.ethrpcClient.NewRequest().SetContext(ctx)
	calls.AddCall(&ethrpc.Call{
		ABI:    PoolV2ABI,
		Target: p.Address,
		Method: poolMethodAmpFactor,
		Params: nil,
	}, []interface{}{&ampFactor})
	calls.AddCall(&ethrpc.Call{
		ABI:    PoolV2ABI,
		Target: p.Address,
		Method: poolMethodHaircutRate,
		Params: nil,
	}, []interface{}{&haircutRate})
	calls.AddCall(&ethrpc.Call{
		ABI:    PoolV2ABI,
		Target: p.Address,
		Method: poolMethodStartCovRatio,
		Params: nil,
	}, []interface{}{&startCovRatio})
	calls.AddCall(&ethrpc.Call{
		ABI:    PoolV2ABI,
		Target: p.Address,
		Method: poolMethodEndCovRatio,
		Params: nil,
	}, []interface{}{&endcovRatio})
	for i, token := range p.Tokens {
		calls.AddCall(&ethrpc.Call{
			ABI:    PoolV2ABI,
			Target: p.Address,
			Method: poolMethodAddressOfAsset,
			Params: []interface{}{common.HexToAddress(token.Address)},
		}, []interface{}{&assetData[i].Address})
	}
	if _, err := calls.Aggregate(); err != nil {
		logger.WithFields(logger.Fields{
			"type":    p.Type,
			"address": p.Address,
		}).Errorf("failed to aggregate call")
		return entity.Pool{}, err
	}

	calls = d.ethrpcClient.NewRequest().SetContext(ctx)
	for i, asset := range assetData {
		calls.AddCall(&ethrpc.Call{
			ABI:    PoolV2ABI,
			Target: asset.Address,
			Method: assetMethodCash,
			Params: nil,
		}, []interface{}{&assetData[i].Cash})
		calls.AddCall(&ethrpc.Call{
			ABI:    PoolV2ABI,
			Target: asset.Address,
			Method: assetMethodLiability,
			Params: nil,
		}, []interface{}{&assetData[i].Liability})
		calls.AddCall(&ethrpc.Call{
			ABI:    PoolV2ABI,
			Target: asset.Address,
			Method: assetMethodGetRelativePrice,
			Params: nil,
		}, []interface{}{&assetData[i].RelativePrice})
	}
	if _, err := calls.TryAggregate(); err != nil {
		logger.WithFields(logger.Fields{
			"type":    p.Type,
			"address": p.Address,
		}).Errorf("failed to try aggregate call")

		return entity.Pool{}, err
	}

	var assetMap = make(map[string]Asset)
	var reserves = make([]string, len(p.Tokens))
	for i, token := range p.Tokens {
		assetMap[token.Address] = assetData[i]
		reserves[i] = assetData[i].Liability.String()
	}

	extraByte, err := json.Marshal(Extra{
		HaircutRate:   haircutRate,
		AmpFactor:     ampFactor,
		StartCovRatio: startCovRatio,
		EndCovRatio:   endcovRatio,
		AssetMap:      assetMap,
	})
	if err != nil {
		logger.WithFields(logger.Fields{
			"address": p.Address,
			"type":    p.Type,
			"error":   err,
		}).Errorf("failed to marshal extra data")

		return entity.Pool{}, err
	}

	p.Reserves = reserves
	p.Extra = string(extraByte)
	p.Timestamp = time.Now().Unix()

	logger.WithFields(logger.Fields{
		"address": p.Address,
		"type":    p.Type,
	}).Infof("[%s] Finish getting new state of pool", p.Type)

	return p, nil
}
