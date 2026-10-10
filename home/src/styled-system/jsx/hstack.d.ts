import type { FunctionComponent } from 'react';
import type { HstackProperties } from '../patterns/hstack';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system';

export interface HstackProps extends HstackProperties, DistributiveOmit<HTMLStyledProps<"div">, keyof HstackProperties> {}

export declare const HStack: FunctionComponent<HstackProps>