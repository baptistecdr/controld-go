package controld

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestListPayments(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"payments": [
				  {
					"id": "70a569d4-67fe-4fe8-b198-2502e0575aaa",
					"tx_id": "in_test",
					"sub_id": "sub_test",
					"fingerprint": "fingerprint",
					"date": "2023-01-27",
					"tx_status": 1,
					"tx_refunded": 0,
					"user": "userPK",
					"org": null,
					"product": {
					  "proxy_access": 1,
					  "type": "standard",
					  "priority": 1,
					  "name": "Full Control",
					  "PK": 3
					},
					"price_point": {
					  "product_id": 3,
					  "price": 60,
					  "eur_price": 51,
					  "gbp_price": 45,
					  "chf_price": 47,
					  "cad_price": 83,
					  "aud_price": 83,
					  "jpy_price": 9500,
					  "stripe_id": "price_1",
					  "eur_stripe_id": "price_2",
					  "gbp_stripe_id": "price_3",
					  "chf_stripe_id": "price_4",
					  "cad_stripe_id": "price_5",
					  "aud_stripe_id": "price_6",
					  "jpy_stripe_id": "price_7",
					  "already_billed": 0,
					  "comment": "N/A",
					  "type": "standard",
					  "duration": 12,
					  "PK": 4
					},
					"amount": 40,
					"currency": "usd",
					"currency_amount": 40,
					"balance": 0,
					"method": "credit card",
					"ts": 1674844026,
					"iv_state": null,
					"iv_meta": null,
					"PK": "in_test"
				  }
				]
			  },
			  "success": true
			}
		`)
	}

	mux.HandleFunc("/billing/payments", handler)
	actual, err := client.ListPayments(context.Background())

	want := []Payment{
		{
			PK:          "in_test",
			ID:          "70a569d4-67fe-4fe8-b198-2502e0575aaa",
			TxID:        "in_test",
			SubID:       "sub_test",
			Fingerprint: "fingerprint",
			Date:        Date{time.Date(2023, 1, 27, 0, 0, 0, 0, time.UTC)},
			TxStatus:    1,
			TxRefunded:  IntBool(false),
			User:        "userPK",
			Product: BillingProduct{
				PK:          3,
				Name:        "Full Control",
				Type:        "standard",
				Priority:    1,
				ProxyAccess: IntBool(true),
			},
			PricePoint: PricePoint{
				PK:            4,
				ProductID:     3,
				Type:          "standard",
				Duration:      12,
				Price:         60,
				EURPrice:      51,
				GBPPrice:      45,
				CHFPrice:      47,
				CADPrice:      83,
				AUDPrice:      83,
				JPYPrice:      9500,
				StripeID:      "price_1",
				EURStripeID:   "price_2",
				GBPStripeID:   "price_3",
				CHFStripeID:   "price_4",
				CADStripeID:   "price_5",
				AUDStripeID:   "price_6",
				JPYStripeID:   "price_7",
				AlreadyBilled: 0,
				Comment:       "N/A",
			},
			Amount:         40,
			Currency:       "usd",
			CurrencyAmount: 40,
			Method:         "credit card",
			Ts:             UnixTime{time.Unix(1674844026, 0).UTC()},
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}
}

func TestListSubscriptions(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"subscriptions": [
				  {
					"id": "a26339c2-9f78-4406-8ecf-71d86d06cd7f",
					"sub_id": "sub_test",
					"user": "userPK",
					"org": null,
					"product": {
					  "proxy_access": 1,
					  "type": "standard",
					  "priority": 1,
					  "name": "Full Control",
					  "PK": 3
					},
					"method": "stripe",
					"currency": "usd",
					"amount": 40,
					"currency_amount": 40,
					"started": "2023-01-27",
					"ended": null,
					"cancel_reason": null,
					"price": 4,
					"status": 1,
					"state": "active",
					"next_bill": 1801074421,
					"PK": "sub_test",
					"next_rebill_date": "2027-01-29"
				  }
				]
			  },
			  "success": true
			}
		`)
	}

	mux.HandleFunc("/billing/subscriptions", handler)
	actual, err := client.ListSubscriptions(context.Background())

	nextRebillDate := Date{time.Date(2027, 1, 29, 0, 0, 0, 0, time.UTC)}
	want := []Subscription{
		{
			PK:    "sub_test",
			ID:    "a26339c2-9f78-4406-8ecf-71d86d06cd7f",
			SubID: "sub_test",
			User:  "userPK",
			Product: BillingProduct{
				PK:          3,
				Name:        "Full Control",
				Type:        "standard",
				Priority:    1,
				ProxyAccess: IntBool(true),
			},
			Method:         "stripe",
			Currency:       "usd",
			Amount:         40,
			CurrencyAmount: 40,
			Started:        Date{time.Date(2023, 1, 27, 0, 0, 0, 0, time.UTC)},
			Price:          4,
			Status:         1,
			State:          "active",
			NextBill:       UnixTime{time.Unix(1801074421, 0).UTC()},
			NextRebillDate: &nextRebillDate,
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}
}

