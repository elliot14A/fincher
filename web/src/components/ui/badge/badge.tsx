import type { ComponentChildren, JSX } from 'preact'
import { badgeRecipe } from './badge.css'

export type BadgeProps = JSX.HTMLAttributes<HTMLSpanElement> & {
  variant?: 'success' | 'danger' | 'warning' | 'teal' | 'primary' | 'neutral'
  children?: ComponentChildren
}

export function Badge({ variant = 'neutral', children, class: userClass, ...props }: BadgeProps) {
  const recipeClass = badgeRecipe({ variant })
  const combinedClass = userClass ? `${recipeClass} ${userClass}` : recipeClass

  return (
    <span class={combinedClass} {...props}>
      {children}
    </span>
  )
}
