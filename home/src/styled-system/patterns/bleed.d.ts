import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemProperties, SystemStyleObject } from '../types/system';

export interface BleedProperties {
  block?: SystemProperties["marginBlock"]
  inline?: SystemProperties["marginInline"]
}

type BleedRestStyles = Omit<SystemStyleObject, keyof BleedProperties>

interface BleedStyles extends BleedProperties, BleedRestStyles {}

interface BleedPatternFn {
  (styles?: BleedStyles): string
  raw: (styles?: BleedStyles) => SystemStyleObject
  propKeys: Array<keyof BleedProperties>
}

export declare function bleedRaw(styles?: BleedStyles): SystemStyleObject;

export declare const bleed: BleedPatternFn;