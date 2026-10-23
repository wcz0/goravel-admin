package main

import (
	"os"

	contractsmodify "github.com/goravel/framework/contracts/packages/modify"
	"github.com/goravel/framework/packages"
	"github.com/goravel/framework/packages/modify"
	"github.com/goravel/framework/support/path"
)

func main() {
	setup := packages.Setup(os.Args)
	moduleImport := setup.Paths().Module().Import()

	install := []contractsmodify.Apply{
		// 在 bootstrap/providers.go 中注册 ServiceProvider。
		modify.RegisterProvider(moduleImport, "&admin.ServiceProvider{}"),

		// 写入可选的配置覆盖文件（已存在时不会覆盖）。
		modify.WhenFileNotExists(
			path.Config("admin.go"),
			modify.File(path.Config("admin.go")).Overwrite(Stubs{}.Config()),
		),
	}

	uninstall := []contractsmodify.Apply{
		modify.File(path.Config("admin.go")).Remove(),
		modify.UnregisterProvider(moduleImport, "&admin.ServiceProvider{}"),
	}

	// 迁移无需在此注册：ServiceProvider.Boot 会通过 facades.Schema().Register 自动挂载。
	setup.Install(install...).Uninstall(uninstall...).Execute()
}
