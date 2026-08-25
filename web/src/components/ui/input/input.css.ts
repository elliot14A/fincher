import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'

export const formField = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.xs,
})

export const labelRow = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
})

export const label = style({
  fontSize: vars.fontSize.xs,
  fontWeight: 500,
  color: vars.color.textPrimary,
  letterSpacing: '0.02em',
})

export const requiredIndicator = style({
  color: vars.color.danger,
  marginLeft: vars.space['3xs'],
})

export const optionalText = style({
  fontSize: vars.fontSize['2xs'],
  color: vars.color.textTertiary,
})

const baseInput = style({
  backgroundColor: vars.color.surface,
  border: `1px solid ${vars.color.borderStrong}`,
  borderRadius: vars.radii.sm,
  color: vars.color.textPrimary,
  fontSize: vars.fontSize.sm,
  padding: `${vars.space.sm} ${vars.space.md}`,
  fontFamily: 'inherit',
  width: '100%',
  boxSizing: 'border-box',
  transition: 'border-color 0.15s ease, box-shadow 0.15s ease',
  outline: 'none',
  ':focus': {
    borderColor: vars.color.primary,
    boxShadow: `0 0 0 1px ${vars.color.primary}, 0 0 0 3px ${vars.color.primaryMuted}`,
  },
  '::placeholder': {
    color: vars.color.textTertiary,
  },
  ':disabled': {
    opacity: 0.6,
    cursor: 'not-allowed',
  },
})

export const textInput = style([baseInput])

export const selectInput = style([
  baseInput,
  {
    appearance: 'none',
    backgroundImage: `url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%238D8E96' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='m6 9 6 6 6-6'/%3E%3C/svg%3E")`,
    backgroundRepeat: 'no-repeat',
    backgroundPosition: `right ${vars.space.md} center`,
    paddingRight: vars.space['2xl'],
    cursor: 'pointer',
  },
])

export const helperText = style({
  fontSize: vars.fontSize['2xs'],
  color: vars.color.textTertiary,
  lineHeight: vars.lineHeight.normal,
})

export const errorText = style({
  fontSize: vars.fontSize['2xs'],
  color: vars.color.danger,
  lineHeight: vars.lineHeight.normal,
})
