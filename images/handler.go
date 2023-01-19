package images

import (
	"net/http"

	"github.com/rabee-inc/go-pkg/errcode"
	"github.com/rabee-inc/go-pkg/parameter"
	"github.com/rabee-inc/go-pkg/renderer"
	"gopkg.in/go-playground/validator.v9"
)

type Handler struct {
	repo Repository
	v    *validator.Validate
}

func NewHandler(repo Repository) *Handler {
	v := validator.New()
	return &Handler{
		repo: repo,
		v:    v,
	}
}

// 変換後の画像をアップデートする
func (h *Handler) UpdateByConvertObjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var param struct {
		Key     string    `json:"key"     validate:"required"`
		Objects []*Object `json:"objects" validate:"required"`
	}
	err := parameter.GetJSON(r, &param)
	if err != nil {
		err = errcode.Set(err, http.StatusBadRequest)
		renderer.HandleError(ctx, w, err)
		return
	}

	if err := h.v.Struct(param); err != nil {
		err = errcode.Set(err, http.StatusBadRequest)
		renderer.HandleError(ctx, w, err)
		return
	}

	err = h.repo.UpdateByConvertObjects(ctx, param.Key, param.Objects)
	if err != nil {
		renderer.HandleError(ctx, w, err)
		return
	}

	renderer.Success(ctx, w)
}

// 作成したOGPをアップデートする
func (h *Handler) UpdateByGenerateURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var param struct {
		Key string `json:"key" validate:"required"`
		ID  string `json:"id"  validate:"required"`
		URL string `json:"url" validate:"required"`
	}
	err := parameter.GetJSON(r, &param)
	if err != nil {
		err = errcode.Set(err, http.StatusBadRequest)
		renderer.HandleError(ctx, w, err)
		return
	}

	if err := h.v.Struct(param); err != nil {
		err = errcode.Set(err, http.StatusBadRequest)
		renderer.HandleError(ctx, w, err)
		return
	}

	err = h.repo.UpdateByGenerateURL(ctx, param.Key, param.ID, param.URL)
	if err != nil {
		renderer.HandleError(ctx, w, err)
		return
	}

	renderer.Success(ctx, w)
}
