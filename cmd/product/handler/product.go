package handler

import (
	"context"
	"log"
	"strings"

	"github.com/krissukoco/go-microservices-marketplace/cmd/product/model"
	"github.com/krissukoco/go-microservices-marketplace/internal/statuscode"
	productPb "github.com/krissukoco/go-microservices-marketplace/proto/product"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func variantsModelToPb(variants []*model.ProductVariantGroup) []*productPb.Variant {
	var variantsPb []*productPb.Variant
	for _, variant := range variants {
		items := []*productPb.VariantItem{}
		for _, item := range variant.Values {
			items = append(items, &productPb.VariantItem{
				Value: item.Value,
				Price: item.Price,
				Stock: item.Stock,
			})
		}
		variantsPb = append(variantsPb, &productPb.Variant{
			Name:  variant.Name,
			Items: items,
		})
	}
	return variantsPb
}

func productModelToPb(product *model.Product) *productPb.Product {
	categorySplit := strings.Split(product.Category, "->")
	return &productPb.Product{
		Id:          product.Id,
		StoreId:     product.StoreId,
		Title:       product.Title,
		Slug:        product.Slug,
		Description: product.Description,
		Tags:        strings.Split(product.Tags, ","),
		Price:       product.Price,
		Category: &productPb.Category{
			Name:      categorySplit[len(categorySplit)-1],
			Hierarchy: categorySplit,
		},
		Variants:      variantsModelToPb(product.Variants),
		AverageRating: float32(product.AverageRating), // TODO: Round to 2 decimal place
		TotalReview:   product.TotalReview,
		TotalSold:     product.TotalSold,
		Stock:         product.Stock,
		CreatedAt:     product.CreatedAt,
		UpdatedAt:     product.UpdatedAt,
	}
}

func (s *Server) GetById(ctx context.Context, in *productPb.ProductId) (*productPb.Product, error) {
	var product *model.Product
	err := product.GetById(s.Pg, in.Id)
	if err != nil {
		if err == model.ErrProductNotFound {
			msg := statuscode.StandardErrorMessage(statuscode.ResourceNotFound, err.Error())
			return nil, status.Error(codes.NotFound, msg)
		}
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	return productModelToPb(product), nil
}

func (s *Server) GetBySlug(ctx context.Context, in *productPb.GetBySlugRequest) (*productPb.Product, error) {
	var product *model.Product
	err := product.GetBySlug(s.Pg, in.Slug)
	if err != nil {
		if err == model.ErrProductNotFound {
			msg := statuscode.StandardErrorMessage(statuscode.ResourceNotFound, err.Error())
			return nil, status.Error(codes.NotFound, msg)
		}
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	return productModelToPb(product), nil
}

func (s *Server) GetByStoreId(ctx context.Context, in *productPb.StoreId) (*productPb.ManyProductResponse, error) {
	var productsRes []*productPb.Product
	products, err := model.GetAllProductsByStore(s.Pg, in.Id)
	if err != nil {
		log.Println("[handler.product.GetByStoreId] Failed to get products: ", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	for _, product := range products {
		productsRes = append(productsRes, productModelToPb(product))
	}

	return &productPb.ManyProductResponse{Products: productsRes}, nil
}

func (s *Server) GetByFilters(ctx context.Context, in *productPb.GetByFiltersRequest) (*productPb.ManyProductResponse, error) {
	var productsRes []*productPb.Product
	products, err := model.GetAllProductsBySearch(s.Pg, in.Search, in.Category, int(in.Page), int(in.Limit))
	if err != nil {
		log.Println("[handler.product.GetByFilters] Failed to get products: ", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	for _, product := range products {
		productsRes = append(productsRes, productModelToPb(product))
	}
	return &productPb.ManyProductResponse{
		Products: productsRes,
	}, nil

}

func (s *Server) Create(ctx context.Context, in *productPb.NewProduct) (*productPb.ProductResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("[handler.product.Create] Failed to get metadata")
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	auths := md.Get("authorization")
	if len(auths) == 0 {
		log.Println("[handler.product.Create] Failed to get authorization")
		msg := statuscode.StandardErrorMessage(statuscode.Unauthorized, "Authorization is required")
		return nil, status.Error(codes.Unauthenticated, msg)
	}
	authToken := auths[0]
	// Get userId from auth
	userId, err := s.AuthClient.GetUserIdByToken(ctx, authToken)
	if err != nil {
		log.Println("[handler.product.Create] Failed to get userId from auth")
		msg := statuscode.StandardErrorMessage(statuscode.Unauthorized, "Invalid authorization")
		return nil, status.Error(codes.Unauthenticated, msg)
	}

	// Find storeId by userId
	var store *model.Store
	err = store.FindByUserId(s.Pg, userId)
	if err != nil {
		if err == model.ErrStoreNotFound {
			msg := statuscode.StandardErrorMessage(statuscode.ResourceNotFound, err.Error())
			return nil, status.Error(codes.NotFound, msg)
		}
		log.Println("[handler.product.Create] Failed to find store by userId")
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	// Validations
	if len(in.Title) < 5 {
		msg := statuscode.StandardErrorMessage(statuscode.ProductTitleTooShort, "Title must be at least 5 characters")
		return nil, status.Error(codes.FailedPrecondition, msg)
	}
	if len(in.Title) > 100 {
		msg := statuscode.StandardErrorMessage(statuscode.ProductTitleTooLong, "Title must be maximum 100 characters")
		return nil, status.Error(codes.FailedPrecondition, msg)
	}
	// TODO: Check and validate category
	splitTags := strings.Split(in.Tags, ",")
	if len(splitTags) > 5 {
		msg := statuscode.StandardErrorMessage(statuscode.ProductTagsTooLong, "Tags must be maximum 5")
		return nil, status.Error(codes.FailedPrecondition, msg)
	}
	if len(in.Variants) == 0 && in.Price <= 0 {
		msg := statuscode.StandardErrorMessage(statuscode.ProductPriceInvalid, "Price must be greater than 0")
		return nil, status.Error(codes.FailedPrecondition, msg)
	}
	var lowestPrice int64 = 0
	for _, v := range in.Variants {
		for _, item := range v.Items {
			if item.Price <= 0 {
				msg := statuscode.StandardErrorMessage(statuscode.ProductPriceInvalid, "Price must be greater than 0")
				return nil, status.Error(codes.FailedPrecondition, msg)
			}
			if lowestPrice == 0 || item.Price < lowestPrice {
				lowestPrice = item.Price
			}
		}
	}
	price := lowestPrice
	if price == 0 {
		price = in.Price
	}

	// Create product
	product := &model.Product{
		Id:          model.NewProductId(),
		StoreId:     store.Id,
		Title:       in.Title,
		Description: in.Description,
		Category:    in.Category,
		Tags:        in.Tags,
		Price:       price,
		Stock:       in.Stock,
	}
	err = product.Save(s.Pg)
	if err != nil {
		log.Println("[handler.product.Create] Failed to save product")
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	// Create variants
	for _, v := range in.Variants {
		variantName := v.Name
		for _, item := range v.Items {
			variant := &model.ProductVariant{
				ProductId: product.Id,
				Name:      variantName,
				Value:     item.Value,
				Price:     item.Price,
				Stock:     item.Stock,
			}
			err = variant.Save(s.Pg)
			if err != nil {
				log.Println("[handler.product.Create] Failed to save variant")
				return nil, status.Error(codes.Internal, "Internal server error")
			}
		}
	}
	err = product.FillVariants(s.Pg)
	if err != nil {
		log.Println("[handler.product.Create] Failed to fill variants")
	}

	return &productPb.ProductResponse{Product: productModelToPb(product)}, nil
}

func (s *Server) Delete(ctx context.Context, in *productPb.ProductId) (*productPb.ProductId, error) {
	var product *model.Product
	err := product.Delete(s.Pg)
	if err != nil {
		// TODO: Handle deletion error
		log.Println("[handler.product.Delete] Failed to delete product: ", err)
	}
	return &productPb.ProductId{Id: in.Id}, nil
}

// TODO: Implement Update Product
