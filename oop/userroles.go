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
	HasPermission(Permission) bool
	GetRole() Role
}

type Role string

const (
	BasicRole     Role = "basic_user"
	ModeratorRole Role = "moderator"
	AdminRole     Role = "admin"
)

type Permission string

const (
	Read    Permission = "read"
	Edit    Permission = "edit"
	BanUser Permission = "ban_user"
	Delete  Permission = "delete"
	Manage  Permission = "manage_roles"
)

type BasicUser struct {
	username string
}

func NewBasicUser(username string) *BasicUser {
	return &BasicUser{
		username: username,
	}
}

func (u *BasicUser) GetUsername() string {
	return u.username
}

func (u *BasicUser) HasPermission(perm Permission) bool {
	switch perm {
	case Read:
		return true
	}

	return false
}

func (u *BasicUser) GetRole() Role {
	return BasicRole
}

type Moderator struct {
	*BasicUser
}

func NewModerator(username string) *Moderator {
	return &Moderator{
		BasicUser: NewBasicUser(username),
	}
}

func (u *Moderator) HasPermission(perm Permission) bool {
	switch perm {
	case Edit, BanUser:
		return true
	}

	return u.BasicUser.HasPermission(perm)
}

func (u *Moderator) GetRole() Role {
	return ModeratorRole
}

type Admin struct {
	*Moderator
}

func NewAdmin(username string) *Admin {
	return &Admin{
		Moderator: NewModerator(username),
	}
}

func (u *Admin) HasPermission(perm Permission) bool {
	switch perm {
	case Delete, Manage:
		return true
	}

	return u.Moderator.HasPermission(perm)
}

func (u *Admin) GetRole() Role {
	return AdminRole
}
