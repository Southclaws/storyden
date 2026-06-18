import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemProperties, SystemStyleObject } from '../types/system';

export interface BleedProperties {
  block?: SystemProperties["marginBlock"]
  inline?: SystemProperties["marginInline"]
  className?: string
}

type BleedRestStyles = Omit<SystemStyleObject, keyof BleedProperties>

interface BleedStyles extends BleedProperties, BleedRestStyles {}

interface BleedPatternFn {
  (styles?: BleedStyles): string
  raw: (styles?: BleedStyles) => SystemStyleObject
}

export declare function bleedRaw(styles?: BleedStyles): SystemStyleObject;

export declare const bleed: BleedPatternFn;