package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("snowflake", map[string]any{
		"epoch": config.Env("SNOWFLAKE_EPOCH", "2019-04-01 00:00:00"),
		"worker_id_bit_length": config.Env("SNOWFLAKE_WORKER_ID_BIT_LENGTH", 5),
		/*
		   |--------------------------------------------------------------------------
		   | Snowflake Configuration
		   |--------------------------------------------------------------------------
		   |
		   | Here you may configure the log settings for snowflake.
		   | If you are using multiple servers, please assign unique
		   | ID(1-31) for Snowflake.
		   |
		   | Available Settings: 1-31
		   |
		*/
		"worker_id": config.Env("SNOWFLAKE_WORKER_ID", 1),

		"seq_bit_length": config.Env("SNOWFLAKE_DATACENTER_ID", 1),
	})
}
