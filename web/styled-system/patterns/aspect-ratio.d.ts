import type { PatternRuntimeConfig } from '../types/pattern';
import type { ConditionalValue, SystemStyleObject } from '../types/system';

export interface AspectRatioProperties {
  ratio?: ConditionalValue<number>
  className?: string
}

type AspectRatioRestStyles = Omit<SystemStyleObject, keyof AspectRatioProperties | "aspectRatio">

interface AspectRatioStyles extends AspectRatioProperties, AspectRatioRestStyles {}

interface AspectRatioPatternFn {
  (styles?: AspectRatioStyles): string
  raw: (styles?: AspectRatioStyles) => SystemStyleObject
}

export declare function aspectRatioRaw(styles?: AspectRatioStyles): SystemStyleObject;

export declare const aspectRatio: AspectRatioPatternFn;