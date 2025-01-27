package model

// Database
type Coupon struct {
	Id          int64  `gorm:"primary_key;auto_increment"`
	Username    string `gorm:"type:varchar(20); not null"`
	CouponName  string `gorm:"type:varchar(60); not null"`
	Amount      int64  // the maximum of coupon
	Left        int64  // left coupon
	Stock       int64  // stock of coupon
	Description string `gorm:"type:varchar(60)"` // description
}

type ReqCoupon struct {
	Name        string
	Amount      int64
	Description string
	Stock       int64
}

type ResCoupon struct {
	Name        string `json:"name"`
	Stock       int64  `json:"stock"`
	Description string `json:"description"`
}

// Seller request the coupon return
type SellerResCoupon struct {
	ResCoupon
	Amount int64 `json:"amount"`
	Left   int64 `json:"left"`
}

// client request
type CustomerResCoupon struct {
	ResCoupon
}

// for seller to request the coupons
func ParseSellerResCoupons(coupons []Coupon) []SellerResCoupon {
	var sellerCoupons []SellerResCoupon
	for _, coupon := range coupons {
		sellerCoupons = append(sellerCoupons,
			SellerResCoupon{ResCoupon{coupon.CouponName, coupon.Stock, coupon.Description},
				coupon.Amount, coupon.Left})
	}
	return sellerCoupons
}

// for customer to request coupons
func ParseCustomerResCoupons(coupons []Coupon) []CustomerResCoupon {
	var customerCoupons []CustomerResCoupon
	for _, coupon := range coupons {
		customerCoupons = append(customerCoupons,
			CustomerResCoupon{ResCoupon{coupon.CouponName, coupon.Stock, coupon.Description}})
	}
	return customerCoupons
}
