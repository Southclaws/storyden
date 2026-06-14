import { memo, splitProps } from '../helpers.mjs';
import { createRecipe, mergeRecipes } from './create-recipe.mjs';

const groupFn = /* @__PURE__ */ createRecipe('ui-group', {
  "orientation": "horizontal"
}, [
  {
    "orientation": "horizontal",
    "attached": true,
    "css": {
      "& > *:first-child:not(:only-child)": {
        "borderEndRadius": "0",
        "marginEnd": "-1px"
      },
      "& > *:last-child:not(:only-child)": {
        "borderStartRadius": "0"
      },
      "& > *:not(:first-child):not(:last-child)": {
        "borderRadius": "0",
        "marginEnd": "-1px"
      }
    }
  },
  {
    "orientation": "vertical",
    "attached": true,
    "css": {
      "& > *:first-child:not(:only-child)": {
        "borderBottomRadius": "0",
        "marginBottom": "-1px"
      },
      "& > *:last-child:not(:only-child)": {
        "borderTopRadius": "0"
      },
      "& > *:not(:first-child):not(:last-child)": {
        "borderRadius": "0",
        "marginBottom": "-1px"
      }
    }
  }
])

const groupVariantMap = {
  "orientation": [
    "horizontal",
    "vertical"
  ],
  "attached": [
    "true"
  ],
  "grow": [
    "true"
  ]
}

const groupVariantKeys = Object.keys(groupVariantMap)

export const group = /* @__PURE__ */ Object.assign(memo(groupFn.recipeFn), {
  __recipe__: true,
  __name__: 'group',
  __getCompoundVariantCss__: groupFn.__getCompoundVariantCss__,
  raw: (props) => props,
  variantKeys: groupVariantKeys,
  variantMap: groupVariantMap,
  merge(recipe) {
    return mergeRecipes(this, recipe)
  },
  splitVariantProps(props) {
    return splitProps(props, groupVariantKeys)
  },
  getVariantProps: groupFn.getVariantProps,
})