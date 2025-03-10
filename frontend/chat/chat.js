/*
  Local key storage and retrieval functions
*/
async function storeKey(uuid, keyPair) {
  localStorage.setItem(
    uuid,
    JSON.stringify({
      keyPair: {
        privateKey: await crypto.subtle.exportKey("jwk", keyPair.privateKey),
        publicKey: await crypto.subtle.exportKey("jwk", keyPair.publicKey),
      },
      createdAt: new Date().toISOString(),
    }),
  );
  localStorage.setItem("latestKey", uuid);
}

async function getLatestKey() {
  const uuid = localStorage.getItem("latestKey");
  if (uuid === null) {
    return null;
  }
  return getKey(uuid);
}

async function getKey(uuid) {
  // Retrieving key info from local storage
  const content = localStorage.getItem(uuid);
  if (content === null) {
    return null;
  }

  // Parse key information and import keys
  parsedContent = JSON.parse(content);
  parsedContent.keyPair.privateKey = await crypto.subtle.importKey(
    "jwk",
    parsedContent.keyPair.privateKey,
    { name: "RSA-OAEP", hash: "SHA-256" },
    true,
    ["decrypt"],
  );
  parsedContent.keyPair.publicKey = await crypto.subtle.importKey(
    "jwk",
    parsedContent.keyPair.publicKey,
    { name: "RSA-OAEP", hash: "SHA-256" },
    true,
    ["encrypt"],
  );

  return parsedContent;
}

function getDeviceUUID() {
  const deviceUUID = localStorage.getItem("deviceUUID");
  if (deviceUUID === null) {
    const newDeviceUUID = generateDeviceUUID();
    localStorage.setItem("deviceUUID", newDeviceUUID);
    return newDeviceUUID;
  }
  return deviceUUID;
}

/*
  Server interaction functions
*/
async function registerKey(keyUuid, keyPair) {
  //TODO: Implement
  console.log(
    "Registering key with server:",
    getDeviceUUID(),
    keyUuid,
    keyPair,
  );
  return;
}

/*
  Crypto functions
*/
async function generateKeyPair() {
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

function generateDeviceUUID() {
  return crypto.randomUUID();
}

async function encrypt(publicKey, plainText) {
  const encodedData = new TextEncoder().encode(plainText);
  const cipherText = await window.crypto.subtle.encrypt(
    {
      name: "RSA-OAEP",
    },
    publicKey,
    encodedData,
  );
  return cipherText;
}

async function decrypt(privateKey, cipherText) {
  const plainText = await window.crypto.subtle.decrypt(
    {
      name: "RSA-OAEP",
    },
    privateKey,
    cipherText,
  );
  return new TextDecoder().decode(plainText);
}

/*
  Other helper functions
*/
async function refreshKeyPair() {
  try {
    const [uuid, keyPair] = await generateKeyPair();

    await registerKey(uuid, keyPair);
    await storeKey(uuid, keyPair);

    return [uuid, keyPair];
  } catch (error) {
    console.error("Error generating key pair:", error);
    return null;
  }
}

async function sendMessage(message) {
  //TODO: Implement
  console.log('Sending message', message )
  return;
}

/*
  Hooks
*/
function displayMessage(message) {
  console.log("Message received:", message);
  document.getElementById('convo-div')
    .insertAdjacentHTML('afterend', `<div class="chat chat-start">
      <div class="chat-bubble">${message}</div>
        <div class="chat-footer opacity-50">Delivered</div>
        </div>`)
}

async function refreshIfNecessary() {
  const key = await getLatestKey();
  if (key === null) {
    console.debug("No key found, generating new key pair");
    await refreshKeyPair();
    return;
  }
  if (key.createdAt < new Date(Date.now() - 10000).toISOString()) {
    console.debug(
      "Most recent stored key is older than 10 second. Refreshing.",
    );
    await refreshKeyPair();
    return;
  }
  console.debug("Most recent stored key is still valid. Not refreshing.");
}

// Automatically key if key pair needs refresh on page load
refreshIfNecessary();

getLatestKey()
  .then(async key => {
    console.log("Using key", key);
    await new Promise(r => setTimeout(r, 2000));
    
    encrypt(key.keyPair.publicKey, "Hello, World!")
      .then(
        async (messageRecieved) =>
          await decrypt(key.keyPair.privateKey, messageRecieved),
      )
      .then(displayMessage);
    return key;
  })


function onLoad() {
  
  // Listen for chat input
  const chatInput = document.getElementById('chat-input')
  chatInput.addEventListener('keypress', function(event) {
    if (event.key === 'Enter') {
      const inputMessage = chatInput.value
      chatInput.value = null
      sendMessage(inputMessage);
    }
  });
}
