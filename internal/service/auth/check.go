package auth

import "context"

func (s *serv) Check(ctx context.Context, token string) (bool, error) {
	_, err := verifyToken(token, s.jwtConfig.TokenSecret())
	if err != nil {
		return false, err
	}

	return true, nil
}
