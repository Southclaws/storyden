import { createCssRuntime, hypenateProperty, isObject, withoutSpace } from '../helpers';
import { breakpointKeys, finalizeConditions, sortConditions } from './conditions';

const utilities = "WebkitTextFillColor:wktf-c,accentColor:ac-c,alignContent:ac,alignItems:ai,alignSelf:as,animation:anim,animationComposition:anim-comp,animationDelay:anim-dly,animationDirection:anim-dir,animationDuration:anim-dur,animationFillMode:anim-fm,animationIterationCount:anim-ic,animationName:anim-n,animationPlayState:anim-ps,animationRange:anim-r,animationRangeEnd:anim-re,animationRangeStart:anim-rs,animationState:anim-s,animationTimeline:anim-tl,animationTimingFunction:anim-tmf,appearance:ap,aspectRatio:asp,backdropBlur:bkdp-blur,backdropBrightness:bkdp-brightness,backdropContrast:bkdp-contrast,backdropFilter:bkdp,backdropGrayscale:bkdp-grayscale,backdropHueRotate:bkdp-hue-rotate,backdropInvert:bkdp-invert,backdropOpacity:bkdp-opacity,backdropSaturate:bkdp-saturate,backdropSepia:bkdp-sepia,backfaceVisibility:bfv,background:bg/1,backgroundAttachment:bg-a/bgAttachment,backgroundBlendMode:bg-bm/bgBlendMode,backgroundClip:bg-cp/bgClip,backgroundColor:bg-c/bgColor,backgroundConic:bg-conic/bgConic,backgroundGradient:bg-grad/bgGradient,backgroundImage:bg-i/bgImage,backgroundLinear:bg-linear/bgLinear,backgroundOrigin:bg-o/bgOrigin,backgroundPosition:bg-p/bgPosition,backgroundPositionX:bg-p-x/bgPositionX,backgroundPositionY:bg-p-y/bgPositionY,backgroundRadial:bg-radial/bgRadial,backgroundRepeat:bg-r/bgRepeat,backgroundSize:bg-s/bgSize,blockSize:h-bs,border:bd,borderBlock:bd-y/borderY,borderBlockColor:bd-y-c/borderYColor,borderBlockEnd:bd-be,borderBlockEndColor:bd-be-c,borderBlockEndWidth:bd-be-w,borderBlockStart:bd-bs,borderBlockStartColor:bd-bs-c,borderBlockStartWidth:bd-bs-w,borderBlockWidth:bd-y-w/borderYWidth,borderBottom:bd-b,borderBottomColor:bd-b-c,borderBottomLeftRadius:bdr-bl/roundedBottomLeft,borderBottomRadius:bdr-b/roundedBottom,borderBottomRightRadius:bdr-br/roundedBottomRight,borderBottomWidth:bd-b-w,borderCollapse:bd-cl,borderColor:bd-c,borderEndEndRadius:bdr-ee/roundedEndEnd,borderEndRadius:bdr-e/roundedEnd,borderEndStartRadius:bdr-es/roundedEndStart,borderInline:bd-x/borderX,borderInlineColor:bd-x-c/borderXColor,borderInlineEnd:bd-e/borderEnd,borderInlineEndColor:bd-e-c/borderEndColor,borderInlineEndWidth:bd-e-w/borderEndWidth,borderInlineStart:bd-s/borderStart,borderInlineStartColor:bd-s-c/borderStartColor,borderInlineStartWidth:bd-s-w/borderStartWidth,borderInlineWidth:bd-x-w/borderXWidth,borderLeft:bd-l,borderLeftColor:bd-l-c,borderLeftRadius:bdr-l/roundedLeft,borderLeftWidth:bd-l-w,borderRadius:bdr/rounded,borderRight:bd-r,borderRightColor:bd-r-c,borderRightRadius:bdr-r/roundedRight,borderRightWidth:bd-r-w,borderSpacing:bd-sp,borderSpacingX:bd-sx,borderSpacingY:bd-sy,borderStartEndRadius:bdr-se/roundedStartEnd,borderStartRadius:bdr-s/roundedStart,borderStartStartRadius:bdr-ss/roundedStartStart,borderTop:bd-t,borderTopColor:bd-t-c,borderTopLeftRadius:bdr-tl/roundedTopLeft,borderTopRadius:bdr-t/roundedTop,borderTopRightRadius:bdr-tr/roundedTopRight,borderTopWidth:bd-t-w,borderWidth:bd-w,boxDecorationBreak:bx-db,boxShadow:bx-sh/shadow,boxShadowColor:bx-sh-c/shadowColor,boxSize:size,boxSizing:bx-s,caretColor:ca-c,clipPath:cp-path,color:c,columnGap:cg,container:cq,containerName:cq-n,containerType:cq-t,display:d,divideColor:dvd-c,divideStyle:dvd-s,divideX:dvd-x,divideY:dvd-y,flexBasis:flex-b,flexDirection:flex-d/flexDir,flexGrow:flex-g,flexShrink:flex-sh,focusRingColor:focus-ring-c,focusRingOffset:focus-ring-o,focusRingStyle:focus-ring-s,focusRingWidth:focus-ring-w,focusVisibleRing:focus-v-ring,fontFamily:ff,fontFeatureSettings:ff-s,fontKerning:fk,fontPalette:fp,fontSize:fs,fontSizeAdjust:fs-a,fontSmoothing:fsmt,fontVariant:fv,fontVariantAlternates:fv-alt,fontVariantCaps:fv-caps,fontVariantNumeric:fv-num,fontVariationSettings:fv-s,fontWeight:fw,gradientFrom:grad-from,gradientFromPosition:grad-from-pos,gradientTo:grad-to,gradientToPosition:grad-to-pos,gradientVia:grad-via,gradientViaPosition:grad-via-pos,gridAutoColumns:grid-ac,gridAutoFlow:grid-af,gridAutoRows:grid-ar,gridColumn:grid-c,gridColumnEnd:grid-ce,gridColumnGap:grid-cg,gridColumnStart:grid-cs,gridGap:grid-g,gridRow:grid-r,gridRowGap:grid-rg,gridTemplateColumns:grid-tc,gridTemplateRows:grid-tr,height:h/1,hideBelow:show,hideFrom:hide,hyphens:hy,inlineSize:w-is,insetBlock:inset-y/insetY,insetBlockEnd:inset-be,insetBlockStart:inset-bs,insetInline:inset-x/insetX,insetInlineEnd:inset-e/end/insetEnd,insetInlineStart:inset-s/insetStart/start,justifyContent:jc,letterSpacing:ls,lineClamp:lc,lineHeight:lh,listStyle:li-s,listStyleImage:li-img,listStylePosition:li-pos,listStyleType:li-t,margin:m/1,marginBlock:my/marginY/1,marginBlockEnd:mbe,marginBlockStart:mbs,marginBottom:mb/1,marginInline:mx/marginX/1,marginInlineEnd:me/marginEnd/1,marginInlineStart:ms/marginStart/1,marginLeft:ml/1,marginRight:mr/1,marginTop:mt/1,mask:msk,maskImage:msk-i,maskSize:msk-s,maxBlockSize:max-b,maxHeight:max-h/maxH,maxInlineSize:max-w-is,maxWidth:max-w/maxW,minBlockSize:min-h-bs,minHeight:min-h/minH,minInlineSize:min-w-is,minWidth:min-w/minW,mixBlendMode:mix-bm,objectFit:obj-f,objectPosition:obj-p,opacity:op,outline:ring/1,outlineColor:ring-c/ringColor,outlineOffset:ring-o/ringOffset,outlineWidth:ring-w/ringWidth,overflow:ov,overflowAnchor:ov-a,overflowBlock:ov-b,overflowClipBox:ovcp-bx,overflowClipMargin:ovcp-m,overflowInline:ov-i,overflowWrap:ov-wrap,overflowX:ov-x,overflowY:ov-y,overscrollBehavior:ovs-b,overscrollBehaviorBlock:ovs-bb,overscrollBehaviorInline:ovs-bi,overscrollBehaviorX:ovs-bx,overscrollBehaviorY:ovs-by,padding:p/1,paddingBlock:py/paddingY/1,paddingBlockEnd:pbe,paddingBlockStart:pbs,paddingBottom:pb/1,paddingInline:px/paddingX/1,paddingInlineEnd:pe/paddingEnd/1,paddingInlineStart:ps/paddingStart/1,paddingLeft:pl/1,paddingRight:pr/1,paddingTop:pt/1,position:pos/1,rowGap:rg,scrollBehavior:scr-bhv,scrollMargin:scr-m,scrollMarginBlock:scr-my/scrollMarginY,scrollMarginBlockEnd:scr-mbe,scrollMarginBlockStart:scr-mbt,scrollMarginBottom:scr-mb,scrollMarginInline:scr-mx/scrollMarginX,scrollMarginInlineEnd:scr-me,scrollMarginInlineStart:scr-ms,scrollMarginLeft:scr-ml,scrollMarginRight:scr-mr,scrollMarginTop:scr-mt,scrollPadding:scr-p,scrollPaddingBlock:scr-py/scrollPaddingY,scrollPaddingBlockEnd:scr-pbe,scrollPaddingBlockStart:scr-pbs,scrollPaddingBottom:scr-pb,scrollPaddingInline:scr-px/scrollPaddingX,scrollPaddingInlineEnd:scr-pe,scrollPaddingInlineStart:scr-ps,scrollPaddingLeft:scr-pl,scrollPaddingRight:scr-pr,scrollPaddingTop:scr-pt,scrollSnapAlign:scr-sa,scrollSnapCoordinate:scrs-c,scrollSnapDestination:scrs-d,scrollSnapMargin:scrs-m,scrollSnapMarginBottom:scrs-mb,scrollSnapMarginLeft:scrs-ml,scrollSnapMarginRight:scrs-mr,scrollSnapMarginTop:scrs-mt,scrollSnapPointsX:scrs-px,scrollSnapPointsY:scrs-py,scrollSnapStop:scrs-s,scrollSnapStrictness:scrs-strt,scrollSnapType:scrs-t,scrollSnapTypeX:scrs-tx,scrollSnapTypeY:scrs-ty,scrollTimeline:scrtl,scrollTimelineAxis:scrtl-a,scrollTimelineName:scrtl-n,scrollbar:scr-bar,scrollbarColor:scr-bar-c,scrollbarGutter:scr-bar-g,scrollbarWidth:scr-bar-w,spaceX:sx,spaceY:sy,srOnly:sr,stroke:stk,strokeDasharray:stk-dsh,strokeDashoffset:stk-do,strokeLinecap:stk-lc,strokeLinejoin:stk-lj,strokeMiterlimit:stk-ml,strokeOpacity:stk-op,strokeWidth:stk-w,tableLayout:tbl,textAlign:ta,textDecoration:td,textDecorationColor:td-c,textDecorationStyle:td-s,textDecorationThickness:td-t,textEmphasisColor:te-c,textGradient:txt-grad,textIndent:ti,textOverflow:tov,textShadow:tsh,textShadowColor:tsh-c/textShadowColor,textSizeAdjust:txt-adj,textStyle:textStyle,textTransform:tt,textUnderlineOffset:tu-o,textWrap:tw,touchAction:tch-a,transform:trf,transformBox:trf-b,transformOrigin:trf-o,transformStyle:trf-s,transition:trs,transitionDelay:trs-dly,transitionDuration:trs-dur,transitionProperty:trs-prop,transitionTimingFunction:trs-tmf,translateX:/x,translateY:/y,translateZ:/z,truncate:trunc,userSelect:us,verticalAlign:va,visibility:vis,width:w/1,wordBreak:wb,zIndex:z"

