import * as utils from "./utils.js"

class AuthForm extends HTMLElement {
  connectedCallback() {
    this.innerHTML = `
<form id="auth">
  <label class="form_label">username</label>
  <input type="text" class="form_input" placeholder="casio123" id="username" required>
  <label class="form_label">password</label>
  <input type="text" class="form_input" placeholder="***" id="password" required>
  <input type="submit" id="form_submit" value="enter">
</form>
    `
  }
}

customElements.define("auth-form", AuthForm)
