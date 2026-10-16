<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useConfirm } from 'primevue/useconfirm'
import Checkbox from 'primevue/checkbox'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import DatePicker from 'primevue/datepicker'
import FloatLabel from 'primevue/floatlabel'
import { formatDeadline } from '../../../shared/datetime'
import type { TodoItemEmits, TodoItemProps } from './types'

const { todo, editing = false } = defineProps<TodoItemProps>()
const emit = defineEmits<TodoItemEmits>()
const { t } = useI18n()
// 跨组件共享键（标题/截止时间占位、通用取消按钮）留在全局，用 global composer 直读避免回落告警。
const { t: gt } = useI18n({ useScope: 'global' })
const confirm = useConfirm()

const draftTitle = ref<string>('')
const draftDeadline = ref<Date>(new Date())
// 删除按钮容器引用：作为 ConfirmPopup 锚点（不依赖 event.currentTarget，避免其为 null 时回退到 body 导致弹窗跑到左上角）
const deleteRef = ref<HTMLDivElement>()

// 每个待办的输入框需唯一 id，供 FloatLabel 的 label[for] 关联
const titleId = computed<string>(() => `todo-edit-title-${todo.id}`)
const dateId = computed<string>(() => `todo-edit-date-${todo.id}`)

// 进入编辑态时，用当前项初始化草稿
watch(() => editing, (val: boolean) => {
  if (val) {
    draftTitle.value = todo.title
    draftDeadline.value = new Date(todo.deadline)
  }
}, { immediate: true })

function onToggle(value: boolean): void {
  emit('toggle', todo.id, value)
}

// 完成项不可编辑，点击行仅对未完成项进入编辑态
function onRowClick(): void {
  if (!todo.completed)
    emit('edit', todo)
}

function save(): void {
  const title = draftTitle.value.trim()
  if (!title || !draftDeadline.value)
    return
  emit('save', { id: todo.id, title, deadline: draftDeadline.value.getTime() })
}

// 删除走 ConfirmPopup 二次确认，锚定到删除按钮所在的容器
function onRemoveClick(): void {
  const target = deleteRef.value
  if (!target)
    return
  confirm.require({
    target,
    message: t('deleteConfirm'),
    icon: 'i-lucide-alert-circle',
    acceptLabel: t('delete'),
    rejectLabel: gt('update.cancel'),
    acceptProps: { severity: 'danger', size: 'small' },
    rejectProps: { severity: 'secondary', text: true, size: 'small' },
    accept: () => emit('remove', todo.id),
  })
  // 兜底定位：PrimeVue 4.5.5 的 ConfirmPopup 首次打开不会调用 alignOverlay，
  // 面板会停在 CSS 默认的 top:0/left:0（左上角），这里在渲染后手动对齐到锚点。
  void nextTick((): void => alignToTarget(target))
}

// 将确认面板对齐到锚点下方，越界时夹取/翻转到视口内
function alignToTarget(target: HTMLElement): void {
  const popup = document.querySelector<HTMLElement>('.p-confirmpopup')
  if (!popup)
    return
  const gap = 8
  const rect = target.getBoundingClientRect()
  const pw = popup.offsetWidth
  const ph = popup.offsetHeight
  let left = rect.left
  let top = rect.bottom + gap
  if (left + pw > window.innerWidth - gap)
    left = Math.max(gap, window.innerWidth - gap - pw)
  if (top + ph > window.innerHeight - gap)
    top = Math.max(gap, rect.top - ph - gap)
  popup.style.left = `${left}px`
  popup.style.top = `${top}px`
}
</script>

