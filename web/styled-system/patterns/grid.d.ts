import type { PatternRuntimeConfig } from '../types/pattern';
import type { TokenValue } from '../types/tokens';
import type { ConditionalValue, SystemProperties, SystemStyleObject } from '../types/system';

export interface GridProperties {
  columnGap?: SystemProperties["gap"]
  columns?: ConditionalValue<number>
  gap?: SystemProperties["gap"]
  minChildWidth?: ConditionalValue<TokenValue<"sizes"> | SystemProperties["width"]>
  rowGap?: SystemProperties["gap"]
}

type GridRestStyles = Omit<SystemStyleObject, keyof GridProperties>

interface GridStyles extends GridProperties, GridRestStyles {}

interface GridPatternFn {
  (styles?: GridStyles): string
  raw: (styles?: GridStyles) => SystemStyleObject
  propKeys: Array<keyof GridProperties>
}

export declare function gridRaw(styles?: GridStyles): SystemStyleObject;

export declare const grid: GridPatternFn;