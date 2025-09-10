package servicedb

import "context"

func (s Service) HostMiniatmGetAddress(ctx context.Context, name string) (string, error) {
	address, err := s.repo.HostMiniatmGetAddress(ctx, name)
	if err != nil {
		return "", err
	}

	return address, nil
}
