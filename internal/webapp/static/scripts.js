const params = new URLSearchParams(window.location.search)
const messageId = params.get('message-id')
const mediaIds = params.get('media-id')

document.addEventListener('DOMContentLoaded', () => {
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

class Templates {
  /** @type HTMLTemplateElement | null */
  static #tagTemplate = null

  /**
   * @param {string} tag
   * @returns {DocumentFragment}
   */
  static createTag(tag) {
    const clone = assertInstance(
      Templates.#getTagTemplateContent(),
      DocumentFragment,
    )
    assertInstance(clone.querySelector('.text'), HTMLElement).textContent = tag
    const button = assertInstance(
      clone.querySelector('button'),
      HTMLButtonElement,
    )
    button.addEventListener('click', () => {
      if (Tags.toggleTag(tag)) {
        button.classList.remove('secondary-bg')
      } else {
        button.classList.add('secondary-bg')
      }
    })
    return clone
  }

  /**
   * @returns {Node}
   */
  static #getTagTemplateContent() {
    if (Templates.#tagTemplate === null) {
      Templates.#tagTemplate = assertInstance(
        document.querySelector('template#tag'),
        HTMLTemplateElement,
      )
    }
    return Templates.#tagTemplate.content.cloneNode(true)
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
