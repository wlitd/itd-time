import type { Component } from 'vue'
import { RemindStyle } from '../../core/bindings'
import SquidwardRemind from './SquidwardRemind.vue'

/** 依据提醒样式解析出对应卡片组件；未知/零值/旧版 minimal 均回退章鱼哥。 */
export function resolveRemind(style: RemindStyle): Component {
  switch (style) {
    case RemindStyle.RemindStyleSquidward:
      return SquidwardRemind
    default:
      return SquidwardRemind
  }
}
