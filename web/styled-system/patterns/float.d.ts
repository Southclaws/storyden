import type { PatternRuntimeConfig } from '../types/pattern';
import type { TokenValue } from '../types/tokens';
import type { ConditionalValue, SystemProperties, SystemStyleObject } from '../types/system';

export interface FloatProperties {
  offset?: ConditionalValue<TokenValue<"spacing"> | SystemProperties["top"]>
  offsetX?: ConditionalValue<TokenValue<"spacing"> | SystemProperties["left"]>
  offsetY?: ConditionalValue<TokenValue<"spacing"> | SystemProperties["top"]>
  placement?: ConditionalValue<"bottom-end" | "bottom-start" | "top-end" | "top-start" | "bottom-center" | "top-center" | "middle-center" | "middle-end" | "middle-start">
  className?: string
}

type FloatRestStyles = Omit<SystemStyleObject, keyof FloatProperties>

interface FloatStyles extends FloatProperties, FloatRestStyles {}

interface FloatPatternFn {
  (styles?: FloatStyles): string
  raw: (styles?: FloatStyles) => SystemStyleObject
}

export declare function floatRaw(styles?: FloatStyles): SystemStyleObject;

export declare const float: FloatPatternFn;