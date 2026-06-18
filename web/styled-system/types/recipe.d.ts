import type { ConditionalValue, SystemStyleObject } from './system';

export type RecipeVariantProps<T> = T extends (props?: infer Props) => unknown ? Props : never

export type RecipeVariant<T> = Required<NonNullable<RecipeVariantProps<T>>>

export type RecipeVariantMap<Variant extends object> = {
  [K in keyof Variant]-?: Array<Variant[K]>
}

export type RecipeConfigVariantMap<T> = {
  [K in keyof T]: Array<keyof T[K]>
}

export interface RecipeRuntimeFn<Props extends object = object, Map extends object = object> {
  (props?: Props): string
  __type: Props
  variantMap: Map
  variantKeys: Array<keyof Props>
  raw: (props?: Props) => SystemStyleObject
  splitVariantProps<T extends Record<string, any>>(props: T): [Props, Omit<T, keyof Props>]
  getVariantProps: (props?: Props) => Props
  merge(recipe: RecipeRuntimeFn): RecipeRuntimeFn
}

export type SlotRecord<Slot extends string, Value> = Partial<Record<Slot, Value>>

export interface SlotRecipeRuntimeFn<Slot extends string, Props extends object = object, Map extends object = object> {
  (props?: Props): SlotRecord<Slot, string>
  __type: Props
  __slot: Slot
  variantMap: Map
  variantKeys: Array<keyof Props>
  raw: (props?: Props) => Record<Slot, SystemStyleObject>
  splitVariantProps<T extends Record<string, any>>(props: T): [Props, Omit<T, keyof Props>]
  getVariantProps: (props?: Props) => Props
}

export type StringToBoolean<T> = T extends 'true' | 'false' ? boolean : T

export type RecipeVariantRecord = Record<string, Record<string, SystemStyleObject>>

export type RecipeSelection<T extends RecipeVariantRecord> = {
  [K in keyof T]?: StringToBoolean<keyof T[K]>
}

export type RecipeCompoundSelection<T> = {
  [K in keyof T]?: StringToBoolean<keyof T[K]> | Array<StringToBoolean<keyof T[K]>>
}

export interface RecipeDefinition<T extends RecipeVariantRecord = RecipeVariantRecord> {
  base?: SystemStyleObject
  variants?: T
  defaultVariants?: RecipeSelection<T>
  compoundVariants?: Array<RecipeCompoundSelection<T> & { css: SystemStyleObject }>
}

export interface RecipeCreatorFn {
  <T extends RecipeVariantRecord>(config: RecipeDefinition<T>): RecipeRuntimeFn<RecipeSelection<T>, RecipeConfigVariantMap<T>>
}

export type SlotRecipeVariantRecord<Slot extends string> = Record<string, Record<string, SlotRecord<Slot, SystemStyleObject>>>

export interface SlotRecipeDefinition<Slot extends string = string, T extends SlotRecipeVariantRecord<Slot> = SlotRecipeVariantRecord<Slot>> {
  className?: string
  slots: Slot[]
  base?: SlotRecord<Slot, SystemStyleObject>
  variants?: T
  defaultVariants?: RecipeSelection<T>
  compoundVariants?: Array<RecipeCompoundSelection<T> & { css: SlotRecord<Slot, SystemStyleObject> }>
}

export interface SlotRecipeCreatorFn {
  <Slot extends string, T extends SlotRecipeVariantRecord<Slot>>(config: SlotRecipeDefinition<Slot, T>): SlotRecipeRuntimeFn<Slot, RecipeSelection<T>, RecipeConfigVariantMap<T>>
}