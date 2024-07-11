package providers

import (
	"time"

	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"
	"github.com/yitter/idgenerator-go/idgen"
)

type AppServiceProvider struct {
}

func (receiver *AppServiceProvider) Register(app foundation.Application) {
	// 注册雪花id
	config := facades.Config()
	var options = idgen.NewIdGeneratorOptions(uint16(config.GetInt("snowflake.worker_id")))
	t, _ := time.Parse("2006-01-02 15:04:05", config.GetString("snowflake.epoch"))
	options.BaseTime = t.UnixNano() / 1000000
	idgen.SetIdGenerator(options)
}

func (receiver *AppServiceProvider) Boot(app foundation.Application) {

}
