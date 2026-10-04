import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemStyleObject } from '../types/system';

export interface LstackProperties {}

type LstackRestStyles = Omit<SystemStyleObject, keyof LstackProperties>

interface LstackStyles extends LstackProperties, LstackRestStyles {}

interface LstackPatternFn {
  (styles?: LstackStyles): string
  raw: (styles?: LstackStyles) => SystemStyleObject
  propKeys: Array<keyof LstackProperties>
}

export declare function lstackRaw(styles?: LstackStyles): SystemStyleObject;

export declare const lstack: LstackPatternFn;