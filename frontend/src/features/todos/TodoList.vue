<script setup lang="ts">
import { computed, nextTick, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import DatePicker from 'primevue/datepicker'
import FloatLabel from 'primevue/floatlabel'
import { useAppStore } from '../../core/stores/app'
import type { DisplayTodo } from '../../core/stores/app'
import TodoItem from './TodoItem/TodoItem.vue'
import type { TodoSavePayload } from './TodoItem/types'

const { t } = useI18n()
// 跨组件共享键留在全局 messages；local 作用域读它们会触发回落告警，故显式用 global composer。
const { t: gt } = useI18n({ useScope: 'global' })
const store = useAppStore()

const todos = computed<DisplayTodo[]>(() => store.sortedTodos)
const pendingCount = computed<number>(() => todos.value.filter(todo => !todo.completed).length)

/** 当前处于行内编辑态的待办 id。 */
const editingId = ref<string | null>(null)

/** 新增态：列表顶部展开一条空编辑器。 */
const adding = ref<boolean>(false)
const newTitle = ref<string>('')
const newDeadline = ref<Date>(new Date())

async function startAdd(): Promise<void> {
  editingId.value = null
  newTitle.value = ''
  newDeadline.value = new Date(Date.now() + 60 * 60 * 1000)
  adding.value = true
  await nextTick()
  document.getElementById('todo-new-input')?.focus()
}

function cancelAdd(): void {
  adding.value = false
}

// 头部 + 按钮：未展开则进入新增，已展开则收起
function onAddClick(): void {
  if (adding.value)
    cancelAdd()
  else
    void startAdd()
}

async function saveNew(): Promise<void> {
  const title = newTitle.value.trim()
  if (!title || !newDeadline.value)
    return
  await store.addTodo(title, newDeadline.value.getTime())
  adding.value = false
}

function onEdit(todo: DisplayTodo): void {
  adding.value = false
  editingId.value = todo.id
}

function onCancelEdit(): void {
  editingId.value = null
}

async function onSave(payload: TodoSavePayload): Promise<void> {
  await store.updateTodo(payload.id, payload.title, payload.deadline)
  editingId.value = null
}

async function onToggle(id: string, completed: boolean): Promise<void> {
  await store.setCompleted(id, completed)
}

async function onRemove(id: string): Promise<void> {
  await store.removeTodo(id)
}
</script>

<template>
  <section class="surface-card flex min-h-0 flex-col overflow-hidden">
    <header class="flex items-center justify-between border-b border-[var(--app-border)] px-4 py-3">
      <div class="flex items-center gap-2">
        <span class="text-sm font-bold">{{ t('title') }}</span>
        <span
          v-if="pendingCount > 0"
          class="rounded-full bg-primary/12 px-2 py-0.5 text-xs font-semibold text-primary"
        >
          {{ pendingCount }}
        </span>
      </div>
      <Button size="small" rounded :label="t('add')" @click="onAddClick">
        <template #icon>
          <i class="i-lucide-plus text-base" />
        </template>
      </Button>
    </header>

    <div class="min-h-0 flex-1 overflow-y-auto px-3 py-3">
      <ul v-if="todos.length || adding" class="m-0 flex list-none flex-col gap-2 p-0">
        <!-- 行内新增编辑器：FloatLabel + 右侧勾保存，无取消/保存按钮 -->
        <li v-if="adding" class="w-full rounded-xl bg-[var(--app-surface-2)] p-2.5 ring-1 ring-primary/40">
          <div class="flex items-center gap-2">
            <div class="flex min-w-0 flex-1 flex-col gap-2">
              <FloatLabel variant="on">
                <InputText
                  id="todo-new-input"
                  v-model="newTitle"
                  fluid
                  @keyup.enter="saveNew"
                  @keyup.escape="cancelAdd"
                />
                <label for="todo-new-input">{{ gt('todos.titlePlaceholder') }}</label>
              </FloatLabel>
              <FloatLabel variant="on">
                <DatePicker
                  v-model="newDeadline"
                  input-id="todo-new-date"
                  fluid
                  show-time
                  hour-format="24"
                  :minute-step="5"
                  date-format="yy-mm-dd"
                  append-to="body"
                />
                <label for="todo-new-date">{{ gt('todos.deadlinePlaceholder') }}</label>
              </FloatLabel>
            </div>
            <Button
              severity="success"
              rounded
              class="shrink-0"
              :disabled="!newTitle.trim()"
              :aria-label="t('add')"
              @click="saveNew"
            >
              <template #icon>
                <i class="i-lucide-check text-lg" />
              </template>
            </Button>
          </div>
        </li>

        <TodoItem
          v-for="todo in todos"
          :key="todo.id"
          :todo="todo"
          :editing="editingId === todo.id"
          @toggle="onToggle"
          @edit="onEdit"
          @remove="onRemove"
          @save="onSave"
          @cancel="onCancelEdit"
        />
      </ul>

      <div v-else class="flex h-full flex-col items-center justify-center gap-3 py-10 text-[var(--app-muted)]">
        <span class="grid h-14 w-14 place-items-center rounded-2xl bg-[var(--app-surface-2)] ring-1 ring-[var(--app-border)]">
          <i class="i-lucide-inbox text-2xl" />
        </span>
        <span class="text-sm">{{ t('empty') }}</span>
      </div>
    </div>
  </section>
</template>

<i18n lang="yaml">
zh-CN:
  title: 待办
  add: 添加待办
  empty: 暂无待办事项
en-US:
  title: Todos
  add: Add Todo
  empty: No todos yet
</i18n>
