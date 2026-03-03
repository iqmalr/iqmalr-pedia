package services

import (
	"errors"
	"log"
	"math"

	"github.com/iqmalr-pedia/go-product/internal/clients"
	"github.com/iqmalr-pedia/go-product/internal/dto/request"
	"github.com/iqmalr-pedia/go-product/internal/dto/response"
	"github.com/iqmalr-pedia/go-product/internal/models"
	"github.com/iqmalr-pedia/go-product/internal/repositories"
)

type ReviewService interface {
	CreateReview(productID, userID uint, req request.CreateReviewRequest) (*response.ReviewDetailResponse, error)
	GetReviews(productID uint, req request.ListReviewsRequest) (*response.ReviewListResponse, error)
	UpdateReview(productID, reviewID uint) (*response.ReviewDetailResponse, error)
	DeleteReview(productID, reviewID uint) error
	MarkHelpful(productID, reviewID, userID uint) (*response.ReviewHelpfulResponse, error)
}

type reviewService struct {
	reviewRepo  repositories.ReviewRepository
	productRepo repositories.ProductRepository
	authClient  clients.AuthClient
}

func NewReviewService(reviewRepo repositories.ReviewRepository, productRepo repositories.ProductRepository, authClient clients.AuthClient) ReviewService {
	return &reviewService{
		reviewRepo:  reviewRepo,
		productRepo: productRepo,
		authClient:  authClient,
	}
}

func (s *reviewService) CreateReview(productID, userID uint, req request.CreateReviewRequest) (*response.ReviewDetailResponse, error) {
	_, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	hasReviewed, err := s.reviewRepo.HasUserReviewedProduct(productID, userID)
	if err != nil {
		return nil, err
	}
	if hasReviewed {
		return nil, errors.New("you have already reviewed this product")
	}

	review := &models.ProductReview{
		ProductID:   productID,
		UserID:      userID,
		Rating:      req.Rating,
		Title:       req.Title,
		Comment:     req.Comment,
		OrderItemID: req.OrderItemID,
		IsApproved:  true,
	}

	if err := s.reviewRepo.Create(review); err != nil {
		return nil, err
	}

	return s.buildReviewResponse(review), nil
}

func (s *reviewService) GetReviews(productID uint, req request.ListReviewsRequest) (*response.ReviewListResponse, error) {
	_, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	reviews, total, err := s.reviewRepo.GetByProductID(productID, req.Page, req.Limit, req.Rating, req.Sort, req.Order)
	if err != nil {
		return nil, err
	}

	var data []response.ReviewDetailResponse
	for i := range reviews {
		data = append(data, *s.buildReviewResponse(&reviews[i]))
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &response.ReviewListResponse{
		Data: data,
		Pagination: response.PaginationResponse{
			Page:       req.Page,
			TotalPages: totalPages,
			Total:      int(total),
			Limit:      req.Limit,
		},
	}, nil
}

func (s *reviewService) UpdateReview(productID, reviewID uint) (*response.ReviewDetailResponse, error) {
	review, err := s.reviewRepo.GetByID(reviewID)
	if err != nil {
		return nil, err
	}

	if review.ProductID != productID {
		return nil, errors.New("review does not belong to this product")
	}

	review.IsApproved = !review.IsApproved

	if err := s.reviewRepo.Update(review); err != nil {
		return nil, err
	}

	return s.buildReviewResponse(review), nil
}

func (s *reviewService) DeleteReview(productID, reviewID uint) error {
	review, err := s.reviewRepo.GetByID(reviewID)
	if err != nil {
		return err
	}

	if review.ProductID != productID {
		return errors.New("review does not belong to this product")
	}

	return s.reviewRepo.Delete(reviewID)
}

func (s *reviewService) MarkHelpful(productID, reviewID, userID uint) (*response.ReviewHelpfulResponse, error) {
	review, err := s.reviewRepo.GetByID(reviewID)
	if err != nil {
		return nil, err
	}

	if review.ProductID != productID {
		return nil, errors.New("review does not belong to this product")
	}

	hasMarked, err := s.reviewRepo.HasUserMarkedHelpful(reviewID, userID)
	if err != nil {
		return nil, err
	}
	if hasMarked {
		return nil, errors.New("you have already marked this review as helpful")
	}

	helpful := &models.ReviewHelpful{
		ReviewID: reviewID,
		UserID:   userID,
	}

	if err := s.reviewRepo.CreateHelpful(helpful); err != nil {
		return nil, err
	}

	if err := s.reviewRepo.IncrementHelpfulCount(reviewID); err != nil {
		return nil, err
	}

	return &response.ReviewHelpfulResponse{
		Message:      "Review marked as helpful",
		HelpfulCount: review.HelpfulCount + 1,
	}, nil
}

func (s *reviewService) buildReviewResponse(review *models.ProductReview) *response.ReviewDetailResponse {
	user := response.ReviewUserResponse{ID: review.UserID}
	userResp, err := s.authClient.GetUserByID(review.UserID)
	if err != nil {
		log.Printf("Failed to fetch user %d info: %v", review.UserID, err)
		user.Name = "Unknown User"
	} else {
		user.Name = userResp.Name
	}

	return &response.ReviewDetailResponse{
		ID:                 review.ID,
		ProductID:          review.ProductID,
		UserID:             review.UserID,
		User:               user,
		OrderItemID:        review.OrderItemID,
		Rating:             review.Rating,
		Title:              review.Title,
		Comment:            review.Comment,
		IsVerifiedPurchase: review.IsVerifiedPurchase,
		IsApproved:         review.IsApproved,
		HelpfulCount:       review.HelpfulCount,
		CreatedAt:          review.CreatedAt,
		UpdatedAt:          review.UpdatedAt,
	}
}
