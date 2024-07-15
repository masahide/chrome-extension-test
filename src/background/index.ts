let settings = {
  browser: "",
  profile: "",
};

function loadSettings() {
  chrome.storage.sync.get(["browser", "profile"], (data) => {
    settings.browser = data.browser;
    settings.profile = data.profile;
    updateContextMenus();
  });
}

function updateContextMenus() {
  chrome.contextMenus.update("openAltBrowser", {
    title: `他のブラウで開く(${settings.browser})`,
  });
}
chrome.runtime.onInstalled.addListener(() => {
  loadSettings();
  chrome.contextMenus.create({
    id: "openAltBrowser",
    title: `他のブラウで開く(${settings.browser})`,
  });
});
chrome.action.onClicked.addListener(() => {
  console.log("onCkeck:", settings);
  chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
    let activeTab = tabs[0];
    let url = activeTab.url;
    sendNativeMessage({
      browser: settings.browser,
      profile: settings.profile,
      url: url,
    });
  });
});
chrome.contextMenus.onClicked.addListener((info, tab) => {
  if (info.menuItemId === "openAltBrowser" && tab) {
    sendNativeMessage({
      browser: settings.browser,
      profile: settings.profile,
      url: tab.url,
    });
  }
});
chrome.storage.onChanged.addListener((changes, areaName) => {
  if (areaName === "sync") {
    loadSettings();
  }
});

chrome.runtime.onMessage.addListener((request) => {
  if (request.name === "testOpen") {
    sendNativeMessage({
      browser: request.browser,
      profile: request.profile,
      url: "https://google.com",
    });
  }
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
