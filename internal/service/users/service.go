package users

type service struct {
	userRepo userRepository
}

func New(userRepo userRepository) *service {
	return &service{userRepo: userRepo}
}
