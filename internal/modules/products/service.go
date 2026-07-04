package products

import (
	"context"
	"errors"
	"strings"
)

var ErrInvalidStatus = errors.New("invalid status")

type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) ListPublicProducts(ctx context.Context, query ListProductsQuery) (ListProductsResponse, error) {
	query.Public = true
	return s.repo.ListProducts(ctx, query)
}

func (s Service) ListAdminProducts(ctx context.Context, query ListProductsQuery) (ListProductsResponse, error) {
	if !ValidateProductStatus(query.Status) {
		return ListProductsResponse{}, ErrInvalidStatus
	}
	return s.repo.ListProducts(ctx, query)
}

func (s Service) FindPublicProductBySlug(ctx context.Context, slug string) (ProductResponse, error) {
	product, err := s.repo.FindPublicProductBySlug(ctx, strings.TrimSpace(slug))
	if err != nil {
		return ProductResponse{}, err
	}
	return ProductResponse{Product: product}, nil
}

func (s Service) FindAdminProductByID(ctx context.Context, id uint64) (ProductResponse, error) {
	product, err := s.repo.FindProductByID(ctx, id)
	if err != nil {
		return ProductResponse{}, err
	}
	return ProductResponse{Product: product}, nil
}

func (s Service) CreateProduct(ctx context.Context, request CreateProductRequest, actorID uint64) (ProductResponse, error) {
	normalizeProductRequest(&request)
	product, err := s.repo.CreateProduct(ctx, request, actorID)
	if err != nil {
		return ProductResponse{}, err
	}
	return ProductResponse{Product: product}, nil
}

func (s Service) UpdateProduct(ctx context.Context, id uint64, request UpdateProductRequest, actorID uint64) (ProductResponse, error) {
	normalizeProductRequest(&request)
	product, err := s.repo.UpdateProduct(ctx, id, request, actorID)
	if err != nil {
		return ProductResponse{}, err
	}
	return ProductResponse{Product: product}, nil
}

func (s Service) DeleteProduct(ctx context.Context, id uint64) error {
	return s.repo.DeleteProduct(ctx, id)
}

func (s Service) ListPublicCategories(ctx context.Context) ([]Category, error) {
	return s.repo.ListActiveCategories(ctx)
}

func (s Service) ListAdminCategories(ctx context.Context) ([]Category, error) {
	return s.repo.ListAllCategories(ctx)
}

func (s Service) CreateCategory(ctx context.Context, request CreateCategoryRequest) (CategoryResponse, error) {
	normalizeCategoryRequest(&request)
	category, err := s.repo.CreateCategory(ctx, request)
	if err != nil {
		return CategoryResponse{}, err
	}
	return CategoryResponse{Category: category}, nil
}

func (s Service) UpdateCategory(ctx context.Context, id uint64, request UpdateCategoryRequest) (CategoryResponse, error) {
	normalizeCategoryRequest(&request)
	category, err := s.repo.UpdateCategory(ctx, id, request)
	if err != nil {
		return CategoryResponse{}, err
	}
	return CategoryResponse{Category: category}, nil
}

func (s Service) DeleteCategory(ctx context.Context, id uint64) error {
	return s.repo.DeleteCategory(ctx, id)
}

func normalizeProductRequest(request *CreateProductRequest) {
	request.Name = strings.TrimSpace(request.Name)
	request.Slug = strings.TrimSpace(request.Slug)
	request.SKU = trimOptional(request.SKU)
	request.ShortDescription = trimOptional(request.ShortDescription)
	request.Description = trimOptional(request.Description)
	request.Material = trimOptional(request.Material)
	request.MainImageURL = trimOptional(request.MainImageURL)
	if request.Status == "" {
		request.Status = ProductStatusDraft
	}
}

func normalizeCategoryRequest(request *CreateCategoryRequest) {
	request.Name = strings.TrimSpace(request.Name)
	request.Slug = strings.TrimSpace(request.Slug)
	request.Description = trimOptional(request.Description)
	if request.Status == "" {
		request.Status = CategoryStatusActive
	}
}

func trimOptional(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
