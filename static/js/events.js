const trackedHosts = ['x.com', 'github.com', 'buymeacoffee.com'];

window.addEventListener('load', () => {
	document.querySelectorAll('.social-link').forEach(link => {
		link.addEventListener('click', () => {
			const url = new URL(link.getAttribute('href'));
			if (!trackedHosts.includes(url.host)) return;
			fathom.trackEvent(`social click: ${url.host}`);
		});
	});
});
