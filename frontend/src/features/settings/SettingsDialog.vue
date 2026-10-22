<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import DatePicker from 'primevue/datepicker'
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'
import FloatLabel from 'primevue/floatlabel'
import { useAppStore } from '../../core/stores/app'
import { Coordinator, RemindKind, RemindStyle } from '../../core/bindings'
import { dateToHHMM, hhmmToDate } from '../../shared/datetime'
import squidward from '../../assets/images/squidward.png'
import RepeatSetting from './RepeatSetting/RepeatSetting.vue'

const visible = defineModel<boolean>('visible', { required: true })

const { t } = useI18n()
// 预览提醒标题与提醒窗口共享，留在全局，用 global composer 直读。
const { t: gt } = useI18n({ useScope: 'global' })
const store = useAppStore()

const offTimeDate = computed<Date>({
  get: () => hhmmToDate(store.offWork.offTime),
  set: (d: Date): Promise<void> => store.setOffTime(dateToHHMM(d)),
})

const repeatDays = computed<number[]>({
  get: () => [...(store.offWork.repeatDays ?? [])],
  set: (days: number[]): Promise<void> => store.saveOffWork({ ...store.offWork, repeatDays: days }),
})

const advanceMinutes = computed<number>({
  get: () => store.remind.advanceMinutes,
  set: (v: number): Promise<void> => store.setAdvanceMinutes(v ?? 0),
})

const style = computed<RemindStyle>({
  get: () => store.remind.style,
  set: (v: RemindStyle): Promise<void> => store.setRemindStyle(v),
})

interface StyleOption {
  value: RemindStyle
  label: string
  desc: string
  icon?: string
  image?: string
}

const styleOptions = computed<StyleOption[]>(() => [
  { value: RemindStyle.RemindStyleSquidward, label: t('squidward'), desc: t('squidwardDescription'), image: squidward },
])

function currentOption(value: RemindStyle | null | undefined): StyleOption | undefined {
  return styleOptions.value.find(o => o.value === value)
}

async function preview(): Promise<void> {
  await Coordinator.Preview(RemindKind.RemindKindOffWork, gt('remind.offworkTitle'))
}
</script>

<template>
  <Dialog
    v-model:visible="visible"
    modal
    closable
    append-to="body"
    :style="{ width: 'min(24rem, 92vw)' }"
    :header="t('title')"
  >
    <div class="flex flex-col gap-5 py-3">
      <!-- 下班时间 -->
      <FloatLabel variant="on">
        <DatePicker
          input-id="off-time"
          v-model="offTimeDate"
          time-only
          hour-format="24"
          :minute-step="5"
          append-to="body"
          class="w-full"
        />
        <label for="off-time">{{ t('offTime') }}</label>
      </FloatLabel>

      <!-- 提前提醒 -->
      <FloatLabel variant="on">
        <InputNumber
          input-id="advance-minutes"
          v-model="advanceMinutes"
          :min="0"
          :max="120"
          :use-grouping="false"
          show-buttons
          class="w-full"
        />
        <label for="advance-minutes">{{ t('advanceMinutes') }} (min)</label>
      </FloatLabel>

      <!-- 重复 -->
      <div class="flex flex-col gap-2.5">
        <span class="text-xs font-medium uppercase tracking-wide text-[var(--app-muted)]">{{ t('repeat') }}</span>
        <RepeatSetting v-model="repeatDays" />
      </div>

      <!-- 提醒样式 -->
      <FloatLabel variant="on">
        <Select
          v-model="style"
          :options="styleOptions"
          option-label="label"
          option-value="value"
          append-to="body"
          input-id="remind-style"
          class="w-full"
        >
          <template #value="slotProps">
            <div v-if="slotProps.value" class="flex items-center gap-2">
              <img v-if="currentOption(slotProps.value)?.image" :src="currentOption(slotProps.value)?.image" alt="" class="size-5 shrink-0 rounded-md object-cover">
              <i v-else :class="[currentOption(slotProps.value)?.icon, 'text-base text-[var(--app-muted)]']" />
              <span>{{ currentOption(slotProps.value)?.label }}</span>
            </div>
          </template>
          <template #option="slotProps">
            <div class="flex flex-col gap-0.5 py-0.5">
              <div class="flex items-center gap-2">
                <img v-if="slotProps.option.image" :src="slotProps.option.image" alt="" class="size-5 shrink-0 rounded-md object-cover">
                <i v-else :class="[slotProps.option.icon, 'text-base text-[var(--app-muted)]']" />
                <span class="text-[15px]">{{ slotProps.option.label }}</span>
              </div>
              <span class="pl-7 text-xs text-[var(--app-muted)]">{{ slotProps.option.desc }}</span>
            </div>
          </template>
        </Select>
        <label for="remind-style">{{ t('label') }}</label>
      </FloatLabel>
    </div>

    <!-- 预览提醒：唯一动作，放到弹窗底部右侧 -->
    <template #footer>
      <div class="flex w-full justify-end">
        <Button
          text
          severity="secondary"
          size="small"
          :label="t('previewRemind')"
          @click="preview"
        >
          <template #icon>
            <i class="i-lucide-eye text-base" />
          </template>
        </Button>
      </div>
    </template>
  </Dialog>
</template>

<i18n lang="yaml">
zh-CN:
  title: 设置
  offTime: 下班时间
  advanceMinutes: 下班提前提醒
  repeat: 重复
  label: 提醒样式
  squidward: 我下班了，蟹老板！
  squidwardDescription: 你我皆是章鱼哥
  previewRemind: 预览提醒
en-US:
  title: Settings
  offTime: Off Time
  advanceMinutes: Off-Work Advance
  repeat: Repeat
  label: Remind Style
  squidward: "I'm off work, Mr. Krabs!"
  squidwardDescription: You and I are both Squidward
  previewRemind: Preview Remind
</i18n>
