import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemStyleObject } from '../types/system';

export interface MenuItemColorPaletteProperties {
  className?: string
}

type MenuItemColorPaletteRestStyles = Omit<SystemStyleObject, keyof MenuItemColorPaletteProperties>

interface MenuItemColorPaletteStyles extends MenuItemColorPaletteProperties, MenuItemColorPaletteRestStyles {}

interface MenuItemColorPalettePatternFn {
  (styles?: MenuItemColorPaletteStyles): string
  raw: (styles?: MenuItemColorPaletteStyles) => SystemStyleObject
}

export declare function menuItemColorPaletteRaw(styles?: MenuItemColorPaletteStyles): SystemStyleObject;

export declare const menuItemColorPalette: MenuItemColorPalettePatternFn;