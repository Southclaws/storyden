import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemStyleObject } from '../types/system';

export interface ContainerProperties {}

type ContainerRestStyles = Omit<SystemStyleObject, keyof ContainerProperties>

interface ContainerStyles extends ContainerProperties, ContainerRestStyles {}

interface ContainerPatternFn {
  (styles?: ContainerStyles): string
  raw: (styles?: ContainerStyles) => SystemStyleObject
  propKeys: Array<keyof ContainerProperties>
}

export declare function containerRaw(styles?: ContainerStyles): SystemStyleObject;

export declare const container: ContainerPatternFn;