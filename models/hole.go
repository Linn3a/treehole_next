package models

import (
	"time"
	"treehole_next/config"

	"gorm.io/gorm"
)

type HoleFloor struct {
	FirstFloor *Floor   `json:"first_floor"`
	LastFloor  *Floor   `json:"last_floor"`
	Floors     *[]Floor `json:"floors"`
}

type Hole struct {
	BaseModel
	DivisionID int          `json:"division_id"`
	Tags       []*Tag       `json:"tags" gorm:"many2many:hole_tags"`
	Floors     []Floor      `json:"-"`
	HoleFloor  HoleFloor    `json:"floors" gorm:"-:all"` // return floors
	View       int          `json:"view"`
	Reply      int          `json:"reply"`
	Mapping    IntStringMap `json:"-"`
	Hidden     bool         `json:"hidden"`
}
type Holes []Hole

// AfterFind set default mapping as {}
//goland:noinspection GoUnusedParameter
func (hole *Hole) AfterFind(tx *gorm.DB) (err error) {
	if hole.Mapping == nil {
		hole.Mapping = map[int]string{}
	}
	return
}

// AfterCreate set default mapping as {}
func (hole *Hole) AfterCreate(tx *gorm.DB) (err error) {
	return hole.AfterFind(tx)
}

func (hole *Hole) LoadTags() error {
	var tags []*Tag
	DB.Model(hole).Association("Tags").Find(&tags)
	if tags == nil {
		hole.Tags = []*Tag{}
	} else {
		hole.Tags = tags
	}
	return nil
}

func (hole *Hole) LoadFloors() error {
	var floors []Floor
	result := DB.Where("hole_id = ?", hole.ID).Limit(config.Config.Size).Find(&floors)
	hole.HoleFloor.Floors = &floors
	if result.RowsAffected == 0 {
		return nil
	}

	hole.HoleFloor.FirstFloor = &floors[0]

	if hole.Reply <= config.Config.Size {
		hole.HoleFloor.LastFloor = &floors[result.RowsAffected-1]
	} else {
		var floor Floor
		DB.Where("hole_id = ?", hole.ID).Last(&floor)
		hole.HoleFloor.LastFloor = &floor
	}
	return nil
}

func (hole *Hole) Preprocess() error {
	hole.LoadTags()
	return hole.LoadFloors()
}

func (holes Holes) Preprocess() error {
	// TODO: cache
	for i := 0; i < len(holes); i++ {
		if err := holes[i].Preprocess(); err != nil {
			return err
		}
	}
	return nil
}

func (hole Hole) MakeHiddenQuerySet(isAdmin bool) (tx *gorm.DB) {
	if isAdmin {
		return DB
	} else {
		return DB.Where("hidden = ?", false)
	}
}

func (hole Hole) MakeQuerySet(offset time.Time, size int, isadmin bool) (tx *gorm.DB) {
	return hole.
		MakeHiddenQuerySet(isadmin).
		Where("updated_at < ?", offset).
		Order("updated_at desc").Limit(size)
}
