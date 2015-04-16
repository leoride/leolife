package util

const _EXP_EMAIL = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`

func IsEmailValid(s string) bool {
	if s != "" {
		return true
	} else {
		return false
	}
}
