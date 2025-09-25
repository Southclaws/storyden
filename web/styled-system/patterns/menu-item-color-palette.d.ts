/* eslint-disable */
import type { SystemStyleObject, ConditionalValue } from '../types/index';
import type { Properties } from '../types/csstype';
import type { SystemProperties } from '../types/style-props';
import type { DistributiveOmit } from '../types/system-types';
import type { Tokens } from '../tokens/index';

export interface MenuItemColorPaletteProperties {
   
}

interface MenuItemColorPaletteStyles extends MenuItemColorPaletteProperties, DistributiveOmit<SystemStyleObject, keyof MenuItemColorPaletteProperties > {}

interface MenuItemColorPalettePatternFn {
  (styles?: MenuItemColorPaletteStyles): string
  raw: (styles?: MenuItemColorPaletteStyles) => SystemStyleObject
}

/**
 * A color palette for menu items.
 */
export declare const menuItemColorPalette: MenuItemColorPalettePatternFn;
