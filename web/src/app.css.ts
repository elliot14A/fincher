import { globalStyle } from '@vanilla-extract/css'
import { vars } from './styles/theme.css'
import { fonts } from './styles/tokens'

globalStyle(':root', {
  colorScheme: 'dark',
})

globalStyle('*, *::before, *::after', {
  boxSizing: 'border-box',
  margin: 0,
  padding: 0,
})

globalStyle('html, body', {
  height: '100%',
  width: '100%',
  backgroundColor: vars.color.background,
  color: vars.color.textPrimary,
  fontFamily: fonts.sans,
  fontFeatureSettings: '"cv02", "cv03", "cv04", "cv11", "tnum"',
  fontVariantNumeric: 'tabular-nums',
  WebkitFontSmoothing: 'antialiased',
  MozOsxFontSmoothing: 'grayscale',
})

globalStyle('#app', {
  height: '100%',
  width: '100%',
  display: 'flex',
  flexDirection: 'column',
})

globalStyle('a', {
  color: 'inherit',
  textDecoration: 'none',
})

globalStyle('button', {
  fontFamily: 'inherit',
  border: 'none',
  background: 'none',
  cursor: 'pointer',
})
