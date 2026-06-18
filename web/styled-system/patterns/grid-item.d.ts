import type { PatternRuntimeConfig } from '../types/pattern';
import type { ConditionalValue, SystemStyleObject } from '../types/system';

export interface GridItemProperties {
  colEnd?: ConditionalValue<number>
  colSpan?: ConditionalValue<number>
  colStart?: ConditionalValue<number>
  rowEnd?: ConditionalValue<number>
  rowSpan?: ConditionalValue<number>
  rowStart?: ConditionalValue<number>
  className?: string
}

type GridItemRestStyles = Omit<SystemStyleObject, keyof GridItemProperties>

interface GridItemStyles extends GridItemProperties, GridItemRestStyles {}

interface GridItemPatternFn {
  (styles?: GridItemStyles): string
  raw: (styles?: GridItemStyles) => SystemStyleObject
}

export declare function gridItemRaw(styles?: GridItemStyles): SystemStyleObject;

export declare const gridItem: GridItemPatternFn;