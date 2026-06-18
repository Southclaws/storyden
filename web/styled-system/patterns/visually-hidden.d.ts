import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemStyleObject } from '../types/system';

export interface VisuallyHiddenProperties {
  className?: string
}

type VisuallyHiddenRestStyles = Omit<SystemStyleObject, keyof VisuallyHiddenProperties>

interface VisuallyHiddenStyles extends VisuallyHiddenProperties, VisuallyHiddenRestStyles {}

interface VisuallyHiddenPatternFn {
  (styles?: VisuallyHiddenStyles): string
  raw: (styles?: VisuallyHiddenStyles) => SystemStyleObject
}

export declare function visuallyHiddenRaw(styles?: VisuallyHiddenStyles): SystemStyleObject;

export declare const visuallyHidden: VisuallyHiddenPatternFn;