package wallet

import (
	"context"

	"cloud.google.com/go/firestore"

	"github.com/rabee-inc/go-pkg/cloudfirestore"
	"github.com/rabee-inc/go-pkg/log"
)

// Repository ... ウォレットのリポジトリ
type Repository struct {
	cFirestore *firestore.Client
}

// Item

// GetItem ... アイテムを取得する
func (r *Repository) GetItem(
	ctx context.Context,
	userID string,
	kind ItemKind,
) (*Item, error) {
	q := ItemRef(r.cFirestore).
		Where("user_id", "==", userID).
		Where("kind", "==", kind)
	dst := &Item{}
	exist, err := cloudfirestore.GetByQuery(ctx, q, dst)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	return dst, nil
}

// GetMultiItem ... アイテムを複数取得する
func (r *Repository) GetMultiItem(
	ctx context.Context,
	userID string,
	kinds []ItemKind,
) (map[ItemKind]*Item, error) {
	docRefs := []*firestore.DocumentRef{}
	for _, kind := range kinds {
		id := GenerateItemID(userID, kind)
		docRef := ItemRef(r.cFirestore).Doc(id)
		docRefs = append(docRefs, docRef)
	}
	items := []*Item{}
	err := cloudfirestore.GetMulti(ctx, r.cFirestore, docRefs, &items)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}
	dsts := map[ItemKind]*Item{}
	for _, item := range items {
		dsts[item.Kind] = item
	}
	return dsts, nil
}

// ListItem ... アイテムリストを取得する
func (r *Repository) ListItem(
	ctx context.Context,
	userID string,
) ([]*Item, error) {
	q := ItemRef(r.cFirestore).Where("user_id", "==", userID)
	dsts := []*Item{}
	err := cloudfirestore.ListByQuery(ctx, q, &dsts)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}
	return dsts, nil
}

// SetItem ... アイテムを設定する
func (r *Repository) SetItem(
	ctx context.Context,
	userID string,
	kind ItemKind,
	src *Item,
) (*Item, error) {
	id := GenerateItemID(userID, kind)
	docRef := ItemRef(r.cFirestore).Doc(id)
	err := cloudfirestore.Set(ctx, docRef, src)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}
	return src, nil
}

// ItemHistory

// ListHistoryByCursor ... 履歴リストを取得する
func (r *Repository) ListHistoryByCursor(
	ctx context.Context,
	userID string,
	kinds []ItemKind,
	limit int,
	cursor string,
) ([]*ItemHistory, string, error) {
	q := ItemHistoryRef(r.cFirestore).
		Where("user_id", "==", userID).
		Where("kind", "in", kinds).
		OrderBy("created_at", firestore.Desc)
	var dsnp *firestore.DocumentSnapshot
	var err error
	if cursor != "" {
		dsnp, err = ItemHistoryRef(r.cFirestore).Doc(cursor).Get(ctx)
		if err != nil {
			log.Error(ctx, err)
			return nil, "", err
		}
	}
	dsts := []*ItemHistory{}
	nDsnp, err := cloudfirestore.ListByQueryCursor(ctx, q, limit, dsnp, &dsts)
	if err != nil {
		log.Error(ctx, err)
		return nil, "", err
	}
	var nCursor string
	if nDsnp != nil {
		nCursor = nDsnp.Ref.ID
	}
	return dsts, nCursor, nil
}

// ListHistoryByPeriod ... 履歴リストを期間指定で取得する
func (r *Repository) ListHistoryByPeriod(
	ctx context.Context,
	userID string,
	kinds []ItemKind,
	startAt int64,
	endAt int64,
) ([]*ItemHistory, error) {
	q := ItemHistoryRef(r.cFirestore).
		Where("user_id", "==", userID).
		Where("kind", "in", kinds).
		Where("created_at", ">=", startAt).
		Where("created_at", "<=", endAt)
	dsts := []*ItemHistory{}
	err := cloudfirestore.ListByQuery(ctx, q, &dsts)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}
	return dsts, nil
}

// CreateHistory ... 履歴を作成する
func (r *Repository) CreateHistory(
	ctx context.Context,
	userID string,
	kind ItemKind,
	amount float64,
	data map[string]interface{},
	comment string,
	createdAt int64) (*ItemHistory, error) {
	src := &ItemHistory{
		UserID:    userID,
		Kind:      kind,
		Amount:    amount,
		Comment:   comment,
		CreatedAt: createdAt,
	}
	if data == nil {
		src.Data = map[string]interface{}{}
	} else {
		src.Data = data
	}
	colRef := ItemHistoryRef(r.cFirestore)
	err := cloudfirestore.Create(ctx, colRef, src)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}
	return src, nil
}

// NewRepository ... リポジトリを作成する
func NewRepository(cFirestore *firestore.Client) *Repository {
	return &Repository{
		cFirestore: cFirestore,
	}
}
