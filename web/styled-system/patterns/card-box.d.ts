import type { PatternRuntimeConfig } from '../types/pattern';
import type { ConditionalValue, SystemProperties, SystemStyleObject } from '../types/system';

export interface CardBoxProperties {
  display?: SystemProperties["display"]
  kind?: ConditionalValue<"edge" | "default">
  className?: string
}

type CardBoxRestStyles = Omit<SystemStyleObject, keyof CardBoxProperties>

interface CardBoxStyles extends CardBoxProperties, CardBoxRestStyles {}

interface CardBoxPatternFn {
  (styles?: CardBoxStyles): string
  raw: (styles?: CardBoxStyles) => SystemStyleObject
}

export declare function CardBoxRaw(styles?: CardBoxStyles): SystemStyleObject;

export declare const CardBox: CardBoxPatternFn;