import type { TokenValue } from './tokens';

export type Pretty<T> = { [K in keyof T]: T[K] } & {}

export type DistributiveOmit<T, K extends keyof any> = T extends unknown ? Omit<T, K> : never

export type DistributiveUnion<T, U> = {
  [K in keyof T]: K extends keyof U ? U[K] | T[K] : T[K]
} & DistributiveOmit<U, keyof T>

export type Assign<T, U> = {
  [K in keyof T]: K extends keyof U ? U[K] : T[K]
} & U

export interface Conditions {
  "2xl": string
  "2xlDown": string
  "2xlOnly": string
  "_active": string
  "_after": string
  "_anyPointerCoarse": string
  "_anyPointerFine": string
  "_anyPointerNone": string
  "_atValue": string
  "_autofill": string
  "_backdrop": string
  "_before": string
  "_checked": string
  "_closed": string
  "_collapsed": string
  "_complete": string
  "_containerLarge": string
  "_containerMedium": string
  "_containerSmall": string
  "_current": string
  "_currentPage": string
  "_currentStep": string
  "_dark": string
  "_default": string
  "_detailsOpen": string
  "_disabled": string
  "_dragging": string
  "_empty": string
  "_enabled": string
  "_even": string
  "_expanded": string
  "_file": string
  "_first": string
  "_firstLetter": string
  "_firstLine": string
  "_firstOfType": string
  "_focus": string
  "_focusVisible": string
  "_focusWithin": string
  "_fullscreen": string
  "_grabbed": string
  "_groupActive": string
  "_groupChecked": string
  "_groupDisabled": string
  "_groupExpanded": string
  "_groupFocus": string
  "_groupFocusVisible": string
  "_groupFocusWithin": string
  "_groupHover": string
  "_groupInvalid": string
  "_hidden": string
  "_highContrast": string
  "_highlighted": string
  "_horizontal": string
  "_hover": string
  "_icon": string
  "_inRange": string
  "_incomplete": string
  "_indeterminate": string
  "_inert": string
  "_invalid": string
  "_invertedColors": string
  "_landscape": string
  "_last": string
  "_lastOfType": string
  "_lessContrast": string
  "_light": string
  "_loading": string
  "_ltr": string
  "_marker": string
  "_moreContrast": string
  "_motionReduce": string
  "_motionSafe": string
  "_noscript": string
  "_now": string
  "_odd": string
  "_off": string
  "_on": string
  "_only": string
  "_onlyOfType": string
  "_open": string
  "_optional": string
  "_osDark": string
  "_osLight": string
  "_outOfRange": string
  "_overValue": string
  "_peerActive": string
  "_peerChecked": string
  "_peerDisabled": string
  "_peerExpanded": string
  "_peerFocus": string
  "_peerFocusVisible": string
  "_peerFocusWithin": string
  "_peerHover": string
  "_peerInvalid": string
  "_peerPlaceholderShown": string
  "_placeholder": string
  "_placeholderShown": string
  "_pointerCoarse": string
  "_pointerFine": string
  "_pointerNone": string
  "_portrait": string
  "_pressed": string
  "_print": string
  "_rangeEnd": string
  "_rangeStart": string
  "_readOnly": string
  "_readWrite": string
  "_required": string
  "_rtl": string
  "_scrollbar": string
  "_scrollbarThumb": string
  "_scrollbarTrack": string
  "_selected": string
  "_selection": string
  "_starting": string
  "_target": string
  "_today": string
  "_topmost": string
  "_unavailable": string
  "_underValue": string
  "_userInvalid": string
  "_userValid": string
  "_valid": string
  "_vertical": string
  "_visited": string
  "base": string
  "lg": string
  "lgDown": string
  "lgOnly": string
  "lgTo2xl": string
  "lgToXl": string
  "md": string
  "mdDown": string
  "mdOnly": string
  "mdTo2xl": string
  "mdToLg": string
  "mdToXl": string
  "sm": string
  "smDown": string
  "smOnly": string
  "smTo2xl": string
  "smToLg": string
  "smToMd": string
  "smToXl": string
  "xl": string
  "xlDown": string
  "xlOnly": string
  "xlTo2xl": string
}

export interface Breakpoints {
  "base": string
  "sm": string
  "md": string
  "lg": string
  "xl": string
  "2xl": string
}

export type ContainerName = AnyString
export type CssContainer = ContainerName | `${ContainerName} / inline-size` | `${ContainerName} / size` | AnyString

export type Condition = keyof Conditions

export type ConditionalValue<T> =
  | T
  | Array<T | null>
  | { [K in Condition]?: ConditionalValue<T> }
  | undefined

export type AnyString = string & {}

export type AnyNumber = number & {}

export type CssVars = `var(--${string})`

type WithColorOpacityModifier<T> = [T] extends [string] ? `${T}/${string}` & { __colorOpacityModifier?: true } : never

type ImportantMark = "!" | "!important"
type WhitespaceImportant = ` ${ImportantMark}`
type Important = ImportantMark | WhitespaceImportant
type WithImportant<T> = [T] extends [string] ? `${T}${Important}` & { __important?: true } : never

export type WithEscapeHatch<T> = T | `[${string}]` | WithColorOpacityModifier<T> | WithImportant<T>

export type OnlyKnown<Value> = Value extends boolean ? Value : Value extends `${infer _}` ? Value : never

export type AlignContentValue = WithEscapeHatch<CssGlobals | CssVars | OnlyKnown<"baseline" | "center" | "end" | "first baseline" | "flex-end" | "flex-start" | "last baseline" | "normal" | "safe center" | "safe end" | "safe start" | "space-around" | "space-between" | "space-evenly" | "start" | "stretch">>

export type AlignItemsValue = WithEscapeHatch<CssGlobals | CssVars | OnlyKnown<"anchor-center" | "baseline" | "center" | "end" | "first baseline" | "flex-end" | "flex-start" | "last baseline" | "normal" | "safe center" | "safe end" | "safe start" | "self-end" | "self-start" | "start" | "stretch">>

export type AlignSelfValue = WithEscapeHatch<CssGlobals | CssVars | OnlyKnown<"anchor-center" | "auto" | "baseline" | "center" | "end" | "first baseline" | "flex-end" | "flex-start" | "last baseline" | "normal" | "safe center" | "safe end" | "safe start" | "self-end" | "self-start" | "start" | "stretch">>

export type AnimationCompositionValue = string | number | CssVars | AnyString

export type AnimationDirectionValue = string | number | CssVars | AnyString

export type AnimationFillModeValue = string | number | CssVars | AnyString

export type AnimationIterationCountValue = string | number | CssVars | AnyString

export type AnimationPlayStateValue = string | number | CssVars | AnyString

export type AnimationRangeEndValue = string | number | CssVars | AnyString

export type AnimationRangeStartValue = string | number | CssVars | AnyString

export type AnimationRangeValue = string | number | CssVars | AnyString

export type AnimationStateValue = string | number | CssVars | AnyString

export type AnimationTimelineValue = string | number | CssVars | AnyString

export type AnimationsValue = WithEscapeHatch<CssGlobals | TokenValue<"animations"> | CssVars>

export type AppearanceValue = WithEscapeHatch<PropertyValueMap["appearance"]>

export type AspectRatiosValue = string | number | CssVars | AnyString

export type AssetsValue = PropertyValueMap["listStyleImage"]

export type BackdropBrightnessValue = string | number | CssVars | AnyString

export type BackdropContrastValue = string | number | CssVars | AnyString

export type BackdropFilterValue = WithEscapeHatch<CssGlobals | "auto" | CssVars>

export type BackdropGrayscaleValue = string | number | CssVars | AnyString

export type BackdropHueRotateValue = string | number | CssVars | AnyString

export type BackdropInvertValue = string | number | CssVars | AnyString

export type BackdropOpacityValue = string | number | CssVars | AnyString

export type BackdropSaturateValue = string | number | CssVars | AnyString

export type BackdropSepiaValue = string | number | CssVars | AnyString

export type BackfaceVisibilityValue = string | number | CssVars | AnyString

export type BackgroundAttachmentValue = string | number | CssVars | AnyString

export type BackgroundBlendModeValue = string | number | CssVars | AnyString

export type BackgroundClipValue = string | number | CssVars | AnyString

export type BackgroundConicValue = string | number | CssVars | AnyString

export type BackgroundGradientValue = WithEscapeHatch<CssGlobals | "to-b" | "to-bl" | "to-br" | "to-l" | "to-r" | "to-t" | "to-tl" | "to-tr" | CssVars>

export type BackgroundLinearValue = WithEscapeHatch<CssGlobals | "to-b" | "to-bl" | "to-br" | "to-l" | "to-r" | "to-t" | "to-tl" | "to-tr" | CssVars>

export type BackgroundOriginValue = string | number | CssVars | AnyString

export type BackgroundPositionValue = PropertyValueMap["backgroundPosition"]

export type BackgroundPositionXValue = string | number | CssVars | AnyString

export type BackgroundPositionYValue = string | number | CssVars | AnyString

export type BackgroundRepeatValue = PropertyValueMap["backgroundRepeat"]

export type BackgroundSizeValue = PropertyValueMap["backgroundSize"]

export type BlockSizeValue = WithEscapeHatch<CssDimensionGlobals | TokenValue<"sizes"> | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "2/3" | "2/4" | "2/5" | "2/6" | "3/4" | "3/5" | "3/6" | "4/5" | "4/6" | "5/6" | "auto" | "dvh" | "lvh" | "screen" | "svh" | CssVars>

export type BlursValue = WithEscapeHatch<CssGlobals | TokenValue<"blurs"> | CssVars>

export type BorderCollapseValue = string | number | CssVars | AnyString

export type BorderStyleValue = string | number | CssVars | AnyString

export type BorderStylesValue = WithEscapeHatch<CssGlobals | TokenValue<"borderStyles"> | CssVars>

export type BorderWidthsValue = WithEscapeHatch<CssGlobals | TokenValue<"borderWidths"> | CssVars>

export type BordersValue = WithEscapeHatch<CssGlobals | TokenValue<"borders"> | CssVars>

export type BoxDecorationBreakValue = string | number | CssVars | AnyString

export type BoxSizeValue = WithEscapeHatch<CssDimensionGlobals | TokenValue<"sizes"> | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10/12" | "11/12" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "3/12" | "3/4" | "3/5" | "3/6" | "4/12" | "4/5" | "4/6" | "5/12" | "5/6" | "6/12" | "7/12" | "8/12" | "9/12" | "auto" | "screen" | CssVars>

export type BoxSizingValue = string | number | CssVars | AnyString

export type BreakpointsValue = WithEscapeHatch<CssGlobals | TokenValue<"breakpoints"> | CssVars>

export type BrightnessValue = string | number | CssVars | AnyString

export type ClipPathValue = PropertyValueMap["clipPath"]

export type ColorPaletteValue = WithEscapeHatch<CssGlobals | "accent" | "accent.dark" | "accent.light" | "amber" | "amber.dark" | "amber.light" | "background" | "backgroundGradientH" | "backgroundGradientV" | "black" | "blue" | "blue.dark" | "blue.light" | "border" | "cardBackgroundGradient" | "conicGradient" | "current" | "gray" | "gray.dark" | "gray.light" | "green" | "green.dark" | "green.light" | "interactive" | "interactive.emphasized" | "interactive.selected" | "neutral" | "neutral.dark" | "neutral.light" | "orange" | "orange.dark" | "orange.light" | "overflow-fade" | "pink" | "pink.dark" | "pink.light" | "red" | "red.dark" | "red.light" | "scrim" | "scroll-fade-top" | "selection" | "slate" | "slate.dark" | "slate.light" | "status" | "status.danger" | "status.info" | "status.success" | "status.warning" | "text" | "tomato" | "tomato.dark" | "tomato.light" | "transparent" | "visibility" | "visibility.draft" | "visibility.draft.border" | "visibility.draft.content" | "visibility.draft.surface" | "visibility.published" | "visibility.published.border" | "visibility.published.content" | "visibility.published.surface" | "visibility.review" | "visibility.review.border" | "visibility.review.content" | "visibility.review.surface" | "visibility.unlisted" | "visibility.unlisted.border" | "visibility.unlisted.content" | "visibility.unlisted.surface" | "white" | CssVars>

export type ColorsValue = WithEscapeHatch<CssColorGlobals | TokenValue<"colors"> | CssVars>

export type ContainerNamesValue = WithEscapeHatch<CssGlobals | TokenValue<"containerNames"> | CssVars>

export type ContainerTypeValue = PropertyValueMap["containerType"]

export type ContainerValue = string | number | CssVars | AnyString

export type ContrastValue = string | number | CssVars | AnyString

export type CursorValue = PropertyValueMap["cursor"]

export type DebugValue = WithEscapeHatch<CssGlobals | boolean | CssVars>

export type DisplayValue = WithEscapeHatch<CssGlobals | CssVars | OnlyKnown<"-ms-flexbox" | "-ms-grid" | "-ms-inline-flexbox" | "-ms-inline-grid" | "-webkit-flex" | "-webkit-inline-flex" | "block" | "contents" | "flex" | "flow" | "flow-root" | "grid" | "inline" | "inline-block" | "inline-flex" | "inline-grid" | "inline-list-item" | "inline-table" | "list-item" | "none" | "ruby" | "ruby-base" | "ruby-base-container" | "ruby-text" | "ruby-text-container" | "run-in" | "table" | "table-caption" | "table-cell" | "table-column" | "table-column-group" | "table-footer-group" | "table-header-group" | "table-row" | "table-row-group">>

export type DropShadowsValue = WithEscapeHatch<CssGlobals | TokenValue<"dropShadows"> | CssVars>

export type DurationsValue = WithEscapeHatch<CssGlobals | TokenValue<"durations"> | CssVars>

export type EasingsValue = WithEscapeHatch<CssGlobals | TokenValue<"easings"> | CssVars>

export type FilterValue = WithEscapeHatch<CssGlobals | "auto" | CssVars>

