const params = new URLSearchParams(window.location.search)
const messageId = params.get('message-id')
const mediaIds = params.get('media-id')

if (!messageId) {
  throw new Error('Message is required')
}

document.addEventListener('DOMContentLoaded', () => {
  const persistence = new Persistence(messageId)

  persistence.lastSession.forEach((el) => Tags.toggleTag(el))

  document.querySelectorAll('input[data-type="tag"]').forEach((t) => {
    const tag = assertInstance(t, HTMLInputElement)
    if (persistence.lastSession.includes(tag.name)) {
      tag.checked = true
    }
    tag.addEventListener('change', () => {
      Tags.toggleTag(tag.name)
      persistence.storeTags(Tags.selected)
    })
  })

  assertInstance(
    document.getElementById('callback'),
    HTMLButtonElement,
  ).addEventListener('click', () => {
    Telegram.WebApp.sendData(
      JSON.stringify({
        messageId,
        mediaIds,
        tags: Tags.selected,
      }),
    )
  })
})

class Tags {
  /** @type string[] */
  static selected = []

  static getSelected() {
    return Tags.selected
  }

  /**
   * @param {string} tag
   * @returns {boolean} true if now present
   */
  static toggleTag(tag) {
    const i = this.selected.indexOf(tag)
    if (i === -1) {
      this.selected.push(tag)
      return true
    }
    this.selected.splice(i, 1)
    return false
  }
}

class Persistence {
  /**
   * @typedef persisted
   * @property {string} messageId
   * @property {string[]} tags
   */

  #key = 'ratatosrk-persistant-tags'
  #messageId
  /** @type string[] */
  lastSession

  /** @param {string} messageId */
  constructor(messageId) {
    this.#messageId = messageId
    this.lastSession = this.#getLastSession()
  }

  /** @returns {string[]} */
  #getLastSession() {
    const stored = localStorage.getItem(this.#key)
    if (!stored) {
      return []
    }
    /** @type persisted */
    const parsed = JSON.parse(stored)
    if (
      typeof parsed !== 'object' ||
      parsed.messageId !== this.#messageId ||
      !Array.isArray(parsed.tags) ||
      parsed.tags.some((el) => typeof el !== 'string')
    ) {
      return []
    }
    return parsed.tags
  }

  /** @param {string[]} tags */
  storeTags(tags) {
    localStorage.setItem(
      this.#key,
      JSON.stringify({ messageId: this.#messageId, tags }),
    )
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
