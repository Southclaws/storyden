import type { PatternRuntimeConfig } from '../types/pattern';
import type { TokenValue } from '../types/tokens';
import type { ConditionalValue, SystemProperties, SystemStyleObject } from '../types/system';

export interface CqProperties {
  name?: ConditionalValue<TokenValue<"containerNames"> | SystemProperties["containerName"]>
  type?: SystemProperties["containerType"]
  className?: string
}

type CqRestStyles = Omit<SystemStyleObject, keyof CqProperties>

interface CqStyles extends CqProperties, CqRestStyles {}

interface CqPatternFn {
  (styles?: CqStyles): string
  raw: (styles?: CqStyles) => SystemStyleObject
}

export declare function cqRaw(styles?: CqStyles): SystemStyleObject;

export declare const cq: CqPatternFn;