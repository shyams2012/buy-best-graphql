package interfaces

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.22

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/shyams2012/buy-best/Pagination"
	"github.com/shyams2012/buy-best/graph/generated"
	"github.com/shyams2012/buy-best/graph/model"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
)

// SetProduct is the resolver for the setProduct field.
func (r *mutationResolver) SetProduct(ctx context.Context, data model.ProductObj) (*model.Product, error) {
	if _, err := CheckAuth(ctx, []model.UserRole{"ADMIN", "EMPLOYEES"}); err != nil {
		return nil, err
	}

	var product model.Product
	if data.ID == nil {
		product.ID = uuid.NewString()
	} else {
		if tx := r.DB().First(&product, "id = ?", *data.ID); tx.Error != nil {
			return nil, fmt.Errorf("product not found, id='%s'", *data.ID)
		}
	}
	product.Name = data.Name
	product.Model = data.Model
	product.Description = data.Description
	product.Price = data.Price

	if tx := r.DB().Save(&product); tx.Error != nil {
		log.Print(tx.Error)
		return nil, fmt.Errorf("error saving product")
	}

	return &product, nil
}

// DeleteProduct is the resolver for the deleteProduct field.
func (r *mutationResolver) DeleteProduct(ctx context.Context, id string) (bool, error) {
	users := UserForContext(ctx)
	if users == nil {
		return false, errors.New("need authentication")
	}
	if users.Role != "ADMIN" {
		return false, fmt.Errorf("not authorized")
	}
	return model.DeleteObject(r.DB(), &model.Product{}, id)
}

// CreatePaymentIntent is the resolver for the createPaymentIntent field.
func (r *mutationResolver) CreatePaymentIntent(ctx context.Context, data model.StripePaymentData) (*model.PaymentIntent, error) {
	stripe.Key = os.Getenv("STRIPE_KEY")
	var product model.Product

	if tx := r.DB().First(&product, "id = ?", data.ProductID); tx.Error != nil {
		return nil, fmt.Errorf("product with id= %v not found", data.ProductID)

	}
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(product.Price)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
	}

	pi, err := paymentintent.New(params)
	log.Printf("pi.New: %v", pi.ClientSecret)

	if err != nil {
		log.Printf("pi.New : %v", err)
		return nil, fmt.Errorf("error creating payment intent")
	}

	paymentIntent := model.PaymentIntent{
		ClientSecret: pi.ClientSecret,
	}

	return &paymentIntent, nil
}

// GetProducts is the resolver for the getProducts field.
func (r *queryResolver) GetProducts(ctx context.Context, filter *model.ProductFilter, pagination *model.Pagination) (*model.ProductList, error) {
	var products []*model.Product
	// if filter == nil {
	// 	if tx := r.DB().Find(&product); tx.Error != nil {
	// 		log.Print(tx.Error)
	// 		return nil, tx.Error
	// 	}
	// } else {
	// 	if tx := r.DB().Where("price between ? and ?", *filter.Min, *filter.Max).Find(&product); tx.Error != nil {
	// 		log.Print(tx.Error)
	// 		return nil, tx.Error
	// 	}
	// }
	// return product, nil

	if pagination != nil {

		tx := r.DBWithFilter(filter).Scopes(Pagination.Paginate(&pagination.Page, &pagination.Limit)).Find(&products)

		if tx.Error != nil {
			log.Print(tx.Error)
			return nil, tx.Error
		}
		var totalRows int64
		tx.Count(&totalRows)
		count := int(totalRows)
		totalpage := &model.ProductPageInfo{TotalPages: count}
		productlists := &model.ProductList{PageInfo: totalpage, Product: products}
		return productlists, nil
	}
	tx := r.DB().Find(&products)
	if tx.Error != nil {
		log.Print(tx.Error)
		return nil, tx.Error
	}

	var totalRows int64
	tx.Count(&totalRows)
	count := int(totalRows)
	totalpage := &model.ProductPageInfo{TotalPages: count}

	productlist := &model.ProductList{PageInfo: totalpage, Product: products}
	return productlist, nil
}

// CompareProducts is the resolver for the compareProducts field.
func (r *queryResolver) CompareProducts(ctx context.Context, ids []*string) ([]*model.Product, error) {
	var products []*model.Product
	var products1 *model.Product
	var products2 *model.Product

	if tx := r.DB().Table("products").Select("name").Where("id = ?", *ids[0]).Scan(&products1); tx.Error != nil {
		return nil, fmt.Errorf("not found")

	}
	if tx := r.DB().Table("products").Select("name").Where("id = ?", *ids[1]).Scan(&products2); tx.Error != nil {
		return nil, fmt.Errorf("not found")

	}

	products = append(products, products1, products2)

	return products, nil
}

// ProductID is the resolver for the productId field.
func (r *transactionResolver) ProductID(ctx context.Context, obj *model.Transaction) (string, error) {
	panic(fmt.Errorf("not implemented: ProductID - productId"))
}

// Transaction returns generated.TransactionResolver implementation.
func (r *Resolver) Transaction() generated.TransactionResolver { return &transactionResolver{r} }

type transactionResolver struct{ *Resolver }
