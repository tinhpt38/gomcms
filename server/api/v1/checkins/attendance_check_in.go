package checkins

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AttendanceCheckInApi struct{}

func (attendanceCheckInApi *AttendanceCheckInApi) CreateAttendanceCheckIn(c *gin.Context) {
	var attendanceCheckIn checkins.AttendanceCheckIn
	err := c.ShouldBindJSON(&attendanceCheckIn)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	attendanceCheckIn.CreatedBy = utils.GetUserID(c)
	err = attendanceCheckInService.CreateAttendanceCheckIn(&attendanceCheckIn)
	if err != nil {
		global.GVA_LOG.Error("thất bại!", zap.Error(err))
		response.FailWithMessage("thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("thành công", c)
}

func (attendanceCheckInApi *AttendanceCheckInApi) DeleteAttendanceCheckIn(c *gin.Context) {
	ID := c.Query("ID")
	userID := utils.GetUserID(c)
	err := attendanceCheckInService.DeleteAttendanceCheckIn(ID, userID)
	if err != nil {
		global.GVA_LOG.Error("thất bại!", zap.Error(err))
		response.FailWithMessage("thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("thành công", c)
}

func (attendanceCheckInApi *AttendanceCheckInApi) DeleteAttendanceCheckInByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	userID := utils.GetUserID(c)
	err := attendanceCheckInService.DeleteAttendanceCheckInByIds(IDs, userID)
	if err != nil {
		global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("Thành công", c)
}

func (attendanceCheckInApi *AttendanceCheckInApi) UpdateAttendanceCheckIn(c *gin.Context) {
	var attendanceCheckIn checkins.AttendanceCheckIn
	err := c.ShouldBindJSON(&attendanceCheckIn)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	attendanceCheckIn.UpdatedBy = utils.GetUserID(c)
	err = attendanceCheckInService.UpdateAttendanceCheckIn(attendanceCheckIn)
	if err != nil {
		global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("Thành công", c)
}

func (attendanceCheckInApi *AttendanceCheckInApi) FindAttendanceCheckIn(c *gin.Context) {
	ID := c.Query("ID")
	reattendanceCheckIn, err := attendanceCheckInService.GetAttendanceCheckIn(ID)
	if err != nil {
		global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:"+err.Error(), c)
		return
	}
	response.OkWithData(reattendanceCheckIn, c)
}

func (attendanceCheckInApi *AttendanceCheckInApi) GetAttendanceCheckInList(c *gin.Context) {
	var pageInfo checkinsReq.AttendanceCheckInSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := attendanceCheckInService.GetAttendanceCheckInInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Thành công", c)
}

func (attendanceCheckInApi *AttendanceCheckInApi) GetAttendanceCheckInLogList(c *gin.Context) {
	var pageInfo checkinsReq.AttendanceCheckInSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := attendanceCheckInService.GetAttendanceCheckInLogInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Thành công", c)
}

func (attendanceCheckInApi *AttendanceCheckInApi) GetAttendanceCheckInPublic(c *gin.Context) {
	response.OkWithDetailed(gin.H{
		"info": "Thông tin API checkins của thành viên không cần xác thực",
	}, "Thành công", c)
}

func (attendanceCheckInApi *AttendanceCheckInApi) CheckinAttendance(c *gin.Context) {
	var checkinReqEncode checkinsReq.CheckinReqEncode
	err := c.ShouldBindJSON(&checkinReqEncode)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	endodePrefix := "E;>YIws8_DdsSMG£sL£@lq8E<(O?Sc5"
	if !strings.HasPrefix(checkinReqEncode.Data, endodePrefix) {
		response.FailWithMessage("Dữ liệu gửi lên của bạn không toàn vẹn", c)
		return
	}

	// Xử lý định dạng mã hóa mới: keyRandom + "_" + randomToken + "_" + encodedData
	parts := strings.SplitN(checkinReqEncode.Data, "_", 3)
	if len(parts) < 3 {
		// Xử lý cách cũ trong trường hợp chỉ có 2 phần (để tương thích ngược)
		parts = strings.SplitN(checkinReqEncode.Data, "_", 2)
		if len(parts) < 2 {
			response.FailWithMessage("Định dạng dữ liệu không hợp lệ", c)
			return
		}
		checkinReqEncode.Data = parts[1]
		global.GVA_LOG.Info("Sử dụng định dạng mã hóa cũ")
	} else {
		// Lấy phần data được mã hóa (phần cuối cùng trong chuỗi đã tách)
		checkinReqEncode.Data = parts[2]
		global.GVA_LOG.Info("Sử dụng định dạng mã hóa mới: version 2.0")
	}

	var checkinReq checkinsReq.CheckinsReq
	decodedData, err := decodeUnicodeBase64(checkinReqEncode.Data)
	if err != nil {
		global.GVA_LOG.Error("Lỗi giải mã dữ liệu: ", zap.Error(err))
		response.FailWithMessage("Dữ liệu gửi lên của bạn không toàn vẹn", c)
		return
	}

	// Kiểm tra JSON hợp lệ và ghi log dữ liệu đã giải mã để debug
	global.GVA_LOG.Debug("Dữ liệu đã giải mã: " + decodedData)

	var rawData interface{}
	if errJSON := json.Unmarshal([]byte(decodedData), &rawData); errJSON != nil {
		global.GVA_LOG.Error("Dữ liệu JSON không hợp lệ: ", zap.Error(errJSON))
		response.FailWithMessage("Dữ liệu không hợp lệ", c)
		return
	}

	// Kiểm tra phiên bản mã hóa
	dataMap, ok := rawData.(map[string]interface{})
	if ok {
		if version, exists := dataMap["_v"]; exists {
			global.GVA_LOG.Info("Phiên bản mã hóa: " + version.(string))

			// Xác minh phiên bản mã hóa
			if version.(string) == "2.0" {
				global.GVA_LOG.Info("Đang sử dụng phiên bản mã hóa 2.0")
			}
		}

		// Kiểm tra timestamp nếu có
		if timestamp, exists := dataMap["timestamp"]; exists {
			global.GVA_LOG.Info("Timestamp: " + timestamp.(string))
		}

		// Kiểm tra và xác minh checksum nếu có
		if checksum, exists := dataMap["_checksum"]; exists {
			// Copy dữ liệu và loại bỏ checksum để tính toán lại
			checksumData := make(map[string]interface{})
			for k, v := range dataMap {
				if k != "_checksum" {
					checksumData[k] = v
				}
			}

			// Log giá trị checksum để debug
			global.GVA_LOG.Info("Checksum từ client: " + checksum.(string))

			// Ở đây có thể thêm logic để tính toán checksum và so sánh nếu cần
			// Ví dụ:
			// calculatedChecksum := generateChecksum(checksumData)
			// if calculatedChecksum != checksum.(string) {
			// 	global.GVA_LOG.Error("Checksum không khớp")
			// 	response.FailWithMessage("Dữ liệu đã bị thay đổi, không toàn vẹn", c)
			// 	return
			// }
		}
	}

	// Giải mã thành đối tượng CheckinsReq
	err = json.Unmarshal([]byte(decodedData), &checkinReq)
	if err != nil {
		global.GVA_LOG.Error("Lỗi chuyển đổi JSON sang đối tượng: ", zap.Error(err))
		response.FailWithMessage("Dữ liệu gửi lên của bạn không toàn vẹn", c)
		return
	}

	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	prefix := "dlu_activities_20422_5BS:W`A8nF<J6Y{V4Nv.r!Je_"
	if !strings.HasPrefix(checkinReq.VisitorId, prefix) {
		response.FailWithMessage("Từ chối điểm danh. Bạn đang điểm danh từ một thiết bị không được phép", c)
		return
	}
	checkinReq.VisitorId = strings.TrimPrefix(checkinReq.VisitorId, prefix)

	result, err := attendanceCheckInService.CheckinAttendance(checkinReq, ip, userAgent)
	if err != nil {
		global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:"+err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// Hàm giải mã base64 và Unicode để khớp với encodeUnicode ở frontend
func decodeUnicodeBase64(s string) (string, error) {
	// Giải mã base64
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	// Chuyển mảng byte thành chuỗi
	str := string(data)

	// Giải mã URL encoding
	result, err := url.QueryUnescape(str)
	if err != nil {
		// Nếu không phải là URL encoded, trả về chuỗi gốc
		// Điều này có thể xảy ra nếu dữ liệu không được mã hóa qua encodeURIComponent
		var decodedStr string
		for i := 0; i < len(str); i++ {
			if str[i] == '%' && i+2 < len(str) {
				// Thử chuyển đổi mã hex sang byte
				b, err := hex2byte(str[i+1 : i+3])
				if err == nil {
					decodedStr += string(b)
					i += 2
				} else {
					decodedStr += string(str[i])
				}
			} else {
				decodedStr += string(str[i])
			}
		}
		return decodedStr, nil
	}

	return result, nil
}

// Hàm chuyển đổi 2 ký tự hex thành byte
func hex2byte(hex string) (byte, error) {
	// Chuyển đổi từ hex sang decimal
	var b byte
	for i := 0; i < len(hex); i++ {
		b = b*16 + hexDigit(hex[i])
	}
	return b, nil
}

// Chuyển đổi một ký tự hex thành giá trị decimal
func hexDigit(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}
