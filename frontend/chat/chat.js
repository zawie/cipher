import { displayMessage, createBubble, updateBubble } from "./bubbles.js";
import {
  encrypt,
  decrypt,
  generateDeviceUUID,
  generateKeyPair,
} from "./crypto.js";
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
  let parsedContent = JSON.parse(content);
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
  console.log("Sending message to", subject);

  const msgGuiId = createBubble("sent");

  const ciphers = await Promise.all(
    await fetch(`/api/key?subject=${subject}`)
      .then((response) => {
        if (!response.ok) {
          console.error("Failed to fetch subject keys:", response);
          updateBubble(msgGuiId, message, "❌ Failed to send message");
          rai;
        }
        return response;
      })
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
    updateBubble(msgGuiId, message, "❌ Failed to send message");
  } else {
    console.log("Message sent:", response);
    updateBubble(msgGuiId, message, "Sent");
  }

  return;
}

function scrollToBottom() {
  // const convoDiv = document.getElementById("convo-div");
  // convoDiv.scrollTop = convoDiv.scrollHeight;
}

window.onload = scrollToBottom;

async function renderExistingMessages(subject) {
  var me = subject == "alice" ? "bob" : "alice";
  const input = `/api/message?subject=${subject}`;
  return await fetch(input)
    .then((response) => response.json())
    .then(async (body) => {
      if (body.messages == null) {
        body.messages = [];
      }
      console.log(`Recieved ${body.messages.length} messages:\n${body}`);
      for (let i = 0; i < body.messages.length; i++) {
        const message = body.messages[i];
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

onLoad();
