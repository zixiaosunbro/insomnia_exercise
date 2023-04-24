package mysql

import "gorm.io/gorm"

type DBOptions struct {
	UserMaster bool
	ForUpdate  bool
	Tx         *gorm.DB
}

type DBOption func(s *DBOptions)

func (o *DBOptions) Apply(fs ...DBOption) {
	for _, f := range fs {
		f(o)
	}
}

func NewDBOptions() *DBOptions {
	return &DBOptions{}
}

// WithMaster using mysql mater instance
func WithMaster() DBOption {
	return func(p *DBOptions) {
		p.UserMaster = true
	}
}

// WithForUpdate select using for update
func WithForUpdate() DBOption {
	return func(p *DBOptions) {
		p.ForUpdate = true
	}
}

// WithTx pass transaction when function call, commit at caller
func WithTx(tx *gorm.DB) DBOption {
	return func(p *DBOptions) {
		p.Tx = tx
	}
}
