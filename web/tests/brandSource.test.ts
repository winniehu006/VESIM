import assert from 'node:assert/strict'
import { readFile } from 'node:fs/promises'
import test from 'node:test'

async function source(path: string) {
  return readFile(new URL(path, import.meta.url), 'utf8')
}

test('production Vue surfaces present the VESIM brand consistently', async () => {
  const [app, shell, login, loading, header, commands, automation, settings, atConsole, zhCN] = await Promise.all([
    source('../src/App.vue'),
    source('../src/layouts/AuthenticatedShell.vue'),
    source('../src/views/Login.vue'),
    source('../src/components/LoadingScreen.vue'),
    source('../src/components/PageHeader.vue'),
    source('../src/components/commands/CommandChat.vue'),
    source('../src/views/AutomaticTasks.vue'),
    source('../src/views/Settings.vue'),
    source('../src/components/DeviceAtTab.vue'),
    source('../src/i18n/zh-CN.ts')
  ])
  const visibleSources = [app, shell, login, loading, header, commands, automation, settings, atConsole]

  for (const content of visibleSources) {
    assert.doesNotMatch(content, /VoHive|VOHIVE/)
  }
  assert.match(shell, /sidebar-brand-icon">V</)
  assert.match(shell, /sidebar-brand-title">VESIM</)
  assert.match(shell, /topbar-product">VESIM</)
  assert.match(login, /identity-mark">V</)
  assert.match(login, /t\('login\.title'\)/)
  assert.match(zhCN, /login:\s*\{[\s\S]*?title:\s*'登录 VESIM'/)
  assert.match(app, /t\('disclaimer\.title'\)/)
  assert.match(app, /<el-input[\s\S]*?v-model="confirmText"/)
  assert.match(app, /<el-button[\s\S]*?rejectDisclaimer/)
  assert.doesNotMatch(app, /<input[\s\S]*?v-model="confirmText"/)
  assert.match(zhCN, /disclaimer:\s*\{[\s\S]*?title:\s*'VESIM 最终用户许可与免责声明'/)
  assert.match(commands, /VESIM 命令会话/)
  assert.match(settings, /title="VESIM Gateway"/)
})

test('frontend integration identifiers use the VESIM namespace', async () => {
  const [app, systemService, websheet, sensitive, phoneSession, settingsStore] = await Promise.all([
    source('../src/App.vue'),
    source('../src/services/system.ts'),
    source('../src/components/CarrierWebsheetDialog.vue'),
    source('../src/composables/useSensitiveVisibility.ts'),
    source('../src/services/phone-session.ts'),
    source('../src/stores/settings.ts')
  ])

  assert.doesNotMatch(app, /disclaimer_agreed_at|shouldShowDisclaimer/)
  assert.match(app, /getDisclaimerStatus/)
  assert.match(systemService, /put<DisclaimerStatus>\('\/settings\/disclaimer'/)
  assert.match(websheet, /vesim-websheet-callback/)
  assert.match(websheet, /vesim-websheet-complete/)
  assert.match(websheet, /vesim-websheet/)
  assert.match(sensitive, /vesim_show_sensitive/)
  assert.match(phoneSession, /vesim_phone_control/)
  assert.match(settingsStore, /x-vesim-signature/)
  assert.doesNotMatch([app, systemService, websheet, sensitive, phoneSession, settingsStore].join('\n'), /vohive/i)
})
