package wallet

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"

	"github.com/rabee-inc/go-pkg/log"
	"github.com/rabee-inc/go-pkg/timeutil"
)

// Service ... ウォレットのサービス
type Service struct {
	repo *Repository
}

// Get ... 取得する
func (s *Service) Get(ctx context.Context, userID string, kind ItemKind) (*Item, error) {
	item, err := s.repo.GetItem(ctx, userID, kind)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}
	if item == nil {
		now := timeutil.NowUnix()
		return &Item{
			UserID:    userID,
			Kind:      kind,
			Amount:    0,
			TotalGive: 0,
			TotalUse:  0,
			CreatedAt: now,
			UpdatedAt: now,
		}, nil
	}
	return item, nil
}

// GetMulti ... 複数種類を取得する
func (s *Service) GetMulti(ctx context.Context, userID string, kinds []ItemKind) (map[ItemKind]*Item, error) {
	items, err := s.repo.ListItem(ctx, userID)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}
	itemMap := map[ItemKind]*Item{}
	now := timeutil.NowUnix()
	for _, kind := range kinds {
		if item := s.getByItem(items, kind); item != nil {
			itemMap[kind] = item
		} else {
			itemMap[kind] = &Item{
				UserID:    userID,
				Kind:      kind,
				Amount:    0,
				TotalGive: 0,
				TotalUse:  0,
				CreatedAt: now,
				UpdatedAt: now,
			}
		}
	}
	return itemMap, nil
}

// GetAmount ... 金額のみ取得する
func (s *Service) GetAmount(ctx context.Context, userID string, kind ItemKind) (float64, error) {
	item, err := s.repo.GetItem(ctx, userID, kind)
	if err != nil {
		log.Error(ctx, err)
		return 0, err
	}
	if item == nil {
		return 0, nil
	}
	return item.Amount, nil
}

// GetMultiAmount ... 複数種類の金額のみ取得する
func (s *Service) GetMultiAmount(ctx context.Context, userID string, kinds []ItemKind) (map[ItemKind]float64, error) {
	items, err := s.repo.ListItem(ctx, userID)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}
	itemMap := map[ItemKind]float64{}
	for _, kind := range kinds {
		if item := s.getByItem(items, kind); item != nil {
			itemMap[kind] = item.Amount
		} else {
			itemMap[kind] = 0
		}
	}
	return itemMap, nil
}

func (s *Service) getByItem(items []*Item, kind ItemKind) *Item {
	for _, item := range items {
		if item.Kind == kind {
			return item
		}
	}
	return nil
}

// Give ... 配布する
func (s *Service) Give(
	ctx context.Context,
	tx *firestore.Transaction,
	userID string,
	amounts map[ItemKind]float64,
	data map[string]interface{},
	comment string) (map[ItemKind]*Item, error) {
	// アイテムを取得
	kinds := []ItemKind{}
	for kind := range amounts {
		kinds = append(kinds, kind)
	}
	itemsMap, err := s.repo.TxGetMultiItem(ctx, tx, userID, kinds)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}

	// アイテム詳細を更新
	updateAmounts, err := s.give(ctx, amounts)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}

	// アイテムを更新
	now := timeutil.NowUnix()
	itemsMap, err = s.updateItems(ctx, tx, userID, kinds, itemsMap, updateAmounts, now)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}

	// 履歴を記録
	for kind, amount := range amounts {
		if amount <= 0 {
			continue
		}
		_, err = s.repo.TxCreateHistory(ctx, tx, userID, kind, amount, data, comment, now)
		if err != nil {
			log.Error(ctx, err)
			return nil, err
		}
	}
	return itemsMap, nil
}

// Use ... 消費する
func (s *Service) Use(
	ctx context.Context,
	tx *firestore.Transaction,
	userID string,
	amounts map[ItemKind]float64,
	data map[string]interface{},
	comment string) (map[ItemKind]*Item, error) {
	now := timeutil.NowUnix()

	// アイテムを取得
	kinds := []ItemKind{}
	for kind := range amounts {
		kinds = append(kinds, kind)
	}
	itemsMap, err := s.repo.TxGetMultiItem(ctx, tx, userID, kinds)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}

	// アイテムを消費
	updateAmounts, err := s.use(ctx, kinds, amounts, itemsMap)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}

	// アイテムを更新
	itemsMap, err = s.updateItems(ctx, tx, userID, kinds, itemsMap, updateAmounts, now)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}

	// 履歴を記録
	for kind, amount := range amounts {
		if amount <= 0 {
			continue
		}
		_, err = s.repo.TxCreateHistory(ctx, tx, userID, kind, -amount, data, comment, now)
		if err != nil {
			log.Error(ctx, err)
			return nil, err
		}
	}
	return itemsMap, nil
}

