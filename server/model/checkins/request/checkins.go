package request

type CheckinsReq struct {
	Email     string   `json:"email" form:"email" gorm:"column:email;comment:Email" binding:"required"`
	Code      string   `json:"code" form:"code" gorm:"column:code;comment:Mã" binding:"required"`
	VisitorId string   `json:"visitorId" form:"visitorId" gorm:"column:visitor_id;comment:ID khách" binding:"required"`
	Lat       *float64 `json:"lat" form:"lat" gorm:"column:lat;comment:Vĩ độ"`
	Lng       *float64 `json:"lng" form:"lng" gorm:"column:lng;comment:Kinh độ"`
	Accuracy  *float64 `json:"accuracy" form:"accuracy" gorm:"column:accuracy;comment:Độ chính xác"`
	FullName  *string  `json:"fullName" form:"fullName"`
}
