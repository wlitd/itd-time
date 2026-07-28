import type { DatePickerProps, DatePickerSize, FormValidationStatus } from 'naive-ui'

export interface FloatDatePickerProps extends /* @vue-ignore */ DatePickerProps {
  placeholder: string
  size?: DatePickerSize
  defaultValue?: number | null
  status?: FormValidationStatus
  disabled?: boolean
  type?: 'date' | 'datetime' | 'month' | 'year' | 'quarter' | 'week'
  panel?: boolean
}

export type FloatDatePickerInst = InstanceType<typeof FloatDatePicker>
