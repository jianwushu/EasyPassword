document.addEventListener('submit', (event) => {
  const form = event.target as HTMLFormElement;
  const passwordInput = form.querySelector('input[type="password"]') as HTMLInputElement;

  if (passwordInput) {
    const usernameInput = form.querySelector('input[type="email"], input[type="text"], input[name="username"], input[name="user"]') as HTMLInputElement;

    if (usernameInput && passwordInput.value) {
      const credentials = {
        account: usernameInput.value,
        password: passwordInput.value,
        website: window.location.href,
        name: document.title,
      };

      console.log('Detected form submission with credentials:', credentials);
      chrome.runtime.sendMessage({
        type: 'SAVE_PASSWORD_REQUEST',
        data: credentials,
      });
    }
  }
});