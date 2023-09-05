/* eslint-disable */
import type { FunctionComponent } from 'react'
import type { AspectRatioProperties } from '../patterns/aspect-ratio';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system-types';

export type AspectRatioProps = AspectRatioProperties & DistributiveOmit<HTMLStyledProps<'div'>, keyof AspectRatioProperties | 'aspectRatio'>


export declare const AspectRatio: FunctionComponent<AspectRatioProps>