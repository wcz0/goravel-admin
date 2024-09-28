package services

type AdminMenuService struct {
	*Service
}

func NewAdminMenuService() *AdminMenuService {
	return &AdminMenuService{
		Service:  NewService(),
	}
}
