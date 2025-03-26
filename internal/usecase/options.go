package usecase

type Option func(*TaskUsecase)

func WithPassword(password string) Option {
	return func(u *TaskUsecase) {
		u.password = password
	}
}
