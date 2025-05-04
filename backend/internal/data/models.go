package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict raised")
)

type Models struct {
	Tokens      TokenModel
	Permissions PermissionModel
	Users       UserModel
	Photo       PhotoModel
	Rating      RatingModel
	Comment     CommentModel
	Like        LikeModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Tokens:      TokenModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Users:       UserModel{DB: db},
		Photo:       PhotoModel{DB: db},
		Rating:      RatingModel{DB: db},
		Comment:     CommentModel{DB: db},
		Like:        LikeModel{DB: db},
	}
}
