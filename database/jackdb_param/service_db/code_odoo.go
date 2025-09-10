package servicedb

import "context"

func (s Service) CodeOdooGetName(ctx context.Context, host, typeTrx string) (string, string, error) {
	name, code, err := s.repo.CodeOdooGetName(ctx, host, typeTrx)
	if err != nil {
		return "", "", err
	}

	return name, code, nil
}
