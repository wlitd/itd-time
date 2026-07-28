import type { Directive } from 'vue'

// 水波纹指令
const ripple: Directive<HTMLElement, boolean | undefined> = {
  mounted(el, binding) {
    el.style.position = 'relative'

    const handleClick = (e: MouseEvent) => {
      if (binding.value === false)
        return

      const rect = el.getBoundingClientRect()
      const x = e.clientX - rect.left
      const y = e.clientY - rect.top

      const radius = Math.sqrt(
        Math.max(
          (x * x + y * y),
          ((rect.width - x) * (rect.width - x) + y * y),
          (x * x + (rect.height - y) * (rect.height - y)),
          ((rect.width - x) * (rect.width - x) + (rect.height - y) * (rect.height - y))
        )
      )

      // 创建裁剪容器
      const clipContainer = document.createElement('div')
      clipContainer.style.position = 'absolute'
      clipContainer.style.inset = '0'
      clipContainer.style.overflow = 'hidden'
      clipContainer.style.borderRadius = getComputedStyle(el).borderRadius
      clipContainer.style.pointerEvents = 'none'

      const ripple = document.createElement('span')
      ripple.style.position = 'absolute'
      ripple.style.borderRadius = '50%'
      ripple.style.left = `${x - radius}px`
      ripple.style.top = `${y - radius}px`
      ripple.style.width = `${radius * 2}px`
      ripple.style.height = `${radius * 2}px`
      ripple.style.background = 'var(--itd-ripple-bg)'
      ripple.style.opacity = '0.4'
      ripple.style.transform = 'scale(0)'
      ripple.style.transition = 'transform 0.6s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.6s ease-out'

      clipContainer.appendChild(ripple)
      el.appendChild(clipContainer)

      // 触发动画
      requestAnimationFrame(() => {
        ripple.style.transform = 'scale(1)'
        ripple.style.opacity = '0'
      })

      setTimeout(() => {
        clipContainer.remove()
      }, 600)
    }

    el.addEventListener('click', handleClick);
    (el as any).__ripple_click_handler__ = handleClick
  },
  unmounted(el: HTMLElement) {
    const handler = (el as any).__ripple_click_handler__
    if (handler) {
      el.removeEventListener('click', handler)
    }
  }
}

export default ripple
