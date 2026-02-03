package oop

/*

# 2. Система управления пользователями и ролями (ООП в Go)

Реализуйте систему управления пользователями с различными ролями и правами доступа, используя принципы ООП: инкапсуляцию, композицию и полиморфизм.

## Задание

### Цель
1. Создать базовый интерфейс `User` с методами:
    - `GetUsername() string` — возвращает имя пользователя.
    - `HasPermission(permission string) bool` — проверяет наличие права доступа.
    - `GetRole() string` — возвращает роль пользователя.
2. Реализовать три типа пользователей:
    - **Обычный пользователь** (`BasicUser`):
        - Может читать данные (`read`), но не может их изменять.
    - **Модератор** (`Moderator`):
        - Наследует права `BasicUser`.
        - Добавляет право `edit` (редактирование данных).
        - Может банить пользователей (`ban_user`).
    - **Администратор** (`Admin`):
        - Наследует права `Moderator`.
        - Добавляет право `delete` (удаление данных).
        - Может управлять ролями (`manage_roles`).

### Требования
- Поля, хранящие права доступа, должны быть **инкапсулированы**.
- Используйте **композицию** для наследования прав.
- Для каждого типа пользователя реализуйте:
    - Конструктор `NewAdmin(username string)`,`NewModerator(username string)`,`NewBasicUser(username string)`.
    - Уникальные права доступа.

*/

type User interface {
	GetUsername() string
	HasPermission(permission string) bool
	GetRole() string
}

type role string

const (
	basicRole     role = "basic_user"
	moderatorRole role = "moderator"
	adminRole     role = "admin"
)

type perm string

const (
	read    perm = "read"
	edit    perm = "edit"
	banUser perm = "ban_user"
	delete  perm = "delete"
	manage  perm = "manage_roles"
)

type BasicUser struct {
	Username string
}

func NewBasicUser(username string) BasicUser {
	return BasicUser{
		Username: username,
	}
}

func (u BasicUser) GetUsername() string {
	return u.Username
}

func (u BasicUser) HasPermission(permission string) bool {
	switch perm(permission) {
	case read:
		return true
	}

	return false
}

func (u BasicUser) GetRole() string {
	return string(basicRole)
}

type Moderator struct {
	BasicUser
}

func NewModerator(username string) Moderator {
	return Moderator{
		BasicUser: NewBasicUser(username),
	}
}

func (u Moderator) HasPermission(permission string) bool {
	switch perm(permission) {
	case edit, banUser:
		return true
	}

	return u.BasicUser.HasPermission(permission)
}

func (u Moderator) GetRole() string {
	return string(moderatorRole)
}

type Admin struct {
	Moderator
}

func NewAdmin(username string) Admin {
	return Admin{
		Moderator: NewModerator(username),
	}
}

func (u Admin) HasPermission(permission string) bool {
	switch perm(permission) {
	case delete, manage:
		return true
	}

	return u.Moderator.HasPermission(permission)
}

func (u Admin) GetRole() string {
	return string(adminRole)
}
