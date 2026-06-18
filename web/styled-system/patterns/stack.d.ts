import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemProperties, SystemStyleObject } from '../types/system';

export interface StackProperties {
  align?: SystemProperties["alignItems"]
  direction?: SystemProperties["flexDirection"]
  gap?: SystemProperties["gap"]
  justify?: SystemProperties["justifyContent"]
  className?: string
}

type StackRestStyles = Omit<SystemStyleObject, keyof StackProperties>

interface StackStyles extends StackProperties, StackRestStyles {}

interface StackPatternFn {
  (styles?: StackStyles): string
  raw: (styles?: StackStyles) => SystemStyleObject
}

export declare function stackRaw(styles?: StackStyles): SystemStyleObject;

export declare const stack: StackPatternFn;