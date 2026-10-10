import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemStyleObject } from '../types/system';

export interface BoxProperties {}

type BoxRestStyles = Omit<SystemStyleObject, keyof BoxProperties>

interface BoxStyles extends BoxProperties, BoxRestStyles {}

interface BoxPatternFn {
  (styles?: BoxStyles): string
  raw: (styles?: BoxStyles) => SystemStyleObject
  propKeys: Array<keyof BoxProperties>
}

export declare function boxRaw(styles?: BoxStyles): SystemStyleObject;

export declare const box: BoxPatternFn;