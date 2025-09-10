package servicedb

import (
	"context"
	"time"

	"github.com/alfianX/jserver/database/jackdb/model"
)

func (s Service) QrNotifSave(ctx context.Context, requestData string) (int64, error) {
	data := model.QrNotif{
		RequestData: requestData,
		CreatedAt:   time.Now(),
	}

	id, err := s.repo.QrNotifSave(ctx, &data)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s Service) QrNotifUpdateResponse(ctx context.Context, responseData string, id int64) error {
	data := model.QrNotif{
		ID:           id,
		ResponseData: responseData,
		UpdatedAt:    time.Now(),
	}

	err := s.repo.QrNotifUpdateResponse(ctx, &data)

	return err
}
