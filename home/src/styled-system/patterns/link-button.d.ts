import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemStyleObject } from '../types/system';

export interface LinkButtonProperties {}

type LinkButtonRestStyles = Omit<SystemStyleObject, keyof LinkButtonProperties>

interface LinkButtonStyles extends LinkButtonProperties, LinkButtonRestStyles {}

interface LinkButtonPatternFn {
  (styles?: LinkButtonStyles): string
  raw: (styles?: LinkButtonStyles) => SystemStyleObject
  propKeys: Array<keyof LinkButtonProperties>
}

export declare function linkButtonRaw(styles?: LinkButtonStyles): SystemStyleObject;

export declare const linkButton: LinkButtonPatternFn;