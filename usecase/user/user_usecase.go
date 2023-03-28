package user

import (
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/infra/repository"
)

type UserUseCase interface {
	CreateUser(dto *CreateUserDTO) (*UserDTO, error)
	UpdateUser(id string, dto *UpadateUserDTO) (*UserDTO, error)
	FindAllUser() ([]*UserDTO, error)
	FindUserByID(id string) (*UserDTO, error)
	RemoveUser(id string) error
}

type UserUseCaseImpl struct {
	UserRepository repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &UserUseCaseImpl{
		UserRepository: userRepo,
	}
}

func (u *UserUseCaseImpl) CreateUser(data *CreateUserDTO) (*UserDTO, error) {

	user, err := entities.NewUser(data.Name, data.Email, data.Password)
	if err != nil {
		return nil, err
	}

	if err := u.UserRepository.Save(user); err != nil {
		return nil, err
	}

	userDto := &UserDTO{}

	return userDto.FromUser(user), nil
}

func (u *UserUseCaseImpl) FindAllUser() ([]*UserDTO, error) {
	users, err := u.UserRepository.FindAll()
	if err != nil {
		return nil, err
	}

	dtoUser := &UserDTO{}
	dtoUsers := dtoUser.FromUserList(users)

	return dtoUsers, nil
}

func (u *UserUseCaseImpl) FindUserByID(id string) (*UserDTO, error) {
	user, err := u.UserRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	dtoUser := &UserDTO{}

	return dtoUser.FromUser(user), nil
}

func (u *UserUseCaseImpl) UpdateUser(id string, dto *UpadateUserDTO) (*UserDTO, error) {
	user, err := u.UserRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := user.Update(dto.Name, dto.Email, dto.Password); err != nil {
		return nil, err
	}

	if err := u.UserRepository.Update(user); err != nil {
		return nil, err
	}

	userDTO := &UserDTO{}

	return userDTO.FromUser(user), nil
}

func (u *UserUseCaseImpl) RemoveUser(id string) error {
	user, err := u.UserRepository.FindByID(id)
	if err != nil {
		return err
	}

	if err := u.UserRepository.Delete(user.ID); err != nil {
		return err
	}

	return nil
}
