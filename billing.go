package controld

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type BillingProduct struct {
	PK          int     `json:"PK"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Priority    int     `json:"priority"`
	ProxyAccess IntBool `json:"proxy_access"`
}

// PricePoint describes a product's price, and is reused both as the
// price_point of a Payment and the price of a Product; ControlD only
// includes Status in the latter.
type PricePoint struct {
	PK            int     `json:"PK"`
	ProductID     int     `json:"product_id"`
	Type          string  `json:"type"`
	Duration      int     `json:"duration"`
	Price         float64 `json:"price"`
	EURPrice      float64 `json:"eur_price"`
	GBPPrice      float64 `json:"gbp_price"`
	CHFPrice      float64 `json:"chf_price"`
	CADPrice      float64 `json:"cad_price"`
	AUDPrice      float64 `json:"aud_price"`
	JPYPrice      float64 `json:"jpy_price"`
	StripeID      string  `json:"stripe_id"`
	EURStripeID   string  `json:"eur_stripe_id"`
	GBPStripeID   string  `json:"gbp_stripe_id"`
	CHFStripeID   string  `json:"chf_stripe_id"`
	CADStripeID   string  `json:"cad_stripe_id"`
	AUDStripeID   string  `json:"aud_stripe_id"`
	JPYStripeID   string  `json:"jpy_stripe_id"`
	AlreadyBilled int     `json:"already_billed"`
	Comment       string  `json:"comment"`
	Status        *int    `json:"status,omitempty"`
}

type Payment struct {
	PK             string         `json:"PK"`
	ID             string         `json:"id"`
	TxID           string         `json:"tx_id"`
	SubID          string         `json:"sub_id"`
	Fingerprint    string         `json:"fingerprint"`
	Date           Date           `json:"date"`
	TxStatus       int            `json:"tx_status"`
	TxRefunded     IntBool        `json:"tx_refunded"`
	User           string         `json:"user"`
	Org            *string        `json:"org"`
	Product        BillingProduct `json:"product"`
	PricePoint     PricePoint     `json:"price_point"`
	Amount         float64        `json:"amount"`
	Currency       string         `json:"currency"`
	CurrencyAmount float64        `json:"currency_amount"`
	Balance        float64        `json:"balance"`
	Method         string         `json:"method"`
	Ts             UnixTime       `json:"ts"`
	IVState        any            `json:"iv_state"`
	IVMeta         any            `json:"iv_meta"`
}

type ListPaymentsBody struct {
	Payments []Payment `json:"payments"`
}

type ListPaymentsResponse struct {
	Body ListPaymentsBody `json:"body"`
	Response
}

// ListPayments returns the billing history of all payments made on the account.
func (api *API) ListPayments(ctx context.Context) ([]Payment, error) {
	uri := buildURI("/billing/payments", nil)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r ListPaymentsResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Payments, nil
}

type Subscription struct {
	PK             string         `json:"PK"`
	ID             string         `json:"id"`
	SubID          string         `json:"sub_id"`
	User           string         `json:"user"`
	Org            *string        `json:"org"`
	Product        BillingProduct `json:"product"`
	Method         string         `json:"method"`
	Currency       string         `json:"currency"`
	Amount         float64        `json:"amount"`
	CurrencyAmount float64        `json:"currency_amount"`
	Started        Date           `json:"started"`
	Ended          *Date          `json:"ended"`
	CancelReason   *string        `json:"cancel_reason"`
	Price          float64        `json:"price"`
	Status         int            `json:"status"`
	State          string         `json:"state"`
	NextBill       UnixTime       `json:"next_bill"`
	NextRebillDate *Date          `json:"next_rebill_date,omitempty"`
}

type ListSubscriptionsBody struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

type ListSubscriptionsResponse struct {
	Body ListSubscriptionsBody `json:"body"`
	Response
}

// ListSubscriptions returns all active and canceled subscriptions associated with the account.
func (api *API) ListSubscriptions(ctx context.Context) ([]Subscription, error) {
	uri := buildURI("/billing/subscriptions", nil)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r ListSubscriptionsResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Subscriptions, nil
}

// ProductSubscription is the subscription nested inside a Product; unlike
// Subscription.Product (a full BillingProduct object), its Product field is
// just the product's PK.
type ProductSubscription struct {
	PK             string   `json:"PK"`
	ID             string   `json:"id"`
	SubID          string   `json:"sub_id"`
	User           string   `json:"user"`
	Org            *string  `json:"org"`
	Product        int      `json:"product"`
	Method         string   `json:"method"`
	Currency       *string  `json:"currency"`
	Amount         float64  `json:"amount"`
	CurrencyAmount *float64 `json:"currency_amount"`
	Started        Date     `json:"started"`
	Ended          *Date    `json:"ended"`
	CancelReason   *string  `json:"cancel_reason"`
	Price          float64  `json:"price"`
	Status         int      `json:"status"`
	State          string   `json:"state"`
	NextBill       UnixTime `json:"next_bill"`
}

type Product struct {
	PK           int                 `json:"PK"`
	Name         string              `json:"name"`
	Type         string              `json:"type"`
	ProxyAccess  IntBool             `json:"proxy_access"`
	Expiry       Date                `json:"expiry"`
	Subscription ProductSubscription `json:"subscription"`
	Price        PricePoint          `json:"price"`
}

type ListProductsBody struct {
	Products []Product `json:"products"`
}

type ListProductsResponse struct {
	Body ListProductsBody `json:"body"`
	Response
}

// ListActiveProducts returns all products currently activated on the account.
func (api *API) ListActiveProducts(ctx context.Context) ([]Product, error) {
	uri := buildURI("/billing/products", nil)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r ListProductsResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Products, nil
}
