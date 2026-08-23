import { recipe } from '@vanilla-extract/recipes'
import { vars } from '#/styles/theme.css'

export const badgeRecipe = recipe({
  base: {
    display: 'inline-flex',
    alignItems: 'center',
    gap: vars.space['2xs'],
    fontWeight: 600,
    borderRadius: vars.radii.xs,
    fontSize: vars.fontSize['3xs'],
    lineHeight: vars.lineHeight.none,
    padding: `3px ${vars.space.xs}`,
    fontVariantNumeric: 'tabular-nums',
  },
  variants: {
    variant: {
      success: {
        backgroundColor: vars.color.successMuted,
        color: vars.color.success,
        border: `1px solid ${vars.color.successBorder}`,
      },
      danger: {
        backgroundColor: vars.color.dangerMuted,
        color: vars.color.danger,
        border: `1px solid ${vars.color.dangerBorder}`,
      },
      warning: {
        backgroundColor: vars.color.warningMuted,
        color: vars.color.warning,
        border: `1px solid ${vars.color.warningBorder}`,
      },
      teal: {
        backgroundColor: vars.color.tealMuted,
        color: vars.color.teal,
        border: `1px solid ${vars.color.tealBorder}`,
      },
      primary: {
        backgroundColor: vars.color.primaryMuted,
        color: vars.color.primary,
        border: `1px solid ${vars.color.primaryBorder}`,
      },
      neutral: {
        backgroundColor: vars.color.surfaceElevated,
        color: vars.color.textSecondary,
        border: `1px solid ${vars.color.border}`,
      },
    },
  },
  defaultVariants: {
    variant: 'neutral',
  },
})
