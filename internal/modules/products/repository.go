package products

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrCategoryNotFound = errors.New("category not found")
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return Repository{db: db}
}

func (r Repository) ListActiveCategories(ctx context.Context) ([]Category, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT id, name, slug, description, status, created_at, updated_at
FROM product_categories
WHERE status = ?
ORDER BY name ASC`, CategoryStatusActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanCategories(rows)
}

func (r Repository) ListAllCategories(ctx context.Context) ([]Category, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT id, name, slug, description, status, created_at, updated_at
FROM product_categories
ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanCategories(rows)
}

func (r Repository) CreateCategory(ctx context.Context, request CreateCategoryRequest) (Category, error) {
	status := request.Status
	if status == "" {
		status = CategoryStatusActive
	}

	result, err := r.db.ExecContext(ctx, `
INSERT INTO product_categories (name, slug, description, status)
VALUES (?, ?, ?, ?)`,
		request.Name,
		request.Slug,
		request.Description,
		status,
	)
	if err != nil {
		return Category{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Category{}, err
	}

	return r.FindCategoryByID(ctx, uint64(id))
}

func (r Repository) UpdateCategory(ctx context.Context, id uint64, request UpdateCategoryRequest) (Category, error) {
	status := request.Status
	if status == "" {
		status = CategoryStatusActive
	}

	result, err := r.db.ExecContext(ctx, `
UPDATE product_categories
SET name = ?, slug = ?, description = ?, status = ?
WHERE id = ?`,
		request.Name,
		request.Slug,
		request.Description,
		status,
		id,
	)
	if err != nil {
		return Category{}, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return Category{}, err
	}
	if affected == 0 {
		return Category{}, ErrCategoryNotFound
	}

	return r.FindCategoryByID(ctx, id)
}

func (r Repository) DeleteCategory(ctx context.Context, id uint64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM product_categories WHERE id = ?", id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrCategoryNotFound
	}

	return nil
}

func (r Repository) FindCategoryByID(ctx context.Context, id uint64) (Category, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, name, slug, description, status, created_at, updated_at
FROM product_categories
WHERE id = ?
LIMIT 1`, id)

	category, err := scanCategory(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Category{}, ErrCategoryNotFound
	}
	return category, err
}

func (r Repository) ListProducts(ctx context.Context, query ListProductsQuery) (ListProductsResponse, error) {
	where, args := buildProductWhere(query)
	limit := normalizeLimit(query.Limit)
	page := normalizePage(query.Page)
	offset := (page - 1) * limit

	countQuery := "SELECT COUNT(1) FROM products p LEFT JOIN product_categories c ON c.id = p.category_id " + where
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return ListProductsResponse{}, err
	}

	selectArgs := append([]any{}, args...)
	selectArgs = append(selectArgs, limit, offset)

	rows, err := r.db.QueryContext(ctx, `
SELECT
    p.id, p.category_id, p.sku, p.name, p.slug, p.short_description, p.description,
    p.price, p.stock, p.material, p.weight_gram, p.length_mm, p.width_mm, p.height_mm,
    p.main_image_url, p.status, p.created_by, p.updated_by, p.created_at, p.updated_at,
    c.id, c.name, c.slug, c.description, c.status, c.created_at, c.updated_at
FROM products p
LEFT JOIN product_categories c ON c.id = p.category_id
`+where+`
ORDER BY p.created_at DESC
LIMIT ? OFFSET ?`, selectArgs...)
	if err != nil {
		return ListProductsResponse{}, err
	}
	defer rows.Close()

	items, err := scanProducts(rows)
	if err != nil {
		return ListProductsResponse{}, err
	}

	return ListProductsResponse{
		Items: items,
		Page:  page,
		Limit: limit,
		Total: total,
	}, nil
}

func (r Repository) FindPublicProductBySlug(ctx context.Context, slug string) (Product, error) {
	return r.findProduct(ctx, "p.slug = ? AND p.status = ?", slug, ProductStatusActive)
}

