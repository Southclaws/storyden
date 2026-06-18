import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemStyleObject } from '../types/system';

export interface WstackProperties {
  className?: string
}

type WstackRestStyles = Omit<SystemStyleObject, keyof WstackProperties>

interface WstackStyles extends WstackProperties, WstackRestStyles {}

interface WstackPatternFn {
  (styles?: WstackStyles): string
  raw: (styles?: WstackStyles) => SystemStyleObject
}

export declare function wstackRaw(styles?: WstackStyles): SystemStyleObject;

export declare const wstack: WstackPatternFn;