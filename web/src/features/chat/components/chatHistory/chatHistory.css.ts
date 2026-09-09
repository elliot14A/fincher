import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

export const wrapper = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space['3xs'],
  marginTop: vars.space.md,
  paddingTop: vars.space.md,
  borderTop: `1px solid ${vars.color.borderSubtle}`,
  minHeight: 0,
  flex: 1,
})

export const sectionLabel = style({
  fontSize: vars.fontSize['3xs'],
  fontFamily: fonts.mono,
  textTransform: 'uppercase',
  letterSpacing: '0.08em',
  color: vars.color.textTertiary,
  padding: `${vars.space['2xs']} ${vars.space.sm}`,
})

export const scrollArea = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '1px',
  overflowY: 'auto',
  minHeight: 0,
})

export const emptyState = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textTertiary,
  padding: `${vars.space.xs} ${vars.space.sm}`,
  lineHeight: vars.lineHeight.snug,
})

export const row = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space['3xs'],
  paddingRight: vars.space['2xs'],
  borderRadius: vars.radii.sm,
  transition: 'background-color 0.1s ease',
  ':hover': {
    backgroundColor: vars.color.surfaceHover,
  },
})

export const rowActive = style({
  backgroundColor: vars.color.surfaceHover,
})

export const rowButton = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
  flex: 1,
  minWidth: 0,
  padding: `${vars.space['2xs']} ${vars.space.sm}`,
  border: 'none',
  backgroundColor: 'transparent',
  color: vars.color.textSecondary,
  cursor: 'pointer',
  textAlign: 'left',
  borderRadius: vars.radii.sm,
  transition: 'color 0.1s ease',
  ':hover': {
    color: vars.color.textPrimary,
  },
})

export const rowTitle = style({
  flex: 1,
  fontSize: vars.fontSize.xs,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

export const rowAction = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  width: '18px',
  height: '18px',
  flexShrink: 0,
  border: 'none',
  backgroundColor: 'transparent',
  borderRadius: vars.radii.xs,
  color: vars.color.textTertiary,
  cursor: 'pointer',
  opacity: 0,
  transition: 'opacity 0.1s ease, color 0.1s ease, background-color 0.1s ease',
  selectors: {
    [`${row}:hover &`]: {
      opacity: 1,
    },
    '&:hover': {
      color: vars.color.danger,
      backgroundColor: vars.color.dangerMuted,
    },
  },
})

export const renameInput = style({
  flex: 1,
  fontSize: vars.fontSize.xs,
  fontFamily: fonts.sans,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.primaryBorder}`,
  borderRadius: vars.radii.xs,
  color: vars.color.textPrimary,
  padding: `${vars.space['3xs']} ${vars.space['2xs']}`,
  outline: 'none',
  minWidth: 0,
})
