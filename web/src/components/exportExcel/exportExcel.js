import {findAttendance } from '@/api/checkins/attendance'
import {getAttendanceCheckInList, getAttendanceCheckInLogList} from '@/api/checkins/attendanceCheckIn'
import {getGroupList} from '@/api/checkins/group'
import {getParticipantListByAttendance} from '@/api/checkins/participant'
import {getAreaList} from '@/api/checkins/area'
import {getConditionList} from '@/api/checkins/condition'
import { saveAs } from 'file-saver'
import * as XLSX from 'xlsx'

export const exportAttendanceData = async (attendanceId, page, pageSize) => {
  try {
    const attendanceCheckInList = await getAttendanceCheckInList({ attendanceId, page, pageSize })
    const areaList = await getAreaList({ attendanceId, page, pageSize })
    const conditionList = await getConditionList({ attendanceId, page, pageSize })
    const groupList = await getGroupList({ attendanceId, page, pageSize })
    const participantListByAttendance = await getParticipantListByAttendance({ attendanceId, page, pageSize })
    const attendance = await findAttendance({ id: attendanceId })

    const wb = XLSX.utils.book_new()

    console.log('attendanceCheckInList', attendanceCheckInList)

    //Tạo mảng các object data
    const dataSheets = [
      { name: 'Chi tiết ', data: [attendance?.data] },
      { name: 'Thành viên ', data: participantListByAttendance?.data?.list || [] },
      { name: 'Nhóm ', data: groupList?.data?.list || [] },
      { name: 'Khu vực ', data: areaList?.data?.list || [] },
      { name: 'Điều kiện ', data: conditionList?.data?.list || [] },
      { name: 'Lịch sử', data: (await getAttendanceCheckInLogList({ attendanceId, page, pageSize }))?.data?.list || [] },
      { name: 'Nhật ký', data: attendanceCheckInList?.data?.list || [] },
    ]

    dataSheets.forEach(({ name, data }) => {
      if (Array.isArray(data) && data.length > 0) {
        const sheet = XLSX.utils.json_to_sheet(data)
        XLSX.utils.book_append_sheet(wb, sheet, name)
      }
    })

    const excelBuffer = XLSX.write(wb, { bookType: 'xlsx', type: 'array' })
    const excelFile = new Blob([excelBuffer], { type: 'application/octet-stream' })
    saveAs(excelFile, `attendance_data_${attendanceId}.xlsx`)
  } catch (error) {
    console.error('Error exporting attendance data:', error)
  }
}

