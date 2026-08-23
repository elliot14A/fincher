import { recipe } from '@vanilla-extract/recipes'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

export const logoRecipe = recipe({
  base: {
    fontFamily: fonts.logo,
    fontWeight: 400,
    color: vars.color.textPrimary,
    lineHeight: vars.lineHeight.snug,
    letterSpacing: '0.01em',
  },
  variants: {
    size: {
      sm: { fontSize: '20px' },
      md: { fontSize: '26px' },
      lg: { fontSize: '30px' },
      xl: { fontSize: '36px' },
    },
  },
  defaultVariants: {
    size: 'md',
  },
})
