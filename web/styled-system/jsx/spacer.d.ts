import type { FunctionComponent } from 'react';
import type { SpacerProperties } from '../patterns/spacer';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system';

export interface SpacerProps extends SpacerProperties, DistributiveOmit<HTMLStyledProps<"div">, keyof SpacerProperties> {}

export declare const Spacer: FunctionComponent<SpacerProps>