import type { FunctionComponent } from 'react';
import type { MenuItemColorPaletteProperties } from '../patterns/menu-item-color-palette';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system';

export interface MenuItemColorPaletteProps extends MenuItemColorPaletteProperties, DistributiveOmit<HTMLStyledProps<"div">, keyof MenuItemColorPaletteProperties> {}

export declare const MenuItemColorPalette: FunctionComponent<MenuItemColorPaletteProps>