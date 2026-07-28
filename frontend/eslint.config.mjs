import antfu from '@antfu/eslint-config'

export default antfu(
  {
    ignores: ['node_modules', 'dist', 'public'],

    // 覆盖antfu的规则
    rules: {
      'style/comma-dangle': ['error', 'never'],
      'style/brace-style': ['error', '1tbs'],
      'style/quote-props': ['error', 'as-needed'],
      'object-shorthand': 'off',
      'antfu/top-level-function': 'off',
      'vue/singleline-html-element-content-newline': 'off',
      'ts/no-unsafe-function-type': 'off',
      'node/prefer-global/process': 'off',
      'no-new': 'off'
    },

    pnpm: true
  }
)
