package postgres

import "github.com/everstake/solana-pools/internal/dao/dmodels"

type SlotTimeCondition struct {
	Limit  int
	Offset int
}

func (db *DB) GetSlotTime(cond *SlotTimeCondition) ([]*dmodels.SlotTime, error) {
	var st []*dmodels.SlotTime
	d := db.DB
	if cond.Limit > 0 {
		d = d.Limit(cond.Limit)
	}
	if cond.Offset > 0 {
		d = d.Offset(cond.Offset)
	}

	d = d.Order("created_at DESC")

	return st, d.Find(&st).Error
}

func (db *DB) CreateSlotTime(slotTime ...*dmodels.SlotTime) error {
	if len(slotTime) == 0 {
		return nil
	}
	return db.Create(&slotTime).Error
}
