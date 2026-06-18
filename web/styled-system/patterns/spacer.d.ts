import type { PatternRuntimeConfig } from '../types/pattern';
import type { TokenValue } from '../types/tokens';
import type { ConditionalValue, SystemStyleObject } from '../types/system';

export interface SpacerProperties {
  size?: ConditionalValue<TokenValue<"spacing">>
  className?: string
}

type SpacerRestStyles = Omit<SystemStyleObject, keyof SpacerProperties>

interface SpacerStyles extends SpacerProperties, SpacerRestStyles {}

interface SpacerPatternFn {
  (styles?: SpacerStyles): string
  raw: (styles?: SpacerStyles) => SystemStyleObject
}

export declare function spacerRaw(styles?: SpacerStyles): SystemStyleObject;

export declare const spacer: SpacerPatternFn;