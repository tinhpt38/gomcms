import { useUserStore } from '@/pinia/modules/user'
import { useRouterStore } from '@/pinia/modules/router'
import getPageTitle from '@/utils/page'
import router from '@/router'
import Nprogress from 'nprogress'
import 'nprogress/nprogress.css'
Nprogress.configure({ showSpinner: false, ease: 'ease', speed: 500 })

const whiteList = ['Login', 'Init', 'Checkin']

const getRouter = async(userStore) => {
  const routerStore = useRouterStore()
  await routerStore.SetAsyncRouter()
  await userStore.GetUserInfo()
  const asyncRouters = routerStore.asyncRouters
  asyncRouters.forEach(asyncRouter => {
    router.addRoute(asyncRouter)
  })
}

async function handleKeepAlive(to) {
  if (to.matched.some(item => item.meta.keepAlive)) {
    if (to.matched && to.matched.length > 2) {
      for (let i = 1; i < to.matched.length; i++) {
        const element = to.matched[i - 1]
        if (element.name === 'layout') {
          to.matched.splice(i, 1)
          await handleKeepAlive(to)
        }
        // Nếu chưa hoàn thành tải theo nhu cầu thì đợi tải
        if (typeof element.components.default === 'function') {
          await element.components.default()
          await handleKeepAlive(to)
        }
      }
    }
  }
}

router.beforeEach(async(to, from) => {
  const routerStore = useRouterStore()
  Nprogress.start()
  const userStore = useUserStore()
  to.meta.matched = [...to.matched]
  handleKeepAlive(to)
  const token = userStore.token
  // Trường hợp nằm trong danh sách trắng
  document.title = getPageTitle(to.meta.title, to)
  if(to.meta.client) {
    return true
  }
  if (whiteList.indexOf(to.name) > -1) {
    if (token) {
      if (!routerStore.asyncRouterFlag && whiteList.indexOf(from.name) < 0) {
        await getRouter(userStore)
      }
      // token có thể giải mã nhưng không tồn tại id người dùng hoặc id vai trò sẽ dẫn đến việc gọi vô hạn
      if (userStore.userInfo?.authority?.defaultRouter != null) {
        if (router.hasRoute(userStore.userInfo.authority.defaultRouter)) {
          return { name: userStore.userInfo.authority.defaultRouter }
        } else {
          return { path: '/layout/404' }
        }
      } else {
        // Đăng xuất người dùng một cách bắt buộc
        userStore.ClearStorage()
        return {
          name: 'Login',
          query: {
            redirect: document.location.hash
          }
        }
      }
    } else {
      return true
    }
  } else {
    // 不在白名单中并且已经登录的时候
    if (token) {
      console.log(sessionStorage.getItem("needCloseAll"))
      if(sessionStorage.getItem("needToHome") === 'true') {
        sessionStorage.removeItem("needToHome")
        return { path: '/'}
      }
      // 添加flag防止多次获取动态路由和栈溢出
      if (!routerStore.asyncRouterFlag && whiteList.indexOf(from.name) < 0) {
        await getRouter(userStore)
        if (userStore.token) {
          if (router.hasRoute(userStore.userInfo.authority.defaultRouter)) {
            return { ...to, replace: true }
          } else {
            return { path: '/layout/404' }
          }
        } else {
          return {
            name: 'Login',
            query: { redirect: to.href }
          }
        }
      } else {
        if (to.matched.length) {
          return true
        } else {
          return { path: '/layout/404' }
        }
      }
    }
    // 不在白名单中并且未登录的时候
    if (!token) {
      return {
        name: 'Login',
        query: {
          redirect: document.location.hash
        }
      }
    }
  }
})

const removeLoading = () => {
  const element = document.getElementById('gva-loading-box');
  if (element) {
    element.remove();
  }
}

router.afterEach(() => {
  // 路由加载完成后关闭进度条
  document.getElementsByClassName('main-cont main-right')[0]?.scrollTo(0, 0)
  removeLoading()
  Nprogress.done()
})

router.onError(() => {
  // 路由发生错误后销毁进度条
  removeLoading()
  Nprogress.remove()
})

