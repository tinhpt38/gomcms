/*
 * gin-vue-admin web框架组
 *
 * */
// 加载网站配置文件夹
import { register } from './global'

export default {
  install: (app) => {
    register(app)
    console.log(`
   Chào mừng đến với Gin-Vue-Admin
    Phiên bản hiện tại: v2.7.2
 Địa chỉ tài liệu tự động hóa mặc định: http://127.0.0.1:${import.meta.env.VITE_SERVER_PORT}/swagger/index.html
 Địa chỉ chạy tệp giao diện người dùng mặc định: http://127.0.0.1:${import.meta.env.VITE_CLI_PORT}
    `)
  }
}
