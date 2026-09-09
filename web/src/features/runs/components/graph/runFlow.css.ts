import { globalStyle, style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'

export const flowShell = style({
  position: 'relative',
  width: '100%',
  height: '220px',
  borderRadius: vars.radii.md,
  border: `1px solid ${vars.color.border}`,
  backgroundColor: vars.color.background,
  overflow: 'hidden',
})

export const flowShellFull = style({
  height: '100%',
  borderRadius: 0,
  border: 'none',
})

export const canvas = style({
  width: '100%',
  height: '100%',
})

export const maximizeBtn = style({
  position: 'absolute',
  top: vars.space.sm,
  right: vars.space.sm,
  zIndex: 5,
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space['2xs'],
  padding: `${vars.space['2xs']} ${vars.space.xs}`,
  borderRadius: vars.radii.xs,
  border: `1px solid ${vars.color.border}`,
  backgroundColor: vars.color.surfaceElevated,
  color: vars.color.textSecondary,
  cursor: 'pointer',
  fontSize: vars.fontSize['3xs'],
  transition: 'color 0.12s ease, background-color 0.12s ease, border-color 0.12s ease',
  ':hover': {
    color: vars.color.textPrimary,
    backgroundColor: vars.color.surfaceHover,
    borderColor: vars.color.borderStrong,
  },
})

globalStyle(`${canvas} .react-flow__background`, {
  backgroundColor: vars.color.background,
})

globalStyle(`${canvas} .react-flow__edge-path`, {
  stroke: vars.color.borderStrong,
  strokeWidth: 1.5,
})

globalStyle(`${canvas} .react-flow__controls`, {
  boxShadow: 'none',
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.xs,
  overflow: 'hidden',
})

globalStyle(`${canvas} .react-flow__controls-button`, {
  backgroundColor: vars.color.surfaceElevated,
  borderBottom: `1px solid ${vars.color.borderSubtle}`,
  color: vars.color.textSecondary,
  fill: vars.color.textSecondary,
})

globalStyle(`${canvas} .react-flow__controls-button:hover`, {
  backgroundColor: vars.color.surfaceHover,
})

globalStyle(`${canvas} .react-flow__handle`, {
  opacity: 0,
})

globalStyle(`${canvas} .react-flow__attribution`, {
  display: 'none',
})
