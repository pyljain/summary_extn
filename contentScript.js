// Listen for messages from the popup
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
    if (request.action === 'getPageContent') {
      // Extract the page content
      const pageContent = document.body.innerText || document.body.textContent;
      sendResponse({ content: pageContent });
    }
  });