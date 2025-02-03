package dbService

import (
	"SecKill/data"
	"SecKill/model"
	"fmt"
)

func GetAllCoupons() ([]model.Coupon, error) {
	// return all the coupons
	var coupons []model.Coupon
	result := data.Db.Find(&coupons)
	return coupons, result.Error
}

// Insert the data of coupon users have.
func UserHasCoupon(userName string, coupon model.Coupon) error {
	return data.Db.Exec(fmt.Sprintf("INSERT IGNORE INTO coupons "+
		"(`username`,`coupon_name`,`amount`,`left`,`stock`,`description`) "+
		"values('%s', '%s', %d, %d, %f, '%s')",
		userName, coupon.CouponName, 1, 1, coupon.Stock, coupon.Description)).Error
}

// the stock of coupon in the database decrease by 1
func DecreaseOneCouponLeft(sellerName string, couponName string) error {
	return data.Db.Exec(fmt.Sprintf("UPDATE coupons c SET c.left=c.left-1 WHERE "+
		"c.username='%s' AND c.coupon_name='%s' AND c.left>0", sellerName, couponName)).Error
}
