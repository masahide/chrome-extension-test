chrome.action.onClicked.addListener(() => {
  chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
    let activeTab = tabs[0];
    let url = activeTab.url;
    sendNativeMessage({ url: url });
  });
});

function sendNativeMessage(message: any) {
  chrome.runtime.sendNativeMessage(
    "com.my_company.my_application",
    message,
    (response) => {
      if (chrome.runtime.lastError) {
        console.error("ERROR:", chrome.runtime.lastError.message);
        return;
      }
      console.log("Received response:", response);
    },
  );
}

/*
chrome.sidePanel
  .setPanelBehavior({ openPanelOnActionClick: true })
  .catch((error) => console.error(error));

function setupContextMenu() {
chrome.contextMenus.create({
  id: "summary-text",
  title: "Summary",
  contexts: ["selection"],
});
}

chrome.runtime.onInstalled.addListener(() => {
setupContextMenu();
});
*/

/*
//アプリから切断されたときの処理
port.onDisconnect.addListener(() => {
  if (chrome.runtime.lastError) {
    console.log("onDisconnect err: " + chrome.runtime.lastError.message);
  }
  console.log("onDisconnect: Disconnected!");
});
*/

/*
const port = chrome.runtime.connectNative("com.my_company.my_application");

//ローカルアプリからメッセージ受信
port.onMessage.addListener((req) => {
  console.log("req : " + JSON.stringify(req));
  if (chrome.runtime.lastError) {
    console.log(
      "onMessage.addListener error: " + chrome.runtime.lastError.message,
    );
  }
});
//ローカルアプリへメッセージ送信
// port.postMessage({ message: "ping", body: "hello from browser extension" });
port.postMessage({ text: "Hello, my_application" });

// 接続を切断し、native messaging host を終了する
setTimeout(() => {
  port.disconnect();
  console.log("Disconnected");
}, 3000);
*/

/*
chrome.contextMenus.onClicked.addListener((data) => {
  chrome.runtime.sendMessage({
    name: "summary-text",
    data: { value: data.selectionText },
  });
});

chrome.runtime.onInstalled.addListener(async () => {
  let manifest = chrome.runtime.getManifest();
  if (manifest && manifest.content_scripts) {
    for (const cs of manifest.content_scripts) {
      for (const tab of await chrome.tabs.query({ url: cs.matches })) {
        if (tab && cs.js && tab.id) {
          chrome.scripting.executeScript({
            files: cs.js,
            target: { tabId: tab.id, allFrames: cs.all_frames },
            injectImmediately: cs.run_at === "document_start",
            // world: cs.world, // uncomment if you use it in manifest.json in Chrome 111+
          });
        }
      }
    }
  }
});
*/
