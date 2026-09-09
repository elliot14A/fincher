import { globalStyle, style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

export const markdown = style({
  fontSize: vars.fontSize.sm,
  lineHeight: vars.lineHeight.normal,
  color: vars.color.textPrimary,
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.md,
})

globalStyle(`${markdown} p`, {
  margin: 0,
})

globalStyle(`${markdown} strong`, {
  color: vars.color.textPrimary,
  fontWeight: 600,
})

globalStyle(`${markdown} ul, ${markdown} ol`, {
  margin: 0,
  paddingLeft: vars.space.xl,
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.xs,
})

globalStyle(`${markdown} li`, {
  lineHeight: vars.lineHeight.snug,
})

globalStyle(`${markdown} li::marker`, {
  color: vars.color.teal,
})

globalStyle(`${markdown} h1, ${markdown} h2, ${markdown} h3, ${markdown} h4`, {
  margin: 0,
  fontSize: vars.fontSize.base,
  fontWeight: 600,
  color: vars.color.textPrimary,
})

globalStyle(`${markdown} a`, {
  color: vars.color.primary,
  textDecoration: 'none',
})

globalStyle(`${markdown} a:hover`, {
  textDecoration: 'underline',
})

globalStyle(`${markdown} code`, {
  fontFamily: fonts.mono,
  fontSize: vars.fontSize.xs,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.borderSubtle}`,
  borderRadius: vars.radii.xs,
  padding: `${vars.space['3xs']} ${vars.space['2xs']}`,
  color: vars.color.teal,
})

globalStyle(`${markdown} pre`, {
  margin: 0,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.borderSubtle}`,
  borderRadius: vars.radii.sm,
  padding: vars.space.md,
  overflowX: 'auto',
})

globalStyle(`${markdown} pre code`, {
  border: 'none',
  backgroundColor: 'transparent',
  padding: 0,
  color: vars.color.textSecondary,
  whiteSpace: 'pre',
})

globalStyle(`${markdown} table`, {
  borderCollapse: 'collapse',
  width: '100%',
  fontSize: vars.fontSize.xs,
})

globalStyle(`${markdown} th, ${markdown} td`, {
  border: `1px solid ${vars.color.borderSubtle}`,
  padding: `${vars.space['2xs']} ${vars.space.sm}`,
  textAlign: 'left',
})

globalStyle(`${markdown} th`, {
  color: vars.color.textSecondary,
  fontWeight: 600,
  backgroundColor: vars.color.surfaceElevated,
})

globalStyle(`${markdown} blockquote`, {
  margin: 0,
  paddingLeft: vars.space.md,
  borderLeft: `2px solid ${vars.color.borderStrong}`,
  color: vars.color.textSecondary,
})
