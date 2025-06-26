package book

import (
	"github.com/take0fit/validationcontext"
	"time"
)

// Book represents a book entity
type Book struct {
	ID          BookID
	Title       Title
	Price       Price
	IsAvailable Available
	Rating      Rating
	PublishedAt PublishedAt
}

// BookID value object
type BookID struct {
	value int
}

//go:generate voauto-gen
func NewBookID(v int, vc *validationcontext.ValidationContext) BookID {
	vc.ValidateMinValue(v, "BookID", 1, "Book ID must be greater than 0")
	return BookID{value: v}
}

func (b BookID) Value() int {
	return b.value
}

// Title value object
type Title struct {
	value string
}

//go:generate voauto-gen
func NewTitle(v string, vc *validationcontext.ValidationContext) Title {
	vc.Required(v, "Title", "title is required", false)
	vc.ValidateMinLength(v, "Title", 1, "title must not be empty")
	vc.ValidateMaxLength(v, "Title", 200, "title must be 200 characters or less")
	return Title{value: v}
}

func (t Title) String() string {
	return t.value
}

// Price value object (uses float64)
type Price struct {
	value float64
}

//go:generate voauto-gen
func NewPrice(v float64, vc *validationcontext.ValidationContext) Price {
	vc.ValidateMinFloatValue(v, "Price", 0.0, "price must be non-negative")
	vc.ValidateMaxFloatValue(v, "Price", 10000.0, "price must be 10,000 or less")
	return Price{value: v}
}

func (p Price) Value() float64 {
	return p.value
}

// Available value object (uses bool)
type Available struct {
	value bool
}

//go:generate voauto-gen
func NewAvailable(v bool, vc *validationcontext.ValidationContext) Available {
	// No specific validation needed for boolean, but we could add business rules
	return Available{value: v}
}

func (a Available) Value() bool {
	return a.value
}

// Rating value object (uses int for rating 1-5)
type Rating struct {
	value int
}

//go:generate voauto-gen
func NewRating(v int, vc *validationcontext.ValidationContext) Rating {
	vc.ValidateMinValue(v, "Rating", 1, "rating must be between 1 and 5")
	vc.ValidateMaxValue(v, "Rating", 5, "rating must be between 1 and 5")
	return Rating{value: v}
}

func (r Rating) Value() int {
	return r.value
}

// PublishedAt value object
type PublishedAt struct {
	value time.Time
}

//go:generate voauto-gen
func NewPublishedAt(v string, vc *validationcontext.ValidationContext) PublishedAt {
	vc.ValidateDate(v, "PublishedAt", "invalid date format")

	parsed, err := time.Parse("2006-01-02", v)
	if err != nil {
		vc.AddError("PublishedAt", "failed to parse date")
		return PublishedAt{value: time.Time{}}
	}

	return PublishedAt{value: parsed}
}

func (p PublishedAt) Value() time.Time {
	return p.value
}

func (p PublishedAt) String() string {
	return p.value.Format("2006-01-02")
}

func NewBook(
	id int,
	title string,
	price float64,
	isAvailable bool,
	rating int,
	publishedAt string,
	vc *validationcontext.ValidationContext,
) Book {
	return Book{
		ID:          NewBookID(id, vc),
		Title:       NewTitle(title, vc),
		Price:       NewPrice(price, vc),
		IsAvailable: NewAvailable(isAvailable, vc),
		Rating:      NewRating(rating, vc),
		PublishedAt: NewPublishedAt(publishedAt, vc),
	}
}
