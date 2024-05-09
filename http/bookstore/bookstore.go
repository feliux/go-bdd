package bookstore

import (
	"errors"
	"fmt"
)

type Store struct {
	books []*Book
}
type Book struct {
	Id     string
	Title  string
	Author string
}

func (s *Store) GetBook(id string) (book *Book, err error) {
	for _, book := range s.books {
		if book.Id == id {
			return book, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Book with id %s not found", id))
}

func (s *Store) AddBook(book *Book) (err error) {
	if s.BookExists(book.Id) {
		return errors.New(fmt.Sprintf("Book with id %s exists", book.Id))
	}
	s.books = append(s.books, book)
	return
}

func (s *Store) BookExists(id string) bool {
	_, err := s.GetBook(id)
	if err != nil {
		return false
	}
	return true
}