func TestListActiveProducts(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"products": [
				  {
					"proxy_access": 1,
					"type": "standard",
					"name": "Full Control",
					"PK": 3,
					"expiry": "2027-01-29",
					"subscription": {
					  "id": "a26339c2-9f78-4406-8ecf-71d86d06cd7f",
					  "sub_id": "sub_test",
					  "user": "userPK",
					  "org": null,
					  "product": 3,
					  "method": "stripe",
					  "currency": null,
					  "amount": 40,
					  "currency_amount": null,
					  "started": "2023-01-27",
					  "ended": null,
					  "cancel_reason": null,
					  "price": 4,
					  "status": 1,
					  "state": "active",
					  "next_bill": 1801074421,
					  "PK": "sub_test"
					},
					"price": {
					  "product_id": 3,
					  "price": 60,
					  "eur_price": 51,
					  "gbp_price": 45,
					  "chf_price": 47,
					  "cad_price": 83,
					  "aud_price": 83,
					  "jpy_price": 9500,
					  "stripe_id": "price_1",
					  "eur_stripe_id": "price_2",
					  "gbp_stripe_id": "price_3",
					  "chf_stripe_id": "price_4",
					  "cad_stripe_id": "price_5",
					  "aud_stripe_id": "price_6",
					  "jpy_stripe_id": "price_7",
					  "already_billed": 0,
					  "comment": "N/A",
					  "type": "standard",
					  "duration": 12,
					  "PK": 4,
					  "status": 1
					}
				  }
				]
			  },
			  "success": true
			}
		`)
	}

	mux.HandleFunc("/billing/products", handler)
	actual, err := client.ListActiveProducts(context.Background())

	priceStatus := 1
	want := []Product{
		{
			PK:          3,
			Name:        "Full Control",
			Type:        "standard",
			ProxyAccess: IntBool(true),
			Expiry:      Date{time.Date(2027, 1, 29, 0, 0, 0, 0, time.UTC)},
			Subscription: ProductSubscription{
				PK:       "sub_test",
				ID:       "a26339c2-9f78-4406-8ecf-71d86d06cd7f",
				SubID:    "sub_test",
				User:     "userPK",
				Product:  3,
				Method:   "stripe",
				Amount:   40,
				Started:  Date{time.Date(2023, 1, 27, 0, 0, 0, 0, time.UTC)},
				Price:    4,
				Status:   1,
				State:    "active",
				NextBill: UnixTime{time.Unix(1801074421, 0).UTC()},
			},
			Price: PricePoint{
				PK:            4,
				ProductID:     3,
				Type:          "standard",
				Duration:      12,
				Price:         60,
				EURPrice:      51,
				GBPPrice:      45,
				CHFPrice:      47,
				CADPrice:      83,
				AUDPrice:      83,
				JPYPrice:      9500,
				StripeID:      "price_1",
				EURStripeID:   "price_2",
				GBPStripeID:   "price_3",
				CHFStripeID:   "price_4",
				CADStripeID:   "price_5",
				AUDStripeID:   "price_6",
				JPYStripeID:   "price_7",
				AlreadyBilled: 0,
				Comment:       "N/A",
				Status:        &priceStatus,
			},
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}
}
