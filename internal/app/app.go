package app

import (
	"log"

	"github.com/DanielTitkov/repertoire/internal/domain"
)

type (
	App struct {
		users []*domain.User
	}
)

func (a *App) CreateUser(u domain.User) error {
	log.Println("app creating user", u)
	a.users = append(a.users, &u)
	// and db logic and such
	return nil
}

func (a *App) GetUsersCount() (int, error) {
	log.Println("app user count")
	return len(a.users), nil
}
