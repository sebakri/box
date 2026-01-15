
window.addEventListener('DOMContentLoaded', () => {
    const codeBlocks = document.querySelectorAll('pre > code');

    codeBlocks.forEach(codeBlock => {
        const copyButton = document.createElement('button');
        copyButton.className = 'copy-code-button';
        copyButton.textContent = 'Copy';

        const preElement = codeBlock.parentNode;
        if (preElement) {
            preElement.style.position = 'relative';
            preElement.appendChild(copyButton);

            copyButton.addEventListener('click', () => {
                const codeToCopy = codeBlock.textContent;
                navigator.clipboard.writeText(codeToCopy).then(() => {
                    copyButton.textContent = 'Copied!';
                    setTimeout(() => {
                        copyButton.textContent = 'Copy';
                    }, 2000);
                }).catch(err => {
                    console.error('Failed to copy text: ', err);
                });
            });
        }
    });
});
