import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemProperties, SystemStyleObject } from '../types/system';

export interface FlexProperties {
  align?: SystemProperties["alignItems"]
  basis?: SystemProperties["flexBasis"]
  direction?: SystemProperties["flexDirection"]
  grow?: SystemProperties["flexGrow"]
  justify?: SystemProperties["justifyContent"]
  shrink?: SystemProperties["flexShrink"]
  wrap?: SystemProperties["flexWrap"]
  className?: string
}

type FlexRestStyles = Omit<SystemStyleObject, keyof FlexProperties>

interface FlexStyles extends FlexProperties, FlexRestStyles {}

interface FlexPatternFn {
  (styles?: FlexStyles): string
  raw: (styles?: FlexStyles) => SystemStyleObject
}

export declare function flexRaw(styles?: FlexStyles): SystemStyleObject;

export declare const flex: FlexPatternFn;