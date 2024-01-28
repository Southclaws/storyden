/* eslint-disable */
import type { FunctionComponent } from 'react'
import type { LinkButtonProperties } from '../patterns/link-button';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system-types';

export interface LinkButtonProps extends LinkButtonProperties, DistributiveOmit<HTMLStyledProps<'div'>, keyof LinkButtonProperties > {}

/** Link button */
export declare const LinkButton: FunctionComponent<LinkButtonProps>