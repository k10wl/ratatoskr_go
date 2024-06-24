const TRANSITION_PROPERTY = '--transition-duration'

const params = new URLSearchParams(window.location.search)
const messageId = params.get('message-id')
const mediaIds = params.get('media-id')

if (!messageId) {
  throw new Error('messageId is required')
}
if (!mediaIds) {
  throw new Error('messageId is required')
}

document.addEventListener('DOMContentLoaded', () => {
  const persistence = new Persistence(messageId)
  const openedGroups = new StringSet(persistence.session.openedGroups)
  const selectedTags = new StringSet(persistence.session.selectedTags)

  const initialTransitionDuration =
    document.documentElement.style.getPropertyValue(TRANSITION_PROPERTY)
  document.documentElement.style.setProperty(TRANSITION_PROPERTY, '0ms')

  const throttledScrollUpdate = throttle(() => {
    persistence.update('scrollY', window.scrollY)
  }, 100)
  window.addEventListener('scroll', () => throttledScrollUpdate())
  window.addEventListener('beforeunload', () =>
    persistence.update('scrollY', window.scrollY),
  )

  document.querySelectorAll('input[data-type="group"]').forEach((g) => {
    const group = assertInstance(g, HTMLInputElement)
    if (persistence.session.openedGroups.includes(group.name)) {
      group.checked = true
    }
    group.addEventListener('change', () => {
      openedGroups.toggle(group.name)
      persistence.update('openedGroups', openedGroups.get())
    })
  })

  document.querySelectorAll('input[data-type="tag"]').forEach((t) => {
    const tag = assertInstance(t, HTMLInputElement)
    const ul = assertInstance(tag.closest('ul'), HTMLUListElement)
    const withGroup = `${ul.id}::${tag.name}`
    if (persistence.session.selectedTags.includes(withGroup)) {
      tag.checked = true
    }
    tag.addEventListener('change', () => {
      selectedTags.toggle(withGroup)
      persistence.update('selectedTags', selectedTags.get())
    })
  })

  window.scrollTo(0, persistence.session.scrollY)
  document.documentElement.style.setProperty(
    TRANSITION_PROPERTY,
    initialTransitionDuration,
  )

  assertInstance(
    document.getElementById('callback'),
    HTMLButtonElement,
  ).addEventListener('click', () => {
    const data = selectedTags.get().reduce((acc, cur) => {
      const [group, tag] = cur.split('::')
      if (acc[group]) {
        acc[group].push(tag)
      } else {
        acc[group] = [tag]
      }
      return acc
    }, /** @type {Record<string, string[]>} */ ({}))
    Telegram.WebApp.sendData(
      JSON.stringify({
        messageId,
        mediaIds,
        data,
      }),
    )
  })
})

class StringSet {
  /** @type Set<string> */
  #selected

  /** @param {string[]} initial */
  constructor(initial = []) {
    this.#selected = new Set(initial)
  }

  /**
   * @param {string} tag
   * @returns {boolean} true if now present
   */
  toggle(tag) {
    if (this.#selected.has(tag)) {
      this.#selected.delete(tag)
      return true
    }
    this.#selected.add(tag)
    return false
  }

  get() {
    return [...this.#selected]
  }
}

class Persistence {
  /**
   * @typedef persisted
   * @property {string} messageId
   * @property {string[]} openedGroups
   * @property {string[]} selectedTags
   * @property {number} scrollY
   */

  #key = 'ratatosrk-persistant-tags'
  /** @type persisted */
  session

  /** @param {string} messageId */
  constructor(messageId) {
    this.session = this.#getLastSession(messageId)
  }

  /**
   * @param {string} messageId
   * @returns {persisted}
   */
  #newSession(messageId) {
    return {
      messageId: messageId,
      openedGroups: [],
      selectedTags: [],
      scrollY: 0,
    }
  }

  /**
   * @param {string} messageId
   * @returns {persisted}
   */
  #getLastSession(messageId) {
    const stored = localStorage.getItem(this.#key)
    if (!stored) {
      return this.#newSession(messageId)
    }
    /** @type persisted */
    const parsed = JSON.parse(stored)
    if (
      typeof parsed !== 'object' ||
      parsed.messageId !== messageId ||
      !Array.isArray(parsed.selectedTags) ||
      parsed.selectedTags.some((el) => typeof el !== 'string')
    ) {
      return this.#newSession(messageId)
    }
    return parsed
  }

  /**
   * @template {keyof persisted} T
   * @param {T} key
   * @param {persisted[T]} value
   */
  update(key, value) {
    this.session[key] = value
    this.#write(this.session)
  }

  /** @param {persisted} data */
  #write(data) {
    localStorage.setItem(this.#key, JSON.stringify(data))
  }
}

/**
 * @template T
 * @returns {T}
 * @param {unknown} obj
 * @param {new (data: any) => T} type
 */
function assertInstance(obj, type) {
  if (obj instanceof type) {
    /** @type {any} */
    const any = obj
    /** @type {T} */
    const t = any
    return t
  }
  throw new Error(`Object ${obj} does not have the right type '${type}'!`)
}

/**
 * @template T
 * @param {(...args: T[]) => any} func
 * @param {number} delay
 */
function throttle(func, delay) {
  let lastCalled = 0
  /** @param {T[]} args */
  return function (...args) {
    const now = performance.now()
    if (now - lastCalled >= delay) {
      func(...args)
      lastCalled = now
    }
  }
}
