import type { FormValidationStatus, TimePickerProps, TimePickerSize } from 'naive-ui'

export interface FloatTimePickerProps extends /* @vue-ignore */ TimePickerProps {
  placeholder: string
  size?: TimePickerSize
  disabled?: boolean
  status?: FormValidationStatus
  defaultValue?: number | null
}

export type FloatTimePickerInst = InstanceType<typeof FloatTimePicker>