export type FlexBasisValue = WithEscapeHatch<CssGlobals | "0" | "0.5" | "1" | "1.5" | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10" | "10.5" | "10/12" | "11" | "11/12" | "12" | "14" | "16" | "2" | "2.5" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "20" | "24" | "28" | "2xl" | "2xs" | "3" | "3.5" | "3/12" | "3/4" | "3/5" | "3/6" | "32" | "36" | "3xl" | "4" | "4.5" | "4/12" | "4/5" | "4/6" | "40" | "44" | "48" | "4xl" | "5" | "5.5" | "5/12" | "5/6" | "52" | "56" | "5xl" | "6" | "6/12" | "60" | "64" | "6xl" | "7" | "7.5" | "7/12" | "72" | "7xl" | "8" | "8/12" | "80" | "8xl" | "9" | "9.5" | "9/12" | "96" | "breakpoint-2xl" | "breakpoint-lg" | "breakpoint-md" | "breakpoint-sm" | "breakpoint-xl" | "fit" | "full" | "layout.commandBar" | "layout.content" | "layout.contentWide" | "layout.drawer" | "layout.form" | "layout.readable" | "layout.sidebar" | "lg" | "max" | "md" | "min" | "prose" | "safeBottom" | "scrollGutter" | "sm" | "viewportHeight" | "xl" | "xs" | CssVars>

export type FlexDirectionValue = string | number | CssVars | AnyString

export type FlexGrowValue = PropertyValueMap["flexGrow"]

export type FlexShrinkValue = PropertyValueMap["flexShrink"]

export type FlexValue = WithEscapeHatch<CssGlobals | "1" | "auto" | "initial" | "none" | CssVars>

export type FloatValue = WithEscapeHatch<CssGlobals | CssVars | OnlyKnown<"end" | "start">>

export type FocusRingValue = WithEscapeHatch<CssGlobals | "inside" | "mixed" | "none" | "outside" | CssVars>

export type FocusVisibleRingValue = WithEscapeHatch<CssGlobals | "inside" | "mixed" | "none" | "outside" | CssVars>

export type FontFeatureSettingsValue = string | number | CssVars | AnyString

export type FontKerningValue = string | number | CssVars | AnyString

export type FontPaletteValue = string | number | CssVars | AnyString

export type FontSizeAdjustValue = PropertyValueMap["fontSizeAdjust"]

export type FontSizesValue = WithEscapeHatch<CssGlobals | TokenValue<"fontSizes"> | CssVars>

export type FontSmoothingValue = WithEscapeHatch<CssGlobals | "antialiased" | "subpixel-antialiased" | CssVars>

export type FontVariantAlternatesValue = string | number | CssVars | AnyString

export type FontVariantCapsValue = PropertyValueMap["fontVariantCaps"] | AnyString

export type FontVariantNumericValue = PropertyValueMap["fontVariantNumeric"]

export type FontVariantValue = string | number | CssVars | AnyString

export type FontVariationSettingsValue = string | number | CssVars | AnyString

export type FontWeightsValue = WithEscapeHatch<CssGlobals | TokenValue<"fontWeights"> | CssVars>

export type FontsValue = WithEscapeHatch<CssGlobals | TokenValue<"fonts"> | CssVars>

export type GradientFromPositionValue = string | number | CssVars | AnyString

export type GradientToPositionValue = string | number | CssVars | AnyString

export type GradientViaPositionValue = string | number | CssVars | AnyString

export type GradientsValue = string | number | CssVars | AnyString

export type GrayscaleValue = string | number | CssVars | AnyString

export type GridAutoColumnsValue = WithEscapeHatch<CssGlobals | "fr" | "max" | "min" | CssVars>

export type GridAutoFlowValue = PropertyValueMap["gridAutoFlow"]

export type GridAutoRowsValue = WithEscapeHatch<CssGlobals | "fr" | "max" | "min" | CssVars>

export type GridColumnEndValue = string | number | CssVars | AnyString

export type GridColumnStartValue = string | number | CssVars | AnyString

export type GridColumnValue = string | number | CssVars | AnyString

export type GridRowValue = string | number | CssVars | AnyString

export type GridTemplateColumnsValue = string | number | CssVars | AnyString

export type GridTemplateRowsValue = string | number | CssVars | AnyString

export type HeightValue = WithEscapeHatch<CssDimensionGlobals | TokenValue<"sizes"> | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "2/3" | "2/4" | "2/5" | "2/6" | "3/4" | "3/5" | "3/6" | "4/5" | "4/6" | "5/6" | "auto" | "dvh" | "lvh" | "screen" | "svh" | CssVars>

export type HueRotateValue = string | number | CssVars | AnyString

export type HyphensValue = PropertyValueMap["hyphens"] | AnyString

export type InlineSizeValue = WithEscapeHatch<CssDimensionGlobals | TokenValue<"sizes"> | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10/12" | "11/12" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "3/12" | "3/4" | "3/5" | "3/6" | "4/12" | "4/5" | "4/6" | "5/12" | "5/6" | "6/12" | "7/12" | "8/12" | "9/12" | "auto" | "screen" | CssVars>

export type InvertValue = string | number | CssVars | AnyString

export type JustifyContentValue = PropertyValueMap["justifyContent"]

export type KeyframesValue = WithEscapeHatch<CssGlobals | "fadeIn" | "fadeOut" | "shimmer" | "targetPulse" | CssVars>

export type LetterSpacingsValue = WithEscapeHatch<CssGlobals | TokenValue<"letterSpacings"> | CssVars>

export type LineClampValue = string | number | CssVars | AnyString

export type LineHeightsValue = WithEscapeHatch<CssGlobals | TokenValue<"lineHeights"> | CssVars>

export type ListStylePositionValue = PropertyValueMap["listStylePosition"] | AnyString

export type ListStyleTypeValue = PropertyValueMap["listStyleType"]

export type ListStyleValue = string | number | CssVars | AnyString

export type MaskClipValue = string | number | CssVars | AnyString

export type MaskCompositeValue = string | number | CssVars | AnyString

export type MaskConicValue = string | number | CssVars | AnyString

export type MaskImageValue = string | number | CssVars | AnyString

export type MaskLinearValue = WithEscapeHatch<CssGlobals | "to-b" | "to-bl" | "to-br" | "to-l" | "to-r" | "to-t" | "to-tl" | "to-tr" | CssVars>

export type MaskModeValue = string | number | CssVars | AnyString

export type MaskOriginValue = string | number | CssVars | AnyString

export type MaskPositionValue = PropertyValueMap["maskPosition"]

export type MaskRadialAtValue = WithEscapeHatch<CssGlobals | "bottom" | "bottom left" | "bottom right" | "center" | "left" | "right" | "top" | "top left" | "top right" | CssVars>

export type MaskRadialShapeValue = WithEscapeHatch<CssGlobals | "circle" | "ellipse" | CssVars>

export type MaskRadialSizeValue = WithEscapeHatch<CssGlobals | "closest-corner" | "closest-side" | "farthest-corner" | "farthest-side" | CssVars>

export type MaskRadialValue = string | number | CssVars | AnyString

export type MaskRepeatValue = PropertyValueMap["maskRepeat"]

export type MaskSizeValue = PropertyValueMap["maskSize"]

export type MaskTypeValue = string | number | CssVars | AnyString

export type MaskValue = string | number | CssVars | AnyString

export type MaxBlockSizeValue = WithEscapeHatch<CssDimensionGlobals | TokenValue<"sizes"> | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "2/3" | "2/4" | "2/5" | "2/6" | "3/4" | "3/5" | "3/6" | "4/5" | "4/6" | "5/6" | "auto" | "dvh" | "lvh" | "screen" | "svh" | CssVars>

export type MaxHeightValue = WithEscapeHatch<CssDimensionGlobals | TokenValue<"sizes"> | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "2/3" | "2/4" | "2/5" | "2/6" | "3/4" | "3/5" | "3/6" | "4/5" | "4/6" | "5/6" | "auto" | "dvh" | "lvh" | "screen" | "svh" | CssVars>

export type MaxInlineSizeValue = WithEscapeHatch<CssDimensionGlobals | TokenValue<"sizes"> | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10/12" | "11/12" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "3/12" | "3/4" | "3/5" | "3/6" | "4/12" | "4/5" | "4/6" | "5/12" | "5/6" | "6/12" | "7/12" | "8/12" | "9/12" | "auto" | "screen" | CssVars>

export type MaxWidthValue = WithEscapeHatch<CssDimensionGlobals | TokenValue<"sizes"> | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10/12" | "11/12" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "3/12" | "3/4" | "3/5" | "3/6" | "4/12" | "4/5" | "4/6" | "5/12" | "5/6" | "6/12" | "7/12" | "8/12" | "9/12" | "auto" | "screen" | CssVars>

export type MinBlockSizeValue = WithEscapeHatch<CssDimensionGlobals | TokenValue<"sizes"> | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "2/3" | "2/4" | "2/5" | "2/6" | "3/4" | "3/5" | "3/6" | "4/5" | "4/6" | "5/6" | "auto" | "dvh" | "lvh" | "screen" | "svh" | CssVars>

export type MinHeightValue = WithEscapeHatch<CssDimensionGlobals | TokenValue<"sizes"> | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "2/3" | "2/4" | "2/5" | "2/6" | "3/4" | "3/5" | "3/6" | "4/5" | "4/6" | "5/6" | "auto" | "dvh" | "lvh" | "screen" | "svh" | CssVars>

export type MinInlineSizeValue = WithEscapeHatch<CssDimensionGlobals | TokenValue<"sizes"> | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10/12" | "11/12" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "3/12" | "3/4" | "3/5" | "3/6" | "4/12" | "4/5" | "4/6" | "5/12" | "5/6" | "6/12" | "7/12" | "8/12" | "9/12" | "auto" | "screen" | CssVars>

export type MinWidthValue = WithEscapeHatch<CssDimensionGlobals | TokenValue<"sizes"> | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10/12" | "11/12" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "3/12" | "3/4" | "3/5" | "3/6" | "4/12" | "4/5" | "4/6" | "5/12" | "5/6" | "6/12" | "7/12" | "8/12" | "9/12" | "auto" | "screen" | CssVars>

export type MixBlendModeValue = WithEscapeHatch<PropertyValueMap["mixBlendMode"]>

export type ObjectFitValue = WithEscapeHatch<PropertyValueMap["objectFit"]>

export type ObjectPositionValue = PropertyValueMap["objectPosition"]

export type OpacityValue = WithEscapeHatch<CssGlobals | TokenValue<"opacity"> | CssVars>

export type OverflowAnchorValue = string | number | CssVars | AnyString

export type OverflowBlockValue = PropertyValueMap["overflowBlock"] | AnyString

export type OverflowClipBoxValue = string | number | CssVars | AnyString

export type OverflowClipMarginValue = PropertyValueMap["overflowClipMargin"]

export type OverflowInlineValue = PropertyValueMap["overflowInline"] | AnyString

export type OverflowValue = WithEscapeHatch<PropertyValueMap["overflow"]>

export type OverflowWrapValue = string | number | CssVars | AnyString

export type OverflowXValue = WithEscapeHatch<PropertyValueMap["overflowX"]>

export type OverflowYValue = PropertyValueMap["overflowY"] | AnyString

export type OverscrollBehaviorBlockValue = string | number | CssVars | AnyString

export type OverscrollBehaviorInlineValue = string | number | CssVars | AnyString

export type OverscrollBehaviorValue = PropertyValueMap["overscrollBehavior"]

export type OverscrollBehaviorXValue = PropertyValueMap["overscrollBehaviorX"] | AnyString

export type OverscrollBehaviorYValue = PropertyValueMap["overscrollBehaviorY"] | AnyString

export type PositionValue = WithEscapeHatch<PropertyValueMap["position"]>

export type RadiiValue = WithEscapeHatch<CssGlobals | TokenValue<"radii"> | CssVars>

export type RotateValue = WithEscapeHatch<CssGlobals | "auto" | "auto-3d" | CssVars>

export type SaturateValue = string | number | CssVars | AnyString

export type ScaleValue = WithEscapeHatch<CssGlobals | "auto" | CssVars>

export type ScaleXValue = string | number | CssVars | AnyString

export type ScaleYValue = string | number | CssVars | AnyString

export type ScrollBehaviorValue = string | number | CssVars | AnyString

export type ScrollSnapAlignValue = PropertyValueMap["scrollSnapAlign"]

export type ScrollSnapCoordinateValue = string | number | CssVars | AnyString

export type ScrollSnapDestinationValue = PropertyValueMap["scrollSnapDestination"]

export type ScrollSnapPointsXValue = string | number | CssVars | AnyString

export type ScrollSnapPointsYValue = string | number | CssVars | AnyString

export type ScrollSnapStopValue = PropertyValueMap["scrollSnapStop"] | AnyString

export type ScrollSnapStrictnessValue = WithEscapeHatch<CssGlobals | "mandatory" | "proximity" | CssVars>

export type ScrollSnapTypeValue = WithEscapeHatch<CssGlobals | "both" | "none" | "x" | "y" | CssVars>

export type ScrollSnapTypeXValue = string | number | CssVars | AnyString

export type ScrollSnapTypeYValue = string | number | CssVars | AnyString

export type ScrollTimelineAxisValue = string | number | CssVars | AnyString

export type ScrollTimelineNameValue = string | number | CssVars | AnyString

export type ScrollTimelineValue = string | number | CssVars | AnyString

export type ScrollbarColorValue = string | number | CssVars | AnyString

export type ScrollbarGutterValue = PropertyValueMap["scrollbarGutter"]

export type ScrollbarValue = WithEscapeHatch<CssGlobals | "hidden" | "visible" | CssVars>

export type ScrollbarWidthValue = WithEscapeHatch<PropertyValueMap["scrollbarWidth"]>

export type SepiaValue = string | number | CssVars | AnyString

export type ShadowsValue = WithEscapeHatch<CssGlobals | TokenValue<"shadows"> | CssVars>

export type SpacingValue = WithEscapeHatch<CssAutoGlobals | TokenValue<"spacing"> | "auto" | CssVars>

export type SrOnlyValue = WithEscapeHatch<CssGlobals | boolean | CssVars>

export type StrokeDasharrayValue = string | number | CssVars | AnyString

