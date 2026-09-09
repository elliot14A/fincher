import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'

export const paragraph = style({
  margin: 0,
  whiteSpace: 'pre-wrap',
  overflowWrap: 'anywhere',
})

export const mention = style({
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space['3xs'],
  margin: '0 1px',
  padding: `0 ${vars.space['2xs']}`,
  borderRadius: vars.radii.xs,
  border: `1px solid ${vars.color.primaryBorder}`,
  backgroundColor: vars.color.primaryMuted,
  color: vars.color.primary,
  fontWeight: 500,
  verticalAlign: 'baseline',
})
