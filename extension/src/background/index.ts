chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.type === 'LOGIN_SUCCESS') {
    chrome.storage.local.set({ token: message.token });
  } else if (message.type === 'LOGOUT') {
    chrome.storage.local.remove('token');
  } else if (message.type === 'GET_TOKEN') {
    chrome.storage.local.get('token', (result) => {
      sendResponse(result.token);
    });
    return true;
  } else if (message.type === 'SAVE_PASSWORD_REQUEST') {
    // 检查用户是否已登录
    chrome.storage.local.get('token', (result) => {
      if (result.token) {
        // 将凭据和操作暂存，以便弹窗可以访问
        chrome.storage.local.set(
          {
            pendingCredential: message.data,
            redirectTo: '/vault/add',
          },
          () => {
            // 打开插件的弹出窗口
            chrome.action.openPopup();
          },
        );
      }
    });
  }
});