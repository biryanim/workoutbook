package auth

import "context"

func (s *serv) Check(ctx context.Context, token string) (int64, bool, error) {
	claims, err := verifyToken(token, s.jwtConfig.TokenSecret())
	if err != nil {
		return 0, false, err
	}

	return claims.UserID, true, nil
}
