import type { FormValidationStatus, SelectProps, SelectSize } from 'naive-ui'
import type { Value } from 'naive-ui/es/select/src/interface'

export interface FloatSelectProps extends /* @vue-ignore */ SelectProps {
  placeholder: string
  size?: SelectSize
  defaultValue?: Value | null
  status?: FormValidationStatus
  disabled?: boolean
}

export type FloatSelectInst = InstanceType<typeof FloatSelect>
