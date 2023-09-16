package wombatlsd_test

import (
	"github.com/KyberNetwork/kyberswap-dex-lib/pkg/entity"
	"github.com/KyberNetwork/kyberswap-dex-lib/pkg/source/pool"
	"github.com/KyberNetwork/kyberswap-dex-lib/pkg/source/wombat/wombatmain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPoolCalcAmountOut(t *testing.T) {
	testCases := []struct {
		name              string
		entityPool        entity.Pool
		tokenAmountIn     pool.TokenAmount
		tokenOut          string
		expectedAmountOut pool.TokenAmount
		expectedErr       error
	}{
		{
			name: "test token0 as tokenIn",
			entityPool: entity.Pool{
				Address:     "",
				Exchange:    "wombat",
				Type:        "wombat-main",
				Reserves:    []string{},
				Tokens:      []*entity.PoolToken{},
				StaticExtra: "",
				Extra:       "",
			},
			tokenAmountIn:     pool.TokenAmount{},
			tokenOut:          "",
			expectedAmountOut: pool.TokenAmount{},
			expectedErr:       nil,
		},
		{
			name: "tes token1 as tokenIn",
			entityPool: entity.Pool{
				Address:  "",
				Exchange: "wombat",
				Type:     "wombat-main",
				Reserves: []string{},
				Tokens:   []*entity.PoolToken{},
				Extra:    "",
			},
			tokenAmountIn:     pool.TokenAmount{},
			tokenOut:          "",
			expectedAmountOut: pool.TokenAmount{},
			expectedErr:       nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pool, err := wombatmain.NewPoolSimulator(tc.entityPool)
			assert.Nil(t, err)
			calcAmountOutResult, err := pool.CalcAmountOut(tc.tokenAmountIn, tc.tokenOut)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedAmountOut, calcAmountOutResult.TokenAmountOut)
		})
	}
}
