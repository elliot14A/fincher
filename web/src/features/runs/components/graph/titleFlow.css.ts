import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'

export const wrap = style({
  width: '100%',
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  minHeight: 0,
})

export const legend = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.md,
  padding: `${vars.space.xs} ${vars.space.sm}`,
  flexShrink: 0,
})

export const legendItem = style({
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space['2xs'],
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
  textTransform: 'uppercase',
  letterSpacing: '0.04em',
})

export const legendSwatchIncident = style({
  width: '10px',
  height: '3px',
  borderRadius: vars.radii.full,
  backgroundColor: vars.color.danger,
})

export const legendSwatchAllocation = style({
  width: '10px',
  height: '3px',
  borderRadius: vars.radii.full,
  backgroundColor: vars.color.primary,
})

export const legendSwatchResolution = style({
  width: '10px',
  height: '3px',
  borderRadius: vars.radii.full,
  backgroundColor: vars.color.success,
})

export const canvasHost = style({
  flex: 1,
  minHeight: 0,
})