func (r Repository) FindProductByID(ctx context.Context, id uint64) (Product, error) {
	return r.findProduct(ctx, "p.id = ?", id)
}

func (r Repository) CreateProduct(ctx context.Context, request CreateProductRequest, actorID uint64) (Product, error) {
	status := request.Status
	if status == "" {
		status = ProductStatusDraft
	}

	result, err := r.db.ExecContext(ctx, `
INSERT INTO products (
    category_id, sku, name, slug, short_description, description, price, stock,
    material, weight_gram, length_mm, width_mm, height_mm, main_image_url,
    status, created_by, updated_by
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		request.CategoryID,
		request.SKU,
		request.Name,
		request.Slug,
		request.ShortDescription,
		request.Description,
		request.Price,
		request.Stock,
		request.Material,
		request.WeightGram,
		request.LengthMM,
		request.WidthMM,
		request.HeightMM,
		request.MainImageURL,
		status,
		actorID,
		actorID,
	)
	if err != nil {
		return Product{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Product{}, err
	}

	return r.FindProductByID(ctx, uint64(id))
}

func (r Repository) UpdateProduct(ctx context.Context, id uint64, request UpdateProductRequest, actorID uint64) (Product, error) {
	status := request.Status
	if status == "" {
		status = ProductStatusDraft
	}

	result, err := r.db.ExecContext(ctx, `
UPDATE products
SET category_id = ?, sku = ?, name = ?, slug = ?, short_description = ?, description = ?,
    price = ?, stock = ?, material = ?, weight_gram = ?, length_mm = ?, width_mm = ?,
    height_mm = ?, main_image_url = ?, status = ?, updated_by = ?
WHERE id = ?`,
		request.CategoryID,
		request.SKU,
		request.Name,
		request.Slug,
		request.ShortDescription,
		request.Description,
		request.Price,
		request.Stock,
		request.Material,
		request.WeightGram,
		request.LengthMM,
		request.WidthMM,
		request.HeightMM,
		request.MainImageURL,
		status,
		actorID,
		id,
	)
	if err != nil {
		return Product{}, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return Product{}, err
	}
	if affected == 0 {
		return Product{}, ErrProductNotFound
	}

	return r.FindProductByID(ctx, id)
}

func (r Repository) DeleteProduct(ctx context.Context, id uint64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrProductNotFound
	}

	return nil
}

func (r Repository) findProduct(ctx context.Context, condition string, args ...any) (Product, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT
    p.id, p.category_id, p.sku, p.name, p.slug, p.short_description, p.description,
    p.price, p.stock, p.material, p.weight_gram, p.length_mm, p.width_mm, p.height_mm,
    p.main_image_url, p.status, p.created_by, p.updated_by, p.created_at, p.updated_at,
    c.id, c.name, c.slug, c.description, c.status, c.created_at, c.updated_at
FROM products p
LEFT JOIN product_categories c ON c.id = p.category_id
WHERE `+condition+`
LIMIT 1`, args...)

	product, err := scanProduct(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Product{}, ErrProductNotFound
	}
	return product, err
}

func buildProductWhere(query ListProductsQuery) (string, []any) {
	conditions := make([]string, 0)
	args := make([]any, 0)

	if query.Public {
		conditions = append(conditions, "p.status = ?")
		args = append(args, ProductStatusActive)
	} else if query.Status != "" {
		conditions = append(conditions, "p.status = ?")
		args = append(args, query.Status)
	}

	if strings.TrimSpace(query.Search) != "" {
		conditions = append(conditions, "(p.name LIKE ? OR p.slug LIKE ? OR p.sku LIKE ?)")
		keyword := "%" + strings.TrimSpace(query.Search) + "%"
		args = append(args, keyword, keyword, keyword)
	}

	if strings.TrimSpace(query.Category) != "" {
		conditions = append(conditions, "c.slug = ?")
		args = append(args, strings.TrimSpace(query.Category))
	}

	if len(conditions) == 0 {
		return "", args
	}

	return "WHERE " + strings.Join(conditions, " AND "), args
}

