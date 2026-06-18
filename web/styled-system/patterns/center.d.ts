import type { PatternRuntimeConfig } from '../types/pattern';
import type { ConditionalValue, SystemStyleObject } from '../types/system';

export interface CenterProperties {
  inline?: ConditionalValue<boolean>
  className?: string
}

type CenterRestStyles = Omit<SystemStyleObject, keyof CenterProperties>

interface CenterStyles extends CenterProperties, CenterRestStyles {}

interface CenterPatternFn {
  (styles?: CenterStyles): string
  raw: (styles?: CenterStyles) => SystemStyleObject
}

export declare function centerRaw(styles?: CenterStyles): SystemStyleObject;

export declare const center: CenterPatternFn;