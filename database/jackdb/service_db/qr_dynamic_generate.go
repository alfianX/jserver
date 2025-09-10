package servicedb

import (
	"context"
	"time"

	"github.com/alfianX/jserver/database/jackdb/model"
)

func (s Service) QrGenerateSave(ctx context.Context, host, requestData string) (int64, error) {
	data := model.QrDynamicGenerate{
		Host:        host,
		RequestData: requestData,
		CreatedAt:   time.Now(),
	}

	id, err := s.repo.QrGenerateSave(ctx, &data)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s Service) QrGenerateUpdateResponse(ctx context.Context, responseData string, id int64) error {
	data := model.QrDynamicGenerate{
		ID:           id,
		ResponseData: responseData,
		UpdatedAt:    time.Now(),
	}

	err := s.repo.QrGenerateUpdateResponse(ctx, &data)

	return err
}
