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
  const body = JSON.stringify({
    deviceUUID: getDeviceUUID(),
    keyUUID: keyUuid,
    publicKey: btoa(
      JSON.stringify(await crypto.subtle.exportKey("jwk", keyPair.publicKey)),
    ),
  });

  const response = await fetch("/api/key", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: body,
  });

  if (!response.ok) {
    console.error("Failed to register key with server:", response);
  } else {
    console.log("Key registered with server:", response);
  }
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

  return btoa(String.fromCharCode(...new Uint8Array(cipherText)));
}

async function decrypt(privateKey, cipherText) {
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

async function sendMessage(message, subject) {
  //TODO: Implement
  console.log("Sending message to", subject);

  const ciphers = await Promise.all(
    await fetch(`/api/key?subject=${subject}`)
      .then((response) => response.json())
      .then(async (subjectKeys) => {
        console.log("Recieved:", subjectKeys);
        await subjectKeys.forEach(async (keyEntry) => {
          keyEntry.PublicKey = await crypto.subtle.importKey(
            "jwk",
            JSON.parse(atob(keyEntry.Encoded)),
            { name: "RSA-OAEP", hash: "SHA-256" },
            true,
            ["encrypt"],
          );
        });

        const myLatestKeyUUID = localStorage.getItem("latestKey");
        const myKey = {
          UUID: myLatestKeyUUID,
          PublicKey: (await getKey(myLatestKeyUUID)).keyPair.publicKey,
        };

        return [myKey, ...subjectKeys].map(async (key) => {
          console.debug("Encrypting with ", key.UUID, key.PublicKey);
          return {
            keyUUID: key.UUID,
            cipher: await encrypt(key.PublicKey, message),
          };
        });
      }),
  );

  console.log("Ciphers:", ciphers);
  const body = JSON.stringify({
    recipient: subject,
    ciphers: ciphers,
  });
  console.log("Sending body:", body);

  const response = await fetch(`/api/message`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: body,
  });

  if (!response.ok) {
    console.error("Failed to send message:", response);
  } else {
    console.log("Message sent:", response);
  }

  return;
}

async function renderExistingMessages(subject) {
  var me = subject == "anya" ? "zawie" : "anya";
  const input = `/api/message?subject=${subject}`;
  return await fetch(input)
    .then((response) => response.json())
    .then(async (body) => {
      console.log(`Recieved ${body.messages.length}`);
      for (const message of body.messages) {
        let msg = "🔒 ENCRYPTED";
        console.debug("Message:", message);
        for (const cipher of message.ciphers) {
          const key = await getKey(cipher.keyUUID);
          if (key === null) {
            console.debug("No private key found for", cipher.keyUUID);
            continue;
          }
          console.debug("Found key corresponding to", cipher.keyUUID);
          try {
            console.debug("Decrypting with", key.keyPair.privateKey);
            msg = await decrypt(key.keyPair.privateKey, cipher.cipher);
            break;
          } catch (e) {
            console.error(
              `Failed to decrypt message with ${cipher.keyUUID}:`,
              e,
            );
          }
        }
        displayMessage(msg, message.sender == me ? "sent" : "received");
      }
    });
}
/*
  Hooks
*/
function displayMessage(message, type) {
  document.getElementById("convo-div").insertAdjacentHTML(
    "afterend",
    `<div class="chat chat-${type == "sent" ? "end" : "start"}">
        <div class="chat-bubble">${message}</div>
        <div class="chat-footer opacity-50"></div>
      </div>`,
  );
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

function getSubject() {
  const queryString = window.location.search;
  const urlParams = new URLSearchParams(queryString);
  return urlParams.get("subject");
}

function setSubjectDisplay(subject) {
  document.getElementById("title").innerHTML = `${subject} Chat`;
  document.getElementById("subject-title").innerHTML = subject;
}

function onLoad() {
  const subject = getSubject();
  setSubjectDisplay(subject);

  console.log("Device UUID", getDeviceUUID());

  // Render existing messags
  renderExistingMessages(subject);

  // Listen for chat input
  const chatInput = document.getElementById("chat-input");
  chatInput.addEventListener("keypress", function (event) {
    if (event.key === "Enter") {
      const inputMessage = chatInput.value;
      chatInput.value = null;
      sendMessage(inputMessage, subject);
    }
  });
}
