import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemStyleObject } from '../types/system';

export interface LstackProperties {
  className?: string
}

type LstackRestStyles = Omit<SystemStyleObject, keyof LstackProperties>

interface LstackStyles extends LstackProperties, LstackRestStyles {}

interface LstackPatternFn {
  (styles?: LstackStyles): string
  raw: (styles?: LstackStyles) => SystemStyleObject
}

export declare function lstackRaw(styles?: LstackStyles): SystemStyleObject;

export declare const lstack: LstackPatternFn;