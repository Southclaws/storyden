/* eslint-disable */
import type { FunctionComponent } from 'react'
import type { FrostedGlassProperties } from '../patterns/frosted-glass';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system-types';

export interface FrostedGlassProps extends FrostedGlassProperties, DistributiveOmit<HTMLStyledProps<'div'>, keyof FrostedGlassProperties > {}

/**
 * A frosted glass effect for overlays, modals, menus, etc. This is most prominently used on the navigation overlays and menus.
 */
export declare const FrostedGlass: FunctionComponent<FrostedGlassProps>