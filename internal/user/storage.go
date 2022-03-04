package user

type Repository interface {
	Get(id string) (User, error)
	Insert(u CreateUserRequest) (string, error)
	Update(id string, u UpdateUserRequest) error
	Delete(id string) error
	GetAll() (UserStore, error)
}
