import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemProperties, SystemStyleObject } from '../types/system';

export interface VstackProperties {
  gap?: SystemProperties["gap"]
  justify?: SystemProperties["justifyContent"]
  className?: string
}

type VstackRestStyles = Omit<SystemStyleObject, keyof VstackProperties>

interface VstackStyles extends VstackProperties, VstackRestStyles {}

interface VstackPatternFn {
  (styles?: VstackStyles): string
  raw: (styles?: VstackStyles) => SystemStyleObject
}

export declare function vstackRaw(styles?: VstackStyles): SystemStyleObject;

export declare const vstack: VstackPatternFn;