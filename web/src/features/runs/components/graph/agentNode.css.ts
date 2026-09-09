import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

export const node = style({
  width: '188px',
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.xs,
  padding: vars.space.sm,
  borderRadius: vars.radii.md,
  border: `1px solid ${vars.color.border}`,
  backgroundColor: vars.color.surfaceElevated,
  cursor: 'pointer',
  transition: 'border-color 0.12s ease, background-color 0.12s ease, transform 0.12s ease',
  ':hover': {
    borderColor: vars.color.borderStrong,
    backgroundColor: vars.color.surfaceHover,
  },
})

export const nodeSelected = style({
  borderColor: vars.color.primary,
  backgroundColor: vars.color.primaryMuted,
})

export const nodeSuccess = style({
  borderTop: `2px solid ${vars.color.success}`,
})

export const nodeDanger = style({
  borderTop: `2px solid ${vars.color.danger}`,
})

export const nodeRunning = style({
  borderTop: `2px solid ${vars.color.warning}`,
})

export const nodeHeader = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
})

export const iconBox = style({
  width: '26px',
  height: '26px',
  borderRadius: vars.radii.sm,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  flexShrink: 0,
  backgroundColor: vars.color.surfaceActive,
  color: vars.color.textSecondary,
})

export const iconBoxSuccess = style({
  backgroundColor: vars.color.successMuted,
  color: vars.color.success,
})

export const iconBoxDanger = style({
  backgroundColor: vars.color.dangerMuted,
  color: vars.color.danger,
})

export const iconBoxRunning = style({
  backgroundColor: vars.color.warningMuted,
  color: vars.color.warning,
})

export const titleGroup = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '1px',
  minWidth: 0,
})

export const title = style({
  fontSize: vars.fontSize.xs,
  fontWeight: 600,
  color: vars.color.textPrimary,
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
})

export const agent = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
})

export const footerRow = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: vars.space.xs,
})

export const latency = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
  fontVariantNumeric: 'tabular-nums',
})

export const stageIndex = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
})
