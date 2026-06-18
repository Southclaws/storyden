import type { PatternRuntimeConfig } from '../types/pattern';
import type { TokenValue } from '../types/tokens';
import type { ConditionalValue, SystemProperties, SystemStyleObject } from '../types/system';

export interface DividerProperties {
  color?: ConditionalValue<TokenValue<"colors"> | SystemProperties["borderColor"]>
  orientation?: ConditionalValue<"horizontal" | "vertical">
  thickness?: ConditionalValue<TokenValue<"sizes"> | SystemProperties["borderWidth"]>
  className?: string
}

type DividerRestStyles = Omit<SystemStyleObject, keyof DividerProperties>

interface DividerStyles extends DividerProperties, DividerRestStyles {}

interface DividerPatternFn {
  (styles?: DividerStyles): string
  raw: (styles?: DividerStyles) => SystemStyleObject
}

export declare function dividerRaw(styles?: DividerStyles): SystemStyleObject;

export declare const divider: DividerPatternFn;