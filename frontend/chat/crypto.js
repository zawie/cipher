/**
 * @typedef {string} KeyId
 * @returns {Promise<[KeyId, CryptoKeyPair]>}
 */
export async function generateKeyPair() {
  const uuid = crypto.randomUUID();
  const keyPair = await window.crypto.subtle.generateKey(
    {
      name: "RSA-OAEP",
      modulusLength: 2048,
      publicExponent: new Uint8Array([0x01, 0x00, 0x01]), // Equivalent to 65537
      hash: "SHA-256",
    },
    true, // extractable
    ["encrypt", "decrypt"], // key usages
  );

  return [uuid, keyPair];
}

/**
 * @typedef {string} DeviceId
 * @returns {DeviceId}- A unique identifier for a device
 */
export function generateDeviceUUID() {
  return crypto.randomUUID();
}

/**
 * @param {CryptoKey} The public key to encrypt with
 * @param {string} The plain text to encrypt
 * @returns {Promise<string>} The encrypted cipher text
 */
export async function encrypt(publicKey, plainText) {
  const encodedData = new TextEncoder().encode(plainText);
  const cipherText = await window.crypto.subtle.encrypt(
    {
      name: "RSA-OAEP",
    },
    publicKey,
    encodedData,
  );

  return btoa(String.fromCharCode(...new Uint8Array(cipherText)));
}

/**
 * @param {CryptoKey} The provate key to decrypt with
 * @param {string} The cipher text to decrypt
 * @returns {Promise<string>} The decrypted plain text
 */
export async function decrypt(privateKey, cipherText) {
  cipherText = new Uint8Array(
    atob(cipherText)
      .split("")
      .map((c) => c.charCodeAt(0)),
  );
  const plainText = await window.crypto.subtle.decrypt(
    {
      name: "RSA-OAEP",
    },
    privateKey,
    cipherText,
  );
  return new TextDecoder().decode(plainText);
}
