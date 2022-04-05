package models

func (user *User) ActivateWallet(activate bool) {
	user.IsActive = activate
}