func scanCategories(rows *sql.Rows) ([]Category, error) {
	categories := make([]Category, 0)
	for rows.Next() {
		category, err := scanCategory(rows)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanCategory(row scanner) (Category, error) {
	var category Category
	var description sql.NullString

	if err := row.Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&description,
		&category.Status,
		&category.CreatedAt,
		&category.UpdatedAt,
	); err != nil {
		return Category{}, err
	}

	category.Description = nullStringPtr(description)
	return category, nil
}

func scanProducts(rows *sql.Rows) ([]Product, error) {
	products := make([]Product, 0)
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func scanProduct(row scanner) (Product, error) {
	var product Product
	var categoryID, weight, length, width, height, createdBy, updatedBy sql.NullInt64
	var sku, shortDescription, description, material, mainImageURL sql.NullString
	var category Category
	var catID sql.NullInt64
	var catName, catSlug, catDescription, catStatus sql.NullString
	var catCreatedAt, catUpdatedAt sql.NullTime

	if err := row.Scan(
		&product.ID,
		&categoryID,
		&sku,
		&product.Name,
		&product.Slug,
		&shortDescription,
		&description,
		&product.Price,
		&product.Stock,
		&material,
		&weight,
		&length,
		&width,
		&height,
		&mainImageURL,
		&product.Status,
		&createdBy,
		&updatedBy,
		&product.CreatedAt,
		&product.UpdatedAt,
		&catID,
		&catName,
		&catSlug,
		&catDescription,
		&catStatus,
		&catCreatedAt,
		&catUpdatedAt,
	); err != nil {
		return Product{}, err
	}

	product.CategoryID = nullUint64Ptr(categoryID)
	product.SKU = nullStringPtr(sku)
	product.ShortDescription = nullStringPtr(shortDescription)
	product.Description = nullStringPtr(description)
	product.Material = nullStringPtr(material)
	product.WeightGram = nullUint64Ptr(weight)
	product.LengthMM = nullUint64Ptr(length)
	product.WidthMM = nullUint64Ptr(width)
	product.HeightMM = nullUint64Ptr(height)
	product.MainImageURL = nullStringPtr(mainImageURL)
	product.CreatedBy = nullUint64Ptr(createdBy)
	product.UpdatedBy = nullUint64Ptr(updatedBy)

	if catID.Valid {
		category.ID = uint64(catID.Int64)
		category.Name = catName.String
		category.Slug = catSlug.String
		category.Description = nullStringPtr(catDescription)
		category.Status = catStatus.String
		category.CreatedAt = catCreatedAt.Time
		category.UpdatedAt = catUpdatedAt.Time
		product.Category = &category
	}

	return product, nil
}

func nullStringPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func nullUint64Ptr(value sql.NullInt64) *uint64 {
	if !value.Valid {
		return nil
	}
	converted := uint64(value.Int64)
	return &converted
}

func normalizePage(page int) int {
	if page < 1 {
		return 1
	}
	return page
}

func normalizeLimit(limit int) int {
	if limit < 1 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func ValidateProductStatus(status string) bool {
	return status == "" || status == ProductStatusDraft || status == ProductStatusActive || status == ProductStatusInactive
}

func ValidateCategoryStatus(status string) bool {
	return status == "" || status == CategoryStatusActive || status == CategoryStatusInactive
}

func DuplicateMessage(err error) string {
	text := strings.ToLower(err.Error())
	switch {
	case strings.Contains(text, "products_slug_unique"):
		return "Product slug already used"
	case strings.Contains(text, "products_sku_unique"):
		return "Product SKU already used"
	case strings.Contains(text, "product_categories_slug_unique"):
		return "Category slug already used"
	default:
		return "Duplicate data"
	}
}

func ForeignKeyMessage(err error) string {
	return fmt.Sprintf("Related data is invalid: %v", err)
}
