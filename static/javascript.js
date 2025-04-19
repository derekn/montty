let ws,
	contentEl,
	connectTimer,
	connectCount = 0

function setConnState(...states) {
	document.body.classList.remove('disconnected', 'connecting')
	document.body.classList.add(...states)
}

function connect() {
	const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
	setConnState('connecting')
	ws = new WebSocket(`${protocol}//${location.host}/ws`)

	ws.onopen = () => {
		setConnState()
		contentEl.innerHTML = ''
		connectCount = 0
	}
	ws.onclose = () => {
		setConnState('disconnected')
		if (connectCount++ < 10) {
			clearTimeout(connectTimer)
			connectTimer = setTimeout(connect, 2000)
		} else {
			console.error(`max connection attempts (${connectCount - 1}) exceeded, stopping`)
		}
	}
	ws.onmessage = (evt) => {
		contentEl.innerHTML += evt.data
		window.scrollTo({ top: document.body.scrollHeight, behavior: 'smooth' })
	}
	ws.onerror = (err) => {
		// console.error('websocket error:', err)
	}
}

document.addEventListener('DOMContentLoaded', () => {
	contentEl = document.getElementById('content')
	connect()
})
