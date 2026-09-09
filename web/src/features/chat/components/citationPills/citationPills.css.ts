import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

export const wrapper = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.xs,
  paddingTop: vars.space.md,
  borderTop: `1px solid ${vars.color.borderSubtle}`,
})

export const label = style({
  fontSize: vars.fontSize['3xs'],
  fontFamily: fonts.mono,
  textTransform: 'uppercase',
  color: vars.color.textTertiary,
  letterSpacing: '0.06em',
})

export const list = style({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'flex-start',
  gap: vars.space.xs,
})

export const pill = style({
  fontSize: vars.fontSize['3xs'],
  fontFamily: fonts.mono,
  color: vars.color.teal,
  backgroundColor: vars.color.tealMuted,
  padding: `${vars.space['2xs']} ${vars.space.sm}`,
  borderRadius: vars.radii.xs,
  border: `1px solid ${vars.color.tealBorder}`,
  display: 'inline-flex',
  alignItems: 'flex-start',
  gap: vars.space.xs,
  maxWidth: '100%',
  lineHeight: vars.lineHeight.snug,
})

export const pillText = style({
  whiteSpace: 'pre-wrap',
  overflowWrap: 'anywhere',
  wordBreak: 'break-word',
})

export const pillIcon = style({
  flexShrink: 0,
  marginTop: vars.space['3xs'],
})
