'use strict';
const {test} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');

const root = path.resolve(__dirname, '..');
const sdk = fs.readFileSync(path.join(root, 'src/types.ts'), 'utf8');
const validator = fs.readFileSync(path.join(root, 'src/utils/validators.ts'), 'utf8');
const native = fs.readFileSync(path.join(root, 'agi-native-language.ts'), 'utf8');

function matches(text, re) {
  const found = text.match(re);
  assert.ok(found, 'Expected language registry declaration not found');
  return found[1];
}
function quotedCodes(text) {
  return [...text.matchAll(/['"]([a-z]{2,3})['"]/g)].map(x => x[1]);
}
test('SDK, validator and native-language engine accept identical language codes', () => {
  const sdkUnion = matches(sdk, /export type SupportedLanguage\s*=([\s\S]*?);/);
  const engineUnion = matches(native, /export type SupportedLanguage\s*=([\s\S]*?);/);
  const validatorArray = matches(validator, /const validLanguages\s*=\s*\[([\s\S]*?)\]/);
  const source = new Set(quotedCodes(sdkUnion));
  assert.deepEqual([...new Set(quotedCodes(engineUnion))].sort(), [...source].sort());
  assert.deepEqual([...new Set(quotedCodes(validatorArray))].sort(), [...source].sort());
  assert.ok(source.has('zip'), 'Zipser German internal code must remain supported');
});
test('native-language map and guidance cover every supported code', () => {
  const union = new Set(quotedCodes(matches(native, /export type SupportedLanguage\s*=([\s\S]*?);/)));
  const map = matches(native, /export const LANGUAGE_MAP:[^=]*=\s*\{([\s\S]*?)\};/);
  const guidance = matches(native, /const guidelines:[^=]*=\s*\{([\s\S]*?)\};/);
  for (const code of union) {
    assert.match(map, new RegExp('(?:^|\\n)\\s*' + code + ':'));
    assert.match(guidance, new RegExp('(?:^|\\n)\\s*' + code + ':'));
  }
});
