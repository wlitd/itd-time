// 集中导出 Wails 自动生成的绑定，供各 feature 以单一路径引用，
// 避免 binding 生成路径散落在业务代码里。
//
// 注意：beta.29 的 TS 绑定把 Go 结构体生成为纯 interface（类型），
// 而枚举/service 命名空间仍是值。isolatedModules 下类型再导出须用 export type。
export {
  SettingsService,
  SystemService,
  TodoService,
  UpdateService,
} from '../../bindings/itdtime/internal/api'

export { Coordinator } from '../../bindings/itdtime/internal/reminder'

// 枚举：值 + 类型，普通再导出。
export { RemindKind, RemindStyle } from '../../bindings/itdtime/internal/domain/models'

// Weekday 是 Go time.Weekday 的映射（Sun=0..Sat=6），落在独立的 time 绑定命名空间。
export { Weekday } from '../../bindings/time/models'

// 结构体：仅类型，须用 export type。
export type {
  AppData,
  OffWorkConfig,
  RemindSettings,
  ThemePrefs,
  Todo,
} from '../../bindings/itdtime/internal/domain/models'

export type { TodoInput, UpdateInfo } from '../../bindings/itdtime/internal/api/models'
