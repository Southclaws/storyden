import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemStyleObject } from '../types/system';

export interface FloatingProperties {}

type FloatingRestStyles = Omit<SystemStyleObject, keyof FloatingProperties>

interface FloatingStyles extends FloatingProperties, FloatingRestStyles {}

interface FloatingPatternFn {
  (styles?: FloatingStyles): string
  raw: (styles?: FloatingStyles) => SystemStyleObject
  propKeys: Array<keyof FloatingProperties>
}

export declare function floatingRaw(styles?: FloatingStyles): SystemStyleObject;

export declare const floating: FloatingPatternFn;