// Exchange ... 両替する
func (s *Service) Exchange(
	ctx context.Context,
	tx *firestore.Transaction,
	userID string,
	fromKind ItemKind,
	toKind ItemKind,
	amount float64,
	data map[string]interface{},
	comment string) (map[ItemKind]*Item, error) {
	now := timeutil.NowUnix()

	// アイテムを取得
	kinds := []ItemKind{fromKind, toKind}
	itemsMap, err := s.repo.TxGetMultiItem(ctx, tx, userID, kinds)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}
	if itemsMap[fromKind] == nil {
		itemsMap[fromKind] = &Item{
			UserID:    userID,
			Kind:      fromKind,
			Amount:    0,
			TotalGive: 0,
			TotalUse:  0,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	if itemsMap[toKind] == nil {
		itemsMap[toKind] = &Item{
			UserID:    userID,
			Kind:      toKind,
			Amount:    0,
			TotalGive: 0,
			TotalUse:  0,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	if amount <= 0 || fromKind == toKind {
		return itemsMap, nil
	}

	// 消費
	fromAmount := map[ItemKind]float64{
		fromKind: amount,
	}
	useUpdateAmounts, err := s.use(ctx, kinds, fromAmount, itemsMap)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}

	// 配布
	toAmount := map[ItemKind]float64{
		toKind: amount,
	}
	giveUpdateAmounts, err := s.give(ctx, toAmount)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}

	updateAmounts := map[ItemKind]float64{}
	for kind, amount := range useUpdateAmounts {
		updateAmounts[kind] = amount
	}
	for kind, amount := range giveUpdateAmounts {
		updateAmounts[kind] = amount
	}
	itemsMap, err = s.updateItems(ctx, tx, userID, kinds, itemsMap, updateAmounts, now)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}

	// 履歴を記録
	_, err = s.repo.TxCreateHistory(ctx, tx, userID, fromKind, -amount, data, comment, now)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}
	_, err = s.repo.TxCreateHistory(ctx, tx, userID, toKind, amount, data, comment, now)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}
	return itemsMap, nil
}

// HistoriesByCursor ... 履歴リストを取得する
func (s *Service) HistoriesByCursor(ctx context.Context, userID string, kinds []ItemKind, limit int, cursor string) ([]*ItemHistory, string, error) {
	if len(kinds) == 0 || limit <= 0 {
		return []*ItemHistory{}, "", nil
	}
	histories, nCursor, err := s.repo.ListHistoryByCursor(ctx, userID, kinds, limit, cursor)
	if err != nil {
		log.Error(ctx, err)
		return nil, "", err
	}
	return histories, nCursor, nil
}

// HistoriesByPeriod ... 期間を指定して履歴リストを取得する
func (s *Service) HistoriesByPeriod(
	ctx context.Context,
	userID string,
	kinds []ItemKind,
	startAt int64,
	endAt int64) ([]*ItemHistory, error) {
	if len(kinds) == 0 || startAt > endAt {
		return []*ItemHistory{}, nil
	}
	histories, err := s.repo.ListHistoryByPeriod(ctx, userID, kinds, startAt, endAt)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}
	return histories, nil
}

func (s *Service) updateItems(ctx context.Context, tx *firestore.Transaction, userID string, kinds []ItemKind, itemsMap map[ItemKind]*Item, amounts map[ItemKind]float64, now int64) (map[ItemKind]*Item, error) {
	for _, kind := range kinds {
		if item, ok := itemsMap[kind]; ok {
			item.UpdatedAt = now
			continue
		}
		item := &Item{
			UserID:    userID,
			Kind:      kind,
			Amount:    0,
			TotalGive: 0,
			TotalUse:  0,
			CreatedAt: now,
			UpdatedAt: now,
		}
		itemsMap[kind] = item
	}

	for kind, item := range itemsMap {
		var amount float64
		var ok bool
		if amount, ok = amounts[kind]; !ok {
			continue
		}
		item.Amount += amount
		if amount > 0 {
			item.TotalGive += amount
		} else {
			item.TotalUse -= amount
		}
		var err error
		item, err = s.repo.TxSetItem(ctx, tx, userID, item.Kind, item)
		if err != nil {
			log.Error(ctx, err)
			return nil, err
		}
	}
	return itemsMap, nil
}

func (s *Service) give(ctx context.Context, amounts map[ItemKind]float64) (map[ItemKind]float64, error) {
	updateAmounts := map[ItemKind]float64{}
	for kind, amount := range amounts {
		if amount <= 0 {
			continue
		}
		updateAmounts[kind] += amount
	}
	return updateAmounts, nil
}

func (s *Service) use(ctx context.Context, kinds []ItemKind, amounts map[ItemKind]float64, itemsMap map[ItemKind]*Item) (map[ItemKind]float64, error) {
	updateAmounts := map[ItemKind]float64{}
	for kind, amount := range amounts {
		var item *Item
		var ok bool
		if item, ok = itemsMap[kind]; !ok {
			continue
		}
		if amount <= 0 {
			continue
		}
		if item.Amount < amount {
			err := log.Warningc(ctx, http.StatusBadRequest, "Not enough amount kind: %s, item_amount: %f, amount: %f", kind, item.Amount, amount)
			return nil, err
		}
		updateAmounts[kind] = -amount
	}
	return updateAmounts, nil
}

func (s *Service) getItemsByKind(items []*Item, kind ItemKind) *Item {
	for _, item := range items {
		if item.Kind == kind {
			return item
		}
	}
	return nil
}

// NewService ... サービスを作成する
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}
