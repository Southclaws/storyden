import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemProperties, SystemStyleObject } from '../types/system';

export interface CircleProperties {
  size?: SystemProperties["width"]
}

type CircleRestStyles = Omit<SystemStyleObject, keyof CircleProperties>

interface CircleStyles extends CircleProperties, CircleRestStyles {}

interface CirclePatternFn {
  (styles?: CircleStyles): string
  raw: (styles?: CircleStyles) => SystemStyleObject
  propKeys: Array<keyof CircleProperties>
}

export declare function circleRaw(styles?: CircleStyles): SystemStyleObject;

export declare const circle: CirclePatternFn;