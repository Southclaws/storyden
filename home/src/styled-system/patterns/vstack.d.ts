import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemProperties, SystemStyleObject } from '../types/system';

export interface VstackProperties {
  gap?: SystemProperties["gap"]
  justify?: SystemProperties["justifyContent"]
}

type VstackRestStyles = Omit<SystemStyleObject, keyof VstackProperties>

interface VstackStyles extends VstackProperties, VstackRestStyles {}

interface VstackPatternFn {
  (styles?: VstackStyles): string
  raw: (styles?: VstackStyles) => SystemStyleObject
  propKeys: Array<keyof VstackProperties>
}

export declare function vstackRaw(styles?: VstackStyles): SystemStyleObject;

export declare const vstack: VstackPatternFn;