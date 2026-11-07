import type { DistributiveOmit, Pretty, SystemStyleObject } from './system';

export type RecipeVariantProps<T> = T extends (props?: infer Props) => unknown ? Props : never

export type RecipeVariant<T> = Required<NonNullable<RecipeVariantProps<T>>>

export type RecipeVariantMap<Variant extends object> = {
  [K in keyof Variant]-?: Array<Exclude<Variant[K], undefined>>
}

export type RecipeConfigVariantMap<T> = {
  [K in keyof T]: Array<keyof T[K]>
}

export interface RecipeRuntimeFn<Props extends object = object, Map extends object = object, Raw = SystemStyleObject> {
  (props?: Props): string
  __type: Props
  variantMap: Map
  variantKeys: Array<keyof Props>
  raw: (props?: Props) => Raw
  splitVariantProps<T extends Props & Record<string, any>>(props: T): [Props, Pretty<DistributiveOmit<T, keyof Props>>]
  getVariantProps: (props?: Props) => Props
  merge(recipe: RecipeRuntimeFn): RecipeRuntimeFn
}

export type ConfigRecipeRuntimeFn<Props extends object = object, Map extends object = object> = RecipeRuntimeFn<Props, Map, Props>

export type SlotRecord<Slot extends string, Value> = Partial<Record<Slot, Value>>

export interface SlotRecipeRuntimeFn<Slot extends string, Props extends object = object, Map extends object = object, Raw = Record<Slot, SystemStyleObject>> {
  (props?: Props): Record<Slot, string>
  __type: Props
  __slot: Slot
  variantMap: Map
  variantKeys: Array<keyof Props>
  raw: (props?: Props) => Raw
  splitVariantProps<T extends Props & Record<string, any>>(props: T): [Props, Pretty<DistributiveOmit<T, keyof Props>>]
  getVariantProps: (props?: Props) => Props
}

export type ConfigSlotRecipeRuntimeFn<Slot extends string, Props extends object = object, Map extends object = object> = SlotRecipeRuntimeFn<Slot, Props, Map, Props>

export type StringToBoolean<T> = T extends 'true' | 'false' ? boolean : T

export type RecipeVariantRecord = Record<string, Record<string, SystemStyleObject>>

export type RecipeSelection<T extends RecipeVariantRecord> = string extends keyof T
  ? {}
  : {
      [K in keyof T]?: StringToBoolean<keyof T[K]> | undefined
    }

export type RecipeCompoundSelection<T> = {
  [K in keyof T]?: StringToBoolean<keyof T[K]> | Array<StringToBoolean<keyof T[K]>> | undefined
}

export interface RecipeDefinition<T extends RecipeVariantRecord = RecipeVariantRecord> {
  base?: SystemStyleObject
  variants?: T
  defaultVariants?: RecipeSelection<T>
  compoundVariants?: Array<RecipeCompoundSelection<T> & { css: SystemStyleObject }>
}

export interface RecipeCreatorFn {
  <T extends RecipeVariantRecord>(config: RecipeDefinition<T>): CvaRuntimeFn<T>
}

export interface CvaRuntimeFn<T extends RecipeVariantRecord> extends RecipeRuntimeFn<RecipeSelection<T>, RecipeConfigVariantMap<T>> {
  config: RecipeDefinition<T>
}

export type SlotRecipeVariantRecord<Slot extends string> = Record<string, Record<string, SlotRecord<Slot, SystemStyleObject>>>

export interface SlotRecipeDefinition<Slot extends string = string, T extends SlotRecipeVariantRecord<Slot> = SlotRecipeVariantRecord<Slot>> {
  className?: string
  slots: readonly Slot[]
  base?: SlotRecord<Slot, SystemStyleObject>
  variants?: T
  defaultVariants?: RecipeSelection<T>
  compoundVariants?: Array<RecipeCompoundSelection<T> & { css: SlotRecord<Slot, SystemStyleObject> }>
}

export interface SlotRecipeCreatorFn {
  <Slot extends string, T extends SlotRecipeVariantRecord<Slot>>(config: SlotRecipeDefinition<Slot, T>): SvaRuntimeFn<Slot, T>
}

export interface SvaRuntimeFn<Slot extends string, T extends SlotRecipeVariantRecord<Slot>> extends SlotRecipeRuntimeFn<Slot, RecipeSelection<T>, RecipeConfigVariantMap<T>> {
  config: SlotRecipeDefinition<Slot, T>
}