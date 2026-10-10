import type { PatternRuntimeConfig } from '../types/pattern';
import type { ConditionalValue, SystemProperties, SystemStyleObject } from '../types/system';

export interface CardProperties {
  display?: SystemProperties["display"]
  kind?: ConditionalValue<"edge" | "default">
}

type CardRestStyles = Omit<SystemStyleObject, keyof CardProperties>

interface CardStyles extends CardProperties, CardRestStyles {}

interface CardPatternFn {
  (styles?: CardStyles): string
  raw: (styles?: CardStyles) => SystemStyleObject
  propKeys: Array<keyof CardProperties>
}

export declare function cardRaw(styles?: CardStyles): SystemStyleObject;

export declare const card: CardPatternFn;