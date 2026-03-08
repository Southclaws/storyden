import { withoutSpace } from '../helpers.mjs';

const conditionsStr = "_target,_checked,_indeterminate,_closed,_open,_detailsOpen,_on,_off,_hidden,_current,_today,_placeholderShown,_collapsed,_containerSmall,_containerMedium,_containerLarge,_hover,_focus,_focusWithin,_focusVisible,_disabled,_active,_visited,_readOnly,_readWrite,_empty,_enabled,_expanded,_highlighted,_complete,_incomplete,_dragging,_before,_after,_firstLetter,_firstLine,_marker,_selection,_file,_backdrop,_first,_last,_only,_even,_odd,_firstOfType,_lastOfType,_onlyOfType,_peerFocus,_peerHover,_peerActive,_peerFocusWithin,_peerFocusVisible,_peerDisabled,_peerChecked,_peerInvalid,_peerExpanded,_peerPlaceholderShown,_groupFocus,_groupHover,_groupActive,_groupFocusWithin,_groupFocusVisible,_groupDisabled,_groupChecked,_groupExpanded,_groupInvalid,_required,_valid,_invalid,_autofill,_inRange,_outOfRange,_placeholder,_pressed,_selected,_grabbed,_underValue,_overValue,_atValue,_default,_optional,_fullscreen,_loading,_currentPage,_currentStep,_unavailable,_rangeStart,_rangeEnd,_now,_topmost,_motionReduce,_motionSafe,_print,_landscape,_portrait,_dark,_light,_osDark,_osLight,_highContrast,_lessContrast,_moreContrast,_ltr,_rtl,_scrollbar,_scrollbarThumb,_scrollbarTrack,_horizontal,_vertical,_icon,_starting,_noscript,_invertedColors,sm,smOnly,smDown,md,mdOnly,mdDown,lg,lgOnly,lgDown,xl,xlOnly,xlDown,2xl,2xlOnly,2xlDown,smToMd,smToLg,smToXl,smTo2xl,mdToLg,mdToXl,mdTo2xl,lgToXl,lgTo2xl,xlTo2xl,base"
const conditions = new Set(conditionsStr.split(','))

const conditionRegex = /^@|&|&$/

export function isCondition(value){
  return conditions.has(value) || conditionRegex.test(value)
}

const underscoreRegex = /^_/
const conditionsSelectorRegex = /&|@/

export function finalizeConditions(paths){
  return paths.map((path) => {
    if (conditions.has(path)){
      return path.replace(underscoreRegex, '')
    }

    if (conditionsSelectorRegex.test(path)){
      return `[${withoutSpace(path.trim())}]`
    }

    return path
  })}

  export function sortConditions(paths){
    return paths.sort((a, b) => {
      const aa = isCondition(a)
      const bb = isCondition(b)
      if (aa && !bb) return 1
      if (!aa && bb) return -1
      return 0
    })
  }