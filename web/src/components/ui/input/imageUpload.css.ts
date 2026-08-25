import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'

export const dropzone = style({
  border: `1px dashed ${vars.color.borderStrong}`,
  borderRadius: vars.radii.md,
  padding: `${vars.space.md} ${vars.space.lg}`,
  backgroundColor: vars.color.surface,
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  cursor: 'pointer',
  transition: 'border-color 0.15s ease, background-color 0.15s ease',
  textAlign: 'center',
  gap: vars.space.xs,
  ':hover': {
    borderColor: vars.color.primary,
    backgroundColor: vars.color.surfaceHover,
  },
})

export const dropzoneActive = style({
  borderColor: vars.color.primary,
  backgroundColor: vars.color.primaryMuted,
})

export const uploadIcon = style({
  color: vars.color.textTertiary,
})

export const uploadPrompt = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textPrimary,
  fontWeight: 500,
})

export const uploadHint = style({
  fontSize: vars.fontSize['2xs'],
  color: vars.color.textTertiary,
})

export const hiddenFileInput = style({
  display: 'none',
})

export const previewContainer = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.md,
  padding: vars.space.sm,
  backgroundColor: vars.color.surface,
  border: `1px solid ${vars.color.borderStrong}`,
  borderRadius: vars.radii.md,
})

export const thumbnail = style({
  width: '48px',
  height: '48px',
  borderRadius: vars.radii.sm,
  objectFit: 'cover',
  border: `1px solid ${vars.color.borderSubtle}`,
})

export const previewDetails = style({
  display: 'flex',
  flexDirection: 'column',
  flex: 1,
  gap: vars.space['3xs'],
  minWidth: 0,
})

export const previewFilename = style({
  fontSize: vars.fontSize.xs,
  fontWeight: 500,
  color: vars.color.textPrimary,
  whiteSpace: 'nowrap',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
})

export const previewSize = style({
  fontSize: vars.fontSize['2xs'],
  color: vars.color.textTertiary,
})

export const removeButton = style({
  background: 'transparent',
  border: 'none',
  color: vars.color.textTertiary,
  cursor: 'pointer',
  padding: vars.space.xs,
  borderRadius: vars.radii.sm,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  ':hover': {
    color: vars.color.danger,
    backgroundColor: vars.color.dangerMuted,
  },
})

export const errorText = style({
  fontSize: vars.fontSize['2xs'],
  color: vars.color.danger,
  marginTop: vars.space['2xs'],
})
