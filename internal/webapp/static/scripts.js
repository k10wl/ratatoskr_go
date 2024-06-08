document.addEventListener('DOMContentLoaded', () => {
  const params = new URLSearchParams(window.location.search)
  params.forEach((v, k) => {
    assertInstance(document.getElementById('url'), HTMLElement).textContent +=
      k + ': ' + v + '\n'
  })
  assertInstance(document.getElementById('url'), HTMLElement).textContent +=
    '\ndata: ' + Telegram.WebApp.initData
})

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
