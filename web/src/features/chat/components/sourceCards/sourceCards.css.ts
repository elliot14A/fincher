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
  letterSpacing: '0.06em',
  color: vars.color.textTertiary,
})

export const list = style({
  display: 'flex',
  flexWrap: 'wrap',
  gap: vars.space.sm,
})

export const card = style({
  display: 'flex',
  alignItems: 'stretch',
  gap: vars.space.sm,
  width: '240px',
  maxWidth: '100%',
  padding: vars.space.xs,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.md,
  overflow: 'hidden',
})

export const poster = style({
  width: '46px',
  height: '66px',
  objectFit: 'cover',
  borderRadius: vars.radii.sm,
  flexShrink: 0,
  backgroundColor: vars.color.surfaceActive,
})

export const body = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space['3xs'],
  minWidth: 0,
  flex: 1,
  paddingRight: vars.space['2xs'],
})

export const name = style({
  fontSize: vars.fontSize.sm,
  fontWeight: 600,
  color: vars.color.textPrimary,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

export const meta = style({
  fontSize: vars.fontSize['3xs'],
  fontFamily: fonts.mono,
  color: vars.color.textTertiary,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

export const statusPill = style({
  marginTop: 'auto',
  alignSelf: 'flex-start',
  fontSize: vars.fontSize['3xs'],
  fontFamily: fonts.mono,
  color: vars.color.teal,
  backgroundColor: vars.color.tealMuted,
  border: `1px solid ${vars.color.tealBorder}`,
  borderRadius: vars.radii.xs,
  padding: `${vars.space['3xs']} ${vars.space['2xs']}`,
})
