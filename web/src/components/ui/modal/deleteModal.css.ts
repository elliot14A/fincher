import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'

export const warningBox = style({
  display: 'flex',
  alignItems: 'flex-start',
  gap: vars.space.md,
  padding: vars.space.md,
  backgroundColor: vars.color.dangerMuted,
  border: `1px solid ${vars.color.dangerBorder}`,
  borderRadius: vars.radii.md,
})

export const warningIcon = style({
  color: vars.color.danger,
  flexShrink: 0,
  marginTop: '2px',
})

export const warningContent = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space['2xs'],
})

export const warningTitle = style({
  fontSize: vars.fontSize.sm,
  fontWeight: 600,
  color: vars.color.textPrimary,
  lineHeight: vars.lineHeight.tight,
})

export const warningText = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textSecondary,
  lineHeight: vars.lineHeight.normal,
})

export const entityBadge = style({
  fontFamily: 'monospace',
  fontSize: vars.fontSize.xs,
  color: vars.color.danger,
  backgroundColor: 'rgba(229, 72, 77, 0.15)',
  padding: `${vars.space['3xs']} ${vars.space.xs}`,
  borderRadius: vars.radii.xs,
  fontWeight: 500,
})

export const promptText = style({
  fontSize: vars.fontSize.sm,
  color: vars.color.textPrimary,
  lineHeight: vars.lineHeight.normal,
})
