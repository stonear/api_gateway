package repository

import (
	"github.com/stonear/api_gateway/db"
	"github.com/stonear/api_gateway/model"
)

type route struct {
	*db.Db
}

func NewRouteRepository(d *db.Db) *route {
	return &route{
		d,
	}
}

type Route interface {
	Get() ([]model.Route, error)
}

func (r *route) Get() ([]model.Route, error) {
	var routes []model.Route
	result := r.Find(&routes)
	if result.Error != nil {
		return nil, result.Error
	}
	return routes, nil
}
