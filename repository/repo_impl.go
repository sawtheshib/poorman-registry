package repository

import (
	"github.com/nduyphuong/reverse-registry/model"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) Interface {
	return &Storage{
		db,
	}
}

func (s *Storage) FindByNameTag(nameWithTag string) (*model.ImageModel, error) {
	var iM model.ImageModel
	query := s.db.Model(&model.ImageModel{})
	query = query.Where("name=?", nameWithTag)
	err := query.Find(&iM).Error
	if err != nil {
		return nil, err
	}
	return &iM, nil
}

func (s *Storage) FindByDigest(hashedIndex string) (*model.ImageModel, error) {
	var iM model.ImageModel
	query := s.db.Model(&model.ImageModel{})
	query = query.Where("hashed_index=?", hashedIndex)
	err := query.Find(&iM).Error
	if err != nil {
		return nil, err
	}
	return &iM, nil
}

func (s *Storage) SaveDigest(nameWithTag string, hashedIndex string) error {
	iM := model.ImageModel{
		Name:         nameWithTag,
		HashedIndex:  hashedIndex,
	}
	if err := s.db.Save(&iM).Error; err != nil {
		return err
	}
	return nil
}
