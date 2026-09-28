package repo

import (
	"context"

	"github.com/LeBaoTai/SDN-RD/internal/shared/model"
)

func (r *Repo) GetAllDevices(ctx context.Context) (*[]model.Device, error) {
	var devices []model.Device

	result := r.dbConn.WithContext(ctx).Find(&devices)
	if result.Error != nil {
		return nil, result.Error
	}
	return &devices, nil
}
