package wombat

import "math/big"

type ExtraMain struct {
	HaircutRate   *big.Int         `json:"haircutRate"`
	AmpFactor     *big.Int         `json:"ampFactor"`
	StartCovRatio *big.Int         `json:"startCovRatio"`
	EndCovRatio   *big.Int         `json:"endCovRatio"`
	AssetMap      map[string]Asset `json:"assetMap"`
}

type Asset struct {
	Cash                    *big.Int
	Liability               *big.Int
	UnderlyingTokenDecimals uint8
	RelativePrice           *big.Int
}

type Gas struct {
	Swap int64
}
