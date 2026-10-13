import js from "@eslint/js";
import tsParser from "@typescript-eslint/parser";
import reactPlugin from "eslint-plugin-react";
import reactHooksPlugin from "eslint-plugin-react-hooks";
import importPlugin from "eslint-plugin-import";
import prettierPlugin from "eslint-plugin-prettier";
import prettier from "eslint-config-prettier";
import { FlatCompat } from "@eslint/eslintrc";

const compat = new FlatCompat();

export default [
  js.configs.recommended,
  ...compat.extends("plugin:@typescript-eslint/recommended"),
  ...compat.extends("plugin:react-hooks/recommended"),
  ...compat.extends("plugin:prettier/recommended"),
  {
    languageOptions: {
      parser: tsParser,
      globals: {
        browser: true,
        node: true,
        es6: true,
      },
      parserOptions: {
        ecmaVersion: 6,
        sourceType: "module",
        ecmaFeatures: {
          legacyDecorators: true,
          modules: true,
        },
      },
    },
    plugins: {
      react: reactPlugin,
      "react-hooks": reactHooksPlugin,
      import: importPlugin,
      prettier: prettierPlugin,
    },
    settings: {
      react: {
        pragma: "React",
        version: "detect",
      },
    },
    rules: {
      ...reactPlugin.configs["recommended"].rules,
      ...reactHooksPlugin.configs["recommended"].rules,

      "prettier/prettier": "error",
      "no-console": 1,
      "lines-between-class-members": [
        "error",
        "always",
        {
          exceptAfterSingleLine: true,
        },
      ],
      "no-else-return": [
        "error",
        {
          allowElseIf: true,
        },
      ],
      "object-curly-spacing": ["error", "always"],
      "prefer-destructuring": [
        "error",
        {
          VariableDeclarator: {
            array: false,
            object: true,
          },
          AssignmentExpression: {
            array: false,
            object: false,
          },
        },
      ],
      "no-trailing-spaces": [
        "error",
        {
          skipBlankLines: true,
        },
      ],
      "comma-dangle": ["error", "always-multiline"],
      "arrow-spacing": [
        "error",
        {
          before: true,
          after: true,
        },
      ],
      "keyword-spacing": "error",
      "space-before-blocks": "error",
      "arrow-body-style": "off",
      "consistent-return": "off",
      "operator-linebreak": "off",
      "no-mixed-operators": "off",
      "implicit-arrow-linebreak": "off",
      "no-confusing-arrow": "off",
      "no-underscore-dangle": "off",
      "react/react-in-jsx-scope" : "off",
      "spaced-comment": "off",
      "default-case": "off",
      "no-continue": "off",
      "guard-for-in": "off",
      "no-restricted-syntax": "off",
      "no-restricted-globals": "off",
      "no-plusplus": "off",
      "no-case-declarations": "off",
      "arrow-parens": ["error", "always"],
      "space-infix-ops": "error",
      "no-return-assign": ["error", "always"],
      "@typescript-eslint/no-explicit-any": "off",
      "@typescript-eslint/ban-ts-ignore": "off",
      "@typescript-eslint/ban-types": "off",
      "@typescript-eslint/explicit-function-return-type": "off",
      "@typescript-eslint/no-use-before-define": "off",
      "@typescript-eslint/no-namespace": "off",
      "@typescript-eslint/interface-name-prefix": "off",
      "@typescript-eslint/no-duplicate-enum-values": "off",
      "@typescript-eslint/no-unused-vars": [
        "error",
        {
          argsIgnorePattern: "^_",
          varsIgnorePattern: "^_",
        },
      ],
      "react/destructuring-assignment": "off",
      "react/state-in-constructor": "off",
      "react/display-name": "off",
      "react/prop-types": "off",
      "react/sort-comp": "off",
      "react/prefer-stateless-function": "off",
      "react/no-did-update-set-state": "off",
      "react/no-children-prop": "error",
      "react/no-deprecated": "error",
      "react/self-closing-comp": "error",
      "react/no-direct-mutation-state": "error",
      "react/require-render-return": "error",
      "react/no-unsafe": "error",
      "react/no-array-index-key": "off",
      "react/static-property-placement": "off",
      "react/jsx-key": "error",
      "react/jsx-curly-brace-presence": [
        "error",
        {
          props: "never",
          children: "never",
        },
      ],
      "react/jsx-closing-bracket-location": [2, "tag-aligned"],
      "react/jsx-boolean-value": [
        "error",
        "never",
        {
          always: ["personal"],
        },
      ],
      "react/jsx-filename-extension": [
        "error",
        {
          extensions: [".tsx", ".jsx"],
        },
      ],
      "react/jsx-first-prop-new-line": ["error", "multiline-multiprop"],
      "react/jsx-fragments": ["error", "syntax"],
      "react/jsx-indent-props": ["error", 2],
      "react/jsx-one-expression-per-line": ["off"],
      "react/jsx-tag-spacing": "off",
      "@typescript-eslint/ban-ts-comment": "off",
      "react/jsx-props-no-spreading": "off",
      "react-hooks/exhaustive-deps": "error",
      "padded-blocks": ["error", "never"],
      "import/order": [
        2,
        {
          groups: ["builtin", "external", "internal", "parent", "sibling", "index"],
          pathGroups: [
            {
              pattern: "~/**/*",
              group: "internal",
            },
            {
              pattern: "*.scss",
              group: "index",
              patternOptions: { matchBase: true },
              position: "after",
            },
          ],
          pathGroupsExcludedImportTypes: [],
          "newlines-between": "always",
        },
      ],
      "import/prefer-default-export": "off",
      "import/no-duplicates": 2,
      "import/first": 2,
    },
  },
  prettier,
]