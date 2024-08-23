import { ElNotification } from 'element-plus'
import { ElMessageBox } from 'element-plus'
import { ElMessage } from 'element-plus'

// Notification: hiển thị thông báo đơn hoặc 1 chuỗi thông báo theo quy trình
export const useNotification = () => {
  const notify = (type, message, duration = 5000, config = {}) => {
    const cfg = { ...config }

    cfg.duration = duration

    ElNotification({
      type: type,
      message: message,
      ...cfg,
    })
  }

  const showSuccessNotification = (message, duration, config) => {
    notify('success', message, duration, config)
  }

  const showWarningNotification = (message, duration, config) => {
    notify('warning', message, duration, config)
  }

  const showInfoNotification = (message, duration, config) => {
    notify('info', message, duration, config)
  }

  const showErrorNotification = (message, duration, config) => {
    notify('error', message, duration, config)
  }

  const showNoMessageNotification = (message, duration, config) => {
    notify('', message, duration, config)
  }

  const multipleNotify = (progress) => {
    if (Array.isArray(progress) && progress.length > 0) {
      progress.forEach((notification, index) => {
        setTimeout(() => {
          ElNotification({
            type: 'info',
            message: notification.message,
          })
        }, index * 1000) // Hiển thị thông báo mỗi 1 giây
      })
    }
  }

  const closeAll = (app) => {
    app.appContext.config.globalProperties.$notify.closeAll()
  }

  return {
    showSuccessNotification,
    showWarningNotification,
    showInfoNotification,
    showErrorNotification,
    showNoMessageNotification,
    closeAll,
    multipleNotify,
  }
}

export const useConfirmBox = () => {
  const confirm = (message, title, type, onOk, onCancel, config) => {
    const cfg = config || {
      confirmButtonText: 'Xác Nhận',
      cancelButtonText: 'Huỷ Bỏ',
      type: type,
    }

    ElMessageBox.confirm(message, title, {
      ...cfg,
    })
      .then(() => {
        if (onOk) onOk()
      })
      .catch(() => {
        if (onCancel) onCancel()
      })
  }

  const dialogError = (message, onOk, onCancel, config) => {
    confirm(message, 'Thông báo lỗi', 'error', onOk, onCancel, config)
  }

  const dialogConfirm = (message, onOk, onCancel, config) => {
    confirm(message, 'Xác nhận', 'success', onOk, onCancel, config)
  }

  return {
    confirm,
    dialogError,
    dialogConfirm,
  }
}

export const useMessage = () => {
  const notify = (message, type, duration = 5000, config = {}) => {
    // config is optional
    const cfg = { ... config }
    cfg.duration = duration

    ElMessage({
      message: message,
      type: type,
      ...cfg
    })
  }

  const showSuccessMessage = (message, duration, config) => {
    notify(message, 'success', duration, config)
  }

  const showErrorMessage = (message, duration, config) => {
    notify(message, 'error', duration, config)
  }
  const showInfoMessage = (message, duration, config) => {
    notify(message, 'info', duration, config)
  }

  const closeAll = (app) => {
    app.appContext.config.globalProperties.$notify.closeAll()
  }

  return { showSuccessMessage, showErrorMessage, showInfoMessage, closeAll }
}
