package mysql

import (
	"errors"
	"github.com/jinzhu/gorm"
	"sync"
)

type sqlPool struct {
	new func() *gorm.DB
	db  []*gorm.DB
	sync.Mutex
}

func newPool(newDb func() *gorm.DB, size int) *sqlPool {
	return &sqlPool{newDb, make([]*gorm.DB, 0, size), sync.Mutex{}}
}

func (s *sqlPool) get() (db *gorm.DB, err error) {
	s.Lock()
	defer s.Unlock()
	//log.Printf("before len:%d", len(s.db))
	if len(s.db) > 0 {
		db = s.db[len(s.db)-1]
		s.db = s.db[:len(s.db)-1]
	} else {
		db = s.new()
	}
	if db == nil {
		return nil, errors.New("db get err")
	}
	if err = db.DB().Ping(); err != nil {
		return nil, err
	}
	//log.Printf("after len:%d", len(s.db))
	return db, nil
}

func (s *sqlPool) put(db *gorm.DB) {
	if db == nil {
		return
	}
	if err := db.DB().Ping(); err != nil {
		return
	}
	s.Lock()
	defer s.Unlock()
	if len(s.db) < cap(s.db) {
		s.db = append(s.db, db)
	} else {
		_ = db.Close()
	}
}
