export type AnimationToken = "bounce" | "ping" | "pulse" | "spin"

export type AspectRatioToken = "golden" | "landscape" | "portrait" | "square" | "ultrawide" | "video" | "wide"

export type BlurToken = "2xl" | "3xl" | "frosted" | "lg" | "md" | "sm" | "xl" | "xs"

export type BorderToken = "none"

export type BreakpointToken = "2xl" | "lg" | "md" | "sm" | "xl"

export type ColorToken = "Mono.ink" | "Mono.slush" | "Primary.campfire" | "Primary.forest" | "Primary.moonlit" | "Primary.saddle" | "Shades.iron" | "Shades.newspaper" | "Shades.slate" | "Shades.stone" | "amber.100" | "amber.200" | "amber.300" | "amber.400" | "amber.50" | "amber.500" | "amber.600" | "amber.700" | "amber.800" | "amber.900" | "amber.950" | "black" | "blue.100" | "blue.200" | "blue.300" | "blue.400" | "blue.50" | "blue.500" | "blue.600" | "blue.700" | "blue.800" | "blue.900" | "blue.950" | "colorPalette" | "colorPalette.100" | "colorPalette.200" | "colorPalette.300" | "colorPalette.400" | "colorPalette.50" | "colorPalette.500" | "colorPalette.600" | "colorPalette.700" | "colorPalette.800" | "colorPalette.900" | "colorPalette.950" | "colorPalette.campfire" | "colorPalette.forest" | "colorPalette.ink" | "colorPalette.iron" | "colorPalette.moonlit" | "colorPalette.newspaper" | "colorPalette.saddle" | "colorPalette.slate" | "colorPalette.slush" | "colorPalette.stone" | "current" | "cyan.100" | "cyan.200" | "cyan.300" | "cyan.400" | "cyan.50" | "cyan.500" | "cyan.600" | "cyan.700" | "cyan.800" | "cyan.900" | "cyan.950" | "emerald.100" | "emerald.200" | "emerald.300" | "emerald.400" | "emerald.50" | "emerald.500" | "emerald.600" | "emerald.700" | "emerald.800" | "emerald.900" | "emerald.950" | "fuchsia.100" | "fuchsia.200" | "fuchsia.300" | "fuchsia.400" | "fuchsia.50" | "fuchsia.500" | "fuchsia.600" | "fuchsia.700" | "fuchsia.800" | "fuchsia.900" | "fuchsia.950" | "gray.100" | "gray.200" | "gray.300" | "gray.400" | "gray.50" | "gray.500" | "gray.600" | "gray.700" | "gray.800" | "gray.900" | "gray.950" | "green.100" | "green.200" | "green.300" | "green.400" | "green.50" | "green.500" | "green.600" | "green.700" | "green.800" | "green.900" | "green.950" | "indigo.100" | "indigo.200" | "indigo.300" | "indigo.400" | "indigo.50" | "indigo.500" | "indigo.600" | "indigo.700" | "indigo.800" | "indigo.900" | "indigo.950" | "lime.100" | "lime.200" | "lime.300" | "lime.400" | "lime.50" | "lime.500" | "lime.600" | "lime.700" | "lime.800" | "lime.900" | "lime.950" | "mauve.100" | "mauve.200" | "mauve.300" | "mauve.400" | "mauve.50" | "mauve.500" | "mauve.600" | "mauve.700" | "mauve.800" | "mauve.900" | "mauve.950" | "mist.100" | "mist.200" | "mist.300" | "mist.400" | "mist.50" | "mist.500" | "mist.600" | "mist.700" | "mist.800" | "mist.900" | "mist.950" | "neutral.100" | "neutral.200" | "neutral.300" | "neutral.400" | "neutral.50" | "neutral.500" | "neutral.600" | "neutral.700" | "neutral.800" | "neutral.900" | "neutral.950" | "olive.100" | "olive.200" | "olive.300" | "olive.400" | "olive.50" | "olive.500" | "olive.600" | "olive.700" | "olive.800" | "olive.900" | "olive.950" | "orange.100" | "orange.200" | "orange.300" | "orange.400" | "orange.50" | "orange.500" | "orange.600" | "orange.700" | "orange.800" | "orange.900" | "orange.950" | "pink.100" | "pink.200" | "pink.300" | "pink.400" | "pink.50" | "pink.500" | "pink.600" | "pink.700" | "pink.800" | "pink.900" | "pink.950" | "purple.100" | "purple.200" | "purple.300" | "purple.400" | "purple.50" | "purple.500" | "purple.600" | "purple.700" | "purple.800" | "purple.900" | "purple.950" | "red.100" | "red.200" | "red.300" | "red.400" | "red.50" | "red.500" | "red.600" | "red.700" | "red.800" | "red.900" | "red.950" | "rose.100" | "rose.200" | "rose.300" | "rose.400" | "rose.50" | "rose.500" | "rose.600" | "rose.700" | "rose.800" | "rose.900" | "rose.950" | "sky.100" | "sky.200" | "sky.300" | "sky.400" | "sky.50" | "sky.500" | "sky.600" | "sky.700" | "sky.800" | "sky.900" | "sky.950" | "slate.100" | "slate.200" | "slate.300" | "slate.400" | "slate.50" | "slate.500" | "slate.600" | "slate.700" | "slate.800" | "slate.900" | "slate.950" | "stone.100" | "stone.200" | "stone.300" | "stone.400" | "stone.50" | "stone.500" | "stone.600" | "stone.700" | "stone.800" | "stone.900" | "stone.950" | "taupe.100" | "taupe.200" | "taupe.300" | "taupe.400" | "taupe.50" | "taupe.500" | "taupe.600" | "taupe.700" | "taupe.800" | "taupe.900" | "taupe.950" | "teal.100" | "teal.200" | "teal.300" | "teal.400" | "teal.50" | "teal.500" | "teal.600" | "teal.700" | "teal.800" | "teal.900" | "teal.950" | "transparent" | "violet.100" | "violet.200" | "violet.300" | "violet.400" | "violet.50" | "violet.500" | "violet.600" | "violet.700" | "violet.800" | "violet.900" | "violet.950" | "white" | "yellow.100" | "yellow.200" | "yellow.300" | "yellow.400" | "yellow.50" | "yellow.500" | "yellow.600" | "yellow.700" | "yellow.800" | "yellow.900" | "yellow.950" | "zinc.100" | "zinc.200" | "zinc.300" | "zinc.400" | "zinc.50" | "zinc.500" | "zinc.600" | "zinc.700" | "zinc.800" | "zinc.900" | "zinc.950"

