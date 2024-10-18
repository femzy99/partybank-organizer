package repositories

import (
	"errors"
	"fmt"
	"github.com/djfemz/organizer-service/partybank-app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var db *gorm.DB

type Repository[T, U any] interface {
	crudRepository[T, U]
}

type crudRepository[T, U any] interface {
	Save(t *T) (*T, error)
	FindById(id U) (*T, error)
	FindAll() ([]*T, error)
	FindAllBy(pageable Pageable) ([]*T, error)
	DeleteById(id U) error
}

type repositoryImpl[T, U any] struct {
	Db *gorm.DB
}

func (r *repositoryImpl[T, U]) Save(t *T) (*T, error) {
	err := r.Db.Save(t).Error
	if err != nil {
		return nil, err
	}
	var id, _ = GetId(*t)
	err = r.Db.Preload(clause.Associations).Where("ID=?", id).First(t).Error
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *repositoryImpl[T, U]) FindById(id U) (*T, error) {
	var t = new(T)
	err := r.Db.Preload(clause.Associations).Where("ID=?", id).First(t).Error
	if err != nil {
		log.Println("err: ", err.Error())
		return nil, err
	}
	return t, nil
}

func (r *repositoryImpl[T, U]) FindAll() ([]*T, error) {
	var organizations []*T
	err := r.Db.Find(&organizations).Error
	if err != nil {
		return nil, err
	}
	return organizations, nil
}

func (r *repositoryImpl[T, U]) FindAllBy(pageable Pageable) ([]*T, error) {
	page := getPage[T](r.Db, pageable)
	if page == nil {
		return nil, errors.New("failed to find page")
	}
	return page.GetElements(), nil
}

func (r *repositoryImpl[T, U]) DeleteById(id U) error {
	err := r.Db.Delete(new(T), id).Error

	if err != nil {
		return err
	}
	return nil
}

func Connect() *gorm.DB {
	if db != nil {
		return db
	}
	port, err := strconv.ParseUint(os.Getenv("DATABASE_PORT"), 10, 64)
	if err != nil {
		log.Fatal("error reading port: ", err)
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d TimeZone=Africa/Lagos", os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_USERNAME"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_NAME"), port)
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.Organizer{}, &models.Event{}, &models.EventStaff{},
		&models.Ticket{}, &models.Series{}, &models.Attendee{})
	if err != nil {
		log.Fatal("error migrating: ", err)
	}
	return db
}

func addEntities(m map[string]any, db *gorm.DB) error {
	for _, v := range m {
		err := db.AutoMigrate(v)
		if err != nil {
			log.Fatal("error migrating: ", err)
		}
	}

	return nil
}

func GetId(T any) (any, error) {
	obj := reflect.ValueOf(T)
	numberOfFields := obj.NumField()
	for index := ZERO; index < numberOfFields; index++ {
		tag := obj.Type().Field(index).Tag
		isPrimaryKeyField := strings.Contains(string(tag), "id")
		if isPrimaryKeyField {
			id, err, done := getPrimaryKey(obj, index)
			if done {
				return id, err
			}
		}
	}
	return nil, errors.New("could not retrieve id")
}

func getPrimaryKey(obj reflect.Value, index int) (any, error, bool) {
	idField := obj.Field(index)
	if idField.CanConvert(reflect.TypeOf(uint64(5))) {
		v := idField.Uint()
		return v, nil, true
	} else if idField.CanConvert(reflect.TypeOf(5)) {
		v := idField.Int()
		return v, nil, true
	} else if idField.CanConvert(reflect.TypeOf("")) {
		return idField.String(), nil, true
	}
	return nil, nil, false
}
