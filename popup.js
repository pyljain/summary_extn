console.log('Popup script loaded');
const summarizeButton = document.getElementById('summarizeButton');
const summaryDiv = document.getElementById('summary');

marked.setOptions({
    breaks: true, // Add line breaks on single linebreaks
    gfm: true,    // GitHub Flavored Markdown
});

document.getElementById('summarizeButton').addEventListener('click', () => {
    console.log('Summarize button clicked');
    // Update button to show loading state
    summarizeButton.disabled = true;
    summarizeButton.innerHTML = `
    <svg class="loading-spinner inline-block w-5 h-5 mr-2" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
    </svg>
    Summarizing...
    `;
    
    // Update summary div to show loading state
    summaryDiv.innerHTML = `
    <div class="flex flex-col items-center justify-center h-full">
        <svg class="loading-spinner w-8 h-8 mb-4 text-indigo-600" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="#4bcffa" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <p class="text-gray-600">Generating summary...</p>
    </div>
    `;
    // Get the active tab
    chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
      const activeTab = tabs[0];
  
      // Inject the content script into the active tab
      chrome.scripting.executeScript(
        {
          target: { tabId: activeTab.id },
          files: ['contentScript.js'],
        },
        () => {
          // After injecting, send a message to the content script
          chrome.tabs.sendMessage(
            activeTab.id,
            { action: 'getPageContent' },
            (response) => {
              if (chrome.runtime.lastError) {
                console.error(chrome.runtime.lastError.message);
                document.getElementById('summary').innerText =
                  'Error retrieving page content.';
                return;
              }
  
              const pageContent = response.content;
  
              // Send the content to the local service
              fetch('http://localhost:9001/summarize', {
                method: 'POST',
                headers: {
                  'Content-Type': 'application/json',
                },
                body: JSON.stringify({ content: pageContent }),
              })
                .then((response) => response.json())
                .then((data) => {
                  // Display the summary in the popup
                  summarizeButton.disabled = false;
                  summarizeButton.innerHTML = 'SUMMARIZE PAGE';
                  const markdownWrapper = document.createElement('div');
                  markdownWrapper.className = 'markdown-body';
                  
                  // Convert markdown to HTML and set it as the inner HTML of the wrapper
                  markdownWrapper.innerHTML = marked.parse(data.summary);
                  
                  // Clear the summary div and append the markdown wrapper
                  summaryDiv.innerHTML = '';
                  summaryDiv.appendChild(markdownWrapper);
                  // document.getElementById('summary').innerText = data.summary;
                })
                .catch((error) => {
                  button.disabled = false;
                  button.innerHTML = 'SUMMARIZE PAGE';
                  console.error('Error:', error);
                  document.getElementById('summary').innerText =
                    'Error fetching summary.';
                });
            }
          );
        }
      );
    });
  });