export type DurationToken = "fast" | "faster" | "fastest" | "normal" | "slow" | "slower" | "slowest"

export type EasingToken = "default" | "in" | "in-out" | "linear" | "out"

export type FontSizeToken = "2xl" | "2xs" | "3xl" | "4xl" | "5xl" | "6xl" | "7xl" | "8xl" | "9xl" | "lg" | "md" | "sm" | "xl" | "xs"

export type FontWeightToken = "black" | "bold" | "extrabold" | "extralight" | "light" | "medium" | "normal" | "semibold" | "thin"

export type FontToken = "gorton" | "hedvig" | "intelone" | "joie" | "mono" | "sans" | "serif" | "worksans"

export type LetterSpacingToken = "normal" | "tight" | "tighter" | "wide" | "wider" | "widest"

export type LineHeightToken = "loose" | "none" | "normal" | "relaxed" | "snug" | "tight"

export type RadiusToken = "2xl" | "3xl" | "4xl" | "full" | "lg" | "md" | "sm" | "xl" | "xs"

export type ShadowToken = "2xl" | "2xs" | "inner" | "inset-2xs" | "inset-sm" | "inset-xs" | "lg" | "md" | "sm" | "xl" | "xs"

export type SizeToken = "0" | "0.5" | "1" | "1.5" | "10" | "11" | "12" | "14" | "16" | "2" | "2.5" | "20" | "24" | "28" | "2xl" | "3" | "3.5" | "32" | "36" | "3xl" | "4" | "4.5" | "40" | "44" | "48" | "4xl" | "5" | "5.5" | "52" | "56" | "5xl" | "6" | "60" | "64" | "6xl" | "7" | "72" | "7xl" | "8" | "80" | "8xl" | "9" | "96" | "breakpoint-2xl" | "breakpoint-lg" | "breakpoint-md" | "breakpoint-sm" | "breakpoint-xl" | "fit" | "full" | "lg" | "max" | "md" | "min" | "prose" | "sm" | "xl" | "xs"

export type SpacingToken = "-0.5" | "-1" | "-1.5" | "-10" | "-11" | "-12" | "-14" | "-16" | "-2" | "-2.5" | "-20" | "-24" | "-28" | "-3" | "-3.5" | "-32" | "-36" | "-4" | "-4.5" | "-40" | "-44" | "-48" | "-5" | "-5.5" | "-52" | "-56" | "-6" | "-60" | "-64" | "-7" | "-72" | "-8" | "-80" | "-9" | "-96" | "0" | "0.5" | "1" | "1.5" | "10" | "11" | "12" | "14" | "16" | "2" | "2.5" | "20" | "24" | "28" | "3" | "3.5" | "32" | "36" | "4" | "4.5" | "40" | "44" | "48" | "5" | "5.5" | "52" | "56" | "6" | "60" | "64" | "7" | "72" | "8" | "80" | "9" | "96"

export interface Tokens {
  animations: AnimationToken
  aspectRatios: AspectRatioToken
  blurs: BlurToken
  borders: BorderToken
  breakpoints: BreakpointToken
  colors: ColorToken
  durations: DurationToken
  easings: EasingToken
  fontSizes: FontSizeToken
  fontWeights: FontWeightToken
  fonts: FontToken
  letterSpacings: LetterSpacingToken
  lineHeights: LineHeightToken
  radii: RadiusToken
  shadows: ShadowToken
  sizes: SizeToken
  spacing: SpacingToken
}

export type Token = `animations.${AnimationToken}` | `aspectRatios.${AspectRatioToken}` | `blurs.${BlurToken}` | `borders.${BorderToken}` | `breakpoints.${BreakpointToken}` | `colors.${ColorToken}` | `durations.${DurationToken}` | `easings.${EasingToken}` | `fontSizes.${FontSizeToken}` | `fontWeights.${FontWeightToken}` | `fonts.${FontToken}` | `letterSpacings.${LetterSpacingToken}` | `lineHeights.${LineHeightToken}` | `radii.${RadiusToken}` | `shadows.${ShadowToken}` | `sizes.${SizeToken}` | `spacing.${SpacingToken}`

export type ColorOpacityModifier = `${number}`

export type ColorOpacityToken = `colors.${ColorToken}/${ColorOpacityModifier}`

export type TokenPath = Token | ColorOpacityToken

export type ColorPalette = "Mono" | "Primary" | "Shades" | "amber" | "black" | "blue" | "current" | "cyan" | "emerald" | "fuchsia" | "gray" | "green" | "indigo" | "lime" | "mauve" | "mist" | "neutral" | "olive" | "orange" | "pink" | "purple" | "red" | "rose" | "sky" | "slate" | "stone" | "taupe" | "teal" | "transparent" | "violet" | "white" | "yellow" | "zinc"

export type TokenValue<T extends string> = T extends keyof Tokens ? Tokens[T] : never