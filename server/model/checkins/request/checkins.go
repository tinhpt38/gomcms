package request

type CheckinsReq struct {
	Email string   `json:"email" form:"email" gorm:"column:email;comment:Email" binding:"required"`
	Code  string   `json:"code" form:"code" gorm:"column:code;comment:Mã" binding:"required"`
	Lat   *float64 `json:"lat" form:"lat" gorm:"column:lat;comment:Vĩ độ" binding:"required"`
	Lng   *float64 `json:"lng" form:"lng" gorm:"column:lng;comment:Kinh độ" binding:"required"`
}
