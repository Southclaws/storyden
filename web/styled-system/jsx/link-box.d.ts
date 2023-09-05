/* eslint-disable */
import type { FunctionComponent } from 'react'
import type { LinkBoxProperties } from '../patterns/link-box';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system-types';

export type LinkBoxProps = LinkBoxProperties & DistributiveOmit<HTMLStyledProps<'div'>, keyof LinkBoxProperties >


export declare const LinkBox: FunctionComponent<LinkBoxProps>