import axios from 'axios' // 引入axios
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/pinia/modules/user'
import router from '@/router/index'
import { ElLoading } from 'element-plus'

const service = axios.create({
  baseURL: import.meta.env.VITE_BASE_API,
  timeout: 99999
})
let activeAxios = 0
let timer
let loadingInstance
const showLoading = (option = {
  target: null,
}) => {
  const loadDom = document.getElementById('gva-base-load-dom')
  activeAxios++
  if (timer) {
    clearTimeout(timer)
  }
  timer = setTimeout(() => {
    if (activeAxios > 0) {
      if (!option.target) option.target = loadDom
      loadingInstance = ElLoading.service(option)
    }
  }, 400)
}

const closeLoading = () => {
  activeAxios--
  if (activeAxios <= 0) {
    clearTimeout(timer)
    loadingInstance && loadingInstance.close()
  }
}
// http request 拦截器
service.interceptors.request.use(
  config => {
    if (!config.donNotShowLoading) {
      showLoading(config.loadingOption)
    }
    const userStore = useUserStore()
    config.headers = {
      'Content-Type': 'application/json',
      'x-token': userStore.token,
      'x-user-id': userStore.userInfo.ID,
      ...config.headers
    }
    return config
  },
  error => {
    if (!error.config.donNotShowLoading) {
      closeLoading()
    }
    ElMessage({
      showClose: true,
      message: error,
      type: 'error'
    })
    return error
  }
)

// http response 拦截器
service.interceptors.response.use(
  response => {
    const userStore = useUserStore()
    if (!response.config.donNotShowLoading) {
      closeLoading()
    }
    if (response.headers['new-token']) {
      userStore.setToken(response.headers['new-token'])
    }
    if (response.data.code === 0 || response.headers.success === 'true') {
      if (response.headers.msg) {
        response.data.msg = decodeURI(response.headers.msg)
      }
      return response.data
    } else {
      ElMessage({
        showClose: true,
        message: response.data.msg || decodeURI(response.headers.msg),
        type: 'error'
      })
      return response.data.msg ? response.data : response
    }
  },
  error => {
    if (!error.config.donNotShowLoading) {
      closeLoading()
    }

    if (!error.response) {
      ElMessageBox.confirm(`
        <p>Phát hiện lỗi trong yêu cầu</p>
        <p>${error}</p>
        `, 'Lỗi yêu cầu', {
        dangerouslyUseHTMLString: true,
        distinguishCancelAndClose: true,
        confirmButtonText: 'Thử lại sau',
        cancelButtonText: 'Hủy'
      })
      return
    }

    switch (error.response.status) {
      case 500:
        ElMessageBox.confirm(`
        <p>Phát hiện lỗi giao diện ${error}</p>
        <p>Mã lỗi<span style="color:red"> 500 </span>：Lỗi này thường xảy ra khi giao diện chưa được đăng ký (hoặc chưa khởi động lại) hoặc đường dẫn yêu cầu (phương thức) không khớp với đường dẫn API (phương thức)--Nếu đây là mã tự động, hãy kiểm tra xem có khoảng trắng không</p>
        `, 'Lỗi giao diện', {
          dangerouslyUseHTMLString: true,
          distinguishCancelAndClose: true,
          confirmButtonText: 'Tôi đã hiểu',
          cancelButtonText: 'Hủy'
        })
          .then(() => {
            const userStore = useUserStore()
            userStore.ClearStorage()
            router.push({ name: 'Login', replace: true })
          })
        break
      case 404:
        ElMessageBox.confirm(`
          <p>Phát hiện lỗi giao diện ${error}</p>
          <p>Mã lỗi<span style="color:red"> 404 </span>：Lỗi này thường xảy ra khi giao diện chưa được đăng ký (hoặc chưa khởi động lại) hoặc đường dẫn yêu cầu (phương thức) không khớp với đường dẫn API (phương thức)--Nếu đây là mã tự động, hãy kiểm tra xem có khoảng trắng không</p>
          `, 'Lỗi giao diện', {
          dangerouslyUseHTMLString: true,
          distinguishCancelAndClose: true,
          confirmButtonText: 'Tôi đã hiểu',
          cancelButtonText: 'Hủy'
        })
        break
      case 401:
        ElMessageBox.confirm(`
          <p>Token không hợp lệ</p>
          <p>Mã lỗi:<span style="color:red"> 401 </span>Thông tin lỗi:${error}</p>
          `, 'Thông tin định danh', {
          dangerouslyUseHTMLString: true,
          distinguishCancelAndClose: true,
          confirmButtonText: 'Đăng nhập lại',
          cancelButtonText: 'Hủy'
        })
          .then(() => {
            const userStore = useUserStore()
            userStore.ClearStorage()
            router.push({ name: 'Login', replace: true })
          })
        break
    }

    return error
  }
)
export default service
