import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemProperties, SystemStyleObject } from '../types/system';

export interface HstackProperties {
  gap?: SystemProperties["gap"]
  justify?: SystemProperties["justifyContent"]
}

type HstackRestStyles = Omit<SystemStyleObject, keyof HstackProperties>

interface HstackStyles extends HstackProperties, HstackRestStyles {}

interface HstackPatternFn {
  (styles?: HstackStyles): string
  raw: (styles?: HstackStyles) => SystemStyleObject
  propKeys: Array<keyof HstackProperties>
}

export declare function hstackRaw(styles?: HstackStyles): SystemStyleObject;

export declare const hstack: HstackPatternFn;