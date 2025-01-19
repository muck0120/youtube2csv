//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -package=$GOPACKAGE -destination=./repository_mock.go

package youtube

import "context"

type IRepository interface {
	FindByID(ctx context.Context, id string) (*Channel, error)
}
