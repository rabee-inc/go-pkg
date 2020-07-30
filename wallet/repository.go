package wallet

import (
	"context"
	"sort"

	"cloud.google.com/go/firestore"

	"github.com/rabee-inc/go-pkg/cloudfirestore"
	"github.com/rabee-inc/go-pkg/log"
)

type Repository struct {
	fCli *firestore.Client
}

// Item

func (r *Repository) GetItem(ctx context.Context, userID string, kind ItemKind) (*Item, error) {
	q := ItemRef(r.fCli).
		Where("user_id", "==", userID).
		Where("kind", "==", kind)
	dst := &Item{}
	exist, err := cloudfirestore.GetByQuery(ctx, q, dst)
	if err != nil {
		log.Errorm(ctx, "cloudfirestore.GetByQuery", err)
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	return dst, nil
}

func (r *Repository) TxGetMultiItem(
	ctx context.Context,
	tx *firestore.Transaction,
	userID string,
	kinds []ItemKind) (map[ItemKind]*Item, error) {
	docRefs := []*firestore.DocumentRef{}
	for _, kind := range kinds {
		id := GenerateItemID(userID, kind)
		docRef := ItemRef(r.fCli).Doc(id)
		docRefs = append(docRefs, docRef)
	}
	items := []*Item{}
	err := cloudfirestore.TxGetMulti(ctx, tx, docRefs, &items)
	if err != nil {
		log.Errorm(ctx, "cloudfirestore.TxGetMulti", err)
		return nil, err
	}
	dsts := map[ItemKind]*Item{}
	for _, item := range items {
		dsts[item.Kind] = item
	}
	return dsts, nil
}

func (r *Repository) ListItem(ctx context.Context, userID string) ([]*Item, error) {
	q := ItemRef(r.fCli).Where("user_id", "==", userID)
	dsts := []*Item{}
	err := cloudfirestore.ListByQuery(ctx, q, &dsts)
	if err != nil {
		log.Errorm(ctx, "cloudfirestore.ListByQuery", err)
		return nil, err
	}
	return dsts, nil
}

func (r *Repository) TxSetItem(
	ctx context.Context,
	tx *firestore.Transaction,
	userID string,
	kind ItemKind,
	src *Item) (*Item, error) {
	id := GenerateItemID(userID, kind)
	docRef := ItemRef(r.fCli).Doc(id)
	err := cloudfirestore.TxSet(ctx, tx, docRef, src)
	if err != nil {
		log.Errorm(ctx, "cloudfirestore.TxSet", err)
		return nil, err
	}
	return src, nil
}

// ItemDetail

func (r *Repository) TxListItemDetail(
	ctx context.Context,
	tx *firestore.Transaction,
	userID string,
	kinds []ItemKind) ([]*ItemDetail, error) {
	q := ItemDetailRef(r.fCli).
		Where("user_id", "==", userID).
		Where("kind", "in", kinds).
		Where("expired", "==", false).
		Where("amount", ">", 0)
	dsts := []*ItemDetail{}
	err := cloudfirestore.TxListByQuery(ctx, tx, q, &dsts)
	if err != nil {
		log.Errorm(ctx, "loudfirestore.TxListByQuery", err)
		return nil, err
	}
	sort.Slice(dsts, func(i, j int) bool {
		return dsts[i].CreatedAt < dsts[j].CreatedAt
	})
	return dsts, nil
}

func (r *Repository) TxCreateItemDetail(
	ctx context.Context,
	tx *firestore.Transaction,
	userID string,
	kind ItemKind,
	amount float64,
	createdAt int64) (*ItemDetail, error) {
	src := &ItemDetail{
		UserID:    userID,
		Kind:      kind,
		Amount:    amount,
		Expired:   false,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
	colRef := ItemDetailRef(r.fCli)
	err := cloudfirestore.TxCreate(ctx, tx, colRef, src)
	if err != nil {
		log.Errorm(ctx, "cloudfirestore.TxCreate", err)
		return nil, err
	}
	return src, nil
}

func (r *Repository) TxUpdateItemDetail(
	ctx context.Context,
	tx *firestore.Transaction,
	src *ItemDetail) (*ItemDetail, error) {
	docRef := ItemDetailRef(r.fCli).Doc(src.ID)
	err := cloudfirestore.TxSet(ctx, tx, docRef, src)
	if err != nil {
		log.Errorm(ctx, "cloudfirestore.TxSet", err)
		return nil, err
	}
	return src, nil
}

// ItemHistory

func (r *Repository) ListHistoryByCursor(
	ctx context.Context,
	userID string,
	kinds []ItemKind,
	limit int,
	cursor string) ([]*ItemHistory, string, error) {
	q := ItemHistoryRef(r.fCli).
		Where("user_id", "==", userID).
		Where("kind", "in", kinds).
		OrderBy("created_at", firestore.Desc)
	var dsnp *firestore.DocumentSnapshot
	var err error
	if cursor != "" {
		dsnp, err = ItemHistoryRef(r.fCli).Doc(cursor).Get(ctx)
		if err != nil {
			log.Errorm(ctx, "Get", err)
			return nil, "", err
		}
	}
	dsts := []*ItemHistory{}
	nDsnp, err := cloudfirestore.ListByQueryCursor(ctx, q, limit, dsnp, &dsts)
	if err != nil {
		log.Errorm(ctx, "cloudfirestore.ListByQueryCursor", err)
		return nil, "", err
	}
	var nCursor string
	if nDsnp != nil {
		nCursor = nDsnp.Ref.ID
	}
	return dsts, nCursor, nil
}

func (r *Repository) ListHistoryByPeriod(
	ctx context.Context,
	userID string,
	kinds []ItemKind,
	startAt int64,
	endAt int64) ([]*ItemHistory, error) {
	q := ItemHistoryRef(r.fCli).
		Where("user_id", "==", userID).
		Where("kind", "in", kinds).
		Where("created_at", ">=", startAt).
		Where("created_at", "<=", endAt)
	dsts := []*ItemHistory{}
	err := cloudfirestore.ListByQuery(ctx, q, &dsts)
	if err != nil {
		log.Errorm(ctx, "cloudfirestore.ListByQuery", err)
		return nil, err
	}
	return dsts, nil
}

func (r *Repository) TxCreateHistory(
	ctx context.Context,
	tx *firestore.Transaction,
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
	colRef := ItemHistoryRef(r.fCli)
	err := cloudfirestore.TxCreate(ctx, tx, colRef, src)
	if err != nil {
		log.Errorm(ctx, "cloudfirestore.TxCreate", err)
		return nil, err
	}
	return src, nil
}

// NewRepository ... リポジトリを作成する
func NewRepository(fCli *firestore.Client) *Repository {
	return &Repository{
		fCli: fCli,
	}
}
