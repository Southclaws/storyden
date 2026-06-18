import { createRecipe } from './runtime';

const groupConfig = {"name":"group","className":"ui-group","defaultVariants":{"orientation":"horizontal"},"compoundVariants":[{"attached":true,"orientation":"horizontal","className":"ui-group--compound__attached_true__orientation_horizontal"},{"attached":true,"orientation":"vertical","className":"ui-group--compound__attached_true__orientation_vertical"}],"variantMap":{"attached":["true"],"grow":["true"],"orientation":["horizontal","vertical"]}}

export const group = /* @__PURE__ */ createRecipe(groupConfig)