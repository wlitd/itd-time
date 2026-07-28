<script setup lang="ts">
import dayjs from 'dayjs'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const todoStore = useTodoStore()
const { sortedTodos } = storeToRefs(todoStore)

/** 当前编辑中的 TODO ID */
const editingId = ref<string | null>(null)

/** 编辑草稿，点击保存时才提交到 store */
const draft = reactive<{ title: string | null, deadline: number | null }>({
  title: '',
  deadline: null
})

/**
 * 是否处于编辑态
 */
const isEditing = (id: string): boolean => editingId.value === id

/**
 * 进入编辑模式（已完成的任务不可编辑，切换目标前自动保存上一个草稿）
 */
function startEdit(todo: DisplayTodo): void {
  if (todo.completed || isEditing(todo.id)) {
    return
  }
  saveEdit()
  editingId.value = todo.id
  draft.title = todo.title
  draft.deadline = todo.deadline
}

/**
 * 保存草稿并退出编辑模式
 */
function saveEdit(): void {
  if (!editingId.value) {
    return
  }
  const todo = todoStore.todos.find(item => item.id === editingId.value)
  if (todo) {
    todoStore.updateTodo({
      ...todo,
      title: (draft.title ?? '').trim(),
      deadline: draft.deadline ?? todo.deadline
    })
  }
  editingId.value = null
}

/**
 * 切换完成状态（若该项正在编辑则先保存）
 */
function toggleDone(todo: DisplayTodo): void {
  if (isEditing(todo.id)) {
    saveEdit()
  }
  todoStore.toggleTodo(todo.id)
}

/**
 * 删除代办
 */
function handleRemove(id: string): void {
  if (isEditing(id)) {
    editingId.value = null
  }
  todoStore.removeTodo(id)
}

/**
 * 格式化展示时间
 */
function formatTime(timestamp: number): string {
  return dayjs(timestamp).format(t('timeFormat'))
}

/**
 * 新增代办：默认进入编辑状态，并滚动至可视区域、聚焦标题输入框
 */
async function scrollAndAdd(): Promise<void> {
  saveEdit()
  todoStore.addTodo({
    title: '',
    deadline: dayjs().add(1, 'hour').startOf('minute').valueOf()
  })

  // store 中 addTodo 为 push，新项在数组末尾
  const newTodo = todoStore.todos.at(-1)
  if (!newTodo) {
    return
  }
  editingId.value = newTodo.id
  draft.title = newTodo.title
  draft.deadline = newTodo.deadline

  await nextTick()
  const targetEl = document.querySelector<HTMLElement>(`[data-id="${newTodo.id}"]`)
  targetEl?.scrollIntoView({ behavior: 'smooth', block: 'center' })
  // 聚焦标题输入框（编辑区第一个 input 即 FloatInput）
  targetEl?.querySelector('input')?.focus()
}

defineExpose({
  scrollAndAdd
})
</script>

<template>
  <!-- 空状态：水平垂直居中 -->
  <div v-if="sortedTodos.length === 0" class="flex-1 flex-center">
    <NEmpty :description="t('empty')" size="large">
      <template #icon>
        <div class="i-lucide:clipboard-list" />
      </template>
    </NEmpty>
  </div>

  <!-- 列表 -->
  <NScrollbar v-else class="flex-1">
    <NList hoverable :show-divider="false">
      <NListItem
        v-for="todo in sortedTodos" :key="todo.id" :data-id="todo.id"
        class="mb-2! rounded-lg transition-opacity duration-300"
        :class="todo.completed ? 'opacity-55' : 'cursor-pointer'" @click="startEdit(todo)"
      >
        <div class="flex items-center gap-3">
          <!-- 完成状态切换 -->
          <NButton text class="shrink-0 text-xl!" @click.stop="toggleDone(todo)">
            <div
              :class="todo.completed
                ? 'i-lucide:circle-check-big text-green-500'
                : 'i-lucide:circle text-gray-300 hover:text-green-400'"
            />
          </NButton>

          <!-- 中间内容 -->
          <div class="flex-1 min-w-0">
            <!-- 编辑模式 -->
            <div v-if="isEditing(todo.id)" class="flex flex-col gap-3 py-2" @click.stop>
              <FloatInput v-model:value="draft.title" :placeholder="t('titlePlaceholder')" size="small" clearable />
              <FloatDatePicker
                v-model:value="draft.deadline" type="datetime" :placeholder="t('deadlinePlaceholder')"
                format="yyyy-MM-dd HH:mm" size="small"
              />
            </div>

            <!-- 展示模式 -->
            <div v-else class="flex flex-col gap-1 overflow-hidden">
              <div class="flex items-center gap-2">
                <NEllipsis class="text-sm min-w-0" :tooltip="{ width: 'trigger' }">
                  <NText
                    :depth="todo.completed ? 3 : 1" :type="todo.isOverdue ? 'error' : 'default'"
                    :strong="todo.isOverdue"
                  >
                    {{ todo.title || t('untitled') }}
                  </NText>
                </NEllipsis>
                <NTag v-if="todo.isOverdue" type="error" size="small" round :bordered="false">
                  {{ t('overdue') }}
                </NTag>
                <NTag v-else-if="todo.completed" type="success" size="small" round :bordered="false">
                  {{ t('done') }}
                </NTag>
              </div>
              <div class="flex items-center gap-1">
                <div class="i-lucide:clock text-xs text-gray-400" />
                <NText depth="3" class="text-xs">
                  {{ formatTime(todo.deadline) }}
                </NText>
              </div>
            </div>
          </div>

          <!-- 右侧操作区 -->
          <div class="shrink-0 flex items-center" @click.stop>
            <!-- 编辑态：保存 -->
            <NButton v-if="isEditing(todo.id)" text type="primary" class="text-lg!" @click="saveEdit">
              <div class="i-lucide:check" />
            </NButton>
            <!-- 展示态：删除（二次确认） -->
            <NPopconfirm v-else @positive-click="handleRemove(todo.id)">
              <template #trigger>
                <NButton text type="error" class="text-lg!">
                  <div class="i-lucide:trash-2" />
                </NButton>
              </template>
              {{ t('deleteConfirm') }}
            </NPopconfirm>
          </div>
        </div>
      </NListItem>
    </NList>
  </NScrollbar>
</template>

<style lang="css" scoped>
/* Naive UI 的 .n-list-item__main 只有 flex: 1 而没有 min-width: 0，
   长标题会撑开容器导致 NEllipsis 无法触发截断，这里补上宽度约束 */
:deep(.n-list-item__main) {
  min-width: 0;
}
</style>

<i18n lang="json">
{
  "zh": {
    "empty": "暂无待办事项",
    "titlePlaceholder": "待办事项",
    "deadlinePlaceholder": "截止时间",
    "untitled": "未命名事项",
    "overdue": "已超时",
    "done": "已完成",
    "deleteConfirm": "确定删除该待办吗？",
    "timeFormat": "M月D日 HH:mm"
  },
  "en": {
    "empty": "No todos yet",
    "titlePlaceholder": "Todo title",
    "deadlinePlaceholder": "Deadline",
    "untitled": "Untitled",
    "overdue": "Overdue",
    "done": "Done",
    "deleteConfirm": "Delete this todo?",
    "timeFormat": "MMM D, HH:mm"
  }
}
</i18n>
