import type { FormValidationStatus, InputNumberProps } from 'naive-ui'
import type { Size } from 'naive-ui/es/input-number/src/interface'

export interface FloatInputNumberProps extends /* @vue-ignore */ InputNumberProps {
  placeholder: string
  size?: Size
  defaultValue?: number | null
  status?: FormValidationStatus
  disabled?: boolean
  buttonPlacement?: 'right' | 'both'
}

export type FloatInputNumberInst = InstanceType<typeof FloatInputNumber>
