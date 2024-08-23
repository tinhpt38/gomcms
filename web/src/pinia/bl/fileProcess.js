import { computed, nextTick, ref } from 'vue'
import { defineStore } from 'pinia'

import { getPercentFileProcess as getPercentFileProcessAPI } from '@/api/config/cfg_file_process'
import { isJsonString } from '@/utils/validator'
import memoize from 'lodash/memoize'

const INTERVAL_TIMEOUT = 3000

export const useFileProcessStore = defineStore('fileProcess', () => {
  const fileProcess = ref({})

  const runningFileProcess = ref({})

  const displayFileProcess = computed(() => {
    if (!runningFileProcess.value) return []
    if (Object.keys(runningFileProcess.value) === 0) return []

    return Object.values(runningFileProcess.value)
  })

  const saveToLocalStorage = () => {
    localStorage.setItem('fileProcess', JSON.stringify(fileProcess.value))
  }

  const initFileProcess = () => {
    const local = localStorage.getItem('fileProcess')
    if (!local) return

    let fileProcessObject = {}

    if (isJsonString(local)) {
      fileProcessObject = JSON.parse(local)
    }

    fileProcess.value = fileProcessObject

    Object.values(fileProcessObject).forEach(fp => {
      addFileProcess(fp, false)
    })
  }
  const initFileProcessMemoized = memoize(initFileProcess)

  const addFileProcess = async(fp, save = true) => {
    const uuid = fp?.uuid || null

    if (!uuid) {
      return
    }

    getFileProcessStatus(uuid)

    const intervalID = setInterval(() => {
      getFileProcessStatus(uuid)
    }, INTERVAL_TIMEOUT)

    runningFileProcess.value = {
      ...runningFileProcess.value,
      [uuid]: {
        ...fp,
        intervalID
      }
    }

    fileProcess.value = {
      ...fileProcess.value,
      [uuid]: fp
    }

    await nextTick()

    if (save) {
      saveToLocalStorage()
    }
  }

  const removeFileProcess = async(fp) => {
    const uuid = fp.uuid

    const runningProcess = runningFileProcess.value[uuid] || null
    if (runningFileProcess.value) {
      clearInterval(runningProcess.intervalID)
    }

    runningFileProcess.value = Object.values(runningFileProcess.value).filter(a => a.uuid !== uuid).reduce((acc, cur) => {
      return {
        ...acc,
        [cur.uuid]: cur
      }
    }, {})

    fileProcess.value = Object.values(fileProcess.value).filter(a => a.uuid !== uuid).reduce((acc, cur) => {
      return {
        ...acc,
        [cur.uuid]: cur
      }
    }, {})

    await nextTick()

    saveToLocalStorage()
  }

  const getFileProcessStatus = async(uuid) => {
    try {
      const res = await getPercentFileProcessAPI({
        uuid: uuid
      })

      if (res.code === 0) {
        const intervalID = runningFileProcess.value[uuid].intervalID || null

        runningFileProcess.value = {
          ...runningFileProcess.value,
          [uuid]: {
            ...res.data,
            intervalID
          }
        }

        if (parseInt(res.data.percent) === 100) {
          if (intervalID) {
            clearInterval(intervalID)
          }
        }
      } else {
        const intervalID = runningFileProcess.value[uuid].intervalID || null

        runningFileProcess.value = {
          ...runningFileProcess.value,
          [uuid]: {
            ...runningFileProcess.value[uuid],
            percent: 100
          }
        }

        if (intervalID) {
          clearInterval(intervalID)
        }
      }
    } catch (error) {
      const intervalID = runningFileProcess.value[uuid].intervalID || null

      if (runningFileProcess.value[uuid].intervalID) {
        clearInterval(runningFileProcess.value[uuid].intervalID)
      }

      runningFileProcess.value = {
        ...runningFileProcess.value,
        [uuid]: {
          ...runningFileProcess.value[uuid],
          intervalID,
          status: 'error',
          msg: 'Xảy ra lỗi khi xử lý tập tin excel'
        }
      }
    }
  }

  return {
    fileProcess,
    runningFileProcess,
    displayFileProcess,
    initFileProcess,
    initFileProcessMemoized,
    addFileProcess,
    removeFileProcess
  }
})
