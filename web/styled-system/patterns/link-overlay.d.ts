import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemStyleObject } from '../types/system';

export interface LinkOverlayProperties {
  className?: string
}

type LinkOverlayRestStyles = Omit<SystemStyleObject, keyof LinkOverlayProperties>

interface LinkOverlayStyles extends LinkOverlayProperties, LinkOverlayRestStyles {}

interface LinkOverlayPatternFn {
  (styles?: LinkOverlayStyles): string
  raw: (styles?: LinkOverlayStyles) => SystemStyleObject
}

export declare function linkOverlayRaw(styles?: LinkOverlayStyles): SystemStyleObject;

export declare const linkOverlay: LinkOverlayPatternFn;