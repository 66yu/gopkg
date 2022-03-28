package errorGroup

type AccessTokenInvalidError struct {
}

func (AccessTokenInvalidError) Error() string {
	return "access token invalid"
}
