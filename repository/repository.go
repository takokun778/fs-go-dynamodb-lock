package repository

import (
	"context"
	"lock/model"
)

type Repository interface {
	Save(context.Context, model.Model) error
	Update(context.Context, model.Model) error
	Find(context.Context, string) (model.Model, error)
	Delete(context.Context, string) error
}
