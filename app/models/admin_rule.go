package models

import "github.com/wcz0/goravel-authz/models"


type AdminRule struct {
	*models.Rule
}

func NewAdminRule() *AdminRule {
	return &AdminRule{
		Rule: &models.Rule{},
	}
}

func (c *AdminRule) TableName() string {
	return "admin_rules"
}

