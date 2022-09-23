package service

import "context"

type Authorizator interface {
	Authorized(ctx context.Context, token string) (bool, error)
}
