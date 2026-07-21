package store

import "github.com/peterkuchinov/The-link-shortener-on-Golang/internal/service"



func NewMemoryStore() service.LinkStore {
	var res service.LinkStore
	return res
}