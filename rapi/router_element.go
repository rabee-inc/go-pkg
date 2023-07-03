package rapi

import (
	"context"
	"net/http"
)

type RouterElement interface {
	GetHandleFunc() http.HandlerFunc
	GetEmptyInput() any
	GetEmptyOutput() any
}

type HandlerMethod[I any] interface {
	RouterElement
	// 共通のリクエストパラメーター受け取り処理をセット
	SetInputFunc(func(ctx context.Context, r *http.Request, param any) error)
	// 共通のバリデーション処理をセット
	SetValidateFunc(func(ctx context.Context, param any) error)
	// エラーをレンダリングする処理をセット
	SetHandleErrorFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error))
	// レスポンスをレンダリングする処理をセット
	SetRenderFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request, output any))
	// エラーをレンダリングする直前にエラーを書き換える処理をセット
	BeforeHandleError(func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) error)
	// 共通リクエストパラメーター受け取り処理の後に必要な処理があればセット
	AfterInput(func(ctx context.Context, r *http.Request, param *I) error)
	// 共通のバリデーション処理の後に必要な処理があればセット
	AfterValidate(func(ctx context.Context, param *I) error)
}
