import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemProperties, SystemStyleObject } from '../types/system';

export interface SquareProperties {
  size?: SystemProperties["width"]
}

type SquareRestStyles = Omit<SystemStyleObject, keyof SquareProperties>

interface SquareStyles extends SquareProperties, SquareRestStyles {}

interface SquarePatternFn {
  (styles?: SquareStyles): string
  raw: (styles?: SquareStyles) => SystemStyleObject
  propKeys: Array<keyof SquareProperties>
}

export declare function squareRaw(styles?: SquareStyles): SystemStyleObject;

export declare const square: SquarePatternFn;