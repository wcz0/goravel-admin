package commands

import (
	"fmt"
	"goravel/app/tools"
	"os"
	"path/filepath"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

type Service struct {
}

//Signature The name and signature of the console command.
func (receiver *Service) Signature() string {
	return "make:service"
}

//Description The console command description.
func (receiver *Service) Description() string {
	return "Create a new service struct file"
}

//Extend The console command extend.
func (receiver *Service) Extend() command.Extend {
	return command.Extend{
		Category: "make",
	}
}

//Handle Execute the console command.
func (receiver *Service) Handle(ctx console.Context) error {
	serviceName := tools.SnakeToPascal(ctx.Argument(0))
	fileName := ctx.Argument(0) + "_service"
	serviceDir := "app/services"
	serviceFilePath := filepath.Join(serviceDir, fileName+".go")
	if _, err := os.Stat(serviceFilePath); err == nil {
		return fmt.Errorf("service file %s already exists", serviceFilePath)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check if service file exists: %v", err)
	}

	content := `
package services

type `+ serviceName +`Service struct {
	*Service
}

func New`+ serviceName +`Service() *`+ serviceName +`Service {
	return &`+ serviceName +`Service{
		Service: NewService(),
	}
}
`
	err := os.WriteFile(serviceFilePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to create service file %s: %v", serviceFilePath, err)
	}

	return nil
}
