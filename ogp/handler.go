package ogp

import (
	"net/http"

	"gopkg.in/go-playground/validator.v9"

	"github.com/rabee-inc/go-pkg/errcode"
	"github.com/rabee-inc/go-pkg/parameter"
	"github.com/rabee-inc/go-pkg/renderer"
)

// Handler ... ハンドラ
type Handler struct {
	repo Repository
}

// UpdateURL ... 作成したOGPをアップデートする
func (h *Handler) UpdateURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Paramを取得
	var param struct {
		Key string `json:"key" validate:"required"`
		ID  string `json:"id"  validate:"required"`
		URL string `json:"url" validate:"required"`
	}
	err := parameter.GetJSON(r, &param)
	if err != nil {
		err = errcode.Set(err, http.StatusBadRequest)
		renderer.HandleError(ctx, w, "parameter.GetJSON", err)
		return
	}

	// Validation
	v := validator.New()
	if err := v.Struct(param); err != nil {
		err = errcode.Set(err, http.StatusBadRequest)
		renderer.HandleError(ctx, w, "v.Struct", err)
		return
	}

	err = h.repo.UpdateURL(ctx, param.Key, param.ID, param.URL)
	if err != nil {
		renderer.HandleError(ctx, w, "h.sSvc.UpdateURL", err)
		return
	}

	// Response
	renderer.Success(ctx, w)
}

// NewHandler ... ハンドラを作成する
func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}
