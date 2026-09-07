import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'

export * from '#/components/ui/sidebar/sidebar.css'

export const complianceBanner = style({
  display: 'flex',
  alignItems: 'flex-start',
  gap: vars.space.sm,
  padding: vars.space.sm,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.xs,
  color: vars.color.textSecondary,
  fontSize: vars.fontSize.xs,
  lineHeight: vars.lineHeight.normal,
})

export const complianceIcon = style({
  color: vars.color.success,
  flexShrink: 0,
  marginTop: '2px',
})
