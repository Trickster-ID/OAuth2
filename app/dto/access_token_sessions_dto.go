package dto

import (
	"oauth2/app/global/model"
	"oauth2/app/models"
)

type GetByRefreshTokenChan struct {
	Data   *models.RefreshTokenSession
	ErrLog *model.ErrorLog
}