<template>
  <!-- 行内编辑态：两个 FloatLabel 输入框 + 右侧勾图标保存，无取消/保存按钮 -->
  <li v-if="editing" class="w-full rounded-xl bg-[var(--app-surface-2)] p-2.5 ring-1 ring-primary/40">
    <div class="flex items-center gap-2">
      <div class="flex min-w-0 flex-1 flex-col gap-2">
        <FloatLabel variant="on">
          <InputText
            :id="titleId"
            v-model="draftTitle"
            fluid
            size="small"
            autofocus
            @keyup.enter="save"
            @keyup.escape="emit('cancel')"
          />
          <label :for="titleId">{{ gt('todos.titlePlaceholder') }}</label>
        </FloatLabel>
        <FloatLabel variant="on">
          <DatePicker
            :input-id="dateId"
            v-model="draftDeadline"
            fluid
            size="small"
            show-time
            hour-format="24"
            :minute-step="5"
            date-format="yy-mm-dd"
            append-to="body"
          />
          <label :for="dateId">{{ gt('todos.deadlinePlaceholder') }}</label>
        </FloatLabel>
      </div>

      <Button
        text
        severity="success"
        rounded
        class="shrink-0"
        :disabled="!draftTitle.trim()"
        :aria-label="t('save')"
        @click="save"
      >
        <template #icon>
          <i class="i-lucide-check text-lg" />
        </template>
      </Button>
    </div>
  </li>

  <!-- 展示态：完成 / 超时 / 普通 三态显著区分 -->
  <li
    v-else
    class="group relative flex w-full items-stretch overflow-hidden rounded-xl ring-1 transition-all duration-200"
    :class="[
      todo.completed
        ? 'bg-primary/8 ring-primary/25'
        : todo.isOverdue
          ? 'bg-danger/10 ring-danger/40'
          : 'bg-[var(--app-surface-2)] ring-[var(--app-border)] hover:ring-primary/40',
      todo.completed ? 'cursor-default' : 'cursor-pointer',
    ]"
    @click="onRowClick"
  >
    <!-- 左侧状态条 -->
    <span
      class="w-1 shrink-0"
      :class="todo.completed ? 'bg-primary opacity-55' : todo.isOverdue ? 'bg-[var(--app-danger)]' : 'bg-primary/40'"
    />

    <div class="flex flex-1 items-center gap-3 py-2.5 pl-3 pr-2">
      <!-- 勾选框 + 文本：完成态整体淡出，模拟 disabled -->
      <div
        class="flex min-w-0 flex-1 items-center gap-3 transition-opacity duration-200"
        :class="todo.completed ? 'opacity-55' : ''"
      >
        <div class="flex items-center" @click.stop>
          <Checkbox
            :model-value="todo.completed"
            binary
            size="small"
            @update:model-value="onToggle"
          />
        </div>

        <div class="flex min-w-0 flex-1 flex-col gap-0.5">
          <div class="flex min-w-0 items-center gap-2">
            <span
              class="truncate text-sm font-medium"
              :class="todo.completed ? 'text-[var(--app-muted)]' : 'text-[var(--app-text)]'"
            >
              {{ todo.title || t('untitled') }}
            </span>
            <!-- 完成 / 超时 tag 紧跟标题之后 -->
            <span
              v-if="todo.completed"
              class="flex shrink-0 items-center gap-1 rounded-full bg-primary/15 px-1.5 py-0.5 text-[10px] font-semibold text-primary"
            >
              <i class="i-lucide-check text-[11px]" />
              {{ t('done') }}
            </span>
            <span
              v-else-if="todo.isOverdue"
              class="flex shrink-0 items-center gap-1 rounded-full bg-danger/15 px-1.5 py-0.5 text-[10px] font-semibold text-[var(--app-danger)]"
            >
              <i class="i-lucide-alert-circle text-[11px]" />
              {{ t('overdue') }}
            </span>
          </div>
          <div class="flex items-center gap-2">
            <span
              class="text-xs"
              :class="todo.isOverdue && !todo.completed ? 'text-[var(--app-danger)] font-medium' : 'text-[var(--app-muted)]'"
            >
              {{ formatDeadline(todo.deadline, t('timeFormat')) }}
            </span>
          </div>
        </div>
      </div>

      <!-- 删除按钮：完成态仍保持可用与正常颜色 -->
      <div ref="deleteRef" class="flex shrink-0 items-center gap-0.5">
        <Button
          severity="danger"
          text
          rounded
          size="small"
          class="opacity-0 transition-opacity duration-200 group-hover:opacity-100 group-focus-within:opacity-100"
          :aria-label="t('deleteConfirm')"
          @click.stop="onRemoveClick"
        >
          <template #icon>
            <i class="i-lucide-trash-2 text-lg" />
          </template>
        </Button>
      </div>
    </div>
  </li>
</template>

<i18n lang="yaml">
zh-CN:
  deleteConfirm: 确定删除该待办吗？
  delete: 删除
  save: 保存
  untitled: 未命名事项
  done: 已完成
  overdue: 已超时
  timeFormat: M月D日 HH:mm
en-US:
  deleteConfirm: Delete this todo?
  delete: Delete
  save: Save
  untitled: Untitled
  done: Done
  overdue: Overdue
  timeFormat: "MMM D, HH:mm"
</i18n>
