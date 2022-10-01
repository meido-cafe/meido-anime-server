package service

import (
	"context"
	"meido-anime-server/internal/common"
	"meido-anime-server/internal/repo"
)

type Service struct {
	repo    *repo.Repo
	bangumi *repo.BangumiRepo
	qb      *repo.Qbittorrent
}

func NewService() *Service {
	return &Service{
		repo:    repo.NewRepo(common.GetSqlite()),
		bangumi: repo.NewBangumiRepo(),
		qb:      repo.NewQbittorrent(),
	}
}

type transactionResult struct {
	Error   error
	TxError error
}

// 事务
func (this *Service) transaction(ctx context.Context, fn func(repo *repo.Repo) error) (ret transactionResult) {
	db := common.GetSqlite()
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		ret.TxError = err
		return
	}
	defer func() {
		if ret.TxError == nil && ret.Error != nil {
			ret.TxError = tx.Rollback()
		}
	}()
	re := repo.NewRepo(tx)

	if funcErr := fn(re); funcErr != nil {
		ret.Error = err
		return
	}

	ret.TxError = tx.Commit()
	return
}
