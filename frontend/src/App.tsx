import { useEffect, useState } from 'react'

type GoogleUser = {
	name: string
	email: string
	picture: string
}

function App() {
	const [user, setUser] = useState<GoogleUser | null>(null)

	useEffect(() => {
		const params = new URLSearchParams(window.location.search)
		const userParam = params.get('user')
		if (userParam) {
			try {
				const parsed = JSON.parse(decodeURIComponent(userParam))
				setUser(parsed)
			} catch (e) {
				console.error('Ошибка парсинга:', e)
			}
		}
	}, [])

	const handleLogin = () => {
		window.location.href = 'http://localhost:8080/login'
	}

	const handleLogout = () => {
		setUser(null)
		window.history.replaceState({}, document.title, '/')
	}

	return (
		<div style={{ padding: '2rem', fontFamily: 'Arial' }}>
			<h1>🔐 Google Login</h1>
			{user ? (
				<div>
					<img src={user.picture} width={80} style={{ borderRadius: '50%' }} />
					<h2>{user.name}</h2>
					<p>{user.email}</p>
					<button onClick={handleLogout}>Выйти</button>
				</div>
			) : (
				<button onClick={handleLogin}>Войти через Google</button>
			)}
		</div>
	)
}

export default App
