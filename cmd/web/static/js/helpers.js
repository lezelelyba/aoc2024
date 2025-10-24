/**
 * Fills template string from session storage
 * 
 * replace {key} placeholders within template with value found under "key" in session storage
 * @param {string} template - template
 * @returns {string} Filled template
 * @throws {Error} If some key is missing from session storage
 */
function fillTemplateFromSession(template) {
  return template.replace(/{(\w+)}/g, (match, key) => {
    const value = sessionStorage.getItem(key);
    if (value === null) {
      throw new Error(`Missing session value for: ${key}`);
    }
    return value;
  });
}