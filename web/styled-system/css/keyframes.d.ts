import type { CssKeyframes } from '../types/system';

export type KeyframesFn = (keyframe: CssKeyframes[string]) => string;

export declare const keyframes: KeyframesFn;