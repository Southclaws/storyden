import type { FunctionComponent } from 'react';
import type { LinkButtonProperties } from '../patterns/link-button';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system';

export interface LinkButtonProps extends LinkButtonProperties, DistributiveOmit<HTMLStyledProps<"div">, keyof LinkButtonProperties> {}

export declare const LinkButton: FunctionComponent<LinkButtonProps>