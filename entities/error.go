package entities

import "errors"

var (
	// ErrBootsServiceUnavailable boots service unavailable
	ErrBootsServiceUnavailable = errors.New("boots_service_unavailable")

	// ErrLoginMemberNotFound login member not found
	ErrLoginMemberNotFound = errors.New("member_not_found")

	// ErrLoginMemberAlreadyExisting login member already existing
	ErrLoginMemberAlreadyExisting = errors.New("You already have an account, please login")
	ErrLoginMemberAlreadyExistingTH = errors.New("มีสมาชิกอยู่ในระบบแล้ว กรุณา login")

	// ErrRegisterBoots register boots
	ErrRegisterBoots = errors.New("register_boots")

	// ErrInvalidCitizen invalid citizen
	ErrInvalidCitizen = errors.New("invalid_citizen")

	// ErrFlashSaleExceed boots flash sales exceed
	ErrFlashSaleExceed = errors.New("boots_flash_sales_exceed")

	// Error invalid otp
	ErrInvalidOTP = errors.New("Invalid OTP code")
	ErrInvalidOTPTH = errors.New("รหัส OTP ไม่ถูกต้อง")
)
