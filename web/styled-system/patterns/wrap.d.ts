import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemProperties, SystemStyleObject } from '../types/system';

export interface WrapProperties {
  align?: SystemProperties["alignItems"]
  columnGap?: SystemProperties["gap"]
  gap?: SystemProperties["gap"]
  justify?: SystemProperties["justifyContent"]
  rowGap?: SystemProperties["gap"]
  className?: string
}

type WrapRestStyles = Omit<SystemStyleObject, keyof WrapProperties>

interface WrapStyles extends WrapProperties, WrapRestStyles {}

interface WrapPatternFn {
  (styles?: WrapStyles): string
  raw: (styles?: WrapStyles) => SystemStyleObject
}

export declare function wrapRaw(styles?: WrapStyles): SystemStyleObject;

export declare const wrap: WrapPatternFn;