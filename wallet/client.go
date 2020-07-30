package wallet

import (
	"context"

	"cloud.google.com/go/firestore"

	"github.com/rabee-inc/go-pkg/log"
)

type Client struct {
	svc  *Service
	fCli *firestore.Client
}

func (c *Client) Get(ctx context.Context, userID string, kind ItemKind) (*Item, error) {
	dst, err := c.svc.Get(ctx, userID, kind)
	if err != nil {
		log.Warningm(ctx, "c.svc.Get", err)
		return nil, err
	}
	return dst, nil
}

func (c *Client) GetMulti(ctx context.Context, userID string, kinds []ItemKind) (map[ItemKind]*Item, error) {
	dsts, err := c.svc.GetMulti(ctx, userID, kinds)
	if err != nil {
		log.Warningm(ctx, "c.svc.GetMulti", err)
		return nil, err
	}
	return dsts, nil
}

func (c *Client) GetAmount(ctx context.Context, userID string, kind ItemKind) (float64, error) {
	dst, err := c.svc.GetAmount(ctx, userID, kind)
	if err != nil {
		log.Warningm(ctx, "c.svc.GetAmount", err)
		return 0, err
	}
	return dst, nil
}

func (c *Client) GetMultiAmount(ctx context.Context, userID string, kinds []ItemKind) (map[ItemKind]float64, error) {
	dsts, err := c.svc.GetMultiAmount(ctx, userID, kinds)
	if err != nil {
		log.Warningm(ctx, "c.svc.GetMultiAmount", err)
		return nil, err
	}
	return dsts, nil
}

func (c *Client) Give(
	ctx context.Context,
	userID string,
	amounts map[ItemKind]float64,
	data map[string]interface{},
	comment string) (map[ItemKind]*Item, error) {
	var dsts map[ItemKind]*Item
	var err error
	err = c.fCli.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		dsts, err = c.svc.Give(ctx, tx, userID, amounts, data, comment)
		if err != nil {
			log.Warningm(ctx, "c.svc.Give", err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Warningm(ctx, "c.fCli.RunTransaction", err)
		return nil, err
	}
	return dsts, nil
}

func (c *Client) Use(
	ctx context.Context,
	userID string,
	amounts map[ItemKind]float64,
	data map[string]interface{},
	comment string) (map[ItemKind]*Item, error) {
	var dsts map[ItemKind]*Item
	var err error
	err = c.fCli.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		dsts, err = c.svc.Use(ctx, tx, userID, amounts, data, comment)
		if err != nil {
			log.Warningm(ctx, "c.svc.Use", err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Warningm(ctx, "c.fCli.RunTransaction", err)
		return nil, err
	}
	return dsts, nil
}

func (c *Client) Exchange(
	ctx context.Context,
	userID string,
	fromKind ItemKind,
	toKind ItemKind,
	amount float64,
	data map[string]interface{},
	comment string) (map[ItemKind]*Item, error) {
	var dsts map[ItemKind]*Item
	var err error
	err = c.fCli.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		dsts, err = c.svc.Exchange(ctx, tx, userID, fromKind, toKind, amount, data, comment)
		if err != nil {
			log.Warningm(ctx, "c.svc.Exchange", err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Warningm(ctx, "c.fCli.RunTransaction", err)
		return nil, err
	}
	return dsts, nil
}

func (c *Client) HistoriesByCursor(
	ctx context.Context,
	userID string,
	kinds []ItemKind,
	limit int,
	cursor string) ([]*ItemHistory, string, error) {
	histories, nCursor, err := c.svc.HistoriesByCursor(ctx, userID, kinds, limit, cursor)
	if err != nil {
		log.Warningm(ctx, "c.svc.HistoriesByCursor", err)
		return nil, "", err
	}
	return histories, nCursor, nil
}

func (c *Client) HistoriesByPeriod(
	ctx context.Context,
	userID string,
	kinds []ItemKind,
	startAt int64,
	endAt int64) ([]*ItemHistory, error) {
	histories, err := c.svc.HistoriesByPeriod(ctx, userID, kinds, startAt, endAt)
	if err != nil {
		log.Warningm(ctx, "c.svc.HistoriesByPeriod", err)
		return nil, err
	}
	return histories, nil
}

// NewClient ... リポジトリを作成する
func NewClient(fCli *firestore.Client) *Client {
	repo := NewRepository(fCli)
	svc := NewService(repo)
	return &Client{
		svc:  svc,
		fCli: fCli,
	}
}
