const params = new URLSearchParams(window.location.search)
const messageId = params.get('message-id')
const mediaIds = params.get('media-id')

document.addEventListener('DOMContentLoaded', () => {
  document.querySelectorAll('input[data-type="tag"]').forEach((t) => {
    const tag = assertInstance(t, HTMLInputElement)
    tag.addEventListener('change', () => Tags.toggleTag(tag.name))
  })
  assertInstance(
    document.getElementById('callback'),
    HTMLButtonElement,
  ).addEventListener('click', () => {
    console.log(Tags.selected)
    return
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
