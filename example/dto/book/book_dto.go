// example/dto/book/book_dto.go
package book

import (
	"github.com/take0fit/validationcontext/example/domain/book/entity/book"
	"github.com/take0fit/validationcontext/voauto"
)

// RequestCreateBook - JSONから来るリクエスト（型変換テスト用）
type RequestCreateBook struct {
	ID          int32   `json:"id"`           // int32 → int に変換
	Title       string  `json:"title"`        // string → string (そのまま)
	Price       int     `json:"price"`        // int → float64 に変換
	IsAvailable string  `json:"is_available"` // "true"/"false" → bool に変換
	Rating      float64 `json:"rating"`       // float64 → int に変換
	PublishedAt string  `json:"published_at"` // string → string (そのまま)
}

// CreateBookDTO - バリデーション済みDTO
type CreateBookDTO struct {
	ID          book.BookID      `vctag:"auto:NewBookID,ID"`
	Title       book.Title       `vctag:"auto:NewTitle,Title"`
	Price       book.Price       `vctag:"auto:NewPrice,Price"`
	IsAvailable book.Available   `vctag:"auto:NewAvailable,IsAvailable"`
	Rating      book.Rating      `vctag:"auto:NewRating,Rating"`
	PublishedAt book.PublishedAt `vctag:"auto:NewPublishedAt,PublishedAt"`
}

func NewCreateBookDTO(req *RequestCreateBook) (*CreateBookDTO, error) {
	return voauto.BindAndValidate[CreateBookDTO](req)
}

// RequestUpdateBook - 部分更新用（オプショナルフィールドのテスト）
type RequestUpdateBook struct {
	Title       *string  `json:"title,omitempty"`
	Price       *float32 `json:"price,omitempty"`        // float32 → float64 に変換
	IsAvailable *int     `json:"is_available,omitempty"` // 0/1 → bool に変換
	Rating      *string  `json:"rating,omitempty"`       // "3" → int に変換
}

// UpdateBookDTO
type UpdateBookDTO struct {
	Title       *book.Title     `vctag:"auto:NewTitle,Title"`
	Price       *book.Price     `vctag:"auto:NewPrice,Price"`
	IsAvailable *book.Available `vctag:"auto:NewAvailable,IsAvailable"`
	Rating      *book.Rating    `vctag:"auto:NewRating,Rating"`
}

func NewUpdateBookDTO(req *RequestUpdateBook) (*UpdateBookDTO, error) {
	return voauto.BindAndValidate[UpdateBookDTO](req)
}
