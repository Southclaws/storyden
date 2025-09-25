/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { DistributiveOmit, Pretty } from '../types/system-types';

interface DatePickerVariant {
  
}

type DatePickerVariantMap = {
  [key in keyof DatePickerVariant]: Array<DatePickerVariant[key]>
}

type DatePickerSlot = "root" | "label" | "clearTrigger" | "content" | "control" | "input" | "monthSelect" | "nextTrigger" | "positioner" | "prevTrigger" | "rangeText" | "table" | "tableBody" | "tableCell" | "tableCellTrigger" | "tableHead" | "tableHeader" | "tableRow" | "trigger" | "viewTrigger" | "viewControl" | "yearSelect" | "presetTrigger" | "view"

export type DatePickerVariantProps = {
  [key in keyof DatePickerVariant]?: ConditionalValue<DatePickerVariant[key]> | undefined
}

export interface DatePickerRecipe {
  __slot: DatePickerSlot
  __type: DatePickerVariantProps
  (props?: DatePickerVariantProps): Pretty<Record<DatePickerSlot, string>>
  raw: (props?: DatePickerVariantProps) => DatePickerVariantProps
  variantMap: DatePickerVariantMap
  variantKeys: Array<keyof DatePickerVariant>
  splitVariantProps<Props extends DatePickerVariantProps>(props: Props): [DatePickerVariantProps, Pretty<DistributiveOmit<Props, keyof DatePickerVariantProps>>]
  getVariantProps: (props?: DatePickerVariantProps) => DatePickerVariantProps
}


export declare const datePicker: DatePickerRecipe