/* eslint-disable */
import type { FunctionComponent } from 'react'
import type { LinkOverlayProperties } from '../patterns/link-overlay';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system-types';

export type LinkOverlayProps = LinkOverlayProperties & DistributiveOmit<HTMLStyledProps<'a'>, keyof LinkOverlayProperties >


export declare const LinkOverlay: FunctionComponent<LinkOverlayProps>