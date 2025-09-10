package servicedb

import "context"

func (s Service) HostDebitGetAddress(ctx context.Context, name string) (string, error) {
	address, err := s.repo.HostDebitGetAddress(ctx, name)
	if err != nil {
		return "", err
	}

	return address, nil
}
