<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import Button from 'primevue/button'
import squidward from '../../assets/images/squidward.png'
import type { RemindCardEmits, RemindCardProps } from './types'

const { title, description } = defineProps<RemindCardProps>()
const emit = defineEmits<RemindCardEmits>()
const { t } = useI18n()
</script>

<template>
  <div class="remind-root flex h-full w-full items-stretch">
    <div class="remind-card flex flex-1 items-center gap-3 rounded-2xl px-4">
      <!-- 头像 -->
      <img :src="squidward" alt="" class="size-10 shrink-0 rounded-full object-cover">
      <!-- 标题 / 描述 -->
      <div class="flex min-w-0 flex-1 flex-col">
        <span class="truncate text-[15px] font-semibold text-[var(--app-text)]">{{ title }}</span>
        <span class="truncate text-xs text-[var(--app-muted)]">{{ description }}</span>
      </div>
      <!-- 挂断 / 接听：圆形图标按钮，文案通过 Tooltip hover 展示（left 方向避免超出窗口被裁切） -->
      <div class="flex shrink-0 items-center gap-2">
        <Button
          rounded
          severity="danger"
          class="!size-8 min-w-0 p-0"
          :aria-label="t('cancelText')"
          v-tooltip.left="t('cancelText')"
          @click="emit('cancel')"
        >
          <template #icon>
            <i class="i-lucide-phone-off text-base" />
          </template>
        </Button>
        <Button
          rounded
          severity="success"
          class="!size-8 min-w-0 p-0"
          :aria-label="t('confirmText')"
          v-tooltip.left="t('confirmText')"
          @click="emit('confirm')"
        >
          <template #icon>
            <i class="i-lucide-phone-call text-base" />
          </template>
        </Button>
      </div>
    </div>
  </div>
</template>

<i18n lang="yaml">
zh-CN:
  confirmText: 朕知道了
  cancelText: 我就硬拖
en-US:
  confirmText: "Yeah, I know"
  cancelText: I'll wing it
</i18n>
