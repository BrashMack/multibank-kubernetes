declare module "*.scss";

declare module '*.png' {
  import { ImageSourcePropType } from 'react-native'
  const content: ImageSourcePropType
  export default content
}

declare module "*.svg" {
  import React from "react";
  import { SvgProps } from "react-native-svg";
  const content: React.FC<SvgProps>;
  export default content;
}

declare module '@eslint/eslintrc' {
  export { FlatCompat } from '@eslint/eslintrc/lib/flat-compat';
  export * from '@eslint/eslintrc/lib/shared/types';
}

declare module 'eslint-plugin-import' {
  import { ESLint } from 'eslint';
  const plugin: ESLint.Plugin;
  export default plugin;
}

declare module '@typescript-eslint/parser' {
  import { Linter } from 'eslint';
  const parser: Linter.ParserModule;
  export = parser;
}
