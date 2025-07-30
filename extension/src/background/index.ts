chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.type === 'SAVE_PASSWORD_REQUEST') {
    // 检查用户是否已登录
    chrome.storage.local.get('auth', (result) => {
      console.log('Received SAVE_PASSWORD_REQUEST:', message.data);
      console.log('Current auth state:', result.auth);
      const authState = result.auth ? JSON.parse(result.auth) : null;
      if (authState && authState.isAuthenticated) {
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
  // 保留其他消息监听器（如果它们仍然需要）
  // 注意：LOGIN_SUCCESS, LOGOUT, GET_TOKEN 现在由 pinia 状态管理处理，
  // 这些消息监听器可能不再需要，除非有其他地方使用它们。
  else if (message.type === 'LOGIN_SUCCESS') {
    // 这个分支可能可以被移除
    chrome.storage.local.set({ token: message.token });
  } else if (message.type === 'LOGOUT') {
    // 这个分支可能可以被移除
    chrome.storage.local.remove('token');
  } else if (message.type === 'GET_TOKEN') {
    // 这个分支可能可以被移除
    chrome.storage.local.get('token', (result) => {
      sendResponse(result.token);
    });
    return true;
  }
});