<template>
  <div class="system">
    <el-form ref="form" :model="config" label-width="240px">
      <!--  System start  -->
      <el-collapse v-model="activeNames">
        <el-collapse-item title="System Configuration" name="1">
          <el-form-item label="Port Value">
            <el-input v-model.number="config.system.addr" />
          </el-form-item>
          <el-form-item label="Database Type">
            <el-select v-model="config.system['db-type']" style="width:100%">
              <el-option value="mysql" />
              <el-option value="pgsql" />
            </el-select>
          </el-form-item>
          <el-form-item label="Oss Type">
            <el-select v-model="config.system['oss-type']" style="width:100%">
              <el-option value="local" />
              <el-option value="qiniu" />
              <el-option value="tencent-cos" />
              <el-option value="aliyun-oss" />
              <el-option value="huawei-obs" />
              <el-option value="cloudflare-r2" />
            </el-select>
          </el-form-item>
          <el-form-item label="Multiple Login Interception">
            <el-checkbox v-model="config.system['use-multipoint']">Enable</el-checkbox>
          </el-form-item>
          <el-form-item label="Enable Redis">
            <el-checkbox v-model="config.system['use-redis']">Enable</el-checkbox>
          </el-form-item>
          <el-form-item label="Rate Limit Count">
            <el-input-number v-model.number="config.system['iplimit-count']" />
          </el-form-item>
          <el-form-item label="Rate Limit Time">
            <el-input-number v-model.number="config.system['iplimit-time']" />
          </el-form-item>
          <el-tooltip content="After modification, please also modify VITE_BASE_PATH in the frontend env environment"
            placement="top-start">
            <el-form-item label="Global Router Prefix">
              <el-input v-model="config.system['router-prefix']" />
            </el-form-item>
          </el-tooltip>
        </el-collapse-item>
        <el-collapse-item title="JWT Signature" name="2">
          <el-form-item label="JWT Signature">
            <el-input v-model="config.jwt['signing-key']" />
          </el-form-item>
          <el-form-item label="Expiration Time">
            <el-input v-model="config.jwt['expires-time']" />
          </el-form-item>
          <el-form-item label="Buffer Time">
            <el-input v-model="config.jwt['buffer-time']" />
          </el-form-item>
          <el-form-item label="Issuer">
            <el-input v-model="config.jwt.issuer" />
          </el-form-item>
        </el-collapse-item>
        <el-collapse-item title="Zap Log Configuration" name="3">
          <el-form-item label="Level">
            <el-input v-model.number="config.zap.level" />
          </el-form-item>
          <el-form-item label="Output">
            <el-input v-model="config.zap.format" />
          </el-form-item>
          <el-form-item label="Log Prefix">
            <el-input v-model="config.zap.prefix" />
          </el-form-item>
          <el-form-item label="Log Folder">
            <el-input v-model="config.zap.director" />
          </el-form-item>
          <el-form-item label="Encoding Level">
            <el-input v-model="config.zap['encode-level']" />
          </el-form-item>
          <el-form-item label="Stack Name">
            <el-input v-model="config.zap['stacktrace-key']" />
          </el-form-item>
          <el-form-item label="Log Retention Time (default in days)">
            <el-input v-model.number="config.zap['retention-day']" />
          </el-form-item>
          <el-form-item label="Show Line">
            <el-checkbox v-model="config.zap['show-line']" />
          </el-form-item>
          <el-form-item label="Log to Console">
            <el-checkbox v-model="config.zap['log-in-console']" />
          </el-form-item>
        </el-collapse-item>
        <el-collapse-item title="Redis Admin Database Configuration" name="4">
          <el-form-item label="Database">
            <el-input v-model.number="config.redis.db" />
          </el-form-item>
          <el-form-item label="Address">
            <el-input v-model="config.redis.addr" />
          </el-form-item>
          <el-form-item label="Password">
            <el-input v-model="config.redis.password" />
          </el-form-item>
        </el-collapse-item>

        <el-collapse-item title="Mongo Database Configuration" name="14">
          <el-form-item label="Collection Name (usually not specified)">
            <el-input v-model="config.mongo.coll" />
          </el-form-item>
          <el-form-item label="MongoDB Options">
            <el-input v-model="config.mongo.options" />
          </el-form-item>
          <el-form-item label="Database Name">
            <el-input v-model="config.mongo.database" />
          </el-form-item>
          <el-form-item label="Username">
            <el-input v-model="config.mongo.username" />
          </el-form-item>
          <el-form-item label="Password">
            <el-input v-model="config.mongo.password" />
          </el-form-item>
          <el-form-item label="Minimum Connection Pool">
            <el-input v-model="config.mongo['min-pool-size']" />
          </el-form-item>
          <el-form-item label="Maximum Connection Pool">
            <el-input v-model="config.mongo['max-pool-size']" />
          </el-form-item>
          <el-form-item label="Socket Timeout">
            <el-input v-model="config.mongo['socket-timeout-ms']" />
          </el-form-item>
          <el-form-item label="Connection Timeout">
            <el-input v-model="config.mongo['socket-timeout-ms']" />
          </el-form-item>
          <el-form-item label="Enable Zap Log">
            <el-checkbox v-model="config.mongo['is-zap']" />
          </el-form-item>
          <el-form-item label="Hosts">
            <template v-for="(item, k) in config.mongo.hosts">
              <div v-for="(_, k2) in item" :key="k2">
                <el-form-item :key="k + k2" :label="k2">
                  <el-input v-model="item[k2]" />
                </el-form-item>
              </div>
            </template>
          </el-form-item>
        </el-collapse-item>

        <el-collapse-item title="Email Configuration" name="5">
          <el-form-item label="Recipient Email">
            <el-input v-model="config.email.to" placeholder="Multiple emails, separated by commas" />
          </el-form-item>
          <el-form-item label="Port">
            <el-input v-model.number="config.email.port" />
          </el-form-item>
          <el-form-item label="Sender Email">
            <el-input v-model="config.email.from" />
          </el-form-item>
          <el-form-item label="Host">
            <el-input v-model="config.email.host" />
          </el-form-item>
          <el-form-item label="Use SSL">
            <el-checkbox v-model="config.email['is-ssl']" />
          </el-form-item>
          <el-form-item label="Secret">
            <el-input v-model="config.email.secret" />
          </el-form-item>
          <el-form-item label="Test Email">
            <el-button @click="email">Test Email</el-button>
          </el-form-item>
        </el-collapse-item>
        <el-collapse-item title="Captcha Configuration" name="7">
          <el-form-item label="Character Length">
            <el-input v-model.number="config.captcha['key-long']" />
          </el-form-item>
          <el-form-item label="Image Width">
            <el-input v-model.number="config.captcha['img-width']" />
          </el-form-item>
          <el-form-item label="Image Height">
            <el-input v-model.number="config.captcha['img-height']" />
          </el-form-item>
        </el-collapse-item>
        <el-collapse-item title="Database Configuration" name="9">
          <template v-if="config.system['db-type'] === 'mysql'">
            <el-form-item label="Username">
              <el-input v-model="config.mysql.username" />
            </el-form-item>
            <el-form-item label="Password">
              <el-input v-model="config.mysql.password" />
            </el-form-item>
            <el-form-item label="Address">
              <el-input v-model="config.mysql.path" />
            </el-form-item>
            <el-form-item label="Database">
              <el-input v-model="config.mysql['db-name']" />
            </el-form-item>
            <el-form-item label="Prefix">
              <el-input v-model="config.mysql['refix']" />
            </el-form-item>
            <el-form-item label="Plural Tables">
              <el-switch v-model="config.mysql['singular']" />
            </el-form-item>
            <el-form-item label="Engine">
              <el-input v-model="config.mysql['engine']" />
            </el-form-item>
            <el-form-item label="maxIdleConns">
              <el-input v-model.number="config.mysql['max-idle-conns']" />
            </el-form-item>
            <el-form-item label="maxOpenConns">
              <el-input v-model.number="config.mysql['max-open-conns']" />
            </el-form-item>
            <el-form-item label="Write Log">
              <el-checkbox v-model="config.mysql['log-zap']" />
            </el-form-item>
            <el-form-item label="Log Mode">
              <el-input v-model="config.mysql['log-mode']" />
            </el-form-item>
          </template>
          <template v-if="config.system['db-type'] === 'pgsql'">
            <el-form-item label="Username">
              <el-input v-model="config.pgsql.username" />
            </el-form-item>
            <el-form-item label="Password">
              <el-input v-model="config.pgsql.password" />
            </el-form-item>
            <el-form-item label="Address">
              <el-input v-model="config.pgsql.path" />
            </el-form-item>
            <el-form-item label="Database">
              <el-input v-model="config.pgsql.dbname" />
            </el-form-item>
            <el-form-item label="Prefix">
              <el-input v-model="config.pgsql['refix']" />
            </el-form-item>
            <el-form-item label="Plural Tables">
              <el-switch v-model="config.pgsql['singular']" />
            </el-form-item>
            <el-form-item label="Engine">
              <el-input v-model="config.pgsql['engine']" />
            </el-form-item>
            <el-form-item label="maxIdleConns">
              <el-input v-model.number="config.pgsql['max-idle-conns']" />
            </el-form-item>
            <el-form-item label="maxOpenConns">
              <el-input v-model.number="config.pgsql['max-open-conns']" />
            </el-form-item>
            <el-form-item label="Write Log">
              <el-checkbox v-model="config.pgsql['log-zap']" />
            </el-form-item>
            <el-form-item label="Log Mode">
              <el-input v-model="config.pgsql['log-mode']" />
            </el-form-item>
          </template>
        </el-collapse-item>

        <el-collapse-item title="OSS Configuration" name="10">
          <template v-if="config.system['oss-type'] === 'local'">
            <h2>Local File Configuration</h2>
            <el-form-item label="Local File Access Path">
              <el-input v-model="config.local.path" />
            </el-form-item>
            <el-form-item label="Local File Storage Path">
              <el-input v-model="config.local['store-path']" />
            </el-form-item>
          </template>
          <template v-if="config.system['oss-type'] === 'qiniu'">
            <h2>Qiniu Upload Configuration</h2>
            <el-form-item label="Storage Region">
              <el-input v-model="config.qiniu.zone" />
            </el-form-item>
            <el-form-item label="Bucket Name">
              <el-input v-model="config.qiniu.bucket" />
            </el-form-item>
            <el-form-item label="CDN Acceleration Domain">
              <el-input v-model="config.qiniu['img-path']" />
            </el-form-item>
            <el-form-item label="Use HTTPS">
              <el-checkbox v-model="config.qiniu['use-https']">Enable</el-checkbox>
            </el-form-item>
            <el-form-item label="Access Key">
              <el-input v-model="config.qiniu['access-key']" />
            </el-form-item>
            <el-form-item label="Secret Key">
              <el-input v-model="config.qiniu['secret-key']" />
            </el-form-item>
            <el-form-item label="Use CDN Upload Acceleration">
              <el-checkbox v-model="config.qiniu['use-cdn-domains']">Enable</el-checkbox>
            </el-form-item>
          </template>
          <template v-if="config.system['oss-type'] === 'tencent-cos'">
            <h2>Tencent Cloud COS Upload Configuration</h2>
            <el-form-item label="Bucket Name">
              <el-input v-model="config['tencent-cos']['bucket']" />
            </el-form-item>
            <el-form-item label="Region">
              <el-input v-model="config['tencent-cos'].region" />
            </el-form-item>
            <el-form-item label="Secret ID">
              <el-input v-model="config['tencent-cos']['secret-id']" />
            </el-form-item>
            <el-form-item label="Secret Key">
              <el-input v-model="config['tencent-cos']['secret-key']" />
            </el-form-item>
            <el-form-item label="Path Prefix">
              <el-input v-model="config['tencent-cos']['path-prefix']" />
            </el-form-item>
            <el-form-item label="Base URL">
              <el-input v-model="config['tencent-cos']['base-url']" />
            </el-form-item>
          </template>
          <template v-if="config.system['oss-type'] === 'aliyun-oss'">
            <h2>Aliyun OSS Upload Configuration</h2>
            <el-form-item label="Region">
              <el-input v-model="config['aliyun-oss'].endpoint" />
            </el-form-item>
            <el-form-item label="Access Key ID">
              <el-input v-model="config['aliyun-oss']['access-key-id']" />
            </el-form-item>
            <el-form-item label="Access Key Secret">
              <el-input v-model="config['aliyun-oss']['access-key-secret']" />
            </el-form-item>
            <el-form-item label="Bucket Name">
              <el-input v-model="config['aliyun-oss']['bucket-name']" />
            </el-form-item>
            <el-form-item label="Bucket URL">
              <el-input v-model="config['aliyun-oss']['bucket-url']" />
            </el-form-item>
          </template>
          <template v-if="config.system['oss-type'] === 'huawei-obs'">
            <h2>Huawei Cloud OBS Upload Configuration</h2>
            <el-form-item label="Path">
              <el-input v-model="config['hua-wei-obs'].path" />
            </el-form-item>
            <el-form-item label="Bucket Name">
              <el-input v-model="config['hua-wei-obs'].bucket" />
            </el-form-item>
            <el-form-item label="Region">
              <el-input v-model="config['hua-wei-obs'].endpoint" />
            </el-form-item>
            <el-form-item label="Access Key">
              <el-input v-model="config['hua-wei-obs']['access-key']" />
            </el-form-item>
            <el-form-item label="Secret Key">
              <el-input v-model="config['hua-wei-obs']['secret-key']" />
            </el-form-item>
          </template>
          <template v-if="config.system['oss-type'] === 'cloudflare-r2'">
            <h2>Cloudflare R2 Upload Configuration</h2>
            <el-form-item label="Path">
              <el-input v-model="config['cloudflare-r2'].path" />
            </el-form-item>
            <el-form-item label="Bucket Name">
              <el-input v-model="config['cloudflare-r2'].bucket" />
            </el-form-item>
            <el-form-item label="Base URL">
              <el-input v-model="config['cloudflare-r2']['base-url']" />
            </el-form-item>
            <el-form-item label="Account ID">
              <el-input v-model="config['cloudflare-r2']['account-id']" />
            </el-form-item>
            <el-form-item label="Access Key ID">
              <el-input v-model="config['cloudflare-r2']['access-key-id']" />
            </el-form-item>
            <el-form-item label="Secret Access Key">
              <el-input v-model="config['cloudflare-r2']['secret-access-key']" />
            </el-form-item>
          </template>

        </el-collapse-item>

        <el-collapse-item title="Excel Upload Configuration" name="11">
          <el-form-item label="Target Directory">
            <el-input v-model="config.excel.dir" />
          </el-form-item>
        </el-collapse-item>

        <el-collapse-item title="Automated Code Configuration" name="12">
          <el-form-item label="Auto Restart (Linux)">
            <el-checkbox v-model="config.autocode['transfer-restart']" />
          </el-form-item>
          <el-form-item label="Root (Project Root Path)">
            <el-input v-model="config.autocode.root" disabled />
          </el-form-item>
          <el-form-item label="Server(Backend code address)">
            <el-input v-model="config.autocode['transfer-restart']" />
          </el-form-item>
          <el-form-item label="SApi(Backend API folder address)">
            <el-input v-model="config.autocode['server-api']" />
          </el-form-item>
          <el-form-item label="SInitialize(Backend Initialize folder)">
            <el-input v-model="config.autocode['server-initialize']" />
          </el-form-item>
          <el-form-item label="SModel(Backend Model file address)">
            <el-input v-model="config.autocode['server-model']" />
          </el-form-item>
          <el-form-item label="SRequest(Backend Request folder address)">
            <el-input v-model="config.autocode['server-request']" />
          </el-form-item>
          <el-form-item label="SRouter(Backend Router folder address)">
            <el-input v-model="config.autocode['server-router']" />
          </el-form-item>
          <el-form-item label="SService(Backend Service folder address)">
            <el-input v-model="config.autocode['server-service']" />
          </el-form-item>
          <el-form-item label="Web(Front-end folder address)">
            <el-input v-model="config.autocode.web" />
          </el-form-item>
          <el-form-item label="WApi(Backend WApi folder address)">
            <el-input v-model="config.autocode['web-api']" />
          </el-form-item>
          <el-form-item label="WForm(Backend WForm folder address)">
            <el-input v-model="config.autocode['web-form']" />
          </el-form-item>
          <el-form-item label="WTable(Backend WTable folder address)">
            <el-input v-model="config.autocode['web-table']" />
          </el-form-item>
        </el-collapse-item>
      </el-collapse>
    </el-form>
    <div class="mt-4">
      <el-button type="primary" @click="update">Update Now</el-button>
      <el-button type="primary" @click="reload">Restart the service (under development)</el-button>
    </div>
  </div>
</template>

<script setup>
import { getSystemConfig, setSystemConfig } from '@/api/system'
import { emailTest } from '@/api/email'
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'

defineOptions({
  name: 'Config'
})

const activeNames = reactive([])
const config = ref({
  system: {
    'iplimit-count': 0,
    'iplimit-time': 0
  },
  jwt: {},
  mysql: {},
  pgsql: {},
  excel: {},
  autocode: {},
  redis: {},
  mongo: {
    coll: '',
    options: '',
    database: '',
    username: '',
    password: '',
    'min-pool-size': '',
    'max-pool-size': '',
    'socket-timeout-ms': '',
    'connect-timeout-ms': '',
    'is-zap': '',
    hosts: [
      {
        host: '',
        port: ''
      }
    ]
  },
  qiniu: {},
  'tencent-cos': {},
  'aliyun-oss': {},
  'hua-wei-obs': {},
  'cloudflare-r2': {},
  captcha: {},
  zap: {},
  local: {},
  email: {},
  timer: {
    detail: {}
  }
})

const initForm = async () => {
  const res = await getSystemConfig()
  if (res.code === 0) {
    config.value = res.data.config
  }
}
initForm()
const reload = () => { }
const update = async () => {
  const res = await setSystemConfig({ config: config.value })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '配置文件Thiết lập thành công'
    })
    await initForm()
  }
}
const email = async () => {
  const res = await emailTest()
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '邮件发送成功'
    })
    await initForm()
  } else {
    ElMessage({
      type: 'error',
      message: '邮件发送失败'
    })
  }
}

</script>

<style lang="scss">
.system {
  @apply bg-white p-9 rounded dark:bg-slate-900;

  h2 {
    @apply p-2.5 my-2.5 text-lg shadow;
  }
}
</style>