const classNameByProp = new Map()
const shorthands = new Map()
if (utilities) {
  utilities.split(",").forEach((utility) => {
    const [prop, meta] = utility.split(":")
    const [className, ...shorthandList] = meta.split("/")
    if (className) classNameByProp.set(prop, className)
    shorthandList.forEach((shorthand) => {
      const key = shorthand === "1" ? className : shorthand
      shorthands.set(key, prop)
    })
  })
}

const resolveShorthand = (prop) => shorthands.get(prop) || prop

const { serializeCss, mergeCss, assignCss } = createCssRuntime({
  hash: false,
  conditions: {
    shift: sortConditions,
    finalize: finalizeConditions,
    breakpoints: { keys: breakpointKeys },
  },
  utility: {
    prefix: null,
    hasShorthand: true,
    toHash(path, hashFn) {
      return hashFn(path.join(":"))
    },
    transform(prop, value) {
      const key = resolveShorthand(prop)
      const propKey = classNameByProp.get(key) || hypenateProperty(key)
      return { className: `${propKey}_${withoutSpace(value)}` }
    },
    resolveShorthand,
  },
})

export const css = /* @__PURE__ */ Object.assign(
  function css(...styles) {
    if (styles.length === 1 && isObject(styles[0])) return serializeCss(styles[0])
    return serializeCss(mergeCss(...styles))
  },
  {
    raw: function cssRaw(...styles) {
      return mergeCss(...styles)
    },
  },
)

export { mergeCss, assignCss }