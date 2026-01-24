package domain

import "errors"

var (
	ErrMovieNotFound     = errors.New("movie not found")
	ErrInvalidMovie      = errors.New("invalid movie")
	ErrInvalidRating     = errors.New("rating should be between 0 and 10")
	ErrInvalidMovieTitle = errors.New("Title length should not exceed 40 chars")
)
