package velodromev2

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"

	"github.com/KyberNetwork/ethrpc"
	"github.com/KyberNetwork/logger"

	"github.com/KyberNetwork/kyberswap-dex-lib/pkg/entity"
)

type PoolTracker struct {
	config       *Config
	ethrpcClient *ethrpc.Client
}

func NewPoolTracker(
	cfg *Config,
	ethrpcClient *ethrpc.Client,
) (*PoolTracker, error) {
	return &PoolTracker{
		config:       cfg,
		ethrpcClient: ethrpcClient,
	}, nil
}

func (d *PoolTracker) GetNewPoolState(ctx context.Context, p entity.Pool) (entity.Pool, error) {
	logger.WithFields(logger.Fields{
		"address": p.Address,
	}).Infof("[%s] Start getting new state of pool", p.Type)

	var (
		reserve Reserves
		fee     *big.Int
	)

	staticExtra, err := extractStaticExtra(p.StaticExtra)
	if err != nil {
		logger.WithFields(logger.Fields{
			"poolAddress": p.Address,
			"error":       err,
		}).Errorf("failed to extract static extra")

		return entity.Pool{}, err
	}

	calls := d.ethrpcClient.NewRequest().SetContext(ctx)

	calls.AddCall(&ethrpc.Call{
		ABI:    pairABI,
		Target: p.Address,
		Method: poolMethodGetReserves,
		Params: nil,
	}, []interface{}{&reserve})

	calls.AddCall(&ethrpc.Call{
		ABI:    factoryABI,
		Target: d.config.FactoryAddress,
		Method: factoryMethodGetFee,
		Params: []interface{}{common.HexToAddress(p.Address), staticExtra.Stable},
	}, []interface{}{&fee})

	if _, err := calls.Aggregate(); err != nil {
		logger.WithFields(logger.Fields{
			"poolAddress": p.Address,
			"error":       err,
		}).Errorf("failed to aggregate to get pool data")

		return entity.Pool{}, err
	}

	swapFee := fee.Int64()

	p.Reserves = entity.PoolReserves{reserve.Reserve0.String(), reserve.Reserve1.String()}
	p.SwapFee = float64(swapFee) / bps
	p.Timestamp = time.Now().Unix()

	logger.WithFields(logger.Fields{
		"address": p.Address,
	}).Infof("[%s] Finish getting new state of pool", p.Type)

	return p, nil
}
