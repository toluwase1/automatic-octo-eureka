package models

func (u *User) ActivateWallet(activate bool) {
	u.IsActive = activate
}
