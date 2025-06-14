package repositories

import (
	"malakashuttle/entities"
	"malakashuttle/utils"

	"gorm.io/gorm"
)

type RouteRepository interface {
	Create(route *entities.Route) error
	FindByID(id uint) (*entities.Route, error)
	Update(route *entities.Route) error
	Delete(id uint) error
	FindAll(params utils.PaginationParams) ([]entities.Route, int64, error)
	CheckDuplicate(originCity, destinationCity string, excludeID *uint) (bool, error)
}

type routeRepository struct {
	db *gorm.DB
}

func NewRouteRepository(db *gorm.DB) RouteRepository {
	return &routeRepository{db: db}
}

// Create creates a new route
func (r *routeRepository) Create(route *entities.Route) error {
	return r.db.Create(route).Error
}

// FindByID finds a route by ID
func (r *routeRepository) FindByID(id uint) (*entities.Route, error) {
	var route entities.Route
	err := r.db.First(&route, id).Error
	if err != nil {
		return nil, err
	}
	return &route, nil
}

// Update updates a route
func (r *routeRepository) Update(route *entities.Route) error {
	return r.db.Save(route).Error
}

// Delete deletes a route
func (r *routeRepository) Delete(id uint) error {
	return r.db.Delete(&entities.Route{}, id).Error
}

// FindAll finds all routes with pagination
func (r *routeRepository) FindAll(params utils.PaginationParams) ([]entities.Route, int64, error) {
	var routes []entities.Route
	var total int64

	// Count total records
	if err := r.db.Model(&entities.Route{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.Scopes(utils.Paginate(params)).Find(&routes).Error
	if err != nil {
		return nil, 0, err
	}

	return routes, total, nil
}

// CheckDuplicate checks if a route with the same origin and destination already exists
func (r *routeRepository) CheckDuplicate(originCity, destinationCity string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Model(&entities.Route{}).
		Where("origin_city = ? AND destination_city = ?", originCity, destinationCity)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}
