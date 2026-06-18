import type { FunctionComponent } from 'react';
import type { FlexProperties } from '../patterns/flex';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system';

export interface FlexProps extends FlexProperties, DistributiveOmit<HTMLStyledProps<"div">, keyof FlexProperties> {}

export declare const Flex: FunctionComponent<FlexProps>