import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'

export const shellRoot = style({
  display: 'flex',
  flexDirection: 'row',
  height: '100vh',
  width: '100vw',
  overflow: 'hidden',
  backgroundColor: vars.color.background,
})

export const shellMainArea = style({
  display: 'flex',
  flexDirection: 'column',
  flex: 1,
  height: '100vh',
  minWidth: 0,
  overflow: 'hidden',
})

export const shellContent = style({
  display: 'flex',
  flex: 1,
  minHeight: 0,
  overflow: 'hidden',
})
