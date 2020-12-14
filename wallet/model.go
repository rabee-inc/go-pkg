package wallet

import (
	"fmt"

	"cloud.google.com/go/firestore"

	"github.com/rabee-inc/go-pkg/encryptutil"
)

// Item ... お財布のアイテム
type Item struct {
	ID        string                 `firestore:"-"          json:"id" cloudfirestore:"id"`
	Ref       *firestore.DocumentRef `firestore:"-"          json:"-"  cloudfirestore:"ref"`
	UserID    string                 `firestore:"user_id"    json:"user_id"`
	Kind      ItemKind               `firestore:"kind"       json:"kind"`
	Amount    float64                `firestore:"amount"     json:"amount"`
	TotalGive float64                `firestore:"total_give" json:"total_give"`
	TotalUse  float64                `firestore:"total_use"  json:"total_use"`
	CreatedAt int64                  `firestore:"created_at" json:"created_at"`
	UpdatedAt int64                  `firestore:"updated_at" json:"updated_at"`
}

// ItemRef ... コレクションの参照を取得
func ItemRef(fCli *firestore.Client) *firestore.CollectionRef {
	return fCli.Collection("wallet_items")
}

// GenerateItemID ... お財布のアイテムIDを作成する
func GenerateItemID(userID string, kind ItemKind) string {
	return encryptutil.ToMD5(fmt.Sprintf("%s::%s", userID, kind))
}

// ItemDetail ... お財布アイテムの詳細
type ItemDetail struct {
	ID        string                 `firestore:"-"          json:"id" cloudfirestore:"id"`
	Ref       *firestore.DocumentRef `firestore:"-"          json:"-"  cloudfirestore:"ref"`
	UserID    string                 `firestore:"user_id"    json:"user_id"`
	Kind      ItemKind               `firestore:"kind"       json:"kind"`
	Amount    float64                `firestore:"amount"     json:"amount"`
	Expired   bool                   `firestore:"expired"    json:"expired"`
	CreatedAt int64                  `firestore:"created_at" json:"created_at"`
	UpdatedAt int64                  `firestore:"updated_at" json:"updated_at"`
}

// ItemDetailRef ... コレクションの参照を取得
func ItemDetailRef(fCli *firestore.Client) *firestore.CollectionRef {
	return fCli.Collection("wallet_item_details")
}

// ItemHistory ... お財布アイテムの配布/消費履歴
type ItemHistory struct {
	ID        string                 `firestore:"-"          json:"id" cloudfirestore:"id"`
	Ref       *firestore.DocumentRef `firestore:"-"          json:"-"  cloudfirestore:"ref"`
	UserID    string                 `firestore:"user_id"    json:"user_id"`
	Kind      ItemKind               `firestore:"kind"       json:"kind"`
	Amount    float64                `firestore:"amount"     json:"amount"`
	Data      map[string]interface{} `firestore:"data"       json:"data"`
	Comment   string                 `firestore:"comment"    json:"comment"`
	CreatedAt int64                  `firestore:"created_at" json:"created_at"`
}

// ItemHistoryRef ... コレクションの参照を取得
func ItemHistoryRef(fCli *firestore.Client) *firestore.CollectionRef {
	return fCli.Collection("wallet_item_histories")
}
