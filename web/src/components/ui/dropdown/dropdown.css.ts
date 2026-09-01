import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

export const menuContainer = style({
  position: 'relative',
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
})

export const menuTrigger = style({
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  width: '28px',
  height: '28px',
  borderRadius: vars.radii.xs,
  backgroundColor: 'transparent',
  border: 'none',
  color: vars.color.textTertiary,
  cursor: 'pointer',
  transition: 'all 0.12s ease',
  ':hover': {
    backgroundColor: vars.color.surfaceHover,
    color: vars.color.textPrimary,
  },
  ':focus-visible': {
    outline: `2px solid ${vars.color.primary}`,
    outlineOffset: '1px',
  },
})

export const menuTriggerActive = style({
  backgroundColor: vars.color.surfaceElevated,
  color: vars.color.textPrimary,
})

export const menuDropdown = style({
  position: 'absolute',
  top: 'calc(100% + 4px)',
  right: 0,
  zIndex: 60,
  minWidth: '170px',
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.borderStrong}`,
  borderRadius: vars.radii.sm,
  padding: '4px',
  display: 'flex',
  flexDirection: 'column',
  gap: '1px',
})

export const menuItem = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.sm,
  width: '100%',
  padding: '6px 10px',
  fontFamily: fonts.sans,
  fontSize: vars.fontSize.xs,
  fontWeight: 450,
  color: vars.color.textSecondary,
  backgroundColor: 'transparent',
  border: 'none',
  borderRadius: vars.radii.xs,
  cursor: 'pointer',
  textAlign: 'left',
  transition: 'background-color 0.12s ease, color 0.12s ease',
  ':hover': {
    backgroundColor: vars.color.surfaceHover,
    color: vars.color.textPrimary,
  },
  ':focus-visible': {
    backgroundColor: vars.color.surfaceHover,
    color: vars.color.textPrimary,
    outline: 'none',
  },
  ':disabled': {
    opacity: 0.45,
    cursor: 'not-allowed',
    pointerEvents: 'none',
  },
})

export const menuItemDanger = style({
  color: vars.color.danger,
  ':hover': {
    backgroundColor: vars.color.dangerMuted,
    color: vars.color.danger,
  },
  ':focus-visible': {
    backgroundColor: vars.color.dangerMuted,
    color: vars.color.danger,
    outline: 'none',
  },
})

export const menuItemIcon = style({
  flexShrink: 0,
  opacity: 0.8,
})

export const menuDivider = style({
  height: '1px',
  backgroundColor: vars.color.borderSubtle,
  margin: '4px 0',
})
