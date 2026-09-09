import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

export const popover = style({
  position: 'absolute',
  bottom: '100%',
  left: 0,
  marginBottom: vars.space.sm,
  zIndex: 30,
  width: '380px',
  maxWidth: '90vw',
  maxHeight: '260px',
  overflowY: 'auto',
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.md,
  boxShadow: '0 12px 32px rgba(0, 0, 0, 0.4)',
  padding: vars.space['3xs'],
})

export const empty = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textTertiary,
  textAlign: 'center',
  padding: `${vars.space.md} ${vars.space.sm}`,
})

export const groupLabel = style({
  fontSize: vars.fontSize['3xs'],
  fontFamily: fonts.mono,
  textTransform: 'uppercase',
  letterSpacing: '0.12em',
  color: vars.color.textTertiary,
  padding: `${vars.space.xs} ${vars.space.sm} ${vars.space['2xs']}`,
})

export const option = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.sm,
  width: '100%',
  padding: `${vars.space.xs} ${vars.space.sm}`,
  border: 'none',
  backgroundColor: 'transparent',
  borderRadius: vars.radii.sm,
  cursor: 'pointer',
  textAlign: 'left',
  color: vars.color.textSecondary,
  transition: 'background-color 0.1s ease, color 0.1s ease',
})

export const optionActive = style({
  backgroundColor: vars.color.primary,
  color: vars.color.textInverse,
})

export const poster = style({
  width: '28px',
  height: '40px',
  objectFit: 'cover',
  borderRadius: vars.radii.xs,
  flexShrink: 0,
  backgroundColor: vars.color.surfaceActive,
})

export const iconBox = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  width: '28px',
  height: '28px',
  borderRadius: vars.radii.xs,
  backgroundColor: vars.color.surfaceActive,
  color: vars.color.textTertiary,
  flexShrink: 0,
})

export const optionBody = style({
  display: 'flex',
  flexDirection: 'column',
  minWidth: 0,
  flex: 1,
})

export const optionName = style({
  fontSize: vars.fontSize.sm,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

export const optionMeta = style({
  fontSize: vars.fontSize['3xs'],
  fontFamily: fonts.mono,
  color: vars.color.textTertiary,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

export const optionMetaActive = style({
  color: vars.color.textInverse,
  opacity: 0.8,
})
