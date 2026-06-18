import type { Token, TokenPath } from '../types/tokens';

interface TokenFn {
  (path: TokenPath, fallback?: string): string
  var: (path: Token, fallback?: string) => string
}

export declare const token: TokenFn;