export type StrokeDashoffsetValue = string | number | CssVars | AnyString

export type StrokeLinecapValue = string | number | CssVars | AnyString

export type StrokeLinejoinValue = PropertyValueMap["strokeLinejoin"] | AnyString

export type StrokeMiterlimitValue = string | number | CssVars | AnyString

export type StrokeOpacityValue = string | number | CssVars | AnyString

export type TableLayoutValue = PropertyValueMap["tableLayout"] | AnyString

export type TextAlignValue = PropertyValueMap["textAlign"] | AnyString

export type TextDecorationStyleValue = PropertyValueMap["textDecorationStyle"] | AnyString

export type TextDecorationThicknessValue = PropertyValueMap["textDecorationThickness"]

export type TextDecorationValue = string | number | CssVars | AnyString

export type TextGradientValue = WithEscapeHatch<CssGlobals | "to-b" | "to-bl" | "to-br" | "to-l" | "to-r" | "to-t" | "to-tl" | "to-tr" | CssVars>

export type TextOverflowValue = PropertyValueMap["textOverflow"]

export type TextSizeAdjustValue = PropertyValueMap["textSizeAdjust"]

export type TextStyleValue = WithEscapeHatch<CssGlobals | "body" | "metadata" | "supporting" | CssVars>

export type TextTransformValue = PropertyValueMap["textTransform"]

export type TextUnderlineOffsetValue = PropertyValueMap["textUnderlineOffset"]

export type TextWrapValue = PropertyValueMap["textWrap"]

export type TouchActionValue = WithEscapeHatch<CssGlobals | CssVars | OnlyKnown<"-ms-manipulation" | "-ms-none" | "-ms-pan-x" | "-ms-pan-y" | "-ms-pinch-zoom" | "auto" | "manipulation" | "none" | "pan-down" | "pan-left" | "pan-right" | "pan-up" | "pan-x" | "pan-y" | "pinch-zoom">>

export type TransformBoxValue = WithEscapeHatch<PropertyValueMap["transformBox"]>

export type TransformOriginValue = string | number | CssVars | AnyString

export type TransformStyleValue = string | number | CssVars | AnyString

export type TransformValue = string | number | CssVars | AnyString

export type TransitionPropertyValue = WithEscapeHatch<CssGlobals | "background" | "colors" | "common" | "position" | "size" | CssVars>

export type TransitionValue = WithEscapeHatch<CssGlobals | "all" | "background" | "colors" | "common" | "opacity" | "position" | "shadow" | "size" | "transform" | CssVars>

export type TranslateValue = WithEscapeHatch<CssGlobals | "auto" | "auto-3d" | CssVars>

export type TranslateXValue = WithEscapeHatch<CssAutoGlobals | TokenValue<"spacing"> | "-1/2" | "-1/3" | "-1/4" | "-2/3" | "-2/4" | "-3/4" | "-full" | "1/2" | "1/3" | "1/4" | "2/3" | "2/4" | "3/4" | "full" | CssVars>

export type TranslateYValue = WithEscapeHatch<CssAutoGlobals | TokenValue<"spacing"> | "-1/2" | "-1/3" | "-1/4" | "-2/3" | "-2/4" | "-3/4" | "-full" | "1/2" | "1/3" | "1/4" | "2/3" | "2/4" | "3/4" | "full" | CssVars>

export type TruncateValue = WithEscapeHatch<CssGlobals | boolean | CssVars>

export type UserSelectValue = WithEscapeHatch<PropertyValueMap["userSelect"]>

export type VerticalAlignValue = PropertyValueMap["verticalAlign"]

export type VisibilityValue = string | number | CssVars | AnyString

export type WidthValue = WithEscapeHatch<CssDimensionGlobals | TokenValue<"sizes"> | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10/12" | "11/12" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "3/12" | "3/4" | "3/5" | "3/6" | "4/12" | "4/5" | "4/6" | "5/12" | "5/6" | "6/12" | "7/12" | "8/12" | "9/12" | "auto" | "screen" | CssVars>

export type WordBreakValue = WithEscapeHatch<PropertyValueMap["wordBreak"]>

export type ZIndexValue = WithEscapeHatch<CssAutoGlobals | TokenValue<"zIndex"> | CssVars>

export type AtRuleType = "media" | "layer" | "container" | "supports" | "page" | "scope" | "starting-style"

export type Selector = `${string}&` | `&${string}` | `@${AtRuleType}${string}`

export type AnySelector = Selector | string

export type CssGlobals = "inherit" | "initial" | "revert" | "revert-layer" | "unset"

export type CssColorGlobals = CssGlobals | "currentColor" | "transparent"

export type CssDimensionGlobals = CssGlobals | "auto" | "fit-content" | "max-content" | "min-content"

export type CssAutoGlobals = CssGlobals | "auto"

export type CssLength = CssDimensionGlobals | (string & {}) | number

export type CssNamedColor = "aliceblue" | "antiquewhite" | "aqua" | "aquamarine" | "azure" | "beige" | "bisque" | "black" | "blanchedalmond" | "blue" | "blueviolet" | "brown" | "burlywood" | "cadetblue" | "chartreuse" | "chocolate" | "coral" | "cornflowerblue" | "cornsilk" | "crimson" | "cyan" | "darkblue" | "darkcyan" | "darkgoldenrod" | "darkgray" | "darkgreen" | "darkgrey" | "darkkhaki" | "darkmagenta" | "darkolivegreen" | "darkorange" | "darkorchid" | "darkred" | "darksalmon" | "darkseagreen" | "darkslateblue" | "darkslategray" | "darkslategrey" | "darkturquoise" | "darkviolet" | "deeppink" | "deepskyblue" | "dimgray" | "dimgrey" | "dodgerblue" | "firebrick" | "floralwhite" | "forestgreen" | "fuchsia" | "gainsboro" | "ghostwhite" | "gold" | "goldenrod" | "gray" | "green" | "greenyellow" | "grey" | "honeydew" | "hotpink" | "indianred" | "indigo" | "ivory" | "khaki" | "lavender" | "lavenderblush" | "lawngreen" | "lemonchiffon" | "lightblue" | "lightcoral" | "lightcyan" | "lightgoldenrodyellow" | "lightgray" | "lightgreen" | "lightgrey" | "lightpink" | "lightsalmon" | "lightseagreen" | "lightskyblue" | "lightslategray" | "lightslategrey" | "lightsteelblue" | "lightyellow" | "lime" | "limegreen" | "linen" | "magenta" | "maroon" | "mediumaquamarine" | "mediumblue" | "mediumorchid" | "mediumpurple" | "mediumseagreen" | "mediumslateblue" | "mediumspringgreen" | "mediumturquoise" | "mediumvioletred" | "midnightblue" | "mintcream" | "mistyrose" | "moccasin" | "navajowhite" | "navy" | "oldlace" | "olive" | "olivedrab" | "orange" | "orangered" | "orchid" | "palegoldenrod" | "palegreen" | "paleturquoise" | "palevioletred" | "papayawhip" | "peachpuff" | "peru" | "pink" | "plum" | "powderblue" | "purple" | "rebeccapurple" | "red" | "rosybrown" | "royalblue" | "saddlebrown" | "salmon" | "sandybrown" | "seagreen" | "seashell" | "sienna" | "silver" | "skyblue" | "slateblue" | "slategray" | "slategrey" | "snow" | "springgreen" | "steelblue" | "tan" | "teal" | "thistle" | "tomato" | "turquoise" | "violet" | "wheat" | "white" | "whitesmoke" | "yellow" | "yellowgreen"

export type CssSystemColor = "AccentColor" | "AccentColorText" | "ActiveText" | "ButtonBorder" | "ButtonFace" | "ButtonText" | "Canvas" | "CanvasText" | "Field" | "FieldText" | "GrayText" | "Highlight" | "HighlightText" | "LinkText" | "Mark" | "MarkText" | "SelectedItem" | "SelectedItemText" | "VisitedText"

export type CssColor = CssColorGlobals | CssNamedColor | CssSystemColor | (string & {})

export type CssAbsoluteSize = "large" | "medium" | "small" | "x-large" | "x-small" | "xx-large" | "xx-small" | "xxx-large"

export type CssBgSize = CssLength | "auto" | "contain" | "cover"

export type CssFontSize = CssLength | CssAbsoluteSize | "larger" | "smaller" | "math"

export type CssFontWeight = CssGlobals | "bold" | "normal" | (string & {}) | number

export type CssLineWidth = CssLength | "medium" | "thick" | "thin"

export type CssLineStyle = CssGlobals | "dashed" | "dotted" | "double" | "groove" | "hidden" | "inset" | "none" | "outset" | "ridge" | "solid"

export type CssLineStyleOpen = CssLineStyle | (string & {})

export type CssRepeatStyle = CssGlobals | "no-repeat" | "repeat" | "repeat-x" | "repeat-y" | "round" | "space" | (string & {})

export type CssFontStretch = CssGlobals | "condensed" | "expanded" | "extra-condensed" | "extra-expanded" | "normal" | "semi-condensed" | "semi-expanded" | "ultra-condensed" | "ultra-expanded" | (string & {})

export type CssOverflow = CssGlobals | "-moz-hidden-unscrollable" | "auto" | "clip" | "hidden" | "overlay" | "scroll" | "visible"

export type CssOverflowOpen = CssOverflow | (string & {})

export type CssOverflowShort = CssGlobals | "auto" | "clip" | "hidden" | "scroll" | "visible"

export type CssOverscrollBehavior = CssGlobals | "auto" | "contain" | "none"

export type CssOverscrollBehaviorOpen = CssOverscrollBehavior | (string & {})

export type CssPosition = CssLength | "bottom" | "center" | "left" | "right" | "top"

export type CssGenericFamily = "-apple-system" | "cursive" | "emoji" | "fangsong" | "fantasy" | "math" | "monospace" | "sans-serif" | "serif" | "system-ui" | "ui-monospace" | "ui-rounded" | "ui-sans-serif" | "ui-serif"

export type CssFontFamily = CssGlobals | CssGenericFamily | (string & {})

export type CssNumeric = CssGlobals | (string & {}) | number

export type CssZIndex = CssAutoGlobals | (string & {}) | number

export type CssAny = CssGlobals | (string & {}) | number

