package filters

import (
	"ielts/internal/domain"
	"ielts/internal/validator"
	"math"
	"strings"
)

type Metadata struct {
	CurrentPage  int
	PageSize     int
	FirstPage    int
	LastPage     int
	TotalRecords int
}
type PhotoSearch struct {
	Title    string
	Author   string
	Category string
	Tags     string
}

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

func (f Filters) Limit() int {
	return f.PageSize
}
func (f Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	v.Check(validator.In(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}

func ValidateRating(v *validator.Validator, rating *domain.Rating) {
	v.Check(rating.PhotoID > 0, "photo_id", "must be a positive integer")
	v.Check(rating.UserID > 0, "user_id", "must be a positive integer")
	v.Check(rating.Score >= 1 && rating.Score <= 5, "score", "must be between 1 and 5")
}

func ValidateComment(v *validator.Validator, comment *domain.Comment) {
	v.Check(comment.PhotoID > 0, "photo_id", "must be a positive integer")
	v.Check(comment.UserID > 0, "user_id", "must be a positive integer")
	v.Check(comment.Content != "", "content", "must not be empty")
	v.Check(len(comment.Content) <= 1000, "content", "must not be more than 1000 bytes long")
}

func ValidateLike(v *validator.Validator, like *domain.Like) {
	v.Check(like.PhotoID > 0, "photo_id", "must be a positive integer")
}

func (f Filters) SortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + f.Sort)
}
func (f Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func CalculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}
	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
