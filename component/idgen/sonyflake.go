package idgen

import (
	"github.com/sony/sonyflake"
	"github.com/soukengo/gopkg/util/strings"
)

type sonyflakeGenerator struct {
	core *sonyflake.Sonyflake
	hash *strings.HashID
}

func NewSonyflakeGenerator() Generator {
	return &sonyflakeGenerator{core: sonyflake.NewSonyflake(sonyflake.Settings{}), hash: strings.NewHashID(string(Sonyflake), 10)}
}

func (g *sonyflakeGenerator) GenID() (id string, err error) {
	num, err := g.core.NextID()
	if err != nil {
		return
	}
	return g.hash.Encode([]int64{int64(num)})
}
