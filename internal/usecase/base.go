package usecase

import (
	"errors"
	"path"
	"runtime"

	"github.com/bbrighter/dreams-api/internal/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BaseUseCase struct {
	logger *zap.Logger
}

func NewBaseUseCase(logger *zap.Logger) BaseUseCase {
	return BaseUseCase{logger: logger}
}

func (b *BaseUseCase) HandleError(err error) bool {
	if err != nil {
		pc, _, _, _ := runtime.Caller(1) // 1 = caller of logError
		funcName := runtime.FuncForPC(pc).Name()
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, entity.ErrorNotFound) {
			b.logger.Info(path.Base(funcName), zap.String("reason", "not found"))
			return false
		}
		b.logger.Error(path.Base(funcName), zap.Error(err))
	}
	return err != nil
}
