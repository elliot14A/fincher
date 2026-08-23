import { recipe } from '@vanilla-extract/recipes'
import { vars } from '#/styles/theme.css'

export const buttonRecipe = recipe({
  base: {
    display: 'inline-flex',
    alignItems: 'center',
    justifyContent: 'center',
    gap: vars.space.xs,
    fontWeight: 500,
    borderRadius: vars.radii.sm,
    transition: 'background-color 0.12s ease, border-color 0.12s ease, color 0.12s ease',
    outline: 'none',
    ':disabled': {
      opacity: 0.5,
      cursor: 'not-allowed',
    },
  },
  variants: {
    variant: {
      primary: {
        backgroundColor: vars.color.primary,
        color: vars.color.textInverse,
        fontWeight: 600,
        ':hover': {
          backgroundColor: vars.color.primaryHover,
        },
      },
      secondary: {
        backgroundColor: vars.color.surfaceElevated,
        color: vars.color.textPrimary,
        border: `1px solid ${vars.color.border}`,
        ':hover': {
          backgroundColor: vars.color.surfaceHover,
          borderColor: vars.color.borderStrong,
        },
      },
      danger: {
        backgroundColor: vars.color.dangerMuted,
        color: vars.color.danger,
        border: `1px solid ${vars.color.dangerBorder}`,
        ':hover': {
          backgroundColor: vars.color.danger,
          color: vars.color.textInverse,
        },
      },
      ghost: {
        backgroundColor: 'transparent',
        color: vars.color.textSecondary,
        ':hover': {
          backgroundColor: vars.color.surfaceHover,
          color: vars.color.textPrimary,
        },
      },
    },
    size: {
      sm: {
        fontSize: vars.fontSize.xs,
        padding: `${vars.space['2xs']} ${vars.space.sm}`,
        height: '28px',
      },
      md: {
        fontSize: vars.fontSize.sm,
        padding: `${vars.space.xs} ${vars.space.md}`,
        height: '32px',
      },
      lg: {
        fontSize: vars.fontSize.base,
        padding: `${vars.space.sm} ${vars.space.lg}`,
        height: '38px',
      },
    },
  },
  defaultVariants: {
    variant: 'secondary',
    size: 'md',
  },
})