export type PropertyValueMap = {
  accentColor: CssColor
  alignContent: CssGlobals | "baseline" | "center" | "end" | "first baseline" | "flex-end" | "flex-start" | "last baseline" | "normal" | "safe center" | "safe end" | "safe start" | "space-around" | "space-between" | "space-evenly" | "start" | "stretch" | (string & {})
  alignItems: CssGlobals | "anchor-center" | "baseline" | "center" | "end" | "first baseline" | "flex-end" | "flex-start" | "last baseline" | "normal" | "safe center" | "safe end" | "safe start" | "self-end" | "self-start" | "start" | "stretch" | (string & {})
  alignSelf: CssGlobals | "anchor-center" | "auto" | "baseline" | "center" | "end" | "first baseline" | "flex-end" | "flex-start" | "last baseline" | "normal" | "safe center" | "safe end" | "safe start" | "self-end" | "self-start" | "start" | "stretch" | (string & {})
  alignmentBaseline: CssGlobals | "alphabetic" | "baseline" | "central" | "ideographic" | "mathematical" | "middle" | "text-after-edge" | "text-before-edge"
  appearance: CssGlobals | "auto" | "button" | "checkbox" | "listbox" | "menulist" | "menulist-button" | "meter" | "none" | "progress-bar" | "radio" | "searchfield" | "textarea" | "textfield"
  backgroundColor: CssColor
  backgroundPosition: CssPosition
  backgroundRepeat: CssRepeatStyle
  backgroundSize: CssBgSize
  baselineShift: CssLength
  blockSize: CssLength
  borderBlockColor: CssColor
  borderBlockEndColor: CssColor
  borderBlockEndStyle: CssLineStyle
  borderBlockEndWidth: CssLineWidth
  borderBlockStartColor: CssColor
  borderBlockStartStyle: CssLineStyle
  borderBlockStartWidth: CssLineWidth
  borderBlockStyle: CssLineStyleOpen
  borderBlockWidth: CssLineWidth
  borderBottomColor: CssColor
  borderBottomLeftRadius: CssLength
  borderBottomRightRadius: CssLength
  borderBottomStyle: CssLineStyle
  borderBottomWidth: CssLineWidth
  borderColor: CssColor
  borderEndEndRadius: CssLength
  borderEndStartRadius: CssLength
  borderImage: CssGlobals | "none" | "repeat" | "round" | "space" | "stretch" | (string & {})
  borderImageWidth: CssLength
  borderInlineColor: CssColor
  borderInlineEndColor: CssColor
  borderInlineEndStyle: CssLineStyle
  borderInlineEndWidth: CssLineWidth
  borderInlineStartColor: CssColor
  borderInlineStartStyle: CssLineStyle
  borderInlineStartWidth: CssLineWidth
  borderInlineStyle: CssLineStyleOpen
  borderInlineWidth: CssLineWidth
  borderLeftColor: CssColor
  borderLeftStyle: CssLineStyle
  borderLeftWidth: CssLineWidth
  borderRadius: CssLength
  borderRightColor: CssColor
  borderRightStyle: CssLineStyle
  borderRightWidth: CssLineWidth
  borderSpacing: CssLength
  borderStartEndRadius: CssLength
  borderStartStartRadius: CssLength
  borderStyle: CssLineStyleOpen
  borderTopColor: CssColor
  borderTopLeftRadius: CssLength
  borderTopRightRadius: CssLength
  borderTopStyle: CssLineStyle
  borderTopWidth: CssLineWidth
  borderWidth: CssLineWidth
  bottom: CssLength
  boxAlign: CssGlobals | "baseline" | "center" | "end" | "start" | "stretch"
  breakAfter: CssGlobals | "all" | "always" | "auto" | "avoid" | "avoid-column" | "avoid-page" | "avoid-region" | "column" | "left" | "page" | "recto" | "region" | "right" | "verso"
  breakBefore: PropertyValueMap["breakAfter"]
  breakInside: CssGlobals | "auto" | "avoid" | "avoid-column" | "avoid-page" | "avoid-region"
  caret: CssGlobals | "auto" | "bar" | "block" | "currentColor" | "underscore" | (string & {})
  caretColor: CssColor
  clear: CssGlobals | "both" | "inline-end" | "inline-start" | "left" | "none" | "right"
  clipPath: CssGlobals | "fill-box" | "margin-box" | "none" | "stroke-box" | "view-box" | (string & {})
  color: CssColor
  colorScheme: CssGlobals | "dark" | "light" | "normal" | (string & {})
  columnGap: CssLength
  columnRuleColor: CssColor
  columnRuleStyle: CssLineStyleOpen
  columnRuleWidth: CssLineWidth
  columnWidth: CssLength
  contain: CssGlobals | "content" | "inline-size" | "layout" | "none" | "paint" | "size" | "strict" | "style" | (string & {})
  containIntrinsicBlockSize: CssLength
  containIntrinsicHeight: CssLength
  containIntrinsicInlineSize: CssLength
  containIntrinsicSize: CssLength
  containIntrinsicWidth: CssLength
  containerType: CssGlobals | "inline-size" | "normal" | "scroll-state" | "size" | (string & {})
  content: CssGlobals | "close-quote" | "no-close-quote" | "no-open-quote" | "none" | "normal" | "open-quote" | (string & {})
  cursor: CssGlobals | "-moz-grab" | "-moz-zoom-in" | "-moz-zoom-out" | "-webkit-grab" | "-webkit-grabbing" | "-webkit-zoom-in" | "-webkit-zoom-out" | "alias" | "all-scroll" | "auto" | "cell" | "col-resize" | "context-menu" | "copy" | "crosshair" | "default" | "e-resize" | "ew-resize" | "grab" | "grabbing" | "help" | "move" | "n-resize" | "ne-resize" | "nesw-resize" | "no-drop" | "none" | "not-allowed" | "ns-resize" | "nw-resize" | "nwse-resize" | "pointer" | "progress" | "row-resize" | "s-resize" | "se-resize" | "sw-resize" | "text" | "vertical-text" | "w-resize" | "wait" | "zoom-in" | "zoom-out" | (string & {})
  display: CssGlobals | "-ms-flexbox" | "-ms-grid" | "-ms-inline-flexbox" | "-ms-inline-grid" | "-webkit-flex" | "-webkit-inline-flex" | "block" | "contents" | "flex" | "flow" | "flow-root" | "grid" | "inline" | "inline-block" | "inline-flex" | "inline-grid" | "inline-list-item" | "inline-table" | "list-item" | "none" | "ruby" | "ruby-base" | "ruby-base-container" | "ruby-text" | "ruby-text-container" | "run-in" | "table" | "table-caption" | "table-cell" | "table-column" | "table-column-group" | "table-footer-group" | "table-header-group" | "table-row" | "table-row-group" | (string & {})
  dominantBaseline: CssGlobals | "alphabetic" | "auto" | "central" | "hanging" | "ideographic" | "mathematical" | "middle" | "text-bottom" | "text-top"
  fieldSizing: CssGlobals | "content" | "fixed"
  flexBasis: CssLength
  flexFlow: CssGlobals | "column" | "column-reverse" | "nowrap" | "row" | "row-reverse" | "wrap" | "wrap-reverse" | (string & {})
  flexGrow: CssNumeric
  flexShrink: CssNumeric
  float: CssGlobals | "inline-end" | "inline-start" | "left" | "none" | "right"
  floodColor: CssColor
  font: CssGlobals | "caption" | "icon" | "menu" | "message-box" | "small-caption" | "status-bar" | (string & {})
  fontFamily: CssFontFamily
  fontSize: CssFontSize
  fontSizeAdjust: CssLength
  fontStretch: CssFontStretch
  fontStyle: CssGlobals | "italic" | "normal" | "oblique" | (string & {})
  fontSynthesis: CssGlobals | "none" | "position" | "small-caps" | "style" | "weight" | (string & {})
  fontVariantCaps: CssGlobals | "all-petite-caps" | "all-small-caps" | "normal" | "petite-caps" | "small-caps" | "titling-caps" | "unicase"
  fontVariantEastAsian: CssGlobals | "full-width" | "jis04" | "jis78" | "jis83" | "jis90" | "normal" | "proportional-width" | "ruby" | "simplified" | "traditional" | (string & {})
  fontVariantLigatures: CssGlobals | "common-ligatures" | "contextual" | "discretionary-ligatures" | "historical-ligatures" | "no-common-ligatures" | "no-contextual" | "no-discretionary-ligatures" | "no-historical-ligatures" | "none" | "normal" | (string & {})
  fontVariantNumeric: CssGlobals | "diagonal-fractions" | "lining-nums" | "normal" | "oldstyle-nums" | "ordinal" | "proportional-nums" | "slashed-zero" | "stacked-fractions" | "tabular-nums" | (string & {})
  fontWeight: CssFontWeight
  fontWidth: CssFontStretch
  gap: CssLength
  gridAutoFlow: CssGlobals | "column" | "dense" | "row" | (string & {})
  gridColumnGap: CssLength
  gridGap: CssLength
  gridRowGap: CssLength
  hangingPunctuation: CssGlobals | "allow-end" | "first" | "force-end" | "last" | "none" | (string & {})
  height: CssLength
  hyphens: CssGlobals | "auto" | "manual" | "none"
  imageRendering: CssGlobals | "-moz-crisp-edges" | "-webkit-optimize-contrast" | "auto" | "crisp-edges" | "pixelated" | "smooth"
  imeMode: CssGlobals | "active" | "auto" | "disabled" | "inactive" | "normal"
  inlineSize: CssLength
  inset: CssLength
  justifyContent: CssGlobals | "center" | "end" | "flex-end" | "flex-start" | "left" | "normal" | "right" | "safe center" | "safe end" | "safe start" | "space-around" | "space-between" | "space-evenly" | "start" | "stretch" | (string & {})
  justifyItems: CssGlobals | "anchor-center" | "baseline" | "center" | "end" | "first baseline" | "flex-end" | "flex-start" | "last baseline" | "left" | "legacy" | "normal" | "right" | "safe center" | "safe end" | "safe start" | "self-end" | "self-start" | "start" | "stretch" | (string & {})
  justifySelf: CssGlobals | "anchor-center" | "auto" | "baseline" | "center" | "end" | "first baseline" | "flex-end" | "flex-start" | "last baseline" | "left" | "normal" | "right" | "safe center" | "safe end" | "safe start" | "self-end" | "self-start" | "start" | "stretch" | (string & {})
  left: CssLength
  letterSpacing: CssLength
  lightingColor: CssColor
  lineBreak: CssGlobals | "anywhere" | "auto" | "loose" | "normal" | "strict"
  lineHeight: CssNumeric
  lineHeightStep: CssLength
  listStyleImage: CssGlobals | "none" | (string & {})
  listStylePosition: CssGlobals | "inside" | "outside"
  listStyleType: PropertyValueMap["listStyleImage"]
  margin: CssLength
  marginBottom: CssLength
  marginLeft: CssLength
  marginRight: CssLength
  marginTop: CssLength
  maskBorder: CssGlobals | "alpha" | "luminance" | "none" | "repeat" | "round" | "space" | "stretch" | (string & {})
  maskBorderWidth: CssLength
  maskPosition: CssPosition
  maskRepeat: CssRepeatStyle
  maskSize: CssBgSize
  maxBlockSize: CssLength
  maxHeight: CssLength
  maxInlineSize: CssLength
  maxWidth: CssLength
  minBlockSize: CssLength
  minHeight: CssLength
  minInlineSize: CssLength
  minWidth: CssLength
  mixBlendMode: CssGlobals | "color" | "color-burn" | "color-dodge" | "darken" | "difference" | "exclusion" | "hard-light" | "hue" | "lighten" | "luminosity" | "multiply" | "normal" | "overlay" | "plus-darker" | "plus-lighter" | "saturation" | "screen" | "soft-light"
  objectFit: CssGlobals | "contain" | "cover" | "fill" | "none" | "scale-down"
  objectPosition: CssPosition
  offsetAnchor: CssPosition
  offsetPosition: CssPosition
  opacity: CssNumeric
  order: CssNumeric
  orphans: CssNumeric
  outlineColor: CssColor
  outlineOffset: CssLength
  outlineStyle: CssGlobals | "auto" | "dashed" | "dotted" | "double" | "groove" | "inset" | "none" | "outset" | "ridge" | "solid"
  outlineWidth: CssLineWidth
  overflow: CssOverflowOpen
  overflowBlock: CssOverflowShort
  overflowClipMargin: CssLength
  overflowInline: CssOverflowShort
  overflowX: CssOverflow
  overflowY: CssOverflow
  overscrollBehavior: CssOverscrollBehaviorOpen
  overscrollBehaviorX: CssOverscrollBehavior
  overscrollBehaviorY: CssOverscrollBehavior
  padding: CssLength
  paddingBottom: CssLength
  paddingLeft: CssLength
  paddingRight: CssLength
  paddingTop: CssLength
  pageBreakAfter: CssGlobals | "always" | "auto" | "avoid" | "left" | "recto" | "right" | "verso"
  pageBreakBefore: PropertyValueMap["pageBreakAfter"]
  perspective: CssLength
  perspectiveOrigin: CssPosition
  placeContent: PropertyValueMap["alignContent"]
  placeItems: PropertyValueMap["alignItems"]
  placeSelf: PropertyValueMap["alignSelf"]
  pointerEvents: CssGlobals | "all" | "auto" | "fill" | "none" | "painted" | "stroke" | "visible" | "visibleFill" | "visiblePainted" | "visibleStroke"
  position: CssGlobals | "-webkit-sticky" | "absolute" | "fixed" | "relative" | "static" | "sticky"
  positionTryOrder: CssGlobals | "most-block-size" | "most-height" | "most-inline-size" | "most-width" | "normal"
  resize: CssGlobals | "block" | "both" | "horizontal" | "inline" | "none" | "vertical"
  right: CssLength
  rowGap: CssLength
  scrollMarginBlock: CssLength
  scrollMarginBlockEnd: CssLength
  scrollMarginBlockStart: CssLength
  scrollMarginBottom: CssLength
  scrollMarginInline: CssLength
  scrollMarginInlineEnd: CssLength
  scrollMarginInlineStart: CssLength
  scrollMarginLeft: CssLength
  scrollMarginRight: CssLength
  scrollMarginTop: CssLength
  scrollPaddingBlock: CssLength
  scrollPaddingBlockEnd: CssLength
  scrollPaddingBlockStart: CssLength
  scrollPaddingBottom: CssLength
  scrollPaddingInline: CssLength
  scrollPaddingInlineEnd: CssLength
  scrollPaddingInlineStart: CssLength
  scrollPaddingLeft: CssLength
  scrollPaddingRight: CssLength
  scrollPaddingTop: CssLength
  scrollSnapAlign: CssGlobals | "center" | "end" | "none" | "start" | (string & {})
  scrollSnapDestination: CssPosition
  scrollSnapStop: CssGlobals | "always" | "normal"
  scrollSnapType: CssGlobals | "block" | "both" | "inline" | "none" | "x" | "y" | (string & {})
  scrollbarGutter: CssGlobals | "auto" | "stable" | "stable both-edges" | (string & {})
  scrollbarWidth: CssGlobals | "auto" | "none" | "thin"
  shapeMargin: CssLength
  shapeOutside: CssGlobals | "border-box" | "content-box" | "margin-box" | "none" | "padding-box" | (string & {})
  speakAs: CssGlobals | "digits" | "literal-punctuation" | "no-punctuation" | "normal" | "spell-out" | (string & {})
  stopColor: CssColor
  strokeColor: CssColor
  strokeLinejoin: CssGlobals | "arcs" | "bevel" | "miter" | "miter-clip" | "round"
  strokeWidth: CssLength
  tabSize: CssLength
  tableLayout: CssGlobals | "auto" | "fixed"
  textAlign: CssGlobals | "-khtml-center" | "-khtml-left" | "-khtml-right" | "-moz-center" | "-moz-left" | "-moz-right" | "-webkit-center" | "-webkit-left" | "-webkit-match-parent" | "-webkit-right" | "center" | "end" | "justify" | "left" | "match-parent" | "right" | "start"
  textAlignLast: CssGlobals | "auto" | "center" | "end" | "justify" | "left" | "right" | "start"
  textAutospace: CssGlobals | "auto" | "ideograph-alpha" | "ideograph-numeric" | "insert" | "no-autospace" | "normal" | "punctuation" | "replace" | (string & {})
  textBox: CssGlobals | "auto" | "cap" | "ex" | "ideographic" | "ideographic-ink" | "none" | "normal" | "text" | "trim-both" | "trim-end" | "trim-start" | (string & {})
  textBoxEdge: CssGlobals | "auto" | "cap" | "ex" | "ideographic" | "ideographic-ink" | "text" | (string & {})
  textDecorationColor: CssColor
  textDecorationLine: CssGlobals | "blink" | "grammar-error" | "line-through" | "none" | "overline" | "spelling-error" | "underline" | (string & {})
  textDecorationSkip: CssGlobals | "box-decoration" | "edges" | "leading-spaces" | "none" | "objects" | "spaces" | "trailing-spaces" | (string & {})
  textDecorationSkipInk: CssGlobals | "all" | "auto" | "none"
  textDecorationStyle: CssGlobals | "dashed" | "dotted" | "double" | "solid" | "wavy"
  textDecorationThickness: CssLength
  textEmphasis: CssGlobals | "circle" | "currentColor" | "dot" | "double-circle" | "filled" | "none" | "open" | "sesame" | "triangle" | (string & {})
  textEmphasisColor: CssColor
  textEmphasisPosition: CssGlobals | "auto" | "over" | "under" | (string & {})
  textEmphasisStyle: CssGlobals | "circle" | "dot" | "double-circle" | "filled" | "none" | "open" | "sesame" | "triangle" | (string & {})
  textIndent: CssLength
  textJustify: CssGlobals | "auto" | "distribute" | "inter-character" | "inter-word" | "none"
  textOrientation: CssGlobals | "mixed" | "sideways" | "sideways-right" | "upright"
  textOverflow: CssGlobals | "clip" | "ellipsis" | (string & {})
  textRendering: CssGlobals | "auto" | "geometricPrecision" | "optimizeLegibility" | "optimizeSpeed"
  textSizeAdjust: CssLength
  textTransform: CssGlobals | "capitalize" | "full-size-kana" | "full-width" | "lowercase" | "math-auto" | "none" | "uppercase" | (string & {})
  textUnderlineOffset: CssLength
  textUnderlinePosition: CssGlobals | "auto" | "from-font" | "left" | "right" | "under" | (string & {})
  textWrap: CssGlobals | "auto" | "balance" | "nowrap" | "pretty" | "stable" | "wrap" | (string & {})
  textWrapMode: CssGlobals | "nowrap" | "wrap"
  textWrapStyle: CssGlobals | "auto" | "balance" | "pretty" | "stable"
  top: CssLength
  touchAction: CssGlobals | "-ms-manipulation" | "-ms-none" | "-ms-pan-x" | "-ms-pan-y" | "-ms-pinch-zoom" | "auto" | "manipulation" | "none" | "pan-down" | "pan-left" | "pan-right" | "pan-up" | "pan-x" | "pan-y" | "pinch-zoom" | (string & {})
  transformBox: CssGlobals | "border-box" | "content-box" | "fill-box" | "stroke-box" | "view-box"
  transitionBehavior: CssGlobals | "allow-discrete" | "normal" | (string & {})
  transitionProperty: CssGlobals | "all" | "none" | (string & {})
  unicodeBidi: CssGlobals | "-moz-isolate" | "-moz-isolate-override" | "-moz-plaintext" | "-webkit-isolate" | "-webkit-isolate-override" | "-webkit-plaintext" | "bidi-override" | "embed" | "isolate" | "isolate-override" | "normal" | "plaintext"
  userSelect: CssGlobals | "-moz-none" | "all" | "auto" | "none" | "text"
  vectorEffect: CssGlobals | "fixed-position" | "non-rotation" | "non-scaling-size" | "non-scaling-stroke" | "none"
  verticalAlign: CssGlobals | "baseline" | "bottom" | "middle" | "sub" | "super" | "text-bottom" | "text-top" | "top" | (string & {})
  viewTimelineInset: CssLength
  whiteSpace: CssGlobals | "-moz-pre-wrap" | "break-spaces" | "collapse" | "normal" | "nowrap" | "pre" | "pre-line" | "pre-wrap" | "preserve" | "preserve-breaks" | "preserve-spaces" | "wrap" | (string & {})
  whiteSpaceCollapse: CssGlobals | "break-spaces" | "collapse" | "preserve" | "preserve-breaks" | "preserve-spaces"
  widows: CssNumeric
  width: CssLength
  wordBreak: CssGlobals | "auto-phrase" | "break-all" | "break-word" | "keep-all" | "normal"
  wordSpacing: CssLength
  wordWrap: CssGlobals | "break-word" | "normal"
  writingMode: CssGlobals | "horizontal-tb" | "sideways-lr" | "sideways-rl" | "vertical-lr" | "vertical-rl"
  zIndex: CssZIndex
  zoom: CssNumeric
}

