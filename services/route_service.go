package services

import (
	"fmt"
	"malakashuttle/dto"
	"malakashuttle/entities"
	"malakashuttle/repositories"
	"malakashuttle/utils"
	"time"
)

type RouteService interface {
	CreateRoute(req dto.RouteRequest) (*dto.RouteResponse, error)
	GetRouteByID(id uint) (*dto.RouteResponse, error)
	UpdateRoute(id uint, req dto.RouteUpdateRequest) (*dto.RouteResponse, error)
	DeleteRoute(id uint) error
	GetAllRoutes(params utils.PaginationParams) (*utils.PaginationResponse, error)
}

type routeService struct {
	routeRepo repositories.RouteRepository
}

func NewRouteService(routeRepo repositories.RouteRepository) RouteService {
	return &routeService{
		routeRepo: routeRepo,
	}
}

func (s *routeService) CreateRoute(req dto.RouteRequest) (*dto.RouteResponse, error) {
	// Check for duplicate route
	isDuplicate, err := s.routeRepo.CheckDuplicate(req.OriginCity, req.DestinationCity, nil)
	if err != nil {
		return nil, utils.NewInternalServerError("Failed to check duplicate route", err)
	}
	if isDuplicate {
		return nil, utils.NewBadRequestErrorWithDetails(
			"Route with the same origin and destination already exists",
			nil,
			req,
		)
	}

	// Validate that origin and destination are different
	if req.OriginCity == req.DestinationCity {
		return nil, utils.NewBadRequestErrorWithDetails(
			"Origin city and destination city cannot be the same",
			nil,
			req,
		)
	}

	// Create route entity
	route := &entities.Route{
		OriginCity:      req.OriginCity,
		DestinationCity: req.DestinationCity,
	}

	// Save to database
	if err := s.routeRepo.Create(route); err != nil {
		return nil, utils.NewInternalServerError("Failed to create route", err)
	}

	// Convert to response DTO
	response := &dto.RouteResponse{
		ID:              route.ID,
		OriginCity:      route.OriginCity,
		DestinationCity: route.DestinationCity,
		CreatedAt:       route.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       route.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func (s *routeService) GetRouteByID(id uint) (*dto.RouteResponse, error) {
	route, err := s.routeRepo.FindByID(id)
	if err != nil {
		return nil, utils.NewNotFoundError(fmt.Sprintf("Route with ID %d not found", id), err)
	}

	response := &dto.RouteResponse{
		ID:              route.ID,
		OriginCity:      route.OriginCity,
		DestinationCity: route.DestinationCity,
		CreatedAt:       route.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       route.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func (s *routeService) UpdateRoute(id uint, req dto.RouteUpdateRequest) (*dto.RouteResponse, error) {
	// Check if route exists
	route, err := s.routeRepo.FindByID(id)
	if err != nil {
		return nil, utils.NewNotFoundError(fmt.Sprintf("Route with ID %d not found", id), err)
	}

	// Validate that origin and destination are different
	if req.OriginCity == req.DestinationCity {
		return nil, utils.NewBadRequestErrorWithDetails(
			"Origin city and destination city cannot be the same",
			nil,
			req,
		)
	}

	// Check for duplicate route (excluding current route)
	isDuplicate, err := s.routeRepo.CheckDuplicate(req.OriginCity, req.DestinationCity, &id)
	if err != nil {
		return nil, utils.NewInternalServerError("Failed to check duplicate route", err)
	}
	if isDuplicate {
		return nil, utils.NewBadRequestErrorWithDetails(
			"Route with the same origin and destination already exists",
			nil,
			req,
		)
	}

	// Update route fields
	route.OriginCity = req.OriginCity
	route.DestinationCity = req.DestinationCity

	// Save changes to database
	if err := s.routeRepo.Update(route); err != nil {
		return nil, utils.NewInternalServerError("Failed to update route", err)
	}

	// Convert to response DTO
	response := &dto.RouteResponse{
		ID:              route.ID,
		OriginCity:      route.OriginCity,
		DestinationCity: route.DestinationCity,
		CreatedAt:       route.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       route.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func (s *routeService) DeleteRoute(id uint) error {
	// Check if route exists
	_, err := s.routeRepo.FindByID(id)
	if err != nil {
		return utils.NewNotFoundError(fmt.Sprintf("Route with ID %d not found", id), err)
	}

	// Delete route
	if err := s.routeRepo.Delete(id); err != nil {
		return utils.NewInternalServerError("Failed to delete route", err)
	}

	return nil
}

func (s *routeService) GetAllRoutes(params utils.PaginationParams) (*utils.PaginationResponse, error) {
	routes, total, err := s.routeRepo.FindAll(params)
	if err != nil {
		return nil, utils.NewInternalServerError("Failed to get routes", err)
	}

	// Convert to response DTOs
	var routeResponses []dto.RouteResponse
	for _, route := range routes {
		routeResponses = append(routeResponses, dto.RouteResponse{
			ID:              route.ID,
			OriginCity:      route.OriginCity,
			DestinationCity: route.DestinationCity,
			CreatedAt:       route.CreatedAt.Format(time.RFC3339),
			UpdatedAt:       route.UpdatedAt.Format(time.RFC3339),
		})
	}

	// Use the existing pagination response builder
	response := utils.CreatePaginationResponse(routeResponses, total, params)

	return &response, nil
}
