import assert from 'node:assert/strict'
import { readFile } from 'node:fs/promises'
import test from 'node:test'

async function asset(path: string) {
  return readFile(new URL(path, import.meta.url), 'utf8')
}

test('browser title and favicon use the VESIM identity', async () => {
  const [html, favicon] = await Promise.all([
    asset('../index.html'),
    asset('../public/favicon.svg')
  ])

  assert.match(html, /<title>VESIM<\/title>/)
  assert.doesNotMatch(html, /VoHive/)
  assert.match(favicon, /id="brand-gradient"/)
  assert.match(favicon, /#1fa574/)
  assert.match(favicon, /M9\.5 9 H13\.2 L16 19\.6 L18\.8 9 H22\.5/)
  assert.doesNotMatch(favicon, /6366f1|a855f7/)
})
