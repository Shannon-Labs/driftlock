module.exports = {
  root: true,
  env: {
    browser: true,
    node: true,
    es2021: true,
  },
  extends: [
    'eslint:recommended',
    'plugin:vue/vue3-essential',
    '@vue/eslint-config-typescript',
  ],
  parserOptions: {
    ecmaVersion: 'latest',
    sourceType: 'module',
  },
  rules: {
    // Warn on console statements but allow console.warn and console.error in dev guards
    'no-console': ['warn', { allow: ['warn', 'error'] }],
    // Require explicit return types is too strict for Vue components
    '@typescript-eslint/explicit-function-return-type': 'off',
    // Allow unused vars starting with underscore (common pattern for intentionally unused)
    '@typescript-eslint/no-unused-vars': ['warn', { argsIgnorePattern: '^_', varsIgnorePattern: '^_' }],
  },
}
