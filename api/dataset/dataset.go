package dataset

import (
	"fmt"
	"github.com/adamydwang/aiops/api/dataset/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"path"
)

type DatasetCenter struct {
	db       *gorm.DB
	rootPath string
}

func (center *DatasetCenter) Init(uri string, rootPath string) error {
	var err error = nil
	center.db, err = gorm.Open(mysql.Open(uri), &gorm.Config{})
	center.rootPath = rootPath
	if err != nil {
		return err
	}
	if err = center.db.AutoMigrate(&dao.Dataset{}); err != nil {
		return err
	}
	if s, err := os.Stat(rootPath); err != nil {
		return err
	} else if !s.IsDir() {
		return fmt.Errorf("rootPath not exits: %s", rootPath)
	}
	return nil
}

func (center *DatasetCenter) CreateDataset(name, creator string) error {
	return center.db.Transaction(func(tx *gorm.DB) error {
		tx = tx.Create(&dao.Dataset{
			Name:    name,
			Creator: creator,
			Status:  dao.StatusCreated,
		})
		if tx.Error != nil {
			return tx.Error
		}
		if err := os.Mkdir(path.Join(center.rootPath, name), os.ModePerm); err != nil {
			return err
		}
		return nil
	})
}

func (center *DatasetCenter) DeleteDataset(name string) error {
	datasets := []dao.Dataset{}
	tx := center.db.Where(&dao.Dataset{Name: name, Status: dao.StatusCreated}).Find(&datasets)
	if len(datasets) == 0 {
		return nil
	}
	tx = tx.Model(&dao.Dataset{Name: name}).Updates(dao.Dataset{Status: dao.StatusDeleted})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (center *DatasetCenter) GetDatasetInfo(name string) (dao.Dataset, error) {
	ds := dao.Dataset{}
	tx := center.db.Where(&dao.Dataset{Name: name, Status: dao.StatusCreated}).First(&ds)
	if tx.Error != nil {
		return ds, tx.Error
	}
	return ds, nil
}

func (center *DatasetCenter) ListDatasets(
	offset, limit int, orderBy string, desc bool, name string) ([]dao.Dataset, error) {
	tx := center.db
	if name != "" {
		tx = center.db.Where(&dao.Dataset{Name: name, Status: dao.StatusCreated})
	}
	orderCond := fmt.Sprintf("%s", orderBy)
	if desc {
		orderCond += fmt.Sprintf(" desc")
	}
	datasets := make([]dao.Dataset)
	tx = tx.Order(orderCond).Offset(offset).Limit(limit).Find(&datasets)
	if tx.Error != nil {
		return datasets, tx.Error
	}
	return datasets, nil
}
