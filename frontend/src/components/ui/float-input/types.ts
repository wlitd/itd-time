import type { FormValidationStatus, InputProps, InputSize } from 'naive-ui'

export interface FloatInputProps extends /* @vue-ignore */ InputProps {
  placeholder: string
  pair?: boolean
  type?: 'text' | 'password' | 'textarea'
  size?: InputSize
  defaultValue?: string
  status?: FormValidationStatus
  disabled?: boolean
}

export type FloatInputInst = InstanceType<typeof FloatInput>
