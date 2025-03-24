const REGISTER_ENDPOINT = `/api/auth/register`;
const LOGIN_ENDPOINT = `/api/auth/login`;

/**
 * An object representing authentication information
 * @typedef {Object} AuthObject
 * @property {string} alias - The signin alias
 * @property {string} privatePassword - The password to use on the client side
 * @property {string} publicPassword - The password to submit to the server
 */

/**
 * Handles sign-up form submission
 *
 * @param {Event} event - The submit event from the form.
 * @returns {boolean} Always returns false to prevent the form from submitting.
 */
function onSignupFormSubmit(event) {
  passAuthToEndpointFromFormEvent(event, REGISTER_ENDPOINT);
  return false; // Return false to prevent form from submitting
}

/**
 * Handles sign-up form submission
 *
 * @param {Event} event - The submit event from the form.
 * @returns {boolean} Always returns false to prevent the form from submitting.
 */
function onLoginFormSubmit(event) {
  passAuthToEndpointFromFormEvent(event, LOGIN_ENDPOINT);
  return false; // Return false to prevent form from submitting
}

/**
 * @param {Event} event - The submit event from the form.
 * @param {string} endpoint - Endpoint to pass auth informaiton to
 * @returns {boolean} Always returns false to prevent the form from submitting.
 */
function passAuthToEndpointFromFormEvent(event, endpoint) {
  event.preventDefault();
  disableButton();

  const auth = new AuthObject(event);
  auth.basicAuthHeader().then(async (authHeader) => {
    const response = await fetch(endpoint, {
      method: "POST",
      headers: {
        Authorization: authHeader,
      },
    });

    console.log(response);
    if (!response.ok) {
      alert(`Error: ${response.statusText}`);
      enableButton();
      return;
    }

    // Redirect to main page on sign-in success
    window.location.href = "/";
  });
}

/*
 * Class representing an authentication (alias + password)
 */
class AuthObject {
  /* Store input */
  #alias;
  #private;
  #public;

  /** Stores resolve functions **/
  #resolvePrivate;
  #resolvePublic;

  /**
   * Turns a form submission event into an auth object
   *
   * @param {Event} event - The submit event from the form.
   */
  constructor(event) {
    // Extract data from the form event
    const form = event.target;
    const formData = new FormData(form);

    this.#alias = formData.get("alias");
    const password = formData.get("password");

    // Promises for signatures
    this.#private = new Promise((resolve) => (this.#resolvePrivate = resolve));
    this.#public = new Promise((resolve) => (this.#resolvePublic = resolve));

    this.#computeSignatures(password);
  }

  /**
   * Computes private and public signatures
   * @private
   */
  async #computeSignatures(password) {
    const encoder = new TextEncoder();
    const encodedPassword = encoder.encode(password);

    const hashBuffer = await window.crypto.subtle.digest(
      "SHA-256",
      encodedPassword,
    );

    const hashArray = new Uint8Array(hashBuffer);
    const midPoint = Math.floor(hashArray.length / 2);

    // Resolve promises
    this.#resolvePrivate(Uint8ArrayToString(hashArray.subarray(0, midPoint)));
    this.#resolvePublic(Uint8ArrayToString(hashArray.subarray(midPoint)));
  }

  /**
   * BasicAuthHeader
   */
  async basicAuthHeader() {
    return `Basic ${btoa(`${this.#alias}:${await this.#public}`)}`;
  }
}

/**
 * Converts a Uint8Array to a string
 * @param {Uint8Array} uint8Array - The byte array to convert
 * @returns {string} - The string representation
 */
function Uint8ArrayToString(uint8Array) {
  return btoa(String.fromCharCode(...uint8Array)); // Base64 encode for readability
}

function disableButton(buttonId = "submit-button") {
  document.getElementById(buttonId).classList.add("btn-disabled");
}

function enableButton(buttonId = "submit-button") {
  document.getElementById(buttonId).classList.remove("btn-disabled");
}