export interface SystemProperties {
  WebkitAppearance?: ConditionalValue<CssAny>
  WebkitBorderBefore?: ConditionalValue<CssAny>
  WebkitBorderBeforeColor?: ConditionalValue<CssAny>
  WebkitBorderBeforeStyle?: ConditionalValue<CssAny>
  WebkitBorderBeforeWidth?: ConditionalValue<CssAny>
  WebkitBoxReflect?: ConditionalValue<CssAny>
  WebkitLineClamp?: ConditionalValue<CssAny>
  WebkitMask?: ConditionalValue<CssAny>
  WebkitMaskAttachment?: ConditionalValue<CssAny>
  WebkitMaskClip?: ConditionalValue<CssAny>
  WebkitMaskComposite?: ConditionalValue<CssAny>
  WebkitMaskImage?: ConditionalValue<CssAny>
  WebkitMaskOrigin?: ConditionalValue<CssAny>
  WebkitMaskPosition?: ConditionalValue<CssAny>
  WebkitMaskPositionX?: ConditionalValue<CssAny>
  WebkitMaskPositionY?: ConditionalValue<CssAny>
  WebkitMaskRepeat?: ConditionalValue<CssAny>
  WebkitMaskRepeatX?: ConditionalValue<CssAny>
  WebkitMaskRepeatY?: ConditionalValue<CssAny>
  WebkitMaskSize?: ConditionalValue<CssAny>
  WebkitOverflowScrolling?: ConditionalValue<CssAny>
  WebkitTapHighlightColor?: ConditionalValue<CssAny>
  WebkitTextFillColor?: ConditionalValue<ColorsValue>
  WebkitTextStroke?: ConditionalValue<CssAny>
  WebkitTextStrokeColor?: ConditionalValue<CssAny>
  WebkitTextStrokeWidth?: ConditionalValue<CssAny>
  WebkitTouchCallout?: ConditionalValue<CssAny>
  WebkitUserModify?: ConditionalValue<CssAny>
  WebkitUserSelect?: ConditionalValue<CssAny>
  accentColor?: ConditionalValue<ColorsValue>
  alignContent?: ConditionalValue<AlignContentValue>
  alignItems?: ConditionalValue<AlignItemsValue>
  alignSelf?: ConditionalValue<AlignSelfValue>
  alignTracks?: ConditionalValue<CssAny>
  alignmentBaseline?: ConditionalValue<PropertyValueMap["alignmentBaseline"] | AnyString>
  all?: ConditionalValue<CssAny>
  anchorName?: ConditionalValue<CssAny>
  anchorScope?: ConditionalValue<CssAny>
  animation?: ConditionalValue<AnimationsValue>
  animationComposition?: ConditionalValue<AnimationCompositionValue>
  animationDelay?: ConditionalValue<DurationsValue>
  animationDirection?: ConditionalValue<AnimationDirectionValue>
  animationDuration?: ConditionalValue<DurationsValue>
  animationFillMode?: ConditionalValue<AnimationFillModeValue>
  animationIterationCount?: ConditionalValue<AnimationIterationCountValue>
  animationName?: ConditionalValue<KeyframesValue | WithEscapeHatch<"none">>
  animationPlayState?: ConditionalValue<AnimationPlayStateValue>
  animationRange?: ConditionalValue<AnimationRangeValue>
  animationRangeEnd?: ConditionalValue<AnimationRangeEndValue>
  animationRangeStart?: ConditionalValue<AnimationRangeStartValue>
  animationTimeline?: ConditionalValue<AnimationTimelineValue>
  animationTimingFunction?: ConditionalValue<EasingsValue>
  appearance?: ConditionalValue<AppearanceValue>
  aspectRatio?: ConditionalValue<AspectRatiosValue>
  backdropFilter?: ConditionalValue<BackdropFilterValue | WithEscapeHatch<"none">>
  backfaceVisibility?: ConditionalValue<BackfaceVisibilityValue>
  background?: ConditionalValue<ColorsValue>
  backgroundAttachment?: ConditionalValue<BackgroundAttachmentValue>
  backgroundBlendMode?: ConditionalValue<BackgroundBlendModeValue>
  backgroundClip?: ConditionalValue<BackgroundClipValue>
  backgroundColor?: ConditionalValue<ColorsValue>
  backgroundImage?: ConditionalValue<AssetsValue>
  backgroundOrigin?: ConditionalValue<BackgroundOriginValue>
  backgroundPosition?: ConditionalValue<BackgroundPositionValue>
  backgroundPositionX?: ConditionalValue<BackgroundPositionXValue>
  backgroundPositionY?: ConditionalValue<BackgroundPositionYValue>
  backgroundRepeat?: ConditionalValue<BackgroundRepeatValue>
  backgroundSize?: ConditionalValue<BackgroundSizeValue>
  baselineShift?: ConditionalValue<PropertyValueMap["baselineShift"]>
  blockSize?: ConditionalValue<BlockSizeValue>
  border?: ConditionalValue<BordersValue>
  borderBlock?: ConditionalValue<BordersValue>
  borderBlockColor?: ConditionalValue<ColorsValue>
  borderBlockEnd?: ConditionalValue<BordersValue>
  borderBlockEndColor?: ConditionalValue<ColorsValue>
  borderBlockEndStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["borderBlockEndStyle"]>>
  borderBlockEndWidth?: ConditionalValue<BorderWidthsValue>
  borderBlockStart?: ConditionalValue<BordersValue>
  borderBlockStartColor?: ConditionalValue<ColorsValue>
  borderBlockStartStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["borderBlockStartStyle"]>>
  borderBlockStartWidth?: ConditionalValue<BorderWidthsValue>
  borderBlockStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["borderBlockStyle"]>>
  borderBlockWidth?: ConditionalValue<BorderWidthsValue>
  borderBottom?: ConditionalValue<BordersValue>
  borderBottomColor?: ConditionalValue<ColorsValue>
  borderBottomLeftRadius?: ConditionalValue<RadiiValue>
  borderBottomRightRadius?: ConditionalValue<RadiiValue>
  borderBottomStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["borderBottomStyle"]>>
  borderBottomWidth?: ConditionalValue<BorderWidthsValue>
  borderCollapse?: ConditionalValue<BorderCollapseValue>
  borderColor?: ConditionalValue<ColorsValue>
  borderEndEndRadius?: ConditionalValue<RadiiValue>
  borderEndStartRadius?: ConditionalValue<RadiiValue>
  borderImage?: ConditionalValue<PropertyValueMap["borderImage"]>
  borderImageOutset?: ConditionalValue<CssAny>
  borderImageRepeat?: ConditionalValue<CssAny>
  borderImageSlice?: ConditionalValue<CssAny>
  borderImageSource?: ConditionalValue<CssAny>
  borderImageWidth?: ConditionalValue<PropertyValueMap["borderImageWidth"]>
  borderInline?: ConditionalValue<BordersValue>
  borderInlineColor?: ConditionalValue<ColorsValue>
  borderInlineEnd?: ConditionalValue<BordersValue>
  borderInlineEndColor?: ConditionalValue<ColorsValue>
  borderInlineEndStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["borderInlineEndStyle"]>>
  borderInlineEndWidth?: ConditionalValue<BorderWidthsValue>
  borderInlineStart?: ConditionalValue<BordersValue>
  borderInlineStartColor?: ConditionalValue<ColorsValue>
  borderInlineStartStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["borderInlineStartStyle"]>>
  borderInlineStartWidth?: ConditionalValue<BorderWidthsValue>
  borderInlineStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["borderInlineStyle"]>>
  borderInlineWidth?: ConditionalValue<BorderWidthsValue>
  borderLeft?: ConditionalValue<BordersValue>
  borderLeftColor?: ConditionalValue<ColorsValue>
  borderLeftStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["borderLeftStyle"]>>
  borderLeftWidth?: ConditionalValue<BorderWidthsValue>
  borderRadius?: ConditionalValue<RadiiValue>
  borderRight?: ConditionalValue<BordersValue>
  borderRightColor?: ConditionalValue<ColorsValue>
  borderRightStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["borderRightStyle"]>>
  borderRightWidth?: ConditionalValue<BorderWidthsValue>
  borderSpacing?: ConditionalValue<SpacingValue>
  borderStartEndRadius?: ConditionalValue<RadiiValue>
  borderStartStartRadius?: ConditionalValue<RadiiValue>
  borderStyle?: ConditionalValue<PropertyValueMap["borderStyle"]>
  borderTop?: ConditionalValue<BordersValue>
  borderTopColor?: ConditionalValue<ColorsValue>
  borderTopLeftRadius?: ConditionalValue<RadiiValue>
  borderTopRightRadius?: ConditionalValue<RadiiValue>
  borderTopStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["borderTopStyle"]>>
  borderTopWidth?: ConditionalValue<BorderWidthsValue>
  borderWidth?: ConditionalValue<BorderWidthsValue>
  bottom?: ConditionalValue<SpacingValue>
  boxAlign?: ConditionalValue<PropertyValueMap["boxAlign"] | AnyString>
  boxDecorationBreak?: ConditionalValue<BoxDecorationBreakValue>
  boxDirection?: ConditionalValue<CssAny>
  boxFlex?: ConditionalValue<CssAny>
  boxFlexGroup?: ConditionalValue<CssAny>
  boxLines?: ConditionalValue<CssAny>
  boxOrdinalGroup?: ConditionalValue<CssAny>
  boxOrient?: ConditionalValue<CssAny>
  boxPack?: ConditionalValue<CssAny>
  boxShadow?: ConditionalValue<ShadowsValue | WithEscapeHatch<"none">>
  boxSizing?: ConditionalValue<BoxSizingValue>
  breakAfter?: ConditionalValue<WithEscapeHatch<PropertyValueMap["breakAfter"]>>
  breakBefore?: ConditionalValue<WithEscapeHatch<PropertyValueMap["breakBefore"]>>
  breakInside?: ConditionalValue<WithEscapeHatch<PropertyValueMap["breakInside"]>>
  captionSide?: ConditionalValue<CssAny>
  caret?: ConditionalValue<PropertyValueMap["caret"]>
  caretColor?: ConditionalValue<ColorsValue>
  caretShape?: ConditionalValue<CssAny>
  clear?: ConditionalValue<WithEscapeHatch<PropertyValueMap["clear"]>>
  clip?: ConditionalValue<CssAny>
  clipPath?: ConditionalValue<ClipPathValue>
  clipRule?: ConditionalValue<CssAny>
  color?: ConditionalValue<ColorsValue>
  colorInterpolation?: ConditionalValue<CssAny>
  colorInterpolationFilters?: ConditionalValue<CssAny>
  colorRendering?: ConditionalValue<CssAny>
  colorScheme?: ConditionalValue<PropertyValueMap["colorScheme"]>
  columnCount?: ConditionalValue<CssAny>
  columnFill?: ConditionalValue<CssAny>
  columnGap?: ConditionalValue<SpacingValue>
  columnRule?: ConditionalValue<CssAny>
  columnRuleColor?: ConditionalValue<PropertyValueMap["columnRuleColor"]>
  columnRuleStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["columnRuleStyle"]>>
  columnRuleWidth?: ConditionalValue<PropertyValueMap["columnRuleWidth"]>
  columnSpan?: ConditionalValue<CssAny>
  columnWidth?: ConditionalValue<PropertyValueMap["columnWidth"]>
  columns?: ConditionalValue<CssAny>
  contain?: ConditionalValue<PropertyValueMap["contain"]>
  containIntrinsicBlockSize?: ConditionalValue<PropertyValueMap["containIntrinsicBlockSize"]>
  containIntrinsicHeight?: ConditionalValue<PropertyValueMap["containIntrinsicHeight"]>
  containIntrinsicInlineSize?: ConditionalValue<PropertyValueMap["containIntrinsicInlineSize"]>
  containIntrinsicSize?: ConditionalValue<PropertyValueMap["containIntrinsicSize"]>
  containIntrinsicWidth?: ConditionalValue<PropertyValueMap["containIntrinsicWidth"]>
  container?: ConditionalValue<CssContainer>
  containerName?: ConditionalValue<ContainerName>
  containerType?: ConditionalValue<ContainerTypeValue>
  content?: ConditionalValue<PropertyValueMap["content"]>
  contentVisibility?: ConditionalValue<CssAny>
  counterIncrement?: ConditionalValue<CssAny>
  counterReset?: ConditionalValue<CssAny>
  counterSet?: ConditionalValue<CssAny>
  cursor?: ConditionalValue<CursorValue>
  cx?: ConditionalValue<CssAny>
  cy?: ConditionalValue<CssAny>
  d?: ConditionalValue<CssAny>
  direction?: ConditionalValue<CssAny>
  display?: ConditionalValue<DisplayValue>
  dominantBaseline?: ConditionalValue<PropertyValueMap["dominantBaseline"] | AnyString>
  emptyCells?: ConditionalValue<CssAny>
  fieldSizing?: ConditionalValue<PropertyValueMap["fieldSizing"] | AnyString>
  fill?: ConditionalValue<ColorsValue | WithEscapeHatch<"none">>
  fillOpacity?: ConditionalValue<CssAny>
  fillRule?: ConditionalValue<CssAny>
  filter?: ConditionalValue<FilterValue | WithEscapeHatch<"none">>
  flex?: ConditionalValue<FlexValue>
  flexBasis?: ConditionalValue<FlexBasisValue | WithEscapeHatch<"auto" | "content" | "fit-content" | "max-content" | "min-content">>
  flexDirection?: ConditionalValue<FlexDirectionValue>
  flexFlow?: ConditionalValue<PropertyValueMap["flexFlow"]>
  flexGrow?: ConditionalValue<FlexGrowValue>
  flexShrink?: ConditionalValue<FlexShrinkValue>
  flexWrap?: ConditionalValue<CssAny>
  float?: ConditionalValue<FloatValue | WithEscapeHatch<"inline-end" | "inline-start" | "left" | "none" | "right">>
  floodColor?: ConditionalValue<PropertyValueMap["floodColor"]>
  floodOpacity?: ConditionalValue<CssAny>
  font?: ConditionalValue<PropertyValueMap["font"]>
  fontFamily?: ConditionalValue<FontsValue>
  fontFeatureSettings?: ConditionalValue<FontFeatureSettingsValue>
  fontKerning?: ConditionalValue<FontKerningValue>
  fontLanguageOverride?: ConditionalValue<CssAny>
  fontOpticalSizing?: ConditionalValue<CssAny>
  fontPalette?: ConditionalValue<FontPaletteValue>
  fontSize?: ConditionalValue<FontSizesValue>
  fontSizeAdjust?: ConditionalValue<FontSizeAdjustValue>
  fontSmooth?: ConditionalValue<CssAny>
  fontStretch?: ConditionalValue<PropertyValueMap["fontStretch"]>
  fontStyle?: ConditionalValue<PropertyValueMap["fontStyle"]>
  fontSynthesis?: ConditionalValue<PropertyValueMap["fontSynthesis"]>
  fontSynthesisPosition?: ConditionalValue<CssAny>
  fontSynthesisSmallCaps?: ConditionalValue<CssAny>
  fontSynthesisStyle?: ConditionalValue<CssAny>
  fontSynthesisWeight?: ConditionalValue<CssAny>
  fontVariant?: ConditionalValue<FontVariantValue>
  fontVariantAlternates?: ConditionalValue<FontVariantAlternatesValue>
  fontVariantCaps?: ConditionalValue<FontVariantCapsValue>
  fontVariantEastAsian?: ConditionalValue<PropertyValueMap["fontVariantEastAsian"]>
  fontVariantEmoji?: ConditionalValue<CssAny>
  fontVariantLigatures?: ConditionalValue<PropertyValueMap["fontVariantLigatures"]>
  fontVariantNumeric?: ConditionalValue<FontVariantNumericValue>
  fontVariantPosition?: ConditionalValue<CssAny>
  fontVariationSettings?: ConditionalValue<FontVariationSettingsValue>
  fontWeight?: ConditionalValue<FontWeightsValue>
  fontWidth?: ConditionalValue<PropertyValueMap["fontWidth"]>
  forcedColorAdjust?: ConditionalValue<CssAny>
  gap?: ConditionalValue<SpacingValue>
  glyphOrientationVertical?: ConditionalValue<CssAny>
  grid?: ConditionalValue<CssAny>
  gridArea?: ConditionalValue<CssAny>
  gridAutoColumns?: ConditionalValue<GridAutoColumnsValue | WithEscapeHatch<"auto" | "max-content" | "min-content">>
  gridAutoFlow?: ConditionalValue<GridAutoFlowValue>
  gridAutoRows?: ConditionalValue<GridAutoRowsValue | WithEscapeHatch<"auto" | "max-content" | "min-content">>
  gridColumn?: ConditionalValue<GridColumnValue>
  gridColumnEnd?: ConditionalValue<GridColumnEndValue>
  gridColumnGap?: ConditionalValue<SpacingValue>
  gridColumnStart?: ConditionalValue<GridColumnStartValue>
  gridGap?: ConditionalValue<SpacingValue>
  gridRow?: ConditionalValue<GridRowValue>
  gridRowEnd?: ConditionalValue<CssAny>
  gridRowGap?: ConditionalValue<SpacingValue>
  gridRowStart?: ConditionalValue<CssAny>
  gridTemplate?: ConditionalValue<CssAny>
  gridTemplateAreas?: ConditionalValue<CssAny>
  gridTemplateColumns?: ConditionalValue<GridTemplateColumnsValue>
  gridTemplateRows?: ConditionalValue<GridTemplateRowsValue>
  hangingPunctuation?: ConditionalValue<PropertyValueMap["hangingPunctuation"]>
  height?: ConditionalValue<HeightValue>
  hyphenateCharacter?: ConditionalValue<CssAny>
  hyphenateLimitChars?: ConditionalValue<CssAny>
  hyphens?: ConditionalValue<HyphensValue>
  imageOrientation?: ConditionalValue<CssAny>
  imageRendering?: ConditionalValue<PropertyValueMap["imageRendering"] | AnyString>
  imageResolution?: ConditionalValue<CssAny>
  imeMode?: ConditionalValue<PropertyValueMap["imeMode"] | AnyString>
  initialLetter?: ConditionalValue<CssAny>
  initialLetterAlign?: ConditionalValue<CssAny>
  inlineSize?: ConditionalValue<InlineSizeValue>
  inset?: ConditionalValue<SpacingValue>
  insetBlock?: ConditionalValue<SpacingValue>
  insetBlockEnd?: ConditionalValue<SpacingValue>
  insetBlockStart?: ConditionalValue<SpacingValue>
  insetInline?: ConditionalValue<SpacingValue>
  insetInlineEnd?: ConditionalValue<SpacingValue>
  insetInlineStart?: ConditionalValue<SpacingValue>
  interpolateSize?: ConditionalValue<CssAny>
  isolation?: ConditionalValue<CssAny>
  justifyContent?: ConditionalValue<JustifyContentValue>
  justifyItems?: ConditionalValue<PropertyValueMap["justifyItems"]>
  justifySelf?: ConditionalValue<PropertyValueMap["justifySelf"]>
  justifyTracks?: ConditionalValue<CssAny>
  left?: ConditionalValue<SpacingValue>
  letterSpacing?: ConditionalValue<LetterSpacingsValue>
  lightingColor?: ConditionalValue<PropertyValueMap["lightingColor"]>
  lineBreak?: ConditionalValue<WithEscapeHatch<PropertyValueMap["lineBreak"]>>
  lineClamp?: ConditionalValue<LineClampValue>
  lineHeight?: ConditionalValue<LineHeightsValue>
  lineHeightStep?: ConditionalValue<PropertyValueMap["lineHeightStep"]>
  listStyle?: ConditionalValue<ListStyleValue>
  listStyleImage?: ConditionalValue<AssetsValue>
  listStylePosition?: ConditionalValue<ListStylePositionValue>
  listStyleType?: ConditionalValue<ListStyleTypeValue>
  margin?: ConditionalValue<SpacingValue>
  marginBlock?: ConditionalValue<SpacingValue>
  marginBlockEnd?: ConditionalValue<SpacingValue>
  marginBlockStart?: ConditionalValue<SpacingValue>
  marginBottom?: ConditionalValue<SpacingValue>
  marginInline?: ConditionalValue<SpacingValue>
  marginInlineEnd?: ConditionalValue<SpacingValue>
  marginInlineStart?: ConditionalValue<SpacingValue>
  marginLeft?: ConditionalValue<SpacingValue>
  marginRight?: ConditionalValue<SpacingValue>
  marginTop?: ConditionalValue<SpacingValue>
  marginTrim?: ConditionalValue<CssAny>
  marker?: ConditionalValue<CssAny>
  markerEnd?: ConditionalValue<CssAny>
  markerMid?: ConditionalValue<CssAny>
  markerStart?: ConditionalValue<CssAny>
  mask?: ConditionalValue<MaskValue>
  maskBorder?: ConditionalValue<PropertyValueMap["maskBorder"]>
  maskBorderMode?: ConditionalValue<CssAny>
  maskBorderOutset?: ConditionalValue<CssAny>
  maskBorderRepeat?: ConditionalValue<CssAny>
  maskBorderSlice?: ConditionalValue<CssAny>
  maskBorderSource?: ConditionalValue<CssAny>
  maskBorderWidth?: ConditionalValue<PropertyValueMap["maskBorderWidth"]>
  maskClip?: ConditionalValue<MaskClipValue>
  maskComposite?: ConditionalValue<MaskCompositeValue>
  maskImage?: ConditionalValue<MaskImageValue>
  maskMode?: ConditionalValue<MaskModeValue>
  maskOrigin?: ConditionalValue<MaskOriginValue>
  maskPosition?: ConditionalValue<MaskPositionValue>
  maskRepeat?: ConditionalValue<MaskRepeatValue>
  maskSize?: ConditionalValue<MaskSizeValue>
  maskType?: ConditionalValue<MaskTypeValue>
  masonryAutoFlow?: ConditionalValue<CssAny>
  mathDepth?: ConditionalValue<CssAny>
  mathShift?: ConditionalValue<CssAny>
  mathStyle?: ConditionalValue<CssAny>
  maxBlockSize?: ConditionalValue<MaxBlockSizeValue | WithEscapeHatch<"none">>
  maxHeight?: ConditionalValue<MaxHeightValue | WithEscapeHatch<"none">>
  maxInlineSize?: ConditionalValue<MaxInlineSizeValue | WithEscapeHatch<"none">>
  maxLines?: ConditionalValue<CssAny>
  maxWidth?: ConditionalValue<MaxWidthValue | WithEscapeHatch<"none">>
  minBlockSize?: ConditionalValue<MinBlockSizeValue>
  minHeight?: ConditionalValue<MinHeightValue>
  minInlineSize?: ConditionalValue<MinInlineSizeValue>
  minWidth?: ConditionalValue<MinWidthValue>
  mixBlendMode?: ConditionalValue<MixBlendModeValue>
  objectFit?: ConditionalValue<ObjectFitValue>
  objectPosition?: ConditionalValue<ObjectPositionValue>
  objectViewBox?: ConditionalValue<CssAny>
  offset?: ConditionalValue<CssAny>
  offsetAnchor?: ConditionalValue<PropertyValueMap["offsetAnchor"]>
  offsetDistance?: ConditionalValue<CssAny>
  offsetPath?: ConditionalValue<CssAny>
  offsetPosition?: ConditionalValue<PropertyValueMap["offsetPosition"]>
  offsetRotate?: ConditionalValue<CssAny>
  opacity?: ConditionalValue<OpacityValue>
  order?: ConditionalValue<PropertyValueMap["order"]>
  orphans?: ConditionalValue<PropertyValueMap["orphans"]>
  outline?: ConditionalValue<BordersValue>
  outlineColor?: ConditionalValue<ColorsValue>
  outlineOffset?: ConditionalValue<SpacingValue>
  outlineStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["outlineStyle"]>>
  outlineWidth?: ConditionalValue<BorderWidthsValue>
  overflow?: ConditionalValue<OverflowValue>
  overflowAnchor?: ConditionalValue<OverflowAnchorValue>
  overflowBlock?: ConditionalValue<OverflowBlockValue>
  overflowClipBox?: ConditionalValue<OverflowClipBoxValue>
  overflowClipMargin?: ConditionalValue<OverflowClipMarginValue>
  overflowInline?: ConditionalValue<OverflowInlineValue>
  overflowWrap?: ConditionalValue<OverflowWrapValue>
  overflowX?: ConditionalValue<OverflowXValue>
  overflowY?: ConditionalValue<OverflowYValue>
  overlay?: ConditionalValue<CssAny>
  overscrollBehavior?: ConditionalValue<OverscrollBehaviorValue>
  overscrollBehaviorBlock?: ConditionalValue<OverscrollBehaviorBlockValue>
  overscrollBehaviorInline?: ConditionalValue<OverscrollBehaviorInlineValue>
  overscrollBehaviorX?: ConditionalValue<OverscrollBehaviorXValue>
  overscrollBehaviorY?: ConditionalValue<OverscrollBehaviorYValue>
  padding?: ConditionalValue<SpacingValue>
  paddingBlock?: ConditionalValue<SpacingValue>
  paddingBlockEnd?: ConditionalValue<SpacingValue>
  paddingBlockStart?: ConditionalValue<SpacingValue>
  paddingBottom?: ConditionalValue<SpacingValue>
  paddingInline?: ConditionalValue<SpacingValue>
  paddingInlineEnd?: ConditionalValue<SpacingValue>
  paddingInlineStart?: ConditionalValue<SpacingValue>
  paddingLeft?: ConditionalValue<SpacingValue>
  paddingRight?: ConditionalValue<SpacingValue>
  paddingTop?: ConditionalValue<SpacingValue>
  page?: ConditionalValue<CssAny>
  pageBreakAfter?: ConditionalValue<PropertyValueMap["pageBreakAfter"] | AnyString>
  pageBreakBefore?: ConditionalValue<PropertyValueMap["pageBreakBefore"] | AnyString>
  pageBreakInside?: ConditionalValue<CssAny>
  paintOrder?: ConditionalValue<CssAny>
  perspective?: ConditionalValue<PropertyValueMap["perspective"]>
  perspectiveOrigin?: ConditionalValue<PropertyValueMap["perspectiveOrigin"]>
  placeContent?: ConditionalValue<PropertyValueMap["placeContent"]>
  placeItems?: ConditionalValue<PropertyValueMap["placeItems"]>
  placeSelf?: ConditionalValue<PropertyValueMap["placeSelf"]>
  pointerEvents?: ConditionalValue<WithEscapeHatch<PropertyValueMap["pointerEvents"]>>
  position?: ConditionalValue<PositionValue>
  positionAnchor?: ConditionalValue<CssAny>
  positionArea?: ConditionalValue<CssAny>
  positionTry?: ConditionalValue<CssAny>
  positionTryFallbacks?: ConditionalValue<CssAny>
  positionTryOrder?: ConditionalValue<PropertyValueMap["positionTryOrder"] | AnyString>
  positionVisibility?: ConditionalValue<CssAny>
  printColorAdjust?: ConditionalValue<CssAny>
  quotes?: ConditionalValue<CssAny>
  r?: ConditionalValue<CssAny>
  resize?: ConditionalValue<WithEscapeHatch<PropertyValueMap["resize"]>>
  right?: ConditionalValue<SpacingValue>
  rotate?: ConditionalValue<RotateValue | WithEscapeHatch<"none">>
  rowGap?: ConditionalValue<SpacingValue>
  rubyAlign?: ConditionalValue<CssAny>
  rubyMerge?: ConditionalValue<CssAny>
  rubyOverhang?: ConditionalValue<CssAny>
  rubyPosition?: ConditionalValue<CssAny>
  rx?: ConditionalValue<CssAny>
  ry?: ConditionalValue<CssAny>
  scale?: ConditionalValue<ScaleValue | WithEscapeHatch<"none">>
  scrollBehavior?: ConditionalValue<ScrollBehaviorValue>
  scrollInitialTarget?: ConditionalValue<CssAny>
  scrollMargin?: ConditionalValue<SpacingValue>
  scrollMarginBlock?: ConditionalValue<SpacingValue>
  scrollMarginBlockEnd?: ConditionalValue<SpacingValue>
  scrollMarginBlockStart?: ConditionalValue<SpacingValue>
  scrollMarginBottom?: ConditionalValue<SpacingValue>
  scrollMarginInline?: ConditionalValue<SpacingValue>
  scrollMarginInlineEnd?: ConditionalValue<SpacingValue>
  scrollMarginInlineStart?: ConditionalValue<SpacingValue>
  scrollMarginLeft?: ConditionalValue<SpacingValue>
  scrollMarginRight?: ConditionalValue<SpacingValue>
  scrollMarginTop?: ConditionalValue<SpacingValue>
  scrollPadding?: ConditionalValue<SpacingValue>
  scrollPaddingBlock?: ConditionalValue<SpacingValue>
  scrollPaddingBlockEnd?: ConditionalValue<SpacingValue>
  scrollPaddingBlockStart?: ConditionalValue<SpacingValue>
  scrollPaddingBottom?: ConditionalValue<SpacingValue>
  scrollPaddingInline?: ConditionalValue<SpacingValue>
  scrollPaddingInlineEnd?: ConditionalValue<SpacingValue>
  scrollPaddingInlineStart?: ConditionalValue<SpacingValue>
  scrollPaddingLeft?: ConditionalValue<SpacingValue>
  scrollPaddingRight?: ConditionalValue<SpacingValue>
  scrollPaddingTop?: ConditionalValue<SpacingValue>
  scrollSnapAlign?: ConditionalValue<ScrollSnapAlignValue>
  scrollSnapCoordinate?: ConditionalValue<ScrollSnapCoordinateValue>
  scrollSnapDestination?: ConditionalValue<ScrollSnapDestinationValue>
  scrollSnapPointsX?: ConditionalValue<ScrollSnapPointsXValue>
  scrollSnapPointsY?: ConditionalValue<ScrollSnapPointsYValue>
  scrollSnapStop?: ConditionalValue<ScrollSnapStopValue>
  scrollSnapType?: ConditionalValue<ScrollSnapTypeValue>
  scrollSnapTypeX?: ConditionalValue<ScrollSnapTypeXValue>
  scrollSnapTypeY?: ConditionalValue<ScrollSnapTypeYValue>
  scrollTimeline?: ConditionalValue<ScrollTimelineValue>
  scrollTimelineAxis?: ConditionalValue<ScrollTimelineAxisValue>
  scrollTimelineName?: ConditionalValue<ScrollTimelineNameValue>
  scrollbarColor?: ConditionalValue<ScrollbarColorValue>
  scrollbarGutter?: ConditionalValue<ScrollbarGutterValue>
  scrollbarWidth?: ConditionalValue<ScrollbarWidthValue>
  shapeImageThreshold?: ConditionalValue<CssAny>
  shapeMargin?: ConditionalValue<PropertyValueMap["shapeMargin"]>
  shapeOutside?: ConditionalValue<PropertyValueMap["shapeOutside"]>
  shapeRendering?: ConditionalValue<CssAny>
  speakAs?: ConditionalValue<PropertyValueMap["speakAs"]>
  stopColor?: ConditionalValue<PropertyValueMap["stopColor"]>
  stopOpacity?: ConditionalValue<CssAny>
  stroke?: ConditionalValue<ColorsValue | WithEscapeHatch<"none">>
  strokeColor?: ConditionalValue<PropertyValueMap["strokeColor"]>
  strokeDasharray?: ConditionalValue<StrokeDasharrayValue>
  strokeDashoffset?: ConditionalValue<StrokeDashoffsetValue>
  strokeLinecap?: ConditionalValue<StrokeLinecapValue>
  strokeLinejoin?: ConditionalValue<StrokeLinejoinValue>
  strokeMiterlimit?: ConditionalValue<StrokeMiterlimitValue>
  strokeOpacity?: ConditionalValue<StrokeOpacityValue>
  strokeWidth?: ConditionalValue<BorderWidthsValue>
  tabSize?: ConditionalValue<PropertyValueMap["tabSize"]>
  tableLayout?: ConditionalValue<TableLayoutValue>
  textAlign?: ConditionalValue<TextAlignValue>
  textAlignLast?: ConditionalValue<PropertyValueMap["textAlignLast"] | AnyString>
  textAnchor?: ConditionalValue<CssAny>
  textAutospace?: ConditionalValue<PropertyValueMap["textAutospace"]>
  textBox?: ConditionalValue<PropertyValueMap["textBox"]>
  textBoxEdge?: ConditionalValue<PropertyValueMap["textBoxEdge"]>
  textBoxTrim?: ConditionalValue<CssAny>
  textCombineUpright?: ConditionalValue<CssAny>
  textDecoration?: ConditionalValue<TextDecorationValue>
  textDecorationColor?: ConditionalValue<ColorsValue>
  textDecorationLine?: ConditionalValue<PropertyValueMap["textDecorationLine"]>
  textDecorationSkip?: ConditionalValue<PropertyValueMap["textDecorationSkip"]>
  textDecorationSkipInk?: ConditionalValue<PropertyValueMap["textDecorationSkipInk"] | AnyString>
  textDecorationStyle?: ConditionalValue<TextDecorationStyleValue>
  textDecorationThickness?: ConditionalValue<TextDecorationThicknessValue>
  textEmphasis?: ConditionalValue<PropertyValueMap["textEmphasis"]>
  textEmphasisColor?: ConditionalValue<ColorsValue>
  textEmphasisPosition?: ConditionalValue<PropertyValueMap["textEmphasisPosition"]>
  textEmphasisStyle?: ConditionalValue<PropertyValueMap["textEmphasisStyle"]>
  textIndent?: ConditionalValue<SpacingValue>
  textJustify?: ConditionalValue<PropertyValueMap["textJustify"] | AnyString>
  textOrientation?: ConditionalValue<PropertyValueMap["textOrientation"] | AnyString>
  textOverflow?: ConditionalValue<TextOverflowValue>
  textRendering?: ConditionalValue<PropertyValueMap["textRendering"] | AnyString>
  textShadow?: ConditionalValue<ShadowsValue | WithEscapeHatch<"none">>
  textSizeAdjust?: ConditionalValue<TextSizeAdjustValue>
  textSpacingTrim?: ConditionalValue<CssAny>
  textTransform?: ConditionalValue<TextTransformValue>
  textUnderlineOffset?: ConditionalValue<TextUnderlineOffsetValue>
  textUnderlinePosition?: ConditionalValue<PropertyValueMap["textUnderlinePosition"]>
  textWrap?: ConditionalValue<TextWrapValue>
  textWrapMode?: ConditionalValue<PropertyValueMap["textWrapMode"] | AnyString>
  textWrapStyle?: ConditionalValue<PropertyValueMap["textWrapStyle"] | AnyString>
  timelineScope?: ConditionalValue<CssAny>
  top?: ConditionalValue<SpacingValue>
  touchAction?: ConditionalValue<TouchActionValue>
  transform?: ConditionalValue<TransformValue>
  transformBox?: ConditionalValue<TransformBoxValue>
  transformOrigin?: ConditionalValue<TransformOriginValue>
  transformStyle?: ConditionalValue<TransformStyleValue>
  transition?: ConditionalValue<TransitionValue>
  transitionBehavior?: ConditionalValue<PropertyValueMap["transitionBehavior"]>
  transitionDelay?: ConditionalValue<DurationsValue>
  transitionDuration?: ConditionalValue<DurationsValue>
  transitionProperty?: ConditionalValue<TransitionPropertyValue | WithEscapeHatch<"all" | "none"> | AnyString>
  transitionTimingFunction?: ConditionalValue<EasingsValue>
  translate?: ConditionalValue<TranslateValue | WithEscapeHatch<"none">>
  unicodeBidi?: ConditionalValue<PropertyValueMap["unicodeBidi"] | AnyString>
  userSelect?: ConditionalValue<UserSelectValue>
  vectorEffect?: ConditionalValue<PropertyValueMap["vectorEffect"] | AnyString>
  verticalAlign?: ConditionalValue<VerticalAlignValue>
  viewTimeline?: ConditionalValue<CssAny>
  viewTimelineAxis?: ConditionalValue<CssAny>
  viewTimelineInset?: ConditionalValue<PropertyValueMap["viewTimelineInset"]>
  viewTimelineName?: ConditionalValue<CssAny>
  viewTransitionClass?: ConditionalValue<CssAny>
  viewTransitionName?: ConditionalValue<CssAny>
  visibility?: ConditionalValue<VisibilityValue>
  whiteSpace?: ConditionalValue<PropertyValueMap["whiteSpace"]>
  whiteSpaceCollapse?: ConditionalValue<PropertyValueMap["whiteSpaceCollapse"] | AnyString>
  widows?: ConditionalValue<PropertyValueMap["widows"]>
  width?: ConditionalValue<WidthValue>
  willChange?: ConditionalValue<CssAny>
  wordBreak?: ConditionalValue<WordBreakValue>
  wordSpacing?: ConditionalValue<PropertyValueMap["wordSpacing"]>
  wordWrap?: ConditionalValue<PropertyValueMap["wordWrap"] | AnyString>
  writingMode?: ConditionalValue<WithEscapeHatch<PropertyValueMap["writingMode"]>>
  x?: ConditionalValue<TranslateXValue>
  y?: ConditionalValue<TranslateYValue>
  zIndex?: ConditionalValue<ZIndexValue>
  zoom?: ConditionalValue<PropertyValueMap["zoom"]>
  animationState?: ConditionalValue<AnimationStateValue>
  backdropBlur?: ConditionalValue<BlursValue>
  backdropBrightness?: ConditionalValue<BackdropBrightnessValue>
  backdropContrast?: ConditionalValue<BackdropContrastValue>
  backdropGrayscale?: ConditionalValue<BackdropGrayscaleValue>
  backdropHueRotate?: ConditionalValue<BackdropHueRotateValue>
  backdropInvert?: ConditionalValue<BackdropInvertValue>
  backdropOpacity?: ConditionalValue<BackdropOpacityValue>
  backdropSaturate?: ConditionalValue<BackdropSaturateValue>
  backdropSepia?: ConditionalValue<BackdropSepiaValue>
  backgroundConic?: ConditionalValue<BackgroundConicValue>
  backgroundGradient?: ConditionalValue<BackgroundGradientValue>
  backgroundLinear?: ConditionalValue<BackgroundLinearValue>
  backgroundRadial?: ConditionalValue<GradientsValue>
  bg?: ConditionalValue<ColorsValue>
  bgAttachment?: ConditionalValue<BackgroundAttachmentValue>
  bgBlendMode?: ConditionalValue<BackgroundBlendModeValue>
  bgClip?: ConditionalValue<BackgroundClipValue>
  bgColor?: ConditionalValue<ColorsValue>
  bgConic?: ConditionalValue<BackgroundConicValue>
  bgGradient?: ConditionalValue<BackgroundGradientValue>
  bgImage?: ConditionalValue<AssetsValue>
  bgLinear?: ConditionalValue<BackgroundLinearValue>
  bgOrigin?: ConditionalValue<BackgroundOriginValue>
  bgPosition?: ConditionalValue<BackgroundPositionValue>
  bgPositionX?: ConditionalValue<BackgroundPositionXValue>
  bgPositionY?: ConditionalValue<BackgroundPositionYValue>
  bgRadial?: ConditionalValue<GradientsValue>
  bgRepeat?: ConditionalValue<BackgroundRepeatValue>
  bgSize?: ConditionalValue<BackgroundSizeValue>
  blur?: ConditionalValue<BlursValue>
  borderBottomRadius?: ConditionalValue<RadiiValue>
  borderEnd?: ConditionalValue<BordersValue>
  borderEndColor?: ConditionalValue<ColorsValue>
  borderEndRadius?: ConditionalValue<RadiiValue>
  borderEndWidth?: ConditionalValue<BorderWidthsValue>
  borderLeftRadius?: ConditionalValue<RadiiValue>
  borderRightRadius?: ConditionalValue<RadiiValue>
  borderSpacingX?: ConditionalValue<SpacingValue>
  borderSpacingY?: ConditionalValue<SpacingValue>
  borderStart?: ConditionalValue<BordersValue>
  borderStartColor?: ConditionalValue<ColorsValue>
  borderStartRadius?: ConditionalValue<RadiiValue>
  borderStartWidth?: ConditionalValue<BorderWidthsValue>
  borderTopRadius?: ConditionalValue<RadiiValue>
  borderX?: ConditionalValue<BordersValue>
  borderXColor?: ConditionalValue<ColorsValue>
  borderXWidth?: ConditionalValue<BorderWidthsValue>
  borderY?: ConditionalValue<BordersValue>
  borderYColor?: ConditionalValue<ColorsValue>
  borderYWidth?: ConditionalValue<BorderWidthsValue>
  boxShadowColor?: ConditionalValue<ColorsValue>
  boxSize?: ConditionalValue<BoxSizeValue>
  brightness?: ConditionalValue<BrightnessValue>
  colorPalette?: ConditionalValue<ColorPaletteValue>
  contrast?: ConditionalValue<ContrastValue>
  debug?: ConditionalValue<DebugValue>
  divideColor?: ConditionalValue<ColorsValue>
  divideStyle?: ConditionalValue<BorderStyleValue>
  divideX?: ConditionalValue<BorderWidthsValue>
  divideY?: ConditionalValue<BorderWidthsValue>
  dropShadow?: ConditionalValue<DropShadowsValue>
  end?: ConditionalValue<SpacingValue>
  flexDir?: ConditionalValue<FlexDirectionValue>
  focusRing?: ConditionalValue<FocusRingValue>
  focusRingColor?: ConditionalValue<ColorsValue>
  focusRingOffset?: ConditionalValue<SpacingValue>
  focusRingStyle?: ConditionalValue<BorderStylesValue>
  focusRingWidth?: ConditionalValue<BorderWidthsValue>
  focusVisibleRing?: ConditionalValue<FocusVisibleRingValue>
  fontSmoothing?: ConditionalValue<FontSmoothingValue>
  gradientFrom?: ConditionalValue<ColorsValue>
  gradientFromPosition?: ConditionalValue<GradientFromPositionValue>
  gradientTo?: ConditionalValue<ColorsValue>
  gradientToPosition?: ConditionalValue<GradientToPositionValue>
  gradientVia?: ConditionalValue<ColorsValue>
  gradientViaPosition?: ConditionalValue<GradientViaPositionValue>
  grayscale?: ConditionalValue<GrayscaleValue>
  h?: ConditionalValue<HeightValue>
  hideBelow?: ConditionalValue<BreakpointsValue>
  hideFrom?: ConditionalValue<BreakpointsValue>
  hueRotate?: ConditionalValue<HueRotateValue>
  insetEnd?: ConditionalValue<SpacingValue>
  insetStart?: ConditionalValue<SpacingValue>
  insetX?: ConditionalValue<SpacingValue>
  insetY?: ConditionalValue<SpacingValue>
  invert?: ConditionalValue<InvertValue>
  m?: ConditionalValue<SpacingValue>
  marginEnd?: ConditionalValue<SpacingValue>
  marginStart?: ConditionalValue<SpacingValue>
  marginX?: ConditionalValue<SpacingValue>
  marginY?: ConditionalValue<SpacingValue>
  maskBottomFrom?: ConditionalValue<SpacingValue>
  maskBottomFromColor?: ConditionalValue<ColorsValue>
  maskBottomTo?: ConditionalValue<SpacingValue>
  maskBottomToColor?: ConditionalValue<ColorsValue>
  maskConic?: ConditionalValue<MaskConicValue>
  maskConicFrom?: ConditionalValue<SpacingValue>
  maskConicFromColor?: ConditionalValue<ColorsValue>
  maskConicTo?: ConditionalValue<SpacingValue>
  maskConicToColor?: ConditionalValue<ColorsValue>
  maskLeftFrom?: ConditionalValue<SpacingValue>
  maskLeftFromColor?: ConditionalValue<ColorsValue>
  maskLeftTo?: ConditionalValue<SpacingValue>
  maskLeftToColor?: ConditionalValue<ColorsValue>
  maskLinear?: ConditionalValue<MaskLinearValue>
  maskLinearFrom?: ConditionalValue<SpacingValue>
  maskLinearFromColor?: ConditionalValue<ColorsValue>
  maskLinearTo?: ConditionalValue<SpacingValue>
  maskLinearToColor?: ConditionalValue<ColorsValue>
  maskRadial?: ConditionalValue<MaskRadialValue>
  maskRadialAt?: ConditionalValue<MaskRadialAtValue>
  maskRadialFrom?: ConditionalValue<SpacingValue>
  maskRadialFromColor?: ConditionalValue<ColorsValue>
  maskRadialShape?: ConditionalValue<MaskRadialShapeValue>
  maskRadialSize?: ConditionalValue<MaskRadialSizeValue>
  maskRadialTo?: ConditionalValue<SpacingValue>
  maskRadialToColor?: ConditionalValue<ColorsValue>
  maskRightFrom?: ConditionalValue<SpacingValue>
  maskRightFromColor?: ConditionalValue<ColorsValue>
  maskRightTo?: ConditionalValue<SpacingValue>
  maskRightToColor?: ConditionalValue<ColorsValue>
  maskTopFrom?: ConditionalValue<SpacingValue>
  maskTopFromColor?: ConditionalValue<ColorsValue>
  maskTopTo?: ConditionalValue<SpacingValue>
  maskTopToColor?: ConditionalValue<ColorsValue>
  maskXFrom?: ConditionalValue<SpacingValue>
  maskXFromColor?: ConditionalValue<ColorsValue>
  maskXTo?: ConditionalValue<SpacingValue>
  maskXToColor?: ConditionalValue<ColorsValue>
  maskYFrom?: ConditionalValue<SpacingValue>
  maskYFromColor?: ConditionalValue<ColorsValue>
  maskYTo?: ConditionalValue<SpacingValue>
  maskYToColor?: ConditionalValue<ColorsValue>
  maxH?: ConditionalValue<MaxHeightValue | WithEscapeHatch<"none">>
  maxW?: ConditionalValue<MaxWidthValue | WithEscapeHatch<"none">>
  mb?: ConditionalValue<SpacingValue>
  me?: ConditionalValue<SpacingValue>
  minH?: ConditionalValue<MinHeightValue>
  minW?: ConditionalValue<MinWidthValue>
  ml?: ConditionalValue<SpacingValue>
  mr?: ConditionalValue<SpacingValue>
  ms?: ConditionalValue<SpacingValue>
  mt?: ConditionalValue<SpacingValue>
  mx?: ConditionalValue<SpacingValue>
  my?: ConditionalValue<SpacingValue>
  p?: ConditionalValue<SpacingValue>
  paddingEnd?: ConditionalValue<SpacingValue>
  paddingStart?: ConditionalValue<SpacingValue>
  paddingX?: ConditionalValue<SpacingValue>
  paddingY?: ConditionalValue<SpacingValue>
  pb?: ConditionalValue<SpacingValue>
  pe?: ConditionalValue<SpacingValue>
  pl?: ConditionalValue<SpacingValue>
  pos?: ConditionalValue<PositionValue>
  pr?: ConditionalValue<SpacingValue>
  ps?: ConditionalValue<SpacingValue>
  pt?: ConditionalValue<SpacingValue>
  px?: ConditionalValue<SpacingValue>
  py?: ConditionalValue<SpacingValue>
  ring?: ConditionalValue<BordersValue>
  ringColor?: ConditionalValue<ColorsValue>
  ringOffset?: ConditionalValue<SpacingValue>
  ringWidth?: ConditionalValue<BorderWidthsValue>
  rotateX?: ConditionalValue<RotateValue>
  rotateY?: ConditionalValue<RotateValue>
  rotateZ?: ConditionalValue<RotateValue>
  rounded?: ConditionalValue<RadiiValue>
  roundedBottom?: ConditionalValue<RadiiValue>
  roundedBottomLeft?: ConditionalValue<RadiiValue>
  roundedBottomRight?: ConditionalValue<RadiiValue>
  roundedEnd?: ConditionalValue<RadiiValue>
  roundedEndEnd?: ConditionalValue<RadiiValue>
  roundedEndStart?: ConditionalValue<RadiiValue>
  roundedLeft?: ConditionalValue<RadiiValue>
  roundedRight?: ConditionalValue<RadiiValue>
  roundedStart?: ConditionalValue<RadiiValue>
  roundedStartEnd?: ConditionalValue<RadiiValue>
  roundedStartStart?: ConditionalValue<RadiiValue>
  roundedTop?: ConditionalValue<RadiiValue>
  roundedTopLeft?: ConditionalValue<RadiiValue>
  roundedTopRight?: ConditionalValue<RadiiValue>
  saturate?: ConditionalValue<SaturateValue>
  scaleX?: ConditionalValue<ScaleXValue>
  scaleY?: ConditionalValue<ScaleYValue>
  scrollMarginX?: ConditionalValue<SpacingValue>
  scrollMarginY?: ConditionalValue<SpacingValue>
  scrollPaddingX?: ConditionalValue<SpacingValue>
  scrollPaddingY?: ConditionalValue<SpacingValue>
  scrollSnapMargin?: ConditionalValue<SpacingValue>
  scrollSnapMarginBottom?: ConditionalValue<SpacingValue>
  scrollSnapMarginLeft?: ConditionalValue<SpacingValue>
  scrollSnapMarginRight?: ConditionalValue<SpacingValue>
  scrollSnapMarginTop?: ConditionalValue<SpacingValue>
  scrollSnapStrictness?: ConditionalValue<ScrollSnapStrictnessValue>
  scrollbar?: ConditionalValue<ScrollbarValue>
  scrollbarThumb?: ConditionalValue<ColorsValue>
  scrollbarTrack?: ConditionalValue<ColorsValue>
  sepia?: ConditionalValue<SepiaValue>
  shadow?: ConditionalValue<ShadowsValue | WithEscapeHatch<"none">>
  shadowColor?: ConditionalValue<ColorsValue>
  spaceX?: ConditionalValue<SpacingValue>
  spaceY?: ConditionalValue<SpacingValue>
  srOnly?: ConditionalValue<SrOnlyValue>
  start?: ConditionalValue<SpacingValue>
  textGradient?: ConditionalValue<TextGradientValue>
  textShadowColor?: ConditionalValue<ColorsValue>
  textStyle?: ConditionalValue<TextStyleValue>
  translateX?: ConditionalValue<TranslateXValue>
  translateY?: ConditionalValue<TranslateYValue>
  translateZ?: ConditionalValue<SpacingValue>
  truncate?: ConditionalValue<TruncateValue>
  w?: ConditionalValue<WidthValue>
  z?: ConditionalValue<SpacingValue>
}

