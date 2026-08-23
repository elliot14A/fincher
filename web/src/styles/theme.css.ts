import { createGlobalTheme } from '@vanilla-extract/css'
import { colors, fontSizes, lineHeights, radii, space } from './tokens'

export const vars = createGlobalTheme(':root', {
  color: colors,
  space,
  radii,
  fontSize: fontSizes,
  lineHeight: lineHeights,
})
