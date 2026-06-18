import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type DatePickerVariant = {}

export type DatePickerVariantProps = {
  [K in keyof DatePickerVariant]?: ConditionalValue<DatePickerVariant[K]>
}

export type DatePickerVariantMap = RecipeVariantMap<DatePickerVariant>

export type DatePickerSlot = "clearTrigger" | "content" | "control" | "input" | "label" | "monthSelect" | "nextTrigger" | "positioner" | "presetTrigger" | "prevTrigger" | "rangeText" | "root" | "table" | "tableBody" | "tableCell" | "tableCellTrigger" | "tableHead" | "tableHeader" | "tableRow" | "trigger" | "view" | "viewControl" | "viewTrigger" | "yearSelect" | "view"

export type DatePickerRecipe = SlotRecipeRuntimeFn<DatePickerSlot, DatePickerVariantProps, DatePickerVariantMap>

export declare const datePicker: DatePickerRecipe;