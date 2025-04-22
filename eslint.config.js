import eslint from '@eslint/js';
import tseslint from 'typescript-eslint';
import angular from 'angular-eslint';
import nxPlugin from '@nx/eslint-plugin';
import prettierConfig from 'eslint-config-prettier';

import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

export default tseslint.config(
  {
    // Common ignores
    ignores: [
      '**/node_modules/**',
      '**/dist/**',
      '**/tmp/**',
      '**/.angular/**',
      '**/.nx/cache/**',
      // Ignore config files in root? Usually not needed for ESLint >= 8.35.0
      // 'eslint.config.js',
      // 'prettier.config.js',
      // 'nx.json',
      // etc.
    ],
  },
  {
    plugins: {
      '@nx': nxPlugin,
    },
    rules: {
      '@nx/enforce-module-boundaries': [
        'error',
        {
          enforceBuildableLibDependency: true,
          allow: [],
          depConstraints: [
            {
              sourceTag: '*',
              onlyDependOnLibsWithTags: ['*'],
            },
          ],
        },
      ],
      'no-extra-semi': 'error', // Apply globally from root config
    },
  },
  // JavaScript files configuration
  {
    files: ['**/*.js', '**/*.jsx'],
    extends: [eslint.configs.recommended],
    rules: {
      // Rules specific to JS files if any, besides recommended
    },
  },
  // Angular specific configurations (for apps/webapp)
  {
    files: ['apps/webapp/**/*.ts'],
    extends: [
      eslint.configs.recommended,
      ...tseslint.configs.recommended,
      ...tseslint.configs.stylistic,
      ...angular.configs.tsRecommended,
    ],
    processor: angular.processInlineTemplates,
    rules: {
      // our project thinks using renaming inputs is ok
      '@angular-eslint/no-input-rename': 'off',
    },
  },
  // Angular template specific configurations (for apps/webapp)
  {
    files: ['apps/webapp/**/*.html'],
    extends: [
      ...angular.configs.templateRecommended, // Recommended rules for templates
      ...angular.configs.templateAccessibility, // Accessibility rules
    ],
    rules: {},
  },
  // Prettier configuration (must be last to override other style rules)
  prettierConfig,
);
