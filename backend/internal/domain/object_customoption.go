package domain

import "context"

type ObjectCustomOption struct {
	ObjectId       string
	CustomOptionId string
	Value          string
}

type ObjectCustomOptionRepository interface {
	GetObjectCustomOptionsByObjectId(ctx context.Context, objectId string) ([]*ObjectCustomOption, error)
	AddObjectCustomOption(ctx context.Context, objectCustomOption *ObjectCustomOption) error
	UpdateObjectCustomOption(ctx context.Context, objectCustomOption *ObjectCustomOption) error
}
