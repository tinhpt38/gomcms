package checkins

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
	"github.com/xuri/excelize/v2"
)

func formatExportTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func formatExportFloat(f *float64) string {
	if f == nil {
		return ""
	}
	return fmt.Sprintf("%v", *f)
}

func sanitizeFileName(name string) string {
	replacer := strings.NewReplacer("/", "-", "\\", "-", ":", "-", "*", "-", "?", "-", "\"", "-", "<", "-", ">", "-", "|", "-")
	return strings.TrimSpace(replacer.Replace(name))
}

func writeExcelRow(f *excelize.File, sheet string, row int, values []interface{}) error {
	for col, value := range values {
		cell, err := excelize.CoordinatesToCellName(col+1, row)
		if err != nil {
			return err
		}
		if err = f.SetCellValue(sheet, cell, value); err != nil {
			return err
		}
	}
	return nil
}

func (attendanceCheckInService *AttendanceCheckInService) ExportAttendanceCheckInExcel(info checkinsReq.AttendanceCheckInSearch) (*bytes.Buffer, string, error) {
	if info.AttendanceId == nil {
		return nil, "", fmt.Errorf("attendanceId là bắt buộc")
	}

	info.Page = 1
	info.PageSize = 0
	list, _, err := attendanceCheckInService.GetAttendanceCheckInInfoList(info)
	if err != nil {
		return nil, "", err
	}

	var attendance checkins.Attendance
	if err = global.GVA_DB.Select("title").First(&attendance, *info.AttendanceId).Error; err != nil {
		return nil, "", err
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Sheet1"
	headers := []interface{}{
		"Ngày giờ", "Thành viên", "Email", "Khu vực", "Nhóm", "IP",
		"Vĩ độ", "Kinh độ", "Độ chính xác", "Client ID", "Số may mắn", "Lượt", "Agent",
	}
	if err = writeExcelRow(f, sheet, 1, headers); err != nil {
		return nil, "", err
	}

	for i, item := range list {
		participantName := ""
		participantEmail := ""
		if item.Participant != nil {
			participantEmail = item.Participant.Email
			if item.Participant.FullName != nil {
				participantName = strings.ReplaceAll(*item.Participant.FullName, "undefined", "")
			}
			if participantName == "" {
				participantName = participantEmail
			}
		}

		areaName := ""
		if item.Area != nil {
			areaName = item.Area.Name
		}

		groupName := ""
		if item.Group != nil {
			groupName = item.Group.Name
		}

		rowValues := []interface{}{
			formatExportTime(item.CheckinDate),
			participantName,
			participantEmail,
			areaName,
			groupName,
			item.IP,
			formatExportFloat(item.Lattidue),
			formatExportFloat(item.Longtidue),
			formatExportFloat(item.Accuracy),
			item.VisitorId,
			item.LuckyNumber,
			item.Counter,
			item.Agent,
		}
		if err = writeExcelRow(f, sheet, i+2, rowValues); err != nil {
			return nil, "", err
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", err
	}

	fileName := fmt.Sprintf("lich-su-diem-danh-%s.xlsx", sanitizeFileName(attendance.Title))
	return buf, fileName, nil
}

func (attendanceCheckInService *AttendanceCheckInService) ExportAttendanceCheckInLogExcel(info checkinsReq.AttendanceCheckInSearch) (*bytes.Buffer, string, error) {
	if info.AttendanceId == nil {
		return nil, "", fmt.Errorf("attendanceId là bắt buộc")
	}

	info.Page = 1
	info.PageSize = 0
	list, _, err := attendanceCheckInService.GetAttendanceCheckInLogInfoList(info)
	if err != nil {
		return nil, "", err
	}

	var attendance checkins.Attendance
	if err = global.GVA_DB.Select("title").First(&attendance, *info.AttendanceId).Error; err != nil {
		return nil, "", err
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Sheet1"
	headers := []interface{}{
		"Ngày giờ", "Họ tên", "Email", "Mã", "IP", "Vĩ độ", "Kinh độ",
		"Độ lệch", "Client ID", "Agent", "Thông báo",
	}
	if err = writeExcelRow(f, sheet, 1, headers); err != nil {
		return nil, "", err
	}

	for i, item := range list {
		fullName := ""
		if item.FullName != nil {
			fullName = *item.FullName
		}

		rowValues := []interface{}{
			formatExportTime(item.CreatedAt),
			fullName,
			item.Email,
			item.Code,
			item.Ip,
			formatExportFloat(item.Lat),
			formatExportFloat(item.Lng),
			formatExportFloat(item.Accuracy),
			item.VisitorId,
			item.Agent,
			item.MessageList,
		}
		if err = writeExcelRow(f, sheet, i+2, rowValues); err != nil {
			return nil, "", err
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", err
	}

	fileName := fmt.Sprintf("nhat-ky-diem-danh-%s.xlsx", sanitizeFileName(attendance.Title))
	return buf, fileName, nil
}
