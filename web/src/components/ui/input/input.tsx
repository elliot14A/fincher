import type { ComponentChildren, JSX } from 'preact'
import {
  errorText,
  formField,
  helperText,
  label as labelClass,
  labelRow,
  optionalText,
  requiredIndicator,
  selectInput,
  textInput,
} from './input.css'

export type FormFieldProps = {
  label: string
  htmlFor?: string
  required?: boolean
  optional?: boolean
  helper?: string
  error?: string
  children: ComponentChildren
}

export function FormField({
  label,
  htmlFor,
  required,
  optional,
  helper,
  error,
  children,
}: FormFieldProps) {
  return (
    <div class={formField}>
      <div class={labelRow}>
        <label htmlFor={htmlFor} class={labelClass}>
          {label}
          {required ? <span class={requiredIndicator}>*</span> : null}
        </label>
        {optional ? <span class={optionalText}>Optional</span> : null}
      </div>

      {children}

      {error ? (
        <span class={errorText}>{error}</span>
      ) : helper ? (
        <span class={helperText}>{helper}</span>
      ) : null}
    </div>
  )
}

export type TextInputProps = JSX.IntrinsicElements['input'] & {
  hasError?: boolean
}

export function TextInput({ hasError, class: className, ...props }: TextInputProps) {
  return <input class={`${textInput} ${className ?? ''}`} {...props} />
}

export type SelectOption = {
  value: string
  label: string
  disabled?: boolean
}

export type SelectInputProps = JSX.IntrinsicElements['select'] & {
  options: SelectOption[]
  hasError?: boolean
}

export function SelectInput({ options, hasError, class: className, ...props }: SelectInputProps) {
  return (
    <select class={`${selectInput} ${className ?? ''}`} {...props}>
      {options.map((opt) => (
        <option key={opt.value} value={opt.value} disabled={opt.disabled}>
          {opt.label}
        </option>
      ))}
    </select>
  )
}
