package repositories

import (
	"errors"

	"github.com/iqmalr-pedia/go-product/internal/models"
	"gorm.io/gorm"
)

type ReviewRepository interface {
	Create(review *models.ProductReview) error
	GetByID(id uint) (*models.ProductReview, error)
	GetByProductID(productID uint, page, limit int, rating *int, sort, order string) ([]models.ProductReview, int64, error)
	Update(review *models.ProductReview) error
	Delete(id uint) error
	HasUserReviewedProduct(productID, userID uint) (bool, error)
	IncrementHelpfulCount(reviewID uint) error
	HasUserMarkedHelpful(reviewID, userID uint) (bool, error)
	CreateHelpful(helpful *models.ReviewHelpful) error
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(review *models.ProductReview) error {
	return r.db.Create(review).Error
}

func (r *reviewRepository) GetByID(id uint) (*models.ProductReview, error) {
	var review models.ProductReview
	err := r.db.Where("id = ?", id).First(&review).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("review not found")
		}
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) GetByProductID(productID uint, page, limit int, rating *int, sort, order string) ([]models.ProductReview, int64, error) {
	var reviews []models.ProductReview
	var total int64

	query := r.db.Model(&models.ProductReview{}).Where("product_id = ?", productID)

	if rating != nil {
		query = query.Where("rating = ?", *rating)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if sort == "" {
		sort = "created_at"
	}
	if order == "" {
		order = "desc"
	}

	offset := (page - 1) * limit
	err := query.Order(sort + " " + order).Offset(offset).Limit(limit).Find(&reviews).Error
	if err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}

func (r *reviewRepository) Update(review *models.ProductReview) error {
	return r.db.Save(review).Error
}

func (r *reviewRepository) Delete(id uint) error {
	return r.db.Delete(&models.ProductReview{}, id).Error
}

func (r *reviewRepository) HasUserReviewedProduct(productID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.ProductReview{}).
		Where("product_id = ? AND user_id = ?", productID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *reviewRepository) IncrementHelpfulCount(reviewID uint) error {
	return r.db.Model(&models.ProductReview{}).
		Where("id = ?", reviewID).
		UpdateColumn("helpful_count", gorm.Expr("helpful_count + 1")).Error
}

func (r *reviewRepository) HasUserMarkedHelpful(reviewID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.ReviewHelpful{}).
		Where("review_id = ? AND user_id = ?", reviewID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *reviewRepository) CreateHelpful(helpful *models.ReviewHelpful) error {
	return r.db.Create(helpful).Error
}
