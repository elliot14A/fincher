import type { ComponentChildren, JSX } from 'preact'
import { buttonRecipe } from './button.css'

export type ButtonProps = JSX.HTMLAttributes<HTMLButtonElement> & {
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost'
  size?: 'sm' | 'md' | 'lg'
  children?: ComponentChildren
  disabled?: boolean
  onClick?: JSX.MouseEventHandler<HTMLButtonElement>
}

export function Button({
  variant = 'secondary',
  size = 'md',
  children,
  class: userClass,
  ...props
}: ButtonProps) {
  const recipeClass = buttonRecipe({ variant, size })
  const combinedClass = userClass ? `${recipeClass} ${userClass}` : recipeClass

  return (
    <button class={combinedClass} {...props}>
      {children}
    </button>
  )
}
