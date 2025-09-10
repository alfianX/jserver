package servicedb

import (
	"context"
	"time"

	"github.com/alfianX/jserver/database/jackdb/model"
)

func (s Service) QrCallbackSave(ctx context.Context, host, requestData string) (int64, error) {
	data := model.QrDynamicCallback{
		Host:        host,
		RequestData: requestData,
		CreatedAt:   time.Now(),
	}

	id, err := s.repo.QrCallbackSave(ctx, &data)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s Service) QrCallbackUpdateResponse(ctx context.Context, responseData string, id int64) error {
	data := model.QrDynamicCallback{
		ID:           id,
		ResponseData: responseData,
		UpdatedAt:    time.Now(),
	}

	err := s.repo.QrCallbackUpdateResponse(ctx, &data)

	return err
}
