package servicedb

import "context"

func (s Service) GetAddressByName(ctx context.Context, name string) (string, error) {
	address, err := s.repo.GetAddressByName(ctx, name)
	if err != nil {
		return "", err
	}

	return address, nil
}

func (s Service) GetAddressQueryByName(ctx context.Context, name string) (string, error) {
	addressQuery, err := s.repo.GetAddressQueryByName(ctx, name)
	if err != nil {
		return "", err
	}

	return addressQuery, nil
}
