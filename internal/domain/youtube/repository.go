package youtube

import "context"

type IRepository interface {
	FindByID(ctx context.Context, id string) (*Channel, error)
}
