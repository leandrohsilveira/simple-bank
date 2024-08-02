/**
 *
 * @param {unknown} value
 * @param {string} message
 * @returns {asserts value}
 */
export function assert(value, message) {
  if (!value) throw new Error(message);
}
