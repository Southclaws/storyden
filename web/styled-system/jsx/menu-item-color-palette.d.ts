/* eslint-disable */
import type { FunctionComponent } from 'react'
import type { MenuItemColorPaletteProperties } from '../patterns/menu-item-color-palette';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system-types';

export interface MenuItemColorPaletteProps extends MenuItemColorPaletteProperties, DistributiveOmit<HTMLStyledProps<'div'>, keyof MenuItemColorPaletteProperties > {}

/**
 * A color palette for menu items.
 */
export declare const MenuItemColorPalette: FunctionComponent<MenuItemColorPaletteProps>