export type CssVarValue = ConditionalValue<CssVars | AnyString | AnyNumber>

export type CssVarProperties = {
  [K in `--${string}`]?: CssVarValue
}

export type NestedStyles = {
  [K in Selector | Condition]?: SystemStyleObject
}

export interface SystemStyleObject extends SystemProperties, CssVarProperties, NestedStyles {}

export type SystemStyleObjectWith<K extends keyof SystemProperties> = Pick<SystemProperties, K> & {
  [S in Selector | Condition]?: SystemStyleObjectWith<K>
}

export type SystemStyleObjectWithout<K extends keyof SystemProperties> = SystemStyleObjectWith<Exclude<keyof SystemProperties, K>>

export interface GlobalStyleObject {
  [selector: string]: SystemStyleObject
}

export interface CssKeyframes {
  [name: string]: {
    [time: string]: SystemStyleObject
  }
}

export interface GlobalFontfaceRule {
  fontFamily: string
  src: string
  fontStyle?: string
  fontWeight?: string | number
  fontDisplay?: "auto" | "block" | "swap" | "fallback" | "optional"
  unicodeRange?: string
  fontFeatureSettings?: string
  fontVariationSettings?: string
  fontStretch?: string
  ascentOverride?: string
  descentOverride?: string
  lineGapOverride?: string
  sizeAdjust?: string
}

export type FontfaceRule = Omit<GlobalFontfaceRule, "fontFamily">

export interface GlobalFontface {
  [name: string]: FontfaceRule | FontfaceRule[]
}

interface WithCss {
  css?: SystemStyleObject | SystemStyleObject[]
}

export type JsxStyleProps = SystemStyleObject & WithCss

export interface PatchedHTMLProps {
  htmlWidth?: string | number
  htmlHeight?: string | number
  htmlTranslate?: "yes" | "no" | undefined
  htmlContent?: string
}

export type OmittedHTMLProps = "color" | "translate" | "transition" | "width" | "height" | "content"

type WithHTMLProps<T> = DistributiveOmit<T, OmittedHTMLProps> & PatchedHTMLProps

export type JsxHTMLProps<T extends Record<string, any>, P extends Record<string, any> = {}> = Assign<WithHTMLProps<T